package try

import (
	"bufio"
	"bytes"
	"github.com/getlantern/mockconn"
	"net"
	"net/http/httptest"
)

// NewHijackableRecorder Creates a new HijackableResponseRecorder that sends the given
// hijackedInputData (okay to leave nil).
func NewHijackableRecorder(hijackedInputData []byte) *HijackableResponseRecorder {
	wrapped := httptest.NewRecorder()
	h := &HijackableResponseRecorder{
		ResponseRecorder: wrapped,
		in:               bufio.NewReader(bytes.NewBuffer(hijackedInputData)),
		out:              bufio.NewWriter(wrapped.Body),
	}
	h.conn = mockconn.New(h.Body, h.in)
	return h
}

// HijackableResponseRecorder is very similar to
// net/http/httputil.ResponseRecorder but also allows hijacking. The results of
// data received after hijacking are available in the Body() just as
// non-hijacked responses.
type HijackableResponseRecorder struct {
	*httptest.ResponseRecorder
	in   *bufio.Reader
	out  *bufio.Writer
	conn *mockconn.Conn
}

func (h *HijackableResponseRecorder) Closed() bool {
	return h.conn.Closed()
}

func (h *HijackableResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return h.conn, bufio.NewReadWriter(h.in, h.out), nil
}
