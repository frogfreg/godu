package fileinfo

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/table"
	"github.com/frogfreg/godu/utilities"
)

type FileInfo struct {
	Name     string
	FileType string
	Size     int
	Children []string
}

var errBadDescriptor = errors.New("bad file descriptor")

func getSize(entry string) (int, error) {
	fi, err := os.Lstat(entry)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// fmt.Printf("Warning: %v\n", err)
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
			// fmt.Printf("Warning: %v\n", err)
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

func FileInfosToRow(fis []FileInfo) []table.Row {
	var rows []table.Row
	for _, fi := range fis {
		rows = append(rows, []string{filepath.Base(fi.Name), fi.FileType, utilities.HumanReadableByteString(fi.Size)})
	}
	return rows
}

func GetRootInfo(root string) ([]FileInfo, error) {
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
					// fmt.Printf("Warning: %v\n", err)
				} else if errors.Is(err, errBadDescriptor) {
					// fmt.Printf("Warning: %v\n", err)
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

func GenerateFileMap(root string) (map[string]FileInfo, error) {
	fileMap := map[string]FileInfo{root: {Name: root, FileType: "dir", Size: 0}}
	f := getMapFillerFunc(fileMap)

	if walkErr := filepath.WalkDir(root, f); walkErr != nil {
		return nil, walkErr
	}

	return fileMap, nil
}

func getMapFillerFunc(m map[string]FileInfo) func(path string, d fs.DirEntry, err error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		parent := filepath.Dir(path)
		if fi, exists := m[parent]; exists {
			fi.Children = append(fi.Children, path)
			m[parent] = fi
		} else {
			m[parent] = FileInfo{Name: parent, FileType: "dir", Size: 0, Children: []string{path}}
		}

		if d.IsDir() {
			m[path] = FileInfo{Name: path, FileType: "dir", Size: 0}
		} else {
			info, err := d.Info()
			if err != nil {
				return err
			}
			m[path] = FileInfo{Name: path, FileType: "file", Size: int(info.Size())}
		}
		return nil
	}
}
