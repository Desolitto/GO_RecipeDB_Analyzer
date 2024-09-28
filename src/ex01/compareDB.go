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

func CompareDatabases(oldCakes, newCakes Recipes) {
	oldCakesMap := make(map[string]Cake)
	newCakesMap := make(map[string]Cake)

	for _, cake := range oldCakes.Cakes {
		oldCakesMap[cake.Name] = cake
	}

	for _, cake := range newCakes.Cakes {
		newCakesMap[cake.Name] = cake
	}

	for _, cake := range newCakes.Cakes {
		if _, exists := oldCakesMap[cake.Name]; !exists {
			fmt.Printf("ADDED cake \"%s\"\n", cake.Name)
		}
	}

	for _, cake := range oldCakes.Cakes {
		if _, exists := newCakesMap[cake.Name]; !exists {
			fmt.Printf("REMOVED cake \"%s\"\n", cake.Name)
		}
	}

	for _, oldCake := range oldCakes.Cakes {
		if newCake, exists := newCakesMap[oldCake.Name]; exists {
			if oldCake.Time != newCake.Time {
				fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", oldCake.Name, newCake.Time, oldCake.Time)
			}
			compareIngredients(oldCake.Name, oldCake.Ingredients, newCake.Ingredients)
		}
	}
}

func compareIngredients(cakeName string, oldIngredients, newIngredients []Ingredient) {
	oldIngredientMap := make(map[string]Ingredient)
	newIngredientMap := make(map[string]Ingredient)

	for _, ingredient := range oldIngredients {
		oldIngredientMap[ingredient.Name] = ingredient
	}

	for _, ingredient := range newIngredients {
		newIngredientMap[ingredient.Name] = ingredient
	}
	for _, newIngredient := range newIngredients {
		if _, exists := oldIngredientMap[newIngredient.Name]; !exists {
			fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", newIngredient.Name, cakeName)
		}
	}

	for _, oldIngredient := range oldIngredients {
		if _, exists := newIngredientMap[oldIngredient.Name]; !exists {
			fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", oldIngredient.Name, cakeName)
		}
	}
	for name, oldIngredient := range oldIngredientMap {
		if newIngredient, exists := newIngredientMap[name]; exists {
			if oldIngredient.Count != newIngredient.Count {
				fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n",
					name, cakeName, newIngredient.Count, oldIngredient.Count)
			}
			if oldIngredient.Unit != newIngredient.Unit {
				if oldIngredient.Unit == "" {
					fmt.Printf("ADDED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", newIngredient.Unit, name, cakeName)
				} else if newIngredient.Unit == "" {
					fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", oldIngredient.Unit, name, cakeName)
				} else {
					fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n",
						name, cakeName, newIngredient.Unit, oldIngredient.Unit)
				}
			}
		}
	}
}

func main() {
	oldFilePath := flag.String("old", "", "Path to the old database file")
	newFilePath := flag.String("new", "", "Path to the new database file")
	flag.Parse()

	if *oldFilePath == "" || *newFilePath == "" {
		fmt.Println("Usage: --old <old_file> --new <new_file>")
		return
	}

	oldReader, err := getReader(*oldFilePath)
	HandleError(err)

	newReader, err := getReader(*newFilePath)
	HandleError(err)

	oldCakes, err := oldReader.Read(*oldFilePath)
	HandleError(err)

	newCakes, err := newReader.Read(*newFilePath)
	HandleError(err)

	CompareDatabases(oldCakes, newCakes)
}

/*
go build compareDB.go
./compareDB --old ../data/original_database.xml --new ../data/stolen_database.json
go run compareDB.go --old ../data/original_database.xml --new ../data/stolen_database.json
*/
/*
ADDED cake "Moonshine Muffin"
REMOVED cake "Blueberry Muffin Cake"
CHANGED cooking time for cake "Red Velvet Strawberry Cake" - "45 min" instead of "40 min"
ADDED ingredient "Coffee beans" for cake  "Red Velvet Strawberry Cake"
REMOVED ingredient "Vanilla extract" for cake  "Red Velvet Strawberry Cake"
CHANGED unit for ingredient "Flour" for cake  "Red Velvet Strawberry Cake" - "mugs" instead of "cups"
CHANGED unit count for ingredient "Strawberries" for cake  "Red Velvet Strawberry Cake" - "8" instead of "7"
REMOVED unit "pieces" for ingredient "Cinnamon" for cake  "Red Velvet Strawberry Cake"
*/ /*
ADDED cake "Moonshine Muffin"
REMOVED cake "Blueberry Muffin Cake"
CHANGED cooking time for cake "Red Velvet Strawberry Cake" - "45 min" instead of "40 min"
ADDED ingredient "Coffee Beans" for cake "Red Velvet Strawberry Cake"
REMOVED ingredient "Vanilla extract" for cake "Red Velvet Strawberry Cake"
CHANGED unit count for ingredient "Strawberries" for cake "Red Velvet Strawberry Cake" - "8" instead of "7"
REMOVED unit "pieces" for ingredient "Cinnamon" for cake "Red Velvet Strawberry Cake"
CHANGED unit count for ingredient "Flour" for cake "Red Velvet Strawberry Cake" - "2" instead of "3"
CHANGED unit for ingredient "Flour" for cake "Red Velvet Strawberry Cake" - "mugs" instead of "cups"
*/ /*
ADDED cake "Moonshine Muffin"
REMOVED cake "Blueberry Muffin Cake"
CHANGED cooking time for cake "Red Velvet Strawberry Cake" - "45 min" instead of "40 min"
ADDED ingredient "Coffee Beans" for cake "Red Velvet Strawberry Cake"
REMOVED ingredient "Vanilla extract" for cake "Red Velvet Strawberry Cake"
CHANGED unit count for ingredient "Flour" for cake "Red Velvet Strawberry Cake" - "2" instead of "3"
CHANGED unit for ingredient "Flour" for cake "Red Velvet Strawberry Cake" - "mugs" instead of "cups"
CHANGED unit count for ingredient "Strawberries" for cake "Red Velvet Strawberry Cake" - "8" instead of "7"
REMOVED unit "pieces" for ingredient "Cinnamon" for cake "Red Velvet Strawberry Cake"
*/
