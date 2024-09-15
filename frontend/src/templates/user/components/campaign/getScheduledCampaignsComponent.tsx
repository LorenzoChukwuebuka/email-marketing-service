import { useEffect, useState, useRef } from "react";
import CreateCampaignComponent from "./createCampaignComponent";
import useCampaignStore, { Campaign } from "../../../../store/userstore/campaignStore";
import { parseDate } from "../../../../utils/utils";
import { useNavigate } from "react-router-dom";
import { BaseEntity } from "../../../../interface/baseentity.interface";
import Pagination from '../../../../components/Pagination';
import EmptyState from "../../../../components/emptyStateComponent";

const GetScheduledCampaignComponent: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { getScheduledCampaign, scheduledCampaignData, paginationInfo, searchCampaign } = useCampaignStore()
    const navigate = useNavigate()
    const [isLoading, setIsLoading] = useState<boolean>(false)
    const [activeDropdown, setActiveDropdown] = useState<string | null>(null);
    const dropdownRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const fetchCampaign = async () => {
            setIsLoading(true)
            await getScheduledCampaign()
            await new Promise(resolve => setTimeout(resolve, 1000))
            setIsLoading(false)
        }
        fetchCampaign()
    }, [getScheduledCampaign])

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setActiveDropdown(null);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, []);

    const deleteCampaign = async (uuid: string) => {
        console.log(uuid)
        setActiveDropdown(null);
    }

    const handlePageChange = (newPage: number) => {
        getScheduledCampaign(newPage, paginationInfo.page_size);
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
                    {Array.isArray(scheduledCampaignData) && scheduledCampaignData.length > 0 ? (
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
                                    {scheduledCampaignData.map((campaign: any) => (
                                        <tr key={campaign.uuid} className="hover:bg-gray-100">
                                            <td className="py-4 px-4">{campaign.name}</td>
                                            <td className="py-4 px-4">{campaign.status.charAt(0).toUpperCase() + campaign.status.slice(1)}</td>
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
                                                                onClick={() => {
                                                                    navigate(`/user/dash/campaign/edit/${campaign.uuid}`);
                                                                    setActiveDropdown(null);
                                                                }}
                                                            >
                                                                Edit
                                                            </button>
                                                            <button
                                                                className="block w-full px-4 py-2 text-sm text-red-700 hover:bg-gray-100"
                                                                onClick={() => deleteCampaign(campaign.uuid)}
                                                            >
                                                                Delete
                                                            </button>
                                                        </div>
                                                    )}
                                                </div>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                            <Pagination paginationInfo={paginationInfo} handlePageChange={handlePageChange} item="Campaigns" />
                        </>
                    ) : (
                        <div className="py-4 px-4 text-center">
                            <EmptyState
                                title="You have not created any Scheduled Campaign"
                                description="Scheduled campaigns are campaigns that are sent on a later date"
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

export default GetScheduledCampaignComponent