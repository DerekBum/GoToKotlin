package GoToJava

import (
	"bytes"
	"io"
	"sync"
)

const MAXLEN = 134217728 * 4

type Accumulator struct {
	buf    *bytes.Buffer
	writer io.Writer
	waiter *sync.WaitGroup
}

func NewAccumulator(writer io.Writer) *Accumulator {
	return &Accumulator{buf: &bytes.Buffer{}, writer: writer, waiter: &sync.WaitGroup{}}
}

func (a *Accumulator) Write(p []byte) (int, error) {
	n, err := a.buf.Write(p)
	if err != nil {
		return n, err
	}

	if a.buf.Len() >= MAXLEN {
		var cpy = make([]byte, a.buf.Len())
		copy(cpy, a.buf.Bytes())

		a.waiter.Wait()
		a.waiter.Add(1)

		go func() {
			a.writer.Write(cpy)
			a.waiter.Done()
		}()
		a.buf.Reset()
	}

	return n, err
}

func (a *Accumulator) WriteRest() (int, error) {
	a.waiter.Wait()

	n, err := a.writer.Write(a.buf.Bytes())
	a.buf.Reset()

	return n, err
}
