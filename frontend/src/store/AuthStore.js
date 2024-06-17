import { create } from 'zustand'
import axiosInstance from '../utils/api'

const useAuthStore = create((set, get) => ({
  formValues: {
    fullname: '',
    company: '',
    email: '',
    password: '',
    confirmPassword: ''
  },
  error: false,
  success: false,
  errorMessage: '',
  successMessage: '',
  formError: false,
  isLoading: false,
  redirectToOTP: false,

  setFormValues: newFormValues => set({ formValues: newFormValues }),
  setError: newError => set({ error: newError }),
  setSuccess: newSuccess => set({ success: newSuccess }),
  setErrorMessage: newErrorMessage => set({ errorMessage: newErrorMessage }),
  setSuccessMessage: newSuccessMessage =>
    set({ successMessage: newSuccessMessage }),

  setIsLoading: newIsLoading => set({ isLoading: newIsLoading }),
  setRedirectToOTP: newRedirectOTP => set({ redirectToOTP: newRedirectOTP }),

  registerUser: async () => {
    const { setIsLoading, setFormValues, setRedirectToOTP } = get()
    try {
      setIsLoading(true)
      const { formValues } = get()
      let response = await axiosInstance.post('/user-signup', formValues)
      if (response.data.status === true) {
        setRedirectToOTP(true)
        setFormValues({
          fullname: '',
          company: '',
          email: '',
          password: '',
          confirmPassword: ''
        })
      }
    } catch (error) {
      console.log(error)
    } finally {
      get().setIsLoading(false)
    }
  }
}))

export default useAuthStore
