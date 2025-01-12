import { useMemo, useState } from "react";
import { Empty, Pagination } from "antd";
import { useBillingQuery } from '../hooks/useBillingQuery';

const BillingList: React.FC = () => {
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const { data: billingData } = useBillingQuery(currentPage, pageSize)

    const onPageChange = (page: number, size: number) => {
        setCurrentPage(page);
        setPageSize(size);
    };
    const bData = useMemo(() => billingData?.payload.data || [], [billingData])
    return (
        <>
            {Array.isArray(bData) && bData ? (
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

                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200 ">
                                {bData && bData.length > 0 ? (
                                    bData.map((billing: any) => (
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

                                        </tr>
                                    ))
                                ) : (
                                    <tr>
                                        <td colSpan={10} className="py-4 px-4 text-center">
                                            <Empty description={"No Billing Data Available"} />
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                    <div className="mt-4 flex justify-center items-center mb-5">
                        {/* Pagination */}
                        <Pagination
                            current={currentPage}
                            pageSize={pageSize}
                            total={billingData?.payload?.total_count || 0} // Ensure your API returns a total count
                            onChange={onPageChange}
                            showSizeChanger
                            pageSizeOptions={["10", "20", "50", "100"]}
                            showTotal={(total) => `Total ${total} billing record(s)`}
                        />
                    </div>


                </>

            ) : (
                <Empty description="No billing lists" />
            )}





        </>
    );

}

export default BillingList