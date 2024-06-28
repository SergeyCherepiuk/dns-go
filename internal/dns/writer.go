package dns

import (
	"errors"
	"strings"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

type PacketWriter struct {
	buf   []byte
	pos   int
	cache map[int]string
}

var (
	ErrTooManyBytes    = errors.New("too many bytes")
	ErrIndexOutOfBound = errors.New("index out of bound")
)

func NewPacketWriter() *PacketWriter {
	return &PacketWriter{
		buf:   make([]byte, 0),
		pos:   0,
		cache: make(map[int]string),
	}
}

func (w *PacketWriter) WriteUint16(uint16 uint16) error {
	if MaxPacketSize < w.pos+2 {
		return ErrTooManyBytes
	}

	bytes := utils.Uint16ToBytes(uint16)
	w.buf = append(w.buf, bytes[:]...)
	w.pos += 2
	return nil
}

func (w *PacketWriter) WriteUint32(uint32 uint32) error {
	if MaxPacketSize < w.pos+4 {
		return ErrTooManyBytes
	}

	bytes := utils.Uint32ToBytes(uint32)
	w.buf = append(w.buf, bytes[:]...)
	w.pos += 4
	return nil
}

func (w *PacketWriter) WriteByte(byte byte) error {
	if MaxPacketSize < w.pos+1 {
		return ErrTooManyBytes
	}

	w.buf = append(w.buf, byte)
	w.pos += 1
	return nil
}

func (w *PacketWriter) WriteBytes(bytes []byte) error {
	if MaxPacketSize < w.pos+len(bytes) {
		return ErrTooManyBytes
	}

	w.buf = append(w.buf, bytes...)
	w.pos += len(bytes)
	return nil
}

func (w *PacketWriter) WriteDomain(domain string) error {
	var bytes []byte

	subdomains := strings.Split(domain, ".")
	for i, subdomain := range subdomains {
		joined := strings.Join(subdomains[i:], ".")

		pointer, ok := utils.KeyByValue(w.cache, joined)
		if ok {
			pointerBytes := utils.Uint16ToBytes(uint16(pointer))
			bytes = append(bytes, pointerBytes[:]...)
			break
		}

		if joined != "" {
			pointer := (len(domain) - len(joined) + w.pos) | 0b11000000_00000000
			w.cache[pointer] = joined
		}

		size := byte(len(subdomain))
		bytes = append(bytes, size)
		bytes = append(bytes, subdomain...)
	}

	return w.WriteBytes(bytes)
}

func (w *PacketWriter) Bytes() []byte {
	return w.buf
}
