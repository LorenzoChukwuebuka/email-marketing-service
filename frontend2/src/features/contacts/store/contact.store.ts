import { create } from 'zustand';
import { ContactFormValues, EditContactValues, FileCSVType } from '../interface/contact.interface';
import contactApi from '../api/contact.api';
import eventBus from '../../../utils/eventbus';
import { handleError } from '../../../utils/isError';
import { invalidateContactQuery } from '../hooks/useContactQuery';
import { queryClient } from '../../../utils/queryclient';

type ContactState = {
    contactFormValues: ContactFormValues;
    selectedIds: string[];
    selectedCSVFile: FileCSVType;
    editContactValues: EditContactValues;
}

type ContactAsyncState = {
    createContact: () => Promise<void>;
    deleteContact: () => Promise<void>;
    editContact: () => Promise<void>;
    batchContactUpload: () => Promise<void>
}

type ContactActions = {
    setContactFormValues: (newFormValues: ContactFormValues) => void;
    setSelectedId: (newSelectedId: string[]) => void;
    setSelectedCSVFile: (newSelectedFile: FileCSVType) => void;
    setEditContactValues: (newEditContactValues: EditContactValues) => void;
}

type ContactStore = ContactState & ContactActions & ContactAsyncState

const InitializeState: ContactState = {
    contactFormValues: {
        first_name: '',
        last_name: '',
        email: '',
        from: '',
        is_subscribed: false,
    },
    selectedCSVFile: null,
    selectedIds: [],

    editContactValues: {
        uuid: '',
    }
}

const useContactStore = create<ContactStore>((set, get) => ({
    ...InitializeState,

    setContactFormValues: (newFormValues: ContactFormValues) => set({ contactFormValues: newFormValues }),
    setSelectedId: (newSelectedId: string[]) => set({ selectedIds: newSelectedId }),
    setSelectedCSVFile: (newSelectedFile: FileCSVType) => set({ selectedCSVFile: newSelectedFile }),
    setEditContactValues: (newEditContactValues: EditContactValues) => set({ editContactValues: newEditContactValues }),

    createContact: async () => {
        const { contactFormValues } = get()
        try {
            const response = await contactApi.createContact(contactFormValues)
            if (response) {
                eventBus.emit('success', 'Contact has been created successfully')
            }
            invalidateContactQuery(queryClient)
        } catch (error) {
            handleError(error)
        }
    },

    deleteContact: async () => {
        const { selectedIds } = get()

        try {
            if (selectedIds.length > 0) {
                const promises = selectedIds.map((contactId) => {
                    return contactApi.deleteContact(contactId)
                })

                const results = await Promise.all(promises)
                results.every(response => response.payload.status === true)

                eventBus.emit('success', "contact(s) deleted successfully")

            }
            invalidateContactQuery(queryClient)
        } catch (error) {
            handleError(error)
        }
    },

    editContact: async () => {
        const { editContactValues } = get()
        try {
            const response = await contactApi.editContact(editContactValues.uuid, editContactValues)
            if (response) {
                eventBus.emit('success', 'contact has been edited successfully')
            }
            invalidateContactQuery(queryClient)
        } catch (error) {
            handleError(error)
        }
    },

    batchContactUpload: async () => {
        const { selectedCSVFile } = get()
        try {
            const data = new FormData
            data.append('contacts_csv', selectedCSVFile as Blob)
            const response = await contactApi.batchUploadContact(data)
            if (response) {
                eventBus.emit('success', 'contacts upload successful')
            }
            invalidateContactQuery(queryClient)
        } catch (error) {
            handleError(error)
        }
    }

}))


export default useContactStore;