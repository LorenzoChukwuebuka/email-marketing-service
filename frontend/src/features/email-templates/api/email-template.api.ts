import axiosInstance from '../../../utils/api';
import { TemplateResponse, Template } from '../interface/email-templates.interface';
import { ResponseT } from '../../../../../frontend/src/interface/api.interface';
import { BaseEntity } from '../../../../../frontend/src/interface/baseentity.interface';


export class TemplateAPI {
    static async getAllMarketingTemplates(page?: number, pageSize?: number, query?: string): Promise<TemplateResponse> {
        console.log('mkt', page, pageSize)
        const response = await axiosInstance.get<TemplateResponse>('/templates/get-all-marketing-templates', {
            params: {
                page: page,
                page_size: pageSize,
                search: query
            }
        });
        return response.data;
    }

    static async getAllTransactionalTemplates(page?: number, pageSize?: number, query?: string): Promise<TemplateResponse> {
        console.log('trc', page, pageSize)
        const response = await axiosInstance.get<TemplateResponse>('/templates/get-all-transactional-templates', {
            params: {
                page: page,
                page_size: pageSize,
                search: query
            }
        });
        return response.data;
    }

    static async createTemplate(formValues: Omit<Template, "user_id">): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>("/templates/create-marketing-template", formValues);
        return response.data;
    }

    static async getSingleTransactionalTemplate(uuid: string) {
        const response = await axiosInstance.get<TemplateResponse>(`/templates/get-transactional-template/${uuid}`);
        return response.data.payload as (Template & BaseEntity);
    }

    static async getSingleMarketingTemplate(uuid: string) {
        const response = await axiosInstance.get<TemplateResponse>(`/templates/get-marketing-template/${uuid}`);
        return response.data.payload as (Template & BaseEntity);
    }

    static async updateTemplate(uuid: string, updatedTemplate: Template & BaseEntity) {
        const response = await axiosInstance.put<ResponseT>(`/templates/update-template/${uuid}`, updatedTemplate);
        return response.data;
    }

    static async deleteTemplate(uuid: string) {
        const response = await axiosInstance.delete<ResponseT>(`/templates/delete-template/${uuid}`);
        return response.data;
    }

    static async sendTestMail(sendEmailTestValues: {
        email_address: string;
        template_id: string;
        subject: string;
    }) {
        const response = await axiosInstance.post<ResponseT>("/templates/send-test-mails", sendEmailTestValues);
        return response.data;
    }
}