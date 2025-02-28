import Cookies from "js-cookie";
import React, { useEffect, useState } from "react";
import { useNavigate } from 'react-router-dom';
import RecentCampaigns from "../components/dashboard/recentcampaigns";
import { useMailCalcQuery } from "../hooks/useMailDataQuery";
import ContactsDashboard from "../components/dashboard/contactdash";
import OverviewStats from "../components/dashboard/overviewstats";

interface UserDetails {
    fullname: string;
}
interface CookieData {
    details: UserDetails;
}

const UserDashIndexTemplate: React.FC = () => {
    const [userName, setUserName] = useState<string>("");
    const navigate = useNavigate();
    const { data: mailData } = useMailCalcQuery();

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

    const handleSendCampaign = () => navigate('/app/campaign');
    const handleCreateContact = () => navigate('/app/contacts');
    const handleCreateEmailTemplate = () => navigate('/app/templates');
    const handleUpgrade = () => navigate('/app/billing');

    return (
        <div className="mt-2 p-4 sm:p-6">
            <div className="flex flex-col sm:flex-row items-center sm:justify-between mb-4">
                <h2 className="text-xl sm:text-2xl font-bold text-center sm:text-left">
                    Welcome {userName}
                    {mailData?.payload?.plan?.toLowerCase() === 'free' && (
                        <button
                            className="ml-2 sm:ml-4 px-3 py-1 text-sm text-blue-700 border-blue-700 border rounded-md transition-colors hover:bg-blue-700 hover:text-white"
                            onClick={handleUpgrade}
                        >
                            Upgrade
                        </button>
                    )}
                </h2>
                <div className="flex flex-wrap justify-center sm:justify-end gap-2 mt-4 sm:mt-0">
                    <button className="px-4 py-2 text-sm text-blue-700 font-semibold border border-blue-700 rounded-lg transition-transform transform hover:scale-105 hover:shadow-lg" onClick={handleSendCampaign}>
                        Create Campaign <i className="bi bi-arrow-up-right-square"></i>
                    </button>
                    <button className="px-4 py-2 text-sm text-blue-700 font-semibold border border-blue-700 rounded-lg transition-transform transform hover:scale-105 hover:shadow-lg" onClick={handleCreateContact}>
                        Add Contact <i className="bi bi-arrow-up-right-square"></i>
                    </button>
                    <button className="px-4 py-2 text-sm text-blue-700 font-semibold border border-blue-700 rounded-lg transition-transform transform hover:scale-105 hover:shadow-lg" onClick={handleCreateEmailTemplate}>
                        Create Template <i className="bi bi-arrow-up-right-square"></i>
                    </button>
                </div>
            </div>

            <RecentCampaigns />
            <ContactsDashboard />
            <OverviewStats />
        </div>
    );
};

export default UserDashIndexTemplate;
