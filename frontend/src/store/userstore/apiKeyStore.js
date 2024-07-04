import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import eventBus from '../../utils/eventBus'

const useAPIKeyStore = create((set, get) => ({
  apiKeyData: null,
  isLoading: false,

  setAPIKeyData: newAPIkeyData => set({ apiKeyData: newAPIkeyData }),
  setIsLoading: newIsLoading => set({ isLoading: newIsLoading }),

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
    const { setIsLoading } = get()
    try {
      setIsLoading(true)
      let response = await axiosInstance.post('/generate-apikey')
      return response.data.payload
    } catch (error) {
      console.log(error)
    } finally {
      get().setIsLoading(false)
    }
  }
}))

export default useAPIKeyStore
