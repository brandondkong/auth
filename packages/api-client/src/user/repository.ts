import { api } from "../lib/router";
import type { ApiResponse } from "../types/response";
import type { User } from "./index";

export async function getUser(): Promise<ApiResponse<User>> {
    return api.get("/user").json<ApiResponse<User>>();
}
