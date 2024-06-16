package dns

import (
	"testing"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

func TestHeaderConstants(t *testing.T) {
	entries := make(utils.DiffEntries, 0)

	entries = append(entries, utils.Diff(HeaderSize, 12)...)

	entries = append(entries, utils.Diff(PacketTypeQuery, 0)...)
	entries = append(entries, utils.Diff(PacketTypeResponse, 1)...)

	entries = append(entries, utils.Diff(OpcodeQuery, 0)...)
	entries = append(entries, utils.Diff(OpcodeIQuery, 1)...)
	entries = append(entries, utils.Diff(OpcodeStatus, 2)...)

	entries = append(entries, utils.Diff(ResponseCodeNoError, 0)...)
	entries = append(entries, utils.Diff(ResponseCodeFormatError, 1)...)
	entries = append(entries, utils.Diff(ResponseCodeServerFailure, 2)...)
	entries = append(entries, utils.Diff(ResponseCodeNameError, 3)...)
	entries = append(entries, utils.Diff(ResponseCodeNotImplemented, 4)...)
	entries = append(entries, utils.Diff(ResponseCodeRefused, 5)...)

	if len(entries) > 0 {
		t.Fatal(entries.String())
	}
}

func TestMarshalHeaderQueryPacket(t *testing.T) {
	header := Header{
		ID:                           0x4edb,
		PacketType:                   PacketTypeQuery,
		Opcode:                       OpcodeQuery,
		AuthoritativeAnswer:          false,
		Truncated:                    false,
		RecursionDesired:             true,
		RecursionAvailable:           false,
		AuthenticData:                true,
		CheckingDisabled:             false,
		ResponseCode:                 ResponseCodeNoError,
		QuestionSectionSize:          1,
		AnswerSectionSize:            0,
		AuthorityRecordsSectionSize:  0,
		AdditionalRecordsSectionSize: 0,
	}

	expectedBytes := [HeaderSize]byte{
		0x4e, 0xdb, 0x01, 0x20, 0x00, 0x01, 0x00, 0x00,
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

	expectedBytes := [HeaderSize]byte{
		0x4e, 0xdb, 0x81, 0x80, 0x00, 0x01, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00,
	}

	actualBytes := MarshalHeader(header)

	if actualBytes != expectedBytes {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalHeaderQueryPacket(t *testing.T) {
	bytes := [HeaderSize]byte{
		0x4e, 0xdb, 0x01, 0x20, 0x00, 0x01, 0x00, 0x00,
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
		AuthenticData:                true,
		CheckingDisabled:             false,
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
	bytes := [HeaderSize]byte{
		0x4e, 0xdb, 0x81, 0x80, 0x00, 0x01, 0x00, 0x01,
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
