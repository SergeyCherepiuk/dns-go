package dns

import (
	"testing"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

func TestMarshalHeaderQueryPacket(t *testing.T) {
	header := Header{
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

	expectedBytes := [12]byte{
		0x4e, 0xdb, 0x01, 0x00,
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	actualBytes := MarshalHeader(header)

	if actualBytes != expectedBytes {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestMarshalHeaderResponsePacket(t *testing.T) {
	header := Header{
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

	expectedBytes := [12]byte{
		0x4e, 0xdb, 0x81, 0x80,
		0x00, 0x01, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00,
	}

	actualBytes := MarshalHeader(header)

	if actualBytes != expectedBytes {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalHeaderQueryPacket(t *testing.T) {
	bytes := [12]byte{
		0x4e, 0xdb, 0x01, 0x00,
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
		entries := utils.Diff(actualHeader, expectedHeader)
		t.Fatal(entries.String())
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
		entries := utils.Diff(actualHeader, expectedHeader)
		t.Fatal(entries.String())
	}
}
