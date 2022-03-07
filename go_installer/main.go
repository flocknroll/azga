package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/flocknroll/azga/go_installer/addtxtcontentgo"
	"github.com/flocknroll/azga/go_installer/msfstools"
	"github.com/flocknroll/azga/go_installer/utils"
)

type SourceType string

const (
	Local SourceType = "local"
	Git   SourceType = "git"
	Http  SourceType = "http"
)

type InstallerConfig struct {
	Config []ConfigEntry
}

type ConfigEntry struct {
	SourceType SourceType
	SourcePath string
	DestPath   string
	CreateFile bool
}

func main() {
	// Configuration
	var config InstallerConfig
	var configPath, msfsPath string

	flag.StringVar(&configPath, "config", "config.json", "Path to the JSON config file")
	flag.StringVar(&msfsPath, "package-path", "", "Path to the MSFS package folder")
	flag.Parse()

	if msfsPath == "" {
		path, ok := msfstools.GetPackageFolderPath()

		if !ok {
			log.Fatal("MSFS package folder not found")
		} else {
			msfsPath = path
		}
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	// Installation
	for _, ce := range config.Config {
		switch ce.SourceType {
		case Local:
			dest := filepath.Join(msfsPath, ce.DestPath)

			if ce.CreateFile {
				log.Printf("Copying %s -> %s\n", ce.SourcePath, dest)
				utils.CopyFile(ce.SourcePath, dest)
			} else {
				log.Printf("Checking %s -> %s\n", ce.SourcePath, dest)
				addtxtcontentgo.AddContent(ce.SourcePath, dest)
			}

		case Http:
			tmpFilePath := utils.DownloadContent(ce.SourcePath)
			defer os.Remove(tmpFilePath)

			dest := filepath.Join(msfsPath, ce.DestPath)

			if ce.CreateFile {
				log.Printf("Copying %s -> %s\n", ce.SourcePath, dest)
				utils.CopyFile(tmpFilePath, dest)
			} else {
				log.Printf("Checking %s -> %s\n", ce.SourcePath, dest)
				addtxtcontentgo.AddContent(tmpFilePath, dest)
			}

		case Git:
			panic("Not implemented")

		default:
			log.Fatal("Invalid source type")
		}
	}
}
