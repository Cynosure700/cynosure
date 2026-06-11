import { FormEvent, MouseEvent, useEffect, useMemo, useRef, useState } from "react";
import { api, type ChatMessage, type Conversation, type MCPServer, type MCPTransport, type Skill, type User } from "./api";

type AuthMode = "login" | "register";
type SidePanel = "capabilities" | "mcp" | null;

type MCPForm = {
    id: string;
    name: string;
    transport: MCPTransport;
    command: string;
    args: string;
    env: string;
    url: string;
    headers: string;
    enabled: boolean;
};

const emptyMCPForm: MCPForm = { id: "", name: "", transport: "sse", command: "", args: "", env: "", url: "", headers: "", enabled: true };

type ContentBlock =
    | { type: "paragraph"; lines: string[] }
    | { type: "list"; ordered: boolean; items: string[] }
    | { type: "code"; language: string; content: string }
    | { type: "table"; header: string[]; rows: string[][] };

function isTableRow(line: string): boolean {
    return /^\s*\|.*\|\s*$/.test(line);
}

function isTableSeparator(line: string): boolean {
    return /^\s*\|?\s*:?-{1,}:?\s*(\|\s*:?-{1,}:?\s*)+\|?\s*$/.test(line);
}

function parseTableRow(line: string): string[] {
    let trimmed = line.trim();
    if (trimmed.startsWith("|")) trimmed = trimmed.slice(1);
    if (trimmed.endsWith("|")) trimmed = trimmed.slice(0, -1);
    return trimmed.split("|").map((cell) => cell.trim());
}

function renderInlineContent(text: string) {
    const segments = text.split(/(`[^`]+`|\*\*[^*]+\*\*)/g);
    return segments.map((segment, index) => {
        if (segment.startsWith("`") && segment.endsWith("`") && segment.length >= 2) {
            return <code key={`inline-${index}`}>{segment.slice(1, -1)}</code>;
        }
        if (segment.startsWith("**") && segment.endsWith("**") && segment.length >= 4) {
            return <strong key={`inline-${index}`}>{segment.slice(2, -2)}</strong>;
        }
        return segment;
    });
}

function parseMessageContent(content: string): ContentBlock[] {
    const normalized = content.replace(/\r\n/g, "\n");
    const lines = normalized.split("\n");
    const blocks: ContentBlock[] = [];
    let paragraphLines: string[] = [];
    let listItems: string[] = [];
    let listOrdered = false;
    let inCodeBlock = false;
    let codeLanguage = "";
    let codeLines: string[] = [];

    const flushParagraph = () => {
        if (paragraphLines.length > 0) {
            blocks.push({ type: "paragraph", lines: [...paragraphLines] });
            paragraphLines = [];
        }
    };

    const flushList = () => {
        if (listItems.length > 0) {
            blocks.push({ type: "list", ordered: listOrdered, items: [...listItems] });
            listItems = [];
        }
    };

    const flushCodeBlock = () => {
        blocks.push({ type: "code", language: codeLanguage, content: codeLines.join("\n") });
        codeLines = [];
        codeLanguage = "";
    };

    for (let index = 0; index < lines.length; index += 1) {
        const line = lines[index];
        const codeFence = line.match(/^```\s*(.*)$/);
        if (codeFence) {
            flushParagraph();
            flushList();
            if (inCodeBlock) {
                flushCodeBlock();
                inCodeBlock = false;
            } else {
                inCodeBlock = true;
                codeLanguage = codeFence[1]?.trim() ?? "";
            }
            continue;
        }

        if (inCodeBlock) {
            codeLines.push(line);
            continue;
        }

        const nextLine = lines[index + 1];
        if (isTableRow(line) && nextLine !== undefined && isTableSeparator(nextLine)) {
            flushParagraph();
            flushList();
            const header = parseTableRow(line);
            const rows: string[][] = [];
            index += 2;
            while (index < lines.length && isTableRow(lines[index]) && !isTableSeparator(lines[index])) {
                rows.push(parseTableRow(lines[index]));
                index += 1;
            }
            index -= 1;
            blocks.push({ type: "table", header, rows });
            continue;
        }

        const orderedMatch = line.match(/^\s*\d+\.\s+(.*)$/);
        const unorderedMatch = line.match(/^\s*[-*+]\s+(.*)$/);

        if (orderedMatch || unorderedMatch) {
            flushParagraph();
            const ordered = Boolean(orderedMatch);
            if (listItems.length > 0 && listOrdered !== ordered) {
                flushList();
            }
            listOrdered = ordered;
            listItems.push((orderedMatch ?? unorderedMatch)?.[1]?.trim() ?? "");
            continue;
        }

        if (line.trim() === "") {
            flushParagraph();
            flushList();
            continue;
        }

        if (listItems.length > 0) {
            flushList();
        }
        paragraphLines.push(line);
    }

    if (inCodeBlock) {
        flushCodeBlock();
    }
    flushParagraph();
    flushList();
    return blocks;
}

function renderAssistantContent(content: string) {
    return parseMessageContent(content).map((block, index) => {
        if (block.type === "code") {
            return (
                <div key={`code-${index}`} className="message-code-block">
                    {block.language ? <span className="message-code-language">{block.language}</span> : null}
                    <pre>
                        <code>{block.content}</code>
                    </pre>
                </div>
            );
        }

        if (block.type === "list") {
            const ListTag = block.ordered ? "ol" : "ul";
            return (
                <ListTag key={`list-${index}`} className="message-list">
                    {block.items.map((item, itemIndex) => (
                        <li key={`item-${index}-${itemIndex}`}>{renderInlineContent(item)}</li>
                    ))}
                </ListTag>
            );
        }

        if (block.type === "table") {
            return (
                <div key={`table-${index}`} className="message-table-wrapper">
                    <table className="message-table">
                        <thead>
                            <tr>
                                {block.header.map((cell, cellIndex) => (
                                    <th key={`th-${index}-${cellIndex}`}>{renderInlineContent(cell)}</th>
                                ))}
                            </tr>
                        </thead>
                        <tbody>
                            {block.rows.map((row, rowIndex) => (
                                <tr key={`tr-${index}-${rowIndex}`}>
                                    {block.header.map((_, cellIndex) => (
                                        <td key={`td-${index}-${rowIndex}-${cellIndex}`}>{renderInlineContent(row[cellIndex] ?? "")}</td>
                                    ))}
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            );
        }

        return (
            <p key={`paragraph-${index}`} className="message-paragraph">
                {block.lines.map((line, lineIndex) => (
                    <span key={`line-${index}-${lineIndex}`}>
                        {lineIndex > 0 ? <br /> : null}
                        {renderInlineContent(line)}
                    </span>
                ))}
            </p>
        );
    });
}

function formatElapsed(ms?: number): string | null {
    if (typeof ms !== "number" || ms < 0) return null;
    if (ms < 1000) return `${ms} ms`;
    return `${(ms / 1000).toFixed(1)} s`;
}

function formatTokens(tokens?: number): string | null {
    if (typeof tokens !== "number" || tokens < 0) return null;
    if (tokens < 1024) return `${tokens}`;
    return `${(tokens / 1024).toFixed(1)}K`;
}

function formatContextLabel(message?: ChatMessage | null): string | null {
    if (!message) return null;
    const tokens = formatTokens(message.context_tokens);
    if (!tokens) return null;
    const budget = message.context_budget;
    if (typeof budget === "number" && budget > 0 && typeof message.context_tokens === "number") {
        const percent = Math.round((message.context_tokens / budget) * 100);
        return `上下文 ${tokens} / ${formatTokens(budget)} (${percent}%)`;
    }
    return `上下文 ${tokens} tokens`;
}

function renderMessageMeta(message: ChatMessage) {
    const elapsed = formatElapsed(message.elapsed_ms);
    const hasTools = typeof message.tool_call_count === "number";
    if (!elapsed && !hasTools) return null;

    return (
        <div className="message-meta">
            {elapsed && <span>耗时 {elapsed}</span>}
            {hasTools && <span>工具 {message.tool_call_count} 次</span>}
        </div>
    );
}

export function App() {
    const [authMode, setAuthMode] = useState<AuthMode>("login");
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [conversations, setConversations] = useState<Conversation[]>([]);
    const [activeConversationId, setActiveConversationId] = useState<string>("");
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [skills, setSkills] = useState<Skill[]>([]);
    const [mcpServers, setMcpServers] = useState<MCPServer[]>([]);
    const [mcpForm, setMcpForm] = useState<MCPForm>(emptyMCPForm);
    const [mcpTestResult, setMcpTestResult] = useState<{ id: string; ok: boolean; message: string } | null>(null);
    const [chatInput, setChatInput] = useState("");
    const [sending, setSending] = useState(false);
    const [updatingMemory, setUpdatingMemory] = useState(false);
    const [creatingConversation, setCreatingConversation] = useState(false);
    const [newConversationTitle, setNewConversationTitle] = useState("");
    const [renamingConversationId, setRenamingConversationId] = useState("");
    const [renameTitle, setRenameTitle] = useState("");
    const [savingRename, setSavingRename] = useState(false);
    const [openConversationMenu, setOpenConversationMenu] = useState<{ conversationId: string; left: number; top: number } | null>(null);
    const [deletingConversationId, setDeletingConversationId] = useState("");
    const [sidePanel, setSidePanel] = useState<SidePanel>(null);
    const [authForm, setAuthForm] = useState({ email: "", username: "", login: "", password: "" });
    const [skillForm, setSkillForm] = useState({ id: "", name: "", description: "", content: "", status: "draft" as Skill["status"] });
    const [expandedReasoningIds, setExpandedReasoningIds] = useState<Set<string>>(new Set());
    const [manualClosedReasoningIds, setManualClosedReasoningIds] = useState<Set<string>>(new Set());
    const activeConversationIdRef = useRef("");
    const expandedReasoningIdsRef = useRef<Set<string>>(new Set());
    const manualClosedReasoningIdsRef = useRef<Set<string>>(new Set());
    const elapsedTimerRef = useRef<number | null>(null);

    const activeConversation = useMemo(
        () => (Array.isArray(conversations) ? conversations : []).find((item) => item.id === activeConversationId) ?? null,
        [conversations, activeConversationId],
    );

    const menuConversation = useMemo(
        () => (Array.isArray(conversations) ? conversations : []).find((item) => item.id === openConversationMenu?.conversationId) ?? null,
        [conversations, openConversationMenu],
    );

    const enabledSkillCount = useMemo(
        () => skills.filter((skill) => skill.status === "enabled").length,
        [skills],
    );

    const conversationContextLabel = useMemo(() => {
        for (let index = messages.length - 1; index >= 0; index -= 1) {
            const message = messages[index];
            if (message.role === "assistant" && typeof message.context_tokens === "number") {
                return formatContextLabel(message);
            }
        }
        return null;
    }, [messages]);

    useEffect(() => {
        return () => {
            if (elapsedTimerRef.current !== null) {
                window.clearInterval(elapsedTimerRef.current);
                elapsedTimerRef.current = null;
            }
        };
    }, []);

    useEffect(() => {
        api.me()
            .then((result) => {
                setUser(result.user);
            })
            .catch(() => undefined)
            .finally(() => setLoading(false));
    }, []);

    useEffect(() => {
        if (!user) return;
        void refreshAll();
    }, [user]);

    useEffect(() => {
        activeConversationIdRef.current = activeConversationId;
        resetReasoningState();
        if (!activeConversationId) {
            setMessages([]);
            return;
        }
        const requestConversationId = activeConversationId;
        setMessages([]);
        void api
            .getConversation(requestConversationId)
            .then((result) => {
                if (activeConversationIdRef.current !== requestConversationId) return;
                setMessages(Array.isArray(result.messages) ? result.messages : []);
            })
            .catch((err) => {
                if (activeConversationIdRef.current !== requestConversationId) return;
                setError(err.message);
            });
    }, [activeConversationId]);

    useEffect(() => {
        const title = activeConversation?.title?.trim();
        document.title = title ? `${title} · Link Agent` : user ? `Link Agent · ${user.username}` : "Link Agent";
    }, [activeConversation, user]);

    useEffect(() => {
        if (!openConversationMenu) return;

        const closeMenu = () => setOpenConversationMenu(null);
        window.addEventListener("click", closeMenu);
        window.addEventListener("resize", closeMenu);
        window.addEventListener("scroll", closeMenu, true);
        return () => {
            window.removeEventListener("click", closeMenu);
            window.removeEventListener("resize", closeMenu);
            window.removeEventListener("scroll", closeMenu, true);
        };
    }, [openConversationMenu]);

    function updateReasoningSet(
        ref: { current: Set<string> },
        setter: (value: Set<string>) => void,
        updater: (current: Set<string>) => Set<string>,
    ) {
        const next = updater(ref.current);
        if (next === ref.current) return;
        ref.current = next;
        setter(next);
    }

    function addReasoningSetValue(ref: { current: Set<string> }, setter: (value: Set<string>) => void, id: string) {
        updateReasoningSet(ref, setter, (current) => {
            if (current.has(id)) return current;
            const next = new Set(current);
            next.add(id);
            return next;
        });
    }

    function removeReasoningSetValue(ref: { current: Set<string> }, setter: (value: Set<string>) => void, id: string) {
        updateReasoningSet(ref, setter, (current) => {
            if (!current.has(id)) return current;
            const next = new Set(current);
            next.delete(id);
            return next;
        });
    }

    function migrateReasoningSetValue(ref: { current: Set<string> }, setter: (value: Set<string>) => void, fromId: string, toId?: string) {
        if (!toId || fromId === toId) return;
        updateReasoningSet(ref, setter, (current) => {
            if (!current.has(fromId)) return current;
            const next = new Set(current);
            next.delete(fromId);
            next.add(toId);
            return next;
        });
    }

    function resetReasoningState() {
        const nextExpanded = new Set<string>();
        const nextManualClosed = new Set<string>();
        expandedReasoningIdsRef.current = nextExpanded;
        manualClosedReasoningIdsRef.current = nextManualClosed;
        setExpandedReasoningIds(nextExpanded);
        setManualClosedReasoningIds(nextManualClosed);
    }

    function handleReasoningToggle(reasoningKey: string, open: boolean) {
        if (open) {
            addReasoningSetValue(expandedReasoningIdsRef, setExpandedReasoningIds, reasoningKey);
            removeReasoningSetValue(manualClosedReasoningIdsRef, setManualClosedReasoningIds, reasoningKey);
            return;
        }
        removeReasoningSetValue(expandedReasoningIdsRef, setExpandedReasoningIds, reasoningKey);
        addReasoningSetValue(manualClosedReasoningIdsRef, setManualClosedReasoningIds, reasoningKey);
    }

    async function refreshAll() {
        const [conversationResult, skillResult, mcpResult] = await Promise.all([api.listConversations(), api.listSkills(), api.listMCPServers()]);
        const nextConversations = Array.isArray(conversationResult.conversations) ? conversationResult.conversations : [];
        const nextSkills = Array.isArray(skillResult.skills) ? skillResult.skills : [];

        setConversations(nextConversations);
        setSkills(nextSkills);
        setMcpServers(Array.isArray(mcpResult.mcp_servers) ? mcpResult.mcp_servers : []);
        if (!activeConversationId && nextConversations.length > 0) {
            setActiveConversationId(nextConversations[0].id);
        }
    }

    async function handleAuthSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault();
        setError("");
        try {
            const result =
                authMode === "login"
                    ? await api.login({ login: authForm.login, password: authForm.password })
                    : await api.register({ email: authForm.email, username: authForm.username, password: authForm.password });
            setUser(result.user);
        } catch (err) {
            setError(err instanceof Error ? err.message : "认证失败");
        }
    }

    async function handleLogout() {
        await api.logout();
        setUser(null);
        setConversations([]);
        setSkills([]);
        setMessages([]);
        setActiveConversationId("");
    }

    async function handleToggleMemory() {
        if (!user || updatingMemory) return;
        const nextEnabled = !user.memory_enabled;
        setUpdatingMemory(true);
        setError("");
        setUser((prev) => (prev ? { ...prev, memory_enabled: nextEnabled } : prev));
        try {
            await api.updateMemoryPreference(nextEnabled);
        } catch (err) {
            setUser((prev) => (prev ? { ...prev, memory_enabled: !nextEnabled } : prev));
            setError(err instanceof Error ? err.message : "更新记忆设置失败");
        } finally {
            setUpdatingMemory(false);
        }
    }

    async function handleCreateConversation() {
        setCreatingConversation(true);
        setError("");
        try {
            const result = await api.createConversation(newConversationTitle);
            setConversations((prev) => [result.conversation, ...prev]);
            setActiveConversationId(result.conversation.id);
            setMessages([]);
            setChatInput("");
            setNewConversationTitle("");
        } catch (err) {
            setError(err instanceof Error ? err.message : "创建会话失败");
        } finally {
            setCreatingConversation(false);
        }
    }

    async function handleSendMessage(event: FormEvent<HTMLFormElement>) {
        event.preventDefault();
        if (!activeConversationId || !chatInput.trim()) return;
        setSending(true);
        setError("");
        const content = chatInput.trim();
        const requestConversationId = activeConversationId;
        const streamingAssistantId = `streaming-${Date.now()}`;
        setChatInput("");
        setMessages((prev) => [...prev, { role: "user", content }, { id: streamingAssistantId, role: "assistant", content: "", reasoning_content: "" }]);

        const startedAt = Date.now();
        const stopElapsedTimer = () => {
            if (elapsedTimerRef.current !== null) {
                window.clearInterval(elapsedTimerRef.current);
                elapsedTimerRef.current = null;
            }
        };
        const clearStreamingElapsed = () => {
            stopElapsedTimer();
            setMessages((prev) => prev.map((message) => (
                message.id === streamingAssistantId ? { ...message, elapsed_ms: undefined } : message
            )));
        };
        stopElapsedTimer();
        elapsedTimerRef.current = window.setInterval(() => {
            if (activeConversationIdRef.current !== requestConversationId) return;
            setMessages((prev) => prev.map((message) => (
                message.id === streamingAssistantId ? { ...message, elapsed_ms: Date.now() - startedAt } : message
            )));
        }, 200);

        try {
            await api.streamConversation(requestConversationId, content, {
                onAssistantDelta: (payload) => setMessages((prev) => prev.map((message) => (
                    activeConversationIdRef.current === requestConversationId && message.id === streamingAssistantId ? { ...message, content: message.content + payload.content } : message
                ))),
                onReasoningDelta: (payload) => {
                    if (activeConversationIdRef.current !== requestConversationId) return;
                    if (payload.content && !manualClosedReasoningIdsRef.current.has(streamingAssistantId)) {
                        addReasoningSetValue(expandedReasoningIdsRef, setExpandedReasoningIds, streamingAssistantId);
                    }
                    setMessages((prev) => prev.map((message) => (
                        message.id === streamingAssistantId ? { ...message, reasoning_content: `${message.reasoning_content ?? ""}${payload.content}` } : message
                    )));
                },
                onMeta: (payload) => {
                    if (activeConversationIdRef.current !== requestConversationId) return;
                    setMessages((prev) => prev.map((message) => (
                        message.id === streamingAssistantId
                            ? { ...message, tool_call_count: payload.tool_call_count, context_tokens: payload.context_tokens, context_budget: payload.context_budget }
                            : message
                    )));
                },
                onAssistant: (payload) => {
                    stopElapsedTimer();
                    setMessages((prev) => {
                        if (activeConversationIdRef.current !== requestConversationId) return prev;
                        migrateReasoningSetValue(expandedReasoningIdsRef, setExpandedReasoningIds, streamingAssistantId, payload.message_id);
                        migrateReasoningSetValue(manualClosedReasoningIdsRef, setManualClosedReasoningIds, streamingAssistantId, payload.message_id);
                        const hasStreamingMessage = prev.some((message) => message.id === streamingAssistantId);
                        if (!hasStreamingMessage) {
                            return [...prev, { id: payload.message_id, role: "assistant", content: payload.content, reasoning_content: payload.reasoning_content, tool_call_count: payload.tool_call_count, context_tokens: payload.context_tokens, context_budget: payload.context_budget }];
                        }
                        return prev.map((message) => (
                            message.id === streamingAssistantId
                                ? { ...message, id: payload.message_id ?? message.id, content: payload.content, reasoning_content: payload.reasoning_content ?? message.reasoning_content, elapsed_ms: undefined, tool_call_count: payload.tool_call_count, context_tokens: payload.context_tokens, context_budget: payload.context_budget }
                                : message
                        ));
                    });
                },
                onError: (message) => {
                    clearStreamingElapsed();
                    if (activeConversationIdRef.current !== requestConversationId) {
                        setSending(false);
                        return;
                    }
                    setError(message);
                    setMessages((prev) => prev.filter((item) => item.id !== streamingAssistantId || item.content.trim() || item.reasoning_content?.trim()));
                    setSending(false);
                },
                onDone: () => {
                    stopElapsedTimer();
                    setSending(false);
                },
            });
            await refreshAll();
        } catch (err) {
            clearStreamingElapsed();
            if (activeConversationIdRef.current !== requestConversationId) {
                setSending(false);
                return;
            }
            setError(err instanceof Error ? err.message : "发送消息失败");
            setMessages((prev) => prev.filter((item) => item.id !== streamingAssistantId || item.content.trim() || item.reasoning_content?.trim()));
            setSending(false);
        }
    }

    async function handleSkillSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault();
        setError("");
        try {
            if (skillForm.id) {
                await api.updateSkill(skillForm.id, skillForm);
            } else {
                await api.createSkill(skillForm);
            }
            setSkillForm({ id: "", name: "", description: "", content: "", status: "draft" });
            const result = await api.listSkills();
            setSkills(Array.isArray(result.skills) ? result.skills : []);
        } catch (err) {
            setError(err instanceof Error ? err.message : "保存 skill 失败");
        }
    }

    async function toggleSkillStatus(skill: Skill) {
        const nextStatus: Skill["status"] = skill.status === "enabled" ? "disabled" : "enabled";
        const result = await api.patchSkillStatus(skill.id, nextStatus);
        setSkills((prev) => prev.map((item) => (item.id === skill.id ? result.skill : item)));
    }

    async function deleteSkill(skillId: string) {
        await api.deleteSkill(skillId);
        setSkills((prev) => prev.filter((item) => item.id !== skillId));
    }

    function parseLines(value: string): string[] {
        return value.split("\n").map((line) => line.trim()).filter((line) => line.length > 0);
    }

    function parseKeyValues(value: string): Record<string, string> {
        const result: Record<string, string> = {};
        for (const line of parseLines(value)) {
            const idx = line.indexOf("=");
            if (idx <= 0) continue;
            result[line.slice(0, idx).trim()] = line.slice(idx + 1).trim();
        }
        return result;
    }

    function formatKeyValues(record: Record<string, string> | null | undefined): string {
        if (!record) return "";
        return Object.entries(record).map(([k, v]) => `${k}=${v}`).join("\n");
    }

    async function handleMCPSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault();
        setError("");
        const payload = {
            name: mcpForm.name,
            transport: mcpForm.transport,
            command: mcpForm.command,
            args: parseLines(mcpForm.args),
            env: parseKeyValues(mcpForm.env),
            url: mcpForm.url,
            headers: parseKeyValues(mcpForm.headers),
            enabled: mcpForm.enabled,
        };
        try {
            if (mcpForm.id) {
                await api.updateMCPServer(mcpForm.id, payload);
            } else {
                await api.createMCPServer(payload);
            }
            setMcpForm(emptyMCPForm);
            const result = await api.listMCPServers();
            setMcpServers(Array.isArray(result.mcp_servers) ? result.mcp_servers : []);
        } catch (err) {
            setError(err instanceof Error ? err.message : "保存 MCP 服务器失败");
        }
    }

    function editMCPServer(server: MCPServer) {
        setMcpForm({
            id: server.id,
            name: server.name,
            transport: server.transport,
            command: server.command,
            args: (server.args ?? []).join("\n"),
            env: formatKeyValues(server.env),
            url: server.url,
            headers: formatKeyValues(server.headers),
            enabled: server.enabled,
        });
        setMcpTestResult(null);
    }

    async function toggleMCPServer(server: MCPServer) {
        const result = await api.patchMCPServerEnabled(server.id, !server.enabled);
        setMcpServers((prev) => prev.map((item) => (item.id === server.id ? result.mcp_server : item)));
    }

    async function deleteMCPServer(id: string) {
        await api.deleteMCPServer(id);
        setMcpServers((prev) => prev.filter((item) => item.id !== id));
        if (mcpForm.id === id) setMcpForm(emptyMCPForm);
    }

    async function testMCPServer(id: string) {
        setMcpTestResult(null);
        try {
            const result = await api.testMCPServer(id);
            setMcpTestResult({
                id,
                ok: result.ok,
                message: result.ok ? `发现 ${result.tools?.length ?? 0} 个工具：${(result.tools ?? []).join(", ")}` : result.error ?? "连接失败",
            });
        } catch (err) {
            setMcpTestResult({ id, ok: false, message: err instanceof Error ? err.message : "测试失败" });
        }
    }

    async function handleDeleteConversation(conversationId: string) {
        setOpenConversationMenu(null);
        setDeletingConversationId(conversationId);
        setError("");
        try {
            await api.deleteConversation(conversationId);
            const nextConversations = conversations.filter((item) => item.id !== conversationId);
            setConversations(nextConversations);
            if (conversationId === activeConversationId) {
                setActiveConversationId(nextConversations[0]?.id ?? "");
                setMessages([]);
                setChatInput("");
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : "删除会话失败");
        } finally {
            setDeletingConversationId("");
        }
    }

    function startRenameConversation(conversation: Conversation) {
        setError("");
        setOpenConversationMenu(null);
        setRenamingConversationId(conversation.id);
        setRenameTitle(conversation.title);
    }

    function toggleConversationMenu(event: MouseEvent<HTMLButtonElement>, conversationId: string) {
        event.stopPropagation();
        const rect = event.currentTarget.getBoundingClientRect();
        const menuHeight = 96;
        const gap = 8;
        setOpenConversationMenu((current) =>
            current?.conversationId === conversationId
                ? null
                : {
                    conversationId,
                    left: rect.right + gap,
                    top: rect.bottom + gap + menuHeight > window.innerHeight ? Math.max(gap, rect.top - menuHeight - gap) : rect.bottom + gap,
                },
        );
    }

    function cancelRenameConversation() {
        setRenamingConversationId("");
        setRenameTitle("");
        setSavingRename(false);
    }

    async function handleRenameConversation(event: FormEvent<HTMLFormElement>, conversationId: string) {
        event.preventDefault();
        setSavingRename(true);
        setError("");
        try {
            const result = await api.renameConversation(conversationId, renameTitle);
            setConversations((prev) => prev.map((item) => (item.id === conversationId ? result.conversation : item)));
            cancelRenameConversation();
        } catch (err) {
            setError(err instanceof Error ? err.message : "重命名会话失败");
            setSavingRename(false);
        }
    }

    if (loading) {
        return <div className="center">加载中...</div>;
    }

    if (!user) {
        return (
            <div className="auth-shell">
                <div className="auth-visual" aria-hidden="true">
                    <div className="auth-visual-card">
                        <div className="visual-glass-panel">
                            <span className="eyebrow">Link Agent</span>
                            <strong>Reason<span className="visual-dot">·</span>Act<span className="visual-dot">·</span>Remember</strong>
                            <p>统一编排多轮对话、Skill 能力与受控工具调用，构建可审计、可扩展的智能协作入口。</p>
                        </div>
                        <ul className="visual-features">
                            <li>
                                <span className="visual-feature-tag">Reason</span>
                                <div>
                                    <strong>多轮推理</strong>
                                    <span>顺着上下文逐步拆解问题，给出可追溯的思考过程。</span>
                                </div>
                            </li>
                            <li>
                                <span className="visual-feature-tag">Act</span>
                                <div>
                                    <strong>受控执行</strong>
                                    <span>编排 Skill 与 MCP 工具调用，每一步都在可审计的边界内。</span>
                                </div>
                            </li>
                            <li>
                                <span className="visual-feature-tag">Remember</span>
                                <div>
                                    <strong>沉淀记忆</strong>
                                    <span>自动记住关键上下文，让下一次对话从已有进度继续。</span>
                                </div>
                            </li>
                        </ul>
                        <div className="visual-footnote">
                            <span>可审计</span>
                            <span>可扩展</span>
                            <span>团队协作</span>
                        </div>
                    </div>
                </div>
                <div className="auth-card">
                    <div className="auth-copy">
                        <span className="eyebrow auth-product-name">Link Agent</span>
                        <h1>欢迎回来</h1>
                        <p>登录后即可继续对话、分析问题、沉淀上下文，并通过能力编排提升复杂任务处理效率。</p>
                    </div>
                    <div className="auth-tabs">
                        <button className={authMode === "login" ? "active" : ""} onClick={() => setAuthMode("login")}>继续对话</button>
                        <button className={authMode === "register" ? "active" : ""} onClick={() => setAuthMode("register")}>创建账号</button>
                    </div>
                    <form onSubmit={handleAuthSubmit} className="stack">
                        {authMode === "register" && (
                            <>
                                <input placeholder="邮箱" value={authForm.email} onChange={(e) => setAuthForm((prev) => ({ ...prev, email: e.target.value }))} />
                                <input placeholder="用户名" value={authForm.username} onChange={(e) => setAuthForm((prev) => ({ ...prev, username: e.target.value }))} />
                            </>
                        )}
                        {authMode === "login" && (
                            <input placeholder="邮箱或用户名" value={authForm.login} onChange={(e) => setAuthForm((prev) => ({ ...prev, login: e.target.value }))} />
                        )}
                        <input type="password" placeholder="密码" value={authForm.password} onChange={(e) => setAuthForm((prev) => ({ ...prev, password: e.target.value }))} />
                        <button type="submit">{authMode === "login" ? "进入聊天" : "创建账号并进入"}</button>
                    </form>
                    {error && <div className="error">{error}</div>}
                </div>
            </div>
        );
    }

    return (
        <div className="app-shell">
            <aside className="sidebar">
                <div className="sidebar-top">
                    <div className="sidebar-brand">
                        <div className="brand-row">
                            <span className="brand-mark">L</span>
                            <div>
                                <strong className="brand-title">Link Agent</strong>
                            </div>
                        </div>
                    </div>
                    <div className="sidebar-header">
                        <div>
                            <strong>{user.username}</strong>
                            <div className="muted">{user.email}</div>
                        </div>
                        <button onClick={() => void handleLogout()}>退出</button>
                    </div>
                    <div className="memory-toggle-row">
                        <div className="memory-toggle-copy">
                            <strong>记忆功能</strong>
                            <span className="muted">{user.memory_enabled ? "已开启，自动记住并复用上下文" : "已关闭，不注入也不提取记忆"}</span>
                        </div>
                        <button
                            type="button"
                            className={user.memory_enabled ? "memory-switch on" : "memory-switch"}
                            role="switch"
                            aria-checked={user.memory_enabled}
                            aria-label="切换记忆功能"
                            disabled={updatingMemory}
                            onClick={() => void handleToggleMemory()}
                        >
                            <span className="memory-switch-thumb" />
                        </button>
                    </div>
                    <div className="sidebar-insights">
                        <div>
                            <strong>{conversations.length}</strong>
                            <span>会话</span>
                        </div>
                        <div>
                            <strong>{enabledSkillCount}</strong>
                            <span>启用能力</span>
                        </div>
                    </div>
                    <div className="current-session-card">
                        <span className="eyebrow">current thread</span>
                        <strong>{activeConversation?.title ?? "等待选择会话"}</strong>
                        <p>{messages.length > 0 ? `${messages.length} 条上下文消息已载入` : "创建或选择会话后开始对话"}</p>
                    </div>
                </div>

                <section className="sidebar-section conversation-section">
                    <div className="section-title conversation-title-row">
                        <div>
                            <h2>会话</h2>
                            <div className="section-subtitle">最近的聊天会出现在这里</div>
                        </div>
                        <button type="submit" form="new-conversation-form" className="conversation-add-button" disabled={creatingConversation}>
                            {creatingConversation ? "..." : "+"}
                        </button>
                    </div>
                    <form
                        id="new-conversation-form"
                        className="new-conversation-form"
                        onSubmit={(event) => {
                            event.preventDefault();
                            void handleCreateConversation();
                        }}
                    >
                        <input
                            placeholder="新会话名称（可选）"
                            value={newConversationTitle}
                            onChange={(event) => setNewConversationTitle(event.target.value)}
                            disabled={creatingConversation}
                        />
                    </form>
                    <div className="list">
                        {conversations.map((conversation) => {
                            const isRenaming = renamingConversationId === conversation.id;

                            return (
                                <div key={conversation.id} className="list-entry">
                                    <div className={isRenaming ? "list-row editing" : "list-row"}>
                                        {isRenaming ? (
                                            <form className="inline-rename-form" onSubmit={(event) => void handleRenameConversation(event, conversation.id)}>
                                                <input
                                                    value={renameTitle}
                                                    onChange={(event) => setRenameTitle(event.target.value)}
                                                    onKeyDown={(event) => {
                                                        if (event.key === "Escape") cancelRenameConversation();
                                                    }}
                                                    placeholder="输入新的会话名称"
                                                    disabled={savingRename}
                                                    autoFocus
                                                />
                                                <div className="inline-rename-actions">
                                                    <button type="submit" disabled={savingRename}>{savingRename ? "保存中..." : "保存"}</button>
                                                    <button type="button" className="secondary-toggle" onClick={cancelRenameConversation} disabled={savingRename}>取消</button>
                                                </div>
                                            </form>
                                        ) : (
                                            <>
                                                <button
                                                    className={conversation.id === activeConversationId ? "list-item active" : "list-item"}
                                                    onClick={() => setActiveConversationId(conversation.id)}
                                                >
                                                    <span className="conversation-copy">
                                                        <strong>{conversation.title}</strong>
                                                    </span>
                                                </button>
                                                <div className="conversation-menu">
                                                    <button
                                                        className="conversation-menu-trigger"
                                                        type="button"
                                                        aria-label={`打开 ${conversation.title} 的操作菜单`}
                                                        onClick={(event) => toggleConversationMenu(event, conversation.id)}
                                                    >
                                                        •••
                                                    </button>
                                                </div>
                                            </>
                                        )}
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                </section>
            </aside>

            <main className={sidePanel ? "main-content with-side-panel" : "main-content"}>
                <section className="chat-panel">
                    <div className="panel-header">
                        <div className="panel-header-copy">
                            <span className="eyebrow">assistant-first chat</span>
                            <h2>{activeConversation?.title ?? "开始一段新对话"}</h2>
                            <p className="muted">{conversationContextLabel ?? "当前会话暂无上下文信息"}</p>
                        </div>
                        <div className="agent-orbit" aria-hidden="true">
                            <span />
                            <span />
                            <span />
                        </div>
                        <div className="panel-actions">
                            <button className={sidePanel === "capabilities" ? "secondary-toggle active" : "secondary-toggle"} onClick={() => setSidePanel((current) => current === "capabilities" ? null : "capabilities")}>能力</button>
                            <button className={sidePanel === "mcp" ? "secondary-toggle active" : "secondary-toggle"} onClick={() => setSidePanel((current) => current === "mcp" ? null : "mcp")}>MCP</button>
                        </div>
                    </div>
                    <div className="messages">
                        {messages.length === 0 ? (
                            <div className="empty-state">
                                <div className="empty-copy">
                                    <span className="eyebrow">ready when you are</span>
                                    <h3>今天想聊点什么？</h3>
                                    <p>你可以直接提问、让我帮你分析问题、整理思路、起草内容，或者继续讨论代码与技术话题。</p>
                                    <div className="prompt-chips">
                                        <span>分析需求</span>
                                        <span>生成方案</span>
                                        <span>整理文档</span>
                                        <span>排查代码</span>
                                    </div>
                                </div>
                                <div className="empty-visual" aria-hidden="true">
                                    <div className="neural-card">
                                        <span className="pulse-dot" />
                                        <strong>Link Agent</strong>
                                        <small>ready to help</small>
                                    </div>
                                </div>
                            </div>
                        ) : messages.map((message, index) => {
                            const reasoningKey = message.id ?? `${message.role}-${index}`;
                            return (
                                <div key={`${message.role}-${index}`} className={`message ${message.role}`}>
                                    <div className="message-header">
                                        <span className="message-role">{message.role === "user" ? "你" : "助手"}</span>
                                        {message.role === "assistant" && renderMessageMeta(message)}
                                    </div>
                                    {message.role === "assistant" && message.reasoning_content?.trim() && (
                                        <details
                                            className="message-reasoning"
                                            open={expandedReasoningIds.has(reasoningKey)}
                                            onToggle={(event) => handleReasoningToggle(reasoningKey, event.currentTarget.open)}
                                        >
                                            <summary>推理过程</summary>
                                            <div>{message.reasoning_content}</div>
                                        </details>
                                    )}
                                    <div className="message-content">
                                        {message.role === "assistant" ? (message.content.trim() ? renderAssistantContent(message.content) : <span className="muted">正在生成...</span>) : message.content}
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                    <form onSubmit={handleSendMessage} className="composer">
                        <div className="composer-head">
                            <span className="eyebrow">message</span>
                            <span className="muted">聚焦输入，不展示多余操作</span>
                        </div>
                        <textarea
                            placeholder="给我发消息，告诉我你现在想解决什么问题。"
                            value={chatInput}
                            onChange={(e) => setChatInput(e.target.value)}
                            onKeyDown={(event) => {
                                if (event.key === "Enter" && !event.shiftKey && !event.nativeEvent.isComposing) {
                                    event.preventDefault();
                                    event.currentTarget.form?.requestSubmit();
                                }
                            }}
                            disabled={!activeConversationId || sending}
                        />
                        <div className="composer-footer">
                            <span className="muted">输入问题，按发送开始对话。</span>
                            <button type="submit" disabled={!activeConversationId || sending}>{sending ? "思考中..." : "发送消息"}</button>
                        </div>
                    </form>
                </section>

                {sidePanel === "capabilities" && (
                    <section className="right-panel">
                        <div className="panel-box">
                            <div className="section-title">
                                <div>
                                    <h3>我的能力</h3>
                                    <div className="section-subtitle">常用能力集中管理</div>
                                </div>
                                <button className="secondary-toggle" onClick={() => setSidePanel(null)}>收起</button>
                            </div>
                            <div className="skill-list">
                                {skills.length === 0 ? <div className="muted">你还没有添加任何个人能力</div> : skills.map((skill) => (
                                    <div key={skill.id} className="skill-card">
                                        <div className="skill-head">
                                            <strong>{skill.name}</strong>
                                            <span className={`status ${skill.status}`}>{skill.status}</span>
                                        </div>
                                        <p>{skill.description || "无描述"}</p>
                                        <div className="actions">
                                            <button onClick={() => setSkillForm(skill)}>编辑</button>
                                            <button onClick={() => void toggleSkillStatus(skill)}>{skill.status === "enabled" ? "禁用" : "启用"}</button>
                                            <button onClick={() => void deleteSkill(skill.id)}>删除</button>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>

                        <div className="panel-box">
                            <div className="panel-box-header">
                                <h3>{skillForm.id ? "编辑能力" : "创建能力"}</h3>
                                <div className="section-subtitle">能力配置保持在侧边，不打断主聊天流</div>
                            </div>
                            <form onSubmit={handleSkillSubmit} className="stack">
                                <input placeholder="能力名称" value={skillForm.name} onChange={(e) => setSkillForm((prev) => ({ ...prev, name: e.target.value }))} />
                                <input placeholder="描述" value={skillForm.description} onChange={(e) => setSkillForm((prev) => ({ ...prev, description: e.target.value }))} />
                                <select value={skillForm.status} onChange={(e) => setSkillForm((prev) => ({ ...prev, status: e.target.value as Skill["status"] }))}>
                                    <option value="draft">draft</option>
                                    <option value="enabled">enabled</option>
                                    <option value="disabled">disabled</option>
                                </select>
                                <textarea placeholder="能力内容（Markdown / Prompt）" value={skillForm.content} onChange={(e) => setSkillForm((prev) => ({ ...prev, content: e.target.value }))} rows={10} />
                                <button type="submit">{skillForm.id ? "保存能力" : "创建能力"}</button>
                            </form>
                        </div>
                        {error && <div className="error">{error}</div>}
                    </section>
                )}

                {sidePanel === "mcp" && (
                    <section className="right-panel">
                        <div className="panel-box">
                            <div className="section-title">
                                <div>
                                    <h3>MCP 服务器</h3>
                                    <div className="section-subtitle">连接 MCP 服务器以扩展工具能力</div>
                                </div>
                                <button className="secondary-toggle" onClick={() => setSidePanel(null)}>收起</button>
                            </div>
                            <div className="skill-list">
                                {mcpServers.length === 0 ? <div className="muted">你还没有配置任何 MCP 服务器</div> : mcpServers.map((server) => (
                                    <div key={server.id} className="skill-card">
                                        <div className="skill-head">
                                            <strong>{server.name}</strong>
                                            <span className={`status ${server.enabled ? "enabled" : "disabled"}`}>{server.transport}</span>
                                        </div>
                                        <p>{server.url || "无地址"}</p>
                                        <div className="actions">
                                            <button onClick={() => editMCPServer(server)}>编辑</button>
                                            <button onClick={() => void toggleMCPServer(server)}>{server.enabled ? "禁用" : "启用"}</button>
                                            <button onClick={() => void testMCPServer(server.id)}>测试连接</button>
                                            <button onClick={() => void deleteMCPServer(server.id)}>删除</button>
                                        </div>
                                        {mcpTestResult?.id === server.id && (
                                            <div className={mcpTestResult.ok ? "muted" : "error"}>{mcpTestResult.message}</div>
                                        )}
                                    </div>
                                ))}
                            </div>
                        </div>

                        <div className="panel-box">
                            <div className="panel-box-header">
                                <h3>{mcpForm.id ? "编辑 MCP 服务器" : "添加 MCP 服务器"}</h3>
                                <div className="section-subtitle">支持 sse / streamable 两种连接方式</div>
                            </div>
                            <form onSubmit={handleMCPSubmit} className="stack">
                                <input placeholder="名称（用于工具命名空间）" value={mcpForm.name} onChange={(e) => setMcpForm((prev) => ({ ...prev, name: e.target.value }))} />
                                <select value={mcpForm.transport} onChange={(e) => setMcpForm((prev) => ({ ...prev, transport: e.target.value as MCPTransport }))}>
                                    <option value="sse">sse</option>
                                    <option value="streamable">streamable</option>
                                </select>
                                <input placeholder="服务地址 URL" value={mcpForm.url} onChange={(e) => setMcpForm((prev) => ({ ...prev, url: e.target.value }))} />
                                <textarea placeholder="请求头（每行 KEY=VALUE）" value={mcpForm.headers} onChange={(e) => setMcpForm((prev) => ({ ...prev, headers: e.target.value }))} rows={3} />
                                <label className="mcp-enabled-row">
                                    <input type="checkbox" checked={mcpForm.enabled} onChange={(e) => setMcpForm((prev) => ({ ...prev, enabled: e.target.checked }))} />
                                    <span>启用</span>
                                </label>
                                <button type="submit">{mcpForm.id ? "保存" : "添加"}</button>
                                {mcpForm.id && <button type="button" onClick={() => setMcpForm(emptyMCPForm)}>取消编辑</button>}
                            </form>
                        </div>
                        {error && <div className="error">{error}</div>}
                    </section>
                )}
            </main>
            {openConversationMenu && menuConversation && (
                <div
                    className="conversation-menu-panel"
                    style={{ left: openConversationMenu.left, top: openConversationMenu.top }}
                    onClick={(event) => event.stopPropagation()}
                >
                    <button
                        className="menu-action"
                        onClick={() => startRenameConversation(menuConversation)}
                        disabled={savingRename || deletingConversationId === menuConversation.id}
                    >
                        重命名会话
                    </button>
                    <button
                        className="menu-action danger"
                        onClick={() => void handleDeleteConversation(menuConversation.id)}
                        disabled={deletingConversationId === menuConversation.id || savingRename}
                    >
                        {deletingConversationId === menuConversation.id ? "删除中..." : "删除会话"}
                    </button>
                </div>
            )}
        </div>
    );
}
