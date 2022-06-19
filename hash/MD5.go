package hash

import (
	"bytes"
	"fmt"
	"math"
)

//Does MD5 Hashing
//Based off of rfc 1321: www.ietf.org/rfc/rfc1321.txt
//Other Source help for specifics: rosettacode.org/wiki/MD5/Implementation#Go, https://en.wikipedia.org/wiki/MD5

//All shift amounts
var shift = []int{
	7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
	5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20,
	4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
	6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21,
}

//Hashes the input and returns the hashed input
//All in Little Endian
func Hash(toHash []byte) (string, error) {
	buf := bytes.NewBuffer(toHash)
	appendPadding(buf)
	appendLength(buf, len(toHash)*8)

	//Initialize the table
	table := [64]uint32{}
	for i := range table {
		table[i] = uint32(4294967296 * math.Abs(math.Sin(float64(i+1))))
	}

	//Initialize the A,B,C,D 32-bit registers
	A := uint32(0x67452301)
	B := uint32(0xEFCDAB89)
	C := uint32(0x98BADCFE)
	D := uint32(0x10325476)

	xBuf := [16]uint32{}

	fmt.Println(buf.Len())
	for i := 0; i < (buf.Len() / 64); i++ {
		temp := buf.Next(64)
		for i := range xBuf {
			xBuf[i] = (uint32(temp[i*4+3]) << 24) | (uint32(temp[i*4+2]) << 16) | (uint32(temp[i*4+1]) << 8) | (uint32(temp[i*4]))
		}

		ATemp := A
		BTemp := B
		CTemp := C
		DTemp := D
		//Rounds
		for j := 0; j < 64; j++ {
			var F uint32
			var g int
			if j < 16 {
				F = (BTemp & CTemp) | ((^BTemp) & DTemp)
				g = j
			} else if j < 32 {
				F = (DTemp & BTemp) | ((^DTemp) & CTemp)
				g = (5*j + 1) % 16
			} else if j < 48 {
				F = BTemp ^ CTemp ^ DTemp
				g = (3*j + 5) % 16
			} else {
				F = CTemp ^ (BTemp | (^DTemp))
				g = (7 * j) % 16
			}
			F = F + ATemp + table[j] + xBuf[g]
			ATemp = DTemp
			DTemp = CTemp
			CTemp = BTemp
			BTemp = BTemp + leftRotate(F, shift[j])
			fmt.Println(ATemp, BTemp, CTemp, DTemp)
		}
		A += ATemp
		B += BTemp
		C += CTemp
		D += DTemp
		fmt.Printf("A:%X\nB:%X\nC:%X\nD:%X\n", A, B, C, D)
	}

	fmt.Printf("Little A:%X\n", toLittleEndian(A))

	outputInt := []uint32{
		toLittleEndian(A),
		toLittleEndian(B),
		toLittleEndian(C),
		toLittleEndian(D),
	}

	return uintToString(outputInt), nil
}

//Appends the bits to the input so that it is 448 mod 512 aka 64 bits shy of a 512 bit multiple
//<< 9 is * 512, >> 9 is / 512
func appendPadding(buf *bytes.Buffer) {
	//Add a 1 as the last bit
	buf.WriteByte(0b10000000)

	//Add as many zeroes as needed to bring the length to 64 bits less than 512 bits
	for buf.Len()%64 != 56 {
		buf.WriteByte(0)
	}
}

//Appends the length of the original message as a 2^64 number
//Little Endian
func appendLength(buf *bytes.Buffer, lenOg int) {
	len64 := uint64(lenOg)
	buf.WriteByte(byte((len64 << 56) >> 56))
	buf.WriteByte(byte((len64 << 48) >> 56))
	buf.WriteByte(byte((len64 << 40) >> 56))
	buf.WriteByte(byte((len64 << 32) >> 56))
	buf.WriteByte(byte((len64 << 24) >> 56))
	buf.WriteByte(byte((len64 << 16) >> 56))
	buf.WriteByte(byte((len64 << 8) >> 56))
	buf.WriteByte(byte(len64 >> 56))
}

//Converts the Big Endian number to Little Endian
func toLittleEndian(bigEnd uint32) uint32 {
	fmt.Println(bigEnd)
	return (bigEnd << 24) | (bigEnd >> 8 << 24 >> 8) | (bigEnd >> 16 << 24 >> 16) | (bigEnd >> 24)
}

//Converts an array of uint32 to hexadecimal string
func uintToString(arr []uint32) string {
	str := ""
	for _, elem := range arr {
		str += fmt.Sprintf("%X", elem)
	}
	return str
}

func leftRotate(toShift uint32, shiftBy int) uint32 {
	return (toShift << shiftBy) | (toShift >> (32 - shiftBy))
}
