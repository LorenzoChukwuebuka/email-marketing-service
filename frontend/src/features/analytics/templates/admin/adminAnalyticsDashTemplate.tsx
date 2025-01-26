import { useState } from "react";
import { HelmetProvider, Helmet } from "react-helmet-async";
import SystemMonitorDashboard from "../../components/admin/admindashanalyticscomponent";


type TabType = "System" | "App"

const AdminAnalyticsDashTemplate: React.FC = () => {

    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "System" || storedTab === "App") ? storedTab : "App";
    });

    const handleTabChange = (tab: TabType) => {
        setActiveTab(tab);
        localStorage.setItem("activeTab", tab);
    };
    return (
        <HelmetProvider>

            <Helmet title={(() => {
                // Define a mapping of activeTab values to titles
                const getTitle = () => {
                    switch (activeTab) {
                        case "App":
                            return "App Analytics";
                        case "System":
                            return "System Analytics";

                        default:
                            return "";
                    }
                };

                // Return the computed title
                return getTitle();
            })()} />



            <div className="p-6 max-w-full">
                <nav className="flex space-x-8 border-b">
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "App"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("App")}
                    >
                        App Analytics
                    </button>
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "System"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("System")}
                    >

                        System Analytics

                    </button>


                </nav>

                {activeTab === "App" && <>  d </>}

                {activeTab === "System" && <SystemMonitorDashboard />}




            </div>

        </HelmetProvider>
    )


}


export default AdminAnalyticsDashTemplate