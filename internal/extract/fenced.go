package extract

import (
	"fmt"
	"strings"
)

func Block(raw, language string) (string, error) {
	startMarker := "```" + language
	lines := strings.Split(raw, "\n")
	inBlock := false
	var out []string
	for _, line := range lines {
		if !inBlock && strings.TrimSpace(line) == startMarker {
			inBlock = true
			continue
		}
		if inBlock && strings.TrimSpace(line) == "```" {
			return strings.TrimSpace(strings.Join(out, "\n")), nil
		}
		if inBlock {
			out = append(out, line)
		}
	}
	return "", fmt.Errorf("missing fenced block: %s", language)
}
