import React, { useEffect, useState } from 'react';
import useAdminUserStore, { AdminUserDetails } from '../../../../../store/admin/AdminUser.store';
import Pagination from '../../../../../components/Pagination';
import useDebounce from '../../../../../hooks/useDebounce';
import { useNavigate } from 'react-router-dom';

const AllUsersTable: React.FC = () => {
    const { getAllUsers, userdetailsData, paginationInfo, searchUser } = useAdminUserStore()
    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query

    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay

    const navigate = useNavigate()

    const handleSearch = (query: string) => {
        setSearchQuery(query)
    }

    const handlePageChange = (newPage: number) => {
        getAllUsers(newPage, paginationInfo.page_size);
    };

    const handleToggle = async (userId: string, field: 'verified' | 'blocked') => {
        // Confirm the toggle action
        const isConfirmed = confirm(`Are you sure you want to toggle ${field} status for this user?`);

        if (!isConfirmed) return; // Exit if user cancels the confirmation

        // Get the current list of users
        const updatedUsers = (userdetailsData as AdminUserDetails[]).map(user => {
            if (user.uuid === userId) {
                // Toggle the selected field
                return { ...user, [field]: !user[field] };
            }
            return user;
        });

        // Update the store with the toggled data
        useAdminUserStore.setState({ userdetailsData: updatedUsers });

        // Make an API call to update the user's status on the backend
        try {
            if (field === 'verified') {
                await useAdminUserStore.getState().verifyUser(userId);
            } else if (field === 'blocked') {
                const isBlocked = updatedUsers.find(user => user.uuid === userId)?.blocked;
                if (isBlocked) {
                    await useAdminUserStore.getState().blockUser(userId);
                } else {
                    await useAdminUserStore.getState().unBlockUser(userId);
                }
            }
        } catch (error) {
            console.error(`Error toggling ${field} status: `, error);
            // Revert the state if the API call fails
            useAdminUserStore.getState().getAllUsers();
        }
    };

    useEffect(() => {
        const fetchUsers = async () => {
            await getAllUsers()
        }
        fetchUsers()
    }, [getAllUsers])


    // Effect for debounced search
    useEffect(() => {
        if (debouncedSearchQuery !== "") {
            searchUser(debouncedSearchQuery, "All");
        } else {
            getAllUsers(); // Reset to all contacts when search query is empty
        }
    }, [debouncedSearchQuery, searchUser, getAllUsers]);

    return (
        <>
            <div className="flex justify-between items-center rounded-md p-2 bg-white mt-10">
                <div className="ml-3">
                    <input
                        type="text"
                        placeholder="Search..."
                        className="bg-gray-100 px-4 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-300"
                        onChange={(e) => handleSearch(e.target.value)}
                    />
                </div>
            </div>
            <div className="overflow-x-auto mt-4">
                <h1 className="font-semibold text-lg mt-4 mb-4"> User List </h1>
                <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Full Name
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Email
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Company
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Phone Number
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Verified
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Blocked
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Created At
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Verified At
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"> Details</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                        {Array.isArray(userdetailsData) && userdetailsData && userdetailsData.length > 0 ? (
                            userdetailsData.map((user, index) => {
                                return (
                                    <tr key={index} className="hover:bg-gray-100">
                                        <td className="py-4 px-4">{user?.fullname}</td>
                                        <td className="py-4 px-4">{user?.email}</td>
                                        <td className="py-4 px-4">{user?.company || 'N/A'}</td>
                                        <td className="py-4 px-4">{user?.phonenumber || 'N/A'}</td>
                                        <td className="py-4 px-4">
                                            <label className="inline-flex items-center">
                                                <input
                                                    type="checkbox"
                                                    className="form-checkbox h-5 w-5 text-blue-600"
                                                    checked={user?.verified}
                                                    onChange={() => handleToggle(user.uuid, 'verified')}
                                                />
                                                <span className="ml-2">{user?.verified ? 'Verified' : 'Not Verified'}</span>
                                            </label>
                                        </td>
                                        <td className="py-4 px-4">
                                            <label className="inline-flex items-center">
                                                <input
                                                    type="checkbox"
                                                    className="form-checkbox h-5 w-5 text-blue-600"
                                                    checked={user?.blocked}
                                                    onChange={() => handleToggle(user.uuid, 'blocked')}
                                                />
                                                <span className="ml-2">{user?.blocked ? 'Blocked' : 'Not Blocked'}</span>
                                            </label>
                                        </td>
                                        <td className="py-4 px-4">{user.created_at && new Date(user.created_at).toLocaleString('en-US', {
                                            timeZone: 'UTC',
                                            year: 'numeric',
                                            month: 'long',
                                            day: 'numeric',
                                            hour: 'numeric',
                                            minute: 'numeric',
                                            second: 'numeric'
                                        }) || "Not available"}</td>
                                        <td className="py-4 px-4">{user.verified_at && new Date(user.verified_at).toLocaleString('en-US', {
                                            timeZone: 'UTC',
                                            year: 'numeric',
                                            month: 'long',
                                            day: 'numeric',
                                            hour: 'numeric',
                                            minute: 'numeric',
                                            second: 'numeric'
                                        }) || "Not verified"}</td>
                                        <td className="py-4 px-4" onClick={() => navigate("/zen/dash/users/detail/" + user.uuid)}>   <i className="bi bi-eye"></i></td>
                                    </tr>
                                );
                            })
                        ) : (
                            <tr>
                                <td colSpan={8} className="py-4 px-4 text-center">
                                    No user data available
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>

                <Pagination paginationInfo={paginationInfo} handlePageChange={handlePageChange} item={"All Users"} />
            </div>
        </>
    )
}

export default AllUsersTable;
