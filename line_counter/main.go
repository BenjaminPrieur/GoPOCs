package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(Info("Starting LineCounter"))
	fmt.Println(Fata("don't kill pinguin"))
	fmt.Println(Info(os.Args))
	args := os.Args[1]
	all := 0

	err := filepath.Walk(args, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && !strings.Contains(info.Name(), "test") {

			ext := filepath.Ext(path)
			if ext == ".go" {

				jj, _ := os.ReadFile(path)
				g, _ := LineCounter(string(jj))
				all += g

				var u string
				z := strconv.Itoa(g)
				if g < 200 {
					u = Info(z)
				} else if g >= 200 && g < 1000 {
					u = Warn(z)
				} else if g >= 1000 {
					u = Fata(z)
				}
				fmt.Printf("File %s have %s lines of code\n", path, u)
			}
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	fmt.Println("\nAll lines", Fata(all))
}

func LineCounter(r string) (int, error) {

	var count int
	rr := strings.Split(r, "\n")
	var rrr []string
	cc := 0
	comstart := false
	for _, str := range rr {
		if strings.HasPrefix(str, "/*") {
			comstart = true
			continue
		}

		if strings.Contains(str, "*/") {
			comstart = false
			cc += 1
		}

		if str != "" && !strings.HasPrefix(str, "//") && !comstart {
			rrr = append(rrr, str)
		}
	}

	count = len(rrr) - cc

	return count, nil
}
