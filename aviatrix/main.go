package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: terratest-converter <test-file>")
		os.Exit(1)
	}

	testFile := os.Args[1]
	testFileContents, err := ioutil.ReadFile(testFile)
	if err != nil {
		fmt.Printf("Error reading test file %s: %s", testFile, err)
		os.Exit(1)
	}

	re := regexp.MustCompile(`func\s+(\(.*?\)\s+)?(\w+)\(.*?\)`)
	terratestFileContents := re.ReplaceAllStringFunc(string(testFileContents), func(match string) string {
		// Extract the function name and signature from the matched string
		re := regexp.MustCompile(`func\s+(\(.*?\)\s+)?(\w+)\(.*?\)`)
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}
		signature := submatches[1]
		functionName := submatches[2]

		// Create the Terratest test function
		terratestTestFunc := fmt.Sprintf("func Test%s(t *testing.T) {\n\t%s(t)\n}\n", functionName, functionName)

		// Add the signature to the Terratest test function
		if signature != "" {
			terratestTestFunc = fmt.Sprintf("func %s Test%s(t *testing.T) {\n\t%s(t)\n}\n", signature, functionName, functionName)
		}

		return terratestTestFunc
	})

	// Write the Terratest file
	terratestFileName := strings.TrimSuffix(testFile, ".go") + "_terratest.go"
	err = ioutil.WriteFile(terratestFileName, []byte(terratestFileContents), 0644)
	if err != nil {
		fmt.Printf("Error writing Terratest file %s: %s", terratestFileName, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted test file %s to Terratest file %s\n", testFile, terratestFileName)
}
