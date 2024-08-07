import { useEffect, useState } from 'react';
import { Check } from 'lucide-react';
import useDailyUserMailSentCalc from '../../../../store/userstore/userDashStore';
import userPlanStore from '../../../../store/userstore/planStore';
import { PlanData } from '../../../../store/admin/planStore';
import PaymentComponent from './paymentcomponent';
import useBillingStore from '../../../../store/userstore/billingStore';

const PricingPlans = () => {
    const { mailData } = useDailyUserMailSentCalc()
    const { fetchPlans, planData } = userPlanStore()
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [selectedPlan, setSelectedPlan] = useState<PlanData | null>(null);
    const { setPaymentValues, initializePayment } = useBillingStore()
    const [loadingPlanId, setLoadingPlanId] = useState<string | null>(null); // State to track the loading plan

    let currentPlan = mailData?.plan

    const handlePlanSelection = async (plan: PlanData) => {
        setSelectedPlan(plan);
        setPaymentValues({
            plan_id: plan.uuid,
            duration: plan.duration,
            amount_to_pay: plan.price,
            payment_method: "Paystack"
        });

        localStorage.setItem('planSubscription', JSON.stringify(plan));
        setLoadingPlanId(plan.uuid); // Set the loading state for this specific plan

        try {
            await initializePayment();
        } finally {
            setLoadingPlanId(null); // Reset loading state after payment initialization
        }
    };

    useEffect(() => {
        fetchPlans()
    }, [fetchPlans])

    const renderPlanCard = (plan: PlanData) => {
        const isCurrentPlan = currentPlan === plan.planname;
        const isLoading = loadingPlanId === plan.uuid; // Check if this specific plan is loading

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
                    onClick={() => handlePlanSelection(plan)}
                >
                    {
                        isLoading
                            ? 'Please wait...'
                            : (isCurrentPlan
                                ? 'Current Plan'
                                : (plan.price === null
                                    ? 'Contact Us'
                                    : 'Choose Plan'
                                )
                            )
                    }

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

        <>
            <div className="container mx-auto mt-5 p-4">

                <h1 className='text-4xl font-semibold mb-10 text-center'> Upgrade your Marketing Plan </h1>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    {planData.map(renderPlanCard)}
                </div>
            </div>


        </>


    );
};

export default PricingPlans;