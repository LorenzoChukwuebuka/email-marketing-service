export type APIResponse<T> = {
    message: string,
    payload: T,
    status: boolean
}



// type Response =
//     | { status: "success", data: object }
//     | { status: "error", error: Error };
