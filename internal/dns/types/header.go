package types

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
