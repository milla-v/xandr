package xgen

import (
	"log"
	"testing"
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

	line, err := MinimalFormat.FormatLine(ur)
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

	line, err := MinimalFormat.FormatLine(ur)
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

	line, err := FullFormat.FormatLine(ur)
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
			{ID: 100, Expiration: 1440, Value: 123},
			{ID: 101, Expiration: 1440, Value: 123},
		},
	}

	line, err := FullFormat.FormatLine(ur)
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

	_, err := FullFormat.FormatLine(ur)
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

	_, err := FullFormat.FormatLine(ur)
	if err == nil {
		t.Fatal("should return error")
	}

	if err.Error() != "seg[1].ID is zero" {
		t.Fatal("invalid error message:", err.Error())
	}
}

func TestFullFormatId(t *testing.T) {
	ur := &UserRecord{
		UID: "12345",
		Segments: []Segment{
			{ID: 100, Code: "CodeTest", MemberID: 123, Expiration: 1440, Value: 123},
		},
	}

	log.Println(ur)
	_, err := FullExternalFormat.FormatLine(ur)

	if err == nil {
		t.Fatal("should return error")
	}
	if err.Error() != "seg[1].Code is empty" || err.Error() != "seg[1].MemberID is zero" {
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

	_, err := FullFormat.FormatLine(ur)
	if err == nil {
		t.Fatal("should return error")
	}

	if err.Error() != "seg[0].Expiration is not in the range [-1, 259200]" {
		t.Fatal("invalid error message:", err.Error())
	}
}

func TestMinimalTextFormater(t *testing.T) {
	min := TextFormater{Sep1: ":", Sep2: ";", Sep3: ":", Sep4: "#", Sep5: "^", SegmentFields: MinimalFormat.SegmentFields}
	if _, err := NewTextFormater(min); err != nil {
		t.Fatal("TestMinimalTextFormater: ", err)
	}
}

// to test SegmentsFields: SEG_ID, SEG_CODE, MEMBER_ID, EXPIRATION, VALUE"
func TestFullTextFormater(t *testing.T) {
	sf := []SegmentFieldName{"SEG_CODE", "MEMBER_ID"}
	log.Println("SF: ", sf)
	full := TextFormater{Sep1: ":", Sep2: ";", Sep3: ":", Sep4: "#", Sep5: "^", SegmentFields: sf}
	if _, err := NewTextFormater(full); err != nil {
		t.Fatal("TestFullTextFormater: ", err)
	}
}
