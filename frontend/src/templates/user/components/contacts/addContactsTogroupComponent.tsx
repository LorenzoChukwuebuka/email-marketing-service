import { useEffect, useState } from "react";
import { Modal } from "../../../../components";
import useContactGroupStore, { ContactGroupData } from "../../../../store/userstore/contactGroupStore";

interface CGProps {
    isOpen: boolean;
    onClose: () => void;
}

const AddContactsToGroupComponent: React.FC<CGProps> = ({ isOpen, onClose }) => {
    const { getAllGroups, setSelectedGroupIds, selectedGroupIds, contactgroupData, addContactToGroup } = useContactGroupStore();

    const [selectedGroup, setSelectedGroup] = useState<string | null>(null);

    const handleGroupSubmit = (e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();

        console.log(selectedGroup, 'component state')
        console.log(selectedGroupIds, 'store state but from the component')
        if (selectedGroup) {
            setSelectedGroupIds([selectedGroup]);
            addContactToGroup();
            setSelectedGroup(null);
            onClose();
        }
    };

    const handleGroupSelect = (uuid: string) => {
        setSelectedGroup(uuid);
    };

    useEffect(() => {
        getAllGroups();
    }, [getAllGroups]);

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Add selected Contact(s) to Group">
            <>
                {/* Search input remains the same */}

                <div className="max-h-60 overflow-y-auto">
                    {contactgroupData && (contactgroupData as ContactGroupData[]).length > 0 && (
                        (contactgroupData as ContactGroupData[]).map((group: ContactGroupData) => (
                            <div key={group.uuid} className="mb-4">
                                <label className="flex items-center space-x-2">
                                    <input
                                        type="radio"
                                        name="group"
                                        className="form-radio text-blue-600"
                                        checked={selectedGroup === group.uuid}
                                        onChange={() => handleGroupSelect(group.uuid)}
                                    />
                                    <span className="font-semibold">{group.group_name} ({group.contacts ? group.contacts.length : 0})</span>
                                </label>
                                <p className="text-sm text-gray-500 ml-6">{group.description}</p>
                            </div>
                        ))
                    )}

                </div>

                <div className="flex justify-end space-x-2">
                    <button
                        type="button"
                        onClick={onClose}
                        className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
                    >
                        Cancel
                    </button>
                    <button
                        type="button"
                        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                        onClick={handleGroupSubmit}
                        disabled={!selectedGroup}
                    >
                        Submit
                    </button>
                </div>
            </>
        </Modal>
    );
}

export default AddContactsToGroupComponent;
