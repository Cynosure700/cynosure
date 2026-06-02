package tools

import (
	"context"
	"fmt"
	"html"
	"sort"
	"strings"

	"nano_cc/internal/sessions"
)

type SkillSnapshot struct {
	UserSkills  *sessions.SkillLoader
	LocalSkills *sessions.SkillLoader
	Merged      *sessions.SkillLoader
}

type LoadedSkill struct {
	Name   string
	Source string
	Entry  *sessions.SkillEntry
}

const skillSnapshotContextKey contextKey = "skill_snapshot"

func NewSkillSnapshot(userSkills, localSkills *sessions.SkillLoader) *SkillSnapshot {
	return &SkillSnapshot{
		UserSkills:  userSkills,
		LocalSkills: localSkills,
		Merged:      sessions.MergeSkillLoaders(localSkills, userSkills),
	}
}

func WithSkillSnapshot(ctx context.Context, snapshot *SkillSnapshot) context.Context {
	return context.WithValue(ctx, skillSnapshotContextKey, snapshot)
}

func SkillSnapshotFromContext(ctx context.Context) (*SkillSnapshot, bool) {
	snapshot, ok := ctx.Value(skillSnapshotContextKey).(*SkillSnapshot)
	return snapshot, ok
}

func (s *SkillSnapshot) LoadSkill(name string) (LoadedSkill, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return LoadedSkill{}, fmt.Errorf("skill name is required")
	}
	if s == nil {
		return LoadedSkill{}, fmt.Errorf("no capabilities are available in this conversation")
	}
	if s.UserSkills != nil {
		entry, err := s.UserSkills.GetEntry(name)
		if err == nil {
			return LoadedSkill{Name: name, Source: "db", Entry: entry}, nil
		}
	}
	if s.LocalSkills != nil {
		entry, err := s.LocalSkills.GetEntry(name)
		if err == nil {
			return LoadedSkill{Name: name, Source: "local", Entry: entry}, nil
		}
	}
	return LoadedSkill{}, fmt.Errorf("unknown skill %q. Available: %s", name, strings.Join(s.availableSkillNames(), ", "))
}

func (s *SkillSnapshot) availableSkillNames() []string {
	seen := map[string]struct{}{}
	for _, loader := range []*sessions.SkillLoader{s.UserSkills, s.LocalSkills} {
		if loader == nil {
			continue
		}
		for name := range loader.Entries() {
			seen[name] = struct{}{}
		}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func handleLoadSkill(ctx context.Context, args map[string]any) (string, error) {
	skillName, _ := args["name"].(string)
	snapshot, _ := SkillSnapshotFromContext(ctx)
	loaded, err := snapshot.LoadSkill(skillName)
	if err != nil {
		return "", err
	}
	return renderLoadedSkill(loaded), nil
}

func renderLoadedSkill(loaded LoadedSkill) string {
	entry := loaded.Entry
	if entry == nil {
		entry = &sessions.SkillEntry{}
	}
	metadata := make([]string, 0, len(entry.Meta))
	keys := make([]string, 0, len(entry.Meta))
	for key := range entry.Meta {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := strings.TrimSpace(entry.Meta[key])
		if value == "" {
			continue
		}
		metadata = append(metadata, key+": "+value)
	}
	return fmt.Sprintf("<skill source=\"%s\" name=\"%s\">\n<metadata>\n%s\n</metadata>\n<content>\n%s\n</content>\n</skill>",
		html.EscapeString(loaded.Source),
		html.EscapeString(loaded.Name),
		html.EscapeString(strings.Join(metadata, "\n")),
		entry.Body,
	)
}
