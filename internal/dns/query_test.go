package dns

import (
	"reflect"
	"testing"

	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

func TestUnmarshalQueryPacketOneQuestion(t *testing.T) {
	bytes := []byte{
		// Header
		0x12, 0x34, 0x01, 0x20, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		// Question 1 (google.com.)
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,
	}

	expectedQueryPackcet := QueryPacket{
		Header: Header{
			ID:                           0x1234,
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
		},
		Questions: []Question{
			{"google.com.", QuestionTypeA, QuestionClassIN},
		},
	}

	actualQueryPacket := UnmarshalQueryPacket(bytes)

	if !reflect.DeepEqual(actualQueryPacket, expectedQueryPackcet) {
		entries := utils.Diff(actualQueryPacket, expectedQueryPackcet)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalQueryPacketTwoQuestion(t *testing.T) {
	bytes := []byte{
		// Header
		0x12, 0x34, 0x01, 0x20, 0x00, 0x02, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		// Question 1 (google.com.)
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,

		// Question 2 (mx.google.com.)
		0x02, 0x6d, 0x78, 0xc0, 0x0c, 0x00, 0x00, 0x0f,
		0x00, 0x01,
	}

	expectedQueryPackcet := QueryPacket{
		Header: Header{
			ID:                           0x1234,
			PacketType:                   PacketTypeQuery,
			Opcode:                       OpcodeQuery,
			AuthoritativeAnswer:          false,
			Truncated:                    false,
			RecursionDesired:             true,
			RecursionAvailable:           false,
			AuthenticData:                true,
			CheckingDisabled:             false,
			ResponseCode:                 ResponseCodeNoError,
			QuestionSectionSize:          2,
			AnswerSectionSize:            0,
			AuthorityRecordsSectionSize:  0,
			AdditionalRecordsSectionSize: 0,
		},
		Questions: []Question{
			{"google.com.", QuestionTypeA, QuestionClassIN},
			{"mx.google.com.", QuestionTypeMX, QuestionClassIN},
		},
	}

	actualQueryPacket := UnmarshalQueryPacket(bytes)

	if !reflect.DeepEqual(actualQueryPacket, expectedQueryPackcet) {
		entries := utils.Diff(actualQueryPacket, expectedQueryPackcet)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalQueryPacketThreeQuestion(t *testing.T) {
	bytes := []byte{
		// Header
		0x12, 0x34, 0x01, 0x20, 0x00, 0x03, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		// Question 1 (google.com.)
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,

		// Question 2 (mx.google.com.)
		0x02, 0x6d, 0x78, 0xc0, 0x0c, 0x00, 0x00, 0x0f,
		0x00, 0x01,

		// Question 3 (com.)
		0xc0, 0x13, 0x00, 0x00, 0x01, 0x00, 0x01,
	}

	expectedQueryPackcet := QueryPacket{
		Header: Header{
			ID:                           0x1234,
			PacketType:                   PacketTypeQuery,
			Opcode:                       OpcodeQuery,
			AuthoritativeAnswer:          false,
			Truncated:                    false,
			RecursionDesired:             true,
			RecursionAvailable:           false,
			AuthenticData:                true,
			CheckingDisabled:             false,
			ResponseCode:                 ResponseCodeNoError,
			QuestionSectionSize:          3,
			AnswerSectionSize:            0,
			AuthorityRecordsSectionSize:  0,
			AdditionalRecordsSectionSize: 0,
		},
		Questions: []Question{
			{"google.com.", QuestionTypeA, QuestionClassIN},
			{"mx.google.com.", QuestionTypeMX, QuestionClassIN},
			{"com.", QuestionTypeA, QuestionClassIN},
		},
	}

	actualQueryPacket := UnmarshalQueryPacket(bytes)

	if !reflect.DeepEqual(actualQueryPacket, expectedQueryPackcet) {
		entries := utils.Diff(actualQueryPacket, expectedQueryPackcet)
		t.Fatal(entries.String())
	}
}
