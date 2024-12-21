import axios, { AxiosError } from "axios";
import eventBus from "./eventbus";

interface ErrorResponse {
    status: string;
    message: string;
    errors?: Record<string, any>
}

export function errResponse(error: any): error is AxiosError<ErrorResponse> {
    return axios.isAxiosError(error) && error.response?.data !== undefined;
}


// Error handler utility
export const handleError = (error: unknown): void => {
    if (errResponse(error)) {
        eventBus.emit('error', error?.response?.data.message);
    } else if (error instanceof Error) {
        eventBus.emit('error', error.message);
    } else {
        console.error("Unknown error:", error);
    }
};

