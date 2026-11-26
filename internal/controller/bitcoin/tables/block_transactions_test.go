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

	chainstorageapi "github.com/coinbase/chainstorage/protos/coinbase/chainstorage"
	sdkmocks "github.com/coinbase/chainstorage/sdk/mocks"

	"github.com/coinbase/chainsformer/internal/config"
	"github.com/coinbase/chainsformer/internal/controller/internal"
	"github.com/coinbase/chainsformer/internal/utils/fxparams"
)

func TestBlockTransactionsTable_SchemaPreservesNativeRecords(t *testing.T) {
	table := NewBlockTransactionsTable(internal.CommonTableParams{Params: fxparams.Params{Config: &config.Config{}, Metrics: tally.NoopScope}})
	require.Equal(t, "table=block_transactions/format=native/encoding=none", table.GetTableName())
	schema := table.GetSchema()
	require.Equal(t, 2, schema.NumFields())
	for i, tc := range []struct {
		name   string
		legacy *arrow.Schema
	}{{"block_record", newBlockSchema()}, {"transaction_record", newTransactionSchema()}} {
		require.Equal(t, tc.name, schema.Field(i).Name)
		require.True(t, schema.Field(i).Nullable)
		nested, ok := schema.Field(i).Type.(*arrow.StructType)
		require.True(t, ok)
		require.Equal(t, tc.legacy.Fields(), nested.Fields(), "names, types, nullability and metadata must match legacy output")
	}
}

// The legacy records are the compatibility oracle for every field; the literal
// checks pin identity, input/output contents and partitioning independently.
func TestBlockTransactionsTable_OneHeaderAndEveryTransaction(t *testing.T) {
	for _, count := range []int{0, 1, 3} {
		t.Run(fmt.Sprintf("%d_transactions", count), func(t *testing.T) {
			native := combinedNativeFixture(count)
			parser := sdkmocks.NewMockParser(gomock.NewController(t))
			raw := newRawBitcoinBlockWithBlob(1024)
			parser.EXPECT().ParseNativeBlock(gomock.Any(), raw).Return(native, nil).Times(1)
			allocator := memory.NewCheckedAllocator(memory.DefaultAllocator)
			defer allocator.AssertSize(t, 0)
			builder := array.NewRecordBuilder(allocator, newBlockTransactionsSchema())
			defer builder.Release()
			require.NoError(t, (blockTransactionsTable{}).TransformBlock(context.Background(), raw, parser, builder, 100))
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

			blockBuilder := array.NewRecordBuilder(allocator, newBlockSchema())
			defer blockBuilder.Release()
			require.NoError(t, (blocksTable{}).transformBlocks(blockBuilder, native.GetBitcoin(), 100))
			legacyBlock := blockBuilder.NewRecord()
			defer legacyBlock.Release()
			assertCombinedColumns(t, blocks, 0, 1, legacyBlock)
			txBuilder := array.NewRecordBuilder(allocator, newTransactionSchema())
			defer txBuilder.Release()
			require.NoError(t, (transactionsTable{}).transformTransactions(txBuilder, native.GetBitcoin(), 100))
			legacyTx := txBuilder.NewRecord()
			defer legacyTx.Release()
			assertCombinedColumns(t, transactions, 1, int64(count+1), legacyTx)

			require.Equal(t, "block-hash", combinedField(blocks, "hash").(*array.String).Value(0))
			require.Equal(t, uint64(812345), combinedField(blocks, "number").(*array.Uint64).Value(0))
			require.Equal(t, uint64(1700000011), combinedField(blocks, "timestamp").(*array.Uint64).Value(0))
			require.Equal(t, uint64(count), combinedField(blocks, "transaction_count").(*array.Uint64).Value(0))
			require.Equal(t, uint64(812300), combinedField(blocks, "_partition_by").(*array.Uint64).Value(0))
			require.Equal(t, uint64(812345), combinedField(blocks, "_repartition_by_range").(*array.Uint64).Value(0))
			blockTxIDs := combinedField(blocks, "transactions").(*array.List).ListValues().(*array.String)
			require.Equal(t, count, blockTxIDs.Len())
			for i := 0; i < count; i++ {
				row := i + 1
				require.Equal(t, fmt.Sprintf("txid-%d", i), blockTxIDs.Value(i))
				require.Equal(t, fmt.Sprintf("txid-%d", i), combinedField(transactions, "hash").(*array.String).Value(row), "hash is the txid, not the witness hash")
				require.Equal(t, uint64(i), combinedField(transactions, "index").(*array.Uint64).Value(row))
				require.Equal(t, i == 0, combinedField(transactions, "is_coinbase").(*array.Boolean).Value(row))
				require.Equal(t, uint64(812345), combinedField(combinedField(transactions, "block").(*array.Struct), "number").(*array.Uint64).Value(row))
				require.Equal(t, uint64(812300), combinedField(transactions, "_partition_by").(*array.Uint64).Value(row))
				require.Equal(t, uint64(812345), combinedField(transactions, "_repartition_by_range").(*array.Uint64).Value(row))
			}
			if count > 0 {
				// One input per transaction, so flattened input i belongs to transaction i.
				inputs := combinedField(transactions, "inputs").(*array.List).ListValues().(*array.Struct)
				require.Equal(t, "coinbase-data", combinedField(inputs, "coinbase").(*array.String).Value(0))
				if count > 1 {
					require.Equal(t, "spender-address", combinedField(inputs, "address").(*array.String).Value(1))
					require.Equal(t, uint64(5000), combinedField(inputs, "value").(*array.Uint64).Value(1))
					require.Equal(t, "prev-txid", combinedField(inputs, "spent_transaction_hash").(*array.String).Value(1))
				}
				outputs := combinedField(transactions, "outputs").(*array.List).ListValues().(*array.Struct)
				require.Equal(t, "receiver-address", combinedField(outputs, "address").(*array.String).Value(0))
				require.Equal(t, uint64(4000), combinedField(outputs, "value").(*array.Uint64).Value(0))
			}
		})
	}
}

func TestBlockTransactionsTable_RejectsInvalidInput(t *testing.T) {
	parseErr := errors.New("parser unavailable")
	noHeader := &chainstorageapi.NativeBlock{Block: &chainstorageapi.NativeBlock_Bitcoin{Bitcoin: &chainstorageapi.BitcoinBlock{}}}
	for _, tc := range []struct {
		name     string
		native   *chainstorageapi.NativeBlock
		parseErr error
		message  string
	}{
		{"parser_failure", nil, parseErr, "parser unavailable"},
		{"missing_native", nil, nil, "bitcoin block"},
		{"missing_bitcoin", &chainstorageapi.NativeBlock{}, nil, "bitcoin block"},
		{"missing_header", noHeader, nil, "header is required"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			parser := sdkmocks.NewMockParser(gomock.NewController(t))
			raw := newRawBitcoinBlockWithBlob(16)
			blob := raw.Blobdata
			parser.EXPECT().ParseNativeBlock(gomock.Any(), raw).Return(tc.native, tc.parseErr).Times(1)
			builder := array.NewRecordBuilder(memory.DefaultAllocator, newBlockTransactionsSchema())
			defer builder.Release()
			err := (blockTransactionsTable{}).TransformBlock(context.Background(), raw, parser, builder, 100)
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

// Transaction 0 is the coinbase; the rest spend a previous output.
func combinedNativeFixture(count int) *chainstorageapi.NativeBlock {
	block := &chainstorageapi.BitcoinBlock{
		Header: &chainstorageapi.BitcoinHeader{
			Hash: "block-hash", Size: 1_000_001, StrippedSize: 900_001, Weight: 3_999_001, Height: 812345,
			Version: 536870912, MerkleRoot: "merkle-root", Time: 1700000011, Nonce: 2_083_236_893,
			Bits: "17053894", Difficulty: "61030681983175.59", ChainWork: "chain-work",
			NumberOfTransactions: uint64(count), PreviousBlockHash: "previous-hash", NextBlockHash: "next-hash",
		},
	}
	for i := 0; i < count; i++ {
		input := &chainstorageapi.BitcoinTransactionInput{
			TransactionId:             "prev-txid",
			FromOutputIndex:           7,
			ScriptSignature:           &chainstorageapi.BitcoinScriptSignature{Assembly: "sig-asm", Hex: "sig-hex"},
			Sequence:                  4294967293,
			TransactionInputWitnesses: []string{"witness-0", "witness-1"},
			FromOutput: &chainstorageapi.BitcoinTransactionOutput{
				ScriptPublicKey: &chainstorageapi.BitcoinScriptPublicKey{Type: "witness_v0_keyhash", Address: "spender-address"},
				Value:           5000,
			},
		}
		if i == 0 {
			input = &chainstorageapi.BitcoinTransactionInput{Coinbase: "coinbase-data", Sequence: 4294967295}
		}
		block.Transactions = append(block.Transactions, &chainstorageapi.BitcoinTransaction{
			TransactionId: fmt.Sprintf("txid-%d", i), Hash: fmt.Sprintf("witness-hash-%d", i),
			Size: 222, VirtualSize: 141, Weight: 561, Version: 2, LockTime: 812344,
			IsCoinbase: i == 0, Index: uint64(i),
			Inputs: []*chainstorageapi.BitcoinTransactionInput{input},
			Outputs: []*chainstorageapi.BitcoinTransactionOutput{{
				Index:           0,
				ScriptPublicKey: &chainstorageapi.BitcoinScriptPublicKey{Assembly: "out-asm", Hex: "out-hex", Type: "witness_v0_keyhash", Address: "receiver-address"},
				Value:           4000,
			}},
			InputCount: 1, OutputCount: 1, InputValue: 5000, OutputValue: 4000, Fee: 1000,
		})
	}
	return &chainstorageapi.NativeBlock{Block: &chainstorageapi.NativeBlock_Bitcoin{Bitcoin: block}}
}

// Dogecoin Core's getblock has no nTx, so the parsed header count is 0.
func TestBlockTransactionsTable_ZeroHeaderCountFallsBackToTransactions(t *testing.T) {
	native := combinedNativeFixture(3)
	native.GetBitcoin().Header.NumberOfTransactions = 0
	parser := sdkmocks.NewMockParser(gomock.NewController(t))
	raw := newRawBitcoinBlockWithBlob(1024)
	parser.EXPECT().ParseNativeBlock(gomock.Any(), raw).Return(native, nil).Times(1)
	allocator := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer allocator.AssertSize(t, 0)
	builder := array.NewRecordBuilder(allocator, newBlockTransactionsSchema())
	defer builder.Release()
	require.NoError(t, (blockTransactionsTable{}).TransformBlock(context.Background(), raw, parser, builder, 100))
	record := builder.NewRecord()
	defer record.Release()

	blocks := record.Column(0).(*array.Struct)
	require.Equal(t, uint64(3), combinedField(blocks, "transaction_count").(*array.Uint64).Value(0))
	txBlock := combinedField(record.Column(1).(*array.Struct), "block").(*array.Struct)
	for row := 1; row < 4; row++ {
		require.Equal(t, uint64(3), combinedField(txBlock, "transaction_count").(*array.Uint64).Value(row), "row %d", row)
	}
}
