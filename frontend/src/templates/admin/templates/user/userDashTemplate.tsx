import { useState } from "react";
import AllUsersTable from "../components/users/allUserTableComponent";
import VerifiedUsersTable from "../components/users/verifiedUsersTableComponent";
import UnVerifiedUsersTable from "../components/users/unverifiedUsersComponent";

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

    </>;
};

export default AdminUserDashTemplate;
