import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { BaseEntity } from '../../interface/baseentity.interface';
import { Contact } from './contactStore';
import { errResponse } from '../../utils/isError';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse } from '../../interface/api.interface';

type FormValues = {
    name: string;
    description: string
}

export type ContactGroupData = Pick<FormValues, 'description'> & BaseEntity & {
    userId: string;
    group_name: string;
    contacts: Omit<Contact, 'group'>[]
}

interface ContactGroupstore {
    isLoading: boolean
    selectedContactIds: string[],
    selectedIds: string[]
    paginationInfo: Omit<PaginatedResponse<ContactGroupData>, 'data'>;
    contactgroupData: ContactGroupData[]
    setIsLoading: (newIsLoading: boolean) => void;
    setContactGroupData: (newgroupData: ContactGroupData[]) => void
    setSelectedContactIds: (newId: string[]) => void
    getAllGroups: (page?: number, pageSize?: number) => Promise<void>;
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<ContactGroupData>, 'data'>) => void;
    setSelectedIds: (newIds: string[]) => void
}

const useContactGroupStore = create<ContactGroupstore>((set, get) => ({
    isLoading: false,
    contactgroupData: [],
    selectedContactIds: [],
    selectedIds: [],

    paginationInfo: {
        total_count: 0,
        total_pages: 0,
        current_page: 1,
        page_size: 10,
    },

    setIsLoading: (newIsLoading) => set({ isLoading: newIsLoading }),
    setContactGroupData: (newgroupData) => set({ contactgroupData: newgroupData }),
    setSelectedContactIds: (newId) => set({ selectedContactIds: newId }),
    setPaginationInfo: (newPaginationInfo) => set({ paginationInfo: newPaginationInfo }),
    setSelectedIds: (newId) => set({ selectedIds: newId }),

    getAllGroups: async (page = 1, pageSize = 10): Promise<void> => {
        try {
            const { setContactGroupData, setPaginationInfo } = get()
            let response = await axiosInstance.get<APIResponse<PaginatedResponse<ContactGroupData>>>(`/get-all-contact-groups?page=${page}&page_size=${pageSize}`)

            const { data, ...paginationInfo } = response.data.payload

            setContactGroupData(data)

            setPaginationInfo(paginationInfo)

        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.message)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    }
}))


export default useContactGroupStore