package msfstools

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Retrieves the MSFS packages installation folder
func GetPackageFolderPath() (string, bool) {
	appDataPath, found := os.LookupEnv("APPDATA")
	if !found {
		return "", false
	}

	msfsOptsPath := filepath.Join(appDataPath, "Microsoft Flight Simulator", "UserCfg.opt")

	f, err := os.OpenFile(msfsOptsPath, os.O_RDONLY, 0)
	if err != nil {
		defer f.Close()
	} else {
		panic(err)
	}

	reader := bufio.NewReader(f)

	var pkgPath string
	found = false
	for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
		if strings.HasPrefix(line, "InstalledPackagesPath") {
			parts := strings.SplitN(line, " ", 2)
			pkgPath = parts[1][1 : len(parts[1])-1] // We remove double quotes around the path
			found = true
			break
		}
	}

	return pkgPath, found
}
