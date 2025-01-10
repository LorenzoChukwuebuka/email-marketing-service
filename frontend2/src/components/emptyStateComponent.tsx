import React from 'react';

interface EmptyStateProps {
    icon?: React.ReactNode;
    title: string;
    description: string;
    buttonText?: string;
    onButtonClick?: () => void;
    className?: string;
}

const EmptyState: React.FC<EmptyStateProps> = ({
    icon,
    title,
    description,
    buttonText,
    onButtonClick,
    className = ''
}) => {
    return (
        <div className={`bg-gray-50 p-8 rounded-lg mb-5 text-center ${className}`}>
            {icon && <div className="mb-4">{icon}</div>}
            <h2 className="text-2xl font-semibold text-gray-700 mb-2">{title}</h2>
            <p className="text-gray-600 mb-6">{description}</p>
            {buttonText && onButtonClick && (
                <button
                    className="bg-blue-600 text-white font-medium py-2 px-6 rounded-md hover:bg-indigo-800 transition duration-300"
                    onClick={onButtonClick}
                >
                    {buttonText}
                </button>
            )}
        </div>
    );
};

export default EmptyState;


// For marketing template
{/* <EmptyState 
  title="No Marketing Template"
  description="Create and easily send marketing templates to your audience"
  buttonText="Create Marketing Template"
  onButtonClick={() => console.log('Create template')}
  icon={<MarketingIcon className="w-12 h-12 mx-auto text-indigo-600" />}
/>

// For empty project list
<EmptyState 
  title="No Projects Yet"
  description="Start creating your first project to get organized"
  buttonText="Create Project"
  onButtonClick={() => console.log('Create project')}
  icon={<FolderIcon className="w-12 h-12 mx-auto text-gray-400" />}
/>

// For empty inbox
<EmptyState 
  title="Your Inbox is Empty"
  description="You're all caught up! No new messages."
  icon={<InboxIcon className="w-12 h-12 mx-auto text-green-500" />}
/> */}