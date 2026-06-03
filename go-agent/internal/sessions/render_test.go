package sessions

import (
	"strings"
	"testing"
)

func TestGetDescriptionsRendersStructuredSkillSummaries(t *testing.T) {
	loader := NewSkillLoader()
	loader.LoadFromEntries(map[string]*SkillEntry{
		"demo": {
			Meta: map[string]string{
				"description": "Demo skill",
				"tags":        "demo,example",
			},
		},
	})

	descriptions := loader.GetDescriptions()

	for _, want := range []string{
		"<skills>",
		"<skill>",
		"<name>demo</name>",
		"<description>Demo skill</description>",
		"<tags>demo,example</tags>",
		"</skill>",
		"</skills>",
	} {
		if !strings.Contains(descriptions, want) {
			t.Fatalf("expected descriptions to contain %q, got %q", want, descriptions)
		}
	}
}

func TestGetDescriptionsProvidesDefaultDescriptionAndEscapesXML(t *testing.T) {
	loader := NewSkillLoader()
	loader.LoadFromEntries(map[string]*SkillEntry{
		"unsafe <skill>": {
			Meta: map[string]string{
				"description": "",
				"tags":        "a & b",
			},
		},
	})

	descriptions := loader.GetDescriptions()

	for _, want := range []string{
		"<name>unsafe &lt;skill&gt;</name>",
		"<description>No description provided.</description>",
		"<tags>a &amp; b</tags>",
	} {
		if !strings.Contains(descriptions, want) {
			t.Fatalf("expected descriptions to contain %q, got %q", want, descriptions)
		}
	}
}
