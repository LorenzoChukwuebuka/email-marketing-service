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

        const editorType = response.payload["editor-type"];
        const templateType = response.payload.type;
        const templateId = response.payload.templateId;
        try {
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
            handleError(error)
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


