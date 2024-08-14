import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { errResponse } from '../../utils/isError';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { BaseEntity } from '../../interface/baseentity.interface';

export type Template = {
    user_id: string;
    template_name: string;
    campaignId: number | null;
    sender_name: string | null;
    from_email: string | null;
    subject: string | null;
    type: string;
    email_html: string;
    email_design: any;
    is_editable: boolean;
    is_published: boolean;
    is_public_template: boolean;
    is_gallery_template: boolean;
    tags: string;
    description: string | null;
    image_Url: string | null;
    is_active: boolean;
    editor_type: string | null;
};

type TemplateStore = {
    formValues: Omit<Template, "user_id">;
    templateData: (Template & BaseEntity)[] | (Template & BaseEntity)
    _templateData: (Template & BaseEntity)[] | (Template & BaseEntity)
    setTemplateData: (newData: (Template & BaseEntity)[] | (Template & BaseEntity)) => void
    _setTemplateData: (newData: (Template & BaseEntity)[] | (Template & BaseEntity)) => void
    currentTemplate: (Template & BaseEntity) | null;
    setCurrentTemplate: (template: (Template & BaseEntity) | null) => void;
    setFormValues: (newFormValues: Omit<Template, "user_id">) => void;
    getAllMarketingTemplates: () => Promise<void>;
    getAllTransactionalTemplates: () => Promise<void>;
    createTemplate: () => Promise<void>
    getSingleTransactionalTemplate: (uuid: string) => Promise<void>
    getSingleMarketingTemplate: (uuid: string) => Promise<void>
    updateTemplate: (uuid: string, updatedTemplate: Template & BaseEntity) => Promise<void>
};


type TemplateResponse = APIResponse<(Template & BaseEntity)[] | (Template & BaseEntity)>

const useTemplateStore = create<TemplateStore>((set, get) => ({
    formValues: {
        template_name: '',
        campaignId: null,
        sender_name: null,
        from_email: null,
        subject: null,
        type: '',
        email_html: '',
        email_design: null,
        is_editable: false,
        is_published: false,
        is_public_template: false,
        is_gallery_template: false,
        tags: '',
        description: null,
        image_Url: null,
        is_active: true,
        editor_type: null,
    },

    templateData: [],
    _templateData: [],
    currentTemplate: null,
    setTemplateData: (newData) => set({ templateData: newData }),
    setFormValues: (newFormValues) => set({ formValues: newFormValues }),
    _setTemplateData: (newData) => set({ _templateData: newData }),
    setCurrentTemplate: (template) => set({ currentTemplate: template }),

    getAllMarketingTemplates: async () => {
        try {

            const response = await axiosInstance.get<TemplateResponse>('/templates/get-all-marketing-templates');
            get().setTemplateData(response.data.payload)

        } catch (error) {
            errResponse(error);
            console.error("Failed to fetch templates", error);
        }
    },
    createTemplate: async () => {
        try {
            const { formValues } = get()
            let response = await axiosInstance.post<ResponseT>("/create-martketing-template", formValues)
            window.location.href = "/editor/1?id=" + response.data.payload.templateId

        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },

    getAllTransactionalTemplates: async () => {
        try {
            const response = await axiosInstance.get<TemplateResponse>('/templates/get-all-transactional-templates');
            get()._setTemplateData(response.data.payload)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },
    getSingleTransactionalTemplate: async (uuid: string) => {
        try {
            const response = await axiosInstance.get<TemplateResponse>("/templates/get-transactional-template/" + uuid)
            get().setCurrentTemplate(response.data.payload as (Template & BaseEntity));
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },
    getSingleMarketingTemplate: async (uuid: string) => {
        try {
            const response = await axiosInstance.get<TemplateResponse>("/templates/get-marketing-template/" + uuid)
            get().setCurrentTemplate(response.data.payload as (Template & BaseEntity));
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },
    updateTemplate: async (uuid: string, updatedTemplate: Template & BaseEntity) => {
        try {

            console.log("you got here shaa")

            const response = await axiosInstance.put<ResponseT>("/templates/update-template/" + uuid, updatedTemplate)
            console.log(response.data)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    }
}));

export default useTemplateStore;
