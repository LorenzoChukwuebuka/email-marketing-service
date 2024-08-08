import { useEffect, useState } from "react";
import PricingPlans from "../components/billing/planpricecomponent";
import BillingList from "../components/billing/billinglistComponent";

const BillingDashTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<"Billing" | "Plans">("Plans");

    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as "Billing" | "Plans");
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
            <h1 className="text-2xl font-bold">{activeTab}</h1>
            <nav className="flex space-x-8 mt-5  border-b">
                <button
                    className={`py-2 border-b-2 ${activeTab === "Plans"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("Plans")}
                >
                    Plans
                </button>

                <button
                    className={`py-2 border-b-2 ${activeTab === "Billing"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("Billing")}
                >
                    Billing
                </button>
            </nav>


            {activeTab === "Plans" && (
                <>
                    <PricingPlans />
                </>
            )}

            {activeTab === "Billing" && (
                <>
                    <BillingList />
                </>
            )}

        </div>

    </>
};

export default BillingDashTemplate;
