import { useState, useEffect } from "react";
import {
    ChangePasswordComponent,
    DeleteAccountComponent,
    ProfileInformationComponent,
} from "../components";

type TabType = "Account Details" | "Change Password" | "Delete Account"

const AccountSettingsTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "Account Details" || storedTab === "Change Password" || storedTab === "Delete Account") ? storedTab : "Account Details";
    });

    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as "Account Details" | "Change Password" | "Delete Account");
        }
    }, []);

    useEffect(() => {
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);


    return (
        <>
            <div className="mb-6 p-4 mt-10">
                <nav className="flex space-x-4 mt-5 border-b">
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Account Details"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => setActiveTab("Account Details")}
                    >
                        Account Details
                    </button>

                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Change Password"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => setActiveTab("Change Password")}
                    >
                        Change Password
                    </button>

                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Delete Account"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => setActiveTab("Delete Account")}
                    >
                        Delete Account
                    </button>
                </nav>
            </div>

            {activeTab === "Delete Account" && (
                <>
                    <DeleteAccountComponent />
                </>
            )}

            {activeTab === "Account Details" && (
                <>
                    <ProfileInformationComponent />
                </>
            )}

            {activeTab === "Change Password" && (
                <>
                    {" "}
                    <ChangePasswordComponent />{" "}
                </>
            )}
        </>
    );
};

export default AccountSettingsTemplate;