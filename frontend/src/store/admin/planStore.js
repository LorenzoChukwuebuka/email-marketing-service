import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import eventBus from '../../utils/eventBus'

const usePlanStore = create((set, get) => ({
  planValues: {
    planname: '',
    duration: '',
    price: '',
    details: '',
    number_of_mails_per_day: ''
  },
  editPlanValues: {
    uuid: '',
    planname: '',
    duration: '',
    price: '',
    details: '',
    number_of_mails_per_day: ''
  },
  isLoading: false,
  planData: [],
  selectedId: [],

  setIsLoading: newIsLoading => set({ isLoading: newIsLoading }),
  setPlanValues: newPlanValues =>
    set(state => ({
      planValues: { ...state.planValues, ...newPlanValues }
    })),
  setLoginValues: newPlanValues => set({ planValues: newPlanValues }),
  setPlanData: newPlanData => set({ planData: newPlanData }),
  setEditPlanValues: newEditPlanValues =>
    set({ editPlanValues: newEditPlanValues }),
  setSelectedId: newSelectedId => set({ selectedId: newSelectedId }),

  createPlan: async () => {
    const { setIsLoading, planValues } = get()

    setIsLoading(true)
    try {
      let response = await axiosInstance.post('/admin/create-plan', planValues)
      if (response.data.status === true) {
        eventBus.emit('success', 'Plan creation was successful')
      }
    } catch (error) {
      eventBus.emit(
        'error',
        error.response.data.payload || 'An unexpected error occured'
      )
    } finally {
      get().setIsLoading(false)
    }
  },
  getPlans: async () => {
    const { setIsLoading, setPlanData } = get()

    try {
      setIsLoading(true)
      let response = await axiosInstance.get('/admin/get-plans')
      setPlanData(response.data.payload)
    } catch (error) {
      eventBus.emit(
        'error',
        error.response.data.payload || 'An unexpected error occured'
      )
    } finally {
      get().setIsLoading(false)
    }
  },
  updatePlan: async () => {
    try {
      let response = await axiosInstance.put(
        '/admin/edit-plan/' + get().editPlanValues.uuid,
        get().editPlanValues
      )
      eventBus.emit('success', response.data.payload)
    } catch (error) {
      eventBus.emit('error', error)
    }
  },

  deletePlan: async () => {
    console.log(get().selectedId)
  }
}))

export default usePlanStore
