import { useState } from "react";
import {
    ChangePasswordComponent,
    DeleteAccountComponent,
    ProfileInformationComponent,
} from "../components";

const AccountSettingsTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<"Account Details" | "Change Password" | "Delete Account">("Account Details");

    console.log("Active Tab:", activeTab);

    return (
        <>
            <div className="mb-6 mt-10">
                <h1 className="text-2xl font-semibold text-base-200">
                    Account Settings
                </h1>
                <nav className="flex space-x-4 mt-5  border-b">
                    <button
                        className={`py-2 border-b-2 ${activeTab === "Account Details"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => setActiveTab("Account Details")}
                    >
                        Account Details
                    </button>

                    <button
                        className={`py-2 border-b-2 ${activeTab === "Change Password"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => setActiveTab("Change Password")}
                    >
                        Change Password
                    </button>

                    <button
                        className={`py-2 border-b-2 ${activeTab === "Delete Account"
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
