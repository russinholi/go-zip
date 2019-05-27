// Copyright 2013, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package go-zip is a wrapper around C library libzip. It tries to mimic standard library "archive/zip" but provides ability to modify existing ZIP archives.

SEE: http://www.nih.at/libzip/index.html
*/
package zip

import (
	"io"

	. "github.com/russinholi/go-zip/c"
)

// Archive provides ability for reading, creating and modifying a ZIP archive.
type Archive struct {
	z *Zip
}

// Open a ZIP archive, create a new one if not exists.
func Open(file string) (a *Archive, err error) {
	a = &Archive{
		z: &Zip{Path: file}}
	err = a.z.Open()
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Close z ZIP archive, and modifications get written to the disk when closing.
func (a *Archive) Close() error {
	return a.z.Close()
}

// Add a file or directory in the ZIP archive.
// If name ends with '/', a directory will be added, and returned w will be nil.
// Otherwise, a file will be added, and w is a Writer to which the file contents should be written.
func (a *Archive) Create(name string) (w io.WriteCloser, err error) {
	if len(name) > 0 && name[len(name)-1] == '/' {
		return a.createDirectory(name)
	}
	return a.createFile(name)
}

// Add a file in the ZIP archive with the supplied comment.
func (a *Archive) CreateFileWithComment(name, comment string) (w io.WriteCloser, err error) {
	return a.createFileWithComment(name, comment)
}

func (a *Archive) createDirectory(name string) (io.WriteCloser, error) {
	_, err := a.z.AddDir(name)
	return nil, err
}

func (a *Archive) createFile(name string) (io.WriteCloser, error) {
	return a.createFileWithComment(name, "")
}

func (a *Archive) createFileWithComment(name, comment string) (io.WriteCloser, error) {
	f, err := newFileWriter(a.z, name, comment)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// ZIP entry count in the ZIP archive.
func (a *Archive) Count() int64 {
	c, err := a.z.GetNumEntries(0)
	if err != nil {
		return 0
	}
	return c
}

// Delete a ZIP entry in the ZIP archive.
func (a *Archive) Delete(name string) error {
	index, err := a.z.NameLocate(name, 0)
	if err != nil {
		return err
	}
	return a.z.Delete(uint64(index))
}

// Rename a ZIP entry in the ZIP archive.
func (a *Archive) Rename(from, to string) error {
	index, err := a.z.NameLocate(from, 0)
	if err != nil {
		return err
	}
	return a.z.Rename(uint64(index), to)
}

// Get a file from the ZIP archive.
func (a *Archive) File(index uint64) (*File, error) {
	h, err := a.z.FileHeader(index)
	if err != nil {
		return nil, err
	}
	return &File{h, a}, nil
}
