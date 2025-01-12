import { create } from 'zustand'
import Cookies from 'js-cookie'
import { AdminLoginValues } from '../interface/admin.auth.interface';
import { handleError } from '../../../utils/isError';
import AdminAuthAPI from '../api/admin.auth.api';

interface AdminAuthState {
    loginValues: AdminLoginValues;
    isLoading: boolean;
    isLoggedIn: boolean;
    setLoginValues: (newLoginValues: AdminLoginValues) => void;
    setIsLoggedIn: (newIsLoggedIn: boolean) => void;
    loginAdmin: () => Promise<void>;
}

const useAdminAuthStore = create<AdminAuthState>((set, get) => ({
    loginValues: {
        email: '',
        password: ''
    },
    isLoading: false,
    isLoggedIn: false,

    setLoginValues: (newLoginValues: AdminLoginValues) => set({ loginValues: newLoginValues }),
    setIsLoggedIn: (newIsLoggedIn: boolean) => set({ isLoggedIn: newIsLoggedIn }),

    loginAdmin: async () => {
        try {
            const { loginValues, setIsLoggedIn } = get()
            const response = await AdminAuthAPI.adminlogin(loginValues)
            if (response.message === 'success') {
                Cookies.set('Cookies', JSON.stringify(response.payload), {
                    expires: 7,
                    sameSite: 'Strict',
                    secure: true
                })
            }
            setIsLoggedIn(true)

        } catch (error) {
            handleError(error)
        }
    }
}))

export default useAdminAuthStore