package http

import (
	"bufio"
	"net"
	"net/http"
)

type codeRespWriter struct {
	http.ResponseWriter
	code int
}

func newCodeResponseWriter(w http.ResponseWriter) *codeRespWriter {
	return &codeRespWriter{
		ResponseWriter: w,
		code:           http.StatusOK,
	}
}

func (w *codeRespWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *codeRespWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h := w.ResponseWriter.(http.Hijacker)
	return h.Hijack()
}
