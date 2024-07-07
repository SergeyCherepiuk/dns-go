package io

import (
	"errors"
	"strings"

	"github.com/SergeyCherepiuk/dns-go/internal/dns/types"
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
	if types.MaxPacketSize < w.pos+2 {
		return ErrTooManyBytes
	}

	bytes := utils.Uint16ToBytes(uint16)
	w.buf = append(w.buf, bytes[:]...)
	w.pos += 2
	return nil
}

func (w *PacketWriter) WriteUint32(uint32 uint32) error {
	if types.MaxPacketSize < w.pos+4 {
		return ErrTooManyBytes
	}

	bytes := utils.Uint32ToBytes(uint32)
	w.buf = append(w.buf, bytes[:]...)
	w.pos += 4
	return nil
}

func (w *PacketWriter) WriteByte(byte byte) error {
	if types.MaxPacketSize < w.pos+1 {
		return ErrTooManyBytes
	}

	w.buf = append(w.buf, byte)
	w.pos += 1
	return nil
}

func (w *PacketWriter) WriteBytes(bytes []byte) error {
	if types.MaxPacketSize < w.pos+len(bytes) {
		return ErrTooManyBytes
	}

	w.buf = append(w.buf, bytes...)
	w.pos += len(bytes)
	return nil
}

func (w *PacketWriter) formatDomain(domain string) []byte {
	var bytes []byte

	subdomains := strings.Split(domain, ".")
	for i, subdomain := range subdomains {
		joined := strings.Join(subdomains[i:], ".")

		pointer, ok := utils.KeyByValue(w.cache, joined)
		if ok {
			pointerWithPrefix := pointer | 0b11000000_00000000
			pointerBytes := utils.Uint16ToBytes(uint16(pointerWithPrefix))
			bytes = append(bytes, pointerBytes[:]...)
			break
		}

		size := byte(len(subdomain))
		bytes = append(bytes, size)
		bytes = append(bytes, subdomain...)
	}

	return bytes
}

func (w *PacketWriter) cacheDomain(domain string) {
	initialDomainLength := len(domain)

	for {
		_, ok := utils.KeyByValue(w.cache, domain)
		if domain != "" && !ok {
			pointer := initialDomainLength - len(domain) + w.pos
			w.cache[pointer] = domain
		}

		index := strings.Index(domain, ".")
		if index == -1 {
			break
		}

		domain = domain[index+1:]
	}
}

func (w *PacketWriter) WriteDomain(domain string) error {
	bytes := w.formatDomain(domain)
	w.cacheDomain(domain)
	return w.WriteBytes(bytes)
}

func (w *PacketWriter) WriteDomainWithLength(domain string) error {
	bytes := w.formatDomain(domain)

	length := uint16(len(bytes))
	err := w.WriteUint16(length)
	if err != nil {
		return err
	}

	w.cacheDomain(domain)
	return w.WriteBytes(bytes)
}

func (w *PacketWriter) Bytes() []byte {
	return w.buf
}
