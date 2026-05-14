import { FormEvent, useEffect, useMemo, useState } from "react";
import { api, type ChatMessage, type Conversation, type Skill, type User } from "./api";

type AuthMode = "login" | "register";
type ToolEvent = { name: string; status: string; result: string };

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
    const [authForm, setAuthForm] = useState({ email: "", username: "", login: "", password: "" });
    const [skillForm, setSkillForm] = useState({ id: "", name: "", description: "", content: "", status: "draft" as Skill["status"] });

    const activeConversation = useMemo(
        () => (Array.isArray(conversations) ? conversations : []).find((item) => item.id === activeConversationId) ?? null,
        [conversations, activeConversationId],
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
        const result = await api.createConversation("新对话");
        setConversations((prev) => [result.conversation, ...prev]);
        setActiveConversationId(result.conversation.id);
        setMessages([]);
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
            onAssistant: (payload) => setMessages((prev) => [...prev, { role: "assistant", content: payload.content }]),
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

    if (loading) {
        return <div className="center">加载中...</div>;
    }

    if (!user) {
        return (
            <div className="auth-shell">
                <div className="auth-card">
                    <h1>nano_cc Chat</h1>
                    <p>登录后即可开始新对话，在需要时使用你的个人能力。</p>
                    <div className="auth-tabs">
                        <button className={authMode === "login" ? "active" : ""} onClick={() => setAuthMode("login")}>登录</button>
                        <button className={authMode === "register" ? "active" : ""} onClick={() => setAuthMode("register")}>注册</button>
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
                        <button type="submit">{authMode === "login" ? "登录" : "注册并登录"}</button>
                    </form>
                    {error && <div className="error">{error}</div>}
                </div>
            </div>
        );
    }

    return (
        <div className="app-shell">
            <aside className="sidebar">
                <div className="sidebar-header">
                    <div>
                        <strong>{user.username}</strong>
                        <div className="muted">{user.email}</div>
                    </div>
                    <button onClick={() => void handleLogout()}>退出</button>
                </div>

                <section>
                    <div className="section-title">
                        <h2>会话</h2>
                        <button onClick={() => void handleCreateConversation()}>新建</button>
                    </div>
                    <div className="list">
                        {conversations.map((conversation) => (
                            <button
                                key={conversation.id}
                                className={conversation.id === activeConversationId ? "list-item active" : "list-item"}
                                onClick={() => setActiveConversationId(conversation.id)}
                            >
                                {conversation.title}
                            </button>
                        ))}
                    </div>
                </section>

                <section>
                    <div className="section-title"><h2>我的能力</h2></div>
                    <div className="skill-list">
                        {skills.map((skill) => (
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
                </section>
            </aside>

            <main className="main-content">
                <section className="chat-panel">
                    <div className="panel-header">
                        <div>
                            <h2>{activeConversation?.title ?? "开始一段新对话"}</h2>
                            <p className="muted">像聊天一样提问，我会优先直接回答。</p>
                        </div>
                    </div>
                    <div className="messages">
                        {messages.length === 0 ? (
                            <div className="empty-state">
                                <h3>今天想聊点什么？</h3>
                                <p>你可以直接提问、请我帮你分析问题、整理思路、起草内容，或者继续讨论代码与技术话题。</p>
                            </div>
                        ) : messages.map((message, index) => (
                            <div key={`${message.role}-${index}`} className={`message ${message.role}`}>
                                <span className="message-role">{message.role === "user" ? "你" : "助手"}</span>
                                <div className="message-content">{message.content}</div>
                            </div>
                        ))}
                    </div>
                    <form onSubmit={handleSendMessage} className="composer">
                        <textarea
                            placeholder="输入任何你想聊的问题，我会尽量直接、清晰地帮助你。"
                            value={chatInput}
                            onChange={(e) => setChatInput(e.target.value)}
                            disabled={!activeConversationId || sending}
                        />
                        <div className="composer-footer">
                            <span className="muted">网页端不执行本地命令或目录操作。</span>
                            <button type="submit" disabled={!activeConversationId || sending}>{sending ? "发送中..." : "发送"}</button>
                        </div>
                    </form>
                </section>

                <section className="right-panel">
                    <div className="panel-box">
                        <h3>运行详情</h3>
                        <div className="tool-events">
                            {toolEvents.length === 0 ? <div className="muted">当前没有工具事件</div> : toolEvents.map((event, index) => (
                                <div key={`${event.name}-${index}`} className="tool-event">
                                    <strong>{event.name}</strong>
                                    <span className={`status ${event.status}`}>{event.status}</span>
                                    <pre>{event.result}</pre>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="panel-box">
                        <h3>{skillForm.id ? "编辑能力" : "创建能力"}</h3>
                        <form onSubmit={handleSkillSubmit} className="stack">
                            <input placeholder="能力名称" value={skillForm.name} onChange={(e) => setSkillForm((prev) => ({ ...prev, name: e.target.value }))} />
                            <input placeholder="描述" value={skillForm.description} onChange={(e) => setSkillForm((prev) => ({ ...prev, description: e.target.value }))} />
                            <select value={skillForm.status} onChange={(e) => setSkillForm((prev) => ({ ...prev, status: e.target.value as Skill["status"] }))}>
                                <option value="draft">draft</option>
                                <option value="enabled">enabled</option>
                                <option value="disabled">disabled</option>
                            </select>
                            <textarea placeholder="能力内容（Markdown / Prompt）" value={skillForm.content} onChange={(e) => setSkillForm((prev) => ({ ...prev, content: e.target.value }))} rows={10} />
                            <button type="submit">{skillForm.id ? "保存修改" : "创建能力"}</button>
                        </form>
                    </div>
                    {error && <div className="error">{error}</div>}
                </section>
            </main>
        </div>
    );
}
