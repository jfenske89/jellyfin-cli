package actions

import (
	"encoding/json"
	"flag"
	"fmt"
	"reflect"
)

// parseFlags accept command-level flags and parse with the global flags
func parseFlags(commandFlags map[string]any) map[string]any {
	flags := map[string]any{
		"json": flag.Bool("json", false, "Output in JSON format"),
	}

	for label := range commandFlags {
		flags[label] = commandFlags[label]
	}

	flag.Parse()

	options := make(map[string]any, 0)

	for label := range flags {
		if flags[label] == nil || !isPointer(flags[label]) {
			continue
		}

		val := reflect.ValueOf(flags[label]).Elem().Interface()

		if label == "json" {
			if val.(bool) {
				// global json flag sets an output option
				options["output"] = "json"
			}
		} else {
			options[label] = val
		}
	}

	return options
}

// isPointer checks if the given value is a pointer
func isPointer(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}

// writeResponse accepts a response, output format, and a textWriter and writes the response
func writeResponse[T any](response T, output string, textWriter func(T)) error {
	if output != "json" {
		textWriter(response)
		return nil
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	fmt.Println(string(jsonBytes))
	return nil
}
