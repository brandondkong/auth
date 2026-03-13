import ky from "ky";
import type { ApiResponse } from "../types/response";

export const api = ky.create({
    prefixUrl: import.meta.env.BACKEND_URL ?? "/api",
});

export function createApi(onAuthFailure: () => void) {
    async function rotateTokens(): Promise<ApiResponse<{ access: string }>> {
       return api.get("/auth/refresh").json<ApiResponse<{ access: string }>>();
    }

    return api.extend({
        hooks: {
            afterResponse: [
                async (req, _, res) => {
                    if (res.status === 401 && req.headers.get("Authorization")) {
                        // Refetch 
                        const access = await rotateTokens();
                        if (!access.success) {
                            onAuthFailure();
                            return;
                        }

                        const headers = new Headers(req.headers);
                        headers.set("Authorization", `Bearer ${access.data.access}`);

                        return ky.retry({
                            request: new Request(req, { headers }),
                            code: "TOKEN_REFRESHED",
                        });
                    }
                }
            ] 
        }
    })
}

