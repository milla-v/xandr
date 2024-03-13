package bss

import (
	"github.com/milla-v/xandr/bss/xgen"
)

const (
	Expired           = -1            // Set Segment.Expiration field to remove user from the segment
	DefaultExpiration = 0             // Segment expiration will be set to member's default
	MaxExpiration     = 180 * 24 * 60 // 180 days in minutes
)

type Segment struct {
	ID         int32
	Code       string
	MemberID   int32
	Expiration int32
	Value      int32
}

type Domain string

const (
	XandrID Domain = ""
	IDFA    Domain = "3"
	AAID    Domain = "8"
)

var domains = map[Domain]string{
	XandrID: "",
	IDFA:    "3",
	AAID:    "8",
}

type UserRecord struct {
	UID    string
	Domain Domain

	Segments []Segment
}

type TextFileParameters = xgen.TextEncoder

/*type TextFileParameter struct {
	Sep1          string // Separator after UID
	Sep2          string // Separator beetween segments
	Sep3          string // Separator between segment fields
	Sep4          string // Separator between segment additions block and segment removals block
	Sep5          string // Separator before domain
	SegmentFields []SegmentFieldName
}*/
