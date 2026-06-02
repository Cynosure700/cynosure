# 推理过程自动展开行为设计

## 1. 背景

当前前端会在 assistant 消息存在 `reasoning_content` 时渲染推理过程折叠块：`web/src/App.tsx:659`。

SSE 流式处理已经支持：

- `reasoning_delta`：追加到临时 assistant 消息的 `reasoning_content`：`web/src/App.tsx:313`
- `assistant_delta`：追加到临时 assistant 消息的 `content`：`web/src/App.tsx:310`
- SSE 事件分发：`web/src/api.ts:153`

但当前推理过程使用原生 `<details>`，没有受控展开状态：`web/src/App.tsx:660`。这导致前端无法区分“正在流式生成的推理过程”和“用户重新打开历史会话看到的推理过程”。

新的目标是：

1. SSE 有推理内容流入时，当前正在生成的推理过程自动展开。
2. 用户重新打开已有会话时，历史推理过程默认不自动展开。
3. 用户仍然可以手动展开或收起推理过程。

## 2. 目标与非目标

### 2.1 目标

1. 流式生成期间，只要收到 `reasoning_delta`，对应 assistant 消息的推理过程自动展开。
2. 历史消息加载完成后，即使存在 `reasoning_content`，推理过程默认保持收起。
3. 用户手动收起正在生成的推理过程后，本轮后续 `reasoning_delta` 不再强制展开，尊重用户选择。
4. 用户手动展开历史推理过程后，保持当前页面内的手动展开状态。
5. 不修改后端 SSE 协议，不新增数据库字段。

### 2.2 非目标

1. 不持久化用户对某条消息推理过程的展开/收起偏好。
2. 不调整推理过程内容的格式化展示。
3. 不改变 `reasoning_content` 的服务端生成、存储和下发逻辑。
4. 不引入新的状态管理库。

## 3. 当前问题分析

当前消息渲染逻辑中，推理过程直接使用：

```tsx
<details className="message-reasoning">
    <summary>推理过程</summary>
    <div>{message.reasoning_content}</div>
</details>
```

原生 `<details>` 的展开状态由浏览器内部维护。因为没有显式 `open` 状态，应用无法基于“是否正在接收 SSE 推理增量”控制展开，也无法在用户手动收起后阻止后续自动展开。

此外，历史会话加载通过 `api.getConversation` 一次性设置 `messages`：`web/src/App.tsx:207`。这些消息不是流式新增，不应触发自动展开。

## 4. 方案对比

### 方案 A：只给 `<details>` 增加 `open={sending}`

做法：当全局 `sending` 为 true 时打开所有推理过程。

优点：实现最简单。

缺点：

- 会展开所有 assistant 消息，而不是只展开当前流式消息。
- 用户手动收起后，后续渲染会再次被 `sending` 打开。
- 无法处理多个会话切换中的精确状态。

结论：不推荐。

### 方案 B：维护“推理展开消息 ID 集合”【推荐】

做法：前端增加一个页面内状态 `expandedReasoningIds`，记录哪些消息的推理过程处于展开状态。

核心规则：

- 新建临时 assistant 消息时，先不展开。
- 收到该临时消息的第一段 `reasoning_delta` 时，将该消息 ID 加入 `expandedReasoningIds`。
- 用户点击 summary 手动展开时，加入集合。
- 用户点击 summary 手动收起时，从集合移除，并把该消息加入 `manualClosedReasoningIds`，表示本轮后续增量不再自动展开。
- `assistant` 终态事件把临时消息 ID 替换为服务端 message ID 时，同步迁移展开状态。
- 历史会话加载时清空这些页面内集合，因此历史推理默认收起。

优点：

- 自动展开只作用于收到 SSE 的当前消息。
- 能尊重用户手动收起。
- 历史消息默认收起，符合需求。
- 不依赖后端协议变更。

缺点：

- 需要把推理折叠块从原生非受控 `<details>` 改成受控组件。
- 需要在临时 ID 替换为真实 ID 时迁移状态。

结论：推荐。

### 方案 C：在 `ChatMessage` 上增加前端临时字段

做法：给消息对象增加如 `reasoningExpanded`、`reasoningTouched` 等 UI 字段。

优点：状态跟随消息对象，渲染时读取直接。

缺点：

- `ChatMessage` 同时承载服务端数据和前端 UI 状态，边界不清晰。
- `api.getConversation`、`streamConversation`、`refreshAll` 等路径都可能覆盖消息对象，容易丢失 UI 状态。
- 后续扩展更多 UI 状态时会污染数据模型。

结论：不推荐。

## 5. 推荐设计

采用方案 B：维护页面内展开状态集合，并把推理过程渲染改为受控行为。

### 5.1 新增状态

在 `App` 中新增两个状态：

```ts
const [expandedReasoningIds, setExpandedReasoningIds] = useState<Set<string>>(new Set());
const [manualClosedReasoningIds, setManualClosedReasoningIds] = useState<Set<string>>(new Set());
```

含义：

- `expandedReasoningIds`：当前页面中处于展开状态的推理过程消息 ID。
- `manualClosedReasoningIds`：用户手动收起过的消息 ID。只用于阻止同一轮 SSE 后续增量再次自动展开。

如果某条 assistant 消息没有服务端 `id`，使用渲染时的稳定 key。流式临时消息已经有 `streaming-${Date.now()}` ID：`web/src/App.tsx:299`，可直接使用。

### 5.2 自动展开规则

在 `onReasoningDelta` 中，当前逻辑会把增量追加到临时 assistant 消息：`web/src/App.tsx:313`。

设计调整为：

1. 先检查 `activeConversationIdRef.current === requestConversationId`，保持当前会话保护逻辑。
2. 当 `payload.content` 非空时：
   - 如果 `streamingAssistantId` 不在 `manualClosedReasoningIds` 中，则加入 `expandedReasoningIds`。
   - 如果用户已经手动收起过，则只追加内容，不自动展开。
3. 不对历史消息触发自动展开，因为历史消息由 `getConversation` 一次性加载，不经过 `onReasoningDelta`。

### 5.3 用户手动控制规则

将推理块封装为内部渲染函数或小组件，例如：

```tsx
<details
    className="message-reasoning"
    open={expandedReasoningIds.has(reasoningKey)}
    onToggle={(event) => handleReasoningToggle(reasoningKey, event.currentTarget.open)}
>
    <summary>推理过程</summary>
    <div>{message.reasoning_content}</div>
</details>
```

`handleReasoningToggle` 规则：

- `open === true`：加入 `expandedReasoningIds`，并从 `manualClosedReasoningIds` 移除。
- `open === false`：从 `expandedReasoningIds` 移除，并加入 `manualClosedReasoningIds`。

注意：受控 `<details>` 的 `onToggle` 可能由程序设置 `open` 触发。实现时需要避免重复 setState 造成无意义渲染；可以在 setter 中比较集合是否已经包含目标 ID。

### 5.4 临时消息 ID 迁移

流式期间使用 `streamingAssistantId` 作为临时 assistant 消息 ID：`web/src/App.tsx:299`。最终 `assistant` 事件可能把它替换为 `payload.message_id`：`web/src/App.tsx:323`。

迁移规则：

1. 如果 `payload.message_id` 存在且不同于 `streamingAssistantId`：
   - `expandedReasoningIds` 中存在 `streamingAssistantId` 时，删除临时 ID，加入真实 ID。
   - `manualClosedReasoningIds` 中存在 `streamingAssistantId` 时，同样迁移到真实 ID。
2. 如果 `payload.message_id` 不存在，继续使用临时 ID。

这样可以避免最终消息 ID 替换后，推理过程突然收起或丢失用户手动状态。

### 5.5 会话切换与历史加载

当 `activeConversationId` 变化时，当前实现会清空 `messages` 并加载历史消息：`web/src/App.tsx:197`。

设计调整为：

1. 在 `activeConversationId` 变化时清空：
   - `expandedReasoningIds`
   - `manualClosedReasoningIds`
2. 历史消息加载后不自动加入 `expandedReasoningIds`。

效果：用户重新打开会话时，历史推理过程默认收起。

### 5.6 发送完成后的行为

`done` 事件或 `assistant` 终态事件不需要主动收起推理过程。

原因：

- 用户正在观看本轮推理过程时，生成结束后保持展开更自然。
- 用户如果不想看，可以手动收起。
- 下次切换会话再回来时，因为状态不持久化，历史推理会默认收起。

## 6. 数据流

### 6.1 新消息流式生成

1. 用户发送消息。
2. 前端创建 user 消息和临时 assistant 消息。
3. 收到第一段 `reasoning_delta`。
4. 前端追加 `reasoning_content`，并自动展开该临时 assistant 的推理过程。
5. 用户如果手动收起，则该消息进入 `manualClosedReasoningIds`。
6. 后续 `reasoning_delta` 继续追加内容，但不再自动展开。
7. 收到 `assistant` 终态事件后，迁移临时 ID 的展开/手动收起状态到真实 message ID。

### 6.2 打开历史会话

1. 用户切换或重新打开会话。
2. 前端清空消息和推理 UI 状态。
3. `getConversation` 返回历史消息。
4. 消息中即使存在 `reasoning_content`，也因为不在 `expandedReasoningIds` 中而默认收起。
5. 用户可手动展开某条历史推理过程。

## 7. 测试设计

当前项目没有前端单元测试框架，建议用以下方式验证：

### 7.1 类型与构建验证

运行：

```bash
npm --prefix web run build
```

验证 TypeScript 类型和生产构建通过。

### 7.2 手工验证用例

1. 新发起一轮会话，并触发模型输出 `reasoning_delta`：
   - 推理过程自动展开。
2. 在推理继续流式输出时，手动收起推理过程：
   - 后续增量继续追加。
   - 推理过程保持收起，不自动弹开。
3. 手动再次展开：
   - 能看到已累计的推理内容。
4. 切换到其他会话，再切回当前会话：
   - 历史推理过程默认收起。
5. 在无 `reasoning_content` 的 assistant 消息上：
   - 不渲染推理过程区域。

## 8. 实施步骤

1. 在 `App` 中新增 `expandedReasoningIds` 和 `manualClosedReasoningIds` 状态。
2. 新增集合更新 helper，避免直接原地修改 `Set`。
3. 在 `activeConversationId` 变化时清空推理 UI 状态。
4. 在 `onReasoningDelta` 中根据规则自动展开当前流式 assistant 消息。
5. 在 `onAssistant` 终态事件中迁移临时 ID 到真实 ID。
6. 将推理 `<details>` 改为受控 `open`，并接入 `onToggle`。
7. 运行构建验证并按手工用例检查行为。

## 9. 风险与边界

1. `Set` 必须通过创建新实例更新，不能原地修改，否则 React 可能不触发渲染。
2. `details` 的 `onToggle` 会在展开状态变化时触发，需要避免状态更新循环。
3. 如果某些历史消息没有 `id`，需要使用渲染 key 兜底；但流式消息必须使用 `streamingAssistantId`，否则无法准确自动展开。
4. 如果用户在流式过程中切换会话，现有 `activeConversationIdRef` 保护应继续阻止旧会话事件影响当前页面：`web/src/App.tsx:307`。

## 10. 推荐结论

建议采用“页面内展开状态集合 + 受控 details”的方案。

该方案能精确满足：SSE 推理增量到来时自动展开、历史会话重新打开时默认收起、用户手动收起后不被后续增量强制展开，并且无需修改后端协议或持久化模型。
