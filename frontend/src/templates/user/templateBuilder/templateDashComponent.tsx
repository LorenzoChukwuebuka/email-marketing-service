import React, { useRef } from 'react';
import EmailEditor, { EditorRef, EmailEditorProps } from 'react-email-editor';

const TemplateBuilderDashComponent: React.FC = () => {
    const emailEditorRef = useRef<EditorRef>(null);

    const exportHtml = () => {
        const unlayer = emailEditorRef.current?.editor;
        unlayer?.exportHtml((data) => {
            const { design, html } = data;
            console.log('exportHtml', html);
        });
    };

    const onReady: EmailEditorProps['onReady'] = (unlayer) => {
        // Editor is ready
        // You can load your template here;
        // The design JSON can be obtained by calling
        // unlayer.loadDesign(callback) or unlayer.exportHtml(callback)

        // const templateJson = { DESIGN_JSON_GOES_HERE };
        // unlayer.loadDesign(templateJson);
    };

    return (
        <div className="h-screen flex flex-col  p-4">
            <button className="mb-4" onClick={exportHtml}>
                Export HTML
            </button>
            <div className="flex-grow">
                <EmailEditor ref={emailEditorRef} onReady={onReady} style={{ height: "100vh" }} />
            </div>
        </div>
    );
};

export default TemplateBuilderDashComponent;