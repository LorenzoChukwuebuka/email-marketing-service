import { QueryClient, useQuery } from "@tanstack/react-query";
import notificationApi from "../api/notification.api";
import { UserNotification } from "../interface/notification.interface";
import { APIResponse } from '../../../frontend/src/interface/api.interface';

export const NOTIFICATIONS_QUERY_KEY = ["notifications_key"] as const


// Query options that will be used consistently across components
export const notificationsQueryOptions = {
    //add all the dynamic query to the key 
    queryKey: [NOTIFICATIONS_QUERY_KEY],
    queryFn: async () => {
        const response = await notificationApi.getNotifications();
        return response ?? [];
    },
    staleTime: 1 * 60 * 1000, // Data is fresh for 5 minutes
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 1 * 60 * 1000, // Refetch every 1 minutes
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
};

// Custom hook for using the employees query
export const useNotificationQuery = () => {
    return useQuery<APIResponse<UserNotification[]>>(
        notificationsQueryOptions
    );
};

// Prefetch function that can be used anywhere
export const prefetchNotifications = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(notificationsQueryOptions);
};

export const invalidateNotification = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: NOTIFICATIONS_QUERY_KEY });
};