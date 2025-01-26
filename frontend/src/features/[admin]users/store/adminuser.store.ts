import { create } from 'zustand'
import { handleError } from '../../../utils/isError';
import eventBus from '../../../utils/eventbus';
import adminusersApi from '../api/adminusers.api';


type AdminUserStore = { 
    blockUser: (userId: string) => Promise<void>
    unBlockUser: (userId: string) => Promise<void>
    verifyUser: (userId: string) => Promise<void>
}

const useAdminUserStore = create<AdminUserStore>(() => ({ 
    blockUser: async (userId) => {
        try {
            const response = await adminusersApi.blockUser(userId)
            if(response){
                eventBus.emit('success', "You have blocked user with id " + userId)
            }
        } catch (error) {
          handleError(error)
        }
    },

    unBlockUser: async (userId) => {
        try {
            const response = await adminusersApi.unBlockUser(userId)
            if (response) {
                eventBus.emit('success', "You have unblocked user with id " + userId)
            }
        } catch (error) {
          handleError(error)
        }
    },

    verifyUser: async (userId) => {
        try {
            const response = await adminusersApi.verifyUser(userId)
            if (response.status === true) {
                eventBus.emit('success', "You have verified user with id " + userId)
            }

        } catch (error) {
          handleError(error)
        }
    },



}))

export default useAdminUserStore