package earg

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os/exec"
)

type ArecordSource struct {
	currentOffset int
	cmd           *exec.Cmd
	stdout        io.ReadCloser
	buf           []byte
	sampleRate    int
}

func NewArecordSource() (*ArecordSource, error) {
	// 'cd' format is 16bit little endian, 44100 Hz, stereo
	ar := &ArecordSource{
		currentOffset: 0,
		sampleRate:    44100,
	}
	ar.buf = make([]byte, 0)
	command := "arecord"
	args := []string{"-f", "cd"}
	ar.cmd = exec.Command(command, args...)
	var err error
	ar.stdout, err = ar.cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Can't get stdout pipe : %w", err)
	}
	err = ar.cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("Can't start command [%s] args [%s] : %w", command, args, err)
	}
	return ar, nil
}

func (ar *ArecordSource) SampleRate() int {
	return 44100
}

func (ar *ArecordSource) CurrentOffset() int {
	return ar.currentOffset
}

func (ar *ArecordSource) Read(b []float64) (int, error) {
	readBuf := make([]byte, len(b)*4) // 4 = 16bit stereo
	numRead, err := ar.stdout.Read(readBuf)
	if err != nil {
		return 0, fmt.Errorf("Error reading from arecord stdout: %w", err)
	}
	ar.buf = append(ar.buf, readBuf[:numRead]...)
	// give as much of ar.buf as possible to 'b' buffer, keep the rest
	bytesPerFloat := 4 // 16bit LE, 2 channels
	bytesToCopy := len(ar.buf)
	if bytesToCopy > len(b)*bytesPerFloat {
		bytesToCopy = len(b) * bytesPerFloat
	}
	twoBytesToFloat64 := func(buf []byte) float64 {
		s := int16(binary.LittleEndian.Uint16(buf))
		return float64(s) / float64(math.MaxInt16)
	}
	for i := 0; i < bytesToCopy; i += bytesPerFloat {
		left := twoBytesToFloat64(ar.buf[i:])
		right := twoBytesToFloat64(ar.buf[i+2:])
		f := left + right/2.0
		b[i/bytesPerFloat] = f
	}
	ar.buf = ar.buf[bytesToCopy:]

	numFloats := bytesToCopy / bytesPerFloat
	ar.currentOffset += numFloats
	return numFloats, nil
}
