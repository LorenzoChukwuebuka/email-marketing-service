import { useMemo, useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { Table } from 'antd';
import { EditOutlined, EyeOutlined } from '@ant-design/icons';
import EmptyState from "../../../../components/emptyStateComponent";
import EditGroupComponent from "./editGroupComponent";
import useContactGroupStore from "../../store/contactgroup.store";
import { ContactGroupData } from '../../interface/contactgroup.interface';
import { APIResponse } from '../../../../../../frontend/src/interface/api.interface';
import { PaginatedResponse } from '../../../../../../frontend/src/interface/pagination.interface';
import LoadingSpinnerComponent from "../../../../components/loadingSpinnerComponent";

interface Props {
    contactgroupData?: APIResponse<PaginatedResponse<ContactGroupData>>
    loading: boolean
    onPageChange: (currentPage: number, pageSize: number) => void
    currentPage: number
    pageSize: number
    refetch: () => void
}

const GetAllContactGroups: React.FC<Props> = ({ contactgroupData, loading, onPageChange, currentPage, pageSize,refetch }) => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [selectedGroup, setSelectedGroup] = useState<ContactGroupData | null>(null);

    // Local state for selection to ensure it works
    const [localSelectedIds, setLocalSelectedIds] = useState<string[]>([]);
    const { selectedGroupIds, setSelectedGroupIds } = useContactGroupStore()
    const cgData = useMemo(() => contactgroupData?.payload?.data || [], [contactgroupData])

    // Sync local state with store state
    useEffect(() => {
        setLocalSelectedIds(selectedGroupIds);
    }, [selectedGroupIds]);

    // Update stores when local selection changes
    useEffect(() => {
        setSelectedGroupIds(localSelectedIds);
    }, [localSelectedIds, setSelectedGroupIds]);

    const openEditModal = (group: ContactGroupData) => {
        setIsModalOpen(true)
        setSelectedGroup(group)
    }

    const columns = [
        {
            title: 'Name',
            dataIndex: 'group_name',
            key: 'group_name',
        },
        {
            title: 'Description',
            dataIndex: 'description',
            key: 'description',
        },
        {
            title: 'Contacts',
            dataIndex: 'contacts',
            key: 'contacts',
            render: (contacts: any[]) => contacts ? contacts.length : 0,
        },
        {
            title: 'Created At',
            dataIndex: 'group_created_at',
            key: 'group_created_at',
            render: (date: string) => {
                return new Date(date).toLocaleString('en-US', {
                    timeZone: 'UTC',
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                    hour: 'numeric',
                    minute: 'numeric',
                    second: 'numeric',
                });
            }
        },
        {
            title: 'Edit',
            key: 'edit',
            render: (_: any, record: ContactGroupData) => (
                <button
                    className="text-gray-400 hover:text-gray-600 p-1"
                    onClick={() => openEditModal(record)}
                >
                    <EditOutlined />
                </button>
            ),
        },
        {
            title: 'View All Contacts',
            key: 'view_contacts',
            render: (_: any, record: ContactGroupData) => (
                <Link
                    to="/app/contacts/view-group"
                    state={{ groupId: record.group_id }}
                    className="text-gray-400 hover:text-gray-600"
                >
                    <EyeOutlined />
                </Link>
            ),
        },
    ];

    // Row selection configuration using local state
    const rowSelection = {
        selectedRowKeys: localSelectedIds,
        onChange: (selectedRowKeys: React.Key[]) => {
            const stringKeys = selectedRowKeys.map(key => String(key));
            setLocalSelectedIds(stringKeys);
        },
        onSelect: (record: ContactGroupData, selected: boolean) => {
            if (selected) {
                setLocalSelectedIds(prev => [...prev, record.group_id]);
            } else {
                setLocalSelectedIds(prev => prev.filter(id => id !== record.group_id));
            }
        },
        onSelectAll: (selected: boolean) => {
            if (selected) {
                const allIds = cgData.map(group => group.group_id);
                setLocalSelectedIds(allIds);
            } else {
                setLocalSelectedIds([]);
            }
        },
    };

    // Pagination configuration
    const paginationConfig = {
        current: currentPage,
        pageSize: pageSize,
        total: contactgroupData?.payload?.total || 0,
        onChange: onPageChange,
        showSizeChanger: true,
        pageSizeOptions: ["10", "20", "50", "100"],
        showTotal: (total: number) => `Total ${total} groups`,
        className: "mt-4 flex justify-center items-center mb-5"
    };

    if (loading) {
        return <LoadingSpinnerComponent />;
    }

    if (!cgData || cgData.length === 0) {
        return (
            <div className="py-4 px-4 text-center">
                <EmptyState
                    title="You have not created any groups"
                    description="Create groups to easily manage your contacts"
                    icon={<i className="bi bi-emoji-frown text-xl"></i>}
                />
            </div>
        );
    }

    return (
        <>
            <div className="mt-8">
                <Table
                    columns={columns}
                    dataSource={cgData}
                    rowKey="group_id"
                    rowSelection={rowSelection}
                    pagination={paginationConfig}
                    className="bg-white rounded-sm"
                    rowClassName="hover:bg-gray-100"
                />

            </div>

            <EditGroupComponent
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                group={selectedGroup as any}
                refetch={refetch} 
            />
        </>
    );
};

export default GetAllContactGroups;