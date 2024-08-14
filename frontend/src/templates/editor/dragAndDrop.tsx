import React, { useRef, useState, useEffect } from 'react';
import EmailEditor, { EditorRef, EmailEditorProps } from 'react-email-editor';
import { useLocation } from 'react-router-dom';
import useTemplateStore, { Template } from '../../store/userstore/templateStore';
import { BaseEntity } from '../../interface/baseentity.interface';
import { useNavigate } from 'react-router-dom';

const DragAndDropEditor: React.FC = () => {
    const emailEditorRef = useRef<EditorRef>(null);
    const [autoSaved, setAutoSaved] = useState<boolean>(false);
    const { getSingleMarketingTemplate, getSingleTransactionalTemplate, currentTemplate, updateTemplate, setCurrentTemplate } = useTemplateStore()

    const navigate = useNavigate()

    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const uuid = queryParams.get('uuid');
    const _type = queryParams.get('type')

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
        console.log("Current template updated:", currentTemplate);
    }, [currentTemplate]);

    useEffect(() => {
        return () => {
            setCurrentTemplate(null)
        }
    }, [])

    const saveDesign = () => {
        console.log("saveDesign function called");
        const unlayer = emailEditorRef.current?.editor;

        unlayer?.exportHtml(async (data) => {
            const { design, html } = data;
            if (uuid && currentTemplate) {
                const updatedTemplate = {
                    ...currentTemplate,
                    email_design: design,
                    email_html: html
                };
                await updateTemplate(uuid, updatedTemplate);
                setAutoSaved(true);
                console.log("Design saved to database!");

                setTimeout(() => setAutoSaved(false), 3000);
            } else {
                console.log("UUID or template is missing", { uuid, currentTemplate });
            }
        });
    };

    const onReady: EmailEditorProps['onReady'] = (unlayer) => {
        console.log("Editor is ready");
        if (currentTemplate && currentTemplate.email_design) {
            unlayer.loadDesign(currentTemplate.email_design);
        }
        unlayer.addEventListener('design:updated', saveDesign);
    };

    if (!currentTemplate) {
        return <div>Loading template...</div>;
    }

    return (
        <div className="h-screen flex flex-col p-4">
            <header className="flex items-center justify-between  bg-gray-100 px-4 h-[10em] py-2">
                <div className="flex items-center">
                    <button className="mr-2 text-gray-600" onClick={() => navigate(-1)}>
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
                    <button className="bg-white text-blue-600 border border-blue-300 px-3 py-1 rounded mr-2">
                        Send Test
                    </button>
                    <button className="bg-navy-900 text-black border-black cursor-pointer px-3 py-1 rounded" onClick={() => saveDesign()}>
                        Save
                    </button>
                </div>
            </header>

            <div className="flex-grow">
                <EmailEditor ref={emailEditorRef} onReady={onReady} style={{ height: "100vh" }} />
            </div>
        </div>
    );
};

export default DragAndDropEditor;
