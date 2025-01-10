import { UserNotification } from "../interface/notification.interface";
import { APIResponse } from '../../../frontend/src/interface/api.interface';
import axiosInstance from "../utils/api";

class NotificationAPI {
    async getNotifications(): Promise<APIResponse<UserNotification[]>> {
        const response = await axiosInstance.get<APIResponse<UserNotification[]>>("/user-notifications");
        return response.data
    }

    async updateReadStatus(): Promise<APIResponse<UserNotification>> {
        const response = await axiosInstance.put<APIResponse<UserNotification>>("/update-read-status");
        return response.data
    }
}

export default new NotificationAPI();