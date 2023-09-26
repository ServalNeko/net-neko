package input

import "os"

type File struct {
	filePath string
}

func NewFile(path string) *File {
	return &File{filePath: path}
}

func (f *File) Read(msgch chan string) error {
	defer close(msgch)

	b, err := os.ReadFile(f.filePath)
	if err != nil {
		return err
	}

	msgch <- string(b)

	return nil
}
