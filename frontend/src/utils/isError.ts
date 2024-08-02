export function isError(value: unknown): value is Error {
    return value instanceof Error ||
        (typeof value === "object" &&
            value !== null &&
            'message' in value &&
            typeof (value as any).message === 'string');
}