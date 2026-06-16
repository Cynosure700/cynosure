# TUI 显示与事件循环修复设计文档

## 背景

当前 TUI 在连续对话时出现两个现象：

1. 首轮回答只显示了开头片段，最终完整回答没有稳定落到界面上。
2. 后续输入可以出现在 transcript 中，但模型回复不再显示，看起来像“无回复”。

截图中第二轮、第三轮用户消息已经渲染，说明键盘输入与 transcript 刷新仍在工作；问题集中在运行中事件消费链路。

## 根因分析

TUI 的一次提问会同时启动两个 Bubble Tea command：

- `waitEvent()`：从 `m.events` 读取 runtime 流式事件。
- `respond()`：调用 runtime，并在结束后直接返回 `done` 事件给 Bubble Tea。

runtime 的最终 `assistant` 事件通过 `EventWriter` 写入 `m.events`，但 `done` 是由另一个 command 直接返回。两个 command 的消息到达顺序没有强约束，因此可能发生：

1. `done` 先被 Bubble Tea 处理，`m.running` 被置为 `false`。
2. 仍在 `m.events` 队列里的同一轮最终 `assistant` 事件没有新的 `waitEvent()` 继续消费。
3. 下一轮提问后 `generation` 自增，新的 `waitEvent()` 先读到上一轮残留事件。
4. 旧事件因 generation 不匹配被丢弃，但当前实现直接返回 `nil` command，没有继续等待本轮新事件。
5. 本轮后续回复事件留在 channel 中不再被消费，表现为后续对话无回复。

这也解释了文字显示不全：首轮流式 delta 已经显示了部分内容，但最终 `assistant` 完整内容可能被 `done` 抢先终止消费链路后滞留在事件队列里。

## 修复方案

采用事件顺序收敛 + 旧事件容错的最小修复：

1. `respond()` 不再把 `done` / `error` 作为独立 command 返回，而是统一写入 `m.events`。
   - runtime 的 `assistant_delta`、`reasoning_delta`、`assistant`、`meta` 与 `done` 共享同一 FIFO channel。
   - 保证最终 `assistant` 先于 `done` 被消费，避免完整回答丢失。
2. `Update(Event)` 遇到 generation 不匹配的旧事件时，如果当前仍在运行，继续返回 `m.waitEvent()`。
   - 即便历史版本或异常路径留下残留事件，也不会打断当前轮事件监听。
3. 保持现有 transcript 渲染结构不做额外 UI 重构，避免扩大改动面。

## 验证设计

新增回归用例：

- `TestModelKeepsWaitingAfterStaleGenerationEvent`：模拟当前轮运行时读到上一轮残留事件，断言旧事件被忽略后仍继续等待并能读到当前 generation 的新事件。
- `TestRespondSendsTerminalEventThroughEventQueue`：断言 `respond()` 不再直接返回终止事件，而是把错误/完成事件写入同一个事件队列，确保事件消费顺序可控。

执行验证：

- 先运行新增用例确认失败，证明测试覆盖了现有问题。
- 修复后运行 `go test ./internal/tui -count=1`。
- 结合截图路径做人工逻辑自检：首轮最终回答应完整替换流式片段，后续轮次即使先消费到旧事件也会继续等待本轮事件。
