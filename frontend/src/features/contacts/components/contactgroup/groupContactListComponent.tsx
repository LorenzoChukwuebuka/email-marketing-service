import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Table, Button, Input, Card, Typography, Space, message } from 'antd';
import { DeleteOutlined, ArrowLeftOutlined, SearchOutlined } from '@ant-design/icons';
import type { ColumnsType, TableProps } from 'antd/es/table';
import useContactGroupStore from "../../store/contactgroup.store";
import { ContactGroupData, ContactInGroup } from '../../interface/contactgroup.interface';
import { useSingleContactGroupQuery } from "../../hooks/useContactGroupQuery";
import LoadingSpinnerComponent from "../../../../components/loadingSpinnerComponent";

const { Title, Text } = Typography;

const GroupContactList: React.FC = () => {
    const navigate = useNavigate();
    const { setSelectedContactIds, setSelectedGroupIds, removeContactFromGroup } = useContactGroupStore();
    const [group, setGroup] = useState<ContactGroupData | null>(null);
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
    const [searchText, setSearchText] = useState<string>('');
    const [filteredContacts, setFilteredContacts] = useState<ContactInGroup[]>([]);

    const location = useLocation();
    const stateData = location.state as { groupId: string };

    const { data: contactgroupData, isLoading, refetch } = useSingleContactGroupQuery(stateData.groupId);
    // Filter contacts based on search text
    useEffect(() => {
        if (!group?.contacts) {
            setFilteredContacts([]);
            return;
        }

        if (!searchText) {
            setFilteredContacts(group.contacts);
        } else {
            const filtered = group.contacts.filter(contact =>
                contact.contact_first_name.toLowerCase().includes(searchText.toLowerCase()) ||
                contact.contact_last_name.toLowerCase().includes(searchText.toLowerCase()) ||
                contact.contact_email.toLowerCase().includes(searchText.toLowerCase())
            );
            setFilteredContacts(filtered);
        }
    }, [group?.contacts, searchText]);

    // Handle row selection
    const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
        setSelectedRowKeys(newSelectedRowKeys);
        setSelectedContactIds(newSelectedRowKeys as string[]);
    };

    const rowSelection: TableProps<ContactInGroup>['rowSelection'] = {
        selectedRowKeys,
        onChange: onSelectChange,
        selections: [
            Table.SELECTION_ALL,
            Table.SELECTION_INVERT,
            Table.SELECTION_NONE,
        ],
    };

    // Handle remove contacts
    const handleRemoveContacts = async () => {
        try {
            const stateData = location.state as { groupId: string };
            setSelectedGroupIds([stateData.groupId]);
            await removeContactFromGroup();
            message.success('Contacts removed successfully');
            setSelectedRowKeys([]);
            refetch()
        } catch (error) {
            console.error(error);
            message.error('Failed to remove contacts');
        }
    };

    // Handle search
    const handleSearch = (value: string) => {
        setSearchText(value);
    };

    // Handle back navigation
    const handleBack = () => {
        navigate(-1);
    };

    // Format date
    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString('en-US', {
            timeZone: 'UTC',
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: 'numeric',
            minute: 'numeric',
            second: 'numeric'
        });
    };

    // Table columns
    const columns: ColumnsType<ContactInGroup> = [
        {
            title: 'First Name',
            dataIndex: 'contact_first_name',
            key: 'contact_first_name',
            sorter: (a, b) => a.contact_first_name.localeCompare(b.contact_first_name),
        },
        {
            title: 'Last Name',
            dataIndex: 'contact_last_name',
            key: 'contact_last_name',
            sorter: (a, b) => a.contact_last_name.localeCompare(b.contact_last_name),
        },
        {
            title: 'Email',
            dataIndex: 'contact_email',
            key: 'contact_email',
            sorter: (a, b) => a.contact_email.localeCompare(b.contact_email),
        },
        {
            title: 'From',
            dataIndex: 'contact_from_origin',
            key: 'contact_from_origin',
            sorter: (a, b) => a.contact_from_origin.localeCompare(b.contact_from_origin),
        },
        {
            title: 'Subscribed',
            dataIndex: 'contact_is_subscribed',
            key: 'contact_is_subscribed',
            render: (subscribed: boolean) => (
                <span style={{ color: subscribed ? '#52c41a' : '#ff4d4f' }}>
                    {subscribed ? 'Yes' : 'No'}
                </span>
            ),
            filters: [
                { text: 'Subscribed', value: true },
                { text: 'Unsubscribed', value: false },
            ],
            onFilter: (value, record) => record.contact_is_subscribed === value,
        },
        {
            title: 'Created At',
            dataIndex: 'contact_created_at',
            key: 'contact_created_at',
            render: (date: string) => formatDate(date),
            sorter: (a, b) => new Date(a.contact_created_at).getTime() - new Date(b.contact_created_at).getTime(),
        },
    ];

    useEffect(() => {
        if (contactgroupData) {
            if (Array.isArray(contactgroupData) && contactgroupData.length > 0) {
                setGroup(contactgroupData[0]);
            } else if (!Array.isArray(contactgroupData)) {
                setGroup(contactgroupData);
            }
        }
    }, [contactgroupData]);

    if (isLoading) {
        return (
            <div style={{ padding: '24px', textAlign: 'center' }}>
                <LoadingSpinnerComponent />
            </div>
        );
    }

    if (!group) {
        return (
            <div style={{ padding: '24px', textAlign: 'center' }}>
                <Text>No group data available.</Text>
            </div>
        );
    }

    return (
        <div style={{ padding: '24px' }}>
            {/* Header */}
            <Card style={{ marginBottom: '24px' }}>
                <Space direction="vertical" size="small" style={{ width: '100%' }}>
                    <Button
                        type="link"
                        icon={<ArrowLeftOutlined />}
                        onClick={handleBack}
                        style={{ padding: 0, height: 'auto' }}
                    >
                        Back
                    </Button>
                    <Title level={2} style={{ margin: 0 }}>
                        Group: {group.group_name}
                    </Title>
                    <Text type="secondary">
                        Description: {group.description}
                    </Text>
                </Space>
            </Card>

            {/* Actions Bar */}
            <Card style={{ marginBottom: '16px' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Space>
                        {selectedRowKeys.length > 0 && (
                            <Button
                                type="primary"
                                danger
                                icon={<DeleteOutlined />}
                                onClick={handleRemoveContacts}
                            >
                                Remove Contact{selectedRowKeys.length > 1 ? 's' : ''} ({selectedRowKeys.length})
                            </Button>
                        )}
                    </Space>
                    <Input
                        placeholder="Search contacts..."
                        prefix={<SearchOutlined />}
                        value={searchText}
                        onChange={(e) => handleSearch(e.target.value)}
                        style={{ width: '300px' }}
                        allowClear
                    />
                </div>
            </Card>

            {/* Table */}
            <Card>
                <Table<ContactInGroup>
                    rowSelection={rowSelection}
                    columns={columns}
                    dataSource={filteredContacts}
                    rowKey="contact_id"
                    pagination={{
                        showSizeChanger: true,
                        showQuickJumper: true,
                        showTotal: (total, range) =>
                            `${range[0]}-${range[1]} of ${total} contacts`,
                    }}
                    scroll={{ x: 'max-content' }}
                    locale={{
                        emptyText: 'No contacts available'
                    }}
                />
            </Card>
        </div>
    );
};

export default GroupContactList;