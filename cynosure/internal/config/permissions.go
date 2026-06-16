package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// bypassPermissionsMode 是 permissions.defaultMode 的特殊取值：命中时所有命令直接放行。
const bypassPermissionsMode = "bypassPermissions"

// Permissions 是工作区级权限配置，存储在 <workspaceRoot>/.cynosure/settings.json 的 permissions 字段。
type Permissions struct {
	DefaultMode  string   `json:"defaultMode,omitempty"`
	AllowedRules []string `json:"allowedRules,omitempty"`
}

// IsBypass 返回 true 表示一切命令直接放行、无需审批。
func (p Permissions) IsBypass() bool {
	return p.DefaultMode == bypassPermissionsMode
}

// Allows 返回 rule 是否已在放行列表中（精确匹配）。
func (p Permissions) Allows(rule string) bool {
	rule = strings.TrimSpace(rule)
	if rule == "" {
		return false
	}
	for _, r := range p.AllowedRules {
		if strings.TrimSpace(r) == rule {
			return true
		}
	}
	return false
}

type workspaceSettings struct {
	Permissions Permissions `json:"permissions"`
}

// WorkspaceCynosureSettingsPath 返回工作区级 settings.json 路径。
func WorkspaceCynosureSettingsPath(workspaceRoot string) string {
	return filepath.Join(strings.TrimSpace(workspaceRoot), ".cynosure", "settings.json")
}

// LoadWorkspacePermissions 实时读取工作区 settings.json 的 permissions 字段。
// 文件不存在时返回零值（即非 bypass、无放行规则），不视为错误。
func LoadWorkspacePermissions(workspaceRoot string) (Permissions, error) {
	path := WorkspaceCynosureSettingsPath(workspaceRoot)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Permissions{}, nil
		}
		return Permissions{}, fmt.Errorf("read workspace settings %s: %w", path, err)
	}
	var settings workspaceSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return Permissions{}, fmt.Errorf("parse workspace settings %s: %w", path, err)
	}
	return settings.Permissions, nil
}

// AppendWorkspaceApprovalRule 把放行规则写入工作区 settings.json 的 permissions.allowedRules，
// 已存在则忽略。保留文件中的其它字段，文件/目录不存在则创建。
func AppendWorkspaceApprovalRule(workspaceRoot, rule string) error {
	rule = strings.TrimSpace(rule)
	if rule == "" {
		return nil
	}
	path := WorkspaceCynosureSettingsPath(workspaceRoot)
	raw := map[string]json.RawMessage{}
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("read workspace settings %s: %w", path, err)
		}
	} else if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parse workspace settings %s: %w", path, err)
	}

	var perms Permissions
	if existing, ok := raw["permissions"]; ok {
		if err := json.Unmarshal(existing, &perms); err != nil {
			return fmt.Errorf("parse workspace permissions %s: %w", path, err)
		}
	}
	if perms.Allows(rule) {
		return nil
	}
	perms.AllowedRules = append(perms.AllowedRules, rule)

	permsJSON, err := json.Marshal(perms)
	if err != nil {
		return fmt.Errorf("marshal permissions: %w", err)
	}
	raw["permissions"] = permsJSON
	out, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal workspace settings: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir workspace settings dir: %w", err)
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return fmt.Errorf("write workspace settings %s: %w", path, err)
	}
	return nil
}
