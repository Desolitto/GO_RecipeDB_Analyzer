package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func HandleError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func CompareFileSystems(oldFilePath, newFilePath string) error {
	oldFiles := make(map[string]struct{})

	file, err := os.Open(oldFilePath)
	if err != nil {
		return fmt.Errorf("error opening old file %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		oldFiles[scanner.Text()] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading old file %v", err)
	}

	file, err = os.Open(newFilePath)
	if err != nil {
		return fmt.Errorf("error opening new file %v", err)
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if _, exists := oldFiles[line]; !exists {
			fmt.Printf("ADDED %s\n", line)
		}
		delete(oldFiles, line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading new file %v", err)
	}

	for file := range oldFiles {
		fmt.Printf("REMOVED %s\n", file)
	}
	return nil
}

func main() {
	oldFilePath := flag.String("old", "", "Path to the old database file")
	newFilePath := flag.String("new", "", "Path to the new database file")
	flag.Parse()

	if *oldFilePath == "" || *newFilePath == "" {
		fmt.Println("Usage: --old <old_file> --new <new_file>")
		return
	}

	err := CompareFileSystems(*oldFilePath, *newFilePath)
	HandleError(err)
}

/*
go build compareFS.go
./compareFS --old snapshot1.txt --new snapshot2.txt
go run compareFS.go --old snapshot1.txt --new snapshot2.txt
*/
