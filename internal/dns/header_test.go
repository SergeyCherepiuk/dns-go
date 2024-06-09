package dns

import "testing"

// TODO: Improve error message by creating diff util function

func TestUnmarshalHeaderQueryPacket(t *testing.T) {
	bytes := [12]byte{
		0x4e, 0xdb, 0x01, 0x20,
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	expectedHeader := Header{
		ID:                           0x4edb,
		PacketType:                   PacketTypeQuery,
		Opcode:                       OpcodeQuery,
		AuthoritativeAnswer:          false,
		Truncated:                    false,
		RecursionDesired:             true,
		RecursionAvailable:           false,
		ResponseCode:                 ResponseCodeNoError,
		QuestionSectionSize:          1,
		AnswerSectionSize:            0,
		AuthorityRecordsSectionSize:  0,
		AdditionalRecordsSectionSize: 0,
	}

	actualHeader := UnmarshalHeader(bytes)

	if actualHeader != expectedHeader {
		t.Fatalf("\nE: %+v\nG: %+v\n", expectedHeader, actualHeader)
	}
}

func TestUnmarshalHeaderResponsePacket(t *testing.T) {
	bytes := [12]byte{
		0x4e, 0xdb, 0x81, 0x80,
		0x00, 0x01, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00,
	}

	expectedHeader := Header{
		ID:                           0x4edb,
		PacketType:                   PacketTypeResponse,
		Opcode:                       OpcodeQuery,
		AuthoritativeAnswer:          false,
		Truncated:                    false,
		RecursionDesired:             true,
		RecursionAvailable:           true,
		ResponseCode:                 ResponseCodeNoError,
		QuestionSectionSize:          1,
		AnswerSectionSize:            1,
		AuthorityRecordsSectionSize:  0,
		AdditionalRecordsSectionSize: 0,
	}

	actualHeader := UnmarshalHeader(bytes)

	if actualHeader != expectedHeader {
		t.Fatalf("\nE: %+v\nG: %+v\n", expectedHeader, actualHeader)
	}
}