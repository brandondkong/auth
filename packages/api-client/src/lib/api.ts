import { useState, useEffect, useCallback } from "react";
import type { ApiResponse } from "../types/response";
import { HTTPError } from "ky";

export interface ExtendedApiResponse<T extends {} = {}> {
    loading: boolean;
    fetchError: string | undefined;
    response: ApiResponse<T> | undefined;
    refetch: () => Promise<void>;
}

export function useApi<T extends {} = {}>(
    callback: () => Promise<ApiResponse<T>>,
    dependencies?: any[]
): ExtendedApiResponse<T> {
    const [loading, setLoading] = useState<boolean>(true);
    const [fetchError, setFetchError] = useState<string | undefined>(undefined);
    const [response, setResponse] = useState<ApiResponse<T> | undefined>(undefined);

    const fetch = async () => {
        setLoading(true);
        setFetchError(undefined);
        try {
            const res = await callback();
            setResponse(res);
        } catch (err) {
            setFetchError(err instanceof Error ? err.message : "Unknown error");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetch();
    }, dependencies ?? []);

    return { loading, fetchError, response, refetch: fetch };
}

export function useLazyApi<T extends {} = {}>(
    callback: () => Promise<ApiResponse<T>>,
    dependencies?: any[]
): ExtendedApiResponse<T> {
    const [loading, setLoading] = useState<boolean>(false);
    const [fetchError, setFetchError] = useState<string | undefined>(undefined);
    const [response, setResponse] = useState<ApiResponse<T> | undefined>(undefined);

    const fetch = useCallback(async () => {
        setLoading(true);
        setFetchError(undefined);
        try {
            const res = await callback();
            setResponse(res);
        } catch (err) {
            if (err instanceof HTTPError) {
                const errorResponse = await err.response.json<ApiResponse<T>>();
                setResponse(errorResponse)
            }
            else {
                setFetchError(err instanceof Error ? err.message : "Unknown error");
            }
        } finally {
            setLoading(false);
        }
    }, dependencies ?? []);

    return { loading, fetchError, response, refetch: fetch };
}

