package j2jwa0

import (
	"context"
	"errors"
	"fmt"

	wa "github.com/tetratelabs/wazero/api"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"
)

var (
	ErrUnableToConvert error = errors.New("unable to convert")
)

type Convert struct{ wa.Function }

func (a Convert) Convert() IO[uint32] {
	return func(ctx context.Context) (uint32, error) {
		results, e := a.Function.Call(ctx)
		if nil != e {
			return 0, fmt.Errorf("%w: %w", ErrUnableToConvert, e)
		}

		if 1 != len(results) {
			return 0, ErrUnableToConvert
		}

		var u uint64 = results[0]
		var i int32 = wa.DecodeI32(u)
		if i < 0 {
			return 0, ErrUnableToConvert
		}
		return uint32(i), nil
	}
}
