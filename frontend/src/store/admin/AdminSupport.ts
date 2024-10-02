import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import eventBus from '../../utils/eventBus'
import Cookies from 'js-cookie'
import { errResponse } from '../../utils/isError';
import { Ticket } from '../userstore/support.store';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse } from '../../interface/api.interface';

type AdminSupportStore = {
    supportData: Ticket[] | Ticket
    paginationInfo: Omit<PaginatedResponse<Ticket>, 'data'>;
    setSupportData: (newData: Ticket | Ticket[]) => void
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<Ticket>, 'data'>) => void;
    getAllTickets: (page?: number, pageSize?: number, search?: string) => Promise<void>
    getAllClosedTickets: (page?: number, pageSize?: number, search?: string) => Promise<void>
    getAllPendingTickets: (page?: number, pageSize?: number, search?: string) => Promise<void>
    replyTicket: (ticketId: string, message: string, files: File[]) => Promise<void>
}

const useAdminSupportStore = create<AdminSupportStore>((set, get) => ({
    supportData: [],
    paginationInfo: {
        total_count: 0,
        total_pages: 0,
        current_page: 1,
        page_size: 10,
    },
    setSupportData: (newData) => set({ supportData: newData }),
    setPaginationInfo: (newData) => set({ paginationInfo: newData }),
    getAllClosedTickets: async (page = 1, pageSize = 10, query = "") => {
        try {
            let response = await axiosInstance.get<APIResponse<PaginatedResponse<Ticket>>>("/admin/support/closed-tickets", {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            })
            if (response.data.status === true) {
                const { data, ...paginationInfo } = response.data.payload;
                get().setSupportData(data);
                get().setPaginationInfo(paginationInfo);
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
    getAllTickets: async (page = 1, pageSize = 10, query = "") => {
        try {
            let response = await axiosInstance.get<APIResponse<PaginatedResponse<Ticket>>>("/admin/support/tickets", {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            })
            if (response.data.status === true) {
                const { data, ...paginationInfo } = response.data.payload;
                get().setSupportData(data);
                get().setPaginationInfo(paginationInfo);
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
    getAllPendingTickets: async (page = 1, pageSize = 10, query = "") => {
        try {
            let response = await axiosInstance.get<APIResponse<PaginatedResponse<Ticket>>>("/admin/support/pending-tickets", {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            })
            if (response.data.status === true) {
                const { data, ...paginationInfo } = response.data.payload;
                get().setSupportData(data);
                get().setPaginationInfo(paginationInfo);
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
    replyTicket: async (ticketId: string, message: string, files: File[]): Promise<void> => {
        try {
            const formData = new FormData();

            // Append the message
            formData.append('message', message);

            // Append the files
            files.forEach((file) => {
                formData.append('file', file);
            });

            const response = await axiosInstance.put<APIResponse<Ticket>>(`/admin/support/reply-ticket/${ticketId}`, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            if (response.data.status === true) {
                eventBus.emit('success', 'Reply sent successfully');
                await new Promise(resolve => setTimeout(resolve, 500))
                window.location.reload()
            } else {
                eventBus.emit('error', 'Failed to send reply');
            }
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },
}))

export default useAdminSupportStore