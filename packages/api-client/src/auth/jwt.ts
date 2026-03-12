import { api } from "../lib/router";
import type { ApiResponse } from "../types/response";

export async function rotateTokens(): Promise<ApiResponse<{ access: string }>> {
   return api.get("/auth/refresh").json<ApiResponse<{ access: string }>>();
}
