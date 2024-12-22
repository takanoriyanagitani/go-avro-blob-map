package j2jwa0

import (
	"context"
	"errors"
	"fmt"

	wa "github.com/tetratelabs/wazero/api"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"
)

var (
	ErrUnableToAllocate error = errors.New("unable to allocate")
)

type Allocate struct{ wa.Function }

func (a Allocate) Allocate(sz uint32) IO[uint32] {
	return func(ctx context.Context) (uint32, error) {
		results, e := a.Function.Call(ctx, wa.EncodeU32(sz))
		if nil != e {
			return 0, fmt.Errorf("%w: %w", ErrUnableToAllocate, e)
		}

		if 1 != len(results) {
			return 0, ErrUnableToAllocate
		}

		var u uint64 = results[0]
		var i int32 = wa.DecodeI32(u)
		if i < 0 {
			return 0, ErrUnableToAllocate
		}
		return uint32(i), nil
	}
}
