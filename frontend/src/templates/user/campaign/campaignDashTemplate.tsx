import { useEffect, useState } from "react";
import GetAllCampaignComponent from "../components/campaign/getAllCampaignsComponent";
import GetScheduledCampaignComponent from "../components/campaign/getScheduledCampaignsComponent";

type TabType = "Campaign" | "Scheduled";

const CampaignDashTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "Campaign" || storedTab === "Scheduled") ? storedTab : "Campaign";
    });

    useEffect(() => {
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);

    const handleTabChange = (tab: TabType) => {
        setActiveTab(tab);
        localStorage.setItem("activeTab", tab);
    };

    return (
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
    );
};

export default CampaignDashTemplate;