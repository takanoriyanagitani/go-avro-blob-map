package main

import (
	"context"
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"strings"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"

	bw "github.com/takanoriyanagitani/go-avro-blob-map/bmap/binary/wasm"
	w2 "github.com/takanoriyanagitani/go-avro-blob-map/bmap/binary/wasm/wazero"

	dh "github.com/takanoriyanagitani/go-avro-blob-map/avro/dec/hamba"
	eh "github.com/takanoriyanagitani/go-avro-blob-map/avro/enc/hamba"
)

var EnvValByKey func(string) IO[string] = Lift(
	func(key string) (string, error) {
		val, found := os.LookupEnv(key)
		switch found {
		case true:
			return val, nil
		default:
			return "", fmt.Errorf("env var %s missing", key)
		}
	},
)

var blobKey IO[string] = EnvValByKey("ENV_BLOB_KEY")

var stdin2avro2maps IO[iter.Seq2[map[string]any, error]] = dh.
	StdinToMapsDefault

var wcfg w2.Config = w2.ConfigDefault
var rtm IO[w2.Runtime] = wcfg.ToRuntime()

var wasmFileDirname IO[string] = EnvValByKey("ENV_WASM_MODULES_DIR")

var wasmSource IO[bw.WasmSource] = Bind(
	wasmFileDirname,
	Lift(func(dname string) (bw.WasmSource, error) {
		return bw.WasmSourceFsLimitedDefault(dname), nil
	}),
)

var wasmModuleId IO[string] = EnvValByKey("ENV_WASM_MODULE_ID")

var wasmBytes IO[[]byte] = Bind(
	wasmSource,
	func(ws bw.WasmSource) IO[[]byte] {
		return Bind(
			wasmModuleId,
			ws,
		)
	},
)

var compiled IO[w2.Compiled] = Bind(
	rtm,
	func(r w2.Runtime) IO[w2.Compiled] {
		return Bind(
			wasmBytes,
			r.ToCompiled,
		)
	},
)
var mdl IO[w2.Module] = Bind(
	compiled,
	func(c w2.Compiled) IO[w2.Module] { return c.ToModule() },
)
var mapw0 IO[w2.MapWazero] = Bind(
	mdl,
	func(m w2.Module) IO[w2.MapWazero] { return m.ToMapper() },
)
var wmap IO[bw.BlobToWasmToBlob] = Bind(
	mapw0,
	Lift(func(mw w2.MapWazero) (bw.BlobToWasmToBlob, error) {
		return mw.ToMapper(), nil
	}),
)

var mapd IO[iter.Seq2[map[string]any, error]] = Bind(
	stdin2avro2maps,
	func(
		original iter.Seq2[map[string]any, error],
	) IO[iter.Seq2[map[string]any, error]] {
		return Bind(
			blobKey,
			func(bk string) IO[iter.Seq2[map[string]any, error]] {
				return Bind(
					wmap,
					func(
						mapper bw.BlobToWasmToBlob,
					) IO[iter.Seq2[map[string]any, error]] {
						return mapper.MapsToMaps(
							original,
							bk,
						)
					},
				)
			},
		)
	},
)

var schemaFilename IO[string] = EnvValByKey("ENV_SCHEMA_FILENAME")

func FilenameToStringLimited(limit int64) func(string) IO[string] {
	return Lift(
		func(filename string) (string, error) {
			f, e := os.Open(filename)
			if nil != e {
				return "", e
			}
			defer f.Close()

			limited := &io.LimitedReader{
				R: f,
				N: limit,
			}
			var buf strings.Builder
			_, e = io.Copy(&buf, limited)
			return buf.String(), e
		},
	)
}

const SchemaFileSizeLimitDefault int64 = 1048576

var schemaContent IO[string] = Bind(
	schemaFilename,
	FilenameToStringLimited(SchemaFileSizeLimitDefault),
)

var stdin2maps2mapd2avro2stdout IO[Void] = Bind(
	schemaContent,
	func(s string) IO[Void] {
		return Bind(
			mapd,
			eh.SchemaToMapsToStdoutDefault(s),
		)
	},
)

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return stdin2maps2mapd2avro2stdout(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
