export type User = {
    id: string;
    email: string;
    username: string;
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

export type ChatMessage = {
    id?: string;
    conversation_id?: string;
    user_id?: string;
    role: "user" | "assistant";
    content: string;
    reasoning_content?: string;
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
        const result = await request<{ conversation: Conversation; messages: ChatMessage[] | null }>(`/api/conversations/${conversationId}`);
        return { ...result, messages: normalizeArray(result.messages) };
    },
    streamConversation: async (
        conversationId: string,
        content: string,
        handlers: {
            onAssistantDelta: (payload: { content: string }) => void;
            onReasoningDelta: (payload: { content: string }) => void;
            onAssistant: (payload: { message_id?: string; content: string; reasoning_content?: string; final?: boolean }) => void;
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
                if (event === "assistant") handlers.onAssistant(payload);
                if (event === "error") handlers.onError(payload.message);
                if (event === "done") handlers.onDone();
            }
        }
    },
};
