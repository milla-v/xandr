package bss

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/linkedin/goavro"
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

	var users []*xgen.UserRecord

	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines[1:] {
		columns := strings.Split(line, ",")
		t.Logf("cols: %+v", columns)
		segID, err := strconv.ParseInt(columns[1], 10, 32)
		log.Println(line)
		if err != nil {
			t.Fatal(err)
		}

		user := &xgen.UserRecord{
			UID: columns[0],
			Segments: []xgen.Segment{
				{ID: 55},
				{ID: int32(segID)},
			},
		}

		users = append(users, user)
	}

	if err := w.Append(users); err != nil {
		t.Fatal(err)
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

	var users []*xgen.UserRecord

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

		user := &xgen.UserRecord{
			UID: columns[0],
			Segments: []xgen.Segment{
				xgen.Segment{
					ID:         int32(segID),
					Expiration: int32(expiration),
					Value:      int32(value),
				},
			},
		}
		users = append(users, user)
	}

	if err := w.Append(users); err != nil {
		t.Fatal(err)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	t.Log("generated file:", out.String())
}

func TestSegmentDataFormatterWithAvroFormat(t *testing.T) {
	const input = `
UID,SegID,Expiration,Value
12345,100,1440,123
12346,101,1440,123`

	const expectedResult = `
map[segments:[map[code: expiration:1440 id:100 member_id:0 timestamp:0 value:123]] uid:map[long:12345]]
map[segments:[map[code: expiration:1440 id:101 member_id:0 timestamp:0 value:123]] uid:map[long:12346]]
`

	var out bytes.Buffer

	w, err := NewSegmentDataFormatter(&out, FormatAvro, nil)
	if err != nil {
		t.Fatal(err)
	}

	var users []*xgen.UserRecord

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

		user := &xgen.UserRecord{
			UID: columns[0],
			Segments: []xgen.Segment{
				xgen.Segment{
					ID:         int32(segID),
					Expiration: int32(expiration),
					Value:      int32(value),
				},
			},
		}

		users = append(users, user)
	}

	if err := w.Append(users); err != nil {
		t.Fatal(err)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	// check result

	var result = "\n"

	ocfr, err := goavro.NewOCFReader(&out)
	if err != nil {
		t.Fatal(err)
	}

	for ocfr.Scan() {
		value, err := ocfr.Read()
		if err != nil {
			t.Fatal(err)
		}
		result += fmt.Sprintf("%+v\n", value)
	}

	if err = ocfr.Err(); err != nil {
		t.Fatal(err)
	}

	if result != expectedResult {
		t.Fatal("\nexpected:", expectedResult, "\nactual  :", result)
	}
}
