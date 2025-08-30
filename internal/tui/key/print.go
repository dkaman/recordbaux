package key

import (
	"strings"

	"github.com/charmbracelet/bubbles/v2/key"
)

func FmtKeymap(bindings []key.Binding) string {
	var entries []string

	for _, binding := range bindings {
		keys, help := binding.Keys(), binding.Help().Desc
		k := strings.Join(keys, "/")
		entry := k + " : "  + help
		entries = append(entries, entry)
	}

	return strings.Join(entries, " | ")
}
