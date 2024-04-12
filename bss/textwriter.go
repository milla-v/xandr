package bss

import (
	"bufio"
	"errors"
	"io"

	"github.com/milla-v/xandr/bss/xgen"
)

type DataFormat string

const (
	FormatText DataFormat = "text"
	FormatAvro DataFormat = "avro"
)

// SegmentDataFormatter provides writing user-segments data to a writer stream according Legacy BSS file format
// or Avro format described on https://learn.microsoft.com/en-us/xandr/bidders/uploading-segment-data-using-bss
type SegmentDataFormatter struct {
	format      DataFormat
	w           *bufio.Writer
	textEncoder *xgen.TextEncoder
}

// NewSegmentDataFormatter creates new BSS text SegmentDataFormatter.
func NewSegmentDataFormatter(w io.Writer, format DataFormat, params *xgen.TextEncoderParameters) (*SegmentDataFormatter, error) {
	df := &SegmentDataFormatter{
		format: format,
		w:      bufio.NewWriter(w),
	}

	if format == FormatText && params == nil {
		return nil, errors.New("text encoder parameters are not specified")
	}

	if format == FormatAvro {
		return nil, errors.New("Avro format is not implemented yet")
	}

	var err error

	df.textEncoder, err = xgen.NewTextEncoder(*params)
	if err != nil {
		return nil, err
	}

	return df, nil
}

// Close flushes buffered text to the writer.
func (df *SegmentDataFormatter) Close() error {
	if err := df.w.Flush(); err != nil {
		return err
	}
	return nil
}

// Append outputs new user record according text format parameters.
func (df *SegmentDataFormatter) Append(ur *xgen.UserRecord) error {
	if df.format == FormatText {
		line, err := df.textEncoder.FormatLine(ur)
		if err != nil {
			return err
		}

		_, err = df.w.WriteString(line + "\n")
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("invalid data format. Only text format is implemented")
}
