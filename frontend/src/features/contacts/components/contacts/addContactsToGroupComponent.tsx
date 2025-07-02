import { useState, useMemo } from "react";
import { Modal, Radio, Button, Empty, Card, Avatar, Badge, Space, Typography, Divider } from "antd";
import { TeamOutlined, UserOutlined, CheckCircleOutlined } from "@ant-design/icons";
import useContactGroupStore from "./../../store/contactgroup.store";
import { useContactGroupQuery } from "../../hooks/useContactGroupQuery";

const { Text, Title } = Typography;

interface CGProps {
    isOpen: boolean;
    onClose: () => void;
}


const AddContactsToGroupComponent: React.FC<CGProps> = ({ isOpen, onClose }) => {
    const { setSelectedGroupIds, addContactToGroup, selectedContactIds } = useContactGroupStore();
    const [currentPage] = useState(1);
    const [pageSize] = useState(20);
    const [loading, setLoading] = useState(false);

    const { data: contactgroupData } = useContactGroupQuery(currentPage, pageSize, undefined);

    const cData = useMemo(() => contactgroupData?.payload.data || [], [contactgroupData]);

    const [selectedGroup, setSelectedGroup] = useState<string>("");

    const handleGroupSubmit = async () => {
        if (selectedGroup) {
            setLoading(true);
            try {
                setSelectedGroupIds([selectedGroup]);
                await addContactToGroup();
                setSelectedGroup("");
                onClose();
            } catch (error) {
                console.error('Error adding contacts to group:', error);
            } finally {
                setLoading(false);
            }
        }
    };

    const handleCancel = () => {
        setSelectedGroup("");
        onClose();
    };

    const handleRadioChange = (e: any) => {
        setSelectedGroup(e.target.value);
    };

    const selectedContactCount = selectedContactIds?.length || 0;

    return (
        <Modal
            title={
                <div className="flex items-center space-x-2">
                    <TeamOutlined className="text-blue-500" />
                    <span>Add Contacts to Group</span>
                </div>
            }
            open={isOpen}
            onCancel={handleCancel}
            width={600}
            footer={[
                <Button key="cancel" onClick={handleCancel}>
                    Cancel
                </Button>,
                <Button
                    key="submit"
                    type="primary"
                    onClick={handleGroupSubmit}
                    disabled={!selectedGroup}
                    loading={loading}
                    icon={<CheckCircleOutlined />}
                >
                    Add to Group
                </Button>
            ]}
        >
            <div className="space-y-6">
                {/* Header Info */}
                <Card size="small" className="bg-blue-50 border-blue-200">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-2">
                            <UserOutlined className="text-blue-500" />
                            <Text>
                                <Badge count={selectedContactCount} className="mr-2" />
                                contact{selectedContactCount !== 1 ? 's' : ''} selected
                            </Text>
                        </div>
                        <Text type="secondary" className="text-sm">
                            Choose a group to add them to
                        </Text>
                    </div>
                </Card>

                <Divider className="my-4" />

                {/* Group Selection */}
                <div>
                    <Title level={5} className="mb-3">
                        Select a Group
                    </Title>

                    <div className="max-h-80 overflow-y-auto pr-2">
                        {cData && cData.length > 0 ? (
                            <Radio.Group
                                onChange={handleRadioChange}
                                value={selectedGroup}
                                className="w-full"
                            >
                                <Space direction="vertical" className="w-full" size="middle">
                                    {cData.map((group) => (
                                        <Card
                                            key={group.group_id}
                                            size="small"
                                            className={`cursor-pointer transition-all duration-200 hover:shadow-md ${selectedGroup === group.group_id
                                                ? 'border-blue-500 bg-blue-50'
                                                : 'border-gray-200 hover:border-blue-300'
                                                }`}
                                            onClick={() => setSelectedGroup(group.group_id)}
                                        >
                                            <Radio value={group.group_id} className="w-full">
                                                <div className="flex items-start space-x-3 ml-2">
                                                    <Avatar
                                                        icon={<TeamOutlined />}
                                                        size="small"
                                                        className="bg-blue-500 flex-shrink-0 mt-1"
                                                    />
                                                    <div className="flex-1 min-w-0">
                                                        <div className="flex items-center justify-between">
                                                            <Text strong className="text-base">
                                                                {group.group_name}
                                                            </Text>
                                                            <Badge
                                                                count={group.contacts ? group.contacts.length : 0}
                                                                showZero
                                                                className="bg-gray-100"
                                                            />
                                                        </div>
                                                        {group.description && (
                                                            <Text
                                                                type="secondary"
                                                                className="text-sm mt-1 block"
                                                            >
                                                                {group.description}
                                                            </Text>
                                                        )}
                                                    </div>
                                                </div>
                                            </Radio>
                                        </Card>
                                    ))}
                                </Space>
                            </Radio.Group>
                        ) : (
                            <div className="text-center py-8">
                                <Empty
                                    image={Empty.PRESENTED_IMAGE_SIMPLE}
                                    description={
                                        <span>
                                            No groups available
                                            <br />
                                            <Text type="secondary" className="text-sm">
                                                Create a group first to organize your contacts
                                            </Text>
                                        </span>
                                    }
                                />
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </Modal>
    );
};

export default AddContactsToGroupComponent;