# 记忆机制调整 设计文档

日期：2026-06-18
范围：`cynosure/internal/agent/runtime`、`cynosure/internal/local`、`cynosure/internal/assistant`、`cynosure/internal/tools`、`cynosure/assets/prompts`

## 1. 背景与目标

当前记忆系统的工作方式：

- **注入**：每轮对话前调用 LLM（`selectRelevantMemories` + `memory_selection.md`），从全部记忆中挑选最多 10 条，渲染进系统提示词 `<memory>` 段。`memory.md` 本身**不注入**，仅作为 `- [name](path) — desc` 索引文件。
- **提取**：每轮结束后异步调用 LLM（`extractMemories` + `memory_extraction.md`）抽取 4 类记忆（preference/feedback/project/reference），逐条 `InsertMemory` 落为 `<name>.md` 文件并重建 `memory.md`。
- **去重/淘汰**：按类型阈值即时触发（`maybeConsolidateType`，preference 20 / feedback 30 / project 30 / reference 30），超阈值时整类喂给 LLM 合并并 `ReplaceMemoriesByUserAndType`。

本次调整目标（对应需求 1–6）：

1. 把 `memory.md` 索引注入系统提示词，最多 200 行、单行（条目）不超过 25KB；超限时仅注入部分并附截断警告。
2. 系统提示词新增"更新/删除过期记忆"的指导，并新增 `update_memory`、`delete_memory` 两个工具；二者必须**同时**操作记忆文件和 `memory.md`。
3. 每轮结束的记忆提取，把已存在记忆注入提示词以避免重复提取（已有 `existing` 注入，需补充 `## Existing memory files` 段式提示）。
4. 选择（注入）记忆不再读 `memory.md` 内容做候选，改为**确定性扫描**记忆目录得到候选集：扫描所有 `.md`（排除 `memory.md`），每个文件只读前 30 行 frontmatter，按 mtime 降序保留最新 200 个；再交给 `selectRelevantMemories`（LLM）做精选——**最多返回 5 个、只选确定有帮助的、不确定不选（宁缺毋滥）、无相关返回空列表**；被选中的记忆**完整文件内容**注入主流程系统提示词，标注**相对时间**（如 `47 days ago`），并对超过 1 天的记忆追加"时点观测、可能过期、需核对当前代码"的说明。
5. 同一会话中不重复注入同一条记忆；若该记忆已在本会话注入过，则**重读其文件**并替换旧内容（捕获文件更新）。
6. 去重/淘汰改为**定期运行**（默认间隔 24h、累积 5+ 次会话触发），运行时同时更新记忆文件与 `memory.md`。

**注意事项**：其余策略不变，不影响现有功能，废弃无用代码。

## 2. 关键决策（已与用户确认）

- **D1（注入结构）**：需求 1 与需求 4 **都注入**。`<memory>` 段包含两个子块：① `memory.md` 索引块（受 200 行 / 25KB / 截断警告约束）；② 确定性扫描得到的记忆条目块（带相对时间）。
- **D2（合并策略）**：原"按类型阈值即时合并"**完全替换**为需求 6 的"定时全量去重"（24h + 累积 5 次会话）。
- **D3（工具寻址）**：`update_memory` / `delete_memory` 用**文件名/相对路径**定位（与 `memory.md` 索引中的 `path` 一致）。
- **D4（会话去重状态）**：放在 **Service 内存态**，按 `conversationID` 索引，记录本会话已注入记忆（文件名 → mtime/内容）。进程重启/恢复会话后重建。
- **D5（定时去重状态）**：放在**记忆目录的元数据文件**（如 `.consolidation_state.json`），记录上次去重时间与累计会话计数，跨进程持久。
- **D6（LLM 选择保留并改造）**：保留 `selectRelevantMemories` + `memory_selection.md`，但**改变其输入与输出**：候选集从"读 `memory.md` / 全量记忆"改为"确定性扫描得到的最多 200 个候选"；选择规则改为**最多 5 个、宁缺毋滥、无相关返回空**；被选中的记忆**完整文件内容**注入主流程系统提示词，并标注相对时间与过期说明。去重环节另有独立 LLM 调用（`consolidateViaLLM`）。

## 3. 总体设计

### 3.1 注入侧（系统提示词 `<memory>` 段）

改造 `selectRelevantMemories`：候选集由确定性扫描产生（替换原"读 `memory.md` / 全量记忆"），LLM 精选后注入被选中记忆的**完整文件内容**。

**数据流**：

```
buildSystemPrompt(ctx, conv, user, snapshot, history, memoryOn)
  └─ if memoryOn:
       memorySection = s.buildMemorySection(ctx, conv.ID, user, history)
            ├─ indexBlock = store.LoadMemoryIndexForPrompt()        // 需求1：memory.md 索引块
            ├─ scanned    = store.ScanRecentMemories()              // 需求4-a：≤200 候选，前30行 frontmatter，mtime 降序
            ├─ picked     = selectRelevantMemories(scanned, history) // 需求4-b：LLM 精选，最多 5，宁缺毋滥，无则空
            ├─ for each picked: 读完整文件内容；按 conv.ID 会话态去重/重读替换 // 需求5
            └─ render → <memory> 段：索引块 + 选中记忆完整内容块（带相对时间 + 过期说明）
```

#### 需求 1：注入 `memory.md` 索引

- 在 `MarkdownMemoryStore` 新增 `LoadMemoryIndexForPrompt() (text string, truncated bool, totalLines int)`：
  - 读取 `memory.md` 全文。
  - 常量：`memoryIndexMaxLines = 200`、`memoryIndexMaxEntryBytes = 25 * 1024`。
  - 逐行累计，**最多保留 200 行**；任一行字节数超过 25KB 时对该行做字节安全截断（按 rune 边界）。
  - 若发生行数截断或任何行被截断，`truncated = true`，并返回真实总行数。
- 渲染时若 `truncated`，在索引块顶部加入警告（结合项目语境的英文提示）：
  ```
  WARNING: MEMORY.md is <totalLines> lines (limit: 200). Only part of it was loaded.
  Keep index entries to one line under ~200 chars; move detail into topic files.
  ```

#### 需求 4-a：确定性扫描（生成候选集）

- 在 `MarkdownMemoryStore` 新增 `ScanRecentMemories() ([]ScannedMemory, error)`：
  - 常量：`FrontmatterMaxLines = 30`、`MaxMemoryFiles = 200`。
  - `filepath.Glob(rootDir/*.md)`，排除 `memory.md`。
  - 每个文件 `os.Stat` 取 `ModTime`；只读前 30 行解析 frontmatter（复用/扩展 `splitFrontMatter`，但限制行数）。
  - 按 mtime 降序排序，保留前 200 个。
  - `ScannedMemory` 字段：`Path`、`Name`、`Description`、`Type`、`ModTime`。（扫描阶段不读 body，仅用 name+description 供 LLM 精选。）

#### 需求 4-b：LLM 精选（selectRelevantMemories 改造）

- `selectRelevantMemories` 输入改为 `ScanRecentMemories()` 的候选集（不再 `ListRelevantMemories` 全量、不再读 `memory.md`）。
- user prompt：把候选按 `[i] (type) name: description` 编号列出（复用现有 `buildSelectionUserPrompt` 风格），加上当前对话上下文。
- `memory_selection.md` 选择规则更新为：
  - 返回最多 **5** 个记忆的索引。
  - 只选**确定有帮助**的；不确定不选——**宁可漏选不可错选**。
  - 无相关记忆时返回空数组 `[]`。
  - 输出仍为 JSON 索引数组（沿用 `parseSelectedIDs` / `pickMemoriesByIndex`，cap 改为 5）。
- 常量 `maxInjectedMemories` 由 10 改为 5。
- 选中的索引映射回候选 `ScannedMemory`，进入下一步"读完整文件 + 会话态去重 + 渲染"。

#### 注入渲染（相对时间 + 过期说明）

- 对每个被选中的记忆，读取其**完整文件内容**（frontmatter 之后的 body，复用 `readMemoryFileLocked`），注入系统提示词。
- **相对时间**：新增 `humanizeRelativeTime(t, now) string`（`just now` / `N minutes ago` / `N hours ago` / `N days ago` …），每条记忆标注其文件 mtime 的相对时间。
- **过期说明**：若被选中记忆中**任一**条 mtime 超过 1 天，在记忆块顶部追加一次说明：
  ```
  Memories are point-in-time observations, not live state — claims about code
  behavior or file:line citations may be outdated. Verify against current code
  before asserting as fact.
  ```
  （为减少冗余，整块出现一次，而非每条重复。）

#### 需求 5：会话内去重 / 重读替换

- `Service` 新增字段：`injectedMemories map[string]map[string]injectedMemoryMeta`（外层 key=conversationID，内层 key=文件相对路径），`injectedMemoryMeta{ModTime time.Time}`，配 `sync.Mutex`。
- 对 LLM 选中的每个记忆文件：
  - 若该会话尚未注入过 → 注入完整内容，记录其 mtime。
  - 若已注入且文件 mtime 未变 → 跳过（不重复注入）。
  - 若已注入但 mtime 变化 → **重读文件**取最新完整内容，替换记录的 mtime，本轮按"更新后内容"注入。
- 说明：系统提示词每轮整体重建，本会话态用于决定"哪些选中条目纳入本轮 `<memory>`"以及标记是否为"更新后重读"。这样实现"不重复注入同一条；已存在则重读替换"的语义。

### 3.2 提取侧（每轮结束，需求 3）

- 现有 `extractMemories` 已注入 `existing`（`buildExtractionUserPrompt`）。补充：在 user prompt 中以需求 3 指定格式追加现有记忆**文件清单**：
  ```
  \n\n## Existing memory files\n\n<existingMemories>\n\nCheck this list before writing — update an existing file rather than creating a duplicate.
  ```
  - `<existingMemories>` 由 `ScanRecentMemories()`（或既有 `ListRelevantMemories`）渲染为"文件名 + 一句话描述"清单。
- 不改变提取的落盘逻辑（仍 `InsertMemory`）。**去掉**提取后立即触发的 `maybeConsolidateType`（移至定时去重，见 3.4）。

### 3.3 更新/删除工具（需求 2）

#### 系统提示词指导

- 在 `DefaultBaseSystemPrompt` 记忆相关说明处新增一行：
  > Update or remove memories that turn out to be wrong or outdated.
- 在 `assistant/prompt.go` 的工具/记忆指引中说明这两个工具用途（仅当工具在本次会话清单中出现时）。

#### 工具定义（`internal/tools/definitions.go`）

- `update_memory`：参数 `path`（必填，记忆文件相对路径，如 `foo.md`）、`name`、`description`、`body`（至少一项），更新指定记忆文件并刷新 `memory.md`。
- `delete_memory`：参数 `path`（必填），删除记忆文件并从 `memory.md` 移除条目。
- 加入 `AllToolSpecs`；默认是否放入 `AllowedTools` 取决于配置（与现有工具一致：仅在配置允许时启用）。

#### 工具执行（关键架构问题）

现有工具走 `agenttools.Dispatch`（`internal/tools/handlers.go`），是**无状态纯函数**，无法访问 `MarkdownMemoryStore`。记忆文件路径在 `~/.cynosure/memory/<workspace-key>/`，**不在工作区内**，且需同步 `memory.md`。

**方案**：在 `runtime` 层拦截这两个工具（类似 `spawn_subagent` / `mcp__` 的特判），由 `Service` 直接调用 store 的新方法，不进 `agenttools.Dispatch`。

- `Service.executeToolCall` 增加分支：`name == "update_memory"` / `"delete_memory"` 时调用 `s.updateMemoryTool(...)` / `s.deleteMemoryTool(...)`。
- 这两个方法调用 store 新增接口：
  - `UpdateMemoryFile(ctx, path string, fields MemoryUpdate) error`
  - `DeleteMemoryFile(ctx, path string) error`
  - 二者内部加锁、校验 `safeRelativeMemoryPath`、写文件后 `rewriteIndexLocked()`（同步 `memory.md`）。
- 在 `conversationStore` 接口与 `local.Store` 上补充对应转发方法。

#### 工具审批

- `delete_memory` 属删除类操作。沿用现有审批链路（`approveToolCall` / `HookManager`），不特殊放行。`update_memory` 视为写操作同样经审批。具体是否需要审批由现有审批策略决定，本设计不改审批语义。

### 3.4 定时去重 / 淘汰（需求 6）

- 删除 `extractMemories` 中对 `maybeConsolidateType` 的调用与 `consolidationThresholds`（D2）。
- 新增 `Service.maybeRunConsolidation(ctx, user)`，在每轮收尾（`scheduleMemoryWork` 内、提取之后）调用：
  - 读取 store 元数据 `ConsolidationState{ LastRunAt time.Time; SessionCount int }`（D5，文件 `.consolidation_state.json`，存于记忆 rootDir）。
  - 每次会话收尾 `SessionCount++`。
  - 触发条件：`SessionCount >= consolidationMinSessions(=5)` **且** `now - LastRunAt >= consolidationInterval(=24h)`。
  - 触发后对四类记忆各自 `consolidateViaLLM` → `ReplaceMemoriesByUserAndType`（已同步 `memory.md`），随后 `LastRunAt=now; SessionCount=0` 落盘。
  - 间隔/计数默认值可由 config 覆盖（新增 `MemoryConsolidationIntervalHours`、`MemoryConsolidationMinSessions`，带默认值，缺省即 24/5）。
- store 新增：`LoadConsolidationState() (ConsolidationState, error)`、`SaveConsolidationState(ConsolidationState) error`。

## 4. 废弃 / 删除清单（需求注意事项）

- `selectRelevantMemories`（memory.go）→ **保留并改造**：候选改为扫描结果、cap 5、注入完整文件内容（非删除）。
- `buildSelectionUserPrompt` → 保留（候选列表来源改为 `ScannedMemory`）；`parseSelectedIDs` 保留；`pickMemoriesByIndex` 保留（cap 5）。
- `assets/prompts/memory_selection.md` → **保留并改写**选择规则（最多 5、宁缺毋滥、无相关返回空）。
- `maybeConsolidateType` + `consolidationThresholds` + `maxPreferenceMemories` 等阈值常量 → 删除（被定时去重替代）。
- `maxInjectedMemories` 由 10 → 5。
- 相关单测（`TestExtractMemories_Consolidates*`）→ 删除（阈值即时合并被废弃）；`TestSelectRelevantMemories_*`、`TestPickMemoriesByIndex_*` → 改写为新候选/新 cap 的行为测试。

保留：`renderMemorySection`/`renderDialogueForMemory`/`extractMemories` 提取主体/`consolidateViaLLM`/会话记忆（conversation_memory.go，与本需求无关，不动）。

## 5. 组件与接口变更汇总

### `internal/local/memory_store.go`
- `LoadMemoryIndexForPrompt() (string, bool, int)`
- `ScanRecentMemories() ([]ScannedMemory, error)` + `ScannedMemory` 类型
- `ReadMemoryFile(ctx, path) (storage.Memory, error)`（注入选中记忆完整内容用，导出 `readMemoryFileLocked`）
- `UpdateMemoryFile(ctx, path, MemoryUpdate)` / `DeleteMemoryFile(ctx, path)`
- `LoadConsolidationState()` / `SaveConsolidationState()` + `ConsolidationState` 类型
- 常量：`FrontmatterMaxLines=30`、`MaxMemoryFiles=200`、`memoryIndexMaxLines=200`、`memoryIndexMaxEntryBytes=25*1024`

### `internal/local/store.go`
- 转发上述新增方法（`s.memory != nil` 时委派）。

### `internal/agent/runtime/runtime.go`
- `conversationStore` 接口补充新增方法。
- `Service` 新增 `injectedMemories` 字段 + mutex。

### `internal/agent/runtime/memory.go`
- 新增 `buildMemorySection(ctx, conversationID, user, history)`、`humanizeRelativeTime`、选中记忆完整内容块渲染（相对时间 + 过期说明 + 会话态去重/重读）。
- 改造 `selectRelevantMemories`：候选改为 `ScanRecentMemories()`，cap 5。
- 删除阈值合并代码（`maybeConsolidateType`/`consolidationThresholds`/阈值常量）；`maxInjectedMemories` 10→5。
- 提取 prompt 增加 `## Existing memory files` 段。

### `internal/agent/runtime/prompt_builder.go`
- `buildSystemPrompt` 改为调用 `buildMemorySection`（需要 `conversation.ID`，签名已含 conversation）。

### `internal/agent/runtime/conversation_memory.go`（scheduleMemoryWork）
- 提取后调用 `maybeRunConsolidation`。

### `internal/agent/runtime/tool_registry.go`
- `executeToolCall` 增加 `update_memory` / `delete_memory` 分支。

### `internal/tools/definitions.go`
- 新增两个工具 spec（仅定义，执行在 runtime 层拦截）。

### `assets/prompts/memory_selection.md`
- 改写选择规则（最多 5、宁缺毋滥、无相关返回空）。

### `internal/assistant/prompt.go`
- `DefaultBaseSystemPrompt` 增加"更新/删除过期记忆"指导。

### `internal/config`
- 可选新增去重间隔/会话阈值配置项（带默认）。

## 6. 测试计划

- `LoadMemoryIndexForPrompt`：行数超 200 截断 + 警告；单条 >25KB 截断；正常不截断。
- `ScanRecentMemories`：排除 `memory.md`；只读前 30 行；mtime 降序；超过 200 截断。
- `selectRelevantMemories`：候选来自扫描结果；cap 5；无相关返回空。
- `humanizeRelativeTime`：各档位边界。
- `buildMemorySection`：包含索引块 + 选中记忆完整内容块；>1 天追加过期说明；会话内重复不再注入；mtime 变化触发重读替换。
- `update_memory` / `delete_memory`：文件被更新/删除且 `memory.md` 同步；非法 path 拒绝；删除经审批。
- 提取 prompt 含 `## Existing memory files`。
- 定时去重：未达条件不触发；达 24h+5 会话触发并更新文件与 `memory.md` 及状态落盘。
- 现有记忆相关测试中被删功能的用例同步删除/改写，`go test ./...` 全绿。

## 7. 风险与缓解

- **工具寻址绕过 Dispatch**：记忆文件在工作区外，沿用 `spawn_subagent` 式 runtime 拦截，保持 `safeRelativeMemoryPath` 校验防目录穿越。
- **会话态内存丢失**：进程重启后会话去重态重建，最坏情况是一次重复注入，无正确性影响。
- **定时去重 LLM 失败**：best-effort，失败保留原数据并记录告警，不阻塞用户响应（沿用现有模式）。
- **截断字节边界**：25KB 截断按 rune 边界，避免破坏多字节字符。
