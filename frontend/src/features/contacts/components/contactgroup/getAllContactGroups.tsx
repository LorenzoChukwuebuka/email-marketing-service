import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import EmptyState from "../../../../components/emptyStateComponent";
import EditGroupComponent from "./editGroupComponent";
import useContactGroupStore from "../../store/contactgroup.store";
import { ContactGroupData } from '../../interface/contactgroup.interface';
import { APIResponse } from '../../../../../../frontend/src/interface/api.interface';
import { PaginatedResponse } from '../../../../../../frontend/src/interface/pagination.interface';
import LoadingSpinnerComponent from "../../../../components/loadingSpinnerComponent";
import { Pagination } from 'antd';

interface Props {
    contactgroupData?: APIResponse<PaginatedResponse<ContactGroupData>>
    loading: boolean
    onPageChange: (currentPage: number, pageSize: number) => void
    currentPage: number
    pageSize: number
}

const GetAllContactGroups: React.FC<Props> = ({ contactgroupData, loading, onPageChange, currentPage, pageSize }) => {

    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [selectedGroup, setSelectedGroup] = useState<ContactGroupData | null>(null);

    const { selectedGroupIds, setSelectedGroupIds } = useContactGroupStore()

    const cgData = useMemo(() => contactgroupData?.payload?.data || [], [contactgroupData])

    const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.checked) {
            const allIds = cgData?.map((contact) => contact.uuid);
            setSelectedGroupIds(allIds as string[]);
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

    return (
        <>
            {loading ? <LoadingSpinnerComponent /> : (<>

                <div className="overflow-x-auto mt-8">
                    <>
                        {cgData && cgData.length > 0 ? (
                            <>
                                <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                                    <thead className="bg-gray-50">
                                        <tr>
                                            <th className="py-3 px-4 text-left">
                                                <input
                                                    type="checkbox"
                                                    className="form-checkbox h-4 w-4 text-blue-600"
                                                    onChange={handleSelectAll}
                                                    checked={selectedGroupIds.length === (cgData?.length ?? 0)}
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
                                        {cgData.map((group: ContactGroupData) => (
                                            <tr key={group.uuid} className="hover:bg-gray-100">
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
                                                    <Link to="/app/contacts/view-group" state={{ groupId: group.uuid }}>
                                                        <i className="bi bi-eye"></i>
                                                    </Link>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>


                                <div className="mt-4 flex justify-center items-center mb-5">
                                    {/* Pagination */}
                                    <Pagination
                                        current={currentPage}
                                        pageSize={pageSize}
                                        total={contactgroupData?.payload?.total_count || 0} // Ensure your API returns a total count
                                        onChange={onPageChange}
                                        showSizeChanger
                                        pageSizeOptions={["10", "20", "50", "100"]}
                                        showTotal={(total) => `Total ${total} groups`}
                                    />
                                </div>
                            </>
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
            </>)}

            <EditGroupComponent isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} group={selectedGroup as ContactGroupData} />
        </>
    );
};

export default GetAllContactGroups;
