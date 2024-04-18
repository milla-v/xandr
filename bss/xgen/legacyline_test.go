package xgen

import (
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	ur := &UserRecord{
		UID:    "12345",
		Domain: "",
		Segments: []Segment{
			{ID: 100},
			{ID: 101},
		},
	}

	enc, err := NewTextEncoder(MinimalFormat)
	if err != nil {
		t.Fatal(err)
	}

	line, err := enc.FormatLine(ur)
	if err != nil {
		t.Fatal(err)
	}

	if line != "12345:100;101" {
		t.Fatal("invalid line:", line)
	}
}

func TestExpired(t *testing.T) {
	ur := &UserRecord{
		UID:    "12345",
		Domain: "",
		Segments: []Segment{
			{ID: 100, Expiration: Expired},
			{ID: 101, Expiration: Expired},
		},
	}

	enc, err := NewTextEncoder(MinimalFormat)
	if err != nil {
		t.Fatal(err)
	}

	line, err := enc.FormatLine(ur)
	if err != nil {
		t.Fatal(err)
	}

	if line != "12345:#100;101" {
		t.Fatal("invalid line:", line)
	}
}

func TestFull(t *testing.T) {
	ur := &UserRecord{
		UID:    "12345",
		Domain: "",
		Segments: []Segment{
			{ID: 100, Expiration: 1440, Value: 123},
			{ID: 101, Expiration: 1440, Value: 123},
		},
	}

	enc, err := NewTextEncoder(FullFormat)
	if err != nil {
		t.Fatal(err)
	}

	line, err := enc.FormatLine(ur)
	if err != nil {
		t.Fatal(err)
	}

	if line != "12345:100:1440:123;101:1440:123" {
		t.Fatal("invalid line:", line)
	}
}

func TestFullIdfa(t *testing.T) {
	ur := &UserRecord{
		UID:    "0000-123123-132123123-3212312",
		Domain: IDFA,
		Segments: []Segment{
			{ID: 100, Expiration: 1440, Value: 123, Timestamp: time.Now().Unix()},
			{ID: 101, Expiration: 1440, Value: 123, Timestamp: time.Now().Unix()},
		},
	}

	enc, err := NewTextEncoder(FullFormat)
	if err != nil {
		t.Fatal(err)
	}

	line, err := enc.FormatLine(ur)
	if err != nil {
		t.Fatal(err)
	}

	if line != "0000-123123-132123123-3212312:100:1440:123;101:1440:123^3" {
		t.Fatal("invalid line:", line)
	}
}

func TestInvalidDomain(t *testing.T) {
	ur := &UserRecord{
		UID:    "0000-123123-132123123-3212312",
		Domain: "aaa",
		Segments: []Segment{
			{ID: 100, Expiration: 1440, Value: 123},
			{ID: 101, Expiration: 1440, Value: 123},
		},
	}

	enc, err := NewTextEncoder(FullFormat)
	if err != nil {
		t.Fatal(err)
	}

	_, err = enc.FormatLine(ur)
	if err == nil {
		t.Fatal("should return error")
	}

	if err.Error() != "invalid domain: aaa" {
		t.Fatal("invalid error message:", err.Error())
	}
}

func TestInvalidSegId(t *testing.T) {
	ur := &UserRecord{
		UID: "12345",
		Segments: []Segment{
			{ID: 100, Expiration: 1440, Value: 123},
			{Expiration: 1440, Value: 123},
		},
	}

	enc, err := NewTextEncoder(FullFormat)
	if err != nil {
		t.Fatal(err)
	}

	_, err = enc.FormatLine(ur)
	if err == nil {
		t.Fatal("should return error")
	}

	if err.Error() != "seg[1].ID is zero" {
		t.Fatal("invalid error message:", err.Error())
	}
}

func TestCodeAndMemberID(t *testing.T) {
	ur := &UserRecord{
		UID: "12345",
		Segments: []Segment{
			{ID: 100, Code: "CodeTest", MemberID: 123, Expiration: 1440, Value: 123},
			{Code: "CodeTest", Expiration: 1440, Value: 123},
		},
	}

	enc, err := NewTextEncoder(FullExternalFormat)
	if err != nil {
		t.Fatal(err)
	}

	_, err = enc.FormatLine(ur)

	t.Log("UR: ", ur)

	if err == nil {
		t.Fatal("should return error")
	}
	if err.Error() != "seg[1].Code is empty" && err.Error() != "seg[1].MemberID is zero" {
		t.Fatal("invalid error message:", err.Error())
	}
}

func TestInvalidExpiration(t *testing.T) {
	ur := &UserRecord{
		UID: "12345",
		Segments: []Segment{
			{ID: 100, Expiration: 181 * 60 * 24, Value: 123},
			{ID: 101, Expiration: 1440, Value: 123},
		},
	}

	enc, err := NewTextEncoder(FullFormat)
	if err != nil {
		t.Fatal(err)
	}

	_, err = enc.FormatLine(ur)
	if err == nil {
		t.Fatal("should return error")
	}

	if err.Error() != "seg[0].Expiration is not in the range [-1, 259200]" {
		t.Fatal("invalid error message:", err.Error())
	}
}

func TestMinimalTextEncoder(t *testing.T) {
	params := TextEncoderParameters{
		Sep1:          ":",
		Sep2:          ";",
		Sep3:          ":",
		Sep4:          "#",
		Sep5:          "^",
		SegmentFields: MinimalFormat.SegmentFields,
	}

	_, err := NewTextEncoder(params)
	if err != nil {
		t.Fatal(err)
	}
}

// to test SegmentsFields: SEG_ID, SEG_CODE, MEMBER_ID, EXPIRATION, VALUE"
func TestFullTextEncoder(t *testing.T) {
	sf := []SegmentFieldName{"SEG_CODE", "MEMBER_ID"}
	t.Log("SF: ", sf)
	params := TextEncoderParameters{
		Sep1:          ":",
		Sep2:          ";",
		Sep3:          ":",
		Sep4:          "#",
		Sep5:          "^",
		SegmentFields: sf}

	_, err := NewTextEncoder(params)
	if err != nil {
		t.Fatal(err)
	}
}
