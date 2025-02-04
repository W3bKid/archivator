package vlc

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type BinaryChunks []BinaryChunk

type BinaryChunk string

type HexChunk string

type HexChunks []HexChunk

type encodingTable map[rune]string

var chunkSize = 8

func Encode(str string) string {
	// prepare text M -> !m
	preparedText := prepareText(str)

	// text to binary -> 1100100
	bStr := encodeBin(preparedText)

	// split binary by chunks (8): bits to bytes -> 10010010 10001101 10111101
	chunks := splitByChunks(bStr, chunkSize)

	// bytes to hex -> 20 30 2c
	return chunks.ToHex().ToString()
}

func (hcs HexChunks) ToString() string {
	const sep = " "

	switch len(hcs) {
	case 0:
		return ""
	case 1:
		return string(hcs[0])
	}

	var buf bytes.Buffer
	buf.WriteString(string(hcs[0]))

	for _, chunk := range hcs[1:] {
		buf.WriteString(sep)
		buf.WriteString(string(chunk))
	}

	return buf.String()
}

func (bcs BinaryChunks) ToHex() HexChunks {
	res := make(HexChunks, 0, len(bcs))

	for _, chunk := range bcs {
		res = append(res, chunk.ToHex())
	}

	return res
}

func (bc BinaryChunk) ToHex() HexChunk {
	num, err := strconv.ParseUint(string(bc), 2, chunkSize)

	if err != nil {
		panic("cant parse binary chunk: " + err.Error())
	}

	res := strings.ToUpper(fmt.Sprintf("%x", num))

	if len(res) == 1 {
		res = "0" + res
	}

	return HexChunk(res)
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

// prepareText prepares text to be fit for encode:
// changes upper case letters to: ! + lower case letter
// i.g.: My name is Ted -> !my name is !ted
func prepareText(str string) string {
	var buf strings.Builder

	for _, ch := range str {
		if unicode.IsUpper(ch) {
			buf.WriteRune('!')
			buf.WriteRune(unicode.ToLower(ch))
		} else {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}

// encodeBinary encodes str into binary codes string without spaces
func encodeBin(str string) string {
	var buf strings.Builder

	for _, ch := range str {
		buf.WriteString(bin(ch))
	}

	return buf.String()
}

func bin(ch rune) string {
	table := getEncodingTable()

	res, ok := table[ch]

	if !ok {
		panic("unknown character: " + string(ch))
	}

	return res
}

func getEncodingTable() encodingTable {
	return encodingTable{
		' ': "11",
		't': "1001",
		'n': "10000",
		's': "0101",
		'r': "01000",
		'd': "00101",
		'!': "001000",
		'c': "000101",
		'm': "000011",
		'g': "0000100",
		'b': "0000010",
		'v': "00000001",
		'k': "0000000001",
		'q': "000000000001",
		'e': "101",
		'o': "10001",
		'a': "011",
		'i': "01001",
		'h': "0011",
		'l': "001001",
		'u': "00011",
		'f': "000100",
		'p': "0000101",
		'w': "0000011",
		'y': "0000001",
		'j': "000000001",
		'x': "00000000001",
		'z': "000000000000",
	}
}
