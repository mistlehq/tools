package main

import (
	"fmt"
	"strings"
)

type argSpec struct {
	takesValue bool
}

type parsedArgs struct {
	positionals []string
	flags       map[string][]string
}

func (parsed parsedArgs) first(name string) string {
	values := parsed.flags[name]
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

// parseArgs accepts interleaved positionals and GNU-style flags in either
// `--name value` or `--name=value` form.
func parseArgs(args []string, specs map[string]argSpec) (parsedArgs, error) {
	parsed := parsedArgs{
		flags: make(map[string][]string),
	}

	for index := 0; index < len(args); index++ {
		arg := args[index]
		if !strings.HasPrefix(arg, "--") {
			parsed.positionals = append(parsed.positionals, arg)
			continue
		}

		nameValue := strings.TrimPrefix(arg, "--")
		parts := strings.SplitN(nameValue, "=", 2)
		name := parts[0]

		spec, ok := specs[name]
		if !ok {
			return parsedArgs{}, fmt.Errorf("unsupported flag: --%s", name)
		}

		if !spec.takesValue {
			parsed.flags[name] = append(parsed.flags[name], "true")
			continue
		}

		if len(parts) == 2 {
			parsed.flags[name] = append(parsed.flags[name], parts[1])
			continue
		}

		index++
		if index >= len(args) {
			return parsedArgs{}, fmt.Errorf("flag --%s requires a value", name)
		}

		parsed.flags[name] = append(parsed.flags[name], args[index])
	}

	return parsed, nil
}
