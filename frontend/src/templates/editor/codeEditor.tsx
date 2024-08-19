import React, { useState, useRef, useEffect } from 'react';
import Editor, { OnMount, OnChange, OnValidate } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import useTemplateStore from '../../store/userstore/templateStore';
import { useNavigate, useLocation } from 'react-router-dom';
import SendTestEmail from '../user/components/templates/sendTestEmail';

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
    const [autoSaved, setAutoSaved] = useState<boolean>(false);
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { getSingleMarketingTemplate, getSingleTransactionalTemplate, currentTemplate, updateTemplate, setCurrentTemplate } = useTemplateStore()

    const navigate = useNavigate();
    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const uuid = queryParams.get('uuid');
    const _type = queryParams.get('type');

    useEffect(() => {
        const fetchData = async () => {
            try {
                if (_type === "t") {
                    await getSingleTransactionalTemplate(uuid as string);
                } else {
                    await getSingleMarketingTemplate(uuid as string);
                }
            } catch (error) {
                console.error("Error fetching template data:", error);
            }
        };

        if (uuid) {
            fetchData();
        }
    }, [uuid, _type]);

    useEffect(() => {
        if (currentTemplate && currentTemplate.email_html) {
            setCode(currentTemplate.email_html);
        }
    }, [currentTemplate]);

    useEffect(() => {
        return () => {
            setCurrentTemplate(null)
        }
    }, [])

    const handleEditorDidMount: OnMount = (editor, monaco) => {
        editorRef.current = editor;
    };

    const handleEditorChange: OnChange = (value, event) => {
        setCode(value ?? '');
        saveCode();
    };

    const handleEditorValidation: OnValidate = (markers) => {
        markers.forEach((marker: any) => console.log('Validation issue:', marker.message));
    };

    const saveCode = async (): Promise<void> => {
        if (uuid && currentTemplate) {
            const updatedTemplate = {
                ...currentTemplate,
                email_html: code
            };
            new Promise(resolve => setTimeout(resolve, 3000));
            await updateTemplate(uuid, updatedTemplate);
            setAutoSaved(true);
            console.log("Code saved to database!");

            setTimeout(() => setAutoSaved(false), 3000);
        } else {
            console.log("UUID or template is missing", { uuid, currentTemplate });
        }
    };

    const options: monaco.editor.IStandaloneEditorConstructionOptions = {
        selectOnLineNumbers: true,
        roundedSelection: false,
        readOnly: false,
        cursorStyle: 'line',
        automaticLayout: true,
        minimap: { enabled: false },
    };

    if (!currentTemplate) {
        return <div>Loading template...</div>;
    }

    const testDesign = () => {
        setIsModalOpen(true)
    }

    const handleNavigate = () => {
        if (_type === "t") {
            navigate("/user/dash/templates")
        } else {
            navigate("/user/dash/marketing")
        }
    }

    return (
        <div className="h-screen flex flex-col">
            <header className="flex items-center justify-between bg-gray-100 px-4 h-[5em] py-2">
                <div className="flex items-center">
                    <button className="mr-2 text-gray-600" onClick={() => handleNavigate()}>
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path fillRule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clipRule="evenodd" />
                        </svg>
                    </button>
                    <div>
                        <h1 className="text-sm font-semibold">{currentTemplate?.template_name}</h1>
                    </div>
                </div>
                <div className="flex items-center space-x-2 text-xs">
                    {autoSaved && (
                        <span className="text-green-600 mr-2">Auto Saved!</span>
                    )}
                    <button className="bg-white text-blue-600 border border-blue-300 px-3 py-1 rounded mr-2" onClick={testDesign}>
                        Send Test
                    </button>
                    <button className="bg-navy-900 text-black border-black cursor-pointer text-sm font-semibold px-3 py-1 rounded" onClick={() => { saveCode(); handleNavigate() }}>
                        Save and exit
                    </button>
                </div>
            </header>

            <Editor
                height="100vh"
                defaultLanguage="html"
                value={code}
                theme="vs-dark"
                options={options}
                onMount={handleEditorDidMount}
                onChange={handleEditorChange}
                onValidate={handleEditorValidation}
            />

            <SendTestEmail isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} template_id={uuid as string} />
                
        </div>
    );
}

export default CodeEditor;