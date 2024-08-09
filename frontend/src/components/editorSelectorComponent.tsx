import React from 'react';

interface EditorCardProps {
    icon: React.ReactNode;
    title: string;
    description: string;
    buttonText: string;
    onClick: () => void;
}

const EditorCard: React.FC<EditorCardProps> = ({ icon, title, description, buttonText, onClick }) => (
    <div className="border border-black rounded-lg  p-6 flex flex-col items-start">
        <div className="text-2xl mb-2">{icon}</div>
        <h2 className="text-xl font-semibold mb-2">{title}</h2>
        <p className="text-gray-600 mb-4">{description}</p>
        <button
            onClick={onClick}
            className="mt-auto bg-white text-gray-800 font-semibold py-2 px-4 border border-gray-400 rounded-full hover:bg-gray-100 transition-colors"
        >
            {buttonText}
        </button>
    </div>
);

const EditorSelection: React.FC = () => {
    return (
        <div className="flex mt-8 items-center justify-center space-x-6">
            <EditorCard
                icon={<span role="img" aria-label="document">📄</span>}
                title="Rich Text Editor"
                description="Use the rich text editor to create simple emails"
                buttonText="Use Rich Text Editor"
                onClick={() => console.log("Rich Text Editor clicked")}
            />
            <EditorCard
                icon={<span role="img" aria-label="code">&#60;&#47;&#62;</span>}
                title="HTML Editor"
                description="Copy and paste your HTML code."
                buttonText="Use HTML Editor"
                onClick={() => console.log("HTML Editor clicked")}
            />
        </div>
    );
};

export default EditorSelection;