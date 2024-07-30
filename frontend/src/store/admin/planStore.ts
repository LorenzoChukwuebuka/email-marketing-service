import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import eventBus from '../../utils/eventBus'

interface Feature {
    name: string;
    identifier: string;
    count_limit: number;
    size_limit: number;
    is_active: boolean;
    description: string;
}

interface PlanValues {
    planname: string;
    duration: string;
    price: number;
    details: string;
    number_of_mails_per_day: string;
    status: string;
    features: Feature[];
}

interface EditPlanValues extends Partial<PlanValues> {
    uuid: string;
}

interface PlanData {
    uuid: string;
    planname: string;
    duration: string;
    price: number;
    number_of_mails_per_day: string;
    details: string;
    status: string;
    features: {
        ID: number;
        CreatedAt: string;
        UpdatedAt: string;
        DeletedAt: string | null;
        uuid: string;
        plan_id: number;
        name: string;
        identifier: string;
        count_limit: number;
        size_limit: number;
        is_active: boolean;
        description: string;
    }[];
    created_at: string;
    updated_at: string;
}

interface PlanState {
    planValues: PlanValues;
    editPlanValues: EditPlanValues;
    isLoading: boolean;
    planData: PlanData[];
    selectedId: string[];
    setIsLoading: (newIsLoading: boolean) => void;
    setPlanValues: (newPlanValues: Partial<PlanValues>) => void;
    setLoginValues: (newPlanValues: PlanValues) => void;
    setPlanData: (newPlanData: PlanData[]) => void;
    setEditPlanValues: (newEditPlanValues: EditPlanValues) => void;
    setSelectedId: (newSelectedId: string[]) => void;
    createPlan: () => Promise<void>;
    getPlans: () => Promise<void>;
    updatePlan: () => Promise<void>;
    deletePlan: () => Promise<void>;
}

const usePlanStore = create<PlanState>((set, get) => ({
    planValues: {
        planname: '',
        duration: '',
        price: 0,
        details: '',
        number_of_mails_per_day: '',
        status: 'active',
        features: []
    },
    editPlanValues: {
        uuid: '',
        planname: '',
        duration: '',
        price: 0,
        details: '',
        number_of_mails_per_day: '',
        status: 'active',
        features: []
    },
    isLoading: false,
    planData: [],
    selectedId: [],

    setIsLoading: (newIsLoading: boolean) => set({ isLoading: newIsLoading }),
    setPlanValues: (newPlanValues: Partial<PlanValues>) =>
        set(state => ({
            planValues: { ...state.planValues, ...newPlanValues }
        })),
    setLoginValues: (newPlanValues: PlanValues) => set({ planValues: newPlanValues }),
    setPlanData: (newPlanData: PlanData[]) => set({ planData: newPlanData }),
    setEditPlanValues: (newEditPlanValues: EditPlanValues) =>
        set({ editPlanValues: newEditPlanValues }),
    setSelectedId: (newSelectedId: string[]) => set({ selectedId: newSelectedId }),

    createPlan: async () => {
        const { setIsLoading, planValues } = get()

        setIsLoading(true)
        try {
            let response = await axiosInstance.post('/admin/create-plan', planValues)
            if (response.data.status === true) {
                eventBus.emit('success', 'Plan creation was successful')
            }
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
    },
    getPlans: async () => {
        const { setIsLoading, setPlanData } = get()

        try {
            setIsLoading(true)
            let response = await axiosInstance.get('/admin/get-plans')
            setPlanData(response.data.payload)
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
    },
    updatePlan: async () => {
        try {
            let response = await axiosInstance.put(
                '/admin/edit-plan/' + get().editPlanValues.uuid,
                get().editPlanValues
            )
            eventBus.emit('success', response.data.payload)
        } catch (error) {
            eventBus.emit('error', error instanceof Error ? error.message : 'An unexpected error occurred')
        }
    },

    deletePlan: async () => {
        try {
            const { selectedId } = get()

            for (let i = 0; i < selectedId.length; i++) {
                let response = await axiosInstance.delete(
                    '/admin/delete-plan/' + selectedId[i]
                )

                eventBus.emit('success', response.data.payload)
            }
        } catch (error) {
            eventBus.emit('error', error instanceof Error ? error.message : 'An unexpected error occurred')
        } finally {
            get().setSelectedId([])
        }
    }
}))

export default usePlanStore