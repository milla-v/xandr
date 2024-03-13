package bss

import (
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/milla-v/xandr/bss/xgen"
)

var MinimalFormat = xgen.MinimalFormat

/*
ur := &UserRecord{
                UID:    "12345",
                Domain: "",
                Segments: []Segment{
                        {ID: 100, Expiration: Expired},
                        {ID: 101, Expiration: Expired},
                },
        }
*/

func TestTextWriter2(t *testing.T) {
	const input = `UID,SegID
12345,100
12346,100
`

	p := TextFileParameters{
		Sep1:          ":",
		Sep2:          ";",
		Sep3:          ",",
		Sep4:          "#",
		Sep5:          "^",
		SegmentFields: []xgen.SegmentFieldName{xgen.SegIdField},
	}

	w, err := NewTextFileWriter("1.txt", p)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines[1:] {
		columns := strings.Split(line, ",")
		t.Logf("cols: %+v", columns)
		segID, err := strconv.ParseInt(columns[1], 10, 32)
		if err != nil {
			t.Fatal(err)
		}

		ur := UserRecord{
			UID: columns[0],
			Segments: []Segment{
				{ID: int32(segID)},
			},
		}

		if err := w.Append(&ur); err != nil {
			t.Fatal(err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	buf, err := os.ReadFile("1.txt")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("generated file:", string(buf))
}

func TestTextWriter(t *testing.T) {

	log.Println("Minimal Format: ", MinimalFormat)
	params := TextFileParameters(MinimalFormat)
	log.Println("Params: ", params)
	w, err := NewTextFileWriter("1.txt", params)
	if err != nil {
		t.Fatal(err)
	}

	var users []UserRecord

	for _, u := range users {
		if err := w.Append(&u); err != nil {
			t.Fatal(err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}
