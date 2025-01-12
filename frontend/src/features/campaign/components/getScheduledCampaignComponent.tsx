import { useEffect, useState, useRef, useMemo } from "react";
import CreateCampaignComponent from "./createCampaignComponent";
import { useNavigate } from "react-router-dom";
import EmptyState from "../../../components/emptyStateComponent";
import { parseDate } from "../../../utils/utils";
import { Pagination } from "antd";
import useDebounce from "../../../hooks/useDebounce";
import { useScheduledCampaignQuery } from "../hooks/useCampaignQuery";

const GetScheduledCampaignComponent: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

    const navigate = useNavigate()

    const [activeDropdown, setActiveDropdown] = useState<string | null>(null);
    const dropdownRef = useRef<HTMLDivElement>(null);
    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay
    const { data: scheduleCampaignData, isLoading } = useScheduledCampaignQuery(currentPage, pageSize, debouncedSearchQuery)

    const onPageChange = (page: number, size: number) => {
        setCurrentPage(page);
        setPageSize(size);
    };

    const handleSearchInput = (query: string) => {
        setSearchQuery(query);
    };

    const schCdata = useMemo(() => scheduleCampaignData?.payload.data || [], [scheduleCampaignData])

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
                            onChange={(e) => handleSearchInput(e.target.value)}
                        />
                    </div>
                </div>

                <div className="overflow-x-auto mt-8">
                    {Array.isArray(schCdata) && schCdata.length > 0 ? (
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
                                    {schCdata.map((campaign: any) => (
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
                            <div className="mt-4 flex justify-center items-center mb-5">
                                {/* Pagination */}
                                <Pagination
                                    current={currentPage}
                                    pageSize={pageSize}
                                    total={scheduleCampaignData?.payload?.total_count || 0} // Ensure your API returns a total count
                                    onChange={onPageChange}
                                    showSizeChanger
                                    pageSizeOptions={["10", "20", "50", "100"]}
                                    showTotal={(total) => `Total ${total} Campaigns`}
                                />
                            </div>
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