import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { errResponse } from '../../utils/isError';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { BaseEntity } from '../../interface/baseentity.interface';

type Template = {
    user_id: string;
    templateName: string;
    campaignId: number | null;
    senderName: string | null;
    fromEmail: string | null;
    subject: string | null;
    type: string;
    emailHtml: string;
    emailDesign: any;
    isEditable: boolean;
    isPublished: boolean;
    isPublicTemplate: boolean;
    isGalleryTemplate: boolean;
    tags: string;
    description: string | null;
    imageUrl: string | null;
    isActive: boolean;
    editorType: string | null;
};

type TemplateStore = {
    formValues: Omit<Template, "user_id">;
    setFormValues: (newFormValues: Omit<Template, "user_id">) => void;
    getAllTemplates: () => Promise<void>;
    createTemplate: () => Promise<void>
};


type TemplateResponse = APIResponse<Template & BaseEntity>

const useTemplateStore = create<TemplateStore>((set, get) => ({
    formValues: {
        templateName: '',
        campaignId: null,
        senderName: null,
        fromEmail: null,
        subject: null,
        type: '',
        emailHtml: '',
        emailDesign: null,
        isEditable: false,
        isPublished: false,
        isPublicTemplate: false,
        isGalleryTemplate: false,
        tags: '',
        description: null,
        imageUrl: null,
        isActive: true,
        editorType: null,
    },

    setFormValues: (newFormValues) => set({ formValues: newFormValues }),

    getAllTemplates: async () => {
        try {
            const response = await axiosInstance.get<TemplateResponse>('/templates');

        } catch (error) {
            errResponse(error);
            console.error("Failed to fetch templates", error);
        }
    },
    createTemplate: async () => {
        try {
            const { formValues } = get()
            let response = await axiosInstance.post<ResponseT>("/create-martketing-template", formValues)
            console.log(response.data.payload)

            window.location.href = "/editor/1?" + response.data.payload.templateId

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
