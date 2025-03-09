package dd

import (
	"fmt"
	"os"
	"runtime"
)

// dd dumps the provided variables and stops the program
func DD(vars ...interface{}) {
	// Get the caller's file and line number
	_, file, line, _ := runtime.Caller(1)

	// Print the file and line number
	fmt.Printf("\n\ndd called in %s (line %d):\n", file, line)

	// Dump the variables
	for i, v := range vars {
		fmt.Printf("Argument %d:\n", i+1)
		fmt.Printf("%+v\n\n", v) // Use %+v to print structs with field names
	}

	// Exit the program
	os.Exit(1)
}