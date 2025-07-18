import { create } from 'zustand'
import { PlanValues, EditPlanValues } from '../interface/plan.interface';
import PlanAPI from '../api/plan.api';
import eventBus from '../../../utils/eventbus';
import { handleError } from '../../../utils/isError';


interface PlanState {
    planValues: Partial<PlanValues>;
    editPlanValues: Partial<EditPlanValues>;
    selectedId: string[];
}


interface PlanActions {
    setPlanValues: (newPlanValues: Partial<PlanValues>) => void;
    setLoginValues: (newPlanValues: PlanValues) => void;
    setEditPlanValues: (newEditPlanValues: EditPlanValues) => void;
    setSelectedId: (newSelectedId: string[]) => void;
}


interface PlanAsyncActions {
    createPlan: () => Promise<void>;
    updatePlan: () => Promise<void>;
    deletePlan: () => Promise<void>;
}


type PlanStore = PlanActions & PlanState & PlanAsyncActions


const InitialState = {
    planValues: {},
    editPlanValues: {},
    selectedId: [],
} satisfies PlanState


const usePlanStore = create<PlanStore>((set, get) => ({


    ...InitialState,
    setPlanValues: (newPlanValues: Partial<PlanValues>) =>
        set(state => ({
            planValues: { ...state.planValues, ...newPlanValues }
        })),
    setLoginValues: (newPlanValues: PlanValues) => set({ planValues: newPlanValues }),

    setEditPlanValues: (newEditPlanValues: EditPlanValues) =>
        set({ editPlanValues: newEditPlanValues }),
    setSelectedId: (newSelectedId: string[]) => set({ selectedId: newSelectedId }),

    createPlan: async () => {
        const { planValues } = get()
        try {
            const response = await PlanAPI.createPlans(planValues as any)
            if (response.payload.status === true) {
                eventBus.emit('success', 'Plan creation was successful')
            }
        } catch (error) {
            handleError(error)
        }

    },

    updatePlan: async () => {
        try {
            const values = get().editPlanValues as any
            const transformed = {
                ...values,
                plan_name: values.name,
            }
            delete transformed.name

            const response = await PlanAPI.updatePlan(transformed)
            eventBus.emit('success', response.status)
        } catch (error) {
            handleError(error)
        }
    },


    deletePlan: async () => {
        try {
            const { selectedId } = get()

            for (let i = 0; i < selectedId.length; i++) {
                const response = await PlanAPI.deletePlan(selectedId[i])
                eventBus.emit('success', response.data.payload)
            }
        } catch (error) {
            handleError(error)
        } finally {
            get().setSelectedId([])
        }
    }
}))

export default usePlanStore