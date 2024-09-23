import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { BaseEntity } from '../../interface/baseentity.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { errResponse } from '../../utils/isError';

type TicketFile = {
    file_name: string
} & BaseEntity

type TicketMessage = {
    user_id: string
    message: string
    is_admin: boolean
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
    files: TicketFile[]
    messages: TicketMessage[];
} & BaseEntity;

type SupportRequestValues = Partial<Ticket>

type SupportTicketStore = {
    supportTicketData: Ticket[] | Ticket
    setSupportData: (newData: Ticket[] | Ticket) => void
    getTickets: () => Promise<void>
    createTicket: (values: SupportRequestValues, files: File[]) => Promise<void>
    getTicketDetails: (uuid: string) => Promise<void>

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

            // Append files
            files.forEach((file, index) => {
                formData.append(`file${index + 1}`, file);
            });

            const response = await axiosInstance.post<APIResponse<Ticket>>("/support/create-ticket", formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            if (response.data.status === true) {
                eventBus.emit('success', 'Ticket created successfully');
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
    }
}))


export default useSupportStore