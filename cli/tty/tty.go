package tty

import (
	"comradequinn/hflow/syncio"
	"io"
	"time"
)

type IO struct {
	out     io.Writer
	err     io.Writer
	in      io.Reader
	syncOut *syncio.Writer
	syncErr *syncio.Writer
	syncIn  *syncio.Reader
}

func New(in io.Reader, readTimeout time.Duration, out, err io.Writer) IO {
	return IO{
		out:     out,
		err:     err,
		in:      in,
		syncOut: syncio.NewWriter(out),
		syncErr: syncio.NewWriter(err),
		syncIn:  syncio.NewReader(in),
	}
}

func (s IO) Out() io.Writer { return s.syncOut }
func (s IO) Err() io.Writer { return s.syncErr }

func (s *IO) ReadFunc(accessKey byte, f func(in io.Reader, out io.Writer)) {
	frames, frame := []string{"|", "/", "-", "|", "/", "-"}, 0

	go func() {
		buffer, read, timeout := make([]byte, 1), make(chan byte), time.NewTicker(time.Millisecond*500)

		go func() {
			for {
				_, _ = s.syncIn.Read(buffer)
				read <- buffer[0]
			}
		}()

		for {
			select {
			case input := <-read:
				if input != accessKey {
					continue
				}

				unlockIO := s.Lock()
				f(s.in, s.out)
				unlockIO()
			case <-timeout.C:
				s.syncOut.Write([]byte("\r"))
				s.syncOut.Write([]byte(frames[frame]))
				s.syncOut.Write([]byte(" proxying..."))

				if frame == (len(frames) - 1) {
					frame = 0
					continue
				}
				frame++
			}
		}
	}()
}

func (s *IO) Lock() func() {
	unlockErr := s.syncErr.Lock()
	unlockOut := s.syncOut.Lock()
	unlockIn := s.syncIn.Lock()

	unlockIO := func() {
		unlockErr()
		unlockOut()
		unlockIn()
	}

	return unlockIO
}
