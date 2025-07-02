import { useNavigate } from 'react-router-dom';
import { useCampaignQuery } from '../../features/campaign/hooks/useCampaignQuery';
import EmptyState from '../emptyStateComponent';
import { useState } from 'react';
import LoadingSpinnerComponent from '../loadingSpinnerComponent';


const RecentCampaigns = () => {
    const navigate = useNavigate();
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [currentPage, _setCurrentPage] = useState(1);
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [pageSize, _setPageSize] = useState(20);
    const { data: campaignData, isLoading } = useCampaignQuery(currentPage, pageSize);

    return (
        <div className="mx-auto mt-5 p-6">
            <div className="flex justify-between items-center mb-4">
                <h2 className="text-2xl font-bold mb-4">Recent campaigns</h2>
                <div>
                    <a href="/app/campaign" className="text-blue-600 mr-4 font-semibold">
                        Go to Campaigns
                    </a>
                </div>
            </div>
            {isLoading ? <LoadingSpinnerComponent /> : (
                <>
                    <div className="rounded-lg">
                        {campaignData?.payload?.data && campaignData.payload.data.length > 0 ? (
                            <>
                                {/*  eslint-disable-next-line @typescript-eslint/no-explicit-any*/}
                                {campaignData.payload.data.slice(0, 3).map((campaign: any, index: number) => (
                                    <div key={index} className="border-b-2 last:border-b-2 p-4">
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
                                                <span>
                                                    {`Last edit ${formatDate(campaign.updated_at || campaign.created_at)}`}
                                                </span>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </>
                        ) : (
                            <EmptyState
                                className="shadow-md"
                                title="You have not created any Campaign"
                                description="Create and easily send marketing emails to your audience"
                                icon={<i className="bi bi-emoji-frown text-xl"></i>}
                                buttonText="Create Campaigns"
                                onButtonClick={() => navigate("/app/campaign")}
                            />
                        )}
                    </div>
                </>
            )}

        </div>
    );
};

// Helper function to format dates
const formatDate = (dateString?: string): string => {
    if (!dateString) return 'Unknown date';

    // Remove the " WAT" part if it exists
    const cleanDate = dateString.replace(" WAT", "");

    try {
        const date = new Date(cleanDate);
        if (isNaN(date.getTime())) return 'Invalid date';

        return date.toLocaleDateString('en-GB', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        });
    } catch {
        return 'Invalid date';
    }
};

export default RecentCampaigns;