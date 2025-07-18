import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
    DatePicker,
    Card,
    Button,
    Typography,
    Space,
    Tag,
    Modal,
    Divider,
    Avatar,
    Tooltip,
    Badge
} from 'antd';
import {
    ArrowLeftOutlined,
    SendOutlined,
    ScheduleOutlined,
    UserOutlined,
    TeamOutlined,
    EditOutlined,
    DesktopOutlined,
    PlusOutlined,
    InfoCircleOutlined,
    EyeOutlined,
    SettingOutlined
} from '@ant-design/icons';
import dayjs from 'dayjs';
import { Helmet, HelmetProvider } from 'react-helmet-async';
import useCampaignStore from '../../store/campaign.store';
import { useSingleCampaignQuery } from '../../hooks/useCampaignQuery';
import LoadingSpinnerComponent from '../../../../components/loadingSpinnerComponent';
import AddSenderComponent from './addSenderCampaignComponent';
import AddCampaignRecipients from './addRecipientComponent';
import AddCampaignSubjectComponent from './addSubjectComponent';
import { CampaignData } from '../../interface/campaign.interface';

const { Title, Text } = Typography;

const EditCampaignForm: React.FC = () => {
    const { id } = useParams<{ id: string }>() as { id: string };
    const { setCurrentCampaignId, setCreateCampaignValues, sendCampaign, updateCampaign } = useCampaignStore();
    const navigate = useNavigate();
    const [isSubjectModalOpen, setIsSubjectModalOpen] = useState<boolean>(false);
    const [isRecipientModalOpen, setIsRecipientModalOpen] = useState<boolean>(false);
    const [isSenderModalOpen, setIsSenderModalOpen] = useState<boolean>(false);
    const [isCalendarOpen, setIsCalendarOpen] = useState<boolean>(false);
    const [campaign, setCampaign] = useState<CampaignData | null>(null);
    const [templatePreview, setTemplatePreview] = useState<string | null>(null);
    const [scheduledDate, setScheduledDate] = useState<Date | null>(null);
    const [isScheduling, setIsScheduling] = useState<boolean>(false);
    const [isSending, setIsSending] = useState<boolean>(false);

    const { data: campaignData, refetch, isLoading } = useSingleCampaignQuery(id);

    // useEffect(() => {
    //     if (campaignData) {
    //         setCampaign(campaignData.payload as CampaignData);
    //         setTemplatePreview(campaignData?.payload.template?.email_html || null);
    //         setScheduledDate(campaignData?.payload.scheduled_at as any);
    //     }
    // }, [campaignData]);

    useEffect(() => {
        if (campaignData) {
            setCampaign(campaignData.payload as CampaignData);

            // More explicit template preview setting
            const template = campaignData?.payload?.template;
            if (template && template.email_html) {
                console.log('Setting template preview from:', template.email_html);
                setTemplatePreview(template.email_html);
            } else {
                console.log('No template HTML found');
                setTemplatePreview(null);
            }

            setScheduledDate(campaignData?.payload.scheduled_at as any);
        }
    }, [campaignData]);

    const handleButtonClick = (item: string) => {
        switch (item) {
            case "Subject":
                setIsSubjectModalOpen(true);
                break;
            case "Design":
                setCurrentCampaignId(id as string);
                setTimeout(() => {
                    navigate("/app/templates");
                }, 1000);
                break;
            case "Recipients":
                setIsRecipientModalOpen(true);
                break;
            default:
                break;
        }
    };

    const scheduleCampaign = async () => {
        if (scheduledDate) {
            setIsScheduling(true);
            try {
                setCreateCampaignValues({ scheduled_at: scheduledDate.toISOString() });
                setIsCalendarOpen(false);
                await updateCampaign(campaign?.id as string);
                await new Promise((resolve) => setTimeout(resolve, 500));
                refetch();
            } catch (error) {
                console.error(error);
            } finally {
                setIsScheduling(false);
            }
        }
    };

    const sendCampgn = async (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
        setIsSending(true);
        try {
            await sendCampaign(campaign?.id as string);
        } catch (error) {
            console.error(error);
        } finally {
            setIsSending(false);
        }
    };

    const campaignSections = [
        {
            key: 'Recipients',
            title: 'Recipients',
            description: 'The people who receive your campaign',
            icon: <TeamOutlined className="text-blue-600" />,
            buttonText: 'Add recipients',
            color: 'blue'
        },
        {
            key: 'Subject',
            title: 'Subject',
            description: 'Add a subject line for this campaign',
            icon: <EditOutlined className="text-green-600" />,
            buttonText: 'Add subject',
            color: 'green'
        },
        {
            key: 'Design',
            title: 'Design',
            description: 'Create your email content',
            icon: <DesktopOutlined className="text-purple-600" />,
            buttonText: 'Start designing',
            color: 'purple'
        }
    ];

    if (isLoading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <LoadingSpinnerComponent />
            </div>
        );
    }

    return (
        <HelmetProvider>
            <Helmet title={`Campaign ${campaign?.name} - CrabMailer`} />
            <div className="min-h-screen bg-gray-50">
                <div className="max-w-4xl mx-auto p-6">
                    {/* Header */}
                    <div className="mb-8">
                        <div className="flex items-center justify-between mb-4">
                            <div className="flex items-center space-x-4">
                                <Button
                                    type="text"
                                    icon={<ArrowLeftOutlined />}
                                    onClick={() => window.history.back()}
                                    className="text-blue-600 hover:text-blue-700 hover:bg-blue-50"
                                />
                                <div className="flex items-center space-x-3">
                                    <Title level={3} className="!mb-0">
                                        {campaign?.name}
                                    </Title>
                                    <Tag color="orange" className="rounded-full px-3 py-1 border-0">
                                        <EditOutlined className="mr-1" />
                                        Draft
                                    </Tag>
                                </div>
                            </div>
                            <Space size="middle">
                                <Button
                                    type="default"
                                    icon={<SendOutlined />}
                                    onClick={sendCampgn}
                                    loading={isSending}
                                    className="border-gray-300 hover:border-blue-500 hover:text-blue-600"
                                >
                                    Send Now
                                </Button>
                                <Button
                                    type="primary"
                                    icon={<ScheduleOutlined />}
                                    onClick={() => setIsCalendarOpen(true)}
                                    className="bg-blue-600 hover:bg-blue-700 border-0"
                                >
                                    Schedule
                                </Button>
                            </Space>
                        </div>
                    </div>

                    {/* Campaign Sections */}
                    <div className="space-y-6">
                        {/* Sender Section */}
                        <Card className="border-0 shadow-sm rounded-lg">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center space-x-4">
                                    <Avatar
                                        size={48}
                                        className="bg-blue-100 flex items-center justify-center"
                                        icon={<UserOutlined className="text-blue-600" />}
                                    />
                                    <div>
                                        <div className="flex items-center space-x-2 mb-1">
                                            <Title level={5} className="!mb-0">
                                                Sender
                                            </Title>
                                            <Badge status="success" />
                                        </div>
                                        <div className="flex items-center space-x-2">
                                            <Text className="font-medium">
                                                {campaign?.sender_from_name ?? "My Company"}
                                            </Text>
                                            <Tooltip title="Review your sender authentication status">
                                                <Button
                                                    type="link"
                                                    size="small"
                                                    icon={<InfoCircleOutlined />}
                                                    className="text-blue-600 px-0"
                                                >
                                                    Review status
                                                </Button>
                                            </Tooltip>
                                        </div>
                                    </div>
                                </div>
                                <Button
                                    type="default"
                                    icon={<SettingOutlined />}
                                    onClick={() => setIsSenderModalOpen(true)}
                                    className="border-gray-300 hover:border-blue-500 hover:text-blue-600"
                                >
                                    Manage Sender
                                </Button>
                            </div>
                        </Card>

                        {/* Other Sections */}
                        {campaignSections.map((section) => (
                            <Card key={section.key} className="border-0 shadow-sm rounded-lg">
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center space-x-4">
                                        <Avatar
                                            size={48}
                                            className={`bg-${section.color}-100 flex items-center justify-center`}
                                            icon={section.icon}
                                        />
                                        <div>
                                            <Title level={5} className="!mb-1">
                                                {section.title}
                                            </Title>
                                            <Text className="text-gray-600">
                                                {section.description}
                                            </Text>
                                        </div>
                                    </div>
                                    <Button
                                        type="default"
                                        icon={<PlusOutlined />}
                                        onClick={() => handleButtonClick(section.key)}
                                        className="border-gray-300 hover:border-blue-500 hover:text-blue-600"
                                    >
                                        {section.buttonText}
                                    </Button>
                                </div>

                                {/* Template Preview */}
                                {section.key === 'Design' && (
                                    <>
                                        <Divider className="my-4" />
                                        <div>
                                            <div className="flex items-center space-x-2 mb-3">
                                                <EyeOutlined className="text-blue-600" />
                                                <Text className="font-medium">Template Preview</Text>
                                            </div>
                                            {templatePreview ? (
                                                <div className="border rounded-lg p-4 bg-white max-h-96 overflow-auto">
                                                    <iframe
                                                        srcDoc={templatePreview}
                                                        title="Template Preview"
                                                        className="w-full h-64 border-0"
                                                        sandbox="allow-scripts"
                                                    />
                                                </div>
                                            ) : (
                                                <div className="border rounded-lg p-4 bg-gray-50 text-center text-gray-500">
                                                    No template preview available
                                                </div>
                                            )}
                                        </div>
                                    </>
                                )}
                            </Card>
                        ))}
                    </div>

                    {/* Schedule Modal */}
                    <Modal
                        title={
                            <div className="flex items-center space-x-2">
                                <ScheduleOutlined className="text-blue-600" />
                                <span>Schedule Campaign</span>
                            </div>
                        }
                        open={isCalendarOpen}
                        onCancel={() => setIsCalendarOpen(false)}
                        footer={[
                            <Button key="cancel" onClick={() => setIsCalendarOpen(false)}>
                                Cancel
                            </Button>,
                            <Button
                                key="schedule"
                                type="primary"
                                onClick={scheduleCampaign}
                                loading={isScheduling}
                                className="bg-blue-600 hover:bg-blue-700 border-0"
                            >
                                Schedule Campaign
                            </Button>
                        ]}
                        width={500}
                        className="rounded-lg"
                    >
                        <div className="py-4">
                            <Text className="text-gray-600 mb-4 block">
                                Select the date and time when you want to send this campaign
                            </Text>
                            <DatePicker
                                value={scheduledDate ? dayjs(scheduledDate) : null}
                                onChange={(date) => {
                                    if (date) setScheduledDate(date.toDate());
                                }}
                                showTime
                                format="YYYY-MM-DD HH:mm:ss"
                                className="w-full"
                                size="large"
                                placeholder="Select date and time"
                            />
                        </div>
                    </Modal>

                    {/* Component Modals */}
                    <AddCampaignSubjectComponent
                        campaign={campaign}
                        isOpen={isSubjectModalOpen}
                        onClose={() => setIsSubjectModalOpen(false)}
                    />
                    <AddCampaignRecipients
                        campaign={campaign}
                        isOpen={isRecipientModalOpen}
                        onClose={() => setIsRecipientModalOpen(false)}
                    />
                    <AddSenderComponent
                        campaign={campaign}
                        isOpen={isSenderModalOpen}
                        onClose={() => setIsSenderModalOpen(false)}
                    />
                </div>
            </div>
        </HelmetProvider>
    );
};

export default EditCampaignForm;