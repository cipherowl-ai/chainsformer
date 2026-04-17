package tables

import (
	"context"
	"runtime"
	"testing"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	chainstorageapi "github.com/coinbase/chainstorage/protos/coinbase/chainstorage"
	sdkmocks "github.com/coinbase/chainstorage/sdk/mocks"
)

func newRawBitcoinBlockWithBlob(blobSize int) *chainstorageapi.Block {
	return &chainstorageapi.Block{
		Blobdata: &chainstorageapi.Block_Bitcoin{
			Bitcoin: &chainstorageapi.BitcoinBlobdata{
				Header: make([]byte, blobSize),
			},
		},
	}
}

func newMinimalBitcoinNativeBlock() *chainstorageapi.NativeBlock {
	return &chainstorageapi.NativeBlock{
		Block: &chainstorageapi.NativeBlock_Bitcoin{
			Bitcoin: &chainstorageapi.BitcoinBlock{
				Header: &chainstorageapi.BitcoinHeader{Hash: "h", Height: 1},
			},
		},
	}
}

type transformFn func(ctx context.Context, block *chainstorageapi.Block, parser *sdkmocks.MockParser, rb *array.RecordBuilder) error

func TestBitcoinTransformBlock_ReleasesBlobdata(t *testing.T) {
	cases := []struct {
		name   string
		schema *arrow.Schema
		fn     transformFn
	}{
		{
			name:   "blocksTable",
			schema: newBlockSchema(),
			fn: func(ctx context.Context, block *chainstorageapi.Block, parser *sdkmocks.MockParser, rb *array.RecordBuilder) error {
				return blocksTable{}.TransformBlock(ctx, block, parser, rb, 0)
			},
		},
		{
			name:   "transactionsTable",
			schema: newTransactionSchema(),
			fn: func(ctx context.Context, block *chainstorageapi.Block, parser *sdkmocks.MockParser, rb *array.RecordBuilder) error {
				return transactionsTable{}.TransformBlock(ctx, block, parser, rb, 0)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			parser := sdkmocks.NewMockParser(ctrl)

			raw := newRawBitcoinBlockWithBlob(1024)
			parser.EXPECT().ParseNativeBlock(gomock.Any(), raw).Return(newMinimalBitcoinNativeBlock(), nil)

			rb := array.NewRecordBuilder(memory.DefaultAllocator, tc.schema)
			defer rb.Release()

			require.NoError(t, tc.fn(context.Background(), raw, parser, rb))
			require.Nil(t, raw.Blobdata, "Blobdata must be nil after TransformBlock so the raw bytes can be GC'd")
		})
	}
}

// BenchmarkBitcoinTransformBlock_MemoryFootprint measures the resident-heap
// delta across a TransformBlock call with a large raw Blobdata attached.
// Run with:
//
//	go test -bench=BenchmarkBitcoinTransformBlock_MemoryFootprint \
//	        -run=^$ -benchtime=3x ./internal/controller/bitcoin/tables/
//
// Compare the resident_delta_B metric between master and this branch. On
// master the delta should be ~blobSize (bytes pinned by raw.Blobdata); on
// this branch it should be near zero (Blobdata was nil'd, GC reclaimed).
func BenchmarkBitcoinTransformBlock_MemoryFootprint(b *testing.B) {
	const blobSize = 64 * 1024 * 1024 // 64 MiB — large enough to dwarf GC noise

	for i := 0; i < b.N; i++ {
		ctrl := gomock.NewController(b)
		parser := sdkmocks.NewMockParser(ctrl)
		raw := newRawBitcoinBlockWithBlob(blobSize)
		parser.EXPECT().ParseNativeBlock(gomock.Any(), raw).Return(newMinimalBitcoinNativeBlock(), nil)
		rb := array.NewRecordBuilder(memory.DefaultAllocator, newBlockSchema())

		runtime.GC()
		var before runtime.MemStats
		runtime.ReadMemStats(&before)

		if err := (blocksTable{}).TransformBlock(context.Background(), raw, parser, rb, 0); err != nil {
			b.Fatal(err)
		}

		runtime.GC()
		var after runtime.MemStats
		runtime.ReadMemStats(&after)

		// raw is still reachable here, simulating the caller's slice reference.
		// If Blobdata was nil'd inside TransformBlock, the blobSize bytes are
		// GC-eligible and the delta should be near zero.
		runtime.KeepAlive(raw)
		b.ReportMetric(float64(int64(after.HeapAlloc)-int64(before.HeapAlloc)), "resident_delta_B")

		rb.Release()
		ctrl.Finish()
	}
}
