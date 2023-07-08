package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Prompt for the path to check
	fmt.Print("Enter the path for the files you want to collect: ")
	pathToCheck, _ := reader.ReadString('\n')
	pathToCheck = strings.TrimSpace(pathToCheck)

	if pathToCheck == "" {
		fmt.Println("No path provided. Exiting...")
		os.Exit(1)
	}

	// Prompt for the extensions to consider
	fmt.Print("\nEnter the extensions to consider (comma-separated), if nothing is entered, all files will be considered: ")
	extensionsInput, _ := reader.ReadString('\n')
	extensionsInput = strings.TrimSpace(extensionsInput)

	var extensions []string
	if extensionsInput != "" {
		extensions = strings.Split(extensionsInput, ",")
	}

	// Remove leading and trailing whitespaces from extensions
	for i := range extensions {
		extensions[i] = strings.TrimSpace(extensions[i])
	}

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting user home directory:", err)
	}
	// csv file path should be $Home/Downloads, and it is named in the format pathToCheck_{DD-MM-YY}.csv
	fileName := fmt.Sprintf("%s_%s.csv", filepath.Base(pathToCheck), time.Now().Format("02-01-2006"))
	csvFilePath := filepath.Join(homeDirectory, "Downloads", fileName)

	headers := []string{"Filename", "Extension"}

	file, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal("Error creating CSV file:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(headers)
	if err != nil {
		log.Fatal("Error writing CSV header:", err)
	}

	err = filepath.Walk(pathToCheck, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			extension := strings.ToLower(filepath.Ext(path))
			if len(extensions) == 0 {
				filename := info.Name()
				err := writer.Write([]string{filename, extension})
				if err != nil {
					return err
				}
			}

			// if an extension list is provided
			for _, ext := range extensions {
				if extension == ext {
					filename := info.Name()

					err := writer.Write([]string{filename, extension})
					if err != nil {
						return err
					}
					break
				}
			}
			fmt.Printf("%s written to %s\n", info.Name(), csvFilePath)
		}

		return nil
	})

	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Println("CSV file created successfully!")
	fmt.Println("CSV file path:", csvFilePath)
}
