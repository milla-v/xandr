package bss

import (
	"bufio"
	"io"
	"os"

	"github.com/milla-v/xandr/bss/xgen"
)

// TextFileWriter provides writing user-segments data to a file according Legacy BSS file format
// described on https://learn.microsoft.com/en-us/xandr/bidders/legacy-bss-file-format
type TextFileWriter struct {
	file      io.Closer
	w         *bufio.Writer
	formatter *xgen.TextFormatter
}

func NewTextFileWriter(fname string, p TextFileParameters) (*TextFileFormatter, error) {
	createdFile, err := os.Create(fname)
	if err != nil {
		return err
	}

	tw := &TextFileWriter{
		w:    bufio.NewWriter(createdFile),
		file: createdFile,
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

func (w *TextFileWriter) Append(ur *UserRecord) error {

}
