package j2jwa0

import (
	"context"
	"errors"

	wa "github.com/tetratelabs/wazero/api"

	. "github.com/takanoriyanagitani/go-avro-blob-map/util"
)

var (
	ErrUnableToCopy error = errors.New("unable to copy input bytes")
)

type Memory struct{ wa.Memory }

func (m Memory) CopyToWasm(offset uint32, data []byte) IO[Void] {
	return func(_ context.Context) (Void, error) {
		var wrote bool = m.Memory.Write(offset, data)
		switch wrote {
		case true:
			return Empty, nil
		default:
			return Empty, ErrUnableToCopy
		}
	}
}

func (m Memory) CopyFromWasm(offset uint32, sz uint32) IO[[]byte] {
	return func(_ context.Context) ([]byte, error) {
		read, ok := m.Memory.Read(offset, sz)
		switch ok {
		case true:
			return read, nil
		default:
			return nil, ErrUnableToCopy
		}
	}
}
