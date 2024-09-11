import { useEffect, useState } from "react";
import CreateCampaignComponent from "./createCampaignComponent";
import useCampaignStore, { Campaign } from "../../../../store/userstore/campaignStore";
import { parseDate } from "../../../../utils/utils";
import { useNavigate } from "react-router-dom";
import { BaseEntity } from "../../../../interface/baseentity.interface";
import Pagination from '../../../../components/Pagination';
import EmptyState from "../../../../components/emptyStateComponent";

const GetAllCampaignComponent: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { getAllCampaigns, campaignData, paginationInfo, deleteCampaign, searchCampaign } = useCampaignStore()
    const [isLoading, setIsLoading] = useState<boolean>(false)
    const navigate = useNavigate()

    useEffect(() => {
        const fetchCampaign = async () => {
            setIsLoading(true)
            await getAllCampaigns()
            await new Promise(resolve => setTimeout(resolve, 1000))
            setIsLoading(false)
        }
        fetchCampaign()
    }, [getAllCampaigns])

    const deleteCamp = async (uuid: string) => {
        const confirmResult = confirm("Do you want to delete campaign?");

        if (confirmResult) {
            await deleteCampaign(uuid)
            await new Promise((resolve) => setTimeout(resolve, 1000))
            await getAllCampaigns()
        } else {
            console.log("Logout canceled");
        }
    }

    const handlePageChange = (newPage: number) => {
        getAllCampaigns(newPage, paginationInfo.page_size);
    };

    const handleSearch = (query: string) => {
        searchCampaign(query);
    };


    return <>

        {isLoading ? (
            <div className="flex items-center justify-center mt-20">
                <span className="loading loading-spinner loading-lg"></span>
            </div>
        ) : (
            <>
                <div className="flex justify-between items-center rounded-md p-2 bg-white mt-10">
                    <div className="space-x-1  h-auto w-full p-2 px-2 ">
                        <button
                            className="bg-gray-300 px-2 py-2 rounded-md transition duration-300"
                            onClick={() => setIsModalOpen(true)}
                        >
                            Create Campaign
                        </button>
                    </div>

                    <div className="ml-3">
                        <input
                            type="text"
                            placeholder="Search..."
                            className="bg-gray-100 px-4 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-300"
                            onChange={(e) => handleSearch(e.target.value)}
                        />
                    </div>
                </div>

                <div className="overflow-x-auto mt-8">
                    {Array.isArray(campaignData) && campaignData.length > 0 ? (
                        <>
                            <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                                <thead className="bg-gray-50">
                                    <tr>
                                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created On</th>
                                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"></th>
                                        <th className="py-3 px-4"></th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-200">
                                    {(campaignData as (BaseEntity & Campaign)[]).map((campaign: any) => {
                                        const isSent = campaign.sent_at !== null;
                                        return (
                                            <tr key={campaign.uuid} className="hover:bg-gray-100">
                                                <td className="py-4 px-4">{campaign.name}</td>
                                                <td className="py-4 px-4">
                                                    <span
                                                        className={`px-2 py- rounded-full text-sm font-medium ${campaign.status === 'draft' ? 'bg-gray-200 text-gray-800' :
                                                            campaign.status === 'Failed' ? 'bg-red-200 text-red-800' :
                                                                campaign.status === 'Sent' ? 'bg-green-200 text-green-800' :
                                                                    campaign.status === 'Processing' ? 'bg-yellow-200 text-yellow-800' : ''
                                                            }`}>
                                                        {campaign.status.charAt(0).toUpperCase() + campaign.status.slice(1)}
                                                    </span>
                                                </td>
                                                <td className="py-4 px-4">
                                                    {parseDate(campaign.created_at).toLocaleString('en-US', {
                                                        timeZone: 'UTC',
                                                        year: 'numeric',
                                                        month: 'long',
                                                        day: 'numeric',
                                                        hour: 'numeric',
                                                        minute: 'numeric',
                                                        second: 'numeric',
                                                    })}
                                                </td>
                                                <td className="py-4 px-4">
                                                    {isSent ? (
                                                        <button
                                                            className="text-gray-800 hover:text-blue-700"
                                                            onClick={() => navigate(`/user/dash/campaign/report/${campaign.uuid}`)}
                                                        >
                                                            <span className="bg-gray-300 rounded-md p-1">View Report</span>
                                                        </button>
                                                    ) : (
                                                        <span className="space-x-5">
                                                            <button
                                                                className="text-gray-400 hover:text-gray-600"
                                                                onClick={() => navigate(`/user/dash/campaign/edit/${campaign.uuid}`)}
                                                            >
                                                                ✏️
                                                            </button>
                                                            <button
                                                                className="text-gray-400 hover:text-gray-600"
                                                                onClick={() => deleteCamp(campaign.uuid)}
                                                            >
                                                                <i className="bi bi-trash text-red-600"></i>
                                                            </button>
                                                        </span>
                                                    )}
                                                </td>
                                            </tr>
                                        );
                                    })}
                                </tbody>
                            </table>
                            <Pagination paginationInfo={paginationInfo} handlePageChange={handlePageChange} item="Campaigns" />
                        </>
                    ) : (
                        <div className="py-4 px-4 text-center">
                            <EmptyState
                                title="You have not created any Campaign"
                                description="Create campaigns to reach your audience"
                                icon={<i className="bi bi-emoji-frown text-xl"></i>}
                            />
                        </div>
                    )}





                </div>
                <CreateCampaignComponent isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />
            </>
        )}


    </>
}

export default GetAllCampaignComponent