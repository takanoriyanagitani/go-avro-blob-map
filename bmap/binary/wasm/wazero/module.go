package j2jwa0

import (
	"context"
	"errors"

	wa "github.com/tetratelabs/wazero/api"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"
)

var (
	ErrInvalidMemory   error = errors.New("invalid memory")
	ErrInvalidFunction error = errors.New("invalid function")
)

type Module struct {
	wa.Module

	Allocate   string
	Convert    string
	InputSize  string
	OutputSize string

	OffsetI string
	OffsetO string
}

func (m Module) GetMemory() IO[Memory] {
	return func(_ context.Context) (Memory, error) {
		var mem wa.Memory = m.Module.Memory()
		switch mem {
		case nil:
			return Memory{}, ErrInvalidMemory
		default:
			return Memory{mem}, nil
		}
	}
}

func (m Module) GetFunction(name string) IO[wa.Function] {
	return func(_ context.Context) (wa.Function, error) {
		var f wa.Function = m.Module.ExportedFunction(name)
		switch f {
		case nil:
			return nil, ErrInvalidFunction
		default:
			return f, nil
		}
	}
}

func (m Module) ToMapper() IO[MapWazero] {
	return Bind(
		m.GetMemory(),
		func(mem Memory) IO[MapWazero] {
			return Bind(
				All(
					m.GetFunction(m.Allocate),
					m.GetFunction(m.Convert),
					m.GetFunction(m.InputSize),
					m.GetFunction(m.OutputSize),

					m.GetFunction(m.OffsetI),
					m.GetFunction(m.OffsetO),
				),
				Lift(func(funcs []wa.Function) (MapWazero, error) {
					return MapWazero{
						Memory: mem,

						Allocate:   Allocate{funcs[0]},
						Convert:    Convert{funcs[1]},
						InputSize:  InputSize{funcs[2]},
						OutputSize: OutputSize{funcs[3]},

						OffsetI: OffsetI{funcs[4]},
						OffsetO: OffsetO{funcs[5]},
					}, nil
				}),
			)
		},
	)
}
