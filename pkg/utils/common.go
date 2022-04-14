package utils

import (
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

// FIXME

func StrfmtToUUID(tmp strfmt.UUID) (uuid.UUID, error) {
	return uuid.Parse(tmp.String())
}

func UUIDToStrfmt(tmp uuid.UUID) strfmt.UUID {
	return strfmt.UUID(tmp.String())

}
