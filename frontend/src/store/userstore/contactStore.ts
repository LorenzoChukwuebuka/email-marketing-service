import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { BaseEntity } from '../../interface/baseentity.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse } from '../../interface/api.interface';

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

type Contact = BaseEntity & ContactFormValues & ContactBase;

interface ContactStore {
    contactFormValues: ContactFormValues;
    contactData: Contact[];
    selectedIds: string[];
    isLoading: boolean;
    paginationInfo: Omit<PaginatedResponse<Contact>, 'data'>;
    setContactData: (newContactData: Contact[]) => void;
    setContactFormValues: (newFormValues: ContactFormValues) => void;
    setIsLoading: (newIsLoading: boolean) => void;
    setSelectedId: (newSelectedId: string[]) => void;
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<Contact>, 'data'>) => void;
    createContact: () => Promise<void>;
    deleteContact: () => Promise<void>;
    editContact: () => Promise<void>;
    getAllContacts: (page?: number, pageSize?: number) => Promise<void>;
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
    setContactData: (newContactData) => set({ contactData: newContactData }),
    setContactFormValues: (newFormValues) => set({ contactFormValues: newFormValues }),
    setSelectedId: (newSelectedId) => set({ selectedIds: newSelectedId }),
    setPaginationInfo: (newPaginationInfo) => set({ paginationInfo: newPaginationInfo }),
    setIsLoading: (newIsLoading) => set({ isLoading: newIsLoading }),

    createContact: async () => {
        try {
            const { setIsLoading, contactFormValues } = get()
            setIsLoading(true)
            const response = await axiosInstance.post('/create-contact', contactFormValues)
            if (response.data.message == "success") {
                eventBus.emit('success', "Contact created successfully")
            }

        } catch (error) {
            eventBus.emit('error', error instanceof Error ? error.message : 'An unexpected error occurred');
        } finally {
            get().setIsLoading(false)
        }
    },
    deleteContact: async () => { /* implementation */ },
    editContact: async () => { /* implementation */ },
    getAllContacts: async (page = 1, pageSize = 10) => {
        try {
            const response = await axiosInstance.get<ContactsAPIResponse>(`/get-all-contacts?page=${page}&page_size=${pageSize}`);
            console.log(response.data)
            const { data, ...paginationInfo } = response.data.payload;
            get().setContactData(data);
            get().setPaginationInfo(paginationInfo);
        } catch (error) {
            eventBus.emit('error', error instanceof Error ? error.message : 'An unexpected error occurred');
        }
    }
}));

export default useContactStore;

