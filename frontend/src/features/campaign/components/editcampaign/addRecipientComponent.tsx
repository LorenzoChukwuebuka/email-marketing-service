import { useEffect, useMemo, useState } from "react";
import { Modal, Checkbox, Tag, Space, Typography, Empty, Badge, Card} from "antd";
import { UserOutlined, TeamOutlined, CloseOutlined } from "@ant-design/icons";
import { useContactGroupQuery } from "./../../../contacts/hooks/useContactGroupQuery";
import { ContactGroupData } from "./../../../contacts/interface/contactgroup.interface";
import useCampaignStore from "./../../store/campaign.store";
import { Campaign } from "./../../interface/campaign.interface";
import { BaseEntity } from "./../../../../interface/baseentity.interface";

const { Title, Text } = Typography;

interface Props {
    isOpen: boolean;
    onClose: () => void;
    campaign: (Campaign & BaseEntity) | null;
}

const AddCampaignRecipients: React.FC<Props> = ({ isOpen, onClose, campaign }) => {
    const { createCampaignGroup } = useCampaignStore();
    const [selectedGroups, setSelectedGroups] = useState<ContactGroupData[]>([]);
    const [loading, setLoading] = useState(false);
    
    const currentPage = 1;
    const pageSize = 2000;

    const { data: contactgroupData, isLoading } = useContactGroupQuery(currentPage, pageSize, undefined);

    const availableGroups = useMemo(() => 
        contactgroupData?.payload.data || [], 
        [contactgroupData]
    );

    // Initialize selected groups based on campaign's existing groups
    useEffect(() => {
        if (availableGroups.length > 0 && campaign?.groups) {
            const preSelectedGroups = availableGroups.filter(group => 
                (campaign?.groups as any).some(campaignGroup => campaignGroup.id === group.group_id)
            );
            setSelectedGroups(preSelectedGroups);
        } else {
            setSelectedGroups([]);
        }
    }, [availableGroups, campaign]);

    const handleGroupToggle = (group: ContactGroupData) => {
        setSelectedGroups(prev => {
            const isSelected = prev.some(g => g.group_id === group.group_id);
            return isSelected
                ? prev.filter(g => g.group_id !== group.group_id)
                : [...prev, group];
        });
    };

    const handleRemoveGroup = (groupId: string) => {
        setSelectedGroups(prev => prev.filter(g => g.group_id !== groupId));
    };

    const handleSubmit = async () => {
        if (!campaign?.id) {
            console.error("Campaign ID is not available");
            return;
        }

        setLoading(true);
        try {
            const groupIds = selectedGroups.map(group => group.group_id);
            await createCampaignGroup(campaign.id, groupIds);
            onClose();
        } catch (error) {
            console.error("Error creating campaign groups:", error);
        } finally {
            setLoading(false);
        }
    };

    const totalContacts = selectedGroups.reduce((sum, group) => 
        sum + (group.contacts?.length || 0), 0
    );

    return (
        <Modal
            open={isOpen}
            onCancel={onClose}
            onOk={handleSubmit}
            title={
                <div className="flex items-center gap-2">
                    <TeamOutlined className="text-blue-500" />
                    <span>Add Recipients to Campaign</span>
                </div>
            }
            okText="Add Recipients"
            cancelText="Cancel"
            okButtonProps={{ 
                disabled: selectedGroups.length === 0,
                loading: loading
            }}
            width={600}
            className="modern-modal"
        >
            <div className="space-y-6">
                {/* Selected Groups Summary */}
                {selectedGroups.length > 0 && (
                    <Card size="small" className="bg-blue-50 border-blue-200">
                        <div className="space-y-3">
                            <div className="flex items-center justify-between">
                                <Title level={5} className="m-0 text-blue-800">
                                    Selected Groups ({selectedGroups.length})
                                </Title>
                                <Badge 
                                    count={totalContacts} 
                                    showZero 
                                    color="blue" 
                                    title="Total contacts"
                                />
                            </div>
                            <div className="flex flex-wrap gap-2">
                                {selectedGroups.map((group) => (
                                    <Tag
                                        key={group.group_id}
                                        closable
                                        onClose={() => handleRemoveGroup(group.group_id)}
                                        closeIcon={<CloseOutlined />}
                                        color="blue"
                                        className="px-3 py-1 text-sm rounded-full"
                                    >
                                        <Space size={4}>
                                            <TeamOutlined />
                                            {group.group_name}
                                            <Badge 
                                                count={group.contacts?.length || 0} 
                                                size="small"
                                                color="rgba(255, 255, 255, 0.8)"
                                                style={{ 
                                                    color: '#1890ff',
                                                    backgroundColor: 'rgba(255, 255, 255, 0.3)'
                                                }}
                                            />
                                        </Space>
                                    </Tag>
                                ))}
                            </div>
                        </div>
                    </Card>
                )}

                {/* Available Groups */}
                <div>
                    <Title level={5} className="mb-4 text-gray-800">
                        Available Contact Groups
                    </Title>
                    
                    <div className="max-h-80 overflow-y-auto space-y-3 p-2">
                        {isLoading ? (
                            <div className="flex justify-center py-8">
                                <Text type="secondary">Loading contact groups...</Text>
                            </div>
                        ) : availableGroups.length > 0 ? (
                            availableGroups.map((group: ContactGroupData) => {
                                const isSelected = selectedGroups.some(g => g.group_id === group.group_id);
                                const contactCount = group.contacts?.length || 0;
                                
                                return (
                                    <Card 
                                        key={group.group_id} 
                                        size="small"
                                        className={`cursor-pointer transition-all duration-200 hover:shadow-md ${
                                            isSelected 
                                                ? 'border-blue-500 bg-blue-50' 
                                                : 'border-gray-200 hover:border-blue-300'
                                        }`}
                                        onClick={() => handleGroupToggle(group)}
                                    >
                                        <div className="flex items-start gap-3">
                                            <Checkbox
                                                checked={isSelected}
                                                onChange={() => handleGroupToggle(group)}
                                                className="mt-1"
                                            />
                                            <div className="flex-1 min-w-0">
                                                <div className="flex items-center justify-between mb-2">
                                                    <Text strong className="text-gray-800">
                                                        {group.group_name}
                                                    </Text>
                                                    <Badge 
                                                        count={contactCount} 
                                                        showZero
                                                        color={isSelected ? 'blue' : 'default'}
                                                        title={`${contactCount} contacts`}
                                                    />
                                                </div>
                                                {group.description && (
                                                    <Text 
                                                        type="secondary" 
                                                        className="text-sm leading-relaxed block"
                                                    >
                                                        {group.description.trim()}
                                                    </Text>
                                                )}
                                                <div className="flex items-center gap-1 mt-2">
                                                    <UserOutlined className="text-gray-400 text-xs" />
                                                    <Text type="secondary" className="text-xs">
                                                        {contactCount} {contactCount === 1 ? 'contact' : 'contacts'}
                                                    </Text>
                                                </div>
                                            </div>
                                        </div>
                                    </Card>
                                );
                            })
                        ) : (
                            <Empty
                                image={Empty.PRESENTED_IMAGE_SIMPLE}
                                description="No contact groups available"
                                className="py-8"
                            />
                        )}
                    </div>
                </div>
            </div>
        </Modal>
    );
};

export default AddCampaignRecipients;