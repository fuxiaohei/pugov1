package packer

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
)

// ReadAll is a copy from ioutil.ReadAll with capacity
func ReadAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}

// Gzip gzip bytes
func Gzip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gGzipPool.GetWriter(&b)
	defer gGzipPool.PutWriter(gz)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Ungzip decode compressed bytes
func Ungzip(data []byte) ([]byte, error) {
	bReader := bytes.NewReader(data)
	r := gGzipPool.GetReader(bReader)
	defer gGzipPool.PutReader(r)
	return ReadAll(r, int64(len(data)*4)) // use four times length to try load ungzip data
}

// thanks gzip pool code from
// https://github.com/ungerik/go-pool/blob/master/gzip.go
var gGzipPool GzipPool

// GzipPool manages a pool of gzip.Writer.
// The pool uses sync.Pool internally.
type GzipPool struct {
	readers sync.Pool
	writers sync.Pool
}

// GetReader returns gzip.Reader from the pool, or creates a new one
// if the pool is empty.
func (pool *GzipPool) GetReader(src io.Reader) (reader *gzip.Reader) {
	if r := pool.readers.Get(); r != nil {
		reader = r.(*gzip.Reader)
		reader.Reset(src)
	} else {
		reader, _ = gzip.NewReader(src)
	}
	return reader
}

// PutReader closes and returns a gzip.Reader to the pool
// so that it can be reused via GetReader.
func (pool *GzipPool) PutReader(reader *gzip.Reader) {
	reader.Close()
	pool.readers.Put(reader)
}

// GetWriter returns gzip.Writer from the pool, or creates a new one
// with gzip.BestCompression if the pool is empty.
func (pool *GzipPool) GetWriter(dst io.Writer) (writer *gzip.Writer) {
	if w := pool.writers.Get(); w != nil {
		writer = w.(*gzip.Writer)
		writer.Reset(dst)
	} else {
		writer, _ = gzip.NewWriterLevel(dst, gzip.BestCompression)
	}
	return writer
}

// PutWriter closes and returns a gzip.Writer to the pool
// so that it can be reused via GetWriter.
func (pool *GzipPool) PutWriter(writer *gzip.Writer) {
	writer.Close()
	pool.writers.Put(writer)
}
