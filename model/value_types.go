package model

import "fmt"

type BytesValue struct {
	bytes int
}

const (
	bytesKiB = 1024
	bytesMiB = bytesKiB * 1024
)

func BytesFromMiB(mib int) *BytesValue {
	return &BytesValue{mib / bytesMiB}
}

// func (v *BytesValue) kibibytes() int {
//   return v.bytes / kib
// }

func (v *BytesValue) mebibytes() int {
	return v.bytes / bytesMiB
}

func (v *BytesValue) ToKubeMebibytes() string {
	return fmt.Sprintf("%dMi", v.mebibytes())
}

type CoresValue struct {
	cores int
}

func (v *CoresValue) millicores() int {
	return v.cores * 1000
}

func (v *CoresValue) ToKubeMillicores() string {
	return fmt.Sprintf("%dm", v.millicores())
}
