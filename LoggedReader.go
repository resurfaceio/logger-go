package logger

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const OVERFLOWED = "{ \"overflowed\": %d }"

type LoggedReader struct {
	io.Reader
	logged []byte
}

func NewLoggedReader(r io.Reader, limit int) (lr *LoggedReader, err error) {
	if r == nil {
		return nil, errors.New("nil reader")
	}
	lr = &LoggedReader{}

	loggedBytes := 0
	overflowed := false

	w := new(bytes.Buffer)
	buf := make([]byte, 1024)

	var n int
	for n > 0 || err == nil {
		n, err = r.Read(buf)
		loggedBytes += n
		if loggedBytes > limit {
			overflowed = true
		} else {
			w.Write(buf[:n])
		}
	}

	if err == io.EOF {
		err = nil
	}

	if err != nil {
		return nil, err
	}

	if overflowed {
		w = bytes.NewBufferString(fmt.Sprintf(OVERFLOWED, loggedBytes))
	}

	lr.logged = w.Bytes()
	lr.Reader = bytes.NewReader(lr.logged)

	return
}
