import { useEffect } from "react";
import useBillingStore from "../../../../store/userstore/billingStore";
import Pagination from "../../../../components/Pagination";

const BillingList: React.FC = () => {
    const { fetchBillingData, billingData, paginationInfo } = useBillingStore()

    const handlePageChange = (newPage: number) => {
        fetchBillingData(newPage, paginationInfo.page_size);
    };


    useEffect(() => {
        fetchBillingData()
    }, [fetchBillingData])

    return (
        <>
            <div className="overflow-x-auto mt-8">
                <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                    <thead className="bg-gray-50">
                        <tr>

                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Amount Paid
                            </th>

                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Duration
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Expiry Date
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Reference
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Transaction ID
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Payment Method
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Plan Name
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Number of Mails Per Day
                            </th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200 ">
                        {billingData && billingData.length > 0 ? (
                            billingData.map((billing: any) => (
                                <tr key={billing.uuid}>

                                    <td className="py-4 px-4">{billing.amount_paid}</td>

                                    <td className="py-4 px-4">{billing.duration}</td>
                                    <td className="py-4 px-4">{new Date(billing.expiry_date).toLocaleString('en-US', {
                                        timeZone: 'UTC',
                                        year: 'numeric',
                                        month: 'long',
                                        day: 'numeric',
                                        hour: 'numeric',
                                        minute: 'numeric',
                                        second: 'numeric'
                                    })}</td>
                                    <td className="py-4 px-4">{billing.reference}</td>
                                    <td className="py-4 px-4">{billing.transaction_id}</td>
                                    <td className="py-4 px-4">{billing.payment_method}</td>
                                    <td className="py-4 px-4">{billing.plan.planname}</td>
                                    <td className="py-4 px-4">{billing.plan.number_of_mails_per_day}</td>
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={10} className="py-4 px-4 text-center">
                                    No contacts available
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>

            <Pagination paginationInfo={paginationInfo} handlePageChange={handlePageChange} item="Billing" />
        </>
    );

}

export default BillingList