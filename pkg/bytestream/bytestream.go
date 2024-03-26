package bytestream

import (
	"bytes"
	"encoding/binary"
	"io"
)

type stream struct {
	stream []byte

	// how many bits are valid in the current byte
	count uint8
}

func newByteReader(b []byte) *stream {
	return &stream{
		stream: b,
		count:  8,
	}
}

func newByteWriter(size int) *stream {
	return &stream{
		stream: make([]byte, size),
		count:  0,
	}
}

// clone returns a copy of the stream
func (b *stream) clone() *stream {
	return &stream{
		stream: append([]byte(nil), b.stream...),
		count:  b.count,
	}
}

func (b *stream) bytes() []byte {
	return b.stream
}

func (b *stream) writeBit(isZero bool) {

	// first check last byte is full
	if b.count == 0 {
		b.stream = append(b.stream, 0)
		b.count = 8
	}

	i := len(b.stream) - 1

	// if isZero is false, set the last bit to 1
	if !isZero {
		b.stream[i] |= 1 << (b.count - 1)
	}

	b.count--
}

func (b *stream) writeByte(val byte) {
	// first check the last byte is full
	if b.count == 0 {
		b.stream = append(b.stream, 0)
		b.count = 8
	}

	// fill up the last byte with the first b.count bits of val
	b.stream[len(b.stream)-1] |= val >> (8 - b.count)
	// add the remaining bits to the next byte
	b.stream = append(b.stream, val<<(b.count))
}

func (b *stream) writeBits(val uint64, n uint8) {
	for i := uint8(0); i < n; i++ {
		b.writeBit(val&(1<<i) == 0)
	}
}

// readBit reads a single bit from the reader stream
func (b *stream) readBit() (bool, error) {
	if len(b.stream) == 0 {
		return false, io.EOF
	}

	if b.count == 0 {
		b.stream = b.stream[1:]
		if len(b.stream) == 0 {
			return false, io.EOF
		}
		b.count = 8
	}

	b.count--
	return b.stream[0]&(1<<b.count) != 0, nil
}

func (b *stream) readByte() (byte, error) {
	if len(b.stream) == 0 {
		return 0, io.EOF
	}

	// current byte is empty, move to the next byte
	if b.count == 0 {
		b.stream = b.stream[1:]
		if len(b.stream) == 0 {
			return 0, io.EOF
		}
		b.count = 8
	}

	if b.count == 8 {
		b.count = 0
		return b.stream[0], nil
	}

	// first get the current byte
	val := b.stream[0]
	b.stream = b.stream[1:]

	if len(b.stream) == 0 {
		return 0, io.EOF
	}

	// then get the remaining bits from the next byte
	val |= b.stream[0] >> b.count
	b.stream[0] <<= (8 - b.count)

	return val, nil
}

func (b *stream) readBits(n uint8) (uint64, error) {
	var val uint64

	for i := uint8(0); i < n; i++ {
		bit, err := b.readBit()
		if err != nil {
			return 0, err
		}
		if bit {
			val |= 1 << i
		}
	}

	return val, nil
}

func (b *stream) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, b.count); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, b.stream); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *stream) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.BigEndian, &b.count); err != nil {
		return err
	}
	b.stream = make([]byte, buf.Len())
	return binary.Read(buf, binary.BigEndian, &b.stream)
}
