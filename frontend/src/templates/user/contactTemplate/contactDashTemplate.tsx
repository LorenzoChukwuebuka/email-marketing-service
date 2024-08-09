import { useEffect, useState } from "react";
import ContactsDashComponent from "../components/contacts/contactDashComponent";
import ContactGroupDash from "../components/contactGroup/contactGroupDashComponent";

const ContactDashTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<"Contact" | "Contact Group">("Contact");

    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as "Contact" | "Contact Group");
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
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Contact"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("Contact")}
                >
                    Contact
                </button>

                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Contact Group"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("Contact Group")}
                >
                    Contact Group
                </button>


            </nav>


            {activeTab === "Contact" && (
                <>
                    <ContactsDashComponent />
                </>
            )}

            {activeTab === "Contact Group" && (
                <>
                    <ContactGroupDash />
                </>
            )}
        </div>


    </>
}

export default ContactDashTemplate