import { useState } from "react";
import { Helmet, HelmetProvider } from "react-helmet-async";
import AllUsersTable from "../../components/allUserTableComponent";
import UnVerifiedUsersTable from "../../components/unverifiedUsersComponent";
import VerifiedUsersTable from "../../components/verifiedUsersTableComponent";


type TabType = "All" | "Verified" | "Unverified";
const AdminUserDashTemplate = () => {
    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "All" || storedTab === "Verified" || storedTab === "Unverified") ? storedTab : "All";
    });

    const handleTabChange = (tab: TabType) => {
        setActiveTab(tab);
        localStorage.setItem("activeTab", tab);
    };

    return <>
        <HelmetProvider>

            <Helmet title={(() => {
                // Define a mapping of activeTab values to titles
                const getTitle = () => {
                    switch (activeTab) {
                        case "All":
                            return "All Users";
                        case "Verified":
                            return "Verified Users";
                        case "Unverified":
                            return "Unverified Users";
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
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "All"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("All")}
                    >
                        All Users
                    </button>
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Verified"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("Verified")}
                    >

                        Verified Users

                    </button>

                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Unverified"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("Unverified")}
                    >
                        Unverified Users
                    </button>
                </nav>

                {activeTab === "All" && <AllUsersTable />}
                {activeTab === "Verified" && <VerifiedUsersTable />}
                {activeTab === "Unverified" && <UnVerifiedUsersTable />}  


            </div>

        </HelmetProvider>

    </>;
};

export default AdminUserDashTemplate;
