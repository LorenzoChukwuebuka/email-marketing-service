import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import eventBus from '../../utils/eventBus'

const useAPIKeyStore = create((set, get) => ({
  apiKeyData: null,
  isLoading: false,
  formValues: { name: '' },

  setAPIKeyData: newAPIkeyData => set({ apiKeyData: newAPIkeyData }),
  setIsLoading: newIsLoading => set({ isLoading: newIsLoading }),
  setFormValues: newFormValue => set({ formValues: newFormValue }),

  getAPIKey: async () => {
    try {
      const { setAPIKeyData } = get()

      let response = await axiosInstance.get('/get-apikey')
      setAPIKeyData(response.data)
    } catch (error) {
      eventBus.emit('error', error)
    }
  },
  generateAPIKey: async () => {
    const { setIsLoading, formValues } = get()
    try {
      setIsLoading(true)
      let response = await axiosInstance.post('/generate-apikey', formValues)
      return response.data.payload
    } catch (error) {
      console.log(error)
      eventBus.emit('error', error.response.data.payload)
    } finally {
      get().setIsLoading(false)
    }
  },
  deleteAPIKey: async apiId => {
    try {
      let response = await axiosInstance.delete('/delete-apikey/' + apiId)
      eventBus.emit('success', response.data.payload)
    } catch (error) {
      eventBus.emit('error', error)
    }
  }
}))

export default useAPIKeyStore
