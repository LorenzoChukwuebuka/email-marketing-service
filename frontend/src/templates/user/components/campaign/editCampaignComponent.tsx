import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import useCampaignStore, { CampaignData } from '../../../../store/userstore/campaignStore';
import AddCampaignSubjectComponent from './addSubjectComponent';
import AddCampaignRecipients from './addRecipientComponent';
import AddSenderComponent from './addSenderComponent';
import { Helmet, HelmetProvider } from 'react-helmet-async';

const EditCampaignForm: React.FC = () => {
    const { id } = useParams<{ id: string }>() as { id: string };
    const { getSingleCampaign, campaignData, resetCampaignData, setCurrentCampaignId, setCreateCampaignValues, sendCampaign, updateCampaign, currentCampaignId } = useCampaignStore();
    const navigate = useNavigate();
    const [isSubjectModalOpen, setIsSubjectModalOpen] = useState<boolean>(false);
    const [isRecipientModalOpen, setIsRecipientModalOpen] = useState<boolean>(false);
    const [isSenderModalOpen, setIsSenderModalOpen] = useState<boolean>(false);
    const [isCalendarOpen, setIsCalendarOpen] = useState<boolean>(false);
    const [campaign, setCampaign] = useState<CampaignData | null>(null);
    const [templatePreview, setTemplatePreview] = useState<string | null>(null);
    const [scheduledDate, setScheduledDate] = useState<Date | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(false)

    useEffect(() => {
        resetCampaignData();

        const fetchData = async () => {
            setIsLoading(true)
            await getSingleCampaign(id);
            await new Promise((resolve) => setTimeout(resolve, 500))
            setIsLoading(false)
        };

        fetchData();
    }, [id, getSingleCampaign, resetCampaignData]);

    useEffect(() => {
        if (campaignData) {
            setCampaign(campaignData as CampaignData);
            //@ts-ignore
            setTemplatePreview(campaignData?.template?.email_html || null);
            //@ts-ignore
            setScheduledDate(campaignData?.scheduled_at || null)
        }
    }, [campaignData]);

    const handleButtonClick = (item: string) => {
        switch (item) {
            case "Subject":
                setIsSubjectModalOpen(true);
                break;
            case "Design":
                setCurrentCampaignId(id as string);
                setTimeout(() => {
                    navigate("/user/dash/templates");
                }, 1000);
                break;
            case "Recipients":
                setIsRecipientModalOpen(true);
                break;
            default:
                break;
        }
    };

    const scheduleCampaign = async () => {
        if (scheduledDate) {
            setCreateCampaignValues({ scheduled_at: scheduledDate.toISOString() });
            setIsCalendarOpen(false);

            await updateCampaign(campaign?.uuid as string)
            await new Promise((resolve) => setTimeout(resolve, 500))
            await getSingleCampaign(id);
        }
    };

    const sendCampgn = async (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
        try {
            await sendCampaign(campaign?.uuid as string);
        } catch (error) {
            console.error(error);
        }
    };






    return (
        <HelmetProvider>

            <Helmet title={`Campaign ${campaign?.name} - CrabMailer`} />
            <main className="p-4">
                {/* Header section */}

                {isLoading ? <div className="flex items-center justify-center mt-20"><span className="loading loading-spinner loading-lg"></span></div> : (


                    <>
                        <div className="flex items-center mb-5">
                            <button className="text-blue-600 mr-2" onClick={() => window.history.back()}>
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                                </svg>
                            </button>
                            <h1 className="text-xl font-semibold mt-2 mb-2">{campaign?.name}</h1>
                            <span className="ml-2 text-gray-500 text-sm">Draft</span>
                            <div className="ml-auto">
                                <button className="bg-white text-gray-700 font-semibold py-2 px-4 border border-gray-300 rounded-md shadow-sm mr-2" onClick={(e) => sendCampgn(e)}>
                                    Send
                                </button>
                                <button
                                    className="bg-black text-white font-semibold py-2 px-4 rounded-md"
                                    onClick={() => setIsCalendarOpen(true)}
                                >
                                    Schedule
                                </button>
                            </div>
                        </div>

                        {/* Main content */}
                        <div className="space-y-4">
                            {/* Sender section */}
                            <div className="border rounded-md p-4">
                                <div className="flex items-center justify-between">
                                    <div>
                                        <div className="flex items-center">
                                            <div className="bg-blue-600 text-white rounded-full w-6 h-6 flex items-center justify-center mr-2">i</div>
                                            <h2 className="text-lg font-semibold">Sender</h2>
                                        </div>
                                        <div className="mt-1">
                                            <span className="font-medium">{campaign?.sender_from_name ?? "My Company"}</span>
                                            <span className="text-blue-600 ml-2 text-sm">Review your sender status</span>
                                        </div>
                                    </div>
                                    <button className="bg-white text-gray-700 font-semibold py-1 px-3 border border-gray-300 rounded-md text-sm" onClick={() => setIsSenderModalOpen(true)}>
                                        Manage sender
                                    </button>
                                </div>
                            </div>

                            {/* Other sections */}
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
                                        <button className="bg-white text-gray-700 font-semibold py-1 px-3 border border-gray-300 rounded-md text-sm" onClick={() => handleButtonClick(item)}>
                                            {item === 'Recipients' && 'Add recipients'}
                                            {item === 'Subject' && 'Add subject'}
                                            {item === 'Design' && 'Start designing'}
                                        </button>
                                    </div>

                                    {/* Conditionally render the template preview */}
                                    {item === 'Design' && templatePreview && (
                                        <div className="mt-4">
                                            <h3 className="text-lg font-semibold">Template Preview</h3>
                                            <div className="border p-4 mt-2" dangerouslySetInnerHTML={{ __html: templatePreview }} />
                                        </div>
                                    )}
                                </div>
                            ))}
                        </div>

                        {isCalendarOpen && (
                            <div className="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center">
                                <div className="bg-white p-4 rounded-md shadow-md">
                                    <h3 className="text-lg font-semibold mb-4">Select Date and Time</h3>
                                    <DatePicker
                                        selected={scheduledDate}
                                        onChange={(date: Date | null, event: React.SyntheticEvent<any> | undefined) => {
                                            if (date) setScheduledDate(date);
                                        }}
                                        showTimeSelect
                                        dateFormat="Pp"
                                        className="border rounded-md p-2 w-full"
                                    />
                                    <div className="mt-4 flex justify-end">
                                        <button className="bg-gray-300 text-gray-700 py-1 px-3 rounded-md mr-2" onClick={() => setIsCalendarOpen(false)}>
                                            Cancel
                                        </button>
                                        <button className="bg-blue-600 text-white py-1 px-3 rounded-md" onClick={scheduleCampaign}>
                                            Schedule
                                        </button>
                                    </div>
                                </div>
                            </div>
                        )}

                        {/* Modals */}
                        <AddCampaignSubjectComponent campaign={campaign} isOpen={isSubjectModalOpen} onClose={() => setIsSubjectModalOpen(false)} />
                        <AddCampaignRecipients campaign={campaign} isOpen={isRecipientModalOpen} onClose={() => setIsRecipientModalOpen(false)} />
                        <AddSenderComponent campaign={campaign} isOpen={isSenderModalOpen} onClose={() => setIsSenderModalOpen(false)} />

                    </>
                )}


            </main>

        </HelmetProvider>
    );
};

export default EditCampaignForm;
