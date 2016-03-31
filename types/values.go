package types

import (
	"fmt"
	"time"
)

// ID is defined as a string pointer in order to check for nil in the context of
// relations. NOTE this may not be best practice.
type ID *string

type Timestamp struct {
	time.Time
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.Time.Format(time.RFC1123))
	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(raw []byte) (err error) {
	str := string(raw)
	noQuotes := str[1 : len(str)-1]
	t.Time, err = time.Parse(time.RFC1123, noQuotes)
	return err
}

func NewTimestamp() *Timestamp {
	return &Timestamp{time.Now().UTC()}
}

// CPU / RAM resource values

// TODO these could probably just be uint types instead of Structs.

type BytesValue struct {
	bytes uint
}

const (
	bytesKiB = 1024
	bytesMiB = bytesKiB * 1024
)

func BytesFromMiB(mib uint) *BytesValue {
	return &BytesValue{mib * bytesMiB}
}

// func (v *BytesValue) kibibytes() uint {
//   return v.bytes / kib
// }

func (v *BytesValue) mebibytes() uint {
	return v.bytes / bytesMiB
}

func (v *BytesValue) ToKubeMebibytes() string {
	return fmt.Sprintf("%dMi", v.mebibytes())
}

type CoresValue struct {
	cores uint
}

func CoresFromMillicores(millicores uint) *CoresValue {
	return &CoresValue{millicores / 1000}
}

func (v *CoresValue) millicores() uint {
	return v.cores * 1000
}

func (v *CoresValue) ToKubeMillicores() string {
	return fmt.Sprintf("%dm", v.millicores())
}
