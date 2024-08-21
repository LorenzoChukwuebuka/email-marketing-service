import React, { useEffect } from 'react';
import EmptyState from '../../../../components/emptyStateComponent';
import { useNavigate } from 'react-router-dom';
import useCampaignStore from '../../../../store/userstore/campaignStore';
import { parseDate } from '../../../../utils/utils';

const RecentCampaigns = () => {
    const navigate = useNavigate()
    const { getAllCampaigns, campaignData } = useCampaignStore()

    useEffect(() => {
        const fetchData = async () => {
            await getAllCampaigns()
        }
        fetchData()
    }, [])

    return (
        <div className=" mx-auto mt-5 p-6">
            <div className="flex justify-between items-center mb-4">
                <h2 className="text-2xl font-bold mb-4">Recent campaigns</h2>
                <div>
                    <a href="/user/dash/campaign" className="text-blue-600 mr-4 font-semibold">Go to Campaigns</a>

                </div>
            </div>
            <div className="border rounded-lg">
                {campaignData && campaignData.length > 0 ? (<>

                    {campaignData.slice(0, 3).map((campaign, index) => (
                        <div key={campaign.uuid} className="border-b-2 last:border-b-2 p-4">
                            <div className="flex justify-between items-center">
                                <span className="font-medium">
                                    {campaign.name} #{index + 1}
                                </span>
                                <div className="flex items-center space-x-4 text-sm text-gray-600">
                                    <span className="flex items-center">
                                        <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                                        </svg>
                                        Draft
                                    </span>
                                    <span>Last edit {parseDate(campaign.updated_at).toLocaleString('en-US', {
                                        timeZone: 'UTC',
                                        year: 'numeric',
                                        month: 'long',
                                        day: 'numeric',
                                        hour: 'numeric',
                                        minute: 'numeric',
                                        second: 'numeric'
                                    })}</span>
                                </div>
                            </div>
                        </div>
                    ))}
                </>) : (
                    <>
                        <EmptyState className='shadow-md'
                            title="You have not created any Campaign"
                            description="Create and easily send marketing emails to your audience"
                            icon={<i className="bi bi-emoji-frown text-xl"></i>}
                            buttonText="Create Campaigns"
                            onButtonClick={() => navigate("/user/dash/campaign")}
                        />


                    </>)}

            </div>
        </div>
    );
};

export default RecentCampaigns;