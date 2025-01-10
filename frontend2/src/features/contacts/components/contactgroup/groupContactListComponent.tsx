import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import useContactGroupStore from "../../store/contactgroup.store";
import { ContactGroupData } from '../../interface/contactgroup.interface';
import { useSingleContactGroupQuery } from "../../hooks/useContactGroupQuery";
import LoadingSpinnerComponent from "../../../../components/loadingSpinnerComponent";


const GroupContactList: React.FC = () => {
    const { setSelectedContactIds, setSelectedGroupIds, selectedContactIds, removeContactFromGroup } = useContactGroupStore();
    const [group, setGroup] = useState<ContactGroupData | null>(null);
    const [selectedIds, setSelectedIds] = useState<string[]>([]);




    const location = useLocation();
    const stateData = location.state as { groupId: string };

    const { data: contactgroupData, isLoading } = useSingleContactGroupQuery(stateData.groupId)

    const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.checked) {
            const allIds = group?.contacts?.map((contact) => contact.uuid) || [];
            setSelectedIds(allIds);
            setSelectedContactIds(allIds);
        } else {
            setSelectedIds([]);
            setSelectedContactIds([]);
        }
    };

    const handleSelect = (uuid: string) => {
        const updatedIds = selectedIds.includes(uuid)
            ? selectedIds.filter((id) => id !== uuid)
            : [...selectedIds, uuid];
        setSelectedIds(updatedIds);
        setSelectedContactIds(updatedIds);
    };

    const handleSubmit = async (e: React.MouseEvent<HTMLButtonElement>) => {
        try {
            e.preventDefault()
            const stateData = location.state as { groupId: string };
            setSelectedGroupIds([stateData.groupId])
            await removeContactFromGroup()
        } catch (error) {
            console.log(error)
        }

    }



    useEffect(() => {
        if (contactgroupData) {
            if (Array.isArray(contactgroupData) && contactgroupData.length > 0) {
                setGroup(contactgroupData[0]);
            } else if (!Array.isArray(contactgroupData)) {
                setGroup(contactgroupData);
            }
        }
    }, [contactgroupData]);

    return (
        <main className="mt-5">
            <div className="p-6">
                {isLoading ? (
                    <LoadingSpinnerComponent />
                ) : group ? (
                    <>
                        <div className="flex items-start flex-col mb-5">
                            <button className="text-blue-600 mr-2" onClick={() => window.history.back()}>
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                                </svg>
                            </button>
                            <h2 className="text-xl font-bold text-gray-800">Group Name: {group.group_name}</h2>
                            <p className="text-md block text-gray-500 p-2 px-2 mt-2">Description: {group.description}</p>
                        </div>

                        <h1 className="text-xl font-semibold mt-5">Contacts</h1>

                        <div className="flex justify-between items-center rounded-md p-2 bg-white mt-10">
                            <div className="space-x-1  h-auto w-full p-2 px-2 ">
                                {selectedContactIds.length > 0 && (
                                    <>
                                        <button
                                            className="bg-red-200 px-4 py-2 rounded-md transition duration-300"
                                            onClick={(e) => handleSubmit(e)}
                                        >
                                            <span className="text-red-500"> Remove Contact(s)  </span>
                                            <i className="bi bi-trash text-red-500"></i>
                                        </button>
                                    </>
                                )}
                            </div>

                            <div className="ml-3">
                                <input
                                    type="text"
                                    placeholder="Search..."
                                    className="bg-gray-100 px-4 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-300"
                                // onChange={(e) => handleSearch(e.target.value)}
                                />
                            </div>

                        </div>

                        <table className="md:min-w-5xl min-w-full w-full mt-5 rounded-sm bg-white">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="py-3 px-4 text-left">
                                        <input
                                            type="checkbox"
                                            className="form-checkbox h-4 w-4 text-blue-600"
                                            onChange={handleSelectAll}
                                            checked={selectedIds.length === (group.contacts?.length ?? 0)}
                                        />
                                    </th>
                                    <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        First Name
                                    </th>
                                    <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Last Name
                                    </th>
                                    <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Email
                                    </th>
                                    <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        From
                                    </th>
                                    <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Created At
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200">
                                {group.contacts && group.contacts.length > 0 ? (
                                    group.contacts.map((contact) => (
                                        <tr key={contact.uuid}>
                                            <td className="py-4 px-4">
                                                <input
                                                    type="checkbox"
                                                    className="form-checkbox h-4 w-4 text-blue-600"
                                                    checked={selectedIds.includes(contact.uuid)}
                                                    onChange={() => handleSelect(contact.uuid)}
                                                />
                                            </td>
                                            <td className="py-4 px-4">{contact.first_name}</td>
                                            <td className="py-4 px-4">{contact.last_name}</td>
                                            <td className="py-4 px-4">{contact.email}</td>
                                            <td className="py-4 px-4">{contact.from}</td>
                                            <td className="py-4 px-4">
                                                {new Date(contact.created_at).toLocaleString('en-US', {
                                                    timeZone: 'UTC',
                                                    year: 'numeric',
                                                    month: 'long',
                                                    day: 'numeric',
                                                    hour: 'numeric',
                                                    minute: 'numeric',
                                                    second: 'numeric'
                                                })}
                                            </td>
                                        </tr>
                                    ))
                                ) : (
                                    <tr>
                                        <td colSpan={6} className="py-4 px-4 text-center">
                                            No contacts available
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </>
                ) : (
                    <p>No group data available.</p>
                )}
            </div>
        </main>
    );
};

export default GroupContactList;