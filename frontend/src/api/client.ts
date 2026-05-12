const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api';

export class ApiError extends Error {
    readonly status: number;
    readonly statusText: string;
    readonly body?: unknown;

    constructor(status: number, statusText: string, body?: unknown) {
        super(`API ${status} ${statusText}`);
        this.name = 'ApiError';
        this.status = status;
        this.statusText = statusText;
        this.body = body;
    }
}

interface RequestOptions {
    signal?: AbortSignal;
}

async function request<T>(
    method: string,
    path: string,
    body?: unknown,
    options: RequestOptions = {}
): Promise<T> {
    const headers: Record<string, string> = { Accept: 'application/json' };
    const init: RequestInit = {
        method,
        headers,
        signal: options.signal,
    };

    if (body !== undefined) {
        headers['Content-Type'] = 'application/json';
        init.body = JSON.stringify(body);
    }

    const response = await fetch(`${BASE_URL}${path}`, init);

    if (!response.ok) {
        let parsedBody: unknown;
        try {
            parsedBody = await response.json();
        } catch {
            parsedBody = undefined;
        }
        throw new ApiError(response.status, response.statusText, parsedBody);
    }

    return response.json() as Promise<T>;
}

export function apiGet<T>(path: string, options: RequestOptions = {}): Promise<T> {
    return request<T>('GET', path, undefined, options);
}

export function apiPost<T>(path: string, body: unknown, options: RequestOptions = {}): Promise<T> {
    return request<T>('POST', path, body, options);
}
