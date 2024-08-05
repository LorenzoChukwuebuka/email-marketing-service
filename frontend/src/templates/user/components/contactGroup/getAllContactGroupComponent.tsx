import { useEffect, useState } from "react";
import { Contact } from "../../../../store/userstore/contactStore";
import useContactGroupStore, { ContactGroupData } from "../../../../store/userstore/contactGroupStore";
import { Link } from "react-router-dom";

const GetAllContactGroups: React.FC = () => {

    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [selectedGroup, setSelectedGroup] = useState<Contact | null>(null);

    const { contactgroupData, selectedIds, setSelectedIds, getAllGroups } = useContactGroupStore()

    const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.checked) {
            const allIds = (contactgroupData as ContactGroupData[]).map((contact) => contact.uuid);
            setSelectedIds(allIds);
        } else {
            setSelectedIds([]);
        }
    };

    const handleSelect = (uuid: string) => {
        if (selectedIds.includes(uuid)) {
            setSelectedIds(selectedIds.filter((id) => id !== uuid));
        } else {
            setSelectedIds([...selectedIds, uuid]);
        }
    }

    const openEditModal = (group: ContactGroupData) => { }

    useEffect(() => {
        getAllGroups();
    }, [getAllGroups]);

    return (
        <>
            <div className="overflow-x-auto mt-8">
                <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="py-3 px-4 text-left">
                                <input
                                    type="checkbox"
                                    className="form-checkbox h-4 w-4 text-blue-600"
                                    onChange={handleSelectAll}
                                    checked={selectedIds.length === ((contactgroupData as ContactGroupData[])?.length ?? 0)}
                                />
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Name
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Description
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                View All Contacts
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Created At
                            </th>
                            <th className="py-3 px-4">
                                Edit
                            </th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200 ">
                        {contactgroupData && (contactgroupData as ContactGroupData[]).length > 0 ? (
                            (contactgroupData as ContactGroupData[]).map((group: ContactGroupData) => (
                                <tr key={group.uuid}>

                                    <td className="py-4 px-4">
                                        <input
                                            type="checkbox"
                                            className="form-checkbox h-4 w-4 text-blue-600"
                                            checked={selectedIds.includes(group.uuid)}
                                            onChange={() => handleSelect(group.uuid)}
                                        />
                                    </td>
                                    <td className="py-4 px-4">{group.group_name}</td>
                                    <td className="py-4 px-4">{group.description}</td>
                                    <td className="py-4 px-4 text-center cursor-pointer">    <Link to="/user/dash/view-group" state={{ groupId: group.uuid }}> <i className="bi bi-eye"></i>  </Link> </td>
                                    <td className="py-4 px-4">{

                                        new Date(group.created_at).toLocaleString('en-US', {
                                            timeZone: 'UTC',
                                            year: 'numeric',
                                            month: 'long',
                                            day: 'numeric',
                                            hour: 'numeric',
                                            minute: 'numeric',
                                            second: 'numeric'
                                        })}</td>
                                    <td className="py-4 px-4">
                                        <button
                                            className="text-gray-400 hover:text-gray-600"
                                            onClick={() => openEditModal(group)}
                                        >
                                            ✏️
                                        </button>
                                    </td>
                                </tr>
                            ))
                        ) :
                            (
                                <tr>
                                    <td colSpan={7} className="py-4 px-4  text-center">
                                        No contacts available
                                    </td>
                                </tr>
                            )
                        }
                    </tbody>
                </table>
            </div>

            {/* Pagination component */}
            {/* EditContact modal */}
        </>
    );
};

export default GetAllContactGroups;
