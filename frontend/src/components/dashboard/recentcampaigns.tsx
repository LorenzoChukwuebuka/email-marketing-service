import { useNavigate } from 'react-router-dom';
import { useCampaignQuery } from '../../features/campaign/hooks/useCampaignQuery';
import EmptyState from '../emptyStateComponent';
import { useState } from 'react';
import LoadingSpinnerComponent from '../loadingSpinnerComponent';
import { Card, Button, Typography, Space, Tag, Divider, Avatar } from 'antd';
import { 
  ArrowRightOutlined, 
  EditOutlined, 
  ClockCircleOutlined,
  MailOutlined,
  PlusOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;

const RecentCampaigns = () => {
    const navigate = useNavigate();
    const [currentPage, _setCurrentPage] = useState(1);
    const [pageSize, _setPageSize] = useState(20);
    const { data: campaignData, isLoading } = useCampaignQuery(currentPage, pageSize);

    const handleViewAllCampaigns = () => {
        navigate("/app/campaign");
    };

    const handleCreateCampaign = () => {
        navigate("/app/campaign");
    };

    return (
        <div className="w-full">
            <Card 
                className="shadow-sm border-0 bg-white rounded-lg"
                bodyStyle={{ padding: '24px' }}
            >
                {/* Header */}
                <div className="flex justify-between items-center mb-6">
                    <div className="flex items-center space-x-3">
                        <div className="w-8 h-8 bg-blue-100 rounded-lg flex items-center justify-center">
                            <MailOutlined className="text-blue-600" />
                        </div>
                        <Title level={4} className="!mb-0 text-gray-900">
                            Recent campaigns
                        </Title>
                    </div>
                    <Button 
                        type="link" 
                        icon={<ArrowRightOutlined />}
                        onClick={handleViewAllCampaigns}
                        className="text-blue-600 hover:text-blue-700 font-medium px-0"
                    >
                        Go to Campaigns
                    </Button>
                </div>

                {/* Content */}
                {isLoading ? (
                    <div className="flex justify-center py-12">
                        <LoadingSpinnerComponent />
                    </div>
                ) : (
                    <>
                        {campaignData?.payload?.data && campaignData.payload.data.length > 0 ? (
                            <div className="space-y-4">
                                {campaignData.payload.data.slice(0, 3).map((campaign: any, index: number) => (
                                    <div key={index}>
                                        <div className="flex items-center justify-between py-4 hover:bg-gray-50 rounded-lg px-4 transition-colors duration-200 cursor-pointer">
                                            <div className="flex items-center space-x-4">
                                                <Avatar 
                                                    size={40}
                                                    className="bg-gradient-to-r from-blue-500 to-purple-600 flex items-center justify-center"
                                                >
                                                    <Text className="text-white font-medium">
                                                        {campaign.name.charAt(0).toUpperCase()}
                                                    </Text>
                                                </Avatar>
                                                <div className="flex flex-col">
                                                    <Text className="font-medium text-gray-900 text-base mb-1">
                                                        {campaign.name}
                                                    </Text>
                                                    <div className="flex items-center space-x-4">
                                                        <Tag 
                                                            icon={<EditOutlined />} 
                                                            color="orange"
                                                            className="border-0 bg-orange-50 text-orange-600 rounded-full px-3 py-1"
                                                        >
                                                            Draft
                                                        </Tag>
                                                        <div className="flex items-center space-x-1 text-gray-500">
                                                            <ClockCircleOutlined className="text-xs" />
                                                            <Text className="text-xs">
                                                                Last edit {formatDate(campaign.updated_at || campaign.created_at)}
                                                            </Text>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                            <Button 
                                                type="text" 
                                                icon={<ArrowRightOutlined />}
                                                className="text-gray-400 hover:text-blue-600 hover:bg-blue-50"
                                                onClick={() => navigate(`/app/campaign/${campaign.id}`)}
                                            />
                                        </div>
                                        {index < campaignData.payload.data.slice(0, 3).length - 1 && (
                                            <Divider className="my-0" />
                                        )}
                                    </div>
                                ))}
                                
                                {/* View All Footer */}
                                <div className="pt-4 border-t border-gray-100">
                                    <Button 
                                        type="text" 
                                        block 
                                        icon={<ArrowRightOutlined />}
                                        onClick={handleViewAllCampaigns}
                                        className="text-blue-600 hover:text-blue-700 hover:bg-blue-50 font-medium h-10"
                                    >
                                        View all campaigns
                                    </Button>
                                </div>
                            </div>
                        ) : (
                            <div className="text-center py-12">
                                <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                                    <MailOutlined className="text-2xl text-gray-400" />
                                </div>
                                <Title level={5} className="text-gray-900 mb-2">
                                    You have not created any Campaign
                                </Title>
                                <Text className="text-gray-500 block mb-6">
                                    Create and easily send marketing emails to your audience
                                </Text>
                                <Button 
                                    type="primary" 
                                    icon={<PlusOutlined />}
                                    onClick={handleCreateCampaign}
                                    className="bg-blue-600 hover:bg-blue-700 border-0 rounded-lg h-10 px-6 font-medium"
                                >
                                    Create Campaign
                                </Button>
                            </div>
                        )}
                    </>
                )}
            </Card>
        </div>
    );
};

// Helper function to format dates
const formatDate = (dateString?: string): string => {
    if (!dateString) return 'Unknown date';

    // Remove the " WAT" part if it exists
    const cleanDate = dateString.replace(" WAT", "");

    try {
        const date = new Date(cleanDate);
        if (isNaN(date.getTime())) return 'Invalid date';

        const now = new Date();
        const diffTime = Math.abs(now.getTime() - date.getTime());
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

        if (diffDays === 0) return 'Today';
        if (diffDays === 1) return 'Yesterday';
        if (diffDays < 7) return `${diffDays} days ago`;
        
        return date.toLocaleDateString('en-GB', {
            day: '2-digit',
            month: 'short',
            year: 'numeric'
        });
    } catch {
        return 'Invalid date';
    }
};

export default RecentCampaigns;