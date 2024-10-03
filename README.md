# Go RecipeDB Analyzer

## Overview

Go RecipeDB Analyzer is a command-line application for analyzing recipe databases stored in XML and JSON formats. This tool reads, compares, and assesses changes in different recipe databases, highlighting modifications such as added, removed, or altered recipes and ingredients.

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Getting Started](#getting-started)
4. [Usage](#usage)
    1. [Exercise 00: Reading](#exercise-00-reading)
    2. [Exercise 01: Assessing Damage](#exercise-01-assessing-damage)
    3. [Exercise 02: Afterparty](#exercise-02-afterparty)
5. [Project Structure](#project-structure)

## Introduction

There are many popular data formats in the world of programming, and Go makes it easy to work with them, particularly XML and JSON. This project simulates a bakery's recipe database stored in XML and compares it with a "stolen" database stored in JSON format. The application can detect various differences between the databases and highlight them.

## Features

- Read recipe databases in both XML and JSON formats.
- Compare two databases and detect changes, such as:
    - Added or removed cakes.
    - Changes in cooking time.
    - Added or removed ingredients.
    - Changes in ingredient quantities or units.
- Command-line interface for ease of use.

## Getting Started

### Prerequisites

- Go programming language installed (version 1.16 or higher recommended)
- Git for version control

### Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/Desolitto/GO_RecipeDB_Analyzer.git
    ```
2. Navigate to the project directory:
    ```bash
    cd GO_RecipeDB_Analyzer
    ```
3. Build the application:
    ```bash
    go build -o readDB ./src/ex00/readDB.go
    go build -o compareDB ./src/ex01/compareDB.go
    go build -o compareFS ./src/ex02/readDB/compareFS.go
    
    ```

## Usage

### Exercise 00: Reading

To read a database file and convert its format (from JSON to XML or vice versa), use the `readDB` command:

```
./readDB -f ./src/data/original_database.xml
./readDB -f ./src/data/stolen_database.json
```
or use Makefile in src/
```
make readDB
```

This will output the contents of the specified file in the opposite format (JSON will be converted to XML, and vice versa).

### Exercise 01: Assessing Damage

To compare the original database with the stolen one and identify changes, use the compareDB command:

```
./compareDB --old ./src/data/original_database.xml --new ./src/data/stolen_database.json
```
or use Makefile in src/
```
make compareDB
```

Expected output format:
```
ADDED cake "Moonshine Muffin"
REMOVED cake "Blueberry Muffin Cake"
CHANGED cooking time for cake "Red Velvet Strawberry Cake" - "45 min" instead of "40 min"
ADDED ingredient "Coffee Beans" for cake "Red Velvet Strawberry Cake"
REMOVED ingredient "Vanilla extract" for cake "Red Velvet Strawberry Cake"
CHANGED unit count for ingredient "Flour" for cake "Red Velvet Strawberry Cake" - "2" instead of "3"
CHANGED unit for ingredient "Flour" for cake "Red Velvet Strawberry Cake" - "mugs" instead of "cups"
CHANGED unit count for ingredient "Strawberries" for cake "Red Velvet Strawberry Cake" - "8" instead of "7"
REMOVED unit "pieces" for ingredient "Cinnamon" for cake "Red Velvet Strawberry Cake"
```

### Exercise 02: Afterparty

To compare two filesystem dumps, use the compareFS command:

```
./compareFS --old ./src/data/snapshot1.txt --new ./src/data/snapshot2.txt
```
or use Makefile in src/
```
make compareFS
```
This will output any added or removed files between the two snapshots.

## Project Structure

```graphql
GO_RecipeDB_Analyzer/
│
├── src/
│   ├── ex00/
│   │   └── readDB.go        # Main command for reading databases
│   ├── ex01/
│   │   └── compareDB.go/    # Main command for comparing recipe database
│   ├── ex02/                
│   │   ├── compareFS.go     # Main command for comparing file system snapshots
│   └── data/                # Sample XML, JSON data files, and snapshots for testing
│       ├── original_database.xml
│       ├── stolen_database.json
│       ├── snapshot1.txt    
│       ├── snapshot2.txt    
│       ├── snapshot3.txt    
│       └── snapshot4.txt    
│
├── .gitignore               # Specifies files and directories to be ignored by Git
├── Makefile                 # Commands for building and running the project
└── README.md                # Project documentation
```


### Explanation of the Directories

- **`go.mod`**: Go module file, which defines dependencies and module information.
- **`src/`**: Contains the source code for the project, organized into subdirectories for each feature.
- **`ex00/`**: Contains `readDB.go`, the main command for reading databases.
- **`ex01/`**: Contains `compareDB.go`, the main command for comparing recipe databases.
- **`ex02/`**: Contains `compareFS.go`, the main command for comparing file system snapshots.
- **`data/`**: Stores sample XML and JSON data files, as well as snapshots for testing.
- **`.gitignore`**: Specifies files and directories to be ignored by Git.
- **`Makefile`**: Contains commands for building and running the project.
- **`README.md`**: Provides documentation for the project, including setup instructions and usage guidelines.
