package vlc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type BinaryChunks []BinaryChunk

type BinaryChunk string

var chunkSize = 8

func NewBinChunks(data []byte) BinaryChunks {
	res := make(BinaryChunks, 0, len(data))

	for _, code := range data {
		res = append(res, NewBinChunk(code))
	}

	return res
}

func NewBinChunk(code byte) BinaryChunk {
	return BinaryChunk(fmt.Sprintf("%08b", code))
}

// Join joins chunks into one line and return string
func (bcs BinaryChunks) Join() string {
	buf := strings.Builder{}

	for _, chunk := range bcs {
		buf.WriteString(string(chunk))
	}

	return buf.String()
}

// splitChunks 00100000000010001101000011 -> 00100000 00001000 11010000 11 -> 00100000 00001000 11010000 11000000
func splitByChunks(str string, chunkSize int) BinaryChunks {
	strlen := utf8.RuneCountInString(str)
	chunkCount := strlen / chunkSize

	if strlen/chunkSize != 0 {
		chunkCount++
	}
	res := make(BinaryChunks, 0, chunkCount)
	var buf strings.Builder
	for i, ch := range str {
		buf.WriteString(string(ch))

		if (i+1)%chunkSize == 0 {
			res = append(res, BinaryChunk(buf.String()))
			buf.Reset()
		}
	}

	if buf.Len() != 0 {
		lastChunk := buf.String()
		lastChunk += strings.Repeat("0", chunkSize-len(lastChunk))
		res = append(res, BinaryChunk(lastChunk))
	}

	return res
}

func (bcs BinaryChunks) Bytes() []byte {
	res := make([]byte, 0, len(bcs))

	for _, bc := range bcs {
		res = append(res, bc.Byte())
	}

	return res
}

func (bc BinaryChunk) Byte() byte {
	num, err := strconv.ParseUint(string(bc), 2, chunkSize)
	if err != nil {
		panic("can't convert binary chunk to byte: " + err.Error())
	}

	return byte(num)
}
