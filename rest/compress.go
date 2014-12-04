package rest

import (
	"bytes"
	"compress/gzip"
	ferr "github.com/jonas747/fortia/error"
)

const (
	ErrNoCompressionSupported          = 1000 // start at 1k to make space for api errors
	ErrNoHandlerForSelectedCompression = 1001
)

var (
	CompressionAlgos         = []string{"gzip"}
	CompressionAlgosHandlers = map[string]func([]byte) []byte{
		"gzip": compressGzip,
	}
)

func Compress(in []byte, acceptedEncodings map[string]bool) ([]byte, ferr.FortiaError) {
	selectedCompression := ""
	for i := 0; i < len(CompressionAlgos); i++ {
		algo := CompressionAlgos[i]
		if _, ok := acceptedEncodings[algo]; ok {
			selectedCompression = algo
		}
	}
	if selectedCompression == "" {
		return in, ferr.Newc("No compression supported", ErrNoCompressionSupported)
	}
	handler, ok := CompressionAlgosHandlers[selectedCompression]
	if !ok {
		return in, ferr.Newc("No handler for selected compression", ErrNoHandlerForSelectedCompression)
	}

	out := handler(in)

	return out, nil
}

func compressGzip(in []byte) []byte {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	total := 0
	for total < len(in) {
		n, _ := writer.Write(in[total:])
		total += n
	}
	writer.Close()
	out := buffer.Bytes()
	return out
}
