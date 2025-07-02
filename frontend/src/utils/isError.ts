import axios, { AxiosError } from "axios";
import eventBus from "./eventbus";

interface ErrorResponse {
    status: boolean;
    message: string;
    error?: string;
    payload?: any;
}

export function errResponse(error: unknown): error is AxiosError<ErrorResponse> {
    return axios.isAxiosError(error) && error.response?.data !== undefined;
}

export const handleError = (error: unknown): void => {
    if (errResponse(error)) {
        const { error: errText, payload, message } = error.response!.data;

        if (errText) {
            eventBus.emit("error", errText);
        } else if (payload) {
            eventBus.emit("error", payload);
        } else {
            eventBus.emit("error", message || "An unknown error occurred.");
        }
    } else if (error instanceof Error) {
        eventBus.emit("error", error.message);
    } else {
        console.error("Unknown error:", error);
    }
};
