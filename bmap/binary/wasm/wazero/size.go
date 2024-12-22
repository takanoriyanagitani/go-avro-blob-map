package j2jwa0

import (
	"context"
	"errors"

	wa "github.com/tetratelabs/wazero/api"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"
)

var (
	ErrUnableToGetSize error = errors.New("unable to get size")
)

type Size struct{ wa.Function }

type InputSize Size
type OutputSize Size

func (i Size) GetSize() IO[uint32] {
	return func(ctx context.Context) (uint32, error) {
		results, e := i.Function.Call(ctx)
		if nil != e {
			return 0, e
		}

		if 1 != len(results) {
			return 0, ErrUnableToGetSize
		}

		var u uint64 = results[0]
		var i int32 = wa.DecodeI32(u)
		if i < 0 {
			return 0, ErrUnableToGetSize
		}
		return uint32(i), nil
	}
}
