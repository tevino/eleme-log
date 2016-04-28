package log

import "bytes"

type fakeWriter struct {
	writed chan bool
	buf    *bytes.Buffer
}

func (f *fakeWriter) Write(p []byte) (n int, err error) {
	f.buf.Write(p)
	f.writed <- true
	return
}

func (f *fakeWriter) String() string {
	return f.buf.String()
}
