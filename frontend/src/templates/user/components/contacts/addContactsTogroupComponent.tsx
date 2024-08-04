import { FormEvent, useEffect, useState } from "react";
import { Modal } from "../../../../components";
import useContactGroupStore, { ContactGroupData } from "../../../../store/userstore/contactGroupStore";

interface CGProps {
    isOpen: boolean;
    onClose: () => void;
}

const AddContactsToGroupComponent: React.FC<CGProps> = ({ isOpen, onClose }) => {

    const { getAllGroups, setSelectedIds, contactgroupData } = useContactGroupStore()
    const [checkedGroups, setCheckedGroups] = useState<Set<string>>(new Set());

    const handleGroupSubmit = (e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault()
        setSelectedIds(Array.from(checkedGroups))
    };

    const handleGroupCheck = (uuid: string) => {
        setCheckedGroups(prev => {
            const newSet = new Set(prev);
            if (newSet.has(uuid)) {
                newSet.delete(uuid);
            } else {
                newSet.add(uuid);
            }
            return newSet;
        });
    };


    useEffect(() => {
        getAllGroups();
    }, [getAllGroups]);


    return <>

        <Modal isOpen={isOpen} onClose={onClose} title="Add selected Contact(s) to Group">
            <>

                <div className="relative mb-4">
                    <input
                        type="text"
                        placeholder="Search"
                        className="w-full pl-8 pr-4 py-2 border rounded-lg"
                    // value={searchTerm}
                    // onChange={(e) => setSearchTerm(e.target.value)}
                    />
                    <svg className="w-5 h-5 text-gray-400 absolute left-2 top-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                </div>


                <div className="max-h-60 overflow-y-auto">
                    {contactgroupData.map((group: ContactGroupData) => (
                        <div key={group.uuid} className="mb-4">
                            <label className="flex items-center space-x-2">
                                <input
                                    type="checkbox"
                                    className="form-checkbox text-blue-600"
                                    checked={checkedGroups.has(group.uuid)}
                                    onChange={() => handleGroupCheck(group.uuid)}
                                />
                                <span className="font-semibold">{group.group_name} ({group.contacts ? group.contacts.length : 0})</span>
                            </label>
                            <p className="text-sm text-gray-500 ml-6">{group.description}</p>
                        </div>
                    ))}
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
                        type="submit"
                        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                        onClick={(e) => handleGroupSubmit(e)}
                    >
                        Submit
                    </button>
                </div>
            </>
        </Modal>

    </>
}

export default AddContactsToGroupComponent