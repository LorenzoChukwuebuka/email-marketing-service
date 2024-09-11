import { useEffect, useState } from "react";
import useContactGroupStore, { ContactGroupData } from "../../../../store/userstore/contactGroupStore";
import { Link } from "react-router-dom";
import EditGroupComponent from "./editGroupComponent";
import EmptyState from "../../../../components/emptyStateComponent";


const GetAllContactGroups: React.FC = () => {

    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [selectedGroup, setSelectedGroup] = useState<ContactGroupData | null>(null);

    const { contactgroupData, selectedGroupIds, setSelectedGroupIds, getAllGroups } = useContactGroupStore()

    const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.checked) {
            const allIds = (contactgroupData as ContactGroupData[]).map((contact) => contact.uuid);
            setSelectedGroupIds(allIds);
        } else {
            setSelectedGroupIds([]);
        }
    };

    const handleSelect = (uuid: string) => {
        if (selectedGroupIds.includes(uuid)) {
            setSelectedGroupIds(selectedGroupIds.filter((id) => id !== uuid));
        } else {
            setSelectedGroupIds([...selectedGroupIds, uuid]);
        }
    }

    const openEditModal = (group: ContactGroupData) => {
        setIsModalOpen(true)
        setSelectedGroup(group)
    }

    useEffect(() => {
        getAllGroups();
    }, [getAllGroups]);

    return (
        <>
            <div className="overflow-x-auto mt-8">
                <>
                    {contactgroupData && (contactgroupData as ContactGroupData[]).length > 0 ? (
                        <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="py-3 px-4 text-left">
                                        <input
                                            type="checkbox"
                                            className="form-checkbox h-4 w-4 text-blue-600"
                                            onChange={handleSelectAll}
                                            checked={selectedGroupIds.length === ((contactgroupData as ContactGroupData[])?.length ?? 0)}
                                        />
                                    </th>
                                    <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Name
                                    </th>
                                    <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Description
                                    </th>
                                    <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Contacts
                                    </th>
                                    <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Created At
                                    </th>
                                    <th className="py-3 px-4">
                                        Edit
                                    </th>
                                    <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        View All Contacts
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200">
                                {(contactgroupData as ContactGroupData[]).map((group: ContactGroupData) => (
                                    <tr key={group.uuid}>
                                        <td className="py-4 px-4">
                                            <input
                                                type="checkbox"
                                                className="form-checkbox h-4 w-4 text-blue-600"
                                                checked={selectedGroupIds.includes(group.uuid)}
                                                onChange={() => handleSelect(group.uuid)}
                                            />
                                        </td>
                                        <td className="py-4 px-4">{group.group_name}</td>
                                        <td className="py-4 px-4">{group.description}</td>
                                        <td className="py-4 px-4">{group.contacts ? group.contacts.length : 0}</td>
                                        <td className="py-4 px-4">
                                            {new Date(group.created_at).toLocaleString('en-US', {
                                                timeZone: 'UTC',
                                                year: 'numeric',
                                                month: 'long',
                                                day: 'numeric',
                                                hour: 'numeric',
                                                minute: 'numeric',
                                                second: 'numeric',
                                            })}
                                        </td>
                                        <td className="py-4 px-4">
                                            <button
                                                className="text-gray-400 hover:text-gray-600"
                                                onClick={() => openEditModal(group)}
                                            >
                                                ✏️
                                            </button>
                                        </td>
                                        <td className="py-4 px-4 text-center cursor-pointer">
                                            <Link to="/user/dash/view-group" state={{ groupId: group.uuid }}>
                                                <i className="bi bi-eye"></i>
                                            </Link>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    ) : (
                        <div className="py-4 px-4 text-center">
                            <EmptyState
                                title="You have not created any groups"
                                description="Create groups to easily manage your contacts"
                                icon={<i className="bi bi-emoji-frown text-xl"></i>}
                            />
                        </div>
                    )}

                </>
            </div>

            <EditGroupComponent isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} group={selectedGroup} />
        </>
    );
};

export default GetAllContactGroups;
