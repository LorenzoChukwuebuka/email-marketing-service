export type APIResponse<T> = { message: string; payload: T; status: true }

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

