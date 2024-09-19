import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { BaseEntity } from '../../interface/baseentity.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { errResponse } from '../../utils/isError';


export type UserNotification = {
    title: string
    read_status: boolean
} & BaseEntity


type UserNotificationStore = {
    notificationsData: UserNotification[]
    setNotificationData: (newData: UserNotification[]) => void
    getUserNotifications: () => Promise<void>
    updateReadStatus: () => Promise<void>

}

const useUserNotificationStore = create<UserNotificationStore>((set, get) => ({
    notificationsData: [],
    setNotificationData: (newData) => set({ notificationsData: newData }),
    getUserNotifications: async () => {
        try {
            let response = await axiosInstance.get<APIResponse<UserNotification[]>>("/user-notifications")

            if (response.data.status === true) {
                get().setNotificationData(response.data.payload)
            }
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data?.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }

    },
    updateReadStatus: async () => {
        try {
            await axiosInstance.put("/update-read-status")
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data?.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    }
}))

export default useUserNotificationStore