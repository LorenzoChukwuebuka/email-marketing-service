import { APIKeyFormValues } from "../interface/apikey.interface";
import axiosInstance from "../../../utils/api";
import { ResponseT } from "../../../interface/api.interface";

class APIKeyAPI {
    static async generateAPIkey(formValues: APIKeyFormValues): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>('/key/apikey/generate', formValues);
        return response.data
    }

    static async deleteAPIKey(apiId: string) {
        let response = await axiosInstance.delete('/key/apikey/delete/' + apiId);
        return response.data
    }

    static async getAPIKey() {
        const response = await axiosInstance.get('/key/apikey/get');
        return response.data
    }
}

export default APIKeyAPI