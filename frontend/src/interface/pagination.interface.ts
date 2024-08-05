export type PaginatedResponse<T> = {
    data: T[];
    total_count: number;
    total_pages: number;
    current_page: number;
    page_size: number;
};



// async function getPaginatedData<T>(url: string, page: number, pageSize: number): Promise<PaginatedResponse<T>> {
//     const { data: apiResponse } = await axiosInstance.get<APIResponse<PaginatedResponse<T>>>(
//         `${url}?page=${page}&page_size=${pageSize}`
//     );

//     if (apiResponse.status) {
//         return apiResponse.payload;
//     } else {
//         throw new Error(apiResponse.message);
//     }
// }

// // Then in your store:
// getAllContacts: async (page = 1, pageSize = 10) => {
//     try {
//         const { data, ...paginationInfo } = await getPaginatedData<Contact>('/get-all-contacts', page, pageSize);
//         get().setContactData(data);
//         get().setPaginationInfo(paginationInfo);
//     } catch (error) {
//         eventBus.emit('error', error instanceof Error ? error.message : 'An unexpected error occurred');
//     }
// }