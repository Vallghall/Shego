package builtin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Vallghall/schego/pkg/mem"
)

// LoadDefinitions returns load primitive definition.
// The load function is created with lexer, parser, and evaluator access via closure.
func LoadDefinitions(loadFn func(string) (mem.Object, error)) []mem.Definition {
	return []mem.Definition{
		{Name: "load", Value: mem.NewPrimitive("load", 1, func(args []mem.Object) (mem.Object, error) {
			// First arg must be a string (file path)
			pathObj, err := mem.AsString(args[0])
			if err != nil {
				return nil, mem.NewTypeError(mem.TypeString, args[0].Type(), "load expects a string path")
			}
			path := pathObj.Value()

			// Search for file: first in current directory, then in core/builtin/
			var fileContent []byte
			var readErr error

			// Try current directory first
			fileContent, readErr = os.ReadFile(path)
			if readErr != nil {
				// Try library path: core/builtin/
				libPath := filepath.Join("core", "builtin", path)
				fileContent, readErr = os.ReadFile(libPath)
				if readErr != nil {
					return nil, fmt.Errorf("load: cannot find file %s (tried %s and %s)", path, path, libPath)
				}
			}

			// Load and evaluate the file
			return loadFn(string(fileContent))
		})},
	}
}
