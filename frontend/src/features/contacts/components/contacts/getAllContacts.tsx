import { useState, useEffect } from "react";
import { Table } from 'antd';
import { EditOutlined } from '@ant-design/icons';
import EmptyState from "../../../../components/emptyStateComponent";
import useContactStore from "../../store/contact.store";
import { ContactAPIResponse } from "../../interface/contact.interface";
import EditContact from './editContactComponent';
import useContactGroupStore from "../../store/contactgroup.store";
import { APIResponse } from "../../../../interface/api.interface";
import { PaginatedResponse } from "../../../../interface/pagination.interface";
import LoadingSpinnerComponent from "../../../../components/loadingSpinnerComponent";

interface Props {
    contactData: APIResponse<PaginatedResponse<ContactAPIResponse>>
    loading: boolean
    currentPage: number
    pageSize: number
    onPageChange: (page: number, size: number) => void
}

const GetAllContacts: React.FC<Props> = ({ contactData, loading, currentPage, pageSize, onPageChange }) => {
    const { selectedIds, setSelectedId } = useContactStore();
    const { setSelectedContactIds } = useContactGroupStore()
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [selectedContact, setSelectedContact] = useState<ContactAPIResponse | null>(null);

    // Local state for selection to ensure it works
    const [localSelectedIds, setLocalSelectedIds] = useState<string[]>([]);
    const cdata = contactData?.payload?.data || []

    // Sync local state with store state
    useEffect(() => {
        setLocalSelectedIds(selectedIds);
    }, [selectedIds]);

    // Update stores when local selection changes
    useEffect(() => {
        setSelectedId(localSelectedIds);
        setSelectedContactIds(localSelectedIds);
    }, [localSelectedIds, setSelectedId, setSelectedContactIds]);

    const openEditModal = (contact: ContactAPIResponse) => {
        setSelectedContact(contact);
        setIsModalOpen(true);
    };

    const closeEditModal = () => {
        setIsModalOpen(false);
        setSelectedContact(null);
    };

    const columns = [
        {
            title: 'First Name',
            dataIndex: 'first_name',
            key: 'first_name',
        },
        {
            title: 'Last Name',
            dataIndex: 'last_name',
            key: 'last_name',
        },
        {
            title: 'Email',
            dataIndex: 'email',
            key: 'email',
        },
        {
            title: 'From',
            dataIndex: 'from_origin',
            key: 'from_origin',
        },
        {
            title: 'Created On',
            dataIndex: 'contact_created_at',
            key: 'contact_created_at',
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
            render: (_: any, record: ContactAPIResponse) => (
                <button
                    className="text-gray-400 hover:text-gray-600 p-1"
                    onClick={() => openEditModal(record)}
                >
                    <EditOutlined />
                </button>
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
        onSelect: (record: ContactAPIResponse, selected: boolean) => {
            if (selected) {
                setLocalSelectedIds(prev => [...prev, record.contact_id]);
            } else {
                setLocalSelectedIds(prev => prev.filter(id => id !== record.contact_id));
            }
        },
        onSelectAll: (selected: boolean) => {
            if (selected) {
                const allIds = cdata.map(contact => contact.contact_id);
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
        total: contactData?.payload?.total || 0,
        onChange: onPageChange,
        showSizeChanger: true,
        pageSizeOptions: ["10", "20", "50", "100"],
        showTotal: (total: number) => `Total ${total} Contacts`,
        className: "mt-4 flex justify-center items-center mb-5"
    };

    if (loading) {
        return <LoadingSpinnerComponent />;
    }

    if (!contactData || !contactData?.payload?.data?.length) {
        return (
            <div className="py-4 px-4 text-center">
                <EmptyState
                    title="You have not created any Contacts"
                    description="Create contacts"
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
                    dataSource={cdata}
                    rowKey="contact_id"
                    rowSelection={rowSelection}
                    pagination={paginationConfig}
                    className="bg-white rounded-sm"
                    rowClassName="hover:bg-slate-100"
                />
            </div>

            <EditContact
                isOpen={isModalOpen}
                onClose={closeEditModal}
                contact={selectedContact}
            />
        </>
    );
};

export default GetAllContacts;