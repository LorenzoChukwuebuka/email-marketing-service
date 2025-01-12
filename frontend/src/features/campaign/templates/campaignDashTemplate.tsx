import { useState } from "react";

import { Helmet, HelmetProvider } from "react-helmet-async";
import useMetadata from "../../../hooks/useMetaData";
import GetAllCampaignComponent from '../components/getAllCampaignComponent';
import GetScheduledCampaignComponent from "../components/getScheduledCampaignComponent";

type TabType = "Campaign" | "Scheduled";

const CampaignDashTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "Campaign" || storedTab === "Scheduled") ? storedTab : "Campaign";
    });

    const handleTabChange = (tab: TabType) => {
        setActiveTab(tab);
        localStorage.setItem("activeTab", tab);
    };

    const metaData = useMetadata("Campaigns")

    return (
        <HelmetProvider>
            <Helmet {...metaData}
                title={activeTab === "Scheduled" ? "Scheduled Campaigns - CrabMailer" : "My Campaigns - CrabMailer"}
            />
            <div className="p-6 max-w-full">
                <nav className="flex space-x-8 border-b">
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Campaign"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("Campaign")}
                    >
                        My Campaigns
                    </button>
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Scheduled"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("Scheduled")}
                    >
                        Scheduled
                    </button>
                </nav>
                {activeTab === "Campaign" && <GetAllCampaignComponent />}
                {activeTab === "Scheduled" && <GetScheduledCampaignComponent />}
            </div>
        </HelmetProvider>
    );
};

export default CampaignDashTemplate;