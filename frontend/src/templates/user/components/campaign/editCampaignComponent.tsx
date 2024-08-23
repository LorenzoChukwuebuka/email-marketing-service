import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import useCampaignStore, { Campaign } from '../../../../store/userstore/campaignStore';
import { BaseEntity } from '../../../../interface/baseentity.interface';
import AddCampaignSubjectComponent from './addSubjectComponent';
import AddCampaignRecipients from './addRecipientComponent';
import AddSenderComponent from './addSenderComponent';

const EditCampaignForm: React.FC = () => {
    const { id } = useParams<{ id: string }>() as { id: string };
    const { getSingleCampaign, campaignData, resetCampaignData } = useCampaignStore()

    const [isSubjectModalOpen, setIsSubjectModalOpen] = useState<boolean>(false);
    const [isRecipientModalOpen, setIsRecipientModalOpen] = useState<boolean>(false);
    const [isSenderModalOpen, setIsSenderModalOpen] = useState<boolean>(false);

    useEffect(() => {
        // Reset campaign data before fetching the new one
        resetCampaignData();

        const fetchData = async () => {
            await getSingleCampaign(id);
        };

        fetchData();
    }, [id, getSingleCampaign, resetCampaignData]);

    const campaign = campaignData as BaseEntity & Campaign;

    const handleButtonClick = (item: string) => {
        switch (item) {
            case "Subject":
                setIsSubjectModalOpen(true);
                break;
            case "Design":
                console.log("hello");
                break;
            case "Recipients":
                setIsRecipientModalOpen(true);
                break;
            default:
                break;
        }
    };

    return (
        <main className='p-4'>
            {/* The rest of your JSX code remains unchanged */}
            <div className="flex items-center mb-5">
                <button className="text-blue-600 mr-2" onClick={() => window.history.back()}>
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                    </svg>
                </button>
                <h1 className="text-xl font-semibold mt-2 mb-2">{campaign?.name}</h1>
                <span className="ml-2 text-gray-500 text-sm">Draft</span>
                <div className="ml-auto">
                    <button className="bg-white text-gray-700 font-semibold py-2 px-4 border border-gray-300 rounded-md shadow-sm mr-2">
                       Send
                    </button>
                    <button className="bg-black text-white font-semibold py-2 px-4 rounded-md">
                        Schedule
                    </button>
                </div>
            </div>

            {/* Your other components and JSX here */}
            <div className="space-y-4">
                <div className="border rounded-md p-4">
                    <div className="flex items-center justify-between">
                        <div>
                            <div className="flex items-center">
                                <div className="bg-blue-600 text-white rounded-full w-6 h-6 flex items-center justify-center mr-2">i</div>
                                <h2 className="text-lg font-semibold">Sender</h2>
                            </div>
                            <div className="mt-1">
                                <span className="font-medium">My Company</span>
                                <span className="text-blue-600 ml-2 text-sm">Review your sender status</span>
                            </div>
                        </div>
                        <button className="bg-white text-gray-700 font-semibold py-1 px-3 border border-gray-300 rounded-md text-sm" onClick={() => setIsSenderModalOpen(true)}>
                            Manage sender
                        </button>
                    </div>
                </div>

                {['Recipients', 'Subject', 'Design'].map((item, index) => (
                    <div key={index} className="border rounded-md p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <div className="flex items-center">
                                    <div className="bg-gray-200 text-gray-500 rounded-full w-6 h-6 flex items-center justify-center mr-2">+</div>
                                    <h2 className="text-lg font-semibold">{item}</h2>
                                </div>
                                <p className="text-gray-500 mt-1">
                                    {item === 'Recipients' && 'The people who receive your campaign'}
                                    {item === 'Subject' && 'Add a subject line for this campaign.'}
                                    {item === 'Design' && 'Create your email content.'}
                                </p>
                            </div>
                            <button className="bg-white text-gray-700 font-semibold py-1 px-3  border border-gray-300 rounded-md text-sm" onClick={() => handleButtonClick(item)}>
                                {item === 'Recipients' && 'Add recipients'}
                                {item === 'Subject' && 'Add subject'}
                                {item === 'Design' && 'Start designing'}
                            </button>
                        </div>
                    </div>
                ))}

                {/* <div className="flex items-center justify-between mt-4">
                    <h2 className="text-lg font-semibold">Additional settings</h2>
                    <button className="bg-white text-gray-700 font-semibold py-1 px-3 border border-gray-300 rounded-md text-sm">
                        Edit settings
                    </button>
                </div> */}
            </div>

            <AddCampaignSubjectComponent campaign={campaign} isOpen={isSubjectModalOpen} onClose={() => setIsSubjectModalOpen(false)} />
            <AddCampaignRecipients campaign={campaign} isOpen={isRecipientModalOpen} onClose={() => setIsRecipientModalOpen(false)} />
            <AddSenderComponent campaign={campaign} isOpen={isSenderModalOpen} onClose={() => setIsSenderModalOpen(false)} />
        </main>
    );
};

export default EditCampaignForm;
