import React, { useState } from 'react';
import { InfoIcon } from 'lucide-react';

interface AuthenticationStatusProps {
    status: string;
    statusColor: 'green' | 'red' | 'yellow';
    tooltipTitle: string;
    tooltipContent: React.ReactNode;
}

const AuthenticationStatus: React.FC<AuthenticationStatusProps> = ({
    status,
    statusColor,
    tooltipTitle,
    tooltipContent
}) => {
    const [showTooltip, setShowTooltip] = useState(false);

    const colorClasses = {
        green: 'bg-green-100 text-green-800',
        red: 'bg-red-100 text-red-800',
        yellow: 'bg-yellow-100 text-yellow-800'
    };

    return (
        <div className="relative inline-flex items-center space-x-2">
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colorClasses[statusColor]}`}>
                {status}
            </span>
            <div
                className="relative"
                onMouseEnter={() => setShowTooltip(true)}
                onMouseLeave={() => setShowTooltip(false)}
            >
                <InfoIcon size={16} className="text-gray-500 cursor-help" />
                {showTooltip && (
                    <div className="absolute left-1/2 transform -translate-x-1/2 bottom-full mb-2 w-64 p-4 bg-white rounded-lg shadow-lg text-sm text-gray-700 z-10">
                        <p className="font-semibold mb-2">{tooltipTitle}</p>
                        {tooltipContent}
                    </div>
                )}
            </div>
        </div>
    );
};

export default AuthenticationStatus;