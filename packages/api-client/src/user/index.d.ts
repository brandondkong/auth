import type { DatabaseModel } from "../types/model";

export interface User extends DatabaseModel {
    email: string;
    name?: string;
}
