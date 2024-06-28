package dns

import (
	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

const HeaderSize = 12

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
	AuthenticData                bool
	CheckingDisabled             bool
	ResponseCode                 ResponseCode
	QuestionSectionSize          uint16
	AnswerSectionSize            uint16
	AuthorityRecordsSectionSize  uint16
	AdditionalRecordsSectionSize uint16
}

func marshalHeader(w *PacketWriter, header Header) error {
	var (
		packetTypeBit          = uint16(header.PacketType) << 15
		opcodeBits             = uint16(header.Opcode) << 11
		authoritativeAnswerBit = uint16(utils.BoolToUint8(header.AuthoritativeAnswer)) << 10
		truncatedBit           = uint16(utils.BoolToUint8(header.Truncated)) << 9
		recursionDesiredBit    = uint16(utils.BoolToUint8(header.RecursionDesired)) << 8
		recursionAvailableBit  = uint16(utils.BoolToUint8(header.RecursionAvailable)) << 7
		authenticDataBit       = uint16(utils.BoolToUint8(header.AuthenticData)) << 5
		checkingDisabledBit    = uint16(utils.BoolToUint8(header.CheckingDisabled)) << 4
		responseCodeBits       = uint16(header.ResponseCode)

		flags = packetTypeBit | opcodeBits | authoritativeAnswerBit | truncatedBit | recursionDesiredBit |
			recursionAvailableBit | authenticDataBit | checkingDisabledBit | responseCodeBits
	)

	var (
		idBits                           = utils.Uint16ToBytes(header.ID)
		flagsBits                        = utils.Uint16ToBytes(flags)
		questionSectionSizeBits          = utils.Uint16ToBytes(header.QuestionSectionSize)
		answerSectionSizeBits            = utils.Uint16ToBytes(header.AnswerSectionSize)
		authorityRecordsSectionSizeBits  = utils.Uint16ToBytes(header.AuthorityRecordsSectionSize)
		additionalRecordsSectionSizeBits = utils.Uint16ToBytes(header.AdditionalRecordsSectionSize)
	)

	bytes := []byte{
		idBits[0], idBits[1],
		flagsBits[0], flagsBits[1],
		questionSectionSizeBits[0], questionSectionSizeBits[1],
		answerSectionSizeBits[0], answerSectionSizeBits[1],
		authorityRecordsSectionSizeBits[0], authorityRecordsSectionSizeBits[1],
		additionalRecordsSectionSizeBits[0], additionalRecordsSectionSizeBits[1],
	}

	return w.WriteBytes(bytes)
}

func unmarshalHeader(r *PacketReader) (Header, error) {
	bytes, err := r.ReadBytes(HeaderSize)
	if err != nil {
		return Header{}, err
	}

	var (
		packetTypeBit          = (bytes[2] >> 7) & 0b00000001
		opcodeBits             = (bytes[2] >> 3) & 0b00001111
		authoritativeAnswerBit = (bytes[2] >> 2) & 0b00000001
		truncatedBit           = (bytes[2] >> 1) & 0b00000001
		recursionDesiredBit    = (bytes[2] >> 0) & 0b00000001
		recursionAvailableBit  = (bytes[3] >> 7) & 0b00000001
		authenticDataBit       = (bytes[3] >> 5) & 0b00000001
		chekingDisabledBit     = (bytes[3] >> 4) & 0b00000001
		responseCodeBits       = (bytes[3] >> 0) & 0b00001111
	)

	header := Header{
		ID:                           utils.BytesToUint16([2]byte(bytes[0:2])),
		PacketType:                   PacketType(packetTypeBit),
		Opcode:                       Opcode(opcodeBits),
		AuthoritativeAnswer:          authoritativeAnswerBit == 1,
		Truncated:                    truncatedBit == 1,
		RecursionDesired:             recursionDesiredBit == 1,
		RecursionAvailable:           recursionAvailableBit == 1,
		AuthenticData:                authenticDataBit == 1,
		CheckingDisabled:             chekingDisabledBit == 1,
		ResponseCode:                 ResponseCode(responseCodeBits),
		QuestionSectionSize:          utils.BytesToUint16([2]byte(bytes[4:6])),
		AnswerSectionSize:            utils.BytesToUint16([2]byte(bytes[6:8])),
		AuthorityRecordsSectionSize:  utils.BytesToUint16([2]byte(bytes[8:10])),
		AdditionalRecordsSectionSize: utils.BytesToUint16([2]byte(bytes[10:12])),
	}

	return header, nil
}
