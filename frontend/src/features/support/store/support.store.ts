import { create } from 'zustand';
import { handleError } from '../../../utils/isError';
import eventBus from '../../../utils/eventbus';
import { SupportRequestValues, Ticket } from '../interface/support.interface';
import { SupportTicket } from '../api/supportticket.api';

type SupportTicketStore = {
    supportTicketData: Ticket[];
    setSupportData: (newData: any) => void;
    getTickets: () => Promise<void>;
    createTicket: (values: SupportRequestValues, files: File[]) => Promise<void>;
    getTicketDetails: (uuid: string) => Promise<void>;
    replyTicket: (ticketId: string, message: string, file: File[]) => Promise<void>;
    closeTicket: (ticketId: string) => Promise<void>;
}

const useSupportStore = create<SupportTicketStore>((set, get) => ({
    supportTicketData: [],
    setSupportData: (newData) => set({ supportTicketData: newData }),

    getTickets: async () => {
        try {
            const data = await SupportTicket.getTickets();
            if (data.status === true) {
                get().setSupportData(data.payload);
            }
        } catch (error) {
            handleError(error)
        }
    },

    createTicket: async (values: SupportRequestValues, files: File[]) => {
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
                formData.append('file', file);
            });

            const data = await SupportTicket.createTicket(formData);
            if (data.status === true) {
                eventBus.emit('success', 'Ticket created successfully');
                await new Promise(resolve => setTimeout(resolve, 500));
                window.location.reload();
            }
        } catch (error) {
            handleError(error)
        }
    },

    getTicketDetails: async (uuid: string) => {
        try {
            const data = await SupportTicket.getTicketDetails(uuid);
            if (data.status === true) {
                get().setSupportData(data.payload);
            }
        } catch (error) {
            handleError(error)
        }
    },

    replyTicket: async (ticketId: string, message: string, files: File[]) => {
        try {
            const formData = new FormData();

            formData.append('message', message);
            files.forEach((file) => {
                formData.append('file', file);
            });

            const data = await SupportTicket.replyTicket(formData,ticketId);
            if (data.status === true) {
                eventBus.emit('success', 'Reply sent successfully');
                await new Promise(resolve => setTimeout(resolve, 500));
                window.location.reload();
            }
        } catch (error) {
            handleError(error)
        }
    },

    closeTicket: async (ticketId: string) => {
        try {
            const data = await SupportTicket.closeTicket(ticketId);
            if (data.status === true) {
                eventBus.emit('success', data.payload);
                await new Promise(resolve => setTimeout(resolve, 500));
                window.location.reload();
            }
        } catch (error) {
            handleError(error)
        }
    }
}));

export default useSupportStore;