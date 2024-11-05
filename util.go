package homeassistant

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func FromBitmap[T any](data byte, result T) T {
	r := reflect.ValueOf(&result).Elem()
	t := r.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if bitmapPos, ok := field.Tag.Lookup("bitmap"); ok {
			pos, err := strconv.Atoi(bitmapPos)
			if err != nil {
				continue
			}
			if bitIsSet(data, pos) {
				r.Field(i).SetBool(true)
			}
		}
	}
	return result
}

func bitIsSet(b byte, pos int) bool {
	return b&(1<<pos) != 0
}

func TcKelvinMirek(value int) int {
	return 1000000 / value
}

func JsonMustMarshal(v interface{}) []byte {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bs
}

func MustParseFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func MustParseByte(s string) byte {
	b, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		panic(err)
	}
	return byte(b)
}

func MustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func ReadBigEndianUInt64(data []byte) (uint64, error) {
	if len(data) > 8 {
		return 0, fmt.Errorf("data is too short, need at least 6 bytes")
	}

	// Create a buffer with extra two bytes for sign extension
	buf := make([]byte, 8)

	// Copy the 6 bytes into the buffer starting at the 2nd byte
	copy(buf[8-len(data):], data)

	// Convert the 8-byte buffer to int64
	var result uint64
	err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &result)
	if err != nil {
		return 0, err
	}

	return result, nil
}
