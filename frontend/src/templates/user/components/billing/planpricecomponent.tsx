import React, { useEffect } from 'react';
import { Check } from 'lucide-react';
import useDailyUserMailSentCalc from '../../../../store/userstore/userDashStore';



const plansData = [
    {
        uuid: "808f0576-25b5-44ba-a027-4ac075dc7b99",
        planname: "Free",
        duration: "month",
        price: 0,
        number_of_mails_per_day: "100",
        details: "Perfect for getting started",
        status: "active",
        features: [
            { name: "100 Emails Per Day" },
            { name: "Basic Analytics" },
            { name: "5 Templates" },
            { name: "Email Support" },
        ],
    },
    {
        uuid: "e7d1cc31-72de-4f3a-b9eb-e74048185ef1",
        planname: "Premium",
        duration: "month",
        price: 30000,
        number_of_mails_per_day: "1000",
        details: "For growing businesses",
        status: "active",
        features: [
            { name: "1,000 Emails Per Day" },
            { name: "Advanced Analytics" },
            { name: "Unlimited Templates" },
            { name: "Priority Support" },
            { name: "Custom Domain" },
        ],
    },
    {
        uuid: "a1b2c3d4-e5f6-4a5b-9c3d-2e1f0a9b8c7d",
        planname: "Business",
        duration: "month",
        price: 50000,
        number_of_mails_per_day: "5000",
        details: "For large scale operations",
        status: "active",
        features: [
            { name: "5,000 Emails Per Day" },
            { name: "Enterprise Analytics" },
            { name: "Dedicated Account Manager" },
            { name: "API Access" },
            { name: "99.99% Uptime SLA" },
        ],
    },
    {
        uuid: "9876fedc-ba98-7654-3210-fedcba987654",
        planname: "Enterprise",
        duration: "month",
        price: null,
        number_of_mails_per_day: "Unlimited",
        details: "Custom solutions for your business",
        status: "active",
        features: [
            { name: "Unlimited Emails" },
            { name: "Custom Integration" },
            { name: "24/7 Phone Support" },
            { name: "On-Premise Deployment Option" },
            { name: "Custom Feature Development" },
        ],
    }
];

const PricingPlans = () => {

    const { mailData } = useDailyUserMailSentCalc()

    let currentPlan = mailData?.plan


    const renderPlanCard = (plan: any) => {
        const isCurrentPlan = currentPlan === plan.planname;
        return (
            <div key={plan.uuid} className={`bg-white rounded-lg p-6  ${isCurrentPlan ? 'border-2 border-blue-500' : ''}`}>
                <h2 className="text-xl font-bold mb-2">{plan.planname}</h2>
                <p className="text-3xl font-bold mb-2">
                    {plan.price === null ? 'Custom' : `₦${plan.price.toLocaleString()}`}
                    <span className="text-sm font-normal">/{plan.duration}</span>
                </p>
                <p className="text-gray-600 mb-4">{plan.details}</p>
                <button
                    className={`w-full py-2 rounded-md mb-4 ${isCurrentPlan ? 'bg-gray-300 text-gray-700' : 'bg-blue-600 text-white'}`}
                    disabled={isCurrentPlan}
                >
                    {isCurrentPlan ? 'Current Plan' : (plan.price === null ? 'Contact Us' : 'Choose Plan')}
                </button>
                <ul className="space-y-2">
                    {plan.features.map((feature: any, index: any) => (
                        <li key={index} className="flex items-center">
                            <Check className="text-green-500 mr-2" size={16} />
                            <span>{feature.name}</span>
                        </li>
                    ))}
                </ul>
            </div>
        );
    };

    return (
        <div className="container mx-auto mt-5 p-4">

            <h1 className='text-2xl font-semibold mb-10 text-center'> Upgrade your Marketing Platform </h1>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {plansData.map(renderPlanCard)}
            </div>
        </div>
    );
};

export default PricingPlans;