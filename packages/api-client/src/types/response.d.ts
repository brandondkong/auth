export interface ApiResponse<T extends {} = {}> {
    success: boolean;
    error: string;
    message: string;
    data: T;
}
