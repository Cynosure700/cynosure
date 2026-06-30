package tools

import (
	"context"
	"fmt"
	"html"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Cynosure700/cynosure/cynosure/internal/sessions"
)

type SkillSnapshot struct {
	UserSkills      *sessions.SkillLoader
	WorkspaceSkills *sessions.SkillLoader
	LocalSkills     *sessions.SkillLoader
	Merged          *sessions.SkillLoader

	// BuiltinMaterializer 在加载来源为 builtin 的 skill 时被调用，把该 skill 的
	// 整棵目录树落盘到 ~/.cynosure/system/skills/<name>/ 并返回落盘 base 目录。
	// 为空时表示不落盘（如单测场景），此时仅返回 skill 正文。
	BuiltinMaterializer func(name string) (baseDir string, err error)
}

type LoadedSkill struct {
	Name    string
	Source  string
	BaseDir string
	Entry   *sessions.SkillEntry
}

const skillSnapshotContextKey contextKey = "skill_snapshot"

func NewSkillSnapshot(userSkills, localSkills *sessions.SkillLoader) *SkillSnapshot {
	return &SkillSnapshot{
		UserSkills:      userSkills,
		WorkspaceSkills: localSkills,
		LocalSkills:     localSkills,
		Merged:          sessions.MergeSkillLoaders(userSkills, localSkills),
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
	if s.Merged != nil {
		entry, err := s.Merged.GetEntry(name)
		if err == nil {
			source := strings.TrimSpace(entry.Source)
			if source == "" {
				source = "local"
			}
			return LoadedSkill{Name: name, Source: source, BaseDir: s.resolveBaseDir(name, source, entry), Entry: entry}, nil
		}
	}
	return LoadedSkill{}, fmt.Errorf("unknown skill %q. Available: %s", name, strings.Join(s.availableSkillNames(), ", "))
}

// resolveBaseDir 计算 skill 的 base 目录（供模型定位 skill 资源），返回原始完整路径。
// 来源为 builtin 时，把整棵目录树落盘到 ~/.cynosure/system/skills/<name>/ 并以其为 base；
// 其它来源使用磁盘上 SKILL.md 的父目录。
func (s *SkillSnapshot) resolveBaseDir(name, source string, entry *sessions.SkillEntry) string {
	if source == "builtin" {
		if s.BuiltinMaterializer != nil {
			if baseDir, err := s.BuiltinMaterializer(name); err == nil && strings.TrimSpace(baseDir) != "" {
				return baseDir
			}
		}
		return ""
	}
	if entry == nil {
		return ""
	}
	path := strings.TrimSpace(entry.Path)
	if path == "" {
		return ""
	}
	return filepath.Dir(path)
}

func (s *SkillSnapshot) availableSkillNames() []string {
	seen := map[string]struct{}{}
	for _, loader := range []*sessions.SkillLoader{s.Merged, s.UserSkills, s.WorkspaceSkills, s.LocalSkills} {
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
	baseDirLine := ""
	if baseDir := strings.TrimSpace(loaded.BaseDir); baseDir != "" {
		baseDirLine = "\nBase directory for this skill: " + baseDir
	}
	return fmt.Sprintf("<skill source=\"%s\" name=\"%s\">\n<metadata>\n%s\n</metadata>\n<content>\n%s\n</content>%s\n</skill>",
		html.EscapeString(loaded.Source),
		html.EscapeString(loaded.Name),
		html.EscapeString(strings.Join(metadata, "\n")),
		entry.Body,
		baseDirLine,
	)
}
