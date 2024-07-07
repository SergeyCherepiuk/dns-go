package serde

import (
	"slices"
	"testing"

	"github.com/SergeyCherepiuk/dns-go/internal/dns/types"
	"github.com/SergeyCherepiuk/dns-go/internal/utils"
)

// TODO: Improve tests

func TestMarshalPacketOneQuestion(t *testing.T) {
	packet := types.Packet{
		Header: types.Header{
			ID:                           0x1234,
			PacketType:                   types.PacketTypeQuery,
			Opcode:                       types.OpcodeQuery,
			AuthoritativeAnswer:          false,
			Truncated:                    false,
			RecursionDesired:             true,
			RecursionAvailable:           false,
			AuthenticData:                true,
			CheckingDisabled:             false,
			ResponseCode:                 types.ResponseCodeNoError,
			QuestionSectionSize:          1,
			AnswerSectionSize:            0,
			AuthorityRecordsSectionSize:  0,
			AdditionalRecordsSectionSize: 0,
		},
		Questions: []types.Question{
			{Domain: "google.com.", Type: types.QuestionTypeA, Class: types.QuestionClassIN},
		},
	}

	expectedBytes := []byte{
		// Header
		0x12, 0x34, 0x01, 0x20, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		// Question 1 (google.com.)
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,
	}

	actualBytes, err := MarshalPacket(packet)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(actualBytes, expectedBytes) {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestMarshalQueryPacketTwoQuestions(t *testing.T) {
	packet := types.Packet{
		Header: types.Header{
			ID:                           0x1234,
			PacketType:                   types.PacketTypeQuery,
			Opcode:                       types.OpcodeQuery,
			AuthoritativeAnswer:          false,
			Truncated:                    false,
			RecursionDesired:             true,
			RecursionAvailable:           false,
			AuthenticData:                true,
			CheckingDisabled:             false,
			ResponseCode:                 types.ResponseCodeNoError,
			QuestionSectionSize:          2,
			AnswerSectionSize:            0,
			AuthorityRecordsSectionSize:  0,
			AdditionalRecordsSectionSize: 0,
		},
		Questions: []types.Question{
			{Domain: "google.com.", Type: types.QuestionTypeA, Class: types.QuestionClassIN},
			{Domain: "mx.google.com.", Type: types.QuestionTypeMX, Class: types.QuestionClassIN},
		},
	}

	expectedBytes := []byte{
		// Header
		0x12, 0x34, 0x01, 0x20, 0x00, 0x02, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		// Question 1 (google.com.)
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,

		// Question 2 (mx.google.com.)
		0x02, 0x6d, 0x78, 0xc0, 0x0c, 0x00, 0x0f, 0x00,
		0x01,
	}

	actualBytes, err := MarshalPacket(packet)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(actualBytes, expectedBytes) {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestMarshalQueryPacketThreeQuestions(t *testing.T) {
	packet := types.Packet{
		Header: types.Header{
			ID:                           0x1234,
			PacketType:                   types.PacketTypeQuery,
			Opcode:                       types.OpcodeQuery,
			AuthoritativeAnswer:          false,
			Truncated:                    false,
			RecursionDesired:             true,
			RecursionAvailable:           false,
			AuthenticData:                true,
			CheckingDisabled:             false,
			ResponseCode:                 types.ResponseCodeNoError,
			QuestionSectionSize:          3,
			AnswerSectionSize:            0,
			AuthorityRecordsSectionSize:  0,
			AdditionalRecordsSectionSize: 0,
		},
		Questions: []types.Question{
			{Domain: "google.com.", Type: types.QuestionTypeA, Class: types.QuestionClassIN},
			{Domain: "mx.google.com.", Type: types.QuestionTypeMX, Class: types.QuestionClassIN},
			{Domain: "com.", Type: types.QuestionTypeA, Class: types.QuestionClassIN},
		},
	}

	expectedBytes := []byte{
		// Header
		0x12, 0x34, 0x01, 0x20, 0x00, 0x03, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		// Question 1 (google.com.)
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,

		// Question 2 (mx.google.com.)
		0x02, 0x6d, 0x78, 0xc0, 0x0c, 0x00, 0x0f, 0x00,
		0x01,

		// Question 3 (com.)
		0xc0, 0x13, 0x00, 0x01, 0x00, 0x01,
	}

	actualBytes, err := MarshalPacket(packet)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(actualBytes, expectedBytes) {
		entries := utils.Diff(actualBytes, expectedBytes)
		t.Fatal(entries.String())
	}
}

func TestUnmarshalQueryPacketOneQuestion(t *testing.T) {
	bytes := []byte{
		// Header
		0x12, 0x34, 0x01, 0x20, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		// Question 1 (google.com.)
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,
	}

	expectedPackcet := types.Packet{
		Header: types.Header{
			ID:                           0x1234,
			PacketType:                   types.PacketTypeQuery,
			Opcode:                       types.OpcodeQuery,
			AuthoritativeAnswer:          false,
			Truncated:                    false,
			RecursionDesired:             true,
			RecursionAvailable:           false,
			AuthenticData:                true,
			CheckingDisabled:             false,
			ResponseCode:                 types.ResponseCodeNoError,
			QuestionSectionSize:          1,
			AnswerSectionSize:            0,
			AuthorityRecordsSectionSize:  0,
			AdditionalRecordsSectionSize: 0,
		},
		Questions: []types.Question{
			{Domain: "google.com.", Type: types.QuestionTypeA, Class: types.QuestionClassIN},
		},
	}

	actualPacket, err := UnmarshalPacket(bytes)
	if err != nil {
		t.Fatal(err)
	}

	entries := utils.Diff(actualPacket, expectedPackcet)
	if len(entries) > 0 {
		t.Fatal(entries.String())
	}
}

func TestUnmarshalQueryPacketTwoQuestions(t *testing.T) {
	bytes := []byte{
		// Header
		0x12, 0x34, 0x01, 0x20, 0x00, 0x02, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		// Question 1 (google.com.)
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,

		// Question 2 (mx.google.com.)
		0x02, 0x6d, 0x78, 0xc0, 0x0c, 0x00, 0x0f, 0x00,
		0x01,
	}

	expectedPackcet := types.Packet{
		Header: types.Header{
			ID:                           0x1234,
			PacketType:                   types.PacketTypeQuery,
			Opcode:                       types.OpcodeQuery,
			AuthoritativeAnswer:          false,
			Truncated:                    false,
			RecursionDesired:             true,
			RecursionAvailable:           false,
			AuthenticData:                true,
			CheckingDisabled:             false,
			ResponseCode:                 types.ResponseCodeNoError,
			QuestionSectionSize:          2,
			AnswerSectionSize:            0,
			AuthorityRecordsSectionSize:  0,
			AdditionalRecordsSectionSize: 0,
		},
		Questions: []types.Question{
			{Domain: "google.com.", Type: types.QuestionTypeA, Class: types.QuestionClassIN},
			{Domain: "mx.google.com.", Type: types.QuestionTypeMX, Class: types.QuestionClassIN},
		},
	}

	actualPacket, err := UnmarshalPacket(bytes)
	if err != nil {
		t.Fatal(err)
	}

	entries := utils.Diff(actualPacket, expectedPackcet)
	if len(entries) > 0 {
		t.Fatal(entries.String())
	}
}

func TestUnmarshalQueryPacketThreeQuestions(t *testing.T) {
	bytes := []byte{
		// Header
		0x12, 0x34, 0x01, 0x20, 0x00, 0x03, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,

		// Question 1 (google.com.)
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03,
		0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,

		// Question 2 (mx.google.com.)
		0x02, 0x6d, 0x78, 0xc0, 0x0c, 0x00, 0x0f, 0x00,
		0x01,

		// Question 3 (com.)
		0xc0, 0x13, 0x00, 0x01, 0x00, 0x01,
	}

	expectedPackcet := types.Packet{
		Header: types.Header{
			ID:                           0x1234,
			PacketType:                   types.PacketTypeQuery,
			Opcode:                       types.OpcodeQuery,
			AuthoritativeAnswer:          false,
			Truncated:                    false,
			RecursionDesired:             true,
			RecursionAvailable:           false,
			AuthenticData:                true,
			CheckingDisabled:             false,
			ResponseCode:                 types.ResponseCodeNoError,
			QuestionSectionSize:          3,
			AnswerSectionSize:            0,
			AuthorityRecordsSectionSize:  0,
			AdditionalRecordsSectionSize: 0,
		},
		Questions: []types.Question{
			{Domain: "google.com.", Type: types.QuestionTypeA, Class: types.QuestionClassIN},
			{Domain: "mx.google.com.", Type: types.QuestionTypeMX, Class: types.QuestionClassIN},
			{Domain: "com.", Type: types.QuestionTypeA, Class: types.QuestionClassIN},
		},
	}

	actualPacket, err := UnmarshalPacket(bytes)
	if err != nil {
		t.Fatal(err)
	}

	entries := utils.Diff(actualPacket, expectedPackcet)
	if len(entries) > 0 {
		t.Fatal(entries.String())
	}
}
