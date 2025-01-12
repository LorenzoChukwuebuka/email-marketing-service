import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { BaseEntity } from '../../interface/baseentity.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { errResponse } from '../../utils/isError';
import { UserDetails } from './AuthStore';
import { Admin } from '../admin/AdminAuthStore';

export type TicketFile = {
    file_name: string
    file_path: string
} & BaseEntity

type TicketMessage = {
    user_id: string
    message: string
    is_admin: boolean
    user: Partial<UserDetails>
    admin: Partial<Admin>
    files: TicketFile[]
} & BaseEntity

export type Ticket = {
    user_id: string;
    name: string;
    email: string;
    subject: string;
    description: string;
    status: string;
    ticket_number: string
    priority: string;
    last_reply: string;
    messages: TicketMessage[];
} & BaseEntity;

type SupportRequestValues = Partial<Ticket>

type SupportTicketStore = {
    supportTicketData: Ticket[] | Ticket
    setSupportData: (newData: Ticket[] | Ticket) => void
    getTickets: () => Promise<void>
    createTicket: (values: SupportRequestValues, files: File[]) => Promise<void>
    getTicketDetails: (uuid: string) => Promise<void>
    replyTicket: (ticketId: string, message: string, file: File[]) => Promise<void>
    closeTicket: (ticketId: string) => Promise<void>
}


const useSupportStore = create<SupportTicketStore>((set, get) => ({
    supportTicketData: [],
    setSupportData: (newData) => set({ supportTicketData: newData }),
    getTickets: async () => {
        try {
            let response = await axiosInstance.get<APIResponse<Ticket[]>>("/support/get-tickets")

            if (response.data.status === true) {
                get().setSupportData(response.data.payload)
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

    createTicket: async (values: SupportRequestValues, files: File[]): Promise<void> => {
        try {
            const formData = new FormData();

            // Append ticket data
            Object.entries(values).forEach(([key, value]) => {
                if (value !== undefined) {
                    formData.append(key, (value as string).toString());
                }
            });

            // Append files using the same 'file' key for each file
            files.forEach((file) => {
                formData.append('file', file);  // Using 'file' instead of 'files[]'
            });

            const response = await axiosInstance.post<APIResponse<Ticket>>("/support/create-ticket", formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            if (response.data.status === true) {
                eventBus.emit('success', 'Ticket created successfully');
                await new Promise(resolve => setTimeout(resolve, 500))
                window.location.reload()
            } else {
                eventBus.emit('error', 'Failed to create ticket');
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

    getTicketDetails: async (uuid: string) => {
        try {
            let response = await axiosInstance.get<APIResponse<Ticket>>("/support/get-ticket/" + uuid)

            if (response.data.status === true) {
                get().setSupportData(response.data.payload)
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

    replyTicket: async (ticketId: string, message: string, files: File[]): Promise<void> => {
        try {
            const formData = new FormData();

            // Append the message
            formData.append('message', message);

            // Append the files
            files.forEach((file) => {
                formData.append('file', file);
            });

            const response = await axiosInstance.put<APIResponse<Ticket>>(`/support/reply-ticket/${ticketId}`, formData, {
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
    closeTicket: async (ticketId: string) => {
        try {
            let response = await axiosInstance.put("/support/close/" + ticketId)
            if (response.data.status === true) {
                eventBus.emit('success', response.data.payload)
                await new Promise(resolve => setTimeout(resolve, 500))
                window.location.reload()
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
    }
}))


export default useSupportStore