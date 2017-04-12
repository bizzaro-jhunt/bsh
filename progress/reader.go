package progress

import (
	"io"
)

type Reader struct {
	Reader io.Reader
	Size   int64

	Draw func(bool, int64, int64, bool)

	started  bool
	offset   int64
	reported int64
}

func (r *Reader) Read(p []byte) (int, error) {
	if !r.started {
		r.Draw(true, 0, r.Size, false)
		r.started = true
	}

	n, err := r.Reader.Read(p)
	r.offset += int64(n)

	if err == nil {
		if r.offset-r.reported > 100/r.Size {
			r.Draw(false, r.offset, r.Size, false)
			r.reported = r.offset
		}
	}
	if err == io.EOF {
		r.Draw(false, r.offset, r.Size, true)
	}

	return n, err
}
