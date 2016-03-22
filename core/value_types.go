package core

import "fmt"

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
