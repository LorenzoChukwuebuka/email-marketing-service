import { APIResponse, ResponseT } from "../../../interface/api.interface"
import { PaginatedResponse } from "../../../interface/pagination.interface"
import axiosInstance from "../../../utils/api"
import { AdminUserDetails, AdminUserStats } from "../interface/adminusers.interface"

class AdminUsersAPI {
    async getAllUsers(page?: number, pageSize?: number, query?: string):Promise<APIResponse<PaginatedResponse<AdminUserDetails>>> {
        const response = await axiosInstance.get('/admin/users/get', {
            params: {
                page: page,
                page_size: pageSize,
                search: query
            }
        })

        return response.data
    }

    async getVerifiedUsers(page?: number, pageSize?: number, query?: string):Promise<APIResponse<PaginatedResponse<AdminUserDetails>>> {
        const response = await axiosInstance.get('/admin/users/get/verified', {
            params: {
                page: page,
                page_size: pageSize,
                search: query
            }
        })

        return response.data
    }

    async getUnverifiedUsers(page?: number, pageSize?: number, query?: string):Promise<APIResponse<PaginatedResponse<AdminUserDetails>>> {
        const response = await axiosInstance.get('/admin/users/get/unverified', {
            params: {
                page: page,
                page_size: pageSize,
                search: query
            }
        })

        return response.data
    }

    async blockUser(userId: string):Promise<ResponseT> {
        const response = await axiosInstance.put(`/admin/users/block/${userId}`)
        return response.data
    }

    async unBlockUser(userId: string):Promise<ResponseT> {
        const response = await axiosInstance.put(`/admin/users/unblock/${userId}`)
        return response.data
    }

    async verifyUser(userId: string):Promise<ResponseT> {
        const response = await axiosInstance.put(`/admin/users/verify/${userId}`)
        return response.data
    }

    async getSingleUser(userId: string):Promise<APIResponse<AdminUserDetails>> {
        const response = await axiosInstance.get(`/admin/users/get/${userId}`)
        return response.data
    }

    async getUserStats(userId:string):Promise<APIResponse<AdminUserStats>> { 
        const response = await axiosInstance.get(`/admin/users/stats/${userId}`)
        return response.data
     }
}

export default new AdminUsersAPI()