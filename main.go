package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var files []string

	imports := make(map[string]bool)
	root := "C:\\Users\\eviatar\\go\\src\\myproj"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".go" {
			return nil
		}
		files = append(files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		file, err := os.Open(file)
		if err != nil {
			fmt.Printf("[*] Failed to open %s\n", file)
		}

		defer file.Close()
		reader := bufio.NewReader(file)

		haveMultipleImports := false
		endOfImports := false
		var myImport string

		for {

			if endOfImports {
				break
			}

			line, _, err := reader.ReadLine()

			if err == io.EOF {
				break
			}

			trimmedLine := strings.Replace(string(line), " ", "", -1)
			trimmedLine = strings.Replace(string(trimmedLine), "\t", "", -1)


			if haveMultipleImports {

				if strings.HasSuffix(trimmedLine, ")"){
					endOfImports = true
					myImport = strings.Replace(trimmedLine, ")", "", -1)
					if myImport != "" {
						if _, ok := imports[myImport]; !ok {
							imports[myImport] = true
						}
					}
					break
				} else {
					if _, ok := imports[trimmedLine]; !ok {
						imports[trimmedLine] = true
					}
				}


			} else {
				if strings.HasPrefix(trimmedLine, "import"){
					if strings.HasPrefix(trimmedLine, "import(") { // at least 1 import
						haveMultipleImports = true
						myImport = strings.Replace(string(trimmedLine), "import(", "", -1)
						if strings.HasSuffix(trimmedLine, ")") { // is exact 1 import
							myImport = strings.Replace(myImport, ")", "", -1)
							if _, ok := imports[myImport]; !ok {
								imports[myImport] = true
							}
							break
						} else {
							if myImport != "" {
								if _, ok := imports[myImport]; !ok {
									imports[myImport] = true
								}
							}
						}

					} else { //only 1 import
						myImport = strings.Replace(string(trimmedLine), "import", "", -1)
						if _, ok := imports[myImport]; !ok {
							imports[myImport] = true
						}
						break
					}
				}
			}

			//fmt.Printf("%s \n", line)
		}
	}

	fmt.Println("[*] Done")


}
