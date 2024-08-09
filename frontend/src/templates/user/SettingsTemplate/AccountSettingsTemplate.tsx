import { useState, useEffect } from "react";
import {
    ChangePasswordComponent,
    DeleteAccountComponent,
    ProfileInformationComponent,
} from "../components";

const AccountSettingsTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState
        <"Account Details" | "Change Password" | "Delete Account"
        >("Account Details");

    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as "Account Details" | "Change Password" | "Delete Account");
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

    console.log("Active Tab:", activeTab);

    return (
        <>
            <div className="mb-6 mt-10">
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