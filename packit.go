package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const version string = "0.0.1"

type archive struct {
	filename string
	excludes []string
}

type FlagSetRegistry struct {
	sets map[string]*flag.FlagSet
}

// .packit file, currently it will be used to add ignore files
// similar to gitignore, there is a command to add a file/dir
// to the packit file, which will simply append that name at the
// end of the .packit file.

func main() {

	flagSets := &FlagSetRegistry{sets: make(map[string]*flag.FlagSet)}

	// Flag set for building the archive
	archiveFlagSet := flag.NewFlagSet("build", flag.ExitOnError)
	flagSets.sets["build"] = archiveFlagSet
	outputFlag := archiveFlagSet.String("o", "", "name of the output file relative to current directory without any extension")
	// TODO:: will implement the format handler later, currenly the focus is on the zip only
	formatFlag := archiveFlagSet.String("f", "zip", "the format of archive you want to create zip/tar")

	// Flag set to handle ignores
	ignoreFlagSet := flag.NewFlagSet("ignore", flag.ExitOnError)
	ignoreListFlag := ignoreFlagSet.Bool("l", false, "Lists all the files which are in the ignore list")
	flagSets.sets["ignore"] = ignoreFlagSet

	if len(os.Args) == 1 {
		printLogo()
		printCommands(flagSets)
		return
	}

	currentDir, err := os.Getwd()

	if err != nil {
		errorMessage("Unable to get the current directory")
		return
	}

	// Building the file name based on the current dir
	currentBase := filepath.Base(currentDir)
	currentBase += "." + *formatFlag

	archive := &archive{filename: currentBase}

	if len(os.Args) > 1 && os.Args[1] == "ignore" {
		if len(os.Args) > 2 {
			ignoreFlagSet.Parse(os.Args[2:])
		}

		if *ignoreListFlag {
			if err := archive.loadExcludes(); err != nil {
				errorMessage(err.Error())
				return
			}

			infoMessage("List of excludes are:")
			for _, v := range archive.excludes {
				fmt.Println(v)
			}
			return
		}

		if err = archive.addExcludes(ignoreFlagSet.Args()); err != nil {
			errorMessage(err.Error())
			return
		}
		return
	}

	if len(os.Args) > 1 && os.Args[1] != "build" {
		fmt.Println("Try using packit build -help")
		errorMessage("Unknown command entered")
		return
	}

	if len(os.Args) > 2 {
		archiveFlagSet.Parse(os.Args[2:])
	}

	// Do we have a filename
	if *outputFlag != "" {
		archive.filename = strings.ReplaceAll(*outputFlag, ".", "_") + "." + *formatFlag
	}

	archive.loadExcludes()

	fmt.Println("Starting to build the archive")
	outPutZip, err := os.Create(archive.filename)

	if err != nil {
		errorMessage(err.Error())
		return
	}

	w := zip.NewWriter(outPutZip)
	defer w.Close()

	archiveDone := make(chan bool)
	go loader(archiveDone)

	archive.createZipFiles(w)
	archiveDone <- false
	successMessage("A zip file has been created " + archive.filename)
}

func (a *archive) createZipFiles(zipWriter *zip.Writer) {
	currentDir, _ := os.Getwd()

	err := filepath.WalkDir(currentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(currentDir, path)
		if err != nil {
			return err
		}

		// We dont want to zip the zip file itself
		if relPath == a.filename {
			return nil
		}

		// We won't archive our excludes file
		if relPath == ".packit" {
			return nil
		}

		// Skip the excludes
		if slices.Contains(a.excludes, relPath) {
			// If dir skip it all
			if d.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if d.IsDir() {
			return nil
		}

		sourceFileW, err := os.Open(path)

		if err != nil {
			return err
		}

		defer sourceFileW.Close()

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
			infoMessage(".packit file does not exists so no file will be ignored")
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

func (a *archive) addExcludes(excludes []string) error {

	if err := a.loadExcludes(); err != nil {
		return err
	}

	excludesToAdd := sliceDiff(a.excludes, excludes)

	currentDir, err := os.Getwd()

	if err != nil {
		return err
	}

	packitPath := filepath.Join(currentDir, ".packit")

	var ignoreString string

	for _, v := range excludesToAdd {
		ignoreString += v + "\n"
	}

	if ignoreString == "" {
		errorMessage("Nothing to add to ignore list")
		return nil
	}

	// TODO:: need to make it hidden on windows as well, as unlike unix based systems, windows does not makes
	// .files hidden
	f, err := os.OpenFile(packitPath, os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write([]byte(ignoreString)); err != nil {
		return err
	}

	fmt.Println(excludes)
	successMessage("Added to the .packit file")

	return nil
}
