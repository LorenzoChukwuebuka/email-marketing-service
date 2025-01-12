import React, { ButtonHTMLAttributes, useState } from 'react';
import useBillingStore from '../../../../store/userstore/billingStore';

const PurchaseSuccess = () => {
    const [isLoading, setIsLoading] = useState(false);

    const amount = JSON.parse(localStorage.getItem("planSubscription") || '{}').price;
    const planSubbed = JSON.parse(localStorage.getItem("planSubscription") || '{}').planname;
    const gateway = "Paystack";
    const { confirmPayment } = useBillingStore()

    const verifyPayment = async (e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault()
        setIsLoading(true)
        try {
            const params = new URLSearchParams(window.location.search);
            const referenceParam = params.get('reference') as string;
            await confirmPayment(referenceParam, gateway)
        } catch (error) {
            console.log(error)
        } finally {
            setIsLoading(false)
     localStorage.removeItem('planSubscription')
        }

    };

    return (
        <div className="min-h-screen flex flex-col items-center justify-center bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">

            <h1 className='text-center text-2xl font-bold mb-4'> {import.meta.env.VITE_API_NAME} </h1>
            <div className="max-w-md w-full bg-white rounded-lg shadow-md overflow-hidden p-8">
                <div className="text-center">
                    <i className="bi bi-check-circle text-5xl text-green-500"></i>
                    <h2 className="mt-4 text-2xl font-bold text-gray-900">
                        Subscription Payment Successful
                    </h2>
                    <p className="mt-2 text-gray-600">Thank you for your payment!</p>
                </div>
                <div className="mt-6">
                    <h3 className="text-lg font-semibold text-gray-800">
                        Transaction Details
                    </h3>
                    <div className="mt-4 bg-gray-50 p-4 rounded-lg shadow-inner">
                        <div className="flex justify-between mt-2">
                            <span className="font-medium text-gray-700">Amount Paid:</span>
                            <span className="text-gray-900">₦ {amount}</span>
                        </div>
                        <div className="flex justify-between mt-2">
                            <span className="font-medium text-gray-700">
                                Plan Subscribed:
                            </span>
                            <span className="text-gray-900">{planSubbed}</span>
                        </div>
                        <div className="flex justify-between mt-2">
                            <span className="font-medium text-gray-700">Date:</span>
                            <span className="text-gray-900">
                                {new Date().toLocaleDateString()}
                            </span>
                        </div>
                        <div className="flex justify-between mt-2">
                            <span className="font-medium text-gray-700">Payment Method:</span>
                            <span className="text-gray-900">{gateway}</span>
                        </div>
                    </div>
                </div>
                <div className="mt-6 flex justify-between">
                    <button
                        onClick={(e) => verifyPayment(e)}
                        className="inline-flex items-center px-4 py-2 bg-blue-500 text-white text-sm font-medium rounded-md hover:bg-blue-600"
                    >
                        {isLoading ? "Please wait ..." : "Done"}
                    </button>
                </div>
            </div>
        </div>
    );
};

export default PurchaseSuccess;
