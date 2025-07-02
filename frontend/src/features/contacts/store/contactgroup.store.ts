import { create } from 'zustand';
import { ContactGroupFormValues, EditGroupValues, AddToGroup } from '../interface/contactgroup.interface';
import { handleError } from '../../../utils/isError';
import { contactGroupAPI } from '../api/contactgroup.api';
import eventBus from '../../../utils/eventbus';
import { invalidateContactGroupQuery } from '../hooks/useContactGroupQuery';
import { queryClient } from '../../../utils/queryclient';

// State Types
type ContactGroupState = {
    selectedContactIds: string[];
    selectedGroupIds: string[];
    formValues: ContactGroupFormValues;
    editValues: EditGroupValues;
}

// Action Types
type ContactGroupActions = {
    setSelectedContactIds: (newId: string[]) => void;
    setFormValues: (newformValue: ContactGroupFormValues) => void;
    setEditValues: (newEditValues: EditGroupValues) => void;
    setSelectedGroupIds: (newIds: string[]) => void;
}

// Async Action Types
type ContactGroupAsyncActions = {
    addContactToGroup: () => Promise<void>;
    removeContactFromGroup: () => Promise<void>;
    createGroup: () => Promise<void>;
    deleteGroup: () => Promise<void>;
    updateGroup: () => Promise<void>;

}

// Combined Store Type
type ContactGroupStore = ContactGroupState & ContactGroupActions & ContactGroupAsyncActions;

// Initial State
const InitialState: ContactGroupState = {
    selectedContactIds: [],
    selectedGroupIds: [],
    formValues: { group_name: '', description: '' },
    editValues: { group_id: '', group_name: '', description: '' },

};

const useContactGroupStore = create<ContactGroupStore>((set, get) => ({
    ...InitialState,
    // Actions
    setSelectedContactIds: (newId) => set({ selectedContactIds: newId }),
    setSelectedGroupIds: (newIds) => set({ selectedGroupIds: newIds }),
    setFormValues: (newformValue) => set({ formValues: newformValue }),
    setEditValues: (newEditValues) => set({ editValues: newEditValues }),

    addContactToGroup: async () => {
        const { selectedContactIds, selectedGroupIds, setSelectedGroupIds } = get();
        try {
            if (selectedGroupIds.length > 0) {
                const groupId = selectedGroupIds[0]
                // Clear selectedGroupIds immediately after using it
                setSelectedGroupIds([])
                const promises = selectedContactIds.map(contactId => {
                    const data = {
                        group_id: groupId,
                        contact_id: contactId
                    } satisfies AddToGroup

                    return contactGroupAPI.addContactToGroup(data)
                })

                const results = await Promise.all(promises)
                results.every(response => response.payload.status === true)
                eventBus.emit('success', "Contacts added to group successfully")
            }
        } catch (error) {
            handleError(error)
        }
    },


    createGroup: async () => {
        const { formValues } = get();
        try {
            const response = await contactGroupAPI.createGroup(formValues)
            if (response) {
                eventBus.emit('success', 'Group has been created successfully')
            }
            invalidateContactGroupQuery(queryClient)
        } catch (error) {
            handleError(error)
        }
    },

    updateGroup: async () => {
        const { editValues } = get();
        try {
            const response = await contactGroupAPI.updateGroup(editValues.group_id, editValues)
            if (response) {
                eventBus.emit('success', 'Group edited successfully')
            }
            invalidateContactGroupQuery(queryClient)
        } catch (error) {
            handleError(error)
        }
    },

    deleteGroup: async () => {
        const { selectedGroupIds } = get();
        try {
            if (selectedGroupIds.length > 0) {
                const promises = selectedGroupIds.map((groupId) => {
                    return contactGroupAPI.deleteGroup(groupId)
                })

                const results = await Promise.all(promises)
                results.every(response => response.payload.status === true)
                eventBus.emit('success', "Group(s) deleted successfully")
            }
            invalidateContactGroupQuery(queryClient)
        } catch (error) {
            handleError(error)
        }
    },

    removeContactFromGroup: async () => {
        const { selectedContactIds, selectedGroupIds, setSelectedGroupIds } = get();
        try {
            if (selectedGroupIds.length > 0) {
                const groupId = selectedGroupIds[0]
                // Clear selectedGroupIds immediately after using it
                setSelectedGroupIds([])
                const promises = selectedContactIds.map(contactId => {
                    const data = {
                        group_id: groupId,
                        contact_id: contactId
                    } satisfies AddToGroup

                    return contactGroupAPI.removeContactFromGroup(data.group_id, data.contact_id)
                })

                const results = await Promise.all(promises)

                results.every(response => response.payload.status === true)
                eventBus.emit('success', 'contact has been removed successfully')
            }

        } catch (error) {
            handleError(error)
        }
    },


}));

export default useContactGroupStore;