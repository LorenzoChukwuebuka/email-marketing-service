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


export default EditorCard;