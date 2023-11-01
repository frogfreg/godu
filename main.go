package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type fileInfo struct {
	name     string
	fileType string
	size     int
}

func getSize(entry string) (int, error) {
	fi, err := os.Stat(entry)
	if err != nil {
		return 0, err
	}

	if !fi.IsDir() {
		return int(fi.Size()), nil
	}

	var sum int

	entries, err := os.ReadDir(entry)
	if err != nil {
		return 0, err
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
		fileType := "file"
		if de.IsDir() {
			fileType = "dir"
		}

		if err != nil {
			return err
		}

		infoList = append(infoList, fileInfo{
			name:     name,
			fileType: fileType,
			size:     size,
		})
	}

	for _, info := range infoList {
		fmt.Printf("%v | %v | %v bytes\n", info.name, info.fileType, info.size)
	}

	return nil
}

func main() {
	if err := showRootInfo("/Users/nyan/Desktop"); err != nil {
		panic(err)
	}

}
