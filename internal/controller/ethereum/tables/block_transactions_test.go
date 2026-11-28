package tables

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally/v4"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/coinbase/chainstorage/protos/coinbase/c3/common"
	chainstorageapi "github.com/coinbase/chainstorage/protos/coinbase/chainstorage"
	sdkmocks "github.com/coinbase/chainstorage/sdk/mocks"

	"github.com/coinbase/chainsformer/internal/controller/internal"
	"github.com/coinbase/chainsformer/internal/utils/fxparams"
)

func TestBlockTransactionsTable_SchemaPreservesNativeRecords(t *testing.T) {
	for _, chain := range []common.Blockchain{common.Blockchain_BLOCKCHAIN_ETHEREUM, common.Blockchain_BLOCKCHAIN_POLYGON, common.Blockchain_BLOCKCHAIN_TRON} {
		t.Run(chain.String(), func(t *testing.T) {
			cfg := newChainConfig(chain)
			table := NewBlockTransactionsTable(internal.CommonTableParams{Params: fxparams.Params{Config: cfg, Metrics: tally.NoopScope}})
			require.Equal(t, "table=block_transactions/format=native/encoding=none", table.GetTableName())
			schema := table.GetSchema()
			require.Equal(t, []string{"block_record", "transaction_record"}, fieldNames(schema))
			for i, legacy := range []*arrow.Schema{newBlockSchema(cfg), newTransactionSchema(cfg)} {
				require.True(t, schema.Field(i).Nullable)
				nested, ok := schema.Field(i).Type.(*arrow.StructType)
				require.True(t, ok)
				require.Equal(t, legacy.Fields(), nested.Fields(), "names, types, nullability and metadata must match legacy output")
			}
			block := schema.Field(0).Type.(*arrow.StructType)
			_, hasWithdrawals := block.FieldByName("withdrawals")
			_, hasRoot := block.FieldByName("withdrawals_root")
			require.Equal(t, chain == common.Blockchain_BLOCKCHAIN_ETHEREUM, hasWithdrawals)
			require.Equal(t, hasWithdrawals, hasRoot)
			tx := schema.Field(1).Type.(*arrow.StructType)
			_, hasMemo := tx.FieldByName("memo")
			require.Equal(t, chain == common.Blockchain_BLOCKCHAIN_TRON, hasMemo)
			receipt, _ := tx.FieldByName("receipt")
			_, hasFee := receipt.Type.(*arrow.StructType).FieldByName("energy_fee")
			require.Equal(t, hasMemo, hasFee)
		})
	}
}

// The legacy records are the compatibility oracle for every field (including
// nested receipts, logs, traces, uncle blocks and withdrawals). Literal checks
// below independently protect the identity, direction and partition contracts.
func TestBlockTransactionsTable_OneHeaderAndEveryTransaction(t *testing.T) {
	for _, chain := range []common.Blockchain{common.Blockchain_BLOCKCHAIN_ETHEREUM, common.Blockchain_BLOCKCHAIN_POLYGON, common.Blockchain_BLOCKCHAIN_TRON} {
		for _, count := range []int{0, 1, 3} {
			t.Run(fmt.Sprintf("%s/%d_transactions", chain, count), func(t *testing.T) {
				cfg := newChainConfig(chain)
				native := combinedNativeFixture(t, count)
				parser := sdkmocks.NewMockParser(gomock.NewController(t))
				raw := newRawEthereumBlockWithBlob(1024)
				parser.EXPECT().ParseNativeBlock(gomock.Any(), raw).Return(native, nil).Times(1)
				allocator := memory.NewCheckedAllocator(memory.DefaultAllocator)
				defer allocator.AssertSize(t, 0)
				builder := array.NewRecordBuilder(allocator, newBlockTransactionsSchema(cfg))
				defer builder.Release()
				require.NoError(t, (blockTransactionsTable{config: cfg}).TransformBlock(context.Background(), raw, parser, builder, 100))
				require.Nil(t, raw.Blobdata)
				record := builder.NewRecord()
				defer record.Release()
				require.Equal(t, int64(count+1), record.NumRows())
				blocks := record.Column(0).(*array.Struct)
				transactions := record.Column(1).(*array.Struct)
				for row := 0; row < count+1; row++ {
					require.Equal(t, row != 0, blocks.IsNull(row), "block validity at row %d", row)
					require.Equal(t, row == 0, transactions.IsNull(row), "transaction validity at row %d", row)
				}
				blockBuilder := array.NewRecordBuilder(allocator, newBlockSchema(cfg))
				defer blockBuilder.Release()
				require.NoError(t, (blocksTable{config: cfg}).transformBlocks(blockBuilder, native.GetEthereum(), 100))
				legacyBlock := blockBuilder.NewRecord()
				defer legacyBlock.Release()
				assertCombinedColumns(t, blocks, 0, 1, legacyBlock)
				txBuilder := array.NewRecordBuilder(allocator, newTransactionSchema(cfg))
				defer txBuilder.Release()
				require.NoError(t, (transactionsTable{config: cfg}).transformTransactions(txBuilder, native.GetEthereum(), 100))
				legacyTx := txBuilder.NewRecord()
				defer legacyTx.Release()
				assertCombinedColumns(t, transactions, 1, int64(count+1), legacyTx)
				require.Equal(t, "block-hash", combinedField(blocks, "hash").(*array.String).Value(0))
				require.Equal(t, uint64(12345), combinedField(blocks, "number").(*array.Uint64).Value(0))
				require.Equal(t, uint64(1700000011), combinedField(blocks, "timestamp").(*array.Uint64).Value(0))
				require.Equal(t, uint64(count), combinedField(blocks, "transaction_count").(*array.Uint64).Value(0))
				require.Equal(t, uint64(12300), combinedField(blocks, "_partition_by").(*array.Uint64).Value(0))
				require.Equal(t, uint64(12345), combinedField(blocks, "_repartition_by_range").(*array.Uint64).Value(0))
				require.Equal(t, "uncle-hash", combinedField(blocks, "uncles").(*array.List).ListValues().(*array.String).Value(0))
				if chain == common.Blockchain_BLOCKCHAIN_ETHEREUM {
					withdrawals := combinedField(blocks, "withdrawals").(*array.List).ListValues().(*array.Struct)
					require.Equal(t, uint64(37), combinedField(withdrawals, "amount").(*array.Uint64).Value(0))
					require.Equal(t, "withdrawals-root", combinedField(blocks, "withdrawals_root").(*array.String).Value(0))
				}
				for i := 0; i < count; i++ {
					row := i + 1
					require.Equal(t, fmt.Sprintf("tx-%d", i), combinedField(transactions, "transaction_hash").(*array.String).Value(row))
					require.Equal(t, uint64(i), combinedField(transactions, "transaction_index").(*array.Uint64).Value(row))
					require.Equal(t, "0x1111111111111111111111111111111111111111", combinedField(transactions, "from_address").(*array.String).Value(row))
					require.Equal(t, "0x2222222222222222222222222222222222222222", combinedField(transactions, "to_address").(*array.String).Value(row))
					require.Equal(t, "987654321", combinedField(transactions, "value_string").(*array.String).Value(row))
					require.Equal(t, uint64(12300), combinedField(transactions, "_partition_by").(*array.Uint64).Value(row))
					if chain == common.Blockchain_BLOCKCHAIN_TRON {
						require.Equal(t, "6d656d6f", combinedField(transactions, "memo").(*array.String).Value(row))
						receipt := combinedField(transactions, "receipt").(*array.Struct)
						require.Equal(t, uint64(113), combinedField(receipt, "fee").(*array.Uint64).Value(row))
						require.Equal(t, uint64(157), combinedField(receipt, "energy_penalty_total").(*array.Uint64).Value(row))
					}
				}
			})
		}
	}
}

func assertCombinedColumns(t *testing.T, nested *array.Struct, start, end int64, legacy arrow.Record) {
	t.Helper()
	require.Equal(t, int(legacy.NumCols()), nested.NumField())
	for i, field := range legacy.Schema().Fields() {
		got := array.NewSlice(nested.Field(i), start, end)
		defer got.Release()
		require.True(t, array.Equal(legacy.Column(i), got), "field %s: want %s, got %s", field.Name, legacy.Column(i), got)
	}
}

func combinedField(record *array.Struct, name string) arrow.Array {
	index, ok := record.DataType().(*arrow.StructType).FieldIdx(name)
	if !ok {
		panic("missing field " + name)
	}
	return record.Field(index)
}

func TestBlockTransactionsTable_RejectsInvalidInput(t *testing.T) {
	parseErr := errors.New("parser unavailable")
	noHeader := &chainstorageapi.NativeBlock{Block: &chainstorageapi.NativeBlock_Ethereum{Ethereum: &chainstorageapi.EthereumBlock{}}}
	badDifficulty := combinedNativeFixture(t, 1)
	badDifficulty.GetEthereum().Header.TotalDifficulty = "invalid"
	for _, tc := range []struct {
		name     string
		native   *chainstorageapi.NativeBlock
		parseErr error
		message  string
	}{
		{"parser_failure", nil, parseErr, "parser unavailable"},
		{"missing_native", nil, nil, "ethereum block"},
		{"missing_ethereum", &chainstorageapi.NativeBlock{}, nil, "ethereum block"},
		{"missing_header", noHeader, nil, "header is required"},
		{"invalid_total_difficulty", badDifficulty, nil, "total difficulty"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := newChainConfig(common.Blockchain_BLOCKCHAIN_ETHEREUM)
			parser := sdkmocks.NewMockParser(gomock.NewController(t))
			raw := newRawEthereumBlockWithBlob(16)
			blob := raw.Blobdata
			parser.EXPECT().ParseNativeBlock(gomock.Any(), raw).Return(tc.native, tc.parseErr).Times(1)
			builder := array.NewRecordBuilder(memory.DefaultAllocator, newBlockTransactionsSchema(cfg))
			defer builder.Release()
			err := (blockTransactionsTable{config: cfg}).TransformBlock(context.Background(), raw, parser, builder, 100)
			require.ErrorContains(t, err, tc.message)
			if tc.parseErr != nil {
				require.ErrorIs(t, err, parseErr)
				require.Same(t, blob, raw.Blobdata, "failed parsing must preserve raw data")
			}
			record := builder.NewRecord()
			defer record.Release()
			require.Zero(t, record.NumRows(), "invalid blocks must not leave partial rows")
		})
	}
}

func combinedNativeFixture(t *testing.T, count int) *chainstorageapi.NativeBlock {
	t.Helper()
	block := &chainstorageapi.EthereumBlock{}
	require.NoError(t, protojson.Unmarshal([]byte(`{
 "header": {"hash":"block-hash","parentHash":"parent-hash","number":"12345","nonce":"nonce","sha3Uncles":"uncles-root","logsBloom":"bloom","transactionsRoot":"transactions-root","stateRoot":"state-root","receiptsRoot":"receipts-root","miner":"0x3333333333333333333333333333333333333333","difficulty":"17","totalDifficulty":"123456789012345","size":"19","extraData":"extra","gasLimit":"21000000","gasUsed":"42001","timestamp":"2023-11-14T22:13:31Z","baseFeePerGas":"23","uncles":["uncle-hash"],"withdrawalsRoot":"withdrawals-root","withdrawals":[{"index":"29","validatorIndex":"31","address":"0x4444444444444444444444444444444444444444","amount":"37"}]},
 "uncles":[{"hash":"uncle-hash","parentHash":"uncle-parent","number":"12344","timestamp":"2023-11-14T22:13:30Z","miner":"uncle-miner","difficulty":"41","gasLimit":"43","gasUsed":"47","baseFeePerGas":"53"}],
 "transactions":[{"hash":"tx-0","index":"0","blockHash":"block-hash","blockNumber":"12345","blockTimestamp":"2023-11-14T22:13:31Z","from":"0x1111111111111111111111111111111111111111","to":"0x2222222222222222222222222222222222222222","nonce":"59","value":"987654321","gas":"61001","gasPrice":"67","input":"0xabcd","type":"2","maxFeePerGas":"71","maxPriorityFeePerGas":"73","priorityFeePerGas":"79","memo":"6d656d6f",
 "receipt":{"transactionHash":"tx-0","transactionIndex":"0","blockHash":"block-hash","blockNumber":"12345","from":"receipt-from","to":"receipt-to","cumulativeGasUsed":"83","gasUsed":"89","contractAddress":"contract-address","logsBloom":"receipt-bloom","root":"receipt-root","type":"2","status":"1","effectiveGasPrice":"97","fee":"113","netFee":"127","netUsage":"131","energyUsage":"137","energyFee":"139","originEnergyUsage":"149","energyUsageTotal":"151","energyPenaltyTotal":"157","logs":[{"logIndex":"101","transactionHash":"tx-0","transactionIndex":"0","blockHash":"block-hash","blockNumber":"12345","address":"log-address","data":"0xcdef","topics":["topic-one","topic-two"],"removed":true}]},
 "flattenedTraces":[{"transactionHash":"tx-0","transactionIndex":"0","blockHash":"block-hash","blockNumber":"12345","from":"trace-from","to":"trace-to","value":"103","input":"trace-input","output":"trace-output","type":"CALL","traceType":"call","callType":"delegatecall","gas":"107","gasUsed":"109","subtraces":"2","traceAddress":["3","5"],"error":"trace-error","status":"0","traceId":"trace-id"}]}]
 }`), block))
	template := block.Transactions[0]
	block.Transactions = nil
	for i := 0; i < count; i++ {
		tx := proto.Clone(template).(*chainstorageapi.EthereumTransaction)
		tx.Hash = fmt.Sprintf("tx-%d", i)
		tx.Index = uint64(i)
		tx.Receipt.TransactionHash = tx.Hash
		tx.Receipt.TransactionIndex = tx.Index
		block.Transactions = append(block.Transactions, tx)
		block.Header.Transactions = append(block.Header.Transactions, tx.Hash)
	}
	return &chainstorageapi.NativeBlock{Block: &chainstorageapi.NativeBlock_Ethereum{Ethereum: block}}
}
