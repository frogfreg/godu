package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type fileInfo struct {
	name     string
	fileType string
	size     int
}

var errBadDescriptor = errors.New("bad file descriptor")

func getSize(entry string) (int, error) {
	fi, err := os.Lstat(entry)
	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Warning: %v\n", err)
			return 0, nil
		} else {
			return 0, err
		}
	}

	if !fi.IsDir() {
		return int(fi.Size()), nil
	}

	var sum int

	entries, err := os.ReadDir(entry)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			fmt.Printf("Warning: %v\n", err)
		} else if strings.Contains(err.Error(), "bad file descriptor") {
			return 0, errBadDescriptor
		} else {
			return 0, err
		}
	}

	for _, e := range entries {
		size, err := getSize(filepath.Join(entry, e.Name()))
		if err != nil {
			return 0, err
		}
		sum += size
	}

	return sum, nil
}

func showRootInfo(root string) error {
	dirEntries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	infoList := []fileInfo{}

	for _, de := range dirEntries {
		name := filepath.Join(root, de.Name())
		size, err := getSize(name)
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				fmt.Printf("Warning: %v\n", err)
			} else if errors.Is(err, errBadDescriptor) {
				fmt.Printf("Warning: %v\n", err)
			} else {
				return err
			}
		}
		fileType := "file"
		if de.IsDir() {
			fileType = "dir"
		}

		infoList = append(infoList, fileInfo{
			name:     name,
			fileType: fileType,
			size:     size,
		})
	}

	slices.SortFunc(infoList, func(a, b fileInfo) int {
		return b.size - a.size
	})

	for _, info := range infoList {
		fmt.Printf("%v | %v | %v bytes\n", info.name, info.fileType, info.size)
	}

	return nil
}

func main() {
	if err := showRootInfo("/"); err != nil {
		panic(err)
	}

}
