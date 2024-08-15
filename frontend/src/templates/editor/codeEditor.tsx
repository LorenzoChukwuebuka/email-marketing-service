import React, { useState, useRef, useEffect } from 'react';
import Editor, { OnMount, OnChange, OnValidate } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import useTemplateStore from '../../store/userstore/templateStore';
import { useNavigate } from 'react-router-dom';

function CodeEditor(): JSX.Element {
    const [code, setCode] = useState<string>(
        `<!doctype html>
        <html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
        <head>
        <!--[if gte mso 9]><xml>  <o:OfficeDocumentSettings>    <o:AllowPNG/>    <o:PixelsPerInch>96</o:PixelsPerInch>  </o:OfficeDocumentSettings></xml><![endif]-->
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta name="x-apple-disable-message-reformatting">  
        <!--[if !mso]><!--><meta http-equiv="X-UA-Compatible" content="IE=edge"><!--<![endif]-->
        </head>
        <body>
        <!-- paste or write your code here -->
        </body>
        </html>
        `
    );
    const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
    const [autoSaved, setAutoSaved] = useState<boolean>(false);
    const { getSingleMarketingTemplate, getSingleTransactionalTemplate, currentTemplate, updateTemplate, setCurrentTemplate } = useTemplateStore()

    useEffect(() => {
        // Load saved code from localStorage when component mounts
        const savedCode = localStorage.getItem('savedCode');
        if (savedCode) {
            setCode(savedCode);
        }
    }, []);

    const navigate = useNavigate()

    const handleEditorDidMount: OnMount = (editor, monaco) => {
        editorRef.current = editor;
    };

    const handleEditorChange: OnChange = (value, event) => {
        setCode(value ?? '');
    };

    const handleEditorValidation: OnValidate = (markers) => {
        markers.forEach((marker) => console.log('Validation issue:', marker.message));
    };

    const runCode = (): void => {
        try {
            eval(code);
        } catch (error) {
            console.error('Error executing code:', error);
        }
    };

    const saveCode = (): void => {
        localStorage.setItem('savedCode', code);
        alert('Code saved successfully!');
    };

    const loadCode = (): void => {
        const savedCode = localStorage.getItem('savedCode');
        if (savedCode) {
            setCode(savedCode);
            if (editorRef.current) {
                editorRef.current.setValue(savedCode);
            }
            alert('Code loaded successfully!');
        } else {
            alert('No saved code found!');
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

    return (
        <div className="h-screen flex flex-col">

            <header className="flex items-center justify-between  bg-gray-100 px-2 h-[5em] py-2">
                <div className="flex items-center">
                    <button className="mr-2 text-gray-600" onClick={() => navigate(-1)}>
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path fillRule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clipRule="evenodd" />
                        </svg>
                    </button>
                    <div>
                        <h1 className="text-sm font-semibold"> hello </h1>
                    </div>
                </div>
                <div className="flex items-center space-x-2 text-xs">

                    {autoSaved && (
                        <span className="text-green-600 mr-2">Auto Saved!</span>
                    )}
                    <button className="bg-white text-blue-600 border border-blue-300 px-3 py-1 rounded mr-2">
                        Send Test
                    </button>
                    <button className="bg-navy-900 text-black border-black cursor-pointer px-3 py-1 rounded" onClick={() => saveCode()}>
                        Save
                    </button>
                </div>
            </header>

            <Editor
                style={{ height: "100vh" }}
                defaultLanguage="html"
                value={code}
                theme="vs-dark"
                options={options}
                onMount={handleEditorDidMount}
                onChange={handleEditorChange}
                onValidate={handleEditorValidation}
            />
            {/* <button onClick={runCode}>Run Code</button>
            <button onClick={saveCode}>Save Code</button>
            <button onClick={loadCode}>Load Code</button> */}
        </div>
    );
}

export default CodeEditor;