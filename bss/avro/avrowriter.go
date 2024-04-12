package avro

import (
	"fmt"
	"io"
	"strconv"

	"github.com/linkedin/goavro/v2"
	"github.com/milla-v/xandr/bss/xgen"

	_ "embed"
)

type UserRecord = xgen.UserRecord

type AvroWriter struct {
	ocfWriter *goavro.OCFWriter
}

// xandr schema from https://learn.microsoft.com/en-us/xandr/bidders/bss-avro-file-format
//
//go:embed xandr_schema.avsc
var xandrSchema string

// NewAvroWriter creates avro writer for generating data in Xandr BSS avro uploading format.
func NewAvroWriter(w io.Writer) (*AvroWriter, error) {
	ocfConfig := goavro.OCFConfig{
		Schema: xandrSchema,
		W:      w,
	}

	ocfWriter, err := goavro.NewOCFWriter(ocfConfig)
	if err != nil {
		return nil, err
	}

	writer := &AvroWriter{
		ocfWriter: ocfWriter,
	}

	return writer, nil
}

func newXandrID(id string) (map[string]interface{}, error) {
	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("xandr is should be int64. error: %w", err)
	}

	return map[string]interface{}{
			"long": n,
		},
		nil
}

func newDeviceID(id string, domain xgen.Domain) (map[string]interface{}, error) {
	// TODO: validate UUID

	return map[string]interface{}{
			"device_id": map[string]interface{}{
				"id":     id,
				"domain": string(domain),
			},
		},
		nil
}

func newSegments(segments []xgen.Segment) ([]map[string]interface{}, error) {
	var list []map[string]interface{}

	for _, segment := range segments {
		// TODO: validate segment fields
		item := map[string]interface{}{
			"id":         segment.ID,
			"value":      segment.Value,
			"expiration": segment.Expiration,
			"code":       segment.Code,
			"member_id":  segment.MemberID,
		}
		list = append(list, item)
	}

	return list, nil
}

func (w *AvroWriter) Append(users []*UserRecord) error {
	var err error
	var records []interface{}

	for _, user := range users {
		var uid map[string]interface{}

		switch user.Domain {
		case xgen.XandrID:
			uid, err = newXandrID(user.UID)
		case xgen.AAID, xgen.IDFA:
			uid, err = newDeviceID(user.UID, user.Domain)
		}

		if err != nil {
			return err
		}

		segments, err := newSegments(user.Segments)
		if err != nil {
			return err
		}

		record := map[string]interface{}{
			"uid":      uid,
			"segments": segments,
		}

		records = append(records, record)
	}

	if err := w.ocfWriter.Append(records); err != nil {
		return err
	}

	return nil
}
