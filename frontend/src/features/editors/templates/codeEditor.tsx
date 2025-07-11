import { useState, useRef, useEffect, useCallback } from 'react';
import { Button, Typography, Spin, message, Badge, Tooltip } from 'antd';
import { ArrowLeftOutlined, SendOutlined, SaveOutlined, } from '@ant-design/icons';
import Editor, { OnMount, OnChange, OnValidate } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import { useNavigate, useLocation } from 'react-router-dom';
import { debounce } from 'lodash';
import SendTestEmail from '../../email-templates/components/sendTestEmail';
import useCampaignStore from '../../campaign/store/campaign.store';
import useTemplateStore from '../../email-templates/store/template.store';
import { useSingleMarketingTemplateQuery } from '../../email-templates/hooks/useMarketingTemplateQuery';
import { useSingleTransactionalTemplateQuery } from '../../email-templates/hooks/useTransactionTemplateQuery';

const { Title } = Typography;

function CodeEditor(): JSX.Element {
    const defaultTemplate = `
    <!doctype html>
    <html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
    <head>
    <!--[if gte mso 9]><xml>
    <o:OfficeDocumentSettings>
    <o:AllowPNG/>
    <o:PixelsPerInch>96</o:PixelsPerInch>
    </o:OfficeDocumentSettings></xml><![endif]-->
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="x-apple-disable-message-reformatting">
    <!--[if !mso]><!--><meta http-equiv="X-UA-Compatible" content="IE=edge"><!--<![endif]-->
    </head>
    <body>
    <!-- paste or write your code here -->
    </body>
    </html>`;

    const [code, setCode] = useState<string>(defaultTemplate);
    const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
    const [saveStatus, setSaveStatus] = useState<'idle' | 'saving' | 'saved' | 'error'>('idle');
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    //@ts-ignore
    const [errorMessage, setErrorMessage] = useState<string | null>(null);

    const {
        currentTemplate,
        updateTemplate,
        setCurrentTemplate
    } = useTemplateStore();

    const {
        updateCampaign,
        setCreateCampaignValues,
        currentCampaignId,
        clearCurrentCampaignId
    } = useCampaignStore();

    const navigate = useNavigate();
    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const uuid = queryParams.get('uuid');
    const _type = queryParams.get('type');

    const transactionalQuery = useSingleTransactionalTemplateQuery(uuid as string);
    const marketingQuery = useSingleMarketingTemplateQuery(uuid as string);

    // Save function with proper error handling
    const saveCode = async (newCode: string): Promise<void> => {
        try {
            setSaveStatus('saving');
            setErrorMessage(null);

            if (currentCampaignId) {
                setCreateCampaignValues({ template_id: uuid as string });
                await updateCampaign(currentCampaignId);
            }

            if (uuid && currentTemplate) {
                const updatedTemplate = {
                    ...currentTemplate,
                    email_html: newCode
                };
                await updateTemplate(uuid, updatedTemplate);
                setSaveStatus('saved');
                message.success('Template saved successfully!');

                // Reset status after 3 seconds
                setTimeout(() => setSaveStatus('idle'), 3000);
            }
        } catch (error) {
            setSaveStatus('error');
            const errorMsg = error instanceof Error ? error.message : 'Failed to save changes';
            setErrorMessage(errorMsg);
            message.error(errorMsg);

            // Reset error after 5 seconds
            setTimeout(() => {
                setSaveStatus('idle');
                setErrorMessage(null);
            }, 5000);
        }
    };

    // Debounced save function - only triggers after 1 second of no changes
    const debouncedSave = useCallback(
        debounce((newCode: string) => saveCode(newCode), 1000),
        [currentTemplate, uuid, currentCampaignId]
    );

    // Cleanup debounce on unmount
    useEffect(() => {
        return () => {
            debouncedSave.cancel();
        };
    }, [debouncedSave]);

    // Template loading effect
    useEffect(() => {
        if (!uuid) return;

        const template = _type === "t"
            ? transactionalQuery.data
            : marketingQuery.data;

        if (template) {
            setCurrentTemplate(template);
        }
    }, [uuid, _type, transactionalQuery.data, marketingQuery.data, setCurrentTemplate]);

    // Code sync effect
    useEffect(() => {
        if (currentTemplate?.email_html) {
            setCode(currentTemplate.email_html);
        }
    }, [currentTemplate]);

    // Cleanup effect
    useEffect(() => {
        return () => {
            setCurrentTemplate(null);
        };
    }, [setCurrentTemplate]);

    const handleEditorDidMount: OnMount = (editor) => {
        editorRef.current = editor;
    };

    const handleEditorChange: OnChange = (value) => {
        const newCode = value ?? '';
        setCode(newCode);
        debouncedSave(newCode);
    };

    const handleEditorValidation: OnValidate = (markers) => {
        markers.forEach((marker: any) => console.log('Validation issue:', marker.message));
    };

    const handleNavigate = () => {
        // Save any pending changes before navigating
        debouncedSave.flush();

        if (currentCampaignId) {
            clearCurrentCampaignId();
            navigate("/app/campaign/edit/" + currentCampaignId);
        } else {
            navigate(`/app/templates/${_type === "t" ? "transactional" : "marketing"}`);
        }
    };

    const renderSaveStatus = () => {
        switch (saveStatus) {
            case 'saving':
                return (
                    <Badge status="processing" text="Saving..." />
                );
            case 'saved':
                return (
                    <Badge status="success" text="Saved!" />
                );
            case 'error':
                return (
                    <Badge status="error" text="Error saving" />
                );
            default:
                return null;
        }
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
                            {currentTemplate.template_name}
                        </Title>
                        <p className="text-sm text-gray-500 mt-1">Code Editor</p>
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
                        onClick={handleNavigate}
                        className="bg-blue-600 hover:bg-blue-700 border-blue-600 hover:border-blue-700"
                    >
                        Save & Exit
                    </Button>
                </div>
            </header>

            {/* Editor Container */}
            <div className="flex-1 relative">
                <Editor
                    height="100%"
                    defaultLanguage="html"
                    value={code}
                    theme="vs-dark"
                    options={{
                        selectOnLineNumbers: true,
                        roundedSelection: false,
                        readOnly: false,
                        cursorStyle: 'line',
                        automaticLayout: true,
                        minimap: { enabled: false },
                        fontSize: 14,
                        lineHeight: 1.5,
                        padding: { top: 16, bottom: 16 },
                        scrollBeyondLastLine: false,
                        smoothScrolling: true,
                        wordWrap: 'on',
                        bracketPairColorization: { enabled: true },
                        guides: {
                            bracketPairs: true,
                            indentation: true,
                        },
                    }}
                    onMount={handleEditorDidMount}
                    onChange={handleEditorChange}
                    onValidate={handleEditorValidation}
                />
            </div>

            {/* Modal */}
            <SendTestEmail
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                template_id={uuid as string}
            />
        </div>
    );
}

export default CodeEditor;