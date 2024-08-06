import { create } from 'zustand'
import { PlanData } from '../admin/planStore';
import eventBus from '../../utils/eventBus';
import { errResponse } from '../../utils/isError';
import axiosInstance from '../../utils/api';
import { APIResponse } from '../../interface/api.interface';

interface uPlanStore {
    fetchPlans: () => Promise<void>;
    planData: PlanData[];
    setPlanData: (newPlanData: PlanData[]) => void;
}

const userPlanStore = create<uPlanStore>((set, get) => ({
    planData: [],
    setPlanData: (newPlanData: PlanData[]) => set({ planData: newPlanData }),

    fetchPlans: async () => {
        try {
            const response = await axiosInstance.get<APIResponse<PlanData[]>>("/get-all-plans");
            const plans = response.data.payload;
            get().setPlanData(plans)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data?.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    }
}));

export default userPlanStore;
