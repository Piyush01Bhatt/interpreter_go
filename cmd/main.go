package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Read and print each line
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")                      // Display prompt
		input, err := reader.ReadString('\n') // Read input until Enter (newline)
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input) // Remove newlines/spaces
		if input == "exit" {             // Exit condition
			fmt.Println("Goodbye!")
			break
		}

		fmt.Println("You entered:", input) // Echo back input
	}
}

func main() {
	fmt.Println("This is the main function")
	args := os.Args[1:]
	fmt.Println(args)

	if len(args) > 1 {
		fmt.Println("(Usage: jlox [script])")
		os.Exit(64)
	}

	if len(args) == 1 {
		filepath := args[0]
		runFile(filepath)
		return
	}

	runPrompt()
}
