package dns

import (
	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

type PacketType uint8

const (
	PacketTypeQuery PacketType = iota
	PacketTypeResponse
)

type Opcode uint8

const (
	OpcodeQuery Opcode = iota
	OpcodeIQuery
	OpcodeStatus
)

type ResponseCode uint8

const (
	ResponseCodeNoError ResponseCode = iota
	ResponseCodeFormatError
	ResponseCodeServerFailure
	ResponseCodeNameError
	ResponseCodeNotImplemented
	ResponseCodeRefused
)

type Header struct {
	ID                           uint16
	PacketType                   PacketType
	Opcode                       Opcode
	AuthoritativeAnswer          bool
	Truncated                    bool
	RecursionDesired             bool
	RecursionAvailable           bool
	ResponseCode                 ResponseCode
	QuestionSectionSize          uint16
	AnswerSectionSize            uint16
	AuthorityRecordsSectionSize  uint16
	AdditionalRecordsSectionSize uint16
}

// TODO: Not implemented
func MarshalHeader(header Header) [12]byte {
	return [12]byte{}
}

func UnmarshalHeader(bytes [12]byte) Header {
	var (
		packetTypeBit          = (bytes[2] >> 7) & 0b00000001
		opcodeBits             = (bytes[2] >> 3) & 0b00001111
		authoritativeAnswerBit = (bytes[2] >> 2) & 0b00000001
		truncatedBit           = (bytes[2] >> 1) & 0b00000001
		recursionDesiredBit    = (bytes[2] >> 0) & 0b00000001
		recursionAvailableBit  = (bytes[3] >> 7) & 0b00000001
		responseCodeBits       = (bytes[3] >> 0) & 0b00001111
	)

	return Header{
		ID:                           utils.BytesToUint16([2]byte(bytes[0:2])),
		PacketType:                   PacketType(packetTypeBit),
		Opcode:                       Opcode(opcodeBits),
		AuthoritativeAnswer:          authoritativeAnswerBit == 1,
		Truncated:                    truncatedBit == 1,
		RecursionDesired:             recursionDesiredBit == 1,
		RecursionAvailable:           recursionAvailableBit == 1,
		ResponseCode:                 ResponseCode(responseCodeBits),
		QuestionSectionSize:          utils.BytesToUint16([2]byte(bytes[4:6])),
		AnswerSectionSize:            utils.BytesToUint16([2]byte(bytes[6:8])),
		AuthorityRecordsSectionSize:  utils.BytesToUint16([2]byte(bytes[8:10])),
		AdditionalRecordsSectionSize: utils.BytesToUint16([2]byte(bytes[10:12])),
	}
}
