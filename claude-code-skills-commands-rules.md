# Claude Code 中 Skills、Commands、Rules 三者区别与加载机制

## 一句话概括

| | Skills | Commands | Rules |
|---|---|---|---|
| **类比** | 专业顾问（需要时请来） | 快捷键/宏 | 法律法规（必须遵守） |
| **本质** | 按需加载的工作流指令包 | 预定义的提示词模板 | 始终生效的行为约束 |

---

## 加载机制（核心区别）

```
会话启动时：
  ① Rules → 全量加载，注入 system prompt（始终占用 token）
  ② CLAUDE.md → 全量加载，注入 system prompt
  ③ Skills → 仅加载 name + description 列表（索引）
  ④ Commands → 注册到斜杠命令系统（不加载内容）

每轮对话时：
  ⑤ Rules → 静默约束 Claude 所有输出
  ⑥ Skills → Claude 根据 description 语义匹配，按需通过 Skill tool 调用完整内容
  ⑦ Commands → 用户手动输入 /xxx 时才加载内容
```

| | Rules | Skills | Commands |
|---|---|---|---|
| **加载时机** | 启动时全量加载 | 按需加载（调用时才读完整内容） | 触发时加载（用户输入 `/xxx` 时） |
| **注入方式** | 注入 system prompt | Skill tool 返回内容 → 注入对话 | 内容作为用户消息注入 |
| **上下文开销** | 高（始终占用 token） | 低（只加载用到的） | 低（只加载触发的） |

---

## 触发方式

| | 触发方式 | 自动/手动 |
|---|---|---|
| **Skills** | Claude 根据 description **自动匹配** + 也可 `/skill-name` 手动调用 | ✅ 可自动 |
| **Commands** | 仅用户输入 `/command-name` | ❌ 仅手动 |
| **Rules** | 无需触发，始终生效 | ✅ 始终生效 |

---

## 文件格式与存放位置

### Skills — 最复杂，有 frontmatter，可带辅助资源

```
.claude/skills/<skill-name>/
├── SKILL.md              # 必需，YAML frontmatter + Markdown
│                         # frontmatter: name, description, allowed-tools
├── references/           # 可选，参考文档（按需读取）
└── agents/               # 可选，子代理配置
```

### Commands — 最简单，一个文件一个命令

```
.claude/commands/<command-name>.md    # 可选 frontmatter + Markdown
                                      # 支持 $ARGUMENTS 占位符
                                      # 子目录用 : 分隔 → db/migrate.md = /db:migrate
```

### Rules — 纯 Markdown，无额外语法

```
.claude/rules/<rule-name>.md          # 纯 Markdown
```

### 层级

| 层级 | Skills | Commands | Rules |
|---|---|---|---|
| 项目级 | `<project>/.claude/skills/` | `<project>/.claude/commands/` | `<project>/.claude/rules/` |
| 用户级 | `~/.claude/skills/` | `~/.claude/commands/` | `~/.claude/rules/` |
| 插件级 | `plugins/cache/.../skills/` | `plugins/cache/.../commands/` | ❌ 不支持 |

---

## 优先级

```
用户显式指令 > CLAUDE.md > Skill 指令 > Rules > 系统默认
项目级 > 插件级 > 用户级（全局）
```

---

## 适合放什么

| | 适合 | 不适合 |
|---|---|---|
| **Skills** | 复杂工作流（TDD、调试、brainstorming）、领域专业知识 | 简单的"必须用 Go"这类硬约束 |
| **Commands** | 重复性操作（review、deploy、fix-issue）、带参数的模板 | 需要自动触发的行为 |
| **Rules** | 项目规范、语言约束、编码标准、禁止事项 | 需要多步骤执行的工作流 |

---

## 简单判断

- 需要 Claude **主动判断何时用** → **Skill**
- 需要**用户手动触发** → **Command**
- 必须**始终遵守** → **Rule**
