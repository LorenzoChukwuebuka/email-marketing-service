import { useEffect } from "react"
import useCampaignStore from "../../../../store/userstore/campaignStore"
import { parseDate } from '../../../../utils/utils';

type Props = { campaignId: string }

const CampaignRecipientComponent: React.FC<Props> = ({ campaignId }) => {

    const { campaignRecipientData, getCampaignRecipients } = useCampaignStore()

    useEffect(() => {
        const fetchData = async () => {
            if (campaignId) {
                await getCampaignRecipients(campaignId)
            }

        }

        fetchData()
    }, [getCampaignRecipients])
    return <>

        <div className="overflow-x-auto mt-8">
            <h1 className="font-semibold text-lg mt-4 mb-4"> Campaign Recipients </h1>
            <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                <thead className="bg-gray-50">
                    <tr>

                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Recipient  Email
                        </th>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Sent At
                        </th>

                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Opened At
                        </th>

                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Open Count
                        </th>

                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Clicked At
                        </th>


                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"> Click Count </th>
                    </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                    {campaignRecipientData && campaignRecipientData.length > 0 ? (
                        campaignRecipientData.map((campaign) => {

                            return (
                                <tr key={campaign.uuid} className="hover:bg-gray-100">
                                    <td className="py-4 px-4">{campaign?.recipient_email}</td>
                                    <td className="py-4 px-4">{parseDate(campaign?.sent_at).toLocaleString('en-US', {
                                        timeZone: 'UTC',
                                        year: 'numeric',
                                        month: 'long',
                                        day: 'numeric',
                                        hour: 'numeric',
                                        minute: 'numeric',
                                        second: 'numeric'
                                    })}</td>
                                    <td className="py-4 px-4">{parseDate(campaign?.opened_at as string).toLocaleString('en-US', {
                                        timeZone: 'UTC',
                                        year: 'numeric',
                                        month: 'long',
                                        day: 'numeric',
                                        hour: 'numeric',
                                        minute: 'numeric',
                                        second: 'numeric'
                                    })}</td>
                                    <td className="py-4 px-4">{campaign?.open_count}</td>
                                    <td className="py-4 px-4">{ campaign?.clicked_at && parseDate(campaign?.clicked_at as string).toLocaleString('en-US', {
                                        timeZone: 'UTC',
                                        year: 'numeric',
                                        month: 'long',
                                        day: 'numeric',
                                        hour: 'numeric',
                                        minute: 'numeric',
                                        second: 'numeric'
                                    }) || "N/A"}</td>
                                    <td className="py-4 px-4">{campaign.click_count}</td>
                                </tr>
                            );
                        })

                    ) : (
                        <tr>
                            <td colSpan={7} className="py-4 px-4  text-center">
                                No  recipients available
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>




        </div>

    </>
}

export default CampaignRecipientComponent