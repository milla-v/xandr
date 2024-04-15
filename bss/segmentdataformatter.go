package bss

import (
	"bufio"
	"errors"
	"io"

	"github.com/milla-v/xandr/bss/avro"
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
	avroEncoder *avro.AvroWriter
}

// NewSegmentDataFormatter creates new BSS text SegmentDataFormatter.
func NewSegmentDataFormatter(w io.Writer, format DataFormat, params *xgen.TextEncoderParameters) (*SegmentDataFormatter, error) {
	var err error

	df := &SegmentDataFormatter{
		format: format,
		w:      bufio.NewWriter(w),
	}

	if format == FormatText && params == nil {
		return nil, errors.New("text encoder parameters are not specified")
	}

	if format == FormatAvro {
		df.avroEncoder, err = avro.NewAvroWriter(df.w)
		if err != nil {
			return nil, err
		}
		return df, err
	}

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

// Append outputs users records to a writer.
func (df *SegmentDataFormatter) Append(users []*xgen.UserRecord) error {
	if df.format == FormatAvro {
		return df.avroEncoder.Append(users)
	}

	for _, user := range users {
		line, err := df.textEncoder.FormatLine(user)
		if err != nil {
			return err
		}

		_, err = df.w.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
