import { APIKeyFormValues } from "../interface/apikey.interface";
import axiosInstance from "../../../utils/api";
import { ResponseT } from "../../../interface/api.interface";

class APIKeyAPI {
    static async generateAPIkey(formValues: APIKeyFormValues): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>('/apikey/generate-apikey', formValues);
        return response.data
    }

    static async deleteAPIKey(apiId: string) {
        let response = await axiosInstance.delete('/apikey/delete-apikey/' + apiId);
        return response.data
    }

    static async getAPIKey() {
        const response = await axiosInstance.get('/apikey/get-apikey');
        return response.data
    }
}

export default APIKeyAPI