package bss

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/milla-v/xandr/bss/xgen"
)

type TextFileParameters = xgen.TextEncoder

// TextFileWriter provides writing user-segments data to a file according Legacy BSS file format
// described on https://learn.microsoft.com/en-us/xandr/bidders/legacy-bss-file-format
type TextFileWriter struct {
	file    io.Closer
	w       *bufio.Writer
	encoder *xgen.TextEncoder
}

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

func NewTextFileWriter(fname string, p TextFileParameters) (*TextFileWriter, error) {
	createdFile, err := os.Create(fname)
	if err != nil {
		return nil, err
	}

	tw := &TextFileWriter{
		w:    bufio.NewWriter(createdFile),
		file: createdFile,
	}

	// create new textencoder, NewTextEncoder should check separators and seg fields.
	tw.encoder, err = xgen.NewTextEncoder(p)
	if err != nil {
		log.Println("NewCheckEncoder result: ", err)
		return nil, err
	} else {
		log.Println("Can write data to a file")
		log.Println("p :", p.SegmentFields)
	}

	return tw, nil
}

func (tw *TextFileWriter) Close() error {
	if err := tw.w.Flush(); err != nil {
		return err
	}
	if err := tw.file.Close(); err != nil {
		return err
	}
	return nil
}

func (w *TextFileWriter) Append(ur *xgen.UserRecord) error {
	var err error
	line, err := w.encoder.FormatLine(ur)
	log.Println("line:", line)

	// use text encoder FormatLine to produce a formatetd line
	// write line to the file

	return err
}
