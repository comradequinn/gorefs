package cmds

import (
	"bufio"
	"gorefs/refs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//GenRefs finds existing references in .go files for which gorefs has an alias configured
func GenRefs(grjc *refs.Content, root string) error {
	rs := make(chan string, 100)
	modified := false

	go findRefs(grjc.Aliases, root, rs)

	for r := range rs {
		if grjc.References[r] == "" {
			modified = true
			grjc.References[r] = "main"
			log.Printf("added %v to gorefs.json with default version of 'main'", r)
		}
	}

	if modified {
		refs.Write(grjc)
	}

	return nil
}

func findRefs(aliases map[string]string, approot string, refs chan string) {
	files := make(chan string, 100)
	done := make(chan bool)

	go getGoFiles(approot, files, done)

	go func() {
		defer close(refs)

		wg := sync.WaitGroup{}

		for file := range files {
			wg.Add(1)
			f := file

			go func() {
				defer wg.Done()
				parseRefs(f, aliases, refs)
			}()
		}

		wg.Wait()
	}()
}

func getGoFiles(approot string, files chan string, done chan bool) {
	defer func() {
		close(files)
		done <- true
	}()

	filepath.Walk(approot, func(filepath string, file os.FileInfo, err error) error {
		if err == nil {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
				files <- filepath
			}
		} else {
			log.Printf("Unable to read file. %v", err)
		}

		return nil
	})
}

func parseRefs(file string, aliases map[string]string, refs chan string) {

	f, err := os.Open(file)

	if err == nil {
		defer f.Close()

		s := bufio.NewScanner(f)

		for s.Scan() {
			text := strings.TrimSpace(strings.ToLower(s.Text()))

			for alias := range aliases {
				if sidx := strings.Index(text, "\""+alias); sidx >= 0 {
					ref := strings.TrimSpace(string([]rune(text)[sidx+1:]))

					if eidx := strings.Index(ref, "\""); eidx >= 0 {
						refs <- strings.TrimSpace(string([]rune(ref)[:eidx]))
					}
				}
			}
		}
	} else {
		log.Printf("Unable to open %v. %v", file, err)
	}
}
