import { useEffect, useState } from "react";
import CreateCampaignComponent from "./createCampaignComponent";
import useCampaignStore, { Campaign } from "../../../../store/userstore/campaignStore";
import { parseDate } from "../../../../utils/utils";
import { useNavigate } from "react-router-dom";
import { BaseEntity } from "../../../../interface/baseentity.interface";

const GetAllCampaignComponent: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { getAllCampaigns, campaignData } = useCampaignStore()
    const navigate = useNavigate()

    useEffect(() => {
        const fetchCampaign = async () => {
            await getAllCampaigns()
        }
        fetchCampaign()
    }, [])



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
                // onChange={(e) => handleSearch(e.target.value)}
                />
            </div>
        </div>


        <div className="overflow-x-auto mt-8">
            <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                <thead className="bg-gray-50">
                    <tr>

                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Name
                        </th>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Status
                        </th>

                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Created On
                        </th>

                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">

                        </th>

                        <th className="py-3 px-4"></th>
                    </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                    {campaignData && (campaignData as (BaseEntity & Campaign)[]).length > 0 ? (
                        (campaignData as (BaseEntity & Campaign)[]).map((campaign: any) => (
                            <tr key={campaign.uuid} className="hover:bg-gray-100">

                                <td className="py-4 px-4">{campaign.name}</td>
                                <td className="py-4 px-4"> {campaign.status.charAt(0).toUpperCase() + campaign.status.slice(1)}</td>
                                <td className="py-4 px-4">
                                    {parseDate(campaign.created_at).toLocaleString('en-US', {
                                        timeZone: 'UTC',
                                        year: 'numeric',
                                        month: 'long',
                                        day: 'numeric',
                                        hour: 'numeric',
                                        minute: 'numeric',
                                        second: 'numeric'
                                    })}
                                </td>

                                <td className="py-4 px-4">
                                    <button
                                        className="text-gray-400 hover:text-gray-600"
                                        onClick={() => navigate(`/user/dash/campaign/edit/${campaign.uuid}`)}
                                    >
                                        ✏️
                                    </button>
                                </td>

                            </tr>
                        ))
                    ) : (
                        <tr>
                            <td colSpan={7} className="py-4 px-4  text-center">
                                No campaigns available
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>
        </div>
        <CreateCampaignComponent isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />
    </>
}

export default GetAllCampaignComponent