import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { BaseEntity } from '../../interface/baseentity.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { errResponse } from '../../utils/isError';

type ContactFormValues = {
    first_name: string
    last_name: string
    email: string
    from: string
    is_subscribed: boolean
}

type Group = {
    group_name: string;
    user_id: string;
    description: string;
} & BaseEntity;

type ContactBase = {
    user_id: string;
    groups: Group[] | null;
}

export type Contact = BaseEntity & ContactFormValues & ContactBase;

type ContactCount = { recent: number; total: number }

type EditContactValues = { uuid: string } & Partial<ContactFormValues>

type ContactEngageCount = {
    engaged: number
    new: number
    total: number
    unsubscribed: number
}

type FileCSVType = null | File;

interface ContactStore {
    contactFormValues: ContactFormValues;
    contactData: Contact[];
    selectedIds: string[];
    isLoading: boolean;
    contactCount: ContactCount
    engagementCount: ContactEngageCount
    selectedCSVFile: FileCSVType;
    editContactValues: EditContactValues;
    paginationInfo: Omit<PaginatedResponse<Contact>, 'data'>;
    setContactData: (newContactData: Contact[]) => void;
    setContactCount: (newData: ContactCount) => void;
    setContactFormValues: (newFormValues: ContactFormValues) => void;
    setIsLoading: (newIsLoading: boolean) => void;
    setSelectedId: (newSelectedId: string[]) => void;
    setSelectedCSVFile: (newSelectedFile: FileCSVType) => void;
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<Contact>, 'data'>) => void;
    setEditContactValues: (newEditContactValues: EditContactValues) => void;
    setEngagementCount: (newData: ContactEngageCount) => void
    createContact: () => Promise<void>;
    deleteContact: () => Promise<void>;
    editContact: () => Promise<void>;
    getAllContacts: (page?: number, pageSize?: number, search?: string) => Promise<void>;
    batchContactUpload: () => Promise<void>
    searchContacts: (query: string) => void;
    getContactCount: () => Promise<void>
    getContactSubEngagement: () => Promise<void>
}

type ContactsAPIResponse = APIResponse<PaginatedResponse<Contact>>;

const useContactStore = create<ContactStore>((set, get) => ({
    contactFormValues: {
        first_name: '',
        last_name: '',
        email: '',
        from: '',
        is_subscribed: false
    },
    contactCount: { recent: 0, total: 0 },
    contactData: [],
    selectedIds: [],
    engagementCount: {
        engaged: 0,
        new: 0,
        total: 0,
        unsubscribed: 0
    },
    paginationInfo: {
        total_count: 0,
        total_pages: 0,
        current_page: 1,
        page_size: 10,
    },
    isLoading: false,
    editContactValues: {
        uuid: '',
        first_name: '',
        last_name: '',
        email: '',
        from: '',
        is_subscribed: false
    },
    selectedCSVFile: null,
    setEditContactValues: (newEditContactValues) => set({ editContactValues: newEditContactValues }),
    setContactData: (newContactData) => set({ contactData: newContactData }),
    setContactFormValues: (newFormValues) => set({ contactFormValues: newFormValues }),
    setSelectedId: (newSelectedId) => set({ selectedIds: newSelectedId }),
    setPaginationInfo: (newPaginationInfo) => set({ paginationInfo: newPaginationInfo }),
    setIsLoading: (newIsLoading) => set({ isLoading: newIsLoading }),
    setSelectedCSVFile: (newSelectedFile) => set({ selectedCSVFile: newSelectedFile }),
    setContactCount: (newData) => set({ contactCount: newData }),
    setEngagementCount: (newData) => set({ engagementCount: newData }),


    createContact: async () => {
        try {
            const { setIsLoading, contactFormValues } = get()
            setIsLoading(true)
            const response = await axiosInstance.post<ResponseT>('/contact/create-contact', contactFormValues)
            if (response.data.message == "success") {
                eventBus.emit('success', "Contact created successfully")
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
            get().setIsLoading(false)
        }
    },
    deleteContact: async () => {
        try {
            const { selectedIds } = get()
            if (selectedIds.length > 0) {
                let promises = selectedIds.map((contactId) => {
                    return axiosInstance.delete<ResponseT>(
                        '/contact/delete-contact/' + contactId
                    )
                })

                const results = await Promise.all(promises)

                const allSuccessful = results.every(response => response.data.status === true)

                if (allSuccessful) {
                    eventBus.emit('success', "Group(s) deleted successfully")
                    await get().getAllContacts()
                } else {
                    eventBus.emit('error', "Some groups  could not be deleted")
                }
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
            get().setSelectedId([])
        }

    },
    editContact: async () => {
        try {
            const { editContactValues } = get()

            console.log(editContactValues)
            let response = await axiosInstance.put<ResponseT>("/contact/update-contact/" + editContactValues.uuid, editContactValues)
            if (response.data.status == true) {
                eventBus.emit('success', "Contact edited successfully")
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
    getAllContacts: async (page = 1, pageSize = 10, query = "") => {
        try {
            const response = await axiosInstance.get<ContactsAPIResponse>('/contact/get-all-contacts', {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            });
            const { data, ...paginationInfo } = response.data.payload;
            get().setContactData(data);
            get().setPaginationInfo(paginationInfo);
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },


    searchContacts: async (query: string) => {
        const { getAllContacts } = get();
        if (!query) {
            await getAllContacts(); // If query is empty, reset to all contacts
            return;
        }

        await getAllContacts(1, 10, query); // Search for contacts with the provided query
    },


    batchContactUpload: async () => {
        try {
            const { setIsLoading, selectedCSVFile } = get()

            setIsLoading(true)

            let data = new FormData

            data.append('contacts_csv', selectedCSVFile as Blob)

            let response = await axiosInstance.post<ResponseT>('/contact/upload-contact-csv', data)

            if (response.data.status == true) {
                eventBus.emit('success', "contacts uploaded successfully")
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
            get().setIsLoading(false)
        }
    },

    getContactCount: async () => {
        try {
            let response = await axiosInstance.get("/contact/get-contact-count")
            get().setContactCount(response.data.payload)
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
    getContactSubEngagement: async () => {
        try {
            let response = await axiosInstance.get<APIResponse<ContactEngageCount>>("/contact/contact-engagement")
            get().setEngagementCount(response.data.payload)
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

export default useContactStore;

