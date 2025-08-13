import React, { useRef, useState, useEffect } from 'react';
import { Button, Typography, Badge, Tooltip, Spin, message } from 'antd';
import { ArrowLeftOutlined, SendOutlined, SaveOutlined } from '@ant-design/icons';
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import { useLocation, useNavigate } from 'react-router-dom';
import useCampaignStore from '../../campaign/store/campaign.store';
import useTemplateStore from '../../email-templates/store/template.store';
import { useSingleMarketingTemplateQuery } from '../../email-templates/hooks/useMarketingTemplateQuery';
import { useSingleTransactionalTemplateQuery } from '../../email-templates/hooks/useTransactionTemplateQuery';
import SendTestEmail from '../../email-templates/components/sendTestEmail';

const { Title } = Typography;

const RichTextEditor: React.FC = () => {
    const quillRef = useRef<ReactQuill>(null);
    const [autoSaved, setAutoSaved] = useState(false);
    const [isSaving, setIsSaving] = useState(false);
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { currentTemplate, updateTemplate, setCurrentTemplate } = useTemplateStore();
    const navigate = useNavigate();
    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const uuid = queryParams.get('uuid');
    const _type = queryParams.get('type');
    const [editorContent, setEditorContent] = useState(currentTemplate?.email_html || '');
    const { updateCampaign, setCreateCampaignValues, currentCampaignId, clearCurrentCampaignId } = useCampaignStore();

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

    useEffect(() => {
        if (currentTemplate?.email_html) {
            setEditorContent(currentTemplate.email_html);
        }
    }, [currentTemplate]);

    useEffect(() => {
        if (!editorContent) return;

        const debounce = setTimeout(() => {
            saveDesign();
        }, 1000);
        return () => clearTimeout(debounce);
    }, [editorContent]);

    const saveDesign = async () => {
        if (isSaving) return; // Prevent multiple saves

        setIsSaving(true);

        try {
            if (currentCampaignId) {
                setCreateCampaignValues({ template_id: uuid as string });
                await updateCampaign(currentCampaignId);
            }

            if (uuid && currentTemplate) {
                const updatedTemplate = {
                    ...currentTemplate,
                    email_html: editorContent
                };
                await updateTemplate(uuid, updatedTemplate);
                setAutoSaved(true);
                message.success('Content saved successfully!');
                console.log("Design saved to database!");
                setTimeout(() => setAutoSaved(false), 3000);
            } else {
                console.log("UUID or template is missing", { uuid, currentTemplate });
            }
        } catch (error) {
            message.error('Failed to save content');
            console.error('Save error:', error);
        } finally {
            setIsSaving(false);
        }
    };

    const handleChange = (content: string) => {
        setEditorContent(content);
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

    // Enhanced Quill modules with better toolbar
    const modules = {
        toolbar: {
            container: [
                [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
                ['bold', 'italic', 'underline', 'strike'],
                ['blockquote', 'code-block'],
                [{ 'list': 'ordered' }, { 'list': 'bullet' }],
                [{ 'script': 'sub' }, { 'script': 'super' }],
                [{ 'indent': '-1' }, { 'indent': '+1' }],
                [{ 'direction': 'rtl' }],
                [{ 'size': ['small', false, 'large', 'huge'] }],
                [{ 'color': [] }, { 'background': [] }],
                [{ 'font': [] }],
                [{ 'align': [] }],
                ['link', 'image', 'video'],
                ['clean']
            ],
        },
        clipboard: {
            matchVisual: false,
        },
        history: {
            delay: 1000,
            maxStack: 100,
            userOnly: true
        }
    };

    const formats = [
        'header', 'font', 'size',
        'bold', 'italic', 'underline', 'strike', 'blockquote',
        'list', 'bullet', 'indent',
        'link', 'image', 'video',
        'color', 'background',
        'align', 'script',
        'code-block', 'direction'
    ];

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
                        <p className="text-sm text-gray-500 mt-1">Rich Text Editor</p>
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
            <div className="flex-1 bg-gray-50 p-4">
                <div className="h-full bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">


                    <style dangerouslySetInnerHTML={{
                        __html: `
                      .ql-toolbar {
                            border-bottom: 1px solid #e5e7eb !important;
                            background: #f9fafb;
                        }
                        
                        .ql-container {
                            border: none !important;
                            font-size: 14px;
                            line-height: 1.6;
                        }
                        
                        .ql-editor {
                            padding: 24px !important;
                            min-height: calc(100vh - 200px);
                        }
                        
                        .ql-editor:focus {
                            outline: none;
                        }
                        
                        .ql-toolbar .ql-formats {
                            margin-right: 16px;
                        }
                        
                        .ql-snow .ql-tooltip {
                            z-index: 1000;
                        }
                        
                        .ql-snow .ql-picker-label {
                            color: #374151;
                        }
                        
                        .ql-snow .ql-stroke {
                            stroke: #6b7280;
                        }
                        
                        .ql-snow .ql-fill {
                            fill: #6b7280;
                        }
                        
                        .ql-snow .ql-active .ql-stroke {
                            stroke: #2563eb;
                        }
                        
                        .ql-snow .ql-active .ql-fill {
                            fill: #2563eb;
                `
                    }} />


                    <ReactQuill
                        ref={quillRef}
                        value={editorContent}
                        onChange={handleChange}
                        modules={modules}
                        formats={formats}
                        theme="snow"
                        placeholder="Start writing your email content..."
                        style={{ height: 'calc(100vh - 160px)' }}
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

export default RichTextEditor;