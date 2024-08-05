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

type EditContactValues = { uuid: string } & Partial<ContactFormValues>


type FileCSVType = null | File;

interface ContactStore {
    contactFormValues: ContactFormValues;
    contactData: Contact[];
    selectedIds: string[];
    isLoading: boolean;
    selectedCSVFile: FileCSVType;
    editContactValues: EditContactValues;
    paginationInfo: Omit<PaginatedResponse<Contact>, 'data'>;
    setContactData: (newContactData: Contact[]) => void;
    setContactFormValues: (newFormValues: ContactFormValues) => void;
    setIsLoading: (newIsLoading: boolean) => void;
    setSelectedId: (newSelectedId: string[]) => void;
    setSelectedCSVFile: (newSelectedFile: FileCSVType) => void;
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<Contact>, 'data'>) => void;
    setEditContactValues: (newEditContactValues: EditContactValues) => void;
    createContact: () => Promise<void>;
    deleteContact: () => Promise<void>;
    editContact: () => Promise<void>;
    getAllContacts: (page?: number, pageSize?: number) => Promise<void>;
    batchContactUpload: () => Promise<void>
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
    contactData: [],
    selectedIds: [],
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

    createContact: async () => {
        try {
            const { setIsLoading, contactFormValues } = get()
            setIsLoading(true)
            const response = await axiosInstance.post<ResponseT>('/create-contact', contactFormValues)
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
                        '/delete-contact/' + contactId
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
            let response = await axiosInstance.put<ResponseT>("/update-contact/" + editContactValues.uuid, editContactValues)
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
    getAllContacts: async (page = 1, pageSize = 10) => {
        try {
            const response = await axiosInstance.get<ContactsAPIResponse>(`/get-all-contacts?page=${page}&page_size=${pageSize}`);
            const { data, ...paginationInfo } = response.data.payload;
            get().setContactData(data);
            get().setPaginationInfo(paginationInfo);
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

    batchContactUpload: async () => {
        try {
            const { setIsLoading, selectedCSVFile } = get()

            setIsLoading(true)

            let data = new FormData

            data.append('contacts_csv', selectedCSVFile as Blob)

            let response = await axiosInstance.post<ResponseT>('/upload-contact-csv', data)

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
    }
}));

export default useContactStore;

