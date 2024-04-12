package bss

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/milla-v/xandr/bss/xgen"
)

var FullFormat = xgen.FullFormat

func TestSegmentDataFormatter2(t *testing.T) {
	const input = `
UID,SegID
12345,100
12346,102
`

	p := xgen.TextEncoderParameters{
		Sep1:          ":",
		Sep2:          ";",
		Sep3:          ",",
		Sep4:          "#",
		Sep5:          "^",
		SegmentFields: []xgen.SegmentFieldName{xgen.SegIdField},
	}

	var out bytes.Buffer

	w, err := NewSegmentDataFormatter(&out, FormatText, &p)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines[1:] {
		columns := strings.Split(line, ",")
		t.Logf("cols: %+v", columns)
		segID, err := strconv.ParseInt(columns[1], 10, 32)
		log.Println(line)
		if err != nil {
			t.Fatal(err)
		}

		ur := xgen.UserRecord{
			UID: columns[0],
			Segments: []xgen.Segment{
				{ID: 55},
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

	t.Log("generated file:", out.String())
}

func TestSegmentDataFormatter(t *testing.T) {
	const input = `
	UID,SegID,Expiration,Value
	12345,100,1440,123
	12346,101,1440,123`

	params := xgen.TextEncoderParameters(FullFormat)

	var out bytes.Buffer

	w, err := NewSegmentDataFormatter(&out, FormatText, &params)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines[1:] {
		columns := strings.Split(line, ",")
		segID, err := strconv.ParseInt(columns[1], 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		expiration, err := strconv.ParseInt(columns[2], 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		value, err := strconv.ParseInt(columns[3], 10, 32)
		if err != nil {
			t.Fatal(err)
		}

		ur := xgen.UserRecord{
			UID: columns[0],
			Segments: []xgen.Segment{
				xgen.Segment{
					ID:         int32(segID),
					Expiration: int32(expiration),
					Value:      int32(value),
				},
			},
		}
		if err := w.Append(&ur); err != nil {
			t.Fatal(err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	t.Log("generated file:", out.String())
}
