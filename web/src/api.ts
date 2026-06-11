export type User = {
    id: string;
    email: string;
    username: string;
    memory_enabled: boolean;
};

export type Skill = {
    id: string;
    user_id: string;
    name: string;
    slug: string;
    description: string;
    content: string;
    status: "draft" | "enabled" | "disabled";
};

export type Conversation = {
    id: string;
    user_id: string;
    title: string;
};

export type MCPTransport = "sse" | "streamable";

export type MCPServer = {
    id: string;
    user_id: string;
    name: string;
    transport: MCPTransport;
    command: string;
    args: string[] | null;
    env: Record<string, string> | null;
    url: string;
    headers: Record<string, string> | null;
    enabled: boolean;
};

export type MCPServerPayload = {
    name: string;
    transport: MCPTransport;
    command: string;
    args: string[];
    env: Record<string, string>;
    url: string;
    headers: Record<string, string>;
    enabled: boolean;
};

export type ChatMessage = {
    id?: string;
    conversation_id?: string;
    user_id?: string;
    role: "user" | "assistant";
    content: string;
    reasoning_content?: string;
    elapsed_ms?: number; // 前端流式计时专用，不来自后端，会话结束后清除
    tool_call_count?: number;
    context_tokens?: number;
    context_budget?: number;
};

export type AssistantMeta = {
    tool_call_count?: number;
    context_tokens?: number;
    context_budget?: number;
};

function conversationTitlePayload(title?: string): { title?: string } {
    const trimmed = title?.trim();
    return trimmed ? { title: trimmed } : {};
}

const API_BASE = (import.meta.env.VITE_API_BASE as string | undefined) ?? "http://localhost:8080";

function normalizeArray<T>(value: T[] | null | undefined): T[] {
    return Array.isArray(value) ? value : [];
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE}${path}`, {
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            ...(init?.headers ?? {}),
        },
        ...init,
    });

    if (!response.ok) {
        const text = await response.text();
        throw new Error(text || `Request failed: ${response.status}`);
    }

    if (response.status === 204) {
        return undefined as T;
    }

    return response.json() as Promise<T>;
}

export const api = {
    register: (payload: { email: string; username: string; password: string }) =>
        request<{ user: User }>("/api/auth/register", { method: "POST", body: JSON.stringify(payload) }),
    login: (payload: { login: string; password: string }) =>
        request<{ user: User }>("/api/auth/login", { method: "POST", body: JSON.stringify(payload) }),
    logout: () => request<{ ok: boolean }>("/api/auth/logout", { method: "POST" }),
    me: () => request<{ user: User }>("/api/me"),
    updateMemoryPreference: (enabled: boolean) =>
        request<{ memory_enabled: boolean }>("/api/me/memory", { method: "PATCH", body: JSON.stringify({ enabled }) }),
    listSkills: async () => {
        const result = await request<{ skills: Skill[] | null }>("/api/skills");
        return { ...result, skills: normalizeArray(result.skills) };
    },
    createSkill: (payload: Pick<Skill, "name" | "description" | "content" | "status">) =>
        request<{ skill: Skill }>("/api/skills", { method: "POST", body: JSON.stringify(payload) }),
    updateSkill: (skillId: string, payload: Pick<Skill, "name" | "description" | "content" | "status">) =>
        request<{ skill: Skill }>(`/api/skills/${skillId}`, { method: "PUT", body: JSON.stringify(payload) }),
    patchSkillStatus: (skillId: string, status: Skill["status"]) =>
        request<{ skill: Skill }>(`/api/skills/${skillId}`, { method: "PATCH", body: JSON.stringify({ status }) }),
    deleteSkill: (skillId: string) => request<{ ok: boolean }>(`/api/skills/${skillId}`, { method: "DELETE" }),
    listMCPServers: async () => {
        const result = await request<{ mcp_servers: MCPServer[] | null }>("/api/mcp-servers");
        return { ...result, mcp_servers: normalizeArray(result.mcp_servers) };
    },
    createMCPServer: (payload: MCPServerPayload) =>
        request<{ mcp_server: MCPServer }>("/api/mcp-servers", { method: "POST", body: JSON.stringify(payload) }),
    updateMCPServer: (id: string, payload: MCPServerPayload) =>
        request<{ mcp_server: MCPServer }>(`/api/mcp-servers/${id}`, { method: "PUT", body: JSON.stringify(payload) }),
    patchMCPServerEnabled: (id: string, enabled: boolean) =>
        request<{ mcp_server: MCPServer }>(`/api/mcp-servers/${id}`, { method: "PATCH", body: JSON.stringify({ enabled }) }),
    deleteMCPServer: (id: string) => request<{ ok: boolean }>(`/api/mcp-servers/${id}`, { method: "DELETE" }),
    testMCPServer: (id: string) =>
        request<{ ok: boolean; tools?: string[]; error?: string }>(`/api/mcp-servers/${id}/test`, { method: "POST" }),
    listConversations: async () => {
        const result = await request<{ conversations: Conversation[] | null }>("/api/conversations");
        return { ...result, conversations: normalizeArray(result.conversations) };
    },
    createConversation: (title?: string) =>
        request<{ conversation: Conversation }>("/api/conversations", { method: "POST", body: JSON.stringify(conversationTitlePayload(title)) }),
    renameConversation: (conversationId: string, title: string) =>
        request<{ conversation: Conversation }>(`/api/conversations/${conversationId}`, {
            method: "PATCH",
            body: JSON.stringify({ title: title.trim() }),
        }),
    deleteConversation: (conversationId: string) =>
        request<{ ok: boolean }>(`/api/conversations/${conversationId}`, { method: "DELETE" }),
    getConversation: async (conversationId: string) => {
        type RawMessage = ChatMessage & { meta?: AssistantMeta | null };
        const result = await request<{ conversation: Conversation; messages: RawMessage[] | null }>(`/api/conversations/${conversationId}`);
        const messages = normalizeArray(result.messages).map((message) => {
            const { meta, ...rest } = message;
            return meta
                ? {
                      ...rest,
                      tool_call_count: meta.tool_call_count,
                      context_tokens: meta.context_tokens,
                      context_budget: meta.context_budget,
                  }
                : rest;
        });
        return { ...result, messages };
    },
    streamConversation: async (
        conversationId: string,
        content: string,
        handlers: {
            onAssistantDelta: (payload: { content: string }) => void;
            onReasoningDelta: (payload: { content: string }) => void;
            onMeta: (payload: AssistantMeta) => void;
            onAssistant: (payload: { message_id?: string; content: string; reasoning_content?: string; final?: boolean } & AssistantMeta) => void;
            onError: (message: string) => void;
            onDone: () => void;
        },
    ) => {
        const response = await fetch(`${API_BASE}/api/conversations/${conversationId}/stream`, {
            method: "POST",
            credentials: "include",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ content }),
        });

        if (!response.ok || !response.body) {
            handlers.onError(await response.text());
            return;
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = "";

        while (true) {
            const { done, value } = await reader.read();
            if (done) {
                handlers.onDone();
                break;
            }

            buffer += decoder.decode(value, { stream: true });
            const chunks = buffer.split("\n\n");
            buffer = chunks.pop() ?? "";

            for (const chunk of chunks) {
                const lines = chunk.split("\n");
                const eventLine = lines.find((line) => line.startsWith("event:"));
                const dataLines = lines.filter((line) => line.startsWith("data:"));
                if (!eventLine || dataLines.length === 0) continue;
                const event = eventLine.replace("event:", "").trim();
                const payload = JSON.parse(dataLines.map((line) => line.replace("data:", "").trim()).join("\n"));
                if (event === "assistant_delta") handlers.onAssistantDelta(payload);
                if (event === "reasoning_delta") handlers.onReasoningDelta(payload);
                if (event === "meta") handlers.onMeta(payload);
                if (event === "assistant") handlers.onAssistant(payload);
                if (event === "error") handlers.onError(payload.message);
                if (event === "done") handlers.onDone();
            }
        }
    },
};
