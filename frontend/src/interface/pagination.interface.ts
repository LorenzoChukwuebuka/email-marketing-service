export type PaginatedResponse<T> = {
    currentPage: number;
    data: T[];
    nextPage: number | null;
    perPage: number;
    prevPage: number | null;
    total: number;
    totalPages: number;
};

export interface PaginationParams {
    page?: number;
    pageSize?: number;
}


