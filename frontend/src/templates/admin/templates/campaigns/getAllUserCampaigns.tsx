import { Helmet, HelmetProvider } from "react-helmet-async"
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import useAdminUserCamapaignStore from "../../../../store/admin/AdminUserCampaignsStore";
import Pagination from '../../../../components/Pagination';
import EmptyState from '../../../../components/emptyStateComponent';
import { parseDate } from '../../../../utils/utils';
import { useEffect, useRef, useState } from "react";

const AdminUserCampaigns: React.FC = () => {
    const { userid } = useParams();
    const { campaignData, getAllUserCampaign, paginationInfo } = useAdminUserCamapaignStore()
    // Get the search params (query parameters)
    const [searchParams] = useSearchParams();
    const username = searchParams.get('username');
    const navigate = useNavigate()
    const [activeDropdown, setActiveDropdown] = useState<string | null>(null);
    const dropdownRef = useRef<HTMLDivElement>(null);

    const handlePageChange = (newPage: number) => {
        getAllUserCampaign(userid as string, newPage, paginationInfo.page_size);
    }

    const handleSearch = (query: string) => { }

    const deleteCamp = (uuid: string) => {
        console.log(uuid)
    }

    const suspendCampaign = (uuid: string) => {
        console.log(uuid)
    }

    const pageTitle = username
        ? `Campaigns for ${username}`
        : `User Campaigns`;


    useEffect(() => {
        getAllUserCampaign(userid as string)
    }, [getAllUserCampaign, userid])

    return (
        <HelmetProvider>

            <Helmet title={pageTitle} />

            <button
                className="text-blue-600 mr-2"
                onClick={() => window.history.back()}
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M10 19l-7-7m0 0l7-7m-7 7h18"
                    />
                </svg>


            </button>

            <div className="flex justify-between items-center rounded-md p-2 bg-white mt-10">


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
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200">
                                {campaignData.map((campaign: any) => {
                                    const isSent = campaign.sent_at !== null;
                                    return (
                                        <tr key={campaign.uuid} className="hover:bg-gray-100">
                                            <td className="py-4 px-4">{campaign.name}</td>
                                            <td className="py-4 px-4">
                                                <span
                                                    className={`px-2 py-1 rounded-full text-sm font-medium ${campaign.status === 'draft' ? 'bg-gray-200 text-gray-800' :
                                                        campaign.status === 'Failed' ? 'bg-red-200 text-red-800' :
                                                            campaign.status === 'Sent' ? 'bg-green-200 text-green-800' :
                                                                campaign.status === 'Sending' ? 'bg-yellow-200 text-yellow-800' : ''
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
                                                        className="text-blue-600 hover:text-blue-700"
                                                        onClick={() => navigate(`/user/dash/campaign/report/${campaign.uuid}`)}
                                                    >
                                                        View Report
                                                    </button>
                                                ) : (
                                                    <div className="relative">
                                                        <button
                                                            className="text-gray-400 hover:text-gray-600"
                                                            onClick={() => setActiveDropdown(activeDropdown === campaign.uuid ? null : campaign.uuid)}
                                                        >
                                                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                                                <path d="M6 10a2 2 0 11-4 0 2 2 0 014 0zM12 10a2 2 0 11-4 0 2 2 0 014 0zM16 12a2 2 0 100-4 2 2 0 000 4z" />
                                                            </svg>
                                                        </button>
                                                        {activeDropdown === campaign.uuid && (
                                                            <div
                                                                ref={dropdownRef}
                                                                className="absolute right-0 mt-2 w-28 bg-white border border-gray-300 rounded-md shadow-lg z-10"
                                                            >
                                                                <button
                                                                    className="block w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                                                                    onClick={() => navigate(`/zen/dash/users/campaign/${campaign.uuid}`)}
                                                                >
                                                                    View
                                                                </button>
                                                                <button
                                                                    className="block w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                                                                    onClick={() => suspendCampaign(campaign.uuid)}

                                                                >
                                                                    Suspend
                                                                </button>
                                                                <button
                                                                    className="block w-full px-4 py-2 text-sm text-red-700 hover:bg-gray-100"
                                                                    onClick={() => deleteCamp(campaign.uuid)}
                                                                >
                                                                    Delete
                                                                </button>
                                                            </div>
                                                        )}
                                                    </div>
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
                            title="User have not created any Campaign"
                            description="Create campaigns to reach your audience"
                            icon={<i className="bi bi-emoji-frown text-xl"></i>}
                        />
                    </div>
                )}
            </div>

        </HelmetProvider>
    )
}

export default AdminUserCampaigns