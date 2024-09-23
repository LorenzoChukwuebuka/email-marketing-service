import { useState } from "react";
import AllSupportTicketComponentTable from "../components/support/allSupportTicketsComponent";

type TabType = "Closed" | "Pending" | "All"

const SupportDashTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "All" || storedTab === "Closed" || storedTab === "Pending") ? storedTab : "All";
    });

    const handleTabChange = (tab: TabType) => {
        setActiveTab(tab);
        localStorage.setItem("activeTab", tab);
    };


    return (
        <>

            <div className="p-6 max-w-full">
                <nav className="flex space-x-8 border-b">
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "All"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("All")}
                    >
                        All Tickets
                    </button>
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Pending"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("Pending")}
                    >

                        Pending Tickets

                    </button>

                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Closed"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("Closed")}
                    >
                        Closed Tickets
                    </button>
                </nav>

                {activeTab === "All" && <AllSupportTicketComponentTable />}




            </div>

        </>
    )
}

export default SupportDashTemplate