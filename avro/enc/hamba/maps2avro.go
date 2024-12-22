package enc

import (
	"context"
	"io"
	"iter"
	"os"

	ha "github.com/hamba/avro/v2"
	ho "github.com/hamba/avro/v2/ocf"

	ab "github.com/takanoriyanagitani/go-avro-blob-map"
	. "github.com/takanoriyanagitani/go-avro-blob-map/util"
)

func MapsToAvroToWriterHamba(
	ctx context.Context,
	m iter.Seq2[map[string]any, error],
	s ha.Schema,
	w io.Writer,
	opts ...ho.EncoderFunc,
) error {
	enc, e := ho.NewEncoderWithSchema(
		s,
		w,
		opts...,
	)
	if nil != e {
		return e
	}
	defer enc.Close()

	for row, e := range m {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if nil != e {
			return e
		}

		e := enc.Encode(row)
		if nil != e {
			return e
		}

		e = enc.Flush()
		if nil != e {
			return e
		}
	}
	return enc.Flush()
}

func CodecConv(c ab.Codec) ho.CodecName {
	switch c {
	case ab.CodecNull:
		return ho.Null
	case ab.CodecDeflate:
		return ho.Deflate
	case ab.CodecSnappy:
		return ho.Snappy
	case ab.CodecZstd:
		return ho.ZStandard
	default:
		return ho.Null
	}
}

func ConfigToOpts(cfg ab.EncodeConfig) []ho.EncoderFunc {
	var blockLen int = cfg.BlockLength
	var codec ab.Codec = cfg.Codec
	var hc ho.CodecName = CodecConv(codec)
	return []ho.EncoderFunc{
		ho.WithBlockLength(blockLen),
		ho.WithCodec(hc),
	}
}

func MapsToAvroToWriter(
	ctx context.Context,
	m iter.Seq2[map[string]any, error],
	schema string,
	w io.Writer,
	cfg ab.EncodeConfig,
) error {
	parsed, e := ha.Parse(schema)
	if nil != e {
		return e
	}

	var opts []ho.EncoderFunc = ConfigToOpts(cfg)
	return MapsToAvroToWriterHamba(
		ctx,
		m,
		parsed,
		w,
		opts...,
	)
}

func ConfigToSchemaToMapsToWriter(
	cfg ab.EncodeConfig,
) func(io.Writer) func(string) func(iter.Seq2[map[string]any, error]) IO[Void] {
	return func(
		w io.Writer,
	) func(string) func(iter.Seq2[map[string]any, error]) IO[Void] {
		return func(
			schema string,
		) func(iter.Seq2[map[string]any, error]) IO[Void] {
			return func(
				m iter.Seq2[map[string]any, error],
			) IO[Void] {
				return func(ctx context.Context) (Void, error) {
					return Empty, MapsToAvroToWriter(
						ctx,
						m,
						schema,
						w,
						cfg,
					)
				}
			}
		}
	}
}

func SchemaToMapsToWriterDefault(
	w io.Writer,
) func(string) func(iter.Seq2[map[string]any, error]) IO[Void] {
	return ConfigToSchemaToMapsToWriter(ab.EncodeConfigDefault)(w)
}

func SchemaToMapsToStdoutDefault(
	schema string,
) func(iter.Seq2[map[string]any, error]) IO[Void] {
	return SchemaToMapsToWriterDefault(os.Stdout)(schema)
}
