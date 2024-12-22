package j2jwa0

import (
	"context"

	wz "github.com/tetratelabs/wazero"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"
)

type Runtime struct {
	wz.Runtime
	wz.ModuleConfig

	Allocate   string
	Convert    string
	InputSize  string
	OutputSize string

	OffsetI string
	OffsetO string
}

func (r Runtime) ToCloser() IO[Void] {
	return func(ctx context.Context) (Void, error) {
		return Empty, r.Runtime.Close(ctx)
	}
}

func (r Runtime) ToCompiled(wasmBytes []byte) IO[Compiled] {
	return func(ctx context.Context) (Compiled, error) {
		compiled, e := r.Runtime.CompileModule(
			ctx,
			wasmBytes,
		)
		return Compiled{
			Runtime:        r.Runtime,
			CompiledModule: compiled,
			ModuleConfig:   r.ModuleConfig,

			Allocate:   r.Allocate,
			Convert:    r.Convert,
			InputSize:  r.InputSize,
			OutputSize: r.OutputSize,

			OffsetI: r.OffsetI,
			OffsetO: r.OffsetO,
		}, e
	}
}
