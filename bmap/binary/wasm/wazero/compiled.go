package j2jwa0

import (
	"context"

	wz "github.com/tetratelabs/wazero"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"
)

type Compiled struct {
	wz.Runtime
	wz.CompiledModule
	wz.ModuleConfig

	Allocate   string
	Convert    string
	InputSize  string
	OutputSize string

	OffsetI string
	OffsetO string
}

func (c Compiled) ToModule() IO[Module] {
	return func(ctx context.Context) (Module, error) {
		mdl, e := c.Runtime.InstantiateModule(
			ctx,
			c.CompiledModule,
			c.ModuleConfig,
		)
		return Module{
			Module: mdl,

			Allocate:   c.Allocate,
			Convert:    c.Convert,
			InputSize:  c.InputSize,
			OutputSize: c.OutputSize,

			OffsetI: c.OffsetI,
			OffsetO: c.OffsetO,
		}, e
	}
}
