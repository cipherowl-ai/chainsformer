package tables

import (
	"context"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"golang.org/x/xerrors"

	api "github.com/coinbase/chainstorage/protos/coinbase/chainstorage"
	"github.com/coinbase/chainstorage/sdk"

	"github.com/coinbase/chainsformer/internal/config"
	"github.com/coinbase/chainsformer/internal/controller/internal"
	"github.com/coinbase/chainsformer/internal/utils/xarrow"
)

type blockTransactionsTable struct{ config *config.Config }

// NewBlockTransactionsTable reads and parses each block once for both outputs.
// A block row is followed by its transaction rows; exactly one nested record is
// populated per row.
func NewBlockTransactionsTable(params internal.CommonTableParams) internal.Table {
	return internal.NewBatchTable(&params,
		internal.NewTableAttributes(internal.TableNameBlockTransactions),
		newBlockTransactionsSchema(),
		blockTransactionsTable{config: params.Config})
}

func newBlockTransactionsSchema() *arrow.Schema {
	f := xarrow.NewSchemaFactory()
	return f.NewSchema(
		f.NewField("block_record", arrow.StructOf(newBlockSchema().Fields()...), "Complete legacy block row, populated only for block records"),
		f.NewField("transaction_record", arrow.StructOf(newTransactionSchema().Fields()...), "Complete legacy transaction row, populated only for transaction records"),
	)
}

func (t blockTransactionsTable) TransformBlock(ctx context.Context, raw *api.Block, parser sdk.Parser, rb *array.RecordBuilder, partitionBySize uint64) error {
	native, err := parser.ParseNativeBlock(ctx, raw)
	if err != nil {
		return xerrors.Errorf("failed to parse raw block to native block: %w", err)
	}
	block := native.GetBitcoin()
	if block == nil {
		return xerrors.New("failed to extract bitcoin block from native block")
	}
	// Checked before any append so an invalid block leaves no partial row.
	if block.GetHeader() == nil {
		return xerrors.New("header is required")
	}
	raw.Blobdata = nil
	blocks := rb.Field(0).(*array.StructBuilder)
	transactions := rb.Field(1).(*array.StructBuilder)
	if err := (blocksTable{config: t.config}).transformBlockWithAppender(xarrow.NewStructRecordAppender(blocks), block, partitionBySize); err != nil {
		return xerrors.Errorf("failed to transform block: %w", err)
	}
	blocks.Append(true)
	transactions.AppendNull()
	return (transactionsTable{config: t.config}).transformTransactionsWithAppender(func() *xarrow.RecordAppender {
		blocks.AppendNull()
		transactions.Append(true)
		return xarrow.NewStructRecordAppender(transactions)
	}, block, partitionBySize)
}
