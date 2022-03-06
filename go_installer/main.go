package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/flocknroll/azga/go_installer/addtxtcontentgo"
	"github.com/flocknroll/azga/go_installer/msfstools"
)

type SourceType string

const (
	Local SourceType = "local"
	Git   SourceType = "git"
)

type InstallerConfig struct {
	Config []ConfigEntry
}

type ConfigEntry struct {
	SourceType SourceType
	SourcePath string
	DestPath   string
}

func main() {
	// fmt.Println(checkContent("test_data/src_found.txt", "test_data/dest.txt"))
	// fmt.Println(checkContent("test_data/src_not_found.txt", "test_data/dest.txt"))
	msfsPath, ok := msfstools.GetPackageFolderPath()

	if !ok {
		log.Fatal("MSFS package folder not found")
		return
	}

	var config InstallerConfig

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	for _, ce := range config.Config {
		switch ce.SourceType {
		case Local:
			dest := filepath.Join(msfsPath, ce.DestPath)
			log.Printf("Checking %s -> %s\n", ce.SourcePath, dest)
			addtxtcontentgo.AddContent(ce.SourcePath, dest)
		case Git:
		default:
			panic("Invalid source type")
		}
	}
}
