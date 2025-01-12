import { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import useCampaignStore from "../store/campaign.store";
import { useCampaignQuery } from "../hooks/useCampaignQuery";
import useDebounce from "../../../hooks/useDebounce";
import EmptyState from "../../../components/emptyStateComponent";
import { Modal, Pagination } from 'antd';
import { parseDate } from '../../../../../frontend/src/utils/utils';
import LoadingSpinnerComponent from "../../../components/loadingSpinnerComponent";
import CreateCampaignComponent from './createCampaignComponent';


const GetAllCampaignComponent: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { deleteCampaign } = useCampaignStore()
    const [activeDropdown, setActiveDropdown] = useState<string | null>(null);
    const dropdownRef = useRef<HTMLDivElement>(null);
    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query
    const navigate = useNavigate()

    const [currentPage, setCurrentPage] = useState(1);

    const [pageSize, setPageSize] = useState(20);

    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay

    const { data: CampaignData, isLoading } = useCampaignQuery(currentPage, pageSize, debouncedSearchQuery)

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

    const deleteCamp = async (uuid: string) => {
        Modal.confirm({
            title: "Are you sure?",
            content: "Do you want to delete campaign?",
            okText: "Yes",
            cancelText: "No",
            onOk: async () => {
                await deleteCampaign(uuid);
                await new Promise(resolve => setTimeout(resolve, 1000));
                location.reload()
            },
        });
        setActiveDropdown(null);
    }


    const onPageChange = (page: number, size: number) => {
        setCurrentPage(page);
        setPageSize(size);
    };

    const cData = CampaignData?.payload.data || []

    const handleSearchInput = (query: string) => {
        setSearchQuery(query);
    };

    return <>
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

        {isLoading ? <LoadingSpinnerComponent /> : (
            <div className="overflow-x-auto mt-8">
                {cData && cData.length > 0 ? (
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

                                {cData.map((campaign: any) => {
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
                                                        onClick={() => navigate(`/app/campaign/report/${campaign.uuid}`)}
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
                                                                    onClick={() => navigate(`/app/campaign/edit/${campaign.uuid}`)}
                                                                >
                                                                    Edit
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

                        <div className="mt-4 flex justify-center items-center mb-5">
                            {/* Pagination */}
                            <Pagination
                                current={currentPage}
                                pageSize={pageSize}
                                total={CampaignData?.payload?.total_count || 0} // Ensure your API returns a total count
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
                            title="You have not created any Campaign"
                            description="Create campaigns to reach your audience"
                            icon={<i className="bi bi-emoji-frown text-xl"></i>}
                        />
                    </div>
                )}
            </div>
        )}


        <CreateCampaignComponent isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />
    </>

}


export default GetAllCampaignComponent