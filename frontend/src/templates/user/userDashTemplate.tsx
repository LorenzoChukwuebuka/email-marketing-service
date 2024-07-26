import Cookies from "js-cookie";
import React, { useEffect, useState } from "react";

interface UserDetails {
    fullname: string;
}

interface CookieData {
    details: UserDetails;
}

const UserDashboardTemplate: React.FC = () => {
    const [userName, setUserName] = useState<string>("");

    useEffect(() => {
        const cookie = Cookies.get("Cookies");
        if (cookie) {
            try {
                const user: CookieData = JSON.parse(cookie);
                setUserName(user.details.fullname);
            } catch (error) {
                console.error("Failed to parse cookie", error);
            }
        }
    }, [])

    
    return (
        <>
            <div className="bg-white rounded-lg shadow-md p-6">
                <h2 className="text-2xl font-bold mb-4">Welcome {userName}</h2>
            </div>

            <div className="p-6 bg-gray-100">
                <div className="flex justify-between mb-8">
                    <ActionCard
                        title="Send Campaign"
                        description="Create a campaign and send marketing mails to your audience easily"
                        icon="📢"
                    />
                    <ActionCard
                        title="Create Contact"
                        description="Add or upload your contacts to your mailing lists"
                        icon="👤"
                    />
                    <ActionCard
                        title="Create Email Template"
                        description="Start a new email template or pick from an existing one"
                        icon="✉️"
                    />
                </div>
            </div>
        </>
    );
};


interface ActionCardInterface {
    title: string
    description: string
    icon: string
}

const ActionCard: React.FC<ActionCardInterface> = ({ title, description, icon }) => (
    <div className="bg-white rounded-lg shadow p-4 w-1/3 mr-4">
        <div className="flex items-center mb-2">
            <span className="text-2xl mr-2">{icon}</span>
            <h3 className="font-semibold">{title}</h3>
        </div>
        <p className="text-sm text-gray-600">{description}</p>
    </div>
);

export default UserDashboardTemplate;
