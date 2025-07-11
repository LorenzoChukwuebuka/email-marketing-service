import React, { useRef, useState, useEffect } from 'react';
import { Button, Typography, Badge, Tooltip, Spin, message } from 'antd';
import { ArrowLeftOutlined, SendOutlined, SaveOutlined  } from '@ant-design/icons';
import EmailEditor, { EditorRef, EmailEditorProps } from 'react-email-editor';
import { useLocation, useNavigate } from 'react-router-dom';
import useCampaignStore from '../../campaign/store/campaign.store';
import SendTestEmail from '../../email-templates/components/sendTestEmail';
import useTemplateStore from '../../email-templates/store/template.store';
import { useSingleTransactionalTemplateQuery } from '../../email-templates/hooks/useTransactionTemplateQuery';
import { useSingleMarketingTemplateQuery } from '../../email-templates/hooks/useMarketingTemplateQuery';

const { Title } = Typography;

const DragAndDropEditor: React.FC = () => {
    const emailEditorRef = useRef<EditorRef>(null);
    const [autoSaved, setAutoSaved] = useState<boolean>(false);
    const [isSaving, setIsSaving] = useState<boolean>(false);
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { currentTemplate, updateTemplate, setCurrentTemplate } = useTemplateStore();
    const { setCreateCampaignValues, updateCampaign, currentCampaignId, clearCurrentCampaignId } = useCampaignStore();
    const navigate = useNavigate();
    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const uuid = queryParams.get('uuid');
    const _type = queryParams.get('type');

    const transactionalQuery = useSingleTransactionalTemplateQuery(uuid as string);
    const marketingQuery = useSingleMarketingTemplateQuery(uuid as string);

    useEffect(() => {
        if (!uuid) return;

        if (_type === "t") {
            if (transactionalQuery.data) {
                setCurrentTemplate(transactionalQuery.data);
            }
        } else {
            if (marketingQuery.data) {
                setCurrentTemplate(marketingQuery.data);
            }
        }
    }, [uuid, _type, transactionalQuery.data, marketingQuery.data, setCurrentTemplate]);

    useEffect(() => {
        return () => {
            setCurrentTemplate(null);
        };
    }, [setCurrentTemplate]);

    const saveDesign = async () => {
        if (isSaving) return; // Prevent multiple saves
        
        setIsSaving(true);
        
        try {
            if (currentCampaignId) {
                setCreateCampaignValues({ template_id: uuid as string });
                await updateCampaign(currentCampaignId);
            }

            const unlayer = emailEditorRef.current?.editor;

            await new Promise<void>((resolve, reject) => {
                unlayer?.exportHtml(async (data) => {
                    try {
                        const { design, html } = data;
                        if (uuid && currentTemplate) {
                            const updatedTemplate = {
                                ...currentTemplate,
                                email_design: design,
                                email_html: html
                            };
                            await updateTemplate(uuid, updatedTemplate);
                            setAutoSaved(true);
                            message.success('Design saved successfully!');
                            console.log("Design saved to database!");

                            setTimeout(() => setAutoSaved(false), 3000);
                            resolve();
                        } else {
                            console.log("UUID or template is missing", { uuid, currentTemplate });
                            reject(new Error("UUID or template is missing"));
                        }
                    } catch (error) {
                        reject(error);
                    }
                });
            });
        } catch (error) {
            message.error('Failed to save design');
            console.error('Save error:', error);
        } finally {
            setIsSaving(false);
        }
    };

    const onReady: EmailEditorProps['onReady'] = (unlayer) => {
        console.log("Editor is ready");
        if (currentTemplate && currentTemplate.email_design) {
            unlayer.loadDesign(currentTemplate.email_design);
        }
        unlayer.addEventListener('design:updated', saveDesign);
    };

    const handleNavigate = () => {
        if (currentCampaignId) {
            clearCurrentCampaignId();
            navigate("/app/campaign/edit/" + currentCampaignId);
        } else {
            navigate(`/app/templates/${_type === "t" ? "transactional" : "marketing"}`);
        }
    };

    const handleSaveAndExit = async () => {
        await saveDesign();
        handleNavigate();
    };

    const renderSaveStatus = () => {
        if (isSaving) {
            return <Badge status="processing" text="Saving..." />;
        }
        if (autoSaved) {
            return <Badge status="success" text="Auto Saved!" />;
        }
        return null;
    };

    if (!currentTemplate) {
        return (
            <div className="h-screen flex items-center justify-center bg-gray-50">
                <div className="text-center">
                    <Spin size="large" />
                    <p className="mt-4 text-gray-600">Loading template...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="h-screen flex flex-col bg-white">
            {/* Modern Header */}
            <header className="flex items-center justify-between bg-white border-b border-gray-200 px-6 py-4 shadow-sm">
                <div className="flex items-center space-x-4">
                    <Button
                        type="text"
                        icon={<ArrowLeftOutlined />}
                        onClick={handleNavigate}
                        className="flex items-center justify-center hover:bg-gray-100"
                    />
                    <div>
                        <Title level={4} className="!m-0 !text-gray-900">
                            {currentTemplate?.template_name}
                        </Title>
                        <p className="text-sm text-gray-500 mt-1">Drag & Drop Editor</p>
                    </div>
                </div>

                <div className="flex items-center space-x-3">
                    {renderSaveStatus()}
                    
                    <Tooltip title="Send test email">
                        <Button
                            type="default"
                            icon={<SendOutlined />}
                            onClick={() => setIsModalOpen(true)}
                            className="hover:border-blue-400 hover:text-blue-600"
                        >
                            Send Test
                        </Button>
                    </Tooltip>
                    
                    <Button
                        type="primary"
                        icon={<SaveOutlined />}
                        onClick={handleSaveAndExit}
                        loading={isSaving}
                        className="bg-blue-600 hover:bg-blue-700 border-blue-600 hover:border-blue-700"
                    >
                        Save & Exit
                    </Button>
                </div>
            </header>

            {/* Editor Container */}
            <div className="flex-1 bg-gray-50">
                <div className="h-full border border-gray-200 rounded-lg mx-4 my-4 overflow-hidden shadow-sm">
                    <EmailEditor 
                        ref={emailEditorRef} 
                        onReady={onReady} 
                        style={{ height: "calc(100vh - 120px)" }}
                        options={{
                            displayMode: 'email',
                            appearance: {
                                theme: 'modern_light',
                                panels: {
                                    tools: {
                                        dock: 'left'
                                    }
                                }
                            }
                        }}
                    />
                </div>
            </div>

            {/* Modal */}
            <SendTestEmail 
                isOpen={isModalOpen} 
                onClose={() => setIsModalOpen(false)} 
                template_id={uuid as string} 
            />
        </div>
    );
};

export default DragAndDropEditor;