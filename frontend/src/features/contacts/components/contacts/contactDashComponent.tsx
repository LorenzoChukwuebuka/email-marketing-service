import { FormEvent, useState } from "react";
import { Button, Input, Modal, Badge, Tooltip, Card } from "antd";
import {
    PlusOutlined,
    ImportOutlined,
    DeleteOutlined,
    TeamOutlined,
    SearchOutlined,
} from "@ant-design/icons";
import useDebounce from "../../../../hooks/useDebounce";
import CreateContact from './createContact';
import useContactStore from "../../store/contact.store";
import GetAllContacts from "./getAllContacts";
import { useContactQuery } from "../../hooks/useContactQuery";
import ContactUpload from "./contactBatchUploadComponent";
import AddContactsToGroupComponent from './addContactsToGroupComponent';
import { APIResponse } from "../../../../interface/api.interface";
import { PaginatedResponse } from "../../../../interface/pagination.interface";
import { ContactAPIResponse } from "../../interface/contact.interface";

const ContactsDashComponent: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { selectedIds, deleteContact } = useContactStore();
    const [importModalOpen, setImportModalOpen] = useState<boolean>(false);
    const [groupModalOpen, setGroupModalOpen] = useState<boolean>(false);
    const [searchQuery, setSearchQuery] = useState<string>("");

    // Debounce the search query
    const debouncedSearchQuery = useDebounce(searchQuery, 300);
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);

    const { data: contacts, isLoading, refetch } = useContactQuery(currentPage, pageSize, debouncedSearchQuery);

    const handleDelete = async (e: FormEvent<HTMLButtonElement>) => {
        e.preventDefault();
        Modal.confirm({
            title: "Delete Contacts",
            content: `Are you sure you want to delete ${selectedIds.length} contact${selectedIds.length > 1 ? 's' : ''}?`,
            okText: "Delete",
            cancelText: "Cancel",
            okType: "danger",
            icon: <DeleteOutlined />,
            onOk: async () => {
                try {
                    await deleteContact();
                    await new Promise(resolve => setTimeout(resolve, 1000));
                    refetch();
                } catch (error) {
                    console.error('Error deleting contacts:', error);
                }
            },
        });
    };

    const onPageChange = (page: number, size: number) => {
        setCurrentPage(page);
        setPageSize(size);
    };

    const addContactToGroup = () => {
        setGroupModalOpen(true);
    };

    const importContact = () => {
        setImportModalOpen(true);
    };

    const handleSearchInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        setSearchQuery(e.target.value);
    };

    const selectedCount = selectedIds.length;

    return (
        <div className="space-y-6 mt-10">
            {/* Header Card */}
            <Card className="shadow-sm">
                <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                    {/* Action Buttons */}
                    <div className="flex flex-wrap items-center gap-2">
                        <Button
                            type="primary"
                            icon={<PlusOutlined />}
                            onClick={() => setIsModalOpen(true)}
                            size="middle"
                        >
                            Create Contact
                        </Button>

                        <Button
                            icon={<ImportOutlined />}
                            onClick={importContact}
                            size="middle"
                        >
                            Import Contacts
                        </Button>

                        {/* Bulk Actions - Show when contacts are selected */}
                        {selectedCount > 0 && (
                            <div className="flex items-center gap-2 ml-4 pl-4 border-l border-gray-200">
                                <Badge count={selectedCount} className="mr-2">
                                    <span className="text-sm text-gray-600">Selected</span>
                                </Badge>

                                <Tooltip title={`Delete ${selectedCount} contact${selectedCount > 1 ? 's' : ''}`}>
                                    <Button
                                        danger
                                        icon={<DeleteOutlined />}
                                        onClick={handleDelete as any}
                                        size="middle"
                                    >
                                        Delete
                                    </Button>
                                </Tooltip>

                                <Tooltip title={`Add ${selectedCount} contact${selectedCount > 1 ? 's' : ''} to group`}>
                                    <Button
                                        type="primary"
                                        ghost
                                        icon={<TeamOutlined />}
                                        onClick={addContactToGroup}
                                        size="middle"
                                    >
                                        Add to Group
                                    </Button>
                                </Tooltip>
                            </div>
                        )}
                    </div>

                    {/* Search Input */}
                    <div className="flex-shrink-0">
                        <Input
                            placeholder="Search contacts..."
                            prefix={<SearchOutlined className="text-gray-400" />}
                            value={searchQuery}
                            onChange={handleSearchInput}
                            allowClear
                            className="w-64"
                            size="middle"
                        />
                    </div>
                </div>
            </Card>

            {/* Contacts Table */}
            <Card className="shadow-sm">
                <GetAllContacts
                    contactData={contacts as APIResponse<PaginatedResponse<ContactAPIResponse>>}
                    loading={isLoading}
                    currentPage={currentPage}
                    pageSize={pageSize}
                    onPageChange={onPageChange}
                />
            </Card>

            {/* Modals */}
            <CreateContact
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                refetch={refetch}
            />

            <ContactUpload
                isOpen={importModalOpen}
                onClose={() => setImportModalOpen(false)}
                refetch={refetch}
            />

            <AddContactsToGroupComponent
                isOpen={groupModalOpen}
                onClose={() => setGroupModalOpen(false)}
            />
        </div>
    );
};

export default ContactsDashComponent;