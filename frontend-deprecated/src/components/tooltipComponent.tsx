import React, { useState } from 'react';
import { HelpCircle } from 'lucide-react';

type Props = {
    text: string
    icon?: any
    iconSize?: number
    iconColor?: any
}

const CustomTooltip = ({ text, icon: Icon = HelpCircle, iconSize = 20, iconColor = "currentColor" }: Props) => {
    const [isTooltipVisible, setIsTooltipVisible] = useState(false);

    return (
        <div className="relative inline-block">
            <Icon
                size={iconSize}
                color={iconColor}
                className="cursor-pointer"
                onMouseEnter={() => setIsTooltipVisible(true)}
                onMouseLeave={() => setIsTooltipVisible(false)}
            />
            {isTooltipVisible && (
                <div className="absolute z-10 px-3 py-2 text-sm font-medium text-white bg-gray-900 rounded-lg shadow-sm tooltip dark:bg-gray-700 -top-10 left-1/2 transform -translate-x-1/2">
                    {text}
                    <div className="tooltip-arrow" data-popper-arrow></div>
                </div>
            )}
        </div>
    );
};

export default CustomTooltip;