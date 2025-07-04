import { Template, SendTestMailValues } from '../interface/email-templates.interface';
import { BaseEntity } from '../../../../../frontend/src/interface/baseentity.interface';
import { create } from 'zustand';
import { handleError } from '../../../utils/isError';
import { TemplateAPI } from '../api/email-template.api';
import eventBus from '../../../utils/eventbus';

type TemplateState = {
    formValues: Omit<Template, "user_id">;
    sendEmailTestValues: SendTestMailValues
    currentTemplate: (Template & BaseEntity) | null;
}

type TemplateActions = {
    setCurrentTemplate: (template: (Template & BaseEntity) | null) => void;
    setFormValues: (newFormValues: Omit<Template, "user_id">) => void;
    setEmailTestValues: (newData: SendTestMailValues) => void
}

type TemplateAsyncActions = {
    updateTemplate: (uuid: string, updatedTemplate: Template & BaseEntity) => Promise<void>
    deleteTemplate: (uuid: string) => Promise<void>
    createTemplate: () => Promise<void>
    sendTestMail: () => Promise<void>
}

const editorPathMap: Record<"drag-and-drop" | "html-editor" | "rich-text", number> = {
    "drag-and-drop": 1,
    "html-editor": 2,
    "rich-text": 3,
};


const InitialState: TemplateState = {
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

    currentTemplate: null,
}

type TemplateStore = TemplateState & TemplateActions & TemplateAsyncActions

const useTemplateStore = create<TemplateStore>((set, get) => ({
    ...InitialState,

    setFormValues: (newFormValues) => set({ formValues: newFormValues }),
    setCurrentTemplate: (template) => set({ currentTemplate: template }),
    setEmailTestValues: (newData) => set({ sendEmailTestValues: newData }),

    createTemplate: async () => {
        const { formValues } = get()
        const response = await TemplateAPI.createTemplate(formValues)

        try {
            const editorType = response.payload.editor_type;
            const templateType = response.payload.type;
            const templateId = response.payload.template_id;

            const typeParam = templateType === "marketing" ? "m" : "t";

            const editorCode = editorPathMap[editorType as keyof typeof editorPathMap];

            if (!editorCode) {
                console.error("Unknown editor type:", editorType);
                return;
            }

            const redirectUrl = `/editor/${editorCode}?type=${typeParam}&uuid=${templateId}`;
            window.location.href = redirectUrl;

        } catch (error) {
            handleError(error);
        }
    },

    updateTemplate: async (uuid, updatedTemplate: Template & BaseEntity) => {
        try {
            const response = await TemplateAPI.updateTemplate(uuid, updatedTemplate)
            if (response) {
                eventBus.emit('success', 'Template updated successfully')
            }
        } catch (error) {
            handleError(error)
        }
    },

    deleteTemplate: async (uuid) => {
        try {
            const response = await TemplateAPI.deleteTemplate(uuid)
            if (response) {
                eventBus.emit('success', 'Template deleted successfully')
            }
        } catch (error) {
            handleError(error)
        }
    },

    sendTestMail: async () => {
        const { sendEmailTestValues } = get()
        try {
            const address = sendEmailTestValues.email_address
            const str = address.split(",")
            if (str.length > 10) {
                eventBus.emit('error', "You exceeded the number email of addresses")
                return
            }
            const response = await TemplateAPI.sendTestMail(sendEmailTestValues)
            if (response.payload.status === true) {
                eventBus.emit('success', "Test mails sent successfully")
            }
        } catch (error) {
            handleError(error)
        }
    }

}))


export default useTemplateStore


