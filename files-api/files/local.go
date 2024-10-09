package files

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	maxFileSize int
	basePath    string
}

func NewLocal(basePath string, maxSize int) (*Local, error) {
	p, err := filepath.Abs(basePath)

	if err != nil {
		return nil, err
	}

	return &Local{maxFileSize: maxSize, basePath: p}, nil
}

func (l *Local) fullPath(path string) string {
	return filepath.Join(l.basePath, path)
}

func (l *Local) Save(path string, contents io.Reader) error {

	fp := l.fullPath(path)

	// get the directory and make sure it exists
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return errors.New("unable to create a directory") // TODO: capture the err
	}

	// if file already exists delete it
	_, err = os.Stat(fp)
	if err == nil {

		err = os.Remove(fp)

		if err != nil {
			return errors.New("unable to delete the file")
		}

	} else if !os.IsNotExist(err) {

		// if it is anything other than not exists err.
		return errors.New("unable to get file info")

	}

	f, err := os.Create(fp)
	if err != nil {
		return errors.New("unable to create file")
	}
	defer f.Close()

	// Write the contents to the new file.
	// TODO: ensure that we are not writing greater than max bytes.

	_, err = io.Copy(f, contents)

	if err != nil {
		return errors.New("unable to write to file")
	}

	return nil
}

func (l *Local) Get(path string) (*os.File, error) {

	fp := l.fullPath(path)

	// open the file
	f, err := os.Open(fp)
	if err != nil {
		return nil, errors.New("unable to open the file")
	}

	return f, nil
}
