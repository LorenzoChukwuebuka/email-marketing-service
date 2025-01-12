import React, { useRef, useState, useEffect } from 'react';
import EmailEditor, { EditorRef, EmailEditorProps } from 'react-email-editor';

const TemplateBuilderComponent: React.FC = () => {
    const emailEditorRef = useRef<EditorRef>(null);
    const [savedDesign, setSavedDesign] = useState<any>(null);

    useEffect(() => {
        const storedDesign = localStorage.getItem('emailDesign');
        if (storedDesign) {
            setSavedDesign(JSON.parse(storedDesign));
        }
    }, []);

    const exportHtml = () => {
        const unlayer = emailEditorRef.current?.editor;
        unlayer?.exportHtml((data) => {
            const { design, html } = data;
            console.log('exportHtml', html);
        });
    };

    const saveDesign = () => {
        const unlayer = emailEditorRef.current?.editor;
        unlayer?.exportHtml((data) => {
            const { design } = data;
            localStorage.setItem('emailDesign', JSON.stringify(design));
            setSavedDesign(design);
        });
    };

    const onReady: EmailEditorProps['onReady'] = (unlayer) => {
        if (savedDesign) {
            unlayer.loadDesign(savedDesign);
        }
    };

    return (
        <div className="h-screen flex flex-col  p-4">
            <div className="mb-4 flex gap-4">
                <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" onClick={exportHtml}>
                    Export HTML
                </button>
                <button className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded" onClick={saveDesign}>
                    Save Design
                </button>
            </div>
            <div className="flex-grow">
                <EmailEditor ref={emailEditorRef} onReady={onReady} style={{ height: "100vh" }} />
            </div>
        </div>
    );
};

export default TemplateBuilderComponent;