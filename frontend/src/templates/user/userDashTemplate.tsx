import Cookies from "js-cookie";
import React, { useEffect, useState } from "react";
import { useNavigate } from 'react-router-dom';
import OverviewStats from "./components/dashboard/overviewStatscomponent";
import RecentCampaigns from "./components/dashboard/recentcampaignscomponent";
import ContactsDashboard from "./components/dashboard/contactsComponent";

interface UserDetails {
    fullname: string;
}
interface CookieData {
    details: UserDetails;
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

    const handleSendCampaign = () => navigate('/user/dash/campaign');
    const handleCreateContact = () => navigate('/user/dash/contacts');
    const handleCreateEmailTemplate = () => navigate('/user/dash/templates');

    return (
        <>
            <div className=" mt-2 p-6 flex items-center">
                <h2 className="text-2xl font-bold mb-2">Welcome {userName}</h2>
                <div className="ml-auto space-x-2 text-blue-700 font-semibold ">
                    <span className=" p-4 w-1/3 mr-4 transition-transform transform hover:scale-105 cursor-pointer hover:shadow-lg" onClick={handleSendCampaign}> Create Campaign <i className="bi bi-arrow-up-right-square"></i> </span>
                    <span className=" p-4 w-1/3 mr-4 transition-transform transform hover:scale-105 cursor-pointer hover:shadow-lg" onClick={handleCreateContact}> Add Contact <i className="bi bi-arrow-up-right-square"></i> </span>
                    <span className=" p-4 w-1/3 mr-4 transition-transform transform hover:scale-105 cursor-pointer hover:shadow-lg" onClick={handleCreateEmailTemplate}> Create Template <i className="bi bi-arrow-up-right-square"></i> </span>
                </div>
            </div>
            <OverviewStats />
            <RecentCampaigns />
            <ContactsDashboard />
        </>
    );
};

export default UserDashboardTemplate;