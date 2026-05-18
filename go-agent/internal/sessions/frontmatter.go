package sessions

import "strings"

func parseFrontmatter(text string) (map[string]string, string) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "---") {
		return map[string]string{}, text
	}

	rest := text[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return map[string]string{}, text
	}

	fm := rest[:idx]
	body := strings.TrimSpace(rest[idx+4:])

	meta := make(map[string]string)
	for _, line := range strings.Split(fm, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			meta[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return meta, body
}
