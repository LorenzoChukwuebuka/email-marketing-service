import { useMemo } from "react"
import { parseDate } from "../../../utils/utils"
import { useAllCampaignStatsQuery } from "../hooks/useAnalyticsQuery"
import { Empty } from 'antd';

const AnalyticsTableComponent: React.FC = () => {
    const { data: allCampaignStatsData } = useAllCampaignStatsQuery()
    const alcsdata = useMemo(() => allCampaignStatsData?.payload || [], [allCampaignStatsData])
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
                    {alcsdata && alcsdata.length > 0 ? (
                        alcsdata.map((campaign, index) => {
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
                                <Empty description="No Data Available" />
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>




        </div>


    </>
}

export default AnalyticsTableComponent