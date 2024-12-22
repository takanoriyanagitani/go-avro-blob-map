package b2wasm2blob

import (
	"bytes"
	"context"
	"errors"
	"io"
	"iter"
	"os"
	"path/filepath"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"
)

var (
	ErrInvalidBlob error = errors.New("invalid blob")
)

type InputBlob []byte
type OutputBlob []byte

type BlobToWasmToBlob func(InputBlob) IO[OutputBlob]

type WasmSource func(moduleId string) IO[[]byte]

const WasmFileSizeMaxDefault int64 = 16777216

func IdToBasenameDefault(moduleId string) string { return moduleId + ".wasm" }

func WasmSourceFsLimited(
	limit int64,
	dirname string,
	id2basename func(string) string,
) WasmSource {
	var buf bytes.Buffer
	return func(moduleId string) IO[[]byte] {
		return func(_ context.Context) ([]byte, error) {
			buf.Reset()

			var basename string = id2basename(moduleId)
			var fullname string = filepath.Join(dirname, basename)
			f, e := os.Open(fullname)
			if nil != e {
				return nil, e
			}
			defer f.Close()
			limited := &io.LimitedReader{
				R: f,
				N: limit,
			}

			_, e = io.Copy(&buf, limited)
			return buf.Bytes(), e
		}
	}
}

func WasmSourceFsLimitedDefault(
	dirname string,
) WasmSource {
	return WasmSourceFsLimited(
		WasmFileSizeMaxDefault,
		dirname,
		IdToBasenameDefault,
	)
}

func (b BlobToWasmToBlob) AnyToBlob(
	input any,
	buf *bytes.Buffer,
) IO[OutputBlob] {
	return func(ctx context.Context) (OutputBlob, error) {
		buf.Reset()

		switch t := input.(type) {
		case []byte:
			return b(t)(ctx)
		case string:
			_, _ = buf.WriteString(t) // error is always nil or OOM
			return b(buf.Bytes())(ctx)
		case nil:
			return nil, nil
		default:
			return nil, ErrInvalidBlob
		}
	}
}

func (b BlobToWasmToBlob) MapsToMaps(
	m iter.Seq2[map[string]any, error],
	blobKey string,
) IO[iter.Seq2[map[string]any, error]] {
	return func(ctx context.Context) (iter.Seq2[map[string]any, error], error) {
		mbuf := map[string]any{}
		var buf bytes.Buffer
		return func(yield func(map[string]any, error) bool) {
			for row, e := range m {
				clear(mbuf)

				if nil != e {
					yield(nil, e)
					return
				}

				for key, val := range row {
					if key == blobKey {
						continue
					}

					mbuf[key] = val
				}

				var blob any = row[blobKey]
				c, e := b.AnyToBlob(blob, &buf)(ctx)
				mbuf[blobKey] = []byte(c)

				if !yield(mbuf, e) {
					return
				}
			}
		}, nil
	}
}
