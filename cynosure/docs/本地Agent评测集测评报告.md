# 本地 Agent 评测集测评报告

## 结论

本评测集为 `cynosure` 本地 coding-agent 设计，参考 SWE-bench / Terminal-Bench 等公开基准的做法，并按三档难度（standard / hard / expert）形成能力梯度。全部用真实 `cynosure` agent 跑通。

- 测评时间：2026-06-27
- 被测 Agent：本机 `cynosure`（模型 `glm-5.2`）
- 评测用例：**60 个**（standard 20 + hard 20 + expert 20）
- 覆盖方面：bugfix 11、feature 12、documentation 9、terminal 10、multi_file 9、refactor 9
- 真实测评结果：**60/60 passed**（standard 20/20 + hard 20/20 + expert 20/20）
- 判分方式：只采信 verifier 的确定性结果，不采信模型最终文本自述
- Baseline 非空测校验：全部 60 个 case 的初始 fixture 在未运行 Agent 时均无法通过 verifier；hard/expert case 加入前逐个用参考解验证过「可解」

## 评测方法

评测集参考行业常用 coding-agent 评测方式设计：

- **SWE-bench 风格**：给 Agent 一个 issue-like 任务和一个独立代码仓库，要求修改代码并通过测试。
- **Terminal-Bench 风格**：给 Agent 一个终端/文件处理任务，最终只通过验证命令检查工作区状态。

质量原则（对标公开 benchmark）：

- **确定性判分**：结果只由测试或 verifier 决定，不由模型回答文本决定。
- **隔离环境**：每个 case 从独立 fixture 复制出临时工作区运行。
- **初始失败**：未运行 Agent 的初始 fixture 必须无法通过 verifier。
- **行为优先**：软件工程类 case 尽量用单元测试 / CLI 输出验证行为，避免过度绑定实现细节。
- **难度梯度**：standard（基础单/双文件）→ hard（更贴近真实仓库）→ expert（跨多文件、长链路推理、隐藏边界陷阱）。

每个 case 都包含：独立工作区 `evals/suite/cases/<case-id>/repo`、提示词（`evals/suite/manifest.json`）、确定性 verifier（`evals/suite/verifiers/*.py`）。运行器 `evals/run_cynosure_eval.py` 复制 fixture 到 `evals/results/<run-id>/workspaces/`，写入 eval 专用 `.cynosure/settings.json`（`bypassPermissions`），通过 `expect` 启动真实 `cynosure --cwd <workspace>`，运行期每 5 秒调用 verifier，一旦通过即结束并记录。

## 覆盖方面

按「方面 × 难度」分布（每个方面在每档至少 3 个）：

| 方面 | standard | hard | expert | 合计 | 覆盖点 |
| --- | ---: | ---: | ---: | ---: | --- |
| bugfix | 3 | 4 | 4 | 11 | 业务逻辑/边界缺陷、库存超卖、文本归一化、区间合并、LRU、semver、分页、RPN、调度器、Decimal 税额、Dijkstra |
| feature | 4 | 4 | 4 | 12 | TTL cache、限流、retry、CSV CLI、表达式求值、令牌桶、拓扑排序、词频、查询构建器、CSV 流、栈式 VM、Trie 路由 |
| documentation | 3 | 3 | 3 | 9 | runbook/OpenAPI/CLI 同步、CHANGELOG、配置表、错误目录、双源对账、CLI 矩阵、指标术语表 |
| terminal | 4 | 3 | 3 | 10 | 日志统计、CSV 聚合、JSONL 汇总、Markdown 索引、日志关联、JSON 汇总、定宽解析、对账、依赖图、INI 合并 |
| multi_file | 3 | 3 | 3 | 9 | config/policy/order 协同、事件溯源、中间件链、状态机、DI 容器、事件总线、ORM 校验 |
| refactor | 3 | 3 | 3 | 9 | 规范化去重、Decimal 重构、格式化 helper、god-function 拆分、校验收敛、策略表、分层、接口抽取、规则去重 |
| **合计** | **20** | **20** | **20** | **60** | 覆盖代码阅读、TDD 修复、Go/Python 跨语言、功能实现、文档同步、终端数据处理、多文件协作、结构性约束 |

## Case 清单

全部 60 个 case（`✓ tests` 表示用单元测试 / `go test` 行为判分，`◆ contract` 表示用 contract verifier 校验文件/CLI 输出）：

| # | Case | 方面 | 难度 | 判分 | 评测能力 |
| --: | --- | --- | --- | --- | --- |
| 1 | `py-ledger-bug` | bugfix | standard | ✓ tests | 修复 Python 加权平均业务逻辑 |
| 2 | `py-inventory-reservation-bug` | bugfix | standard | ◆ contract | 防止库存超卖和非法数量 |
| 3 | `go-slug-bug` | bugfix | standard | ◆ contract | 修复 Go slug 标点/连续分隔符处理 |
| 4 | `go-cache-feature` | feature | standard | ✓ tests | 实现 Go TTL cache |
| 5 | `py-rate-limiter-feature` | feature | standard | ◆ contract | 实现 Python per-user 时间窗限流 |
| 6 | `go-retry-feature` | feature | standard | ◆ contract | 实现 Go retry helper |
| 7 | `py-cli-csv-feature` | feature | standard | ◆ contract | 实现 Python CSV 汇总 CLI |
| 8 | `docs-runbook-update` | documentation | standard | ✓ tests | 以配置/脚本同步 runbook |
| 9 | `docs-api-migration` | documentation | standard | ◆ contract | 以 OpenAPI 同步 API 文档 |
| 10 | `docs-cli-reference` | documentation | standard | ◆ contract | 以源码命令表同步 CLI 文档 |
| 11 | `shell-log-pipeline` | terminal | standard | ✓ tests | 统计 HTTP 日志 |
| 12 | `terminal-csv-aggregate` | terminal | standard | ◆ contract | 聚合 CSV 销售数据 |
| 13 | `terminal-jsonl-errors` | terminal | standard | ◆ contract | 汇总 JSONL 错误日志 |
| 14 | `terminal-markdown-index` | terminal | standard | ◆ contract | 生成 Markdown 文档索引 |
| 15 | `multifile-config-refactor` | multi_file | standard | ✓ tests | Python config/app 多文件重构 |
| 16 | `py-auth-policy-multifile` | multi_file | standard | ◆ contract | Python policy/api 行为协同 |
| 17 | `go-order-total-multifile` | multi_file | standard | ◆ contract | Go order/pricing 协同 |
| 18 | `python-refactor-validation` | refactor | standard | ✓ tests | Python 用户规范化去重 |
| 19 | `py-money-refactor` | refactor | standard | ◆ contract | Python 金额 Decimal 重构 |
| 20 | `go-formatter-refactor` | refactor | standard | ◆ contract | Go 用户格式化 helper 重构 |
| 21 | `hard-py-interval-merge-bug` | bugfix | hard | ◆ contract | 区间合并（未排序/相接/嵌套/不可变输入） |
| 22 | `hard-go-lru-bug` | bugfix | hard | ◆ contract | LRU 缓存淘汰顺序 |
| 23 | `hard-py-semver-bug` | bugfix | hard | ◆ contract | 语义化版本比较（prerelease 规则） |
| 24 | `hard-go-pagination-bug` | bugfix | hard | ◆ contract | 1-based 分页边界 |
| 25 | `hard-py-expr-eval-feature` | feature | hard | ◆ contract | 表达式求值器（禁用 eval） |
| 26 | `hard-go-tokenbucket-feature` | feature | hard | ◆ contract | 令牌桶限流 |
| 27 | `hard-py-graph-topo-feature` | feature | hard | ◆ contract | 确定性拓扑排序 |
| 28 | `hard-go-wordcount-cli-feature` | feature | hard | ◆ contract | 词频统计 CLI |
| 29 | `hard-docs-changelog` | documentation | hard | ◆ contract | 按类型分组生成 CHANGELOG |
| 30 | `hard-docs-config-table` | documentation | hard | ◆ contract | 从 JSON schema 生成配置表 |
| 31 | `hard-docs-error-catalog` | documentation | hard | ◆ contract | 从源码生成错误码目录 |
| 32 | `hard-terminal-log-join` | terminal | hard | ◆ contract | 双日志按 id 关联 |
| 33 | `hard-terminal-json-rollup` | terminal | hard | ◆ contract | 嵌套 JSON 汇总 |
| 34 | `hard-terminal-access-report` | terminal | hard | ◆ contract | 定宽日志解析 |
| 35 | `hard-multifile-py-eventsourcing` | multi_file | hard | ◆ contract | 事件溯源账户 |
| 36 | `hard-multifile-go-middleware` | multi_file | hard | ◆ contract | HTTP 中间件链 |
| 37 | `hard-multifile-py-statemachine` | multi_file | hard | ◆ contract | 文档状态机 |
| 38 | `hard-refactor-py-pipeline` | refactor | hard | ◆ contract | god-function 拆分 |
| 39 | `hard-refactor-go-validation` | refactor | hard | ◆ contract | Go 重复校验收敛 |
| 40 | `hard-refactor-py-strategy` | refactor | hard | ◆ contract | if/elif 改 dispatch 注册表 |
| 41 | `expert-py-rpn-pipeline-bug` | bugfix | expert | ◆ contract | 跨文件 RPN 求值（操作数顺序/多余操作数陷阱） |
| 42 | `expert-go-scheduler-bug` | bugfix | expert | ◆ contract | 优先级任务调度（依赖/环/未知依赖） |
| 43 | `expert-py-decimal-tax-bug` | bugfix | expert | ◆ contract | Decimal 税额（按小计计税+单次舍入） |
| 44 | `expert-go-graph-dijkstra-bug` | bugfix | expert | ◆ contract | Dijkstra（贪心首边陷阱） |
| 45 | `expert-py-query-builder-feature` | feature | expert | ◆ contract | 链式查询构建器（稳定排序/缺失字段） |
| 46 | `expert-go-csv-stream-feature` | feature | expert | ◆ contract | CSV 引号解析 + 短路 pipeline |
| 47 | `expert-py-mini-vm-feature` | feature | expert | ◆ contract | 栈式 VM（跳转/操作数顺序/错误） |
| 48 | `expert-go-trie-router-feature` | feature | expert | ◆ contract | Trie 路由（静态>参数>通配优先级） |
| 49 | `expert-docs-openapi-sync` | documentation | expert | ◆ contract | OpenAPI+handler 双源对账 |
| 50 | `expert-docs-cli-matrix` | documentation | expert | ◆ contract | 命令/全局 flag 合并矩阵 |
| 51 | `expert-docs-metrics-glossary` | documentation | expert | ◆ contract | 指标术语表（单位/类型不一致修正） |
| 52 | `expert-terminal-ledger-recon` | terminal | expert | ◆ contract | 双账本对账（匹配/差异/单边） |
| 53 | `expert-terminal-dep-graph` | terminal | expert | ◆ contract | 依赖图分析（深度/环/最被依赖） |
| 54 | `expert-terminal-ini-merge` | terminal | expert | ◆ contract | INI 分层合并优先级 |
| 55 | `expert-multifile-py-di-container` | multi_file | expert | ◆ contract | DI 容器（递归注入+环检测） |
| 56 | `expert-multifile-go-eventbus` | multi_file | expert | ◆ contract | 事件总线（FIFO+保序退订） |
| 57 | `expert-multifile-py-orm-validate` | multi_file | expert | ◆ contract | ORM 校验（bool 不当 int 陷阱） |
| 58 | `expert-refactor-py-layering` | refactor | expert | ◆ contract | I/O 与计算分层 |
| 59 | `expert-refactor-go-interface-extract` | refactor | expert | ◆ contract | Go 接口抽取解耦 |
| 60 | `expert-refactor-py-dedup-rules` | refactor | expert | ◆ contract | 重复定价规则收敛 |

## 测评结果汇总（60/60）

被测 Agent：本机 `cynosure`（模型 `glm-5.2`）。三档结果合并在下表，结果列均为 `pass`（`exit_reason=verifier_passed`）：

| # | Case | 方面 | 难度 | 结果 | Agent 用时(s) |
| --: | --- | --- | --- | --- | ---: |
| 1 | `py-ledger-bug` | bugfix | standard | pass | 46.23 |
| 2 | `py-inventory-reservation-bug` | bugfix | standard | pass | 41.10 |
| 3 | `go-slug-bug` | bugfix | standard | pass | 38.86 |
| 4 | `go-cache-feature` | feature | standard | pass | 39.26 |
| 5 | `py-rate-limiter-feature` | feature | standard | pass | 67.15 |
| 6 | `go-retry-feature` | feature | standard | pass | 38.94 |
| 7 | `py-cli-csv-feature` | feature | standard | pass | 82.06 |
| 8 | `docs-runbook-update` | documentation | standard | pass | 70.88 |
| 9 | `docs-api-migration` | documentation | standard | pass | 20.37 |
| 10 | `docs-cli-reference` | documentation | standard | pass | 30.57 |
| 11 | `shell-log-pipeline` | terminal | standard | pass | 30.34 |
| 12 | `terminal-csv-aggregate` | terminal | standard | pass | 25.50 |
| 13 | `terminal-jsonl-errors` | terminal | standard | pass | 25.44 |
| 14 | `terminal-markdown-index` | terminal | standard | pass | 30.47 |
| 15 | `multifile-config-refactor` | multi_file | standard | pass | 36.08 |
| 16 | `py-auth-policy-multifile` | multi_file | standard | pass | 41.13 |
| 17 | `go-order-total-multifile` | multi_file | standard | pass | 26.93 |
| 18 | `python-refactor-validation` | refactor | standard | pass | 30.78 |
| 19 | `py-money-refactor` | refactor | standard | pass | 61.86 |
| 20 | `go-formatter-refactor` | refactor | standard | pass | 38.80 |
| 21 | `hard-py-interval-merge-bug` | bugfix | hard | pass | 25.62 |
| 22 | `hard-go-lru-bug` | bugfix | hard | pass | 44.06 |
| 23 | `hard-py-semver-bug` | bugfix | hard | pass | 30.73 |
| 24 | `hard-go-pagination-bug` | bugfix | hard | pass | 27.90 |
| 25 | `hard-py-expr-eval-feature` | feature | hard | pass | 30.77 |
| 26 | `hard-go-tokenbucket-feature` | feature | hard | pass | 38.77 |
| 27 | `hard-py-graph-topo-feature` | feature | hard | pass | 46.24 |
| 28 | `hard-go-wordcount-cli-feature` | feature | hard | pass | 42.58 |
| 29 | `hard-docs-changelog` | documentation | hard | pass | 20.32 |
| 30 | `hard-docs-config-table` | documentation | hard | pass | 25.39 |
| 31 | `hard-docs-error-catalog` | documentation | hard | pass | 25.37 |
| 32 | `hard-terminal-log-join` | terminal | hard | pass | 55.91 |
| 33 | `hard-terminal-json-rollup` | terminal | hard | pass | 25.38 |
| 34 | `hard-terminal-access-report` | terminal | hard | pass | 25.40 |
| 35 | `hard-multifile-py-eventsourcing` | multi_file | hard | pass | 30.80 |
| 36 | `hard-multifile-go-middleware` | multi_file | hard | pass | 22.57 |
| 37 | `hard-multifile-py-statemachine` | multi_file | hard | pass | 35.99 |
| 38 | `hard-refactor-py-pipeline` | refactor | hard | pass | 66.74 |
| 39 | `hard-refactor-go-validation` | refactor | hard | pass | 33.33 |
| 40 | `hard-refactor-py-strategy` | refactor | hard | pass | 92.32 |
| 41 | `expert-py-rpn-pipeline-bug` | bugfix | expert | pass | 62.00 |
| 42 | `expert-go-scheduler-bug` | bugfix | expert | pass | 50.37 |
| 43 | `expert-py-decimal-tax-bug` | bugfix | expert | pass | 31.03 |
| 44 | `expert-go-graph-dijkstra-bug` | bugfix | expert | pass | 34.07 |
| 45 | `expert-py-query-builder-feature` | feature | expert | pass | 169.39 |
| 46 | `expert-go-csv-stream-feature` | feature | expert | pass | 34.37 |
| 47 | `expert-py-mini-vm-feature` | feature | expert | pass | 25.81 |
| 48 | `expert-go-trie-router-feature` | feature | expert | pass | 67.47 |
| 49 | `expert-docs-openapi-sync` | documentation | expert | pass | 30.67 |
| 50 | `expert-docs-cli-matrix` | documentation | expert | pass | 56.11 |
| 51 | `expert-docs-metrics-glossary` | documentation | expert | pass | 20.43 |
| 52 | `expert-terminal-ledger-recon` | terminal | expert | pass | 51.06 |
| 53 | `expert-terminal-dep-graph` | terminal | expert | pass | 45.96 |
| 54 | `expert-terminal-ini-merge` | terminal | expert | pass | 20.44 |
| 55 | `expert-multifile-py-di-container` | multi_file | expert | pass | 25.96 |
| 56 | `expert-multifile-go-eventbus` | multi_file | expert | pass | 34.00 |
| 57 | `expert-multifile-py-orm-validate` | multi_file | expert | pass | 57.20 |
| 58 | `expert-refactor-py-layering` | refactor | expert | pass | 36.34 |
| 59 | `expert-refactor-go-interface-extract` | refactor | expert | pass | 47.84 |
| 60 | `expert-refactor-py-dedup-rules` | refactor | expert | pass | 36.11 |

按难度 × 方面汇总（通过/总数）：

| 方面 | standard | hard | expert | 合计 |
| --- | ---: | ---: | ---: | ---: |
| bugfix | 3/3 | 4/4 | 4/4 | 11/11 |
| feature | 4/4 | 4/4 | 4/4 | 12/12 |
| documentation | 3/3 | 3/3 | 3/3 | 9/9 |
| terminal | 4/4 | 3/3 | 3/3 | 10/10 |
| multi_file | 3/3 | 3/3 | 3/3 | 9/9 |
| refactor | 3/3 | 3/3 | 3/3 | 9/9 |
| **合计** | **20/20** | **20/20** | **20/20** | **60/60** |

- 全部 60 个 case 的 `exit_reason` 均为 `verifier_passed`，无超时。
- 单 case 用时区间：standard 20–82s、hard 20–92s、expert 20–169s；三档累计约 822.8s / 818s / 936.6s。
- 最慢为 `expert-py-query-builder-feature`（169s），说明复杂跨文件实现确实更吃推理。

## 行业标准合理性评估

按难度档评估（每个 case 都满足：确定性判分、隔离环境、初始 fixture 必失败、加入前已验证可解）：

| 难度 | 对标标准 | 合理性评估 |
| --- | --- | --- |
| standard | SWE-bench / Terminal-Bench 风格 | 单/双文件基础任务，覆盖六大方面；软件类用 `unittest`/`go test` 判行为，文档/终端类用 contract verifier 校验产物，初始 fixture 全部失败，非空测。 |
| hard | SWE-bench / Terminal-Bench 风格 | 更贴近真实仓库的算法/数据/重构任务（LRU、令牌桶、semver、拓扑、事件溯源、状态机等），含明确边界用例；加入前逐个用参考解验证可解。 |
| expert | SWE-bench Verified 风格 | 强调跨多文件（2–4 文件）、长链路推理（VM 跳转、Dijkstra 多跳、依赖图深度）与隐藏边界陷阱（贪心首边、`bool` 非 `int`、RPN 多余操作数、按小计单次舍入、INI 覆盖优先级、文档单位/类型不一致）。 |

判分类型分布：`✓ tests`（单元测试 / `go test` 行为判分）11 个，`◆ contract`（contract verifier 校验文件/CLI/结构）49 个。两类都不读取模型最终文本，仅依据 verifier 退出码。

## 运行证据

三档真实测评的结果目录（均含 `results.json`、`REPORT.md`、每个 case 的 `*.transcript.txt` 与 `*.verifier.json`）：

```text
evals/results/20260626-155035   # standard 20/20
evals/results/20260626-173140   # hard 20/20
evals/results/20260627-161129   # expert 20/20
```

示例（单个 case 的产物）：

```text
evals/results/20260627-161129/expert-go-graph-dijkstra-bug.transcript.txt
evals/results/20260627-161129/expert-go-graph-dijkstra-bug.verifier.json
```

Baseline（`--skip-agent`，不驱动 Agent）非空测证据保留在 `evals/results/20260626-154429/`，证明未运行 Agent 的初始 fixture 均不能通过 verifier。

## 复现方式

从 Go module 根目录运行：

```bash
cd /Users/bytedance/golang_pro/cynosure/cynosure

# 1) 校验评测集结构（数量、难度档、每方面≥3、fixture/verifier 路径有效）
go test ./evals -run TestEvaluationSuiteManifestIsComplete -count=1

# 2) 全量真实测评（60 个 case）
python3 evals/run_cynosure_eval.py

# 3) 仅跑某一档（按 difficulty 过滤生成 --case 列表）
python3 evals/run_cynosure_eval.py \
  $(python3 -c "import json;[print('--case',c['id']) for c in json.load(open('evals/suite/manifest.json'))['cases'] if c.get('difficulty')=='expert']")

# 4) 仅跑单个 case
python3 evals/run_cynosure_eval.py --case expert-go-graph-dijkstra-bug

# 5) 不跑 Agent，仅验证初始 fixture 非空测
python3 evals/run_cynosure_eval.py --skip-agent
```

## 评估

评测集覆盖本地 coding agent 的核心能力面（代码阅读、TDD 修复、Go/Python 跨语言、功能实现、文档同步、终端数据处理、多文件协作、结构性约束），并通过 standard/hard/expert 三档形成清晰梯度。全部 60 个 case 都有独立 verifier、初始必失败、加入前已验证可解。本次用真实 `cynosure`（`glm-5.2`）跑通全部三档，合计 **60/60 passed**。

局限：

- 规模 60 个、三档梯度清晰，适合作为本地 smoke/regression/能力分层 benchmark；仍不等价于 SWE-bench Verified 级别的大规模公开基准。
- runner 通过 TUI + `expect` 驱动真实 `cynosure`，能验证真实交互链路，但耗时和 transcript 可读性受 TUI 渲染影响。
- expert 档仍全部通过（区分度尚未触顶），单 case 耗时已明显上升（最慢 169s）；若需进一步压测，可加入更大仓库、多轮交互依赖或需要运行期调试（reproduce→fix→verify）的 case。
