import axiosInstance from "../../../utils/api"
import { APIResponse } from '../../../interface/api.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';
import { Ticket } from '../interface/support.interface';

class AdminSupportAPI {
    async getAllClosedTickets(page?: number, pageSize?: number, query?: string) {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<Ticket>>>("/admin/support/closed-tickets", {
            params: {
                page: page,
                page_size: pageSize,
                search: query
            }
        })

        return response.data
    }

    async getAllTickets(page?: number, pageSize?: number, query?: string) {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<Ticket>>>("/admin/support/tickets", {
            params: {
                page: page,
                page_size: pageSize,
                search: query
            }
        })

        return response.data
    }

    async getAllPendingTickets(page?: number, pageSize?: number, query?: string) {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<Ticket>>>("/admin/support/pending-tickets", {
            params: {
                page: page || undefined,
                page_size: pageSize || undefined,
                search: query || undefined
            }
        })

        return response.data
    }

    async replyTickets(ticketId: string, formData: FormData) {
        const response = await axiosInstance.put<APIResponse<Ticket>>(`/admin/support/reply-ticket/${ticketId}`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });

        return response.data
    }
}

export default new AdminSupportAPI()