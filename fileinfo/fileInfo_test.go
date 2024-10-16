package fileinfo

import (
	"reflect"
	"slices"
	"testing"
)

func TestGetRootInfo(t *testing.T) {
	expected := []FileInfo{{Name: "testfiles/dir1", FileType: "dir", Size: 10}, {Name: "testfiles/file2.txt", FileType: "file", Size: 2}, {Name: "testfiles/file1.txt", FileType: "file", Size: 1}}

	infoList, err := GetRootInfo("testfiles")
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

func TestGenerateFileMap(t *testing.T) {
	m, err := GenerateFileMap(nil, "testfiles")
	if err != nil {
		t.Error(err)
	}
	expected := map[string]FileInfo{
		"testfiles":                     {Name: "testfiles", FileType: "dir", Size: 13, Children: []string{"testfiles/dir1", "testfiles/file1.txt", "testfiles/file2.txt"}, Checked: true},
		"testfiles/dir1":                {Name: "testfiles/dir1", FileType: "dir", Size: 10, Children: []string{"testfiles/dir1/dir2", "testfiles/dir1/file3.txt"}, Checked: true},
		"testfiles/dir1/dir2":           {Name: "testfiles/dir1/dir2", FileType: "dir", Size: 7, Children: []string{"testfiles/dir1/dir2/file4.txt"}, Checked: true},
		"testfiles/dir1/dir2/file4.txt": {Name: "testfiles/dir1/dir2/file4.txt", FileType: "file", Size: 7, Children: []string(nil), Checked: true},
		"testfiles/dir1/file3.txt":      {Name: "testfiles/dir1/file3.txt", FileType: "file", Size: 3, Children: []string(nil), Checked: true},
		"testfiles/file1.txt":           {Name: "testfiles/file1.txt", FileType: "file", Size: 1, Children: []string(nil), Checked: true},
		"testfiles/file2.txt":           {Name: "testfiles/file2.txt", FileType: "file", Size: 2, Children: []string(nil), Checked: true},
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("expected %v, but got %v", expected, m)
	}
}

func TestGetSortedDirs(t *testing.T) {
	m, err := GenerateFileMap(nil, "testfiles")
	if err != nil {
		t.Error(err)
	}

	list := GetSortedDirs(m, "testfiles")
	expected := []FileInfo{{Name: "testfiles/dir1", FileType: "dir", Size: 10}, {Name: "testfiles/file2.txt", FileType: "file", Size: 2}, {Name: "testfiles/file1.txt", FileType: "file", Size: 1}}

	if slices.CompareFunc(expected, list, func(a, b FileInfo) int {
		if a.FileType != b.FileType || a.Name != b.Name || a.Size != b.Size {
			return -1
		}
		return 0
	}) != 0 {
		t.Errorf("slices are not equal. Expected %v, but got %v", expected, list)
	}
}

func TestCleanChildren(t *testing.T) {
	m := map[string]FileInfo{
		"testfiles":      {Name: "testfiles", FileType: "dir", Size: 13, Children: []string{"testfiles/dir1", "testfiles/file1.txt", "testfiles/file2.txt"}, Checked: true},
		"testfiles/dir1": {Name: "testfiles/dir1", FileType: "dir", Size: 10, Children: []string{"testfiles/dir1/dir2", "testfiles/dir1/file3.txt"}, Checked: true},
	}

	expected := map[string]FileInfo{
		"testfiles": {Name: "testfiles", FileType: "dir", Size: 3, Children: []string{"testfiles/file1.txt", "testfiles/file2.txt"}, Checked: true},
	}

	CleanChildren(m, "testfiles/dir1")

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected %v, got %v", expected, m)
	}
}

func BenchmarkGetRootInfo(b *testing.B) {
	for range b.N {
		_, err := GetRootInfo("testfiles")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGetSortedDirs(b *testing.B) {
	for range b.N {
		m, err := GenerateFileMap(nil, "testfiles")
		if err != nil {
			b.Error(err)
		}

		_ = GetSortedDirs(m, "testfiles")
	}
}
