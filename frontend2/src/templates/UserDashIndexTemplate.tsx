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
    const { data: mailData } = useMailCalcQuery()


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
    const handleUpgrade = () => navigate('/app/billing')

    return (
        <>
            <div className=" mt-2 p-6 flex items-center">
                <h2 className="text-2xl font-bold mb-2">
                    Welcome {userName}
                    <span>
                        {mailData?.payload?.plan?.toLowerCase() === 'free' ? (
                            <>
                                <button
                                    className="ml-4 px-3 py-1 text-sm text-blue-700 border-blue-700 border rounded-md transition-colors"
                                    onClick={handleUpgrade}
                                >
                                    Upgrade
                                </button>
                            </>
                        ) : (
                            <span className="ml-2 text-sm font-normal text-gray-600">

                            </span>
                        )}
                    </span>
                </h2>

                <div className="ml-auto space-x-2 text-blue-700 font-semibold ">
                    <span className=" p-4 w-1/3 mr-4 transition-transform transform hover:scale-105 cursor-pointer hover:shadow-lg" onClick={handleSendCampaign}> Create Campaign <i className="bi bi-arrow-up-right-square"></i> </span>
                    <span className=" p-4 w-1/3 mr-4 transition-transform transform hover:scale-105 cursor-pointer hover:shadow-lg" onClick={handleCreateContact}> Add Contact <i className="bi bi-arrow-up-right-square"></i> </span>
                    <span className=" p-4 w-1/3 mr-4 transition-transform transform hover:scale-105 cursor-pointer hover:shadow-lg" onClick={handleCreateEmailTemplate}> Create Template <i className="bi bi-arrow-up-right-square"></i> </span>
                </div>
            </div>

            <RecentCampaigns />
            <ContactsDashboard/>
            <OverviewStats/>
        </>

    );
};

export default UserDashIndexTemplate;