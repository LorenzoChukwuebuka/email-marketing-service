import axiosInstance from "../../../utils/api";
import { APIResponse } from '../../../../../frontend/src/interface/api.interface';
import { EditPlanValues, PlanValues, PlanData } from '../interface/plan.interface';
import { ResponseT } from "../../../interface/api.interface";

class PlanAPI {
    static async getPlans(): Promise<APIResponse<PlanData[]>> {
        const response = await axiosInstance.get<APIResponse<PlanData[]>>("/misc/plan/get");
        return response.data;
    }

    static async createPlans(planValues: PlanValues): Promise<ResponseT> {
        const response = await axiosInstance.post('/admin/plans/create', planValues)
        return response.data
    }

    static async updatePlan(editPlanValues: EditPlanValues): Promise<ResponseT> {
        const response = await axiosInstance.put(
            '/admin/plans/update/' + editPlanValues.id,
            editPlanValues
        )
        return response.data
    }

    static async deletePlan(id: string) {
        const response = await axiosInstance.delete(
            '/admin/plans/delete/' + id
        )
        return response.data
    }

}

export default PlanAPI