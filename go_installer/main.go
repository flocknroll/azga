package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/flocknroll/azga/go_installer/addtxtcontent"
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

// Handle installation process from source to destination
func handleInstall(srcPath string, destPath string, createFile bool) {
	if createFile {
		log.Printf("Copying %s -> %s\n", srcPath, destPath)
		utils.CopyFile(srcPath, destPath)
	} else {
		log.Printf("Checking %s -> %s\n", srcPath, destPath)

		if !addtxtcontent.CheckContent(srcPath, destPath) {
			start, end, _, found := addtxtcontent.CheckDelimitedSection(destPath, "# AZGA DATA START", "# AZGA DATA END")
			if found {
				log.Println("  --> Old version found - deleting")
				addtxtcontent.DeleteLines(destPath, start, end)
			}
			log.Println("  --> Content not found - adding")
			addtxtcontent.AddContent(srcPath, destPath)
		} else {
			log.Println("  --> Content found")
		}
	}
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
			handleInstall(ce.SourcePath, dest, ce.CreateFile)

		case Http:
			tmpFilePath := utils.DownloadContent(ce.SourcePath)
			defer os.Remove(tmpFilePath)
			log.Printf("Downloaded %s -> %s\n", ce.SourcePath, tmpFilePath)

			dest := filepath.Join(msfsPath, ce.DestPath)
			handleInstall(tmpFilePath, dest, ce.CreateFile)

		case Git:
			panic("Not implemented")

		default:
			log.Fatal("Invalid source type")
		}
	}
}
