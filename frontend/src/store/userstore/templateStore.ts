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


export type SendTestMailValues = {
    template_id: string
    email_address: string
    subject: string
}

type TemplateStore = {
    formValues: Omit<Template, "user_id">;
    sendEmailTestValues: SendTestMailValues
    templateData: (Template & BaseEntity)[] | (Template & BaseEntity)
    _templateData: (Template & BaseEntity)[] | (Template & BaseEntity)
    setTemplateData: (newData: (Template & BaseEntity)[] | (Template & BaseEntity)) => void
    _setTemplateData: (newData: (Template & BaseEntity)[] | (Template & BaseEntity)) => void
    setEmailTestValues: (newData: SendTestMailValues) => void
    currentTemplate: (Template & BaseEntity) | null;
    setCurrentTemplate: (template: (Template & BaseEntity) | null) => void;
    setFormValues: (newFormValues: Omit<Template, "user_id">) => void;
    getAllMarketingTemplates: (query?: string) => Promise<void>;
    getAllTransactionalTemplates: (query?: string) => Promise<void>;
    createTemplate: () => Promise<void>
    getSingleTransactionalTemplate: (uuid: string) => Promise<void>
    getSingleMarketingTemplate: (uuid: string) => Promise<void>
    updateTemplate: (uuid: string, updatedTemplate: Template & BaseEntity) => Promise<void>
    deleteTemplate: (uuid: string) => Promise<void>
    resetForm: () => void
    sendTestMail: () => Promise<void>
    searchMarketing: (query: string) => Promise<void>
    searchTransactional: (query: string) => Promise<void>
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

    sendEmailTestValues: {
        email_address: "",
        template_id: "",
        subject: ""
    },

    templateData: [],
    _templateData: [],
    currentTemplate: null,
    setTemplateData: (newData) => set({ templateData: newData }),
    setFormValues: (newFormValues) => set({ formValues: newFormValues }),
    _setTemplateData: (newData) => set({ _templateData: newData }),
    setCurrentTemplate: (template) => set({ currentTemplate: template }),
    setEmailTestValues: (newData) => set({ sendEmailTestValues: newData }),

    getAllMarketingTemplates: async (query = "") => {
        try {
            const response = await axiosInstance.get<TemplateResponse>('/templates/get-all-marketing-templates', {
                params: {
                    search: query || undefined
                }
            });
            get().setTemplateData(response.data.payload)

        } catch (error) {
            errResponse(error);
            console.error("Failed to fetch templates", error);
        }
    },
    createTemplate: async () => {
        try {
            const { formValues } = get()
            let response = await axiosInstance.post<ResponseT>("/templates/create-marketing-template", formValues)

            const editorType = response.data.payload["editor-type"];
            const templateType = response.data.payload.type;
            const templateId = response.data.payload.templateId;

            let redirectUrl = "";
            const typeParam = templateType === "marketing" ? "m" : "t";

            switch (editorType) {
                case "html-editor":
                    redirectUrl = `/editor/2?type=${typeParam}&uuid=${templateId}`;
                    break;
                case "drag-and-drop":
                    redirectUrl = `/editor/1?type=${typeParam}&uuid=${templateId}`;
                    break;
                case "rich-text":
                    redirectUrl = `/editor/3?type=${typeParam}&uuid=${templateId}`
                    break;
                default:
                    console.log("Unknown editor type:", editorType);
                    return;
            }
            window.location.href = redirectUrl;

        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            get().resetForm()
        }
    },

    getAllTransactionalTemplates: async (query = "") => {
        try {
            const response = await axiosInstance.get<TemplateResponse>('/templates/get-all-transactional-templates', {
                params: {
                    search: query || undefined
                }
            });
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
    },

    deleteTemplate: async (uuid: string) => {
        try {
            const response = await axiosInstance.delete<ResponseT>("/templates/delete-template/" + uuid)
            if (response.data.status === true) {
                eventBus.emit("success", "template deleted successfully")
            }
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


    sendTestMail: async () => {
        try {
            const { sendEmailTestValues } = get()
            let address = sendEmailTestValues.email_address
            let str = address.split(",")
            if (str.length > 10) {
                eventBus.emit('error', "You exceeded the number email of addresses")
                return
            }
            let response = await axiosInstance.post<ResponseT>("/templates/send-test-mails", sendEmailTestValues)
            if (response.data.status === true) {
                eventBus.emit('success', "Test mails sent successfully")
            }
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            get().setEmailTestValues({ template_id: "", email_address: "", subject: "" })
        }
    },

    resetForm: () => set({
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
        }
    }),

    searchMarketing: async (query) => {
        const { getAllMarketingTemplates } = get()
        if (!query) {
            await getAllMarketingTemplates()
        }

        await getAllMarketingTemplates(query)

    },
    searchTransactional: async (query) => {
        const { getAllTransactionalTemplates } = get()

        if (!query) {
            await getAllTransactionalTemplates()
        }

        await getAllTransactionalTemplates(query)
    }
}));

export default useTemplateStore;
