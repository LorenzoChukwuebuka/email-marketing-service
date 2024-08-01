export interface APIResponse<T> {
    message: string,
    payload: T,
    status: boolean
}