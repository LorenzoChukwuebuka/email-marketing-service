import Cookies from "js-cookie";
import React, { useEffect, useState } from "react";
import { useNavigate } from 'react-router-dom';
import OverviewStats from "./components/dashboard/overviewStatscomponent";
import RecentCampaigns from "./components/dashboard/recentcampaignscomponent";

interface UserDetails {
    fullname: string;
}

interface CookieData {
    details: UserDetails;
}

interface ActionCardInterface {
    title: string;
    description: string;
    icon: string;
    onClick: () => void;
}

const UserDashboardTemplate: React.FC = () => {
    const [userName, setUserName] = useState<string>("");
    const navigate = useNavigate();

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
    }, []);

    const handleSendCampaign = () => navigate('/send-campaign');
    const handleCreateContact = () => navigate('/user/dash/contacts');
    const handleCreateEmailTemplate = () => navigate('/user/dash/templates');

    return (
        <>
            <div className="bg-white rounded-lg  p-6">
                <h2 className="text-2xl font-bold mb-4">Welcome {userName}</h2>
            </div>

            <div className="p-6 bg-gray-100">
                <div className="flex justify-between mb-5">
                    <ActionCard
                        title="Send Campaign"
                        description="Create a campaign and send marketing mails to your audience easily"
                        icon="📢"
                        onClick={handleSendCampaign}
                    />

                    <ActionCard
                        title="Create Contact"
                        description="Add or upload your contacts to your mailing lists"
                        icon="👤"
                        onClick={handleCreateContact}
                    />

                    <ActionCard
                        title="Create Email Template"
                        description="Start a new email template or pick from an existing one"
                        icon="✉️"
                        onClick={handleCreateEmailTemplate}
                    />
                </div>
            </div>

            <OverviewStats />

            <RecentCampaigns />
        </>
    );
};


const ActionCard: React.FC<ActionCardInterface> = ({ title, description, icon, onClick }) => (
    <div
        className="bg-white rounded-lg cursor-pointer shadow p-4 w-1/3 mr-4 transition-transform transform hover:scale-105 hover:shadow-lg hover:bg-gray-50"
        onClick={onClick}
    >
        <div className="flex items-center mb-2">
            <span className="text-2xl mr-2">{icon}</span>
            <h3 className="font-semibold">{title}</h3>
        </div>
        <p className="text-sm text-gray-600">{description}</p>
    </div>
);


export default UserDashboardTemplate;