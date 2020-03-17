/*--------------------------------------------------------*\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: https://hprose.com                     |
|                                                          |
| io/encoding/slice_encoder.go                             |
|                                                          |
| LastModified: Mar 17, 2020                               |
| Author: Ma Bingyao <andot@hprose.com>                    |
|                                                          |
\*________________________________________________________*/

package encoding

import (
	"reflect"

	"github.com/hprose/hprose-golang/v3/io"
	"github.com/modern-go/reflect2"
)

// SliceEncoder is the implementation of ValueEncoder for *slice.
type SliceEncoder struct{}

var sliceEncoder SliceEncoder

// Encode writes the hprose encoding of v to stream
// if v is already written to stream, it will writes it as reference
func (valenc SliceEncoder) Encode(enc *Encoder, v interface{}) (err error) {
	var ok bool
	if ok, err = enc.WriteReference(v); !ok && err == nil {
		err = valenc.Write(enc, v)
	}
	return
}

// Write writes the hprose encoding of v to stream
// if v is already written to stream, it will writes it as value
func (SliceEncoder) Write(enc *Encoder, v interface{}) (err error) {
	enc.SetReference(v)
	return writeSlice(enc, reflect.ValueOf(v).Elem().Interface())
}

// WriteSlice to encoder
func WriteSlice(enc *Encoder, v interface{}) (err error) {
	enc.AddReferenceCount(1)
	return writeSlice(enc, v)
}

var emptySlice = []byte{io.TagList, io.TagOpenbrace, io.TagClosebrace}

func writeSlice(enc *Encoder, v interface{}) (err error) {
	writer := enc.Writer
	if bytes, ok := v.([]byte); ok {
		return writeBytes(writer, bytes)
	}
	count := (*reflect.SliceHeader)(reflect2.PtrOf(v)).Len
	if count == 0 {
		_, err = writer.Write(emptySlice)
		return
	}
	if err = WriteHead(writer, count, io.TagList); err == nil {
		if err = writeSliceBody(enc, v); err == nil {
			err = WriteFoot(writer)
		}
	}
	return
}

func writeSliceBody(enc *Encoder, v interface{}) error {
	switch v := v.(type) {
	case []uint16:
		return writeUint16SliceBody(enc.Writer, v)
	case []uint32:
		return writeUint32SliceBody(enc.Writer, v)
	case []uint64:
		return writeUint64SliceBody(enc.Writer, v)
	case []uint:
		return writeUintSliceBody(enc.Writer, v)
	case []int8:
		return writeInt8SliceBody(enc.Writer, v)
	case []int16:
		return writeInt16SliceBody(enc.Writer, v)
	case []int32:
		return writeInt32SliceBody(enc.Writer, v)
	case []int64:
		return writeInt64SliceBody(enc.Writer, v)
	case []int:
		return writeIntSliceBody(enc.Writer, v)
	case []bool:
		return writeBoolSliceBody(enc.Writer, v)
	case []float32:
		return writeFloat32SliceBody(enc.Writer, v)
	case []float64:
		return writeFloat64SliceBody(enc.Writer, v)
	case []complex64:
		return writeComplex64SliceBody(enc, v)
	case []complex128:
		return writeComplex128SliceBody(enc, v)
	case []string:
		return writeStringSliceBody(enc, v)
	case []interface{}:
		return writeInterfaceSliceBody(enc, v)
	default:
		return writeOtherSliceBody(enc, v)
	}
}

func writeInt8SliceBody(writer io.Writer, slice []int8) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteInt8(writer, slice[i])
		}
	}
	return
}

func writeInt16SliceBody(writer io.Writer, slice []int16) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteInt16(writer, slice[i])
		}
	}
	return
}

func writeInt32SliceBody(writer io.Writer, slice []int32) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteInt32(writer, slice[i])
		}
	}
	return
}

func writeInt64SliceBody(writer io.Writer, slice []int64) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteInt64(writer, slice[i])
		}
	}
	return
}

func writeIntSliceBody(writer io.Writer, slice []int) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteInt(writer, slice[i])
		}
	}
	return
}

func writeUint16SliceBody(writer io.Writer, slice []uint16) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteUint16(writer, slice[i])
		}
	}
	return
}

func writeUint32SliceBody(writer io.Writer, slice []uint32) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteUint32(writer, slice[i])
		}
	}
	return
}

func writeUint64SliceBody(writer io.Writer, slice []uint64) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteUint64(writer, slice[i])
		}
	}
	return
}

func writeUintSliceBody(writer io.Writer, slice []uint) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteUint(writer, slice[i])
		}
	}
	return
}

func writeBoolSliceBody(writer io.Writer, slice []bool) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteBool(writer, slice[i])
		}
	}
	return
}

func writeFloat32SliceBody(writer io.Writer, slice []float32) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteFloat32(writer, slice[i])
		}
	}
	return
}

func writeFloat64SliceBody(writer io.Writer, slice []float64) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteFloat64(writer, slice[i])
		}
	}
	return
}

func writeComplex64SliceBody(enc *Encoder, slice []complex64) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteComplex64(enc, slice[i])
		}
	}
	return
}

func writeComplex128SliceBody(enc *Encoder, slice []complex128) (err error) {
	for i := range slice {
		if err == nil {
			err = WriteComplex128(enc, slice[i])
		}
	}
	return
}

func writeStringSliceBody(enc *Encoder, slice []string) (err error) {
	for i := range slice {
		if err == nil {
			err = EncodeString(enc, slice[i])
		}
	}
	return
}

func writeInterfaceSliceBody(enc *Encoder, slice []interface{}) (err error) {
	for i := range slice {
		if err == nil {
			err = enc.Encode(slice[i])
		}
	}
	return
}

func writeOtherSliceBody(enc *Encoder, slice interface{}) (err error) {
	t := reflect2.TypeOf(slice).(*reflect2.UnsafeSliceType)
	et := t.Elem()
	ptr := reflect2.PtrOf(slice)
	n := t.UnsafeLengthOf(ptr)
	for i := 0; i < n && err == nil; i++ {
		err = enc.Encode(et.UnsafeIndirect(t.UnsafeGetIndex(ptr, i)))
	}
	return
}