package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/quincy/scoutbook-tools/roster"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	rosterPath, outputPath := configureFlags()

	rosterFile, err := os.Open(rosterPath)
	if err != nil {
		fmt.Printf("Error opening roster file: %v\n", err)
		os.Exit(1)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("Error closing roster file:", err)
		}
	}(rosterFile)

	parser := roster.NewCsvParser()

	scoutbookUsers, err := parser.ParseAdultRoster(rosterFile)
	if err != nil {
		fmt.Printf("Error parsing roster: %v\n", err)
		os.Exit(1)
	}

	var users []user
	for _, sbu := range scoutbookUsers {
		users = append(users, createUser(sbu.FirstName, sbu.LastName, sbu.Email))
	}

	writeOutput(outputPath, users)
}

// removeDuplicates removes duplicate entries where name and email are the same (case-insensitive)
func removeDuplicates(users []user) []user {
	seen := make(map[string]bool)
	var result []user

	for _, user := range users {
		// Create a unique key using lowercase name and email
		key := strings.ToLower(user.Name + "," + user.Email)

		if !seen[key] {
			seen[key] = true
			result = append(result, user)
		}
	}

	return result
}

// sortUsers sorts users alphabetically by name
func sortUsers(users []user) {
	sort.Slice(users, func(i, j int) bool {
		return strings.ToLower(users[i].Name) < strings.ToLower(users[j].Name)
	})
}

// writeOutput writes the output the specified outputPath or to stdout if outputPath is nil
func writeOutput(outputPath string, users []user) {
	// Remove duplicates
	users = removeDuplicates(users)

	// Sort alphabetically
	sortUsers(users)

	var out io.Writer
	if outputPath == "" {
		out = os.Stdout
	} else {
		file, err := os.Create(outputPath)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer func(outputFile *os.File) {
			err := outputFile.Close()
			if err != nil {
				log.Println("Error closing output file:", err)
			}
		}(file)
		out = file
	}

	// Create CSV writer
	writer := csv.NewWriter(out)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Name", "Email"}); err != nil {
		fmt.Printf("Error writing CSV header: %v\n", err)
		os.Exit(1)
	}

	// Write user data
	for _, user := range users {
		if err := writer.Write([]string{user.Name, user.Email}); err != nil {
			fmt.Printf("Error writing user to CSV: %v\n", err)
			os.Exit(1)
		}
	}
}

// createUser is a helper function that creates a user from firstName, lastName, and email
func createUser(firstName, lastName, email string) user {
	return user{Name: fmt.Sprintf("%s %s", firstName, lastName), Email: email}
}

func configureFlags() (string, string) {
	// Define command line flags
	rosterPath := flag.String("roster", "", "Path to adult roster CSV file (required)")
	outputPath := flag.String("output", "", "Path to output CSV file, defaults to stdout")
	flag.Parse()

	// Validate required flags
	if *rosterPath == "" {
		fmt.Println("Error: roster path is required")
		flag.Usage()
		os.Exit(1)
	}
	return *rosterPath, *outputPath
}

type user struct {
	Name  string
	Email string
}
