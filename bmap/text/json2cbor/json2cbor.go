package json2cbor

import (
	"bytes"
	"context"
	"errors"
	"iter"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"
)

var (
	ErrInvalidBlob error = errors.New("invalid blob")
)

type JsonBlob []byte
type CborBlob []byte

type JsonToCbor func(JsonBlob) IO[CborBlob]

func (j JsonToCbor) AnyToCborBlob(
	input any,
	buf *bytes.Buffer,
) IO[CborBlob] {
	return func(ctx context.Context) (CborBlob, error) {
		buf.Reset()

		switch t := input.(type) {
		case []byte:
			return j(t)(ctx)
		case string:
			_, _ = buf.WriteString(t) // error is always nil or OOM
			return j(buf.Bytes())(ctx)
		case nil:
			return nil, nil
		default:
			return nil, ErrInvalidBlob
		}
	}
}

func (j JsonToCbor) MapsToMaps(
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
				c, e := j.AnyToCborBlob(blob, &buf)(ctx)
				mbuf[blobKey] = []byte(c)

				if !yield(mbuf, e) {
					return
				}
			}
		}, nil
	}
}
