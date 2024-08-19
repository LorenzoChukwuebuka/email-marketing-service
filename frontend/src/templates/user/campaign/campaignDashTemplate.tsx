import { useEffect, useState } from "react";


const CampaignDashTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<"Campaign" | "Drafts" | "Scheduled" | "Sent">("Campaign");

    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as "Campaign" | "Drafts");
        }
    }, []);

    useEffect(() => {
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);

    useEffect(() => {
        return () => {
            localStorage.removeItem("activeTab");
        };
    }, []);


    return <>
        <div className="p-6 max-w-full">
            <nav className="flex space-x-8  border-b">
                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Campaign"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("Campaign")}
                >
                    My  Campaigns
                </button>

                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Drafts"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("Drafts")}
                >
                    Drafts
                </button>
            </nav>

            {activeTab === "Campaign" && (
                <>
                    hello
                </>
            )}

            {activeTab === "Drafts" && (
                <>
                    hi
                </>
            )}
        </div>
    </>
}

export default CampaignDashTemplate