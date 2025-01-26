import { create } from 'zustand'
import { handleError } from '../../../utils/isError';
import adminsupportApi from '../api/adminsupport.api';
import eventBus from '../../../utils/eventbus';

type AdminSupportStore = {
    replyTicket: (ticketId: string, message: string, files: File[]) => Promise<void>
}

const useAdminSupportStore = create<AdminSupportStore>(() => ({

    replyTicket: async (ticketId: string, message: string, files: File[]): Promise<void> => {
        try {
            const formData = new FormData();
            // Append the message
            formData.append('message', message);
            // Append the files
            files.forEach((file) => {
                formData.append('file', file);
            });
            const response = await adminsupportApi.replyTickets(ticketId, formData)
            if (response.status === true) {
                eventBus.emit('success', 'Reply sent successfully');
                await new Promise(resolve => setTimeout(resolve, 500))
                window.location.reload()
            } else {
                eventBus.emit('error', 'Failed to send reply');
            }
        } catch (error) {
            handleError(error)
        }
    },
}))

export default useAdminSupportStore