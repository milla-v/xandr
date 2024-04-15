package avro

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/linkedin/goavro/v2"
	"github.com/milla-v/xandr/bss/xgen"
)

func TestAvroWriter(t *testing.T) {
	var out bytes.Buffer

	wr, err := NewAvroWriter(&out)
	if err != nil {
		t.Fatal(err)
	}

	var users []*UserRecord

	user := &UserRecord{
		UID: "12345",
		Segments: []xgen.Segment{
			{ID: 100, Expiration: 180 * 60 * 24, Value: 123},
			{ID: 101, Expiration: 1440, Value: 123},
		},
	}

	users = append(users, user)

	if err := wr.Append(users); err != nil {
		t.Fatal(err)
	}

	ocfr, err := goavro.NewOCFReader(&out)
	if err != nil {
		t.Fatal(err)
	}

	const expectedText = "map[segments:[map[code: expiration:259200 id:100 member_id:0 timestamp:0 value:123] map[code: expiration:1440 id:101 member_id:0 timestamp:0 value:123]] uid:map[long:12345]]"

	for ocfr.Scan() {
		value, err := ocfr.Read()
		if err != nil {
			t.Fatal(err)
		}
		text := fmt.Sprintf("%+v", value)
		if text != expectedText {
			t.Fatal("\nexpected:", expectedText, "\nactual  :", text)
		}
	}

	if err = ocfr.Err(); err != nil {
		t.Fatal(err)
	}
}
