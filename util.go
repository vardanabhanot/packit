package main

import (
	"flag"
	"fmt"
	"time"
)

func errorMessage(msg string) {
	fmt.Printf("\x1b[0;31m%s\x1b[0m\n", msg)
}

func successMessage(msg string) {
	fmt.Printf("\x1b[0;32m%s\x1b[0m\n", msg)
}

func infoMessage(msg string) {
	fmt.Printf("\x1b[0;34m%s\x1b[0m\n", msg)
}

func loader(archiveDone <-chan bool) {
	var isStart bool
	loaderState := []string{"|", "/", "-", "\\", "-"}
	fmt.Print("\x1b[?25l")
	for {
		select {
		case <-archiveDone:
			fmt.Print("\x1b[1D\x1b[K") // Move back and clean up the loader
			fmt.Print("\x1b[?25h")     // Showing back the cursor
			return

		default:
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
}

// Currently this diffs just for slice of stings
// will look into other type when needed, currenly not needed in this project
// Func taken from https://github.com/mtdevs28080617/go-slice-diff/blob/main/slice_diff.go
func sliceDiff(sliceA []string, sliceB []string) []string {
	ma := make(map[string]struct{}, len(sliceA))
	var diffs []string

	for _, ka := range sliceA {
		ma[ka] = struct{}{}
	}

	for _, kb := range sliceB {
		if _, ok := ma[kb]; !ok {
			diffs = append(diffs, kb)
		}
	}

	return diffs
}
func printLogo() {

	const Logo = `
 ____            _    _ _
|  _ \ __ _  ___| | _(_) |_
| |_) / _' |/ __| |/ / | __|
|  __/ (_| | (__|   <| | |_
|_|   \__,_|\___|_|\_\_|\__|`

	fmt.Println(Logo)
	fmt.Println("Packit an archive creator")

}

func printCommands(flagSets *FlagSetRegistry) {
	fmt.Println("\nUsage: packit <command> [flags] [arguments]")
	fmt.Println("\nAvailable Commands:")
	for i, v := range flagSets.sets {
		fmt.Println("\n packit " + i)

		v.VisitAll(func(f *flag.Flag) {
			fmt.Printf("\t-%s\t(Default: %s)\t%s\n", f.Name, f.DefValue, f.Usage)
		})

		if i == "ignore" {
			fmt.Println("\n\tAdding files to ignore")
			fmt.Println("\tpackit ignore file1 file2 file3")
		}
	}
}
