package tables

import (
	"context"
	"testing"

	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/apache/arrow/go/v10/arrow/memory"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	chainstorageapi "github.com/coinbase/chainstorage/protos/coinbase/chainstorage"
	sdkmocks "github.com/coinbase/chainstorage/sdk/mocks"

	"github.com/coinbase/chainsformer/internal/config"
	"github.com/coinbase/chainsformer/internal/controller/internal"
)

func newRawEthereumBlockWithBlob(blobSize int) *chainstorageapi.Block {
	return &chainstorageapi.Block{
		Blobdata: &chainstorageapi.Block_Ethereum{
			Ethereum: &chainstorageapi.EthereumBlobdata{
				Header: make([]byte, blobSize),
			},
		},
	}
}

func newMinimalEthereumNativeBlock() *chainstorageapi.NativeBlock {
	return &chainstorageapi.NativeBlock{
		Block: &chainstorageapi.NativeBlock_Ethereum{
			Ethereum: &chainstorageapi.EthereumBlock{
				Header: &chainstorageapi.EthereumHeader{Hash: "h", Number: 1},
			},
		},
	}
}

func newBlockAndEvent(block *chainstorageapi.Block) *internal.BlockAndEvent {
	return &internal.BlockAndEvent{
		Block: block,
		BlockChainEvent: &chainstorageapi.BlockchainEvent{
			SequenceNum: 1,
			Type:        chainstorageapi.BlockchainEvent_BLOCK_ADDED,
			Block: &chainstorageapi.BlockIdentifier{
				Hash: "h", Height: 1, Tag: 1,
			},
		},
	}
}

func TestEthereumTransformBlock_ReleasesBlobdata_Batch(t *testing.T) {
	cfg := &config.Config{}
	cases := []struct {
		name   string
		schema *arrow.Schema
		fn     func(ctx context.Context, block *chainstorageapi.Block, parser *sdkmocks.MockParser, rb *array.RecordBuilder) error
	}{
		{
			name:   "blocksTable",
			schema: newBlockSchema(cfg),
			fn: func(ctx context.Context, block *chainstorageapi.Block, parser *sdkmocks.MockParser, rb *array.RecordBuilder) error {
				return blocksTable{config: cfg}.TransformBlock(ctx, block, parser, rb, 0)
			},
		},
		{
			name:   "transactionsTable",
			schema: newTransactionSchema(cfg),
			fn: func(ctx context.Context, block *chainstorageapi.Block, parser *sdkmocks.MockParser, rb *array.RecordBuilder) error {
				return transactionsTable{config: cfg}.TransformBlock(ctx, block, parser, rb, 0)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			parser := sdkmocks.NewMockParser(ctrl)

			raw := newRawEthereumBlockWithBlob(1024)
			parser.EXPECT().ParseNativeBlock(gomock.Any(), raw).Return(newMinimalEthereumNativeBlock(), nil)

			rb := array.NewRecordBuilder(memory.DefaultAllocator, tc.schema)
			defer rb.Release()

			require.NoError(t, tc.fn(context.Background(), raw, parser, rb))
			require.Nil(t, raw.Blobdata, "Blobdata must be nil after TransformBlock so the raw bytes can be GC'd")
		})
	}
}

func TestEthereumTransformBlock_ReleasesBlobdata_Stream(t *testing.T) {
	cfg := &config.Config{}
	cases := []struct {
		name   string
		schema *arrow.Schema
		fn     func(ctx context.Context, blockAndEvent *internal.BlockAndEvent, parser *sdkmocks.MockParser, rb *array.RecordBuilder) error
	}{
		{
			name:   "nativeStreamedBlocksTable",
			schema: newStreamedBlocksSchema(cfg),
			fn: func(ctx context.Context, blockAndEvent *internal.BlockAndEvent, parser *sdkmocks.MockParser, rb *array.RecordBuilder) error {
				return nativeStreamedBlocksTable{config: cfg}.TransformBlock(ctx, blockAndEvent, parser, rb, 0)
			},
		},
		{
			name:   "nativeStreamedTransactionsTable",
			schema: newStreamedTransactionSchema(cfg),
			fn: func(ctx context.Context, blockAndEvent *internal.BlockAndEvent, parser *sdkmocks.MockParser, rb *array.RecordBuilder) error {
				return nativeStreamedTransactionsTable{config: cfg}.TransformBlock(ctx, blockAndEvent, parser, rb, 0)
			},
		},
		{
			name:   "rawNativeStreamedTransactionsTable",
			schema: newRawStreamedTransactionSchema(),
			fn: func(ctx context.Context, blockAndEvent *internal.BlockAndEvent, parser *sdkmocks.MockParser, rb *array.RecordBuilder) error {
				return rawNativeStreamedTransactionsTable{config: cfg}.TransformBlock(ctx, blockAndEvent, parser, rb, 0)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			parser := sdkmocks.NewMockParser(ctrl)

			raw := newRawEthereumBlockWithBlob(1024)
			blockAndEvent := newBlockAndEvent(raw)
			parser.EXPECT().ParseNativeBlock(gomock.Any(), raw).Return(newMinimalEthereumNativeBlock(), nil)

			rb := array.NewRecordBuilder(memory.DefaultAllocator, tc.schema)
			defer rb.Release()

			require.NoError(t, tc.fn(context.Background(), blockAndEvent, parser, rb))
			require.Nil(t, raw.Blobdata, "Blobdata must be nil after TransformBlock so the raw bytes can be GC'd")
		})
	}
}
