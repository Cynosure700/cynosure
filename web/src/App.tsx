import { FormEvent, MouseEvent, useEffect, useMemo, useState } from "react";
import { api, type ChatMessage, type Conversation, type Skill, type User } from "./api";

type AuthMode = "login" | "register";
type ToolEvent = { name: string; status: string; result: string };
type SidePanel = "capabilities" | "details" | null;

type ContentBlock =
    | { type: "paragraph"; lines: string[] }
    | { type: "list"; ordered: boolean; items: string[] }
    | { type: "code"; language: string; content: string };

function renderInlineContent(text: string) {
    const segments = text.split(/(`[^`]+`)/g);
    return segments.map((segment, index) => {
        if (segment.startsWith("`") && segment.endsWith("`") && segment.length >= 2) {
            return <code key={`inline-${index}`}>{segment.slice(1, -1)}</code>;
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

    for (const line of lines) {
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

export function App() {
    const [authMode, setAuthMode] = useState<AuthMode>("login");
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [conversations, setConversations] = useState<Conversation[]>([]);
    const [activeConversationId, setActiveConversationId] = useState<string>("");
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [skills, setSkills] = useState<Skill[]>([]);
    const [toolEvents, setToolEvents] = useState<ToolEvent[]>([]);
    const [chatInput, setChatInput] = useState("");
    const [sending, setSending] = useState(false);
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

    const activeConversation = useMemo(
        () => (Array.isArray(conversations) ? conversations : []).find((item) => item.id === activeConversationId) ?? null,
        [conversations, activeConversationId],
    );

    const enabledSkillCount = useMemo(
        () => skills.filter((skill) => skill.status === "enabled").length,
        [skills],
    );

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
        if (!activeConversationId) return;
        void api
            .getConversation(activeConversationId)
            .then((result) => setMessages(Array.isArray(result.messages) ? result.messages : []))
            .catch((err) => setError(err.message));
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

    async function refreshAll() {
        const [conversationResult, skillResult] = await Promise.all([api.listConversations(), api.listSkills()]);
        const nextConversations = Array.isArray(conversationResult.conversations) ? conversationResult.conversations : [];
        const nextSkills = Array.isArray(skillResult.skills) ? skillResult.skills : [];

        setConversations(nextConversations);
        setSkills(nextSkills);
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

    async function handleCreateConversation() {
        setCreatingConversation(true);
        setError("");
        try {
            const result = await api.createConversation(newConversationTitle);
            setConversations((prev) => [result.conversation, ...prev]);
            setActiveConversationId(result.conversation.id);
            setMessages([]);
            setToolEvents([]);
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
        setChatInput("");
        setMessages((prev) => [...prev, { role: "user", content }]);
        setToolEvents([]);

        await api.streamConversation(activeConversationId, content, {
            onTool: (payload) => setToolEvents((prev) => [...prev, payload]),
            onAssistant: (payload) => setMessages((prev) => [...prev, { role: "assistant", content: payload.content, reasoning_content: payload.reasoning_content }]),
            onError: (message) => setError(message),
            onDone: () => setSending(false),
        });

        await refreshAll();
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
                setToolEvents([]);
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
        setOpenConversationMenu((current) =>
            current?.conversationId === conversationId
                ? null
                : {
                    conversationId,
                    left: rect.right + 8,
                    top: rect.top,
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
                            <strong>Reason · Act · Remember</strong>
                            <p>统一编排多轮对话、Skill 能力与受控工具调用，构建可审计、可扩展的智能协作入口。</p>
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
                    <div className="section-title">
                        <div>
                            <h2>会话</h2>
                            <div className="section-subtitle">最近的聊天会出现在这里</div>
                        </div>
                    </div>
                    <form
                        className="stack"
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
                        <button type="submit" disabled={creatingConversation}>{creatingConversation ? "创建中..." : "新对话"}</button>
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
                                                    {openConversationMenu?.conversationId === conversation.id && (
                                                        <div
                                                            className="conversation-menu-panel"
                                                            style={{ left: openConversationMenu.left, top: openConversationMenu.top }}
                                                            onClick={(event) => event.stopPropagation()}
                                                        >
                                                            <button
                                                                className="menu-action"
                                                                onClick={() => startRenameConversation(conversation)}
                                                                disabled={savingRename || deletingConversationId === conversation.id}
                                                            >
                                                                重命名会话
                                                            </button>
                                                            <button
                                                                className="menu-action danger"
                                                                onClick={() => void handleDeleteConversation(conversation.id)}
                                                                disabled={deletingConversationId === conversation.id || savingRename}
                                                            >
                                                                {deletingConversationId === conversation.id ? "删除中..." : "删除会话"}
                                                            </button>
                                                        </div>
                                                    )}
                                                </div>
                                            </>
                                        )}
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                </section>

                <section className="sidebar-section secondary-actions">
                    <div className="section-title">
                        <div>
                            <h2>辅助面板</h2>
                            <div className="section-subtitle">默认收起，需要时再展开</div>
                        </div>
                    </div>
                    <div className="secondary-buttons">
                        <button className={sidePanel === "capabilities" ? "secondary-toggle active" : "secondary-toggle"} onClick={() => setSidePanel((current) => current === "capabilities" ? null : "capabilities")}>能力</button>
                        <button className={sidePanel === "details" ? "secondary-toggle active" : "secondary-toggle"} onClick={() => setSidePanel((current) => current === "details" ? null : "details")}>工具调用</button>
                    </div>
                </section>
            </aside>

            <main className={sidePanel ? "main-content with-side-panel" : "main-content"}>
                <section className="chat-panel">
                    <div className="panel-header">
                        <div className="panel-header-copy">
                            <span className="eyebrow">assistant-first chat</span>
                            <h2>{activeConversation?.title ?? "开始一段新对话"}</h2>
                            <p className="muted">像与通用聊天助手对话一样提问，我会优先直接回答。</p>
                        </div>
                        <div className="agent-orbit" aria-hidden="true">
                            <span />
                            <span />
                            <span />
                        </div>
                        <div className="panel-actions">
                            <button className={sidePanel === "capabilities" ? "secondary-toggle active" : "secondary-toggle"} onClick={() => setSidePanel((current) => current === "capabilities" ? null : "capabilities")}>能力</button>
                            <button className={sidePanel === "details" ? "secondary-toggle active" : "secondary-toggle"} onClick={() => setSidePanel((current) => current === "details" ? null : "details")}>工具调用</button>
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
                        ) : messages.map((message, index) => (
                            <div key={`${message.role}-${index}`} className={`message ${message.role}`}>
                                <span className="message-role">{message.role === "user" ? "你" : "助手"}</span>
                                {message.role === "assistant" && message.reasoning_content?.trim() && (
                                    <details className="message-reasoning">
                                        <summary>推理过程</summary>
                                        <div>{message.reasoning_content}</div>
                                    </details>
                                )}
                                <div className="message-content">
                                    {message.role === "assistant" ? renderAssistantContent(message.content) : message.content}
                                </div>
                            </div>
                        ))}
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
                            disabled={!activeConversationId || sending}
                        />
                        <div className="composer-footer">
                            <span className="muted">输入问题，按发送开始对话。</span>
                            <button type="submit" disabled={!activeConversationId || sending}>{sending ? "思考中..." : "发送消息"}</button>
                        </div>
                    </form>
                </section>

                {sidePanel && (
                    <section className="right-panel">
                        {sidePanel === "details" ? (
                            <div className="panel-box">
                                <div className="section-title">
                                    <div>
                                        <h3>工具调用</h3>
                                        <div className="section-subtitle">仅在需要排查时查看</div>
                                    </div>
                                    <button className="secondary-toggle" onClick={() => setSidePanel(null)}>收起</button>
                                </div>
                                <div className="tool-events">
                                    {toolEvents.length === 0 ? <div className="muted">这轮对话还没有工具调用</div> : toolEvents.map((event, index) => (
                                        <div key={`${event.name}-${index}`} className="tool-event">
                                            <strong>{event.name}</strong>
                                            <span className={`status ${event.status}`}>{event.status}</span>
                                            <pre>{event.result}</pre>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        ) : (
                            <>
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
                            </>
                        )}
                        {error && <div className="error">{error}</div>}
                    </section>
                )}
            </main>
        </div>
    );
}
