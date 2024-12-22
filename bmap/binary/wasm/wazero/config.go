package j2jwa0

import (
	"context"

	wz "github.com/tetratelabs/wazero"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"
)

var RuntimeConfigDefault wz.RuntimeConfig = wz.NewRuntimeConfig()

var ModuleConfigDefault wz.ModuleConfig = wz.NewModuleConfig().
	WithName("")

const (
	AllocateFnNameDefault   string = "b2b_allocate"
	ConvertFnNameDefault    string = "b2b_convert"
	InputSizeFnNameDefault  string = "b2b_input_size"
	OutputSizeFnNameDefault string = "b2b_output_size"
	OffsetNameDefaultI      string = "b2b_offset_i"
	OffsetNameDefaultO      string = "b2b_offset_o"
)

type Config struct {
	wz.RuntimeConfig
	wz.ModuleConfig

	Allocate   string
	Convert    string
	InputSize  string
	OutputSize string

	OffsetI string
	OffsetO string
}

func (c Config) ToRuntime() IO[Runtime] {
	return func(ctx context.Context) (Runtime, error) {
		var rtm wz.Runtime = wz.NewRuntimeWithConfig(ctx, c.RuntimeConfig)
		return Runtime{
			Runtime:      rtm,
			ModuleConfig: c.ModuleConfig,

			Allocate:   c.Allocate,
			Convert:    c.Convert,
			InputSize:  c.InputSize,
			OutputSize: c.OutputSize,

			OffsetI: c.OffsetI,
			OffsetO: c.OffsetO,
		}, nil
	}
}

var ConfigDefault Config = Config{
	RuntimeConfig: RuntimeConfigDefault,
	ModuleConfig:  ModuleConfigDefault,

	Allocate:   AllocateFnNameDefault,
	Convert:    ConvertFnNameDefault,
	InputSize:  InputSizeFnNameDefault,
	OutputSize: OutputSizeFnNameDefault,

	OffsetI: OffsetNameDefaultI,
	OffsetO: OffsetNameDefaultO,
}
