import ky from "ky";

export const api = ky.create({
    prefixUrl: import.meta.env.BACKEND_URL ?? "http://127.0.0.1:5000/api"
});
