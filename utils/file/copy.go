package file

import (
	"io"
	"os"
)

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
// Based on https://stackoverflow.com/a/21061062 (MIT License)
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
