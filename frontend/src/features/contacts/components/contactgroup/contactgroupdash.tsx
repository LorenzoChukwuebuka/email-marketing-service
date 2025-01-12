import { useState } from "react";

import useDebounce from "../../../../hooks/useDebounce";
import useContactGroupStore from "../../store/contactgroup.store";
import CreateGroup from './createGroupComponent';
import { useContactGroupQuery } from "../../hooks/useContactGroupQuery";
import GetAllContactGroups from './getAllContactGroups';
import { ContactGroupData } from "../../interface/contactgroup.interface";
import { APIResponse } from "../../../../interface/api.interface";
import { PaginatedResponse } from "../../../../interface/pagination.interface";
import { Modal } from "antd";

const ContactGroupDash: React.FC = () => {

    const { selectedGroupIds, deleteGroup } = useContactGroupStore()
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query

    // Debounce the search query
    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);

    const { data: groupData, isLoading } = useContactGroupQuery(currentPage, pageSize, debouncedSearchQuery)

    const onPageChange = (page: number, size: number) => {
        setCurrentPage(page);
        setPageSize(size);
    };


    const handleSearch = (query: string) => {
        setSearchQuery(query)
    };

    const deleteGrp = async () => {
        Modal.confirm({
            title: "Are you sure?",
            content: "Do you want to delete groups(s)?",
            okText: "Yes",
            cancelText: "No",
            onOk: async () => {
                await deleteGroup();
                await new Promise(resolve => setTimeout(resolve, 3000));
                location.reload()
            },
        });
    }

    return <>

        <div className="flex justify-between items-center rounded-md p-2 bg-white mt-10">
            <div className="space-x-1  h-auto w-full p-2 px-2 ">
                <button
                    className="bg-gray-300 px-2 py-2 rounded-md transition duration-300"
                    onClick={() => setIsModalOpen(true)}
                >
                    Create Group
                </button>


                {selectedGroupIds.length > 0 && (
                    <>
                        <button
                            className="bg-red-200 px-4 py-2 rounded-md transition duration-300"
                            onClick={() => deleteGrp()}
                        >
                            <span className="text-red-500"> Delete Group </span>
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
                    onChange={(e) => handleSearch(e.target.value)}
                />
            </div>

        </div>

        <GetAllContactGroups contactgroupData={groupData as APIResponse<PaginatedResponse<ContactGroupData>>} loading={isLoading} onPageChange={(currentPage, pageSize) => onPageChange(currentPage, pageSize)} currentPage={currentPage} pageSize={pageSize} />

        <CreateGroup isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />
    </>




}

export default ContactGroupDash