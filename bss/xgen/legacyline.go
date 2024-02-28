package xgen

import (
	"fmt"
	"io"
	"log"
	"strings"
)

const legacyLineTemplate = "{UID}{SEP_1}{SEGMENTS_TO_ADD}{SEP_4}{SEGMENTS_TO_REMOVE}{SEP_5}{DOMAIN}"

type SegmentFieldName string

const (
	SegIdField      SegmentFieldName = "SEG_ID"
	SegCodeField    SegmentFieldName = "SEG_CODE"
	MemberIdField   SegmentFieldName = "MEMBER_ID"
	ExpirationField SegmentFieldName = "EXPIRATION"
	ValueField      SegmentFieldName = "VALUE"
)

type TextFormater struct {
	Sep1          string // Separator after UID
	Sep2          string // Separator beetween segments
	Sep3          string // Separator between segment fields
	Sep4          string // Separator between segment additions block and segment removals block
	Sep5          string // Separator before domain
	SegmentFields []SegmentFieldName
}

var MinimalFormat = TextFormater{
	Sep1:          ":",
	Sep2:          ";",
	Sep3:          ":",
	Sep4:          "#",
	Sep5:          "^",
	SegmentFields: []SegmentFieldName{SegIdField},
}

var FullFormat = TextFormater{
	Sep1: ":",
	Sep2: ";",
	Sep3: ":",
	Sep4: "#",
	Sep5: "^",
	SegmentFields: []SegmentFieldName{
		SegIdField,
		ExpirationField,
		ValueField,
	},
}

var FullExternalFormat = TextFormater{
	Sep1: ":",
	Sep2: ";",
	Sep3: ":",
	Sep4: "#",
	Sep5: "^",
	SegmentFields: []SegmentFieldName{
		SegCodeField,
		MemberIdField,
		ExpirationField,
		ValueField,
	},
}

func (tf *TextFormater) FormatLine(ur *UserRecord) (string, error) {
	if _, ok := domains[ur.Domain]; !ok {
		return "", fmt.Errorf("invalid domain: %s", ur.Domain)
	}

	var b strings.Builder

	b.WriteString(ur.UID)
	b.WriteString(tf.Sep1)

	var adds []Segment
	var rems []Segment

	for _, seg := range ur.Segments {
		if seg.Expiration == Expired {
			rems = append(rems, seg)
		} else {
			adds = append(adds, seg)
		}
	}

	if err := genSegments(&b, tf, adds); err != nil {
		return "", err
	}

	if len(rems) > 0 {
		b.WriteString(tf.Sep4)
		if err := genSegments(&b, tf, rems); err != nil {
			return "", nil
		}
	}

	if ur.Domain != "" {
		b.WriteString(tf.Sep5)
		b.WriteString(string(ur.Domain))
	}

	return b.String(), nil
}

func genSegments(w io.Writer, tf *TextFormater, list []Segment) error {
	const maxValue = 2147483647

	for i, seg := range list {
		for j, sf := range tf.SegmentFields {
			switch sf {
			case SegIdField:
				if seg.ID == 0 {
					return fmt.Errorf("seg[%d].ID is zero", i)
				}
				fmt.Fprintf(w, "%d", seg.ID)
			case SegCodeField:
				if seg.Code == "" {
					return fmt.Errorf("seg[%d].Code is empty", i)
				}
				io.WriteString(w, seg.Code)
			case MemberIdField:
				if seg.MemberID == 0 {
					return fmt.Errorf("seg[%d].MemberID is zero", i)
				}
				fmt.Fprintf(w, "%d", seg.MemberID)
			case ExpirationField:
				if seg.Expiration < -1 || seg.Expiration > MaxExpiration {
					return fmt.Errorf("seg[%d].Expiration is not in the range [-1, %d]", i, MaxExpiration)
				}
				fmt.Fprintf(w, "%d", seg.Expiration)
			case ValueField:
				if seg.Value < 0 || seg.Value > maxValue {
					return fmt.Errorf("seg[%d].Value is not in the range [-1, %d]", i, maxValue)
				}
				fmt.Fprintf(w, "%d", seg.Value)
			}

			if j < len(tf.SegmentFields)-1 {
				io.WriteString(w, tf.Sep3)
			}
		}

		if i < len(list)-1 {
			io.WriteString(w, tf.Sep2)
		}
	}

	return nil
}
