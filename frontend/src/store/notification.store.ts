import { create } from 'zustand';
import { handleError } from '../utils/isError';
import { invalidateNotification } from '../hooks/useNotificationQuery';
import { queryClient } from '../utils/queryclient';
import notificationApi from '../api/notification.api';

type UserNotificationStore = {
    updateReadStatus: () => Promise<void>
}

const useUserNotificationStore = create<UserNotificationStore>(() => ({
    updateReadStatus: async () => {
        try {
            await notificationApi.updateReadStatus()
            new Promise(resolve => setTimeout(resolve, 500))
            invalidateNotification(queryClient)
        } catch (error) {
            handleError(error)
        }
    }
}))

export default useUserNotificationStore