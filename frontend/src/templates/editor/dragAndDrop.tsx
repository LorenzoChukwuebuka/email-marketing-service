import React, { useRef, useState, useEffect } from 'react';
import EmailEditor, { EditorRef, EmailEditorProps } from 'react-email-editor';

const DragAndDropEditor: React.FC = () => {
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
            {/* <div className="mb-4 flex gap-4">
                <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" onClick={exportHtml}>
                    Export HTML
                </button>
                <button className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded" onClick={saveDesign}>
                    Save Design
                </button>
            </div> */}

            <header className="flex items-center justify-between bg-gray-100 px-4 py-2">
                <div className="flex items-center">
                    <button className="mr-2 text-gray-600">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path fillRule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clipRule="evenodd" />
                        </svg>
                    </button>
                    <div>
                        <h1 className="text-sm font-semibold">Colorado Beasley</h1>
                        <p className="text-xs text-gray-500">ID: 6924 - Draft</p>
                    </div>
                </div>
                <div className="flex items-center text-xs">
                    <span className="text-green-600 mr-2">Auto Saved!</span>
                    <button className="bg-white text-blue-600 border border-blue-300 px-3 py-1 rounded mr-2">
                        Send Test
                    </button>
                    <button className="bg-navy-900 text-white px-3 py-1 rounded">
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