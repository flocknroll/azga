package main

import (
	"log"
	"path/filepath"

	"github.com/flocknroll/azga/go_installer/addtxtcontentgo"
	"github.com/flocknroll/azga/go_installer/msfstools"
)

func main() {
	// fmt.Println(checkContent("test_data/src_found.txt", "test_data/dest.txt"))
	// fmt.Println(checkContent("test_data/src_not_found.txt", "test_data/dest.txt"))

	var filesList = map[string]string{
		"Airports.txt": "Community/aerosoft-crj/Data/NavData/Airports.txt",
	}

	msfsPath, ok := msfstools.GetPackageFolderPath()

	if !ok {
		log.Fatal("MSFS package folder not found")
		return
	}

	for src, dest := range filesList {
		dest = filepath.Join(msfsPath, dest)
		log.Printf("Checking %s -> %s\n", src, dest)
		addtxtcontentgo.AddContent(src, dest)
	}
}
