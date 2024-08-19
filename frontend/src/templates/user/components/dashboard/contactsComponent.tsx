import React, { useEffect } from 'react';
import useContactStore from '../../../../store/userstore/contactStore';

const ContactsDashboard: React.FC = () => {

    const { getContactCount,contactCount } = useContactStore()

    useEffect(() => {
        const fetchCount = async () => {
            await getContactCount()
        }
        fetchCount()
    }, [])


    return (
        <div className="w-full mx-auto mb-10 mt-10 p-4">
            <div className="flex justify-between items-center mb-4">
                <h1 className="text-2xl font-bold">Contacts</h1>

            </div>

            <div className="flex bg-white rounded-lg shadow overflow-hidden">
                <div className="flex-1 p-6 border-r">
                    <div className="flex items-center mb-2">
                        <span className="text-4xl font-bold">{contactCount.total}</span>
                        <span className="ml-2 text-gray-600">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                            </svg>
                        </span>
                    </div>
                    <p className="text-sm text-gray-600">Total contacts</p>
                </div>
                <div className="flex-1 p-6">
                    <div className="flex items-center mb-2">
                        <span className="text-4xl font-bold">{contactCount.recent}</span>
                        <span className="ml-2 text-green-600">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                            </svg>
                        </span>
                    </div>
                    <p className="text-sm text-gray-600">New contacts over the last 30 days</p>
                </div>
            </div>
        </div>
    );
};

export default ContactsDashboard;