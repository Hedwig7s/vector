package coord3

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Hedwig7s/vector"
)

const componentCount = 3

// Write writes the components as binary. Only the fixed size unsigned types
// (uint8, uint16, uint32, uint64) are supported, uint is not as its size
// depends on the platform.
func (v Coord[T]) Write(out io.Writer, endian binary.ByteOrder) (err error) {
	switch vv := any(v).(type) {
	case Uint8:
		_, err = out.Write([]byte{vv.x, vv.y, vv.z})
		return

	case Uint16:
		bytes := make([]byte, 2*componentCount)
		endian.PutUint16(bytes[0:], vv.x)
		endian.PutUint16(bytes[2:], vv.y)
		endian.PutUint16(bytes[4:], vv.z)
		_, err = out.Write(bytes)
		return

	case Uint32:
		bytes := make([]byte, 4*componentCount)
		endian.PutUint32(bytes[0:], vv.x)
		endian.PutUint32(bytes[4:], vv.y)
		endian.PutUint32(bytes[8:], vv.z)
		_, err = out.Write(bytes)
		return

	case Uint64:
		bytes := make([]byte, 8*componentCount)
		endian.PutUint64(bytes[0:], vv.x)
		endian.PutUint64(bytes[8:], vv.y)
		endian.PutUint64(bytes[16:], vv.z)
		_, err = out.Write(bytes)
		return
	}

	panic(fmt.Errorf("write unimplemented type: %#v", v))
}

func Read[T vector.Unsigned](in io.Reader, endian binary.ByteOrder) (v Coord[T], err error) {
	switch any(v).(type) {
	case Uint8:
		vv, err := ReadUint8(in)
		return any(vv).(Coord[T]), err

	case Uint16:
		vv, err := ReadUint16(in, endian)
		return any(vv).(Coord[T]), err

	case Uint32:
		vv, err := ReadUint32(in, endian)
		return any(vv).(Coord[T]), err

	case Uint64:
		vv, err := ReadUint64(in, endian)
		return any(vv).(Coord[T]), err
	}

	panic(fmt.Errorf("read unimplemented type: %#v", v))
}

func ReadUint8(in io.Reader) (Coord[uint8], error) {
	buf := make([]byte, componentCount)
	_, err := io.ReadFull(in, buf)
	return Coord[uint8]{
		x: buf[0],
		y: buf[1],
		z: buf[2],
	}, err
}

func ReadUint16(in io.Reader, endian binary.ByteOrder) (Coord[uint16], error) {
	buf := make([]byte, componentCount*2)
	_, err := io.ReadFull(in, buf)
	return Coord[uint16]{
		x: endian.Uint16(buf[0:]),
		y: endian.Uint16(buf[2:]),
		z: endian.Uint16(buf[4:]),
	}, err
}

func ReadUint32(in io.Reader, endian binary.ByteOrder) (Coord[uint32], error) {
	buf := make([]byte, componentCount*4)
	_, err := io.ReadFull(in, buf)
	return Coord[uint32]{
		x: endian.Uint32(buf[0:]),
		y: endian.Uint32(buf[4:]),
		z: endian.Uint32(buf[8:]),
	}, err
}

func ReadUint64(in io.Reader, endian binary.ByteOrder) (Coord[uint64], error) {
	buf := make([]byte, componentCount*8)
	_, err := io.ReadFull(in, buf)
	return Coord[uint64]{
		x: endian.Uint64(buf[0:]),
		y: endian.Uint64(buf[8:]),
		z: endian.Uint64(buf[16:]),
	}, err
}
