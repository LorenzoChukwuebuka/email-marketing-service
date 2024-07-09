import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import eventBus from '../../utils/eventBus'
import Cookies from 'js-cookie'

const useAdminAuthStore = create((set, get) => ({
  loginValues: {
    email: '',
    password: ''
  },
  isLoading: false,
  isLoggedIn: false,
  setIsLoading: newIsLoading => set({ isLoading: newIsLoading }),
  setLoginValues: newLoginValues => set({ loginValues: newLoginValues }),
  setIsLoggedIn: newIsLoggedIn => set({ isLoggedIn: newIsLoggedIn }),

  loginAdmin: async () => {
    try {
      const { setIsLoading, loginValues, setLoginValues, setIsLoggedIn } = get()
      setIsLoading(true)

      let response = await axiosInstance.post('admin/admin-login', loginValues)

      if (response.data.message === 'success') {
        //save the user Credentials to a cookie
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
      eventBus.emit(
        'error',
        error.response.data.payload || 'An unexpected error occured'
      )
    } finally {
      get().setIsLoading(false)
    }
  }
}))

export default useAdminAuthStore
