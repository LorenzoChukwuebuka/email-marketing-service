import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import eventBus from '../../utils/eventBus'
import Cookies from 'js-cookie'
import { errResponse } from '../../utils/isError';
import { UserDetails } from '../userstore/AuthStore';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse, ResponseT } from '../../interface/api.interface';


export type AdminUserDetails = UserDetails & { verified_at: string }

export type UserSearchType = "All" | "Verified" | "Unverified"


type AdminUserStats = {
    total_contacts: number,
    total_campaigns: number,
    total_templates: number,
    total_campaigns_sent: number,
   // total_subscriptions: number,
    total_groups: number
}

type AdminUserStore = {
    userdetailsData: AdminUserDetails[] | AdminUserDetails
    userStatsData: AdminUserStats
    setUserStats: (newData: AdminUserStats) => void
    setUserDetailsData: (newData: AdminUserDetails[] | AdminUserDetails) => void
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<AdminUserDetails>, 'data'>) => void;
    getAllUsers: (page?: number, pageSize?: number, search?: string) => Promise<void>
    getVerifiedUsers: (page?: number, pageSize?: number, search?: string) => Promise<void>
    getUnverifiedUsers: (page?: number, pageSize?: number, search?: string) => Promise<void>
    blockUser: (userId: string) => Promise<void>
    unBlockUser: (userId: string) => Promise<void>
    verifyUser: (userId: string) => Promise<void>
    getSingleUser: (userId: string) => Promise<void>
    paginationInfo: Omit<PaginatedResponse<AdminUserDetails>, 'data'>;
    searchUser: (query?: string, type?: UserSearchType) => void
    getUserStat: (userId: string) => Promise<void>
}

const useAdminUserStore = create<AdminUserStore>((set, get) => ({
    userdetailsData: [],
    paginationInfo: {
        total_count: 0,
        total_pages: 0,
        current_page: 1,
        page_size: 10,
    },

    userStatsData: {
        total_contacts: 0,
        total_campaigns: 0,
        total_templates: 0,
        total_campaigns_sent: 0,
       // total_subscriptions: 0,
        total_groups: 0
    },

    setUserStats: (newData) => set({ userStatsData: newData }),

    setUserDetailsData: (newData) => {
        set({ userdetailsData: newData })
    },
    setPaginationInfo: (newPaginationInfo) => set({ paginationInfo: newPaginationInfo }),

    getAllUsers: async (page = 1, pageSize = 10, query = "") => {
        try {
            const response = await axiosInstance.get('/admin/users/users', {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            })
            const { data, ...paginationInfo } = response.data.payload;
            get().setUserDetailsData(data);
            get().setPaginationInfo(paginationInfo);
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

    getVerifiedUsers: async (page = 1, pageSize = 10, query = "") => {
        try {
            const response = await axiosInstance.get('/admin/users/verified-users', {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            })
            const { data, ...paginationInfo } = response.data.payload;
            get().setUserDetailsData(data);
            get().setPaginationInfo(paginationInfo);
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

    getUnverifiedUsers: async (page = 1, pageSize = 10, query = "") => {
        try {
            const response = await axiosInstance.get('/admin/users/unverified-users', {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            })
            const { data, ...paginationInfo } = response.data.payload;
            get().setUserDetailsData(data);
            get().setPaginationInfo(paginationInfo);
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

    blockUser: async (userId) => {
        try {
            await axiosInstance.put(`/admin/users/block/${userId}`)
            eventBus.emit('success', "You have blocked user with id " + userId)
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

    unBlockUser: async (userId) => {
        try {
            let response = await axiosInstance.put<ResponseT>(`/admin/users/unblock/${userId}`)
            if (response.data.status === true) {
                eventBus.emit('success', "You have unblocked user with id " + userId)
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

    verifyUser: async (userId) => {
        try {
            let response = await axiosInstance.put<ResponseT>(`/admin/users/verify/${userId}`)
            if (response.data.status === true) {
                eventBus.emit('success', "You have verified user with id " + userId)
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

    getSingleUser: async (userId) => {
        try {
            const response = await axiosInstance.get<APIResponse<AdminUserDetails>>(`/admin/users/user/${userId}`)
            get().setUserDetailsData(response.data.payload)
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

    getUserStat: async (userId) => {
        try {
            let response = await axiosInstance.get<APIResponse<AdminUserStats>>(`/admin/users/stats/${userId}`)
            get().setUserStats(response.data.payload)
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

    searchUser: async (query?: string, type: UserSearchType = 'All') => {
        const { getAllUsers, getVerifiedUsers, getUnverifiedUsers } = get();

        if (!query) {
            // If no query is provided, reset to the normal fetch based on the type
            switch (type) {
                case 'Verified':
                    await getVerifiedUsers();
                    break;
                case 'Unverified':
                    await getUnverifiedUsers();
                    break;
                default:
                    await getAllUsers();
            }
            return;
        }

        // Handle search with query
        switch (type) {
            case 'Verified':
                await getVerifiedUsers(1, 10, query);
                break;
            case 'Unverified':
                await getUnverifiedUsers(1, 10, query);
                break;
            default:
                await getAllUsers(1, 10, query);
        }
    }

}))

export default useAdminUserStore