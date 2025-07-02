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
import { Button, Input, Space, Badge, Tooltip } from "antd";
import { PlusOutlined, DeleteOutlined, SearchOutlined } from "@ant-design/icons";

const ContactGroupDash: React.FC = () => {
    const { selectedGroupIds, deleteGroup } = useContactGroupStore()
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query

    // Debounce the search query
    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);

    const { data: groupData, isLoading, refetch } = useContactGroupQuery(currentPage, pageSize, debouncedSearchQuery)

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
        <div className="flex justify-between items-center rounded-lg p-4 bg-white mt-10 shadow-sm border">
            <Space size="middle" className="flex-1">
                <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    size="large"
                    onClick={() => setIsModalOpen(true)}
                    className="shadow-sm"
                >
                    Create Group
                </Button>

                {selectedGroupIds.length > 0 && (
                    <Tooltip title={`Delete ${selectedGroupIds.length} selected group(s)`}>
                        <Button
                            danger
                            icon={<DeleteOutlined />}
                            size="large"
                            onClick={() => deleteGrp()}
                            className="shadow-sm"
                        >
                            Delete Group
                            <Badge
                                count={selectedGroupIds.length}
                                style={{ marginLeft: '8px' }}
                                showZero={false}
                            />
                        </Button>
                    </Tooltip>
                )}
            </Space>

            <Input
                placeholder="Search groups..."
                prefix={<SearchOutlined className="text-gray-400" />}
                allowClear
                size="large"
                style={{ width: 280 }}
                className="shadow-sm"
                onChange={(e) => handleSearch(e.target.value)}
            />
        </div>
        <GetAllContactGroups
            contactgroupData={groupData as APIResponse<PaginatedResponse<ContactGroupData>>}
            loading={isLoading}
            onPageChange={(currentPage, pageSize) => onPageChange(currentPage, pageSize)}
            currentPage={currentPage}
            pageSize={pageSize}
            refetch={refetch} />
        <CreateGroup isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} refetch={refetch} />
    </>
}

export default ContactGroupDash