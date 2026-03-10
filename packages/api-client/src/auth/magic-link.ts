import type { ApiResponse } from "../types/response"
import { api } from "../lib/router"

export async function generateMagicLink(email: string): Promise<ApiResponse> {
    return await api.post("auth/magic-link", {
        json: { email },
    }).json<ApiResponse>();
}

export async function verifyMagicLink(token: string): Promise<ApiResponse> {
    return await api.get(`auth/magic-link/callback?token=${token}`)
    .json<ApiResponse>();
}
