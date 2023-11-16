package fileinfo

import (
	"slices"
	"testing"
)

func TestGetRootInfo(t *testing.T) {
	expected := []fileInfo{{name: "testfiles/dir1", fileType: "dir", size: 10}, {name: "testfiles/file2.txt", fileType: "file", size: 2}, {name: "testfiles/file1.txt", fileType: "file", size: 1}}

	infoList, err := getRootInfo("testfiles")
	if err != nil {
		t.Error(err)
	}

	if slices.CompareFunc(expected, infoList, func(a, b fileInfo) int {
		if a.fileType != b.fileType || a.name != b.name || a.size != b.size {
			return -1
		}
		return 0
	}) != 0 {
		t.Errorf("slices are not equal. Expected %v, but got %v", expected, infoList)
	}
}
