import { APIResponse } from "../../../interface/api.interface";
import axiosInstance from "../../../utils/api"
import { SMTPKeyDATA, SMTPKeyFormValues } from "../interface/smtpkey.interface"

class SMTPKeyAPI {
    static async getSmtpKeys():Promise<APIResponse<SMTPKeyDATA>> {
        const response = await axiosInstance.get<APIResponse<SMTPKeyDATA>>('/smtpkey/get-smtp-keys')
        return response.data
    }

    static async createSMTPKey(formValues: SMTPKeyFormValues) {
        const response = await axiosInstance.post('/smtpkey/create-smtp-key', formValues);
        return response.data
    }

    static async deleteSMTPKey(id: string) {
        const response = await axiosInstance.delete(`/smtpkey/delete-smtp-key/${id}`);
        return response.data
    }

    static async generateSMTPKey() {
        const response = await axiosInstance.put('/smtpkey/generate-new-smtp-master-password');
        return response.data
    }
}




export default SMTPKeyAPI