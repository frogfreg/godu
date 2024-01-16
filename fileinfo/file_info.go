package fileinfo

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/bubbles/table"
)

type FileInfo struct {
	Name     string
	FileType string
	Size     int
}

// var errBadDescriptor = errors.New("bad file descriptor")

func FileInfosToRow(fis []FileInfo) []table.Row {
	var rows []table.Row
	for _, fi := range fis {
		rows = append(rows, []string{filepath.Base(fi.Name), fi.FileType, fmt.Sprintf("%v", fi.Size)})
	}
	return rows
}

func getFileInfo(root string, d os.DirEntry) (FileInfo, error) {

	if !d.IsDir() {
		info, err := d.Info()
		if err != nil {
			return FileInfo{}, err
		}

		return FileInfo{Name: root, FileType: "file", Size: int(info.Size())}, nil
	}

	fi := FileInfo{Name: root, FileType: "dir", Size: 0}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil && !os.IsPermission(err) {
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		fi.Size += int(info.Size())

		return nil
	})

	if err != nil {
		return fi, err
	}

	return fi, nil
}

func GetRootInfo(root string) ([]FileInfo, error) {
	var fis []FileInfo

	entries, err := os.ReadDir(root)
	if err != nil {
		return fis, err
	}

	fiChan := make(chan FileInfo)
	errChan := make(chan error)

	for _, e := range entries {
		e := e
		go func() {
			fi, err := getFileInfo(filepath.Join(root, e.Name()), e)
			if err != nil {
				errChan <- err
				return
			}
			fiChan <- fi
		}()
	}

	for i := 0; i < len(entries); i++ {
		select {
		case err := <-errChan:
			return fis, err
		case fi := <-fiChan:
			fis = append(fis, fi)
		}
	}

	slices.SortFunc(fis, func(a, b FileInfo) int {
		return b.Size - a.Size
	})

	return fis, nil
}
