package j2jwa0

import (
	"context"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"

	bw "github.com/takanoriyanagitani/go-avro-blob-map/bmap/binary/wasm"
)

type MapWazero struct {
	Memory

	Allocate
	Convert
	InputSize
	OutputSize

	OffsetI
	OffsetO
}

func (m MapWazero) CopyToWasm(original []byte) IO[Void] {
	return func(ctx context.Context) (Void, error) {
		var sz int = len(original)

		_, e := m.Allocate.Allocate(uint32(sz))(ctx)
		if nil != e {
			return Empty, e
		}

		offset, e := Offset(m.OffsetI).GetOffset()(ctx)
		if nil != e {
			return Empty, e
		}

		_, e = m.Memory.CopyToWasm(offset, original)(ctx)
		return Empty, e
	}
}

func (m MapWazero) ToMapper() bw.BlobToWasmToBlob {
	return func(original bw.InputBlob) IO[bw.OutputBlob] {
		return func(ctx context.Context) (bw.OutputBlob, error) {
			_, e := m.CopyToWasm(original)(ctx)
			if nil != e {
				return nil, e
			}

			_, e = m.Convert.Convert()(ctx)
			if nil != e {
				return nil, e
			}

			outputOffset, e := Offset(m.OffsetO).GetOffset()(ctx)
			if nil != e {
				return nil, e
			}

			outputSize, e := Size(m.OutputSize).GetSize()(ctx)
			if nil != e {
				return nil, e
			}

			return m.Memory.CopyFromWasm(
				outputOffset,
				outputSize,
			)(ctx)
		}
	}
}
