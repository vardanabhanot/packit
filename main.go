package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type archive struct {
	filename string
	excludes []string
}

// .packit file, currently it will be used to add ignore files
// similar to gitignore, there is a command to add a file/dir
// to the packit file, which will simply append that name at the
// end of the .packit file.

func main() {

	outputFlag := flag.String("o", "", "name of the output relative to current directory without any extension")
	formatFlag := flag.String("f", "zip", "the format of archive you want to create zip/tar")
	// TODO:: will implement the format handler later, currenly the focus is on the zip only

	flag.Parse()

	archive := &archive{filename: "test.zip"}

	// Do we have a filename
	if *outputFlag != "" {
		archive.filename = strings.ReplaceAll(*outputFlag, ".", "_") + "." + *formatFlag
		fmt.Println(archive)
	}

	archive.loadExcludes()

	outPutZip, err := os.Create(archive.filename)

	if err != nil {
		errorMessage(err.Error())
		return
	}

	w := zip.NewWriter(outPutZip)
	defer w.Close()

	//go func(w *zip.Writer) {
	archive.createZipFiles(w)
	//}(w)
	successMessage("A zip file has been created " + archive.filename)
}

func (a *archive) createZipFiles(zipWriter *zip.Writer) {
	currentDir, _ := os.Getwd()

	err := filepath.WalkDir(currentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		sourceFileW, err := os.Open(path)

		if err != nil {
			return err
		}

		defer sourceFileW.Close()

		relPath, err := filepath.Rel(currentDir, path)
		if err != nil {
			return err
		}

		// We dont want to zip the zip file itself
		if relPath == a.filename {
			return nil
		}

		zipFileW, err := zipWriter.Create(relPath)

		if err != nil {
			return err
		}

		if _, err = io.Copy(zipFileW, sourceFileW); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		errorMessage(err.Error())
	}
}

func (a *archive) loadExcludes() error {
	currentDir, err := os.Getwd()

	if err != nil {
		fmt.Println("Unable to read the current dir")
	}

	packitPath := filepath.Join(currentDir, ".packit")

	if _, err = os.Stat(packitPath); err != nil {
		if os.IsNotExist(err) {
			infoMessage(".packit file does not exists so no file will be ignored\n")
			return nil
		} else {
			return err
		}
	}

	// Packit file just includes ingore folders so it
	// should not be to big to not be able to read in single go
	packitContent, err := os.ReadFile(packitPath)

	if err != nil {
		return err
	}

	packitContentString := string(packitContent)
	a.excludes = strings.Split(packitContentString, "\n")

	return nil
}

func errorMessage(msg string) {
	fmt.Printf("\x1b[0;31m%s\x1b[0m", msg)
}

func successMessage(msg string) {
	fmt.Printf("\x1b[0;32m%s\x1b[0m", msg)
}

func infoMessage(msg string) {
	fmt.Printf("\x1b[0;34m%s\x1b[0m", msg)
}

func loader() {
	var isStart bool
	for true {
		loaderState := []string{"|", "/", "-", "\\", "-"}
		for _, v := range loaderState {
			time.Sleep(200 * time.Millisecond)
			if v == "|" {
				if isStart {
					fmt.Println("\x1b[?25h" + v)
				} else {
					fmt.Printf("\x1b[1D%s", v)
				}

				continue
			}
			fmt.Printf("\x1b[1D%s", v)
		}

	}
}
