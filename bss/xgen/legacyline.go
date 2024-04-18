package xgen

import (
	"fmt"
	"io"
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
	TimestampField  SegmentFieldName = "TIMESTAMP"
	NotAllowed                       = "[](){}$\\/|?*+-"
)

type TextEncoderParameters struct {
	Sep1          string // Separator after UID
	Sep2          string // Separator beetween segments
	Sep3          string // Separator between segment fields
	Sep4          string // Separator between segment additions block and segment removals block
	Sep5          string // Separator before domain
	SegmentFields []SegmentFieldName
}

type TextEncoder struct {
	parameters TextEncoderParameters
}

var MinimalFormat = TextEncoderParameters{
	Sep1:          ":",
	Sep2:          ";",
	Sep3:          ":",
	Sep4:          "#",
	Sep5:          "^",
	SegmentFields: []SegmentFieldName{SegIdField},
}

var FullFormat = TextEncoderParameters{
	Sep1: ":",
	Sep2: ";",
	Sep3: ":",
	Sep4: "#",
	Sep5: "^",
	SegmentFields: []SegmentFieldName{
		SegIdField,
		ExpirationField,
		ValueField,
		TimestampField,
	},
}

var FullExternalFormat = TextEncoderParameters{
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
		TimestampField,
	},
}

func (tf *TextEncoder) FormatLine(ur *UserRecord) (string, error) {
	if _, ok := domains[ur.Domain]; !ok {
		return "", fmt.Errorf("invalid domain: %s", ur.Domain)
	}

	var b strings.Builder

	b.WriteString(ur.UID)
	b.WriteString(tf.parameters.Sep1)

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
		b.WriteString(tf.parameters.Sep4)
		if err := genSegments(&b, tf, rems); err != nil {
			return "", nil
		}
	}

	if ur.Domain != "" {
		b.WriteString(tf.parameters.Sep5)
		b.WriteString(string(ur.Domain))
	}

	return b.String(), nil
}

func genSegments(w io.Writer, tf *TextEncoder, list []Segment) error {
	const maxValue = 2147483647

	for i, seg := range list {
		for j, sf := range tf.parameters.SegmentFields {
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
			case TimestampField:
				fmt.Fprintf(w, "%d", seg.Timestamp)
			}

			if j < len(tf.parameters.SegmentFields)-1 {
				io.WriteString(w, tf.parameters.Sep3)
			}
		}

		if i < len(list)-1 {
			io.WriteString(w, tf.parameters.Sep2)
		}
	}

	return nil
}

func NewTextEncoder(parameters TextEncoderParameters) (*TextEncoder, error) {
	sp := []string{
		parameters.Sep1,
		parameters.Sep2,
		parameters.Sep3,
		parameters.Sep4,
		parameters.Sep5,
	}

	var tf TextEncoder
	var err error

	if err = checkSegments(parameters.SegmentFields); err != nil {
		return nil, err
	}

	if err = checkSeparators(sp); err != nil {
		return nil, err
	}

	tf.parameters.Sep1 = parameters.Sep1
	tf.parameters.Sep2 = parameters.Sep2
	tf.parameters.Sep3 = parameters.Sep3
	tf.parameters.Sep4 = parameters.Sep4
	tf.parameters.Sep5 = parameters.Sep5
	tf.parameters.SegmentFields = parameters.SegmentFields

	return &tf, nil
}

func checkSegments(sf []SegmentFieldName) error {
	var err error
	var segIDfound bool
	var segCodeFound bool
	var memberIDfound bool

	//start check segmentFields
	for _, s := range sf {
		if s == SegIdField {
			segIDfound = true
		}
		if s == SegCodeField {
			segCodeFound = true
		}
		if s == MemberIdField {
			memberIDfound = true
		}
	}
	//check if at least  SEG_ID or SEG_CODE was choosen
	if segIDfound == false && segCodeFound == false {
		return fmt.Errorf("Choose at least  SEG_ID or SEG_CODE")
	}
	// check if SEG_CODE or SEG_ID included but not both.
	if segIDfound == true && segCodeFound == true {
		return fmt.Errorf("You may include SEG_CODE or SEG_ID but not both")
	}
	// if SEG_CODE present, MEMBER_ID should be choosen too
	if segCodeFound == true && memberIDfound == false {
		return fmt.Errorf("If SEG_CODE present, MEMBER_ID should be choosen too")
	}

	return err
}

func checkSeparators(sp []string) error {
	for i, s := range sp {
		if len(s) != 1 && s != "\t" && s != " " {
			return fmt.Errorf("sep%d should be a single character", i+1)
		}
		if strings.ContainsAny(s, NotAllowed) {
			return fmt.Errorf("sep%d: symbols "+NotAllowed+" are not allowed as a separators", i+1)
		}
	}
	return nil
}
