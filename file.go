// Copyright 2013, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zip

import (
	"io"
	"os"

	. "github.com/russinholi/go-zip/c"
)

// File provides ability to read a file in the ZIP archive.
type File struct {
	*FileHeader
	a *Archive
}

// Open returns a ReadCloser that provides access to the File's contents.
func (f *File) Open() (rc io.ReadCloser, err error) {
	fr, err := f.a.z.Fopen(f.Name, 0)
	if err != nil {
		return nil, err
	}
	return &fileReader{fr}, nil
}

type fileWriter struct {
	rpipe *os.File
	wpipe *os.File
	done  chan error
}

func newFileWriter(z *Zip, name, comment string) (w *fileWriter, err error) {
	w = &fileWriter{
		rpipe: nil,
		wpipe: nil,
		done:  make(chan error)}
	w.rpipe, w.wpipe, err = os.Pipe()
	if err != nil {
		return nil, err
	}

	err = z.AddFd(name, comment, w.rpipe.Fd())
	if err != nil {
		w.wpipe.Close()
		w.rpipe.Close()
		return nil, err
	}

	return w, nil
}

func (w *fileWriter) Write(p []byte) (nn int, err error) {
	return w.wpipe.Write(p)
}

func (w *fileWriter) Close() error {
	w.wpipe.Close()
	w.wpipe = nil
	err := <-w.done // wait for the archive finishing writing.
	w.rpipe.Close()
	w.rpipe = nil
	return err
}

type fileReader struct {
	f *ZipFile
}

func (r *fileReader) Read(b []byte) (int, error) {
	n, err := r.f.Read(b)
	return int(n), err
}

func (r *fileReader) Close() error {
	return r.f.Close()
}
