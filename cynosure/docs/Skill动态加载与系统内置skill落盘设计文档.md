# Skill 动态加载与系统内置 Skill 落盘设计文档

## 1. 背景与目标

当前 Cynosure 的 skill 机制在启动时一次性加载（`bootstrap.go:loadSkills`），把内置（嵌入二进制）、用户级（`~/.cynosure/skills`）、工作区级（`<workspace>/.cynosure/skills`）三类 skill 合并成一个静态 `SkillLoader`，挂在 `Service.BuiltinSkills` 上。系统提示词每轮虽然会重建，但 skill 列表来源是这个**启动时固定**的快照；`<skills>` 段落原先是系统提示词的顶层段落；`<memory>` 段落原先是系统提示词末尾的顶层段落；`load_skill` 返回的内容不含 skill 所在路径；`/skills` 命令会把每个 skill 的真实磁盘路径展示给用户。

本次需求要调整六处机制，且**只增不改、不影响其他现有策略**：

1. **每轮动态加载非系统 skill**：除系统内置 skill 外，每轮对话都实时重新加载一次用户级 / 工作区级 skill 到系统提示词，使得对话期间新增 / 修改 / 删除磁盘 skill 能即时生效（系统内置 skill 仍来自嵌入二进制，无需重载）。
2. **skill 段落移入 `<system-reminder>`**：`<skills>` 摘要段落从系统提示词顶层移动到临时 user `<system-reminder>` 内部。
3. **`load_skill` 返回 base 路径**：调用 `load_skill` 后，除返回 skill 正文外，还要显式标明该 skill 的 base 目录，例如 `Base directory for this skill: ~/.cynosure/skills/skill-creator`。
4. **系统内置 skill 落盘 + 路径标明**：当 `load_skill` 加载的是「系统内置」skill 时，把该 skill 的**整个目录树**落盘到 `~/.cynosure/system/skills/<name>/`（若目录已存在则跳过写入），并在返回内容里以该落盘路径作为 base 目录。
5. **`/skills` 只显示来源**：TUI 中 `/skills` 输出不再展示每个 skill 的磁盘 `path`，只展示其来源类型（当前项目 / 家目录 / 系统内置 三种）。
6. **memory 段落移入 `<system-reminder>`**：`<memory>` 段落从系统提示词末尾的顶层移动到临时 user `<system-reminder>` 内部（仅调整渲染位置，记忆子系统机制不变）。

### 非目标（保持不变）

- 不改变 skill 的合并优先级（workspace > user > builtin）。
- 不改变 `load_skill` 工具的 schema、工具名、注册方式。
- 不改变 frontmatter 解析、`canonicalSkillName`、`GetDescriptions` 的 XML 结构。
- 不改变子 agent 的 skill 隔离语义（子 agent 复用父快照）。
- 不改变记忆的提取 / 选择 / 断点等子系统机制——本次仅调整 memory 段落的**渲染位置**（移入临时 user `<system-reminder>`），记忆内容的生成、注入逻辑与守卫文案（`memorySectionGuidance`）均保持不变。
- 不改变压缩、审批、超时等其他运行期机制。

## 2. 现状梳理（关键代码位置）

| 关注点 | 位置 | 现状 |
| --- | --- | --- |
| 启动时加载并合并三类 skill | `internal/local/bootstrap.go:118 loadSkills` | builtin(FS, `source="builtin"`) + user(`source="user"`) + workspace(`source="workspace"`) 合并成一个 `SkillLoader`，挂到 `Service.BuiltinSkills` |
| 每轮构建 skill 快照 | `internal/agent/runtime/prompt_builder.go:60 buildSkillSnapshot` | 直接返回 `NewSkillSnapshot(nil, s.BuiltinSkills)`，**静态** |
| 每轮构建系统提示词 | `internal/agent/runtime/conversation_flow.go:95-100` | 每轮调用 `buildSkillSnapshot` + `buildSystemPrompt` |
| `<skills>` 段落渲染 | `internal/assistant/prompt.go:190-202` | 顶层段落，位于 `<tools>` 之后 |
| `<memory>` 段落渲染 | `internal/assistant/prompt.go:208-210` | 顶层段落，位于系统提示词末尾（`<git-status>` 之后）；含 `memorySectionGuidance` + 记忆正文 |
| `<system-reminder>` 段落渲染 | `internal/assistant/prompt.go:215-224 renderSystemReminder` | 含 `current_day` 与 cynosureMd 上下文 |
| `load_skill` 返回内容 | `internal/tools/load_skill.go:95 renderLoadedSkill` | 返回 `<skill source name><metadata><content>`，**无路径** |
| skill 来源标记 | `internal/sessions/skill.go:18 SkillEntry.Source` | `builtin` / `user` / `workspace` |
| `/skills` 渲染 | `internal/tui/app.go:639 renderSkillDetails` | 展示 name + `[source]` + description + `path:` |
| 内置 skill 嵌入 FS | `assets/embed.go:33 BuiltinSkillsFS` | `//go:embed all:skills` 子树，含 scripts/references/assets/agents |
| 家目录路径助手 | `internal/config/local_config.go:113 CynosureSkillsDir` | `~/.cynosure/skills` |

来源字符串目前有三种取值：`builtin`、`user`、`workspace`，分别由 `bootstrap.go:loadSkills` 设定。

## 3. 设计方案

### 3.1 总体策略

采用「**运行期持有 skill 加载所需的目录配置 + 每轮重载非内置 skill**」方案：

- 在 `Service` 上新增字段保存「内置 skill 加载器（嵌入 FS，启动时固定）」与「用户 / 工作区 skill 目录路径」。每轮 `buildSkillSnapshot` 时，**重新**从用户目录、工作区目录扫描 skill，与固定的内置加载器合并，得到本轮快照。
- skill 摘要段落与 memory 段落移入 `<system-reminder>`。
- `load_skill` 在返回内容中追加 base 目录行；当来源为 `builtin` 时先落盘整棵目录树到 `~/.cynosure/system/skills/<name>/`，并以该落盘目录作为 base。
- `/skills` 渲染改为显示「来源类型」中文标签，移除 `path`。

为什么选「每轮重载」而不是「文件监听 / mtime 缓存」：项目既有惯例是每轮重建系统提示词（git status、cynosureMd、memory 都是每轮注入），skill 重载与之一致；skill 数量少、扫描成本低（仅 `WalkDir` 找 `SKILL.md` + 读 frontmatter），无需引入 watcher 复杂度与缓存失效问题。这与用户「每轮对话都实时动态加载一次」的要求精确对齐。

### 3.2 内置 skill 区分：来源即真相

「系统内置 skill」的判定**完全依据 `SkillEntry.Source == "builtin"`**。内置 skill 来自嵌入二进制（`assets.BuiltinSkillsFS()`），其 `Path` 是 FS 内的相对路径（如 `skill-creator/SKILL.md`），不是真实磁盘路径。因此：

- 内置加载器在启动时构建一次即可（嵌入内容不会变），保存在 `Service` 上，无需每轮重载。
- 用户级 / 工作区级 skill 每轮从磁盘重载。
- 合并仍遵循 workspace > user > builtin 的同名覆盖；若某内置 skill 被同名用户 / 工作区 skill 覆盖，则该名字的来源变为 user/workspace，落盘逻辑不触发（落盘只针对最终来源为 builtin 的 skill）。

### 3.3 运行期数据结构调整

**`Service`（`internal/agent/runtime/runtime.go`）新增字段：**

```go
type Service struct {
    // ... 现有字段 ...
    BuiltinSkills *sessions.SkillLoader // 嵌入二进制内置 skill，启动时固定（语义不变）
    UserSkillsDir string                // ~/.cynosure/skills（每轮重载来源）
}
```

> 说明：`BuiltinSkills` 字段名保持不变，但**语义收敛为「仅内置」**——不再把 user/workspace 合并进来。原先 `bootstrap.go:loadSkills` 的三方合并拆分为：内置加载器单独保存；用户目录路径保存为字符串供每轮重载；工作区目录由 `Cfg.WorkspaceRoot` 推导（`config.WorkspaceCynosureSkillsDir`），无需额外字段。因此实际只新增 `UserSkillsDir` 一个字段（外加一个 setter）。

**新增 setter：**

```go
func (s *Service) SetUserSkillsDir(dir string) { s.UserSkillsDir = dir }
```

`SetBuiltinSkills` 保持不变，但调用方传入的将是**仅内置**的加载器。

### 3.4 每轮重载：`buildSkillSnapshot`

`internal/agent/runtime/prompt_builder.go` 改为：

```go
func (s *Service) buildSkillSnapshot(ctx context.Context, userID string) (*agenttools.SkillSnapshot, error) {
    // 每轮从磁盘重载 user/workspace skill；内置 skill 用启动时固定的加载器。
    userAndWorkspace, err := sessions.LoadSkillsFromDirs([]sessions.SkillDir{
        {Path: s.UserSkillsDir, Source: "user"},
        {Path: config.WorkspaceCynosureSkillsDir(s.Cfg.WorkspaceRoot), Source: "workspace"},
    })
    if err != nil {
        return nil, err
    }
    merged := sessions.MergeSkillLoaders(s.BuiltinSkills, userAndWorkspace)
    return agenttools.NewSkillSnapshot(nil, merged), nil
}
```

- `LoadSkillsFromDirs` 已存在，对不存在的目录静默返回空（`LoadAllFromDirWithSource` 内 `os.IsNotExist` 已处理）。
- 合并顺序保证 workspace > user > builtin 覆盖语义不变。
- 出错（如目录读权限异常）向上返回；`conversation_flow.go:95-98` 已有错误处理路径。

> 子 agent 复用父快照（`subagent.go:66 buildSubagentProfile(kind, parent.User, parent.Skills)`），因此父快照每轮重载后，子 agent 自动获得最新 skill，无需额外改动。

### 3.5 skill 段落与 memory 段落移入 `<system-reminder>`

`internal/assistant/prompt.go` 调整：

- **skills**：`BuildSystemPrompt` 当前在 `<tools>` 之后单独 `renderTag("skills", ...)`。改为把 skills 摘要体并入 `<system-reminder>`，删除顶层独立的 `<skills>` 段落渲染。
- **memory**：`BuildSystemPrompt` 当前在末尾单独 `renderTag("memory", memorySectionGuidance + memory)`。改为把「记忆守卫文案 + 记忆正文」并入 `<system-reminder>`，删除顶层独立的 `<memory>` 段落渲染。
- 具体做法：`renderSystemReminder` 增加 `skillDescriptions string` 与 `memorySection string` 两个入参；在 `current_day` 之后依次追加 skills 区块、cynosureMd 区块、memory 区块。
- `BuildSystemPrompt` 中原 `<git-status>` 段落保持顶层不变（git status 不属于本次调整范围）。

**调整后 system prompt 顶层段落顺序：** `identity → workspace → tools → git-status`（原末尾的 `<memory>` 与 `<tools>` 后的 `<skills>` 顶层段落均移除，内容并入 `<system-reminder>`；`<system-reminder>` 作为紧跟 system message 的临时 user message 注入，不再属于 system prompt 字符串）。

**调整后 `<system-reminder>` 内部结构（自上而下）：**

```
<system-reminder>
current_day: 2026-06-25

<skills 区块>
以下技能只提供摘要。
重要规则：
- ...
- 使用或遵循某个技能前，先用 `load_skill` 以精确的技能名加载其完整说明。
- ...
可用技能：
<skills>...</skills>

在回答用户问题时，你可以参考以下上下文：   ← cynosureMd（如有）
# cynosureMd
...

<memory 区块>                                ← 记忆（如有）
记忆（Memory）段可能同时包含过往记忆索引和真实有效记忆...（memorySectionGuidance）
...记忆正文...
</system-reminder>
```

> 设计取舍：`current_day`、skills、cynosureMd、memory 都属于「运行期动态注入的提醒 / 上下文类信息」，并入同一 `<system-reminder>` 符合需求「skill 段落、memory 段落都放到 `<system-reminder>` 段落中」。段落内顺序定为 current_day → skills → cynosureMd → memory：记忆守卫文案明确「记忆不具事实优先级、须服从当前上下文」，放在最后可降低其对前序运行期事实信息的干扰。skills 区块、cynosureMd 区块的相对顺序与原实现保持一致。

**空值处理：** skills 摘要为空时不渲染 skills 区块；memory 为空时不渲染 memory 区块（含守卫文案）——沿用现有「memory 段为空则不注入 `memorySectionGuidance`」语义；当 current_day、skills、cynosureMd、memory 四者全空时，`<system-reminder>` 段落整体省略。

> 影响的既有测试：
> - `internal/assistant/prompt_test.go:TestBuildSystemPromptUsesLoadedBasePromptAndAppendsDynamicSections` 断言了顶层 `<skills>` / `<memory>` 段落，以及 `<tools>` → `<git-status>` → `<memory>` 的顺序。这些断言将按新结构更新（skills、memory 均位于临时 user `<system-reminder>` 内；system prompt 顶层不再有 `<skills>`、`<memory>` 或 `<system-reminder>`）。
> - `TestBuildSystemPromptOmitsMemoryGuidanceWhenMemorySectionIsEmpty`、`TestBuildSystemPromptOmitsEmptyCynosureMarkdownContext` 等空值断言需复核：memory 为空时不应出现 `memorySectionGuidance`，但 `<system-reminder>` 可能因 current_day 等仍存在。
> - `internal/agent/runtime/runtime_test.go` 中校验 `<system-reminder>`、记忆相关的断言同步复核。
> - 一律**以新实现为准修正测试断言**，符合「以代码为事实修复单测」的项目惯例。

### 3.6 `load_skill` 返回 base 路径 + 内置 skill 落盘

#### 3.6.1 base 路径来源

为每个 skill 计算 base 目录（skill 根目录，即 `SKILL.md` 所在目录）：

- **user / workspace skill**：`SkillEntry.Path` 是真实磁盘路径（如 `/Users/<you>/.cynosure/skills/skill-creator/SKILL.md`），base = `filepath.Dir(Path)`，**原样展示完整路径，不缩写家目录**。
- **builtin skill**：`SkillEntry.Path` 是嵌入 FS 内相对路径（如 `skill-creator/SKILL.md`），无真实磁盘路径。落盘后 base = `~/.cynosure/system/skills/<name>` 的完整绝对路径（同样原样展示，不缩写）。

#### 3.6.2 内置 skill 落盘

新增落盘逻辑（建议放在 `internal/sessions` 或新文件 `internal/sessions/skill_materialize.go`，因为它需要访问嵌入 FS 子树并写盘）：

```go
// MaterializeBuiltinSkill 把内置 skill 的整棵目录树从嵌入 FS 落盘到 destRoot/<name>/。
// 若目标目录已存在则跳过（不覆盖用户可能的本地改动）。返回落盘后的 base 目录。
func MaterializeBuiltinSkill(fsys fs.FS, name, destRoot string) (string, error)
```

- 入参：嵌入 FS 子树（`assets.BuiltinSkillsFS()`）、skill 名（= FS 内顶层目录名 = `canonicalSkillName`）、目标根（`~/.cynosure/system/skills`）。
- **已存在则跳过**：`os.Stat(dest)` 成功即直接返回 `dest`，不写入（符合用户选择「已存在则跳过」）。
- 否则递归遍历 `fsys` 下该 skill 子树（`fs.WalkDir(fsys, name, ...)`），逐文件 `MkdirAll` + `WriteFile` 到 `destRoot/<相对路径>`，保留目录结构（scripts/references/assets/agents 等）。
- 写入采用「先写临时目录再原子 rename」或「逐文件写入」二选一；考虑到 skill 子树可能含多文件，采用**逐文件写入**并在失败时清理半成品目录，避免遗留不完整状态。
- 返回 base 目录 `destRoot/<name>`（真实磁盘路径）。

**FS 内 skill 子目录名定位：** 内置 skill 的 `SkillEntry.Path` 形如 `skill-creator/SKILL.md`，其顶层目录名即落盘子目录名。但 `canonicalSkillName` 可能取自 frontmatter `name`（与目录名不一定相同）。为稳妥，落盘按 **FS 内 `SKILL.md` 的父目录**（`filepath.Dir(entry.Path)` 在 FS 语义下，即 `path.Dir`）作为源子树根与目标子目录名，保证目录树完整对应；返回的 base 用该目录名。若 frontmatter `name` 与目录名不同，以 FS 目录名作为落盘目录（与磁盘结构一致），日志记录二者差异以便排查。

> 关于 `~/.cynosure/system/skills` 的新路径助手：在 `internal/config/local_config.go` 新增 `CynosureSystemSkillsDir() (string, error)`，返回 `~/.cynosure/system/skills`，与既有 `CynosureSkillsDir` 风格一致。

#### 3.6.3 `load_skill` 流程改造

`internal/tools/load_skill.go`：

- `LoadedSkill` 结构新增 `BaseDir string` 字段。
- `SkillSnapshot.LoadSkill(name)` 命中条目后：
  - 若 `source == "builtin"`：调用落盘逻辑，得到 `~/.cynosure/system/skills/<name>`，作为 `BaseDir`。
  - 否则：`BaseDir = filepath.Dir(entry.Path)`（真实磁盘路径）。
- `renderLoadedSkill` 在 `<content>` 之后、`</skill>` 之前追加一行 base 目录说明：

```
<skill source="builtin" name="skill-creator">
<metadata>
...
</metadata>
<content>
...skill 正文...
</content>
Base directory for this skill: /Users/<you>/.cynosure/system/skills/skill-creator
</skill>
```

- 展示路径时原样输出原始完整路径，不缩写家目录（即不把 `$HOME` 前缀替换为 `~`）。

**落盘依赖注入：** `SkillSnapshot` 当前不持有嵌入 FS 与目标根路径。方案：在 `NewSkillSnapshot` 或快照构造处注入一个「内置 skill 物化器」函数 / 接口，避免 `tools` 包直接依赖 `assets` 与 `config`（保持依赖方向干净）。具体地：

- 在 `SkillSnapshot` 增加可选字段 `BuiltinMaterializer func(name string) (baseDir string, err error)`。
- `runtime` 层（`buildSkillSnapshot`）构造快照时注入该闭包：闭包内部调用 `sessions.MaterializeBuiltinSkill(assets.BuiltinSkillsFS(), name, systemSkillsDir)`。
- `tools.LoadSkill` 命中 builtin 且 `BuiltinMaterializer != nil` 时调用它得到 BaseDir；为空（如单测未注入）时退化为不落盘，BaseDir 用 FS 相对路径或留空（不影响正文返回）。

> 这样 `tools` 包零新增对 `assets`/`config` 的依赖，落盘细节封装在 `sessions` + `runtime` 注入。

#### 3.6.4 错误处理

- 落盘失败（磁盘权限、空间不足等）：记录 warning 日志，BaseDir 退化为 FS 相对路径或 `~/.cynosure/system/skills/<name>`（标称路径），**不阻断** skill 正文返回——保证 `load_skill` 主功能（返回工作流说明）不受落盘副作用影响。

### 3.7 `/skills` 只显示来源（去掉 path）

`internal/tui/app.go:renderSkillDetails` 调整：

- 移除 `path:` 行。
- `[source]` 原样字符串（`builtin`/`user`/`workspace`）映射为中文来源标签：

| Source 值 | 显示标签 |
| --- | --- |
| `workspace` | `当前项目` |
| `user` | `家目录` |
| `builtin` | `系统内置` |

- 渲染格式（每行一条，符合用户「列表每行一条」偏好）：

```
已加载 Skills：3 个
- skill-creator [系统内置] Create and scaffold skills...
- project-helper [当前项目] Project-specific helper...
- my-tool [家目录] Personal tooling...
```

> `SkillSummary` 仍可保留 `Path` 字段（数据结构不动），仅渲染层不展示。`SessionInfo.Skills` 注入的 summaries 来自 `bundle.Skills`（启动快照）；由于 `/skills` 展示的是「已加载」概览且需求只要求显示来源，保持注入启动快照即可（来源类型不随每轮重载变化）。

> 影响的既有测试：`internal/tui/events_test.go:TestHandleSkillsCommandShowsSkillDetails` 断言输出含 `path`，将更新为断言含中文来源标签且**不含** path。

### 3.8 bootstrap 装配调整

`internal/local/bootstrap.go`：

- `loadSkills` 拆分语义：
  - 内置加载器：`sessions.LoadSkillsFromFS(assets.BuiltinSkillsFS(), "builtin")`，传给 `runtimeService.SetBuiltinSkills(builtinOnly)`。
  - 用户目录：`runtimeService.SetUserSkillsDir(userSkillsDir)`。
  - 工作区目录由 `Cfg.WorkspaceRoot` 推导，无需额外字段。
- `Bundle.Skills`（用于 `/skills` 展示的启动快照）：仍由「内置 + user + workspace」合并出 summaries（启动时算一次即可，来源类型固定）。即 bootstrap 仍调用一次三方合并用于 `Summaries()`，但 runtime 每轮重载用于系统提示词。
- `Bundle.SkillCount` 含义不变。

> 取舍：`/skills` 的 Bundle.Skills 与系统提示词的每轮快照是两条路径。Bundle.Skills 反映启动时全集（含来源），满足展示需求；系统提示词每轮重载反映最新磁盘状态，满足动态加载需求。二者不冲突。

## 4. 涉及文件清单

| 文件 | 改动 |
| --- | --- |
| `internal/agent/runtime/runtime.go` | 新增 `UserSkillsDir` 字段与 `SetUserSkillsDir`；`BuiltinSkills` 语义收敛为仅内置 |
| `internal/agent/runtime/prompt_builder.go` | `buildSkillSnapshot` 改为每轮重载 user/workspace + 合并内置；注入 BuiltinMaterializer |
| `internal/assistant/prompt.go` | skills 摘要与 memory 段落移入 `<system-reminder>`；删除顶层 `<skills>` 与 `<memory>` 段落 |
| `internal/tools/load_skill.go` | `LoadedSkill.BaseDir`；`SkillSnapshot.BuiltinMaterializer`；`renderLoadedSkill` 追加 base 行（原始完整路径，不缩写家目录） |
| `internal/sessions/skill_materialize.go`（新增） | `MaterializeBuiltinSkill` 落盘整棵目录树（已存在则跳过） |
| `internal/config/local_config.go` | 新增 `CynosureSystemSkillsDir()` → `~/.cynosure/system/skills` |
| `internal/local/bootstrap.go` | 拆分内置 / 用户 / 工作区装配；注入 UserSkillsDir |
| `internal/tui/app.go` | `renderSkillDetails` 去掉 path，来源转中文标签 |
| `assets/system_prompt.md` & `DefaultBaseSystemPrompt`（`prompt.go`） | 同步「skills 摘要、记忆均位于 `<system-reminder>`」与「load_skill 返回含 base 目录」的描述（提示词镜像，两处都改） |
| `README.md` | 同步 skill 机制说明（动态加载、系统内置落盘路径、/skills 来源展示）；如有描述 memory 段落位置处一并同步 |

## 5. 测试计划

**单元测试（新增 / 更新，以代码为事实修正旧断言）：**

1. `prompt_test.go`：更新顺序断言——skills 摘要与 memory 正文均出现在临时 user `<system-reminder>` 内、不再是顶层 `<skills>` / `<memory>` 段落；system prompt 顶层段落顺序为 `identity → workspace → tools → git-status`；空 skills/memory 时 system-reminder 仍可仅含 current_day；memory 为空时不出现 `memorySectionGuidance`。
2. `runtime_test.go`：
   - `buildSkillSnapshot` 每轮重载：在 user/workspace 目录写入 skill 文件后调用，验证快照含该 skill；删除文件后再次调用，验证快照不再含它（动态性）。
   - builtin 来源 skill 经 `load_skill` 返回内容含 `Base directory for this skill: ~/.cynosure/system/skills/<name>`，且 `~/.cynosure/system/skills/<name>/SKILL.md` 及子目录文件已落盘。
   - 落盘「已存在则跳过」：预置目标目录含哨兵文件，调用后哨兵文件保留（未被覆盖）。
3. `load_skill` 渲染测试：user/workspace skill 返回 base = 其真实目录（原始完整路径，不缩写家目录）；builtin 注入 materializer 返回系统落盘目录的完整路径。
4. `events_test.go`：`/skills` 输出含 `当前项目`/`家目录`/`系统内置` 标签，**不含** 磁盘 path。
5. `sessions` 新增 `MaterializeBuiltinSkill` 单测：整棵目录树落盘（含嵌套子目录）；已存在跳过；目标根不可写时返回错误。

**集成验证：**

- `go build ./...`、`go vet ./...`、`go test ./... -count=1`（关注 `internal/assistant`、`internal/agent/runtime`、`internal/tools`、`internal/sessions`、`internal/tui`、`internal/config`）。
- 手动：真实运行 TUI，`load_skill skill-creator` 观察返回含落盘 base 路径且 `~/.cynosure/system/skills/skill-creator/` 出现完整目录树；`/skills` 只显示来源标签；对话中新增 `<ws>/.cynosure/skills/foo/SKILL.md` 后下一轮 `<system-reminder>` 内出现该 skill 摘要。

## 6. 兼容性与风险

- **来源字符串契约**：落盘与中文标签都依赖 `source` 取值 `builtin`/`user`/`workspace`。这是 bootstrap 设定的内部契约，单测覆盖，风险低。
- **每轮磁盘扫描开销**：skill 数量通常个位数，`WalkDir` + 读 frontmatter 成本可忽略；与每轮重建系统提示词的既有开销同量级。
- **落盘副作用**：仅在 `load_skill` 命中 builtin 且目标不存在时写盘；「已存在跳过」避免覆盖用户改动；失败不阻断主功能。
- **提示词镜像**：`assets/system_prompt.md` 与 `prompt.go:DefaultBaseSystemPrompt` 须同步描述，避免漂移（项目既有约束）。
- **memory 段落仅改渲染位置**：`<memory>` 内容（守卫文案 + 记忆正文）从顶层末尾移入 `<system-reminder>`，记忆子系统的提取 / 选择 / 断点 / 注入逻辑零改动；需复核依赖「顶层 `<memory>` 段落」位置的既有断言。
- **不影响子 agent / 记忆子系统 / 压缩 / 审批**：改动集中在 skill 加载、提示词 skills+memory 段落渲染位置、load_skill 渲染、/skills 渲染几处，记忆生成与其余运行期机制零改动。

## 7. 实施顺序（实现阶段参考）

1. `config.CynosureSystemSkillsDir` + `sessions.MaterializeBuiltinSkill`（+ 单测）。
2. `runtime.Service` 字段与 `buildSkillSnapshot` 每轮重载（+ 单测）。
3. `prompt.go` skills 段落与 memory 段落移入 `<system-reminder>`（+ 更新 prompt 单测）。
4. `load_skill.go` base 路径 + materializer 注入（+ 单测）。
5. `tui/app.go` `/skills` 来源标签（+ 更新 events 单测）。
6. `bootstrap.go` 装配拆分。
7. 同步 `system_prompt.md` / `DefaultBaseSystemPrompt` / `README.md`。
8. 全量 `go build/vet/test`，手动 TUI 验证。
