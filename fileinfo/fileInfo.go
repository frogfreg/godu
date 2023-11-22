package fileinfo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type FileInfo struct {
	Name     string
	FileType string
	Size     int
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

func getRootInfo(root string) ([]FileInfo, error) {
	infoList := []FileInfo{}

	dirEntries, err := os.ReadDir(root)
	if err != nil {
		return infoList, err
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error)

	for _, de := range dirEntries {
		de := de
		wg.Add(1)
		go func() {
			defer wg.Done()
			name := filepath.Join(root, de.Name())
			size, err := getSize(name)
			if err != nil {
				if errors.Is(err, os.ErrPermission) {
					fmt.Printf("Warning: %v\n", err)
				} else if errors.Is(err, errBadDescriptor) {
					fmt.Printf("Warning: %v\n", err)
				} else {
					errChan <- err
					return
				}
			}
			errChan <- nil
			fileType := "file"
			if de.IsDir() {
				fileType = "dir"
			}

			mu.Lock()
			infoList = append(infoList, FileInfo{
				Name:     name,
				FileType: fileType,
				Size:     size,
			})
			mu.Unlock()
		}()
	}

	for range dirEntries {
		if err := <-errChan; err != nil {
			return infoList, err
		}
	}

	wg.Wait()

	slices.SortFunc(infoList, func(a, b FileInfo) int {
		return b.Size - a.Size
	})

	return infoList, err
}
