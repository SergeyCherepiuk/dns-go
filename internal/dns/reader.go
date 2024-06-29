package dns

import (
	"errors"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

type PacketReader struct {
	buf []byte
	pos int
}

var (
	ErrInvalidPacketSize = errors.New("invalid packet size")
	ErrNotEnoughBytes    = errors.New("not enough bytes")
)

func newPacketReader(bytes []byte) (*PacketReader, error) {
	if len(bytes) < 12 || len(bytes) > 512 {
		return nil, ErrInvalidPacketSize
	}
	return &PacketReader{bytes, 0}, nil
}

func (r *PacketReader) ReadUint16() (uint16, error) {
	if len(r.buf) < r.pos+2 {
		return 0, ErrNotEnoughBytes
	}

	bytes := [2]byte(r.buf[r.pos : r.pos+2])
	uint16 := utils.BytesToUint16(bytes)
	r.pos += 2
	return uint16, nil
}

func (r *PacketReader) ReadUint32() (uint32, error) {
	if len(r.buf) < r.pos+4 {
		return 0, ErrNotEnoughBytes
	}

	bytes := [4]byte(r.buf[r.pos : r.pos+4])
	uint32 := utils.BytesToUint32(bytes)
	r.pos += 4
	return uint32, nil
}

func (r *PacketReader) ReadByte() (byte, error) {
	if len(r.buf) < r.pos+1 {
		return 0, ErrNotEnoughBytes
	}

	byte := r.buf[r.pos]
	r.pos += 1
	return byte, nil
}

func (r *PacketReader) ReadByteAt(pos int) (byte, error) {
	if len(r.buf) < pos+1 {
		return 0, ErrNotEnoughBytes
	}

	byte := r.buf[pos]
	return byte, nil
}

func (r *PacketReader) ReadBytes(n int) ([]byte, error) {
	if len(r.buf) < r.pos+n {
		return nil, ErrNotEnoughBytes
	}

	bytes := r.buf[r.pos : r.pos+n]
	r.pos += n
	return bytes, nil
}

func (r *PacketReader) ReadBytesAt(n int, pos int) ([]byte, error) {
	if len(r.buf) < pos+n {
		return nil, ErrNotEnoughBytes
	}

	bytes := r.buf[pos : pos+n]
	return bytes, nil
}

func (r *PacketReader) ReadDomain() (string, error) {
	domain := make([]byte, 0)
	pointer := -1

	for {
		var (
			size byte
			err  error
		)

		if pointer >= 0 {
			size, err = r.ReadByteAt(pointer)
			pointer += 1
		} else {
			size, err = r.ReadByte()
		}

		if err != nil {
			return "", err
		}

		if size == 0 {
			break
		}

		if size&0b11000000 == 0b11000000 {
			var secondSizeByte byte

			if pointer >= 0 {
				secondSizeByte, err = r.ReadByteAt(pointer)
				pointer += 1
			} else {
				secondSizeByte, err = r.ReadByte()
			}

			if err != nil {
				return "", err
			}

			pointerBytes := [2]byte{byte(size) & 0b00111111, secondSizeByte}
			pointer = int(utils.BytesToUint16(pointerBytes))
			continue
		}

		var bytes []byte

		if pointer >= 0 {
			bytes, err = r.ReadBytesAt(int(size), pointer)
			pointer += int(size)
		} else {
			bytes, err = r.ReadBytes(int(size))
		}

		if err != nil {
			return "", err
		}

		domain = append(domain, bytes...)
		domain = append(domain, '.')
	}

	return string(domain), nil
}
