export type BaseEntity = {
    id:string
    company_id?: string;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
}

export type WithPrefix<T, P extends string> = {
  [K in keyof T as `${P}${K & string}`]: T[K];
};