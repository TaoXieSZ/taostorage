package bytestream

import (
	"io"
	"testing"
)

func TestNewByteWriter(t *testing.T) {
	// create a new byte writer
	bw := newByteWriter(1)

	// check the count
	if bw.count != 0 {
		t.Errorf("expected count to be 0, got %d", bw.count)
	}
}

func TestReadBitEOF(t *testing.T) {
	// create a new byte writer
	bw := newByteWriter(1)

	_, err := bw.readBit()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestReadBitEOF2(t *testing.T) {
	bw := newByteReader([]byte{1})
	bw.count = 0
	_, err := bw.readBit()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestReadByteEOF1(t *testing.T) {
	b := newByteWriter(1)
	_, err := b.readByte()
	if err != io.EOF {
		t.Errorf("Unexpected value: %v\n", err)
	}
}

func TestReadByteEOF2(t *testing.T) {
	b := newByteReader([]byte{1})
	b.count = 0
	_, err := b.readByte()
	if err != io.EOF {
		t.Errorf("Unexpected value: %v\n", err)
	}
}

func TestReadByteEOF3(t *testing.T) {
	b := newByteReader([]byte{1})
	b.count = 16
	_, err := b.readByte()
	if err != io.EOF {
		t.Errorf("Unexpected value: %v\n", err)
	}
}

func TestReadBitsEOF(t *testing.T) {
	b := newByteReader([]byte{1})
	_, err := b.readBits(9)
	if err != io.EOF {
		t.Errorf("Unexpected value: %v\n", err)
	}
}

func TestUnmarshalBinaryErr(t *testing.T) {
	b := &stream{}
	err := b.UnmarshalBinary([]byte{})
	if err == nil {
		t.Errorf("An error was expected\n")
	}
}
