package fileinfo

import (
	"slices"
	"testing"
)

func TestGetRootInfo(t *testing.T) {
	expected := []FileInfo{{Name: "testfiles/dir1", FileType: "dir", Size: 10}, {Name: "testfiles/file2.txt", FileType: "file", Size: 2}, {Name: "testfiles/file1.txt", FileType: "file", Size: 1}}

	infoList, err := getRootInfo("testfiles")
	if err != nil {
		t.Error(err)
	}

	if slices.CompareFunc(expected, infoList, func(a, b FileInfo) int {
		if a.FileType != b.FileType || a.Name != b.Name || a.Size != b.Size {
			return -1
		}
		return 0
	}) != 0 {
		t.Errorf("slices are not equal. Expected %v, but got %v", expected, infoList)
	}
}
