import { useEffect } from "react"
import useCampaignStore from "../../../store/userstore/campaignStore"

import { parseDate } from "../../../utils/utils"

const AnalyticsTableComponent: React.FC = () => {
    const { getAllCampaignStats, allCampaignStatsData } = useCampaignStore()

    useEffect(() => {
        const fetchData = async () => {
            await getAllCampaignStats()
        }
        fetchData()
    }, [getAllCampaignStats])
    return <>

        <div className="overflow-x-auto mt-8">
            <h1 className="font-semibold text-lg mt-4 mb-4"> Email Campaign Analytics </h1>
            <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                <thead className="bg-gray-50">
                    <tr>

                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Campaign Name
                        </th>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Bounces
                        </th>

                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Recipients
                        </th>

                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Opened
                        </th>

                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Clicked
                        </th>


                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"> Sent Date </th>
                    </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                    {allCampaignStatsData && allCampaignStatsData.length > 0 ? (
                        allCampaignStatsData.map((campaign, index) => {

                            return (
                                <tr key={index} className="hover:bg-gray-100">
                                    <td className="py-4 px-4">{campaign?.name}</td>
                                    <td className="py-4 px-4">{campaign?.bounces}</td>
                                    <td className="py-4 px-4">{campaign.recipients}</td>
                                    <td className="py-4 px-4">{campaign?.opened}</td>
                                    <td className="py-4 px-4">{campaign?.clicked}</td>
                                    <td className="py-4 px-4">{campaign.sent_date && parseDate(campaign?.sent_date as string).toLocaleString('en-US', {
                                        timeZone: 'UTC',
                                        year: 'numeric',
                                        month: 'long',
                                        day: 'numeric',
                                        hour: 'numeric',
                                        minute: 'numeric',
                                        second: 'numeric'
                                    }) || "Not sent"}</td>
                                </tr>
                            );
                        })

                    ) : (
                        <tr>
                            <td colSpan={7} className="py-4 px-4  text-center">
                                No  stats available
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>




        </div>


    </>
}

export default AnalyticsTableComponent