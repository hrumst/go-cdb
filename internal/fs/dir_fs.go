package fs

import (
	"os"
)

type dir struct {
	dirPath string
}

func NewDir(dirPath string) *dir {
	return &dir{
		dirPath: dirPath,
	}
}

func (d dir) filepath(filename string) string {
	return d.dirPath + string(os.PathSeparator) + filename
}

func (d dir) WriteSync(filename string, data []byte) (os.FileInfo, error) {
	filePath := d.filepath(filename)
	fd, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer fd.Close()

	if err != nil {
		return nil, err
	}
	if _, err := fd.Write(data); err != nil {
		return nil, err
	}
	if err := fd.Sync(); err != nil {
		return nil, err
	}
	return fd.Stat()
}

func (d dir) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(d.filepath(filename))
}

func (d dir) FilesStats() ([]os.FileInfo, error) {
	files, err := os.ReadDir(d.dirPath)
	if err != nil {
		return nil, err
	}

	result := make([]os.FileInfo, 0, len(files))
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}
		result = append(result, fileInfo)
	}
	return result, nil
}
