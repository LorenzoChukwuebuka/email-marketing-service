import { APIResponse } from "../../../interface/api.interface";
import axiosInstance from "../../../utils/api"
import { SMTPKeyDATA, SMTPKeyFormValues } from "../interface/smtpkey.interface"

class SMTPKeyAPI {
    static async getSmtpKeys():Promise<APIResponse<SMTPKeyDATA>> {
        const response = await axiosInstance.get<APIResponse<SMTPKeyDATA>>('/key/smtpkey/get')
        return response.data
    }

    static async createSMTPKey(formValues: SMTPKeyFormValues) {
        const response = await axiosInstance.post('/key/smtpkey/create', formValues);
        return response.data
    }

    static async deleteSMTPKey(id: string) {
        const response = await axiosInstance.delete(`/key/smtpkey/delete/${id}`);
        return response.data
    }

    static async generateSMTPKey() {
        const response = await axiosInstance.post('/key/smtpkey/generate-new-masterkey');
        return response.data
    }
}


export default SMTPKeyAPI