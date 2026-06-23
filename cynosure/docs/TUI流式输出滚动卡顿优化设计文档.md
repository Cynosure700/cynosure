# TUI 流式输出滚动卡顿优化设计文档

## 背景

底部正在流式输出（模型逐 token 回答）时，用鼠标滚轮往上滑动会出现明显问题：

1. 滚动先无响应（延迟），过一会儿才突然跳动若干行（追帧）。
2. 整个滚动过程十分卡顿，无法平滑地往上查看历史。

该现象只在「边输出边滚动」时出现；空闲时滚动顺滑。说明问题集中在流式输出期间的渲染与事件消费链路，而非滚动逻辑本身。

## 根因分析

TUI 是 Bubble Tea 架构：所有消息（键盘、鼠标、流式事件、tick）都在**单一 `Update` goroutine** 上串行处理。流式输出期间存在三层叠加问题：

1. **每个 `assistant_delta` 都触发整屏全量重渲染。**
   `Update(Event)` 处理完每个事件后无条件调用 `refreshViewport`：

   ```
   refreshViewport → renderTranscript → renderMessages
       → 遍历全部 m.messages
       → 每条 assistant 消息都跑一次 m.renderer.Render()（glamour 完整 markdown 解析）
       + colorizeFileReferencesWithRestore（正则）+ wrapText
   ```

   runtime 在 `internal/agent/runtime/conversation_flow.go` 以极高频率发出 `assistant_delta`。**每来一个增量就把全部历史消息（含所有历史 assistant 块的 glamour 渲染）重算一遍**——单次开销 O(历史长度)，历史越长越慢。

2. **昂贵的全量渲染与鼠标滚轮抢占同一个串行队列。**
   `tea.MouseMsg`（滚轮）和洪水般的 `assistant_delta` 排在同一条 `Update` 消息流里。当 `Update` 忙于全量渲染时，滚轮消息只能排队 → 表现为「延迟」；积压一批后趁输出间隙被一次性消费 → 表现为「突然追帧、卡顿」。

3. **往上滚动时仍在做全量重渲染。**
   往上滚后 `autoFollow=false`，`refreshViewport` 不再 `GotoBottom`，但 `SetContent` 的全量重渲染照跑。因此恰恰在「边输出边往上看」时最卡。

## 优化方案

两项互补、低风险的改动，分别削减渲染**次数**与每次渲染**成本**，均不改变现有 UI 结构与交互语义。

### 方案 A —— 事件合并（coalescing）

将「每事件一渲染」改为「每批一渲染」：

- 把原 `Update(Event)` 内联的事件处理逻辑拆出 `applyEvent(Event)`：只应用单个事件到 Model 状态，**不触发渲染**；非当前 `generation` 的事件直接忽略。
- 新增 `drainPendingEvents()`：用 `select { case <-m.events: default: }` 非阻塞地把 channel 中**已就绪**的事件全部取出并 `applyEvent`。
- `Update(Event)` 流程变为：`applyEvent(当前事件)` → `drainPendingEvents()` → **只调用一次 `refreshViewport`** → 视运行状态决定是否继续 `waitEvent()`。

效果：一轮模型输出产生的大量连续 `assistant_delta` 被压成一帧，渲染次数从「每个 token 一次」降到「每批一次」，`Update` 循环迅速腾出来响应鼠标滚轮，消除排队延迟与追帧抖动。`done`/`error` 等终止事件若混在批次中，也会在本批 `applyEvent` 时正确清除 `running`，drain 结束后不再继续等待。

### 方案 B —— 按消息渲染缓存

将「每帧重算全部历史」改为「只重算变化的那条」：

- 新增 `messageRenderCache`（`map[string]string`，指针语义，使 Bubble Tea 的 Model 值拷贝之间共享同一份缓存）。
- `renderMessages` 改走 `renderCachedMessage`：缓存 key 由 `messageRenderKey` 生成，覆盖所有影响渲染输出的字段 —— `(宽度, role, content, reasoning_content, 以及 ToolCall 的各字段)`。
  - 历史消息内容不变 → 命中缓存，跳过 glamour/正则/wrap。
  - 流式时只有最后一条 assistant 在变 → 仅它 miss 并重算。
  - 单帧渲染成本从 O(全部历史) 降到 O(1 条)。
- `refreshViewport` 中加 `pruneRenderCache`：按当前消息列表的 live key 集合清理过期缓存项。流式内容增长每次都会产生新 key、窗口宽度变化也会使旧 key 失效，pruning 保证缓存不会无限膨胀。

渲染输出是 `(消息内容, 宽度)` 的纯函数，因此按内容缓存是安全的；宽度变化（`WindowSizeMsg`）天然产生新 key，整体等价于失效重建。

### 协同效果

A 砍掉渲染**次数**，B 砍掉每次渲染**成本**。两者叠加后，流式输出期间 `Update` 循环负载大幅下降，鼠标滚轮不再被饿死，向上滚动即时跟手。

## 验证设计

新增回归用例（`internal/tui/events_test.go`）：

- `TestUpdateCoalescesPendingStreamingEventsIntoOneFrame`：预填多个已就绪 `assistant_delta`，断言它们在一次 `Update` 内被合并应用、channel 被清空、内容正确拼接。
- `TestUpdateStopsDrainingAtTerminalDoneEvent`：drain 到 `done` 时应停止 `running`，且不再返回继续等待的 command。
- `TestRenderCacheReusesUnchangedMessageRenders`：篡改缓存值后再次渲染应返回被篡改值，证明命中缓存而非重算。
- `TestRenderCacheInvalidatesWhenMessageContentChanges`：内容变更后应重算出新内容，不返回过期缓存。
- `TestRefreshViewportPrunesStaleRenderCacheEntries`：模拟流式内容增长，断言旧 key 在下一次刷新时被清理、缓存项数收敛到当前消息数。

同时更新因合并行为改变而过时的 `TestModelKeepsWaitingAfterStaleGenerationEvent`：合并后 fresh 事件在同一次 `Update` 内被 drain 应用（而非延迟到 `cmd()`），按「以代码为事实」修正旧断言。

执行验证：

- `go build ./...`、`go vet ./internal/tui/` 通过。
- `go test ./internal/tui/ -count=1`：除预先存在、与本次改动无关的 `TestTodoWriteToolMessageRendersCheckboxList`（todo 复选框颜色断言漂移）外全部通过。
- 结合现有滚动相关用例（`TestMouseWheelUpScrollsTranscriptHistory`、`TestManualScrollUpPreventsAutoScrollOnNewContent`、`TestBottomPositionAutoFollowsNewContent`、`TestScrollingBackToBottomRestoresAutoFollow`）确认 autoFollow 与滚动语义未受影响。

## 影响范围

- 仅改动 `internal/tui/app.go`（事件分发、渲染缓存）与 `internal/tui/events_test.go`（用例）。
- 不改变事件协议、UI 结构、键鼠交互语义与 autoFollow 行为。
- 缓存为进程内、随会话生命周期存在，pruning 保证内存有界。
