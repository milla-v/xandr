package bss

import (
	"bufio"
	"io"

	"github.com/milla-v/xandr/bss/xgen"
)

// TextFileWriter provides writing user-segments data to a file according Legacy BSS file format
// described on https://learn.microsoft.com/en-us/xandr/bidders/legacy-bss-file-format
type TextFileWriter struct {
	w       *bufio.Writer
	encoder *xgen.TextEncoder
}

func NewTextFileWriter(w io.Writer, params xgen.TextEncoderParameters) (*TextFileWriter, error) {
	tw := &TextFileWriter{
		w: bufio.NewWriter(w),
	}

	var err error

	tw.encoder, err = xgen.NewTextEncoder(params)
	if err != nil {
		return nil, err
	}

	return tw, nil
}

func (tw *TextFileWriter) Close() error {
	if err := tw.w.Flush(); err != nil {
		return err
	}
	return nil
}

func (w *TextFileWriter) Append(ur *xgen.UserRecord) error {
	line, err := w.encoder.FormatLine(ur)
	if err != nil {
		return err
	}

	_, err = w.w.WriteString(line + "\n")
	if err != nil {
		return err
	}

	return nil
}
