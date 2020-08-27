package main

import (
	"os"
	"io"
	"fmt"
	"strings"
)

const hex = "0123456789abcdef"

var (
	cols uint8		/* octets per line */
	nogo uint8		/* no of grouped octets */
)

func init() {
	set_defaults()
}

func set_defaults() {
	cols = 16
	nogo = 2
}

/*
xxx("./0123456789:;<=", 16, 1) returns:
	"2e 2f 30 31 32 33 34 35  36 37 38 39 3a 3b 3c 3d"
	                         ^ extra space

xxx("./0123456789:;<=", 16, 2) returns:
	"2e2f 3031 3233 3435  3637 3839 3a3b 3c3d"
	 ^ grouped octets

xxx("10", 16, 2) returns:
	"3130                                    "
	                                        ^ full length

xxx("./0123456789:;<=", 16, 4) returns:
	"2e2f3031 32333435  36373839 3a3b3c3d"
	 ^ big endian by default

*/
func hexdata(data []byte, cols uint8, nogo uint8) []byte {
	var (
		length uint32
		buf []byte
		j uint32
	)

	length = uint32(2 * cols + cols / nogo)
	buf = make([]byte, length)

	for i, c := range data {
		if i > 0 && i % int(nogo) == 0 {
			buf[j] = ' '; j++
		}

		if i == int(cols) / 2 {
			buf[j] = ' '; j++
		}

		buf[j] = hex[c >> 4]; j++
		buf[j] = hex[c & 0x0f]; j++
	}

	for ; j < length; j++ {
		buf[j] = ' '
	}

	return buf
}

func toChar(b byte) byte {
	if 32 <= b && b <= 126 {
		return b
	}
	return '.'
}

/*
Output format:
	00000010: 2e2f 3031 3233 3435  3637 3839 3a3b 3c3d  ./0123456789:;<=
	^ offset  ^ grouped octets    ^ extra space         ^ ASCII of line

*/
func hex_string(data []byte, off uint32) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%08x: ", off))	/* offset */

	b.Write(hexdata(data, cols, nogo))
	b.WriteString("  ")

	for _, c := range data {
		b.WriteByte(toChar(c))		/* ASCII of line */
	}
	b.WriteByte('\n')

	return b.String()
}

func process(r io.Reader, w io.Writer) {
	var off uint32 = 0
	buf := make([]byte, cols)

	for {
		n, err := io.ReadFull(r, buf)
		if n == 0 {
			fmt.Println("io.ReadFull:", err)
			break
		}

		s := hex_string(buf[:n], off)
		w.Write([]byte(s))
		off += uint32(n)

		if err == io.ErrUnexpectedEOF {
			break
		}
	}
}

func main() {
	set_defaults()
	process(os.Stdin, os.Stdout)
}

