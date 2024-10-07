import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { BaseEntity } from '../../interface/baseentity.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { errResponse } from '../../utils/isError';

export interface Sender extends BaseEntity {
    name: string;
    email: string;
    verified: boolean
    is_signed: boolean
}

export interface SenderFormValues {
    name: string;
    email: string;
}

interface SenderStore {
    senderData: Sender[] | Sender;
    isLoading: boolean;
    error: string | null;
    senderFormValues: SenderFormValues;
    currentPage: number;
    totalPages: number;
    paginationInfo: Omit<PaginatedResponse<Sender>, 'data'>;
    setSenderFormValues: (values: SenderFormValues) => void;
    createSender: () => Promise<void>;
    setSenders: (newData: Sender | Sender[]) => void
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<Sender>, 'data'>) => void;
    getSenders: (page?: number, pageSize?: number, search?: string) => Promise<void>;
    updateSender: (id: string, data: Partial<Sender>) => Promise<void>;
    deleteSender: (id: string) => Promise<void>;
}

const useSenderStore = create<SenderStore>((set, get) => ({
    senderData: [],
    isLoading: false,
    error: null,
    senderFormValues: {
        name: '',
        email: '',
    },
    currentPage: 1,
    totalPages: 1,

    paginationInfo: {
        total_count: 0,
        total_pages: 0,
        current_page: 1,
        page_size: 10,
    },

    setSenderFormValues: (newData) => set({ senderFormValues: newData }),
    setSenders: (newData) => set({ senderData: newData }),
    setPaginationInfo: (newPaginationInfo) => set({ paginationInfo: newPaginationInfo }),

    createSender: async () => {
        set({ isLoading: true, error: null });
        try {
            let response = await axiosInstance.post<ResponseT>("/sender/create-sender", get().senderFormValues)
            if (response.data.status === true) {
                eventBus.emit('success', "Sender has been created successfully")
            }
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            set({ isLoading: false });
        }
    },

    getSenders: async (page = 1, pageSize = 10, query = "") => {
        try {
            let response = await axiosInstance.get<APIResponse<PaginatedResponse<Sender>>>("/sender/get-all-senders", {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            })

            const { data, ...paginationInfo } = response.data.payload;
            get().setSenders(data)
            get().setPaginationInfo(paginationInfo)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            set({ isLoading: false });
        }
    },

    updateSender: async (id: string, data: Partial<Sender>) => {

        try {
            let response = await axiosInstance.put<ResponseT>("/sender/update-sender/" + id, data)

            if (response.data.status === true) {
                eventBus.emit('success', "Sender updated successfully")
            }
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }

    },

    deleteSender: async (id: string) => {
        try {
            let response = await axiosInstance.delete<ResponseT>("/sender/delete-sender/" + id)
            if (response.data.status === true) {
                eventBus.emit("success", "sender deleted successfully")
            }
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },
}));

export default useSenderStore;