package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Ingredient struct {
	Name  string      `json:"ingredient_name" xml:"itemname"`
	Count json.Number `json:"ingredient_count" xml:"itemcount"`
	Unit  string      `json:"ingredient_unit,omitempty" xml:"itemunit,omitempty"`
}

type Cake struct {
	Name        string       `json:"name" xml:"name"`
	Time        string       `json:"time" xml:"stovetime"`
	Ingredients []Ingredient `json:"ingredients" xml:"ingredients>item"`
}

type Recipes struct {
	Cakes []Cake `json:"cake" xml:"cake"`
}

func HandleError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

type DBReader interface {
	Read(str string) (Recipes, error)
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return byteValue, nil
}

type JSONReader struct{}

func (r JSONReader) Read(filePath string) (Recipes, error) {
	var db Recipes
	byteValue, err := readFile(filePath)
	if err != nil {
		return db, err
	}
	cleanedJSON := removeJSONComments(byteValue)

	err = json.Unmarshal(cleanedJSON, &db)
	if err != nil {
		return db, errors.New("failed to deserialize JSON")
	}
	return db, nil
}

func removeJSONComments(data []byte) []byte {
	var buffer bytes.Buffer
	inComment := false

	for i := 0; i < len(data); i++ {
		if data[i] == '/' && i+1 < len(data) && data[i+1] == '/' {
			inComment = true
			i++
			continue
		}

		if inComment && (data[i] == '\n' || data[i] == '\r') {
			inComment = false
		}

		if !inComment {
			buffer.WriteByte(data[i])
		}
	}

	return buffer.Bytes()
}

type XMLReader struct{}

func (r XMLReader) Read(filePath string) (Recipes, error) {
	var db Recipes
	byteValue, err := readFile(filePath)
	if err != nil {
		return db, err
	}

	err = xml.Unmarshal(byteValue, &db)
	if err != nil {
		return db, errors.New("failed to deserialize XML")
	}

	return db, nil
}

func getReader(filePath string) (DBReader, error) {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".json":
		return JSONReader{}, nil
	case ".xml":
		return XMLReader{}, nil
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}
}

func printDB(db Recipes, format string) error {
	var data []byte
	var err error

	if format == "json" {
		data, err = json.MarshalIndent(db, "", "    ")
	} else if format == "xml" {
		data, err = xml.MarshalIndent(db, "", "    ")
	}
	if err != nil {
		return errors.New("failed to marshal data")
	}

	fmt.Println(string(data))
	return nil
}

func readPath() string {
	var f = flag.String("f", "", "read path")
	flag.Parse()
	return *f
}

func main() {
	filePath := readPath()
	if filePath == "" {
		fmt.Println("Please specify the file path using the -f option")
		return
	}
	reader, err := getReader(filePath)
	HandleError(err)

	db, err := reader.Read(filePath)
	HandleError(err)

	var format string
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".json" {
		format = "xml"
	} else if ext == ".xml" {
		format = "json"
	}
	err = printDB(db, format)
	HandleError(err)

}

/*
go build readDB.go
./readDB -f ../data/original_database.xml
./readDB -f ../data/stolen_database.json
go run readDB.go -f ../data/original_database.xml
go run readDB.go -f ../data/stolen_database.json
*/
