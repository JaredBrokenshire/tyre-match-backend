package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"tyre-match-backend/db/migrations/process"
)

var confirmationMessage = "\nAre you sure you want to update the database? Y or N:  "

func init() {
	flag.Bool("confirm", false, "should we ask for confirmation?")
}

func main() {
	confirmRequired := os.Args[0] == "--confirm"
	if confirmRequired == false {
		process.Run()
		return
	}

	// Confirmation is required
	if askForConfirmation() {
		log.Print()
		process.Run()
	}
}

func askForConfirmation() bool {
	fmt.Print(confirmationMessage)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Print(confirmationMessage)
		return askForConfirmation()
	}
}

// You might want to put the following two functions in a separate utility package.

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
