/*--------------------------------------------------------*\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: https://hprose.com                     |
|                                                          |
| io/encoding/encode.go                                    |
|                                                          |
| LastModified: Mar 21, 2020                               |
| Author: Ma Bingyao <andot@hprose.com>                    |
|                                                          |
\*________________________________________________________*/

package encoding

import (
	"math"
	"math/big"
	"reflect"
	"strconv"

	"github.com/hprose/hprose-golang/v3/io"
	"github.com/modern-go/reflect2"
)

const (
	digits = "0123456789"
	digit2 = "" +
		"0001020304050607080910111213141516171819" +
		"2021222324252627282930313233343536373839" +
		"4041424344454647484950515253545556575859" +
		"6061626364656667686970717273747576777879" +
		"8081828384858687888990919293949596979899"
	digit3 = "" +
		"000001002003004005006007008009010011012013014015016017018019" +
		"020021022023024025026027028029030031032033034035036037038039" +
		"040041042043044045046047048049050051052053054055056057058059" +
		"060061062063064065066067068069070071072073074075076077078079" +
		"080081082083084085086087088089090091092093094095096097098099" +
		"100101102103104105106107108109110111112113114115116117118119" +
		"120121122123124125126127128129130131132133134135136137138139" +
		"140141142143144145146147148149150151152153154155156157158159" +
		"160161162163164165166167168169170171172173174175176177178179" +
		"180181182183184185186187188189190191192193194195196197198199" +
		"200201202203204205206207208209210211212213214215216217218219" +
		"220221222223224225226227228229230231232233234235236237238239" +
		"240241242243244245246247248249250251252253254255256257258259" +
		"260261262263264265266267268269270271272273274275276277278279" +
		"280281282283284285286287288289290291292293294295296297298299" +
		"300301302303304305306307308309310311312313314315316317318319" +
		"320321322323324325326327328329330331332333334335336337338339" +
		"340341342343344345346347348349350351352353354355356357358359" +
		"360361362363364365366367368369370371372373374375376377378379" +
		"380381382383384385386387388389390391392393394395396397398399" +
		"400401402403404405406407408409410411412413414415416417418419" +
		"420421422423424425426427428429430431432433434435436437438439" +
		"440441442443444445446447448449450451452453454455456457458459" +
		"460461462463464465466467468469470471472473474475476477478479" +
		"480481482483484485486487488489490491492493494495496497498499" +
		"500501502503504505506507508509510511512513514515516517518519" +
		"520521522523524525526527528529530531532533534535536537538539" +
		"540541542543544545546547548549550551552553554555556557558559" +
		"560561562563564565566567568569570571572573574575576577578579" +
		"580581582583584585586587588589590591592593594595596597598599" +
		"600601602603604605606607608609610611612613614615616617618619" +
		"620621622623624625626627628629630631632633634635636637638639" +
		"640641642643644645646647648649650651652653654655656657658659" +
		"660661662663664665666667668669670671672673674675676677678679" +
		"680681682683684685686687688689690691692693694695696697698699" +
		"700701702703704705706707708709710711712713714715716717718719" +
		"720721722723724725726727728729730731732733734735736737738739" +
		"740741742743744745746747748749750751752753754755756757758759" +
		"760761762763764765766767768769770771772773774775776777778779" +
		"780781782783784785786787788789790791792793794795796797798799" +
		"800801802803804805806807808809810811812813814815816817818819" +
		"820821822823824825826827828829830831832833834835836837838839" +
		"840841842843844845846847848849850851852853854855856857858859" +
		"860861862863864865866867868869870871872873874875876877878879" +
		"880881882883884885886887888889890891892893894895896897898899" +
		"900901902903904905906907908909910911912913914915916917918919" +
		"920921922923924925926927928929930931932933934935936937938939" +
		"940941942943944945946947948949950951952953954955956957958959" +
		"960961962963964965966967968969970971972973974975976977978979" +
		"980981982983984985986987988989990991992993994995996997998999"
)

var minInt64Buf = []byte("-9223372036854775808")

func toBytes(i uint64, buf []byte) (off int) {
	off = len(buf)
	var q, p uint64
	for i >= 100 {
		q = i / 1000
		p = (i - (q * 1000)) * 3
		i = q
		off -= 3
		copy(buf[off:off+3], digit3[p:p+3])
	}
	if i >= 10 {
		q = i / 100
		p = (i - (q * 100)) * 2
		i = q
		off -= 2
		copy(buf[off:off+2], digit2[p:p+2])
	}
	if i > 0 {
		off--
		buf[off] = digits[i]
	}
	return
}

func writeInt64(writer io.BytesWriter, i int64) (err error) {
	if i >= 0 {
		return writeUint64(writer, uint64(i))
	}
	if i == math.MinInt64 {
		_, err = writer.Write(minInt64Buf)
		return err
	}
	var u uint64 = uint64(-i)
	var buf [20]byte
	off := toBytes(u, buf[:]) - 1
	buf[off] = '-'
	_, err = writer.Write(buf[off:])
	return
}

func writeUint64(writer io.BytesWriter, i uint64) (err error) {
	if (i >= 0) && (i <= 9) {
		return writer.WriteByte(digits[i])
	}
	var buf [20]byte
	off := toBytes(i, buf[:])
	_, err = writer.Write(buf[off:])
	return
}

// WriteInt64 to writer
func WriteInt64(writer io.BytesWriter, i int64) (err error) {
	if (i >= 0) && (i <= 9) {
		return writer.WriteByte(digits[i])
	}
	var tag = io.TagInteger
	if (i < math.MinInt32) || (i > math.MaxInt32) {
		tag = io.TagLong
	}
	if err = writer.WriteByte(tag); err == nil {
		if err = writeInt64(writer, i); err == nil {
			err = writer.WriteByte(io.TagSemicolon)
		}
	}
	return
}

// WriteUint64 to writer
func WriteUint64(writer io.BytesWriter, i uint64) (err error) {
	if (i >= 0) && (i <= 9) {
		return writer.WriteByte(digits[i])
	}
	var tag = io.TagInteger
	if i > math.MaxInt32 {
		tag = io.TagLong
	}
	if err = writer.WriteByte(tag); err == nil {
		if err = writeUint64(writer, i); err == nil {
			err = writer.WriteByte(io.TagSemicolon)
		}
	}
	return
}

// WriteInt32 to writer
func WriteInt32(writer io.BytesWriter, i int32) (err error) {
	if (i >= 0) && (i <= 9) {
		return writer.WriteByte(digits[i])
	}
	if err = writer.WriteByte(io.TagInteger); err == nil {
		if err = writeInt64(writer, int64(i)); err == nil {
			err = writer.WriteByte(io.TagSemicolon)
		}
	}
	return
}

// WriteUint32 to writer
func WriteUint32(writer io.BytesWriter, i uint32) (err error) {
	return WriteUint64(writer, uint64(i))
}

// WriteInt16 to writer
func WriteInt16(writer io.BytesWriter, i int16) (err error) {
	return WriteInt32(writer, int32(i))
}

// WriteUint16 to writer
func WriteUint16(writer io.BytesWriter, i uint16) (err error) {
	if (i >= 0) && (i <= 9) {
		return writer.WriteByte(digits[i])
	}
	if err = writer.WriteByte(io.TagInteger); err == nil {
		if err = writeUint64(writer, uint64(i)); err == nil {
			err = writer.WriteByte(io.TagSemicolon)
		}
	}
	return
}

// WriteInt8 to writer
func WriteInt8(writer io.BytesWriter, i int8) (err error) {
	return WriteInt32(writer, int32(i))
}

// WriteUint8 to writer
func WriteUint8(writer io.BytesWriter, i uint8) (err error) {
	return WriteUint16(writer, uint16(i))
}

// WriteInt to writer
func WriteInt(writer io.BytesWriter, i int) (err error) {
	return WriteInt64(writer, int64(i))
}

// WriteUint to writer
func WriteUint(writer io.BytesWriter, i uint) (err error) {
	return WriteUint64(writer, uint64(i))
}

// WriteNil to writer
func WriteNil(writer io.BytesWriter) (err error) {
	return writer.WriteByte(io.TagNull)
}

// WriteBool to writer
func WriteBool(writer io.BytesWriter, b bool) (err error) {
	if b {
		return writer.WriteByte(io.TagTrue)
	}
	return writer.WriteByte(io.TagFalse)
}

func writeFloat(writer io.BytesWriter, f float64, bitSize int) (err error) {
	if f != f {
		return writer.WriteByte(io.TagNaN)
	}
	if f > math.MaxFloat64 {
		if err = writer.WriteByte(io.TagInfinity); err == nil {
			err = writer.WriteByte(io.TagPos)
		}
		return
	}
	if f < -math.MaxFloat64 {
		if err = writer.WriteByte(io.TagInfinity); err == nil {
			err = writer.WriteByte(io.TagNeg)
		}
		return
	}
	if err = writer.WriteByte(io.TagDouble); err == nil {
		var buf [24]byte
		if _, err = writer.Write(strconv.AppendFloat(buf[:0], f, 'g', -1, bitSize)); err == nil {
			err = writer.WriteByte(io.TagSemicolon)
		}
	}
	return
}

// WriteFloat32 to writer
func WriteFloat32(writer io.BytesWriter, f float32) error {
	return writeFloat(writer, float64(f), 32)
}

// WriteFloat64 to writer
func WriteFloat64(writer io.BytesWriter, f float64) error {
	return writeFloat(writer, f, 64)
}

func utf16Length(str string) (n int) {
	length := len(str)
	n = length
	c := 0
	for i := 0; i < length; i++ {
		a := str[i]
		if c == 0 {
			if (a & 0xe0) == 0xc0 {
				c = 1
				n--
			} else if (a & 0xf0) == 0xe0 {
				c = 2
				n -= 2
			} else if (a & 0xf8) == 0xf0 {
				c = 3
				n -= 2
			} else if (a & 0x80) == 0x80 {
				return -1
			}
		} else {
			if (a & 0xc0) != 0x80 {
				return -1
			}
			c--
		}
	}
	if c != 0 {
		return -1
	}
	return n
}

func writeBinary(writer io.BytesWriter, bytes []byte, length int) (err error) {
	if length > 0 {
		err = writeUint64(writer, uint64(length))
	}
	if err == nil {
		if err = writer.WriteByte(io.TagQuote); err == nil {
			if _, err = writer.Write(bytes); err == nil {
				err = writer.WriteByte(io.TagQuote)
			}
		}
	}
	return
}

func writeBytes(writer io.BytesWriter, bytes []byte) (err error) {
	if err = writer.WriteByte(io.TagBytes); err == nil {
		err = writeBinary(writer, bytes, len(bytes))
	}
	return
}

func writeString(writer io.BytesWriter, s string, length int) (err error) {
	if length < 0 {
		return writeBytes(writer, reflect2.UnsafeCastString(s))
	}
	if err = writer.WriteByte(io.TagString); err == nil {
		err = writeBinary(writer, reflect2.UnsafeCastString(s), length)
	}
	return
}

// WriteHead to writer, n is the count of elements in list or map
func WriteHead(writer io.BytesWriter, n int, tag byte) (err error) {
	if err = writer.WriteByte(tag); err == nil {
		if n > 0 {
			err = writeUint64(writer, uint64(n))
		}
		if err == nil {
			err = writer.WriteByte(io.TagOpenbrace)
		}
	}
	return
}

// WriteObjectHead to writer, r is the reference number of struct
func WriteObjectHead(writer io.BytesWriter, r int) (err error) {
	if err = writer.WriteByte(io.TagObject); err == nil {
		if err = writeUint64(writer, uint64(r)); err == nil {
			err = writer.WriteByte(io.TagOpenbrace)
		}
	}
	return
}

// WriteFoot of list or map to writer
func WriteFoot(writer io.BytesWriter) error {
	return writer.WriteByte(io.TagClosebrace)
}

func writeComplex(enc *Encoder, r float64, i float64, bitSize int) (err error) {
	writer := enc.Writer
	if i == 0 {
		return writeFloat(writer, r, bitSize)
	}
	enc.AddReferenceCount(1)
	if err = WriteHead(writer, 2, io.TagList); err == nil {
		if err = writeFloat(writer, r, bitSize); err == nil {
			if err = writeFloat(writer, i, bitSize); err == nil {
				err = WriteFoot(writer)
			}
		}
	}
	return
}

// WriteComplex64 to enc.Writer
func WriteComplex64(enc *Encoder, c complex64) error {
	return writeComplex(enc, float64(real(c)), float64(imag(c)), 32)
}

// WriteComplex128 to enc.Writer
func WriteComplex128(enc *Encoder, c complex128) error {
	return writeComplex(enc, real(c), imag(c), 64)
}

// WriteBigInt to writer
func WriteBigInt(writer io.BytesWriter, i *big.Int) (err error) {
	if err = writer.WriteByte(io.TagLong); err == nil {
		if _, err = writer.Write(reflect2.UnsafeCastString(i.String())); err == nil {
			err = writer.WriteByte(io.TagSemicolon)
		}
	}
	return
}

// WriteBigFloat to writer
func WriteBigFloat(writer io.BytesWriter, f *big.Float) (err error) {
	if err = writer.WriteByte(io.TagDouble); err == nil {
		var buf [32]byte
		if _, err = writer.Write(f.Append(buf[:0], 'g', -1)); err == nil {
			err = writer.WriteByte(io.TagSemicolon)
		}
	}
	return
}

// WriteBigRat to enc.Writer
func WriteBigRat(enc *Encoder, r *big.Rat) (err error) {
	if r.IsInt() {
		return WriteBigInt(enc.Writer, r.Num())
	}
	enc.AddReferenceCount(1)
	s := r.String()
	return writeString(enc.Writer, s, len(s))
}

// WriteError to encoder
func WriteError(enc *Encoder, e error) (err error) {
	if err = enc.Writer.WriteByte(io.TagError); err == nil {
		enc.AddReferenceCount(1)
		s := e.Error()
		err = writeString(enc.Writer, s, utf16Length(s))
	}
	return
}

// EncodeReference to enc
func EncodeReference(valenc ValueEncoder, enc *Encoder, v interface{}) (err error) {
	if reflect2.IsNil(v) {
		return WriteNil(enc.Writer)
	}
	var ok bool
	if ok, err = enc.WriteReference(v); !ok && err == nil {
		err = valenc.Write(enc, v)
	}
	return
}

// SetReference to enc
func SetReference(enc *Encoder, v interface{}) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		enc.SetReference(v)
	} else {
		enc.AddReferenceCount(1)
	}
}