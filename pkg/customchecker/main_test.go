// Модуль кастомного анализатора

package customchecker

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOsExit tests the OsExit function
func TestOsExit(t *testing.T) {
	src := `
package main

import "os"

func main() {
    os.Exit(0)
}
`
	// Create a new file set
	fset := token.NewFileSet()

	// Parse the source code
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	assert.NoError(t, err)

	// Call OsExit and get the position of os.Exit
	pos := OsExit(file)

	// Check that the position is not nil
	assert.NotNil(t, pos)
}
