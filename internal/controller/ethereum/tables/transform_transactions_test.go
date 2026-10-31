package tables

import (
	"testing"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/stretchr/testify/require"

	"github.com/coinbase/chainstorage/protos/coinbase/c3/common"
	chainstorageapi "github.com/coinbase/chainstorage/protos/coinbase/chainstorage"

	"github.com/coinbase/chainsformer/internal/config"
)

const tronTestMemo = "6be029c1fa65e2b3a98e23b325b5d4da1c1ee32b39c7ddf2a20ed685776e280b"

func newChainConfig(blockchain common.Blockchain) *config.Config {
	return &config.Config{Chain: config.ChainConfig{Blockchain: blockchain}}
}

// newTronBlockWithMemo returns a block whose first transaction carries a memo and whose second doesn't.
func newTronBlockWithMemo() *chainstorageapi.EthereumBlock {
	return &chainstorageapi.EthereumBlock{
		Header: &chainstorageapi.EthereumHeader{Hash: "h", Number: 87573393},
		Transactions: []*chainstorageapi.EthereumTransaction{
			{Hash: "4b1345bbc3610321d041619b80e52f16ddfde9108b814ac5af0b7c1ea9240446", Index: 0, Type: 1, Memo: tronTestMemo, Receipt: &chainstorageapi.EthereumTransactionReceipt{}},
			{Hash: "61cba69bd9b6fb70692aca4f7cdfc45ccc95822da2c00e9b88c5101c6a603ed3", Index: 1, Type: 1, Receipt: &chainstorageapi.EthereumTransactionReceipt{}},
		},
	}
}

func fieldNames(schema *arrow.Schema) []string {
	names := make([]string, 0, len(schema.Fields()))
	for _, field := range schema.Fields() {
		names = append(names, field.Name)
	}
	return names
}

func TestTransactionSchema_Memo(t *testing.T) {
	tron := newChainConfig(common.Blockchain_BLOCKCHAIN_TRON)
	for _, schema := range []*arrow.Schema{newTransactionSchema(tron), newStreamedTransactionSchema(tron)} {
		names := fieldNames(schema)
		require.Equal(t, []string{"traces", "memo", "_partition_by", "_repartition_by_range"}, names[len(names)-4:])
	}

	ethereum := newChainConfig(common.Blockchain_BLOCKCHAIN_ETHEREUM)
	for _, schema := range []*arrow.Schema{newTransactionSchema(ethereum), newStreamedTransactionSchema(ethereum)} {
		require.False(t, schema.HasField("memo"))
		names := fieldNames(schema)
		require.Equal(t, []string{"traces", "_partition_by", "_repartition_by_range"}, names[len(names)-3:])
	}
}

func TestTransformTransactions_Memo(t *testing.T) {
	type transformFn func(cfg *config.Config, rb *array.RecordBuilder, block *chainstorageapi.EthereumBlock) error
	batch := func(cfg *config.Config, rb *array.RecordBuilder, block *chainstorageapi.EthereumBlock) error {
		return transactionsTable{config: cfg}.transformTransactions(rb, block, 0)
	}
	stream := func(cfg *config.Config, rb *array.RecordBuilder, block *chainstorageapi.EthereumBlock) error {
		event := &chainstorageapi.BlockchainEvent{SequenceNum: 7, Type: chainstorageapi.BlockchainEvent_BLOCK_ADDED}
		return nativeStreamedTransactionsTable{config: cfg}.transformStreamedTransactions(rb, block, event, 0)
	}
	cases := []struct {
		name            string
		blockchain      common.Blockchain
		newSchema       func(cfg *config.Config) *arrow.Schema
		transform       transformFn
		wantMemo        bool
		wantRepartition uint64
	}{
		{"tron/transactionsTable", common.Blockchain_BLOCKCHAIN_TRON, newTransactionSchema, batch, true, 87573393},
		{"tron/nativeStreamedTransactionsTable", common.Blockchain_BLOCKCHAIN_TRON, newStreamedTransactionSchema, stream, true, 7},
		{"ethereum/transactionsTable", common.Blockchain_BLOCKCHAIN_ETHEREUM, newTransactionSchema, batch, false, 87573393},
		{"ethereum/nativeStreamedTransactionsTable", common.Blockchain_BLOCKCHAIN_ETHEREUM, newStreamedTransactionSchema, stream, false, 7},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cfg := newChainConfig(tc.blockchain)
			rb := array.NewRecordBuilder(memory.DefaultAllocator, tc.newSchema(cfg))
			defer rb.Release()

			require.NoError(t, tc.transform(cfg, rb, newTronBlockWithMemo()))
			record := rb.NewRecord()
			defer record.Release()
			require.Equal(t, int64(2), record.NumRows())

			hashes := record.Column(record.Schema().FieldIndices("transaction_hash")[0]).(*array.String)
			require.Equal(t, "4b1345bbc3610321d041619b80e52f16ddfde9108b814ac5af0b7c1ea9240446", hashes.Value(0))
			require.Equal(t, "61cba69bd9b6fb70692aca4f7cdfc45ccc95822da2c00e9b88c5101c6a603ed3", hashes.Value(1))

			repartition := record.Column(record.Schema().FieldIndices("_repartition_by_range")[0]).(*array.Uint64)
			require.Equal(t, tc.wantRepartition, repartition.Value(0))
			require.Equal(t, tc.wantRepartition, repartition.Value(1))

			if !tc.wantMemo {
				require.Empty(t, record.Schema().FieldIndices("memo"))
				return
			}
			memos := record.Column(record.Schema().FieldIndices("memo")[0]).(*array.String)
			require.Equal(t, tronTestMemo, memos.Value(0))
			require.False(t, memos.IsNull(1))
			require.Equal(t, "", memos.Value(1))
		})
	}
}
