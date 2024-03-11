package bss

import (
	"io"

	"github.com/milla-v/xandr/bss/xgen"
)

// TextFileWriter provides writing user-segments data to a file according Legacy BSS file format
// described on https://learn.microsoft.com/en-us/xandr/bidders/legacy-bss-file-format
type TextFileWriter struct {
	w         io.WriteCloser
	formatter *xgen.TextFormatter
}

func NewTextFileWriter(fname string, p TextFileParameters) *TextFileFormatter {
	return &TextFileWriter{}
}

func (w *TextFileWriter) Close() error {
	return nil
}

func (w *TextFileWriter) Append(ur *UserRecord) error {
}
