# load_skill 工具设计文档

## 背景

当前 Web Runtime 已经存在 `load_skill` 工具，并且默认会在未显式配置工具白名单时开启：

- 工具定义位于 `internal/tools/definitions.go`
- 默认工具名由 `internal/web/runtime/tool_registry.go` 中的 `defaultWebAllowedTool = "load_skill"` 控制
- 当前会话开始前会通过 `buildSkillSnapshot` 从数据库读取当前用户 enabled skill，并和本地 builtin skill 合并
- 当前 `load_skill` 执行逻辑通过合并后的 `SkillLoader.GetContent(name)` 返回 skill body

本次目标是在现有基础上明确并强化 `load_skill` 的语义：它是一个默认开启的工具，用来加载指定 skill 的全部信息；查找顺序必须是先查当前用户数据库 skill，查不到再查本地 skill。

## 目标

1. `load_skill` 默认开启。
2. 工具描述明确为：加载 skill 的全部信息。
3. 执行 `load_skill` 时按优先级查找：
   1. 当前用户数据库中的 enabled skills
   2. 本地 builtin skills
4. 如果数据库和本地存在同名 / 同 slug skill，数据库版本优先。
5. 返回内容包含完整 skill 信息，而不只是 body：来源、名称、描述、路径或 ID、正文内容等。
6. 保持单轮对话内 skill snapshot 一致：每次用户发送消息前刷新数据库 skill，本轮内 `system prompt` 和 `load_skill` 使用同一份 snapshot。

## 非目标

1. 不在工具调用时反复访问数据库；数据库刷新仍发生在每次用户消息进入 runtime 时。
2. 不新增前端 API。
3. 不改变 skill 创建、编辑、删除接口。
4. 不引入本地 skill 热重载；本地 builtin skill 仍按服务启动时加载的结果为准。

## 当前问题

当前实现已经通过 `MergeSkillLoaders(s.BuiltinSkills, buildDBSkillLoader(skills))` 达到了“DB 覆盖本地同名 skill”的效果，但存在三个不足：

1. `ToolContext` 只持有合并后的 `Loader`，执行时无法显式知道 skill 来自数据库还是本地。
2. `GetContent` 只返回 `<skill>` 包裹的 body，没有返回完整元信息。
3. 工具描述偏泛化：`Load a skill by name to access specialized knowledge`，没有表达“加载全部信息”和“DB 优先，本地兜底”。

## 方案对比

### 方案 A：继续只使用合并后的 SkillLoader

做法：保留当前合并逻辑，只调整工具描述和 `GetContent` 输出。

优点：改动最小。

缺点：无法明确返回 source；也不容易在测试中证明执行路径是“DB 先、本地后”，只能通过覆盖结果间接证明。

### 方案 B：引入显式 SkillSnapshot / SkillResolver（推荐）

做法：会话开始时构建一个结构化 snapshot，分别保存用户 DB loader、本地 loader 和合并 loader。`system prompt` 使用合并 loader 渲染描述；`load_skill` 使用 resolver 先查 DB loader，失败后查本地 loader。

优点：语义清晰；测试容易覆盖；返回结果可以明确标记 `source=db` 或 `source=local`；符合“先数据库，后本地”的要求。

缺点：需要小幅调整 runtime 内部类型。

### 方案 C：工具执行时实时查数据库，再查本地

做法：每次 `load_skill` 调用都根据 userID 查询数据库，查不到再查本地。

优点：工具调用时获取的是最新 DB 状态。

缺点：破坏单轮响应内 skill snapshot 一致性；增加 DB 查询次数；和当前“每次用户消息前刷新一次 skill snapshot”的设计冲突。

## 推荐方案

采用方案 B：引入显式 `SkillSnapshot` / `SkillResolver`。

原因：它既满足“每次用户消息前从数据库刷新”的既有设计，也能在工具执行时严格表达“先 DB，后本地”的查找顺序，并且可以返回完整元信息。

## 详细设计

### 1. SkillSnapshot

新增 runtime 内部结构：

```go
type SkillSnapshot struct {
    UserSkills  *sessions.SkillLoader
    LocalSkills *sessions.SkillLoader
    Merged      *sessions.SkillLoader
}
```

职责：

- `UserSkills`：由 `ListEnabledSkillsByUser` 查询数据库后构建。
- `LocalSkills`：服务启动时加载的 builtin skill loader。
- `Merged`：用于 system prompt 展示可用 skill 描述，DB 同名 skill 覆盖本地 skill。

### 2. 构建流程

替换当前 `buildSkillSnapshot(ctx, userID) (*sessions.SkillLoader, error)` 为：

```go
func (s *Service) buildSkillSnapshot(ctx context.Context, userID string) (*SkillSnapshot, error)
```

构建逻辑：

1. 调用 `Store.ListEnabledSkillsByUser(ctx, userID)` 获取用户 DB skill。
2. `buildDBSkillLoader(skills)` 得到 `UserSkills`。
3. `LocalSkills = s.BuiltinSkills`。
4. `Merged = MergeSkillLoaders(LocalSkills, UserSkills)`，确保 DB 覆盖本地同名 skill。

### 3. system prompt 使用

`buildSystemPrompt` 使用 `snapshot.Merged` 渲染 skill 描述：

```go
SkillDescriptions: snapshot.Merged.GetDescriptions()
```

这样用户看到的可用 skill 列表仍然是合并后的结果，且同名时展示 DB 版本。

### 4. load_skill 工具上下文

调整 `ToolContext`：

```go
type ToolContext struct {
    User         storage.User
    Conversation storage.Conversation
    Skills       *SkillSnapshot
}
```

保留兼容策略：如果为了减少一次性改动，也可以先保留 `Loader` 字段，并新增 `Skills` 字段；执行 `load_skill` 时优先使用 `Skills`，没有时退回 `Loader`。

### 5. load_skill 查找规则

新增 resolver 方法：

```go
func (s *SkillSnapshot) LoadSkill(name string) (LoadedSkill, error)
```

查找顺序：

1. `UserSkills.GetEntry(name)` 成功则返回 `source=db`。
2. `LocalSkills.GetEntry(name)` 成功则返回 `source=local`。
3. 都失败则返回包含可用 skill 列表的错误。

### 6. 返回格式

建议 `load_skill` 返回结构化文本，便于模型读取，也便于调试：

```xml
<skill source="db" name="skill-name">
<metadata>
description: ...
path: db://skills/{id}
tags: ...
</metadata>
<content>
...
</content>
</skill>
```

本地 skill 示例：

```xml
<skill source="local" name="skill-name">
<metadata>
description: ...
path: /absolute/path/to/SKILL.md
tags: ...
</metadata>
<content>
...
</content>
</skill>
```

说明：

- `source` 用于明确命中来源。
- `metadata` 输出 `SkillEntry.Meta` 中的全部字段，保证“全部信息”。
- `content` 输出完整正文。
- 继续追加 runtime paths note，保持当前 `APP_HOME / COMMAND_BIN_DIR / WORKSPACE_ROOT` 等上下文行为。

### 7. 工具定义

更新 `internal/tools/definitions.go` 中 `load_skill` 描述：

```text
Load the full information of a skill by name. It first looks up the current user's enabled database skills, then falls back to local builtin skills.
```

参数保持不变：

```json
{
  "name": "skill-name"
}
```

### 8. 默认开启

保持现有默认逻辑：

- `defaultWebAllowedTool = "load_skill"`
- `LoadWebConfig` 默认 `WEB_ALLOWED_TOOLS` 包含 `load_skill`

如果用户没有配置任何工具，Web Runtime 只暴露 `load_skill`。

## 测试设计

### 单元测试

1. `load_skill` 默认开启
   - 不配置 `WebAllowedTools`
   - 断言 `RegisteredTools(cfg)` 只包含 `load_skill`

2. DB skill 优先
   - 本地和 DB 都有同名 skill
   - 调用 `load_skill`
   - 断言返回 `source="db"`，正文为 DB 内容

3. DB 查不到时 fallback 本地
   - DB 无该 skill，本地存在
   - 调用 `load_skill`
   - 断言返回 `source="local"`，正文为本地内容

4. 返回完整元信息
   - 构造包含 description、tags 等 meta 的 skill
   - 断言输出包含 metadata 字段和完整 content

5. 未知 skill 报错
   - DB 和本地都不存在
   - 断言错误包含 unknown skill 和可用 skill 列表

### 集成测试

1. 会话开始时刷新 DB skill snapshot
   - 第一次返回 DB skill A
   - 第二次更新 store 返回 DB skill B
   - 断言两次会话 loop 使用不同 snapshot

2. system prompt 和 load_skill 使用同一 snapshot
   - 当前轮构建 snapshot 后，prompt 描述和工具加载结果一致

## 兼容性

1. 外部 API 不变。
2. 工具名和参数不变。
3. 默认开启行为不变。
4. 只增强返回信息和查找语义。
5. 当前已有用户自定义 skill 每次消息前刷新的行为保持不变。

## 实施步骤

1. 新增 `SkillSnapshot` 与 resolver。
2. 调整 `buildSkillSnapshot` 返回结构化 snapshot。
3. 调整 `buildSystemPrompt` 使用 `snapshot.Merged`。
4. 调整 `ToolContext` 和 `load_skill` 执行逻辑。
5. 更新 `load_skill` 工具描述。
6. 补充 DB 优先、本地 fallback、完整信息输出、默认开启测试。
7. 运行 `go test ./...`。

## 风险与规避

1. 风险：调整 `ToolContext` 影响现有工具执行测试。
   - 规避：保留兼容字段或集中修改 runtime 测试 helper。

2. 风险：返回完整 metadata 后输出变长。
   - 规避：metadata 只输出 `SkillEntry.Meta`、source、path，不额外读取文件系统内容。

3. 风险：DB 与本地同名 skill 的 prompt 描述和工具加载不一致。
   - 规避：prompt 使用 `Merged`，工具使用 DB-first resolver，二者优先级一致。

## 结论

推荐实现显式 `SkillSnapshot` / `SkillResolver`。该方案能清晰满足“默认开启、加载全部信息、先数据库再本地”的需求，同时保持当前每次用户消息前刷新 skill 的一致性模型。
