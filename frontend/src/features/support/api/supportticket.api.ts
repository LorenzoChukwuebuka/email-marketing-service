import axiosInstance from '../../../utils/api';
import { APIResponse } from '../../../../../frontend/src/interface/api.interface';
import { Ticket } from "./../interface/support.interface"

export class SupportTicket {
    static async getTickets(): Promise<APIResponse<Ticket[]>> {
        const response = await axiosInstance.get<APIResponse<Ticket[]>>("/support/get");
        return response.data;
    }

    static async createTicket(formData: FormData): Promise<APIResponse<Ticket>> {
        const response = await axiosInstance.post<APIResponse<Ticket>>(
            "/support/create",
            formData,
            {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            }
        );
        return response.data;
    }

    static async getTicketDetails(uuid: string): Promise<APIResponse<Ticket>> {
        const response = await axiosInstance.get<APIResponse<Ticket>>(`/support/get/${uuid}`);
        return response.data;
    }

    static async replyTicket(formData: FormData, ticketId: string): Promise<APIResponse<Ticket>> {
        const response = await axiosInstance.put<APIResponse<Ticket>>(
            `/support/reply/${ticketId}`,
            formData,
            {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            }
        );
        return response.data;
    }

    static async closeTicket(ticketId: string): Promise<APIResponse<any>> {
        const response = await axiosInstance.put(`/support/close/${ticketId}`);
        return response.data;
    }
}