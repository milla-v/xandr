package xgen

import "testing"

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
