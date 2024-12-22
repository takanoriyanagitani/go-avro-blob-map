package j2cbor

import (
	"bytes"
	"context"
	"encoding/json"

	ac "github.com/fxamacker/cbor/v2"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"

	tj "github.com/takanoriyanagitani/go-avro-blob-map/bmap/text/json2cbor"
)

func MapToBuffer(
	m map[string]any,
	b *bytes.Buffer,
) error {
	b.Reset()
	return ac.MarshalToBuffer(m, b)
}

func JsonToCborStd(
	j tj.JsonBlob,
	mbuf map[string]any,
	buf *bytes.Buffer,
) IO[tj.CborBlob] {
	return func(_ context.Context) (tj.CborBlob, error) {
		clear(mbuf)
		e := json.Unmarshal(j, &mbuf)
		if nil != e {
			return nil, e
		}
		e = MapToBuffer(mbuf, buf)
		return buf.Bytes(), e
	}
}

func JsonToCborStdNew() tj.JsonToCbor {
	var mbuf map[string]any = map[string]any{}
	var buf bytes.Buffer
	return func(j tj.JsonBlob) IO[tj.CborBlob] {
		return JsonToCborStd(
			j,
			mbuf,
			&buf,
		)
	}
}
