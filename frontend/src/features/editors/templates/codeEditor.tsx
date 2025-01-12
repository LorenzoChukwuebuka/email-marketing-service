import { useState, useRef, useEffect, useCallback } from 'react';
import Editor, { OnMount, OnChange, OnValidate } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import { useNavigate, useLocation } from 'react-router-dom';
import { debounce } from 'lodash';
import SendTestEmail from '../../email-templates/components/sendTestEmail';
import useCampaignStore from '../../campaign/store/campaign.store';
import useTemplateStore from '../../email-templates/store/template.store';
import { useSingleMarketingTemplateQuery } from '../../email-templates/hooks/useMarketingTemplateQuery';
import { useSingleTransactionalTemplateQuery } from '../../email-templates/hooks/useTransactionTemplateQuery';

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
    `

    const [code, setCode] = useState<string>(defaultTemplate);
    const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
    const [saveStatus, setSaveStatus] = useState<'idle' | 'saving' | 'saved' | 'error'>('idle');
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
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

                // Reset status after 3 seconds
                setTimeout(() => setSaveStatus('idle'), 3000);
            }
        } catch (error) {
            setSaveStatus('error');
            setErrorMessage(error instanceof Error ? error.message : 'Failed to save changes');

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
        // Consider integrating validation errors with save status
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

    if (!currentTemplate) {
        return <div>Loading template...</div>;
    }

    return (
        <div className="h-screen flex flex-col">
            <header className="flex items-center justify-between bg-gray-100 px-4 h-[5em] py-2">
                <div className="flex items-center">
                    <button
                        className="mr-2 text-gray-600"
                        onClick={handleNavigate}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path fillRule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clipRule="evenodd" />
                        </svg>
                    </button>
                    <h1 className="text-sm font-semibold">{currentTemplate.template_name}</h1>
                </div>
                <div className="flex items-center space-x-2 text-xs">
                    {saveStatus === 'saving' && (
                        <span className="text-blue-600">Saving...</span>
                    )}
                    {saveStatus === 'saved' && (
                        <span className="text-green-600">Saved!</span>
                    )}
                    {saveStatus === 'error' && (
                        <span className="text-red-600">{errorMessage}</span>
                    )}
                    <button
                        className="bg-white text-blue-600 border border-blue-300 px-3 py-1 rounded"
                        onClick={() => setIsModalOpen(true)}
                    >
                        Send Test
                    </button>
                    <button
                        className="bg-navy-900 text-black border-black font-semibold px-3 py-1 rounded"
                        onClick={handleNavigate}
                    >
                        Save and exit
                    </button>
                </div>
            </header>

            <Editor
                height="100vh"
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
                }}
                onMount={handleEditorDidMount}
                onChange={handleEditorChange}
                onValidate={handleEditorValidation}
            />

            <SendTestEmail
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                template_id={uuid as string}
            />
        </div>
    );
}

export default CodeEditor;