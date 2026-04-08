package argparse

import (
	"fmt"
	"strings"
)

type Spec struct {
	TakesValue bool
}

type Parsed struct {
	Positionals []string
	Flags       map[string][]string
}

func (parsed Parsed) First(name string) string {
	values := parsed.Flags[name]
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (parsed Parsed) Has(name string) bool {
	return len(parsed.Flags[name]) > 0
}

// Parse accepts interleaved positionals and GNU-style flags in either
// `--name value` or `--name=value` form.
func Parse(args []string, specs map[string]Spec) (Parsed, error) {
	parsed := Parsed{
		Flags: make(map[string][]string),
	}

	for index := 0; index < len(args); index++ {
		arg := args[index]
		if !strings.HasPrefix(arg, "--") {
			parsed.Positionals = append(parsed.Positionals, arg)
			continue
		}

		nameValue := strings.TrimPrefix(arg, "--")
		parts := strings.SplitN(nameValue, "=", 2)
		name := parts[0]

		spec, ok := specs[name]
		if !ok {
			return Parsed{}, fmt.Errorf("unsupported flag: --%s", name)
		}

		if !spec.TakesValue {
			parsed.Flags[name] = append(parsed.Flags[name], "true")
			continue
		}

		if len(parts) == 2 {
			parsed.Flags[name] = append(parsed.Flags[name], parts[1])
			continue
		}

		index++
		if index >= len(args) {
			return Parsed{}, fmt.Errorf("flag --%s requires a value", name)
		}

		parsed.Flags[name] = append(parsed.Flags[name], args[index])
	}

	return parsed, nil
}
