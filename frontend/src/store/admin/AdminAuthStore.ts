import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import eventBus from '../../utils/eventBus'
import Cookies from 'js-cookie'

interface LoginValues {
    email: string;
    password: string;
}

interface AdminAuthState {
    loginValues: LoginValues;
    isLoading: boolean;
    isLoggedIn: boolean;
    setIsLoading: (newIsLoading: boolean) => void;
    setLoginValues: (newLoginValues: LoginValues) => void;
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
    setIsLoading: (newIsLoading: boolean) => set({ isLoading: newIsLoading }),
    setLoginValues: (newLoginValues: LoginValues) => set({ loginValues: newLoginValues }),
    setIsLoggedIn: (newIsLoggedIn: boolean) => set({ isLoggedIn: newIsLoggedIn }),

    loginAdmin: async () => {
        try {
            const { setIsLoading, loginValues, setLoginValues, setIsLoggedIn } = get()
            setIsLoading(true)

            let response = await axiosInstance.post('admin/admin-login', loginValues)

            if (response.data.message === 'success') {
                Cookies.set('Cookies', JSON.stringify(response.data.payload), {
                    expires: 7,
                    sameSite: 'Strict',
                    secure: true
                })
            }

            setIsLoggedIn(true)

            setLoginValues({
                email: '',
                password: ''
            })
        } catch (error) {
            if (error instanceof Error) {
                eventBus.emit(
                    'error',
                    (error as any).response?.data?.payload || 'An unexpected error occurred'
                )
            } else {
                eventBus.emit('error', 'An unexpected error occurred')
            }
        } finally {
            get().setIsLoading(false)
        }
    }
}))

export default useAdminAuthStore