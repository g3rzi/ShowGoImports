package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ImportMap map[string]bool

func main() {
	var directoryFlag string
	flag.StringVar(&directoryFlag, "dir", "d", "Specify the Go project directory")
	flag.Parse()

	//directoryFlag = "C:\\Users\\eviatar\\go\\src\\mymy"

	if _, err := os.Stat(directoryFlag); os.IsNotExist(err) {
		fmt.Println("[*] No directory was specified, exiting.")
		os.Exit(1)
	}

	var files []string

	imports := make(ImportMap)

	root := directoryFlag
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
		weirdImportSyntax := false
		endOfImports := false
		var myImport string

		for {
			line, _, err := reader.ReadLine()

			if err == io.EOF {
				break
			}

			if endOfImports {
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


			} else if weirdImportSyntax {
				if strings.HasPrefix(trimmedLine, "(") {

				}

			}else {
				if strings.HasPrefix(trimmedLine, "import"){
					if strings.HasPrefix(trimmedLine, "import(") { // at least 1 import
						addImportFromMultipleImportType(&haveMultipleImports, trimmedLine, &imports)
/*
						haveMultipleImports = true
						myImport = strings.Replace(string(trimmedLine), "import(", "", -1)
						if strings.HasSuffix(trimmedLine, ")") { // is exact 1 import
							myImport = strings.Replace(myImport, ")", "", -1)
							if _, ok := imports[myImport]; !ok {
								imports[myImport] = true
							}
							haveMultipleImports = false
						} else {
							if myImport != "" {
								if _, ok := imports[myImport]; !ok {
									imports[myImport] = true
								}
							}
						}
*/
					} else if strings.HasPrefix(trimmedLine, "import\"") { //only 1 import
						myImport = strings.Replace(string(trimmedLine), "import", "", -1)
						if _, ok := imports[myImport]; !ok {
							imports[myImport] = true
						}
					} else {
						weirdImportSyntax = true
					}

				}
			}

			//fmt.Printf("%s \n", line)
		}
	}

	//fmt.Println("[*] Done")
	for importFromFile := range imports {
		fmt.Printf("%s\n", importFromFile)
	}

}

func addImportFromMultipleImportType(haveMultipleImports *bool, trimmedLine string, imports *ImportMap) {
	*haveMultipleImports = true
	myImport := strings.Replace(string(trimmedLine), "import(", "", -1)
	if strings.HasSuffix(trimmedLine, ")") { // is exact 1 import
		myImport = strings.Replace(myImport, ")", "", -1)
		if _, ok := (*imports)[myImport]; !ok {
			(*imports)[myImport] = true
		}
		*haveMultipleImports = false
	} else {
		if myImport != "" {
			if _, ok := (*imports)[myImport]; !ok {
				(*imports)[myImport] = true
			}
		}
	}
}
