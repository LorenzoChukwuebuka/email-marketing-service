export type APIResponse<T> = { message: string; payload: T; status: true }


export class APIError extends Error {
    status: number;
    payload: any;

    constructor(message: string, status: number, payload: any) {
        super(message);
        this.name = 'APIError';
        this.status = status;
        this.payload = payload;
    }
}


type JSONValue =
    | string
    | number
    | boolean
    | null
    | JSONObject;

  interface JSONObject {
    [key: string]: JSONValue;
}


export type ResponseT = APIResponse<JSONObject>

