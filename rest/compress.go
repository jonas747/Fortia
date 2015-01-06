package rest

import (
	"bytes"
	"compress/gzip"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
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

func Compress(in []byte, r *Request) ([]byte, errors.FortiaError) {
	selectedCompression := ""
	for i := 0; i < len(CompressionAlgos); i++ {
		algo := CompressionAlgos[i]
		if _, ok := r.AcceptedEncodings[algo]; ok {
			selectedCompression = algo
		}
	}
	if selectedCompression == "" {
		return in, NewRestError(r, errorcodes.ErrorCode_ClientSupportsNoCompression, "No compression supported", nil)
	}

	handler, ok := CompressionAlgosHandlers[selectedCompression]
	if !ok {
		return in, NewRestError(r, errorcodes.ErrorCode_NoHandlerForSelectedCompression, "", nil)
	}

	out := handler(in)
	r.RW.Header().Set("Content-Encoding", selectedCompression)

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
