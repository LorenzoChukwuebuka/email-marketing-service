import axios, { AxiosError } from "axios";

export function isError(value: unknown): value is Error {
    return value instanceof Error ||
        (typeof value === "object" &&
            value !== null &&
            'message' in value &&
            typeof (value as any).message === 'string');
}




interface ErrorResponse {
    status: string;
    message: string;
    payload?: Record<string, any>
}



export function errResponse(error: any): error is AxiosError<ErrorResponse> {
    return axios.isAxiosError(error) && error.response?.data !== undefined;
}
