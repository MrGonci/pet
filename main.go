package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Pattern struct {
	Name  string `json:"name"`
	Regex string `json:"regex"`
}

var Red = "\033[31m"
var Green = "\033[32m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Reset = "\033[0m"

type Config struct {
	Patterns []Pattern `json:"patterns"`
}

func main() {
	dir := flag.String("d", ".", "Directory to search for files.")
	verbose := flag.Bool("v", false, "Verbose mode.")
	output := flag.String("o", "", "Output file to save the results.")
	flag.Parse()

	fmt.Println(Blue + ` 
		|\---/|
		| o_o |	   Pattern Enumeration Tool
		 \_^_/ 
	`)

	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf(Red+"Failed to load config: %v", err)
	}

	var outFile *os.File
	var omitColors bool
	if *output != "" {
		outFile, err = os.Create(*output)
		if err != nil {
			log.Fatalf(Red+"Failed to create output file: %v", err)
		}
		defer outFile.Close()
		omitColors = true
		fmt.Printf(Green+"Saved on file: %s\n", *output)
	} else {
		outFile = os.Stdout
	}

	fmt.Fprintf(outFile, Cyan+"Searching patterns in directory '%s'\n", *dir)

	err = searchPatternsInDir(*dir, config, *verbose, outFile, omitColors)
	if err != nil {
		log.Fatalf(Red+"Failed to search patterns: %v", err)
	}
}

func loadConfig(filename string) (*Config, error) {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func searchPatternsInDir(dir string, config *Config, verbose bool, outFile *os.File, omitColors bool) error {
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(outFile, Red+"Error accessing path %s: %v\n", path, err)
			return nil
		}
		if !d.IsDir() {
			err := searchPatternsInFile(path, config, verbose, outFile, omitColors)
			if err != nil {
				fmt.Fprintf(outFile, "Failed to process file %s: %v\n", path, err)
			}
		}
		return nil
	})

	fmt.Println(Reset)

	return err
}

func searchPatternsInFile(filePath string, config *Config, verbose bool, outFile *os.File, omitColors bool) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := string(fileData)
	lineList := splitLines(lines)

	for _, pattern := range config.Patterns {
		re, err := regexp.Compile(pattern.Regex)
		if err != nil {
			return err
		}
		matches := re.FindAllStringIndex(string(fileData), -1)
		filePath = strings.Replace(filepath.ToSlash(filePath), "//", "/", -1)

		if len(matches) > 0 {
			header := fmt.Sprintf(Magenta+"\n[!] Pattern Found '%s' in %s\n", pattern.Name, filePath)
			if verbose {
				header = fmt.Sprintf(Magenta+"\nThe file %s has the pattern '%s' with %d matches\n", filePath, pattern.Name, len(matches))
			}
			writeOutput(outFile, header, omitColors)

			for _, match := range matches {
				line := findLineFromIndex(lineList, match[0])
				matchDetails := fmt.Sprintf(Green+"  Line %d: %s\n", line, string(fileData[match[0]:match[1]]))
				writeOutput(outFile, matchDetails, omitColors)
			}
		}
	}

	return nil
}

func splitLines(text string) []string {
	return strings.Split(text, "\n")
}

func findLineFromIndex(lines []string, index int) int {
	lineStart := 0
	for lineNumber, line := range lines {
		lineLength := len(line) + 1
		if index < lineStart+lineLength {
			return lineNumber + 1
		}
		lineStart += lineLength
	}
	return -1
}

func writeOutput(outFile *os.File, text string, omitColors bool) {
	if omitColors {
		text = stripColors(text)
	}
	fmt.Fprint(outFile, text)
}

func stripColors(text string) string {
	colorRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return colorRegex.ReplaceAllString(text, "")
}
