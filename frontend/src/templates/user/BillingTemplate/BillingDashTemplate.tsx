import { useEffect, useState } from "react";
import PricingPlans from "../components/billing/planpricecomponent";
import BillingList from "../components/billing/billinglistComponent";
import { Helmet, HelmetProvider } from "react-helmet-async";

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


    return <>

        <HelmetProvider>

            <Helmet title={activeTab === "Billing" ? "Billing - CrabMailer" : "Plans - CrabMailer"} />
            <div className="p-6 max-w-full">

                <nav className="flex space-x-8 border-b">
                    <button
                        className={`py-2 border-b-2 text-xl font-semibold ${activeTab === "Plans"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => setActiveTab("Plans")}
                    >
                        Plans
                    </button>

                    <button
                        className={`py-2 border-b-2 text-xl font-semibold ${activeTab === "Billing"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => setActiveTab("Billing")}
                    >
                        Billing History
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

        </HelmetProvider>

    </>
};

export default BillingDashTemplate;
