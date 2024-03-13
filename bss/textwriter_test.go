package bss

import (
	"log"
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
