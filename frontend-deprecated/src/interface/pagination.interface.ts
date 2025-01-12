export type PaginatedResponse<T> = {
    data: T[];
    total_count: number;
    total_pages: number;
    current_page: number;
    page_size: number;
};
