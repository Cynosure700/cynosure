package sessions

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

// MaterializeBuiltinSkill 把内置 skill 的整棵目录树从嵌入 FS 落盘到 destRoot/<name>/。
// 若目标目录已存在则跳过写入（不覆盖用户可能的本地改动），直接返回该目录。
// 返回落盘后的 base 目录（真实磁盘路径）。
func MaterializeBuiltinSkill(fsys fs.FS, name, destRoot string) (string, error) {
	if fsys == nil {
		return "", fmt.Errorf("builtin skill filesystem is nil")
	}
	if name == "" {
		return "", fmt.Errorf("builtin skill name is required")
	}
	dest := filepath.Join(destRoot, name)

	// 已存在则跳过：不覆盖用户对落盘副本的本地改动。
	if _, err := os.Stat(dest); err == nil {
		return dest, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat materialized skill dir %s: %w", dest, err)
	}

	info, err := fs.Stat(fsys, name)
	if err != nil {
		return "", fmt.Errorf("locate builtin skill %q: %w", name, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("builtin skill %q is not a directory", name)
	}

	if err := copyEmbeddedSubtree(fsys, name, dest); err != nil {
		// 失败时清理半成品目录，避免遗留不完整的落盘副本。
		_ = os.RemoveAll(dest)
		return "", err
	}
	return dest, nil
}

// copyEmbeddedSubtree 把 fsys 中 srcRoot 子树逐文件复制到 destDir，保留目录结构。
func copyEmbeddedSubtree(fsys fs.FS, srcRoot, destDir string) error {
	return fs.WalkDir(fsys, srcRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel := relFromFSRoot(srcRoot, p)
		target := filepath.Join(destDir, filepath.FromSlash(rel))
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create dir for %s: %w", target, err)
		}
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return fmt.Errorf("read embedded %s: %w", p, err)
		}
		if err := os.WriteFile(target, data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", target, err)
		}
		return nil
	})
}

// relFromFSRoot 返回 fs 路径 p 相对于 srcRoot 的部分（slash 分隔）。
func relFromFSRoot(srcRoot, p string) string {
	if p == srcRoot {
		return "."
	}
	rel := p[len(srcRoot):]
	for len(rel) > 0 && rel[0] == '/' {
		rel = rel[1:]
	}
	return path.Clean(rel)
}
