import React, { useRef, useState } from 'react';
import { Helmet, HelmetProvider } from "react-helmet-async";
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';

const AdminSendUsersMail: React.FC = () => {
    const quillRef = useRef<ReactQuill>(null);
    const [editorContent, setEditorContent] = useState('');

    const handleChange = (content: string) => {
        setEditorContent(content);
    };

    const handleSend = () => {
        // Implement your send logic here
        console.log("Sending email with content:", editorContent);
    };

    return (
        <HelmetProvider>
            <Helmet title="Email Users" />
            <div className="h-screen flex flex-col p-4">
                <h1 className="text-lg font-semibold mb-4">Send Email to Users</h1>
                <div className="flex-grow mb-4">
                    <ReactQuill
                        ref={quillRef}
                        value={editorContent}
                        onChange={handleChange}
                        modules={{
                            toolbar: [
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
                            clipboard: {
                                matchVisual: false,
                            },
                        }}
                        style={{ height: 'calc(100vh - 12rem)', overflowY: 'auto' }}
                    />
                </div>
                <button
                    className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    onClick={handleSend}
                >
                    Send Email
                </button>
            </div>
        </HelmetProvider>
    );
};

export default AdminSendUsersMail;