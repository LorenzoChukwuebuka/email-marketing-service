import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { BaseEntity } from '../../interface/baseentity.interface';
import { Contact } from './contactStore';
import { errResponse } from '../../utils/isError';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse, ResponseT } from '../../interface/api.interface';

type FormValues = {
    group_name: string;
    description: string
}

type AddToGroup = {
    group_id: string;
    contact_id: string
}

export type ContactGroupData = FormValues & BaseEntity & {
    userId: string;
    contacts: Omit<Contact, 'group'>[]
}

type EditGroupValues = FormValues & { uuid: string }

interface ContactGroupstore {
    isLoading: boolean
    selectedContactIds: string[];
    formValues: FormValues;
    selectedIds: string[]
    paginationInfo: Omit<PaginatedResponse<ContactGroupData>, 'data'>;
    contactgroupData: ContactGroupData[] | ContactGroupData
    setIsLoading: (newIsLoading: boolean) => void;
    setContactGroupData: (newgroupData: ContactGroupData[] | ContactGroupData) => void
    setSelectedContactIds: (newId: string[]) => void
    setFormValues: (newformValue: FormValues) => void
    getAllGroups: (page?: number, pageSize?: number) => Promise<void>;
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<ContactGroupData>, 'data'>) => void;
    setSelectedIds: (newIds: string[]) => void;
    addContactToGroup: () => Promise<void>;
    getSingleGroup: (uuid: string) => Promise<void>
    createGroup: () => Promise<void>
    deleteGroup: () => Promise<void>
    updateGroup: () => Promise<void>
}

const useContactGroupStore = create<ContactGroupstore>((set, get) => ({
    isLoading: false,
    contactgroupData: [],
    selectedContactIds: [],
    selectedIds: [],
    formValues: { group_name: '', description: "" },

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
    setFormValues: (newformValue) => set({ formValues: newformValue }),

    getAllGroups: async (page = 1, pageSize = 10): Promise<void> => {
        try {
            const { setContactGroupData, setPaginationInfo } = get()
            let response = await axiosInstance.get<APIResponse<PaginatedResponse<ContactGroupData>>>(`/get-all-contact-groups?page=${page}&page_size=${pageSize}`)

            const { data, ...paginationInfo } = response.data.payload

            setContactGroupData(data)

            setPaginationInfo(paginationInfo)

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
    addContactToGroup: async () => {
        try {
            const { setIsLoading, selectedContactIds, selectedIds } = get()

            setIsLoading(true)

            if (selectedIds.length > 0) {
                const promises = selectedContactIds.map(contactId => {
                    const data = {
                        group_id: selectedIds[0],
                        contact_id: contactId
                    } satisfies AddToGroup

                    return axiosInstance.post<ResponseT>("/add-contact-to-group", data)
                })

                const results = await Promise.all(promises)

                const allSuccessful = results.every(response => response.data.status === true)

                if (allSuccessful) {
                    eventBus.emit('success', "Contacts added to group successfully")
                } else {
                    eventBus.emit('error', "Some contacts could not be added to the group")
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
            const { setIsLoading, setSelectedIds, setSelectedContactIds } = get()

            setIsLoading(false)
            setSelectedContactIds([])
            setSelectedIds([])
        }
    },

    getSingleGroup: async (uuid: string) => {
        try {
            const { setIsLoading, setContactGroupData } = get()
            setIsLoading(true)
            let response = await axiosInstance.get<APIResponse<ContactGroupData>>('/get-single-group/' + uuid)
            const groupData = response.data.payload
            setContactGroupData(Array.isArray(groupData) ? groupData : [groupData].filter(Boolean))

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

    createGroup: async () => {
        try {
            const { formValues, setIsLoading } = get()
            setIsLoading(true)
            let response = await axiosInstance.post<ResponseT>("/create-contact-group", formValues)
            if (response.data.status === true) {
                eventBus.emit('success', "Group created successfully")
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
            get().setFormValues({ group_name: "", description: "" })
        }
    },

    updateGroup: async () => { },
    deleteGroup: async () => { }
}))


export default useContactGroupStore