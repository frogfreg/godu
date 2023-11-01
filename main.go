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
	fi,err := os.Stat(entry)
	if err != nil{
		return err
	}
	
	if !fi.IsDir(){
		size, err := fi.Size()
		if err != nil {
			return 0, err
		}
		return ingo
	}

	var sum int
}

func showRootInfo(root string)error{
	

	dirEntries,err := os.ReadDir()
	if err != nil{
		return err
	}

	infoList := []fileInfo{}

	for _, de := range dirEntries{
		name := filepath.Join(root,de.Name()) 
		size, err := getSize(name)
		if err != nil{
			return err
		}

		infoList = append(infoList, fileInfo{
			name:name,
			fileType: de.IsDir(),
			size: size
		})
	}
}

func main() {
	if err := showRootInfo("/mnt/f/Desktop"); err != nil {
		panic(err)
	}

}
