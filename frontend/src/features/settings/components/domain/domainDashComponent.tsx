import { useMemo, useState } from "react";
import useDebounce from "../../../../hooks/useDebounce";
import { useNavigate } from "react-router-dom";
import EmptyState from "../../../../components/emptyStateComponent";
import useDomainStore from "../../store/domain.store";
import { useDomainQuery } from "../../hooks/useDomainQuery";
import { Modal, Pagination } from 'antd';

const DomainDashboardComponent: React.FC = () => {
    const { deleteDomain } = useDomainStore()
    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query
    // Debounce the search query
    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);

    const onPageChange = (page: number, size: number) => {
        setCurrentPage(page);
        setPageSize(size);
    };
    const navigate = useNavigate()
    const { data: domainData, refetch } = useDomainQuery(currentPage, pageSize, debouncedSearchQuery)
    const dData = useMemo(() => domainData?.payload?.data || [], [domainData])
    const handleSearchInput = (query: string) => {
        setSearchQuery(query);
    };

    const handleRedirect = (id: string) => {
        navigate("/app/settings/domain/records/" + id)
    }

    const handleDelete = async (id: string) => {
        Modal.confirm({
            title: "Are you sure?",
            content: "Do you want to delete domain?",
            okText: "Yes",
            cancelText: "No",
            onOk: async () => {
                await deleteDomain(id);
                await new Promise(resolve => setTimeout(resolve, 1000));
                refetch()
            },
        });
    }

    return (
        <>
            <div className="container mx-auto p-6">
                <h1 className="text-2xl font-bold mb-4">Domains</h1>
                <p className="mb-6 text-gray-600">
                    An email domain is the part of your email address that comes after the @ symbol. It helps your recipients recognise your brand and trust your emails. For better deliverability results, domain must be authenticated.
                </p>
                <a href="#" className="text-blue-500 underline mb-4 inline-block">Learn what a DKIM or DMARC is.</a>

                <div className="w-52 mb-4">
                    <input
                        type="text"
                        placeholder="Search domain by name"
                        className="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        onChange={(e) => handleSearchInput(e.target.value)}
                        value={searchQuery}
                    />
                </div>
                {Array.isArray(dData) && dData.length > 0 ? (
                    <>
                        <table className="min-w-full bg-white rounded-lg ">
                            <thead>
                                <tr>
                                    <th className="py-3 px-6 bg-gray-100 text-left text-sm font-semibold text-gray-700">Domain name</th>
                                    <th className="py-3 px-6 bg-gray-100 text-left text-sm font-semibold text-gray-700">Domain status</th>
                                    <th className="py-3 px-6 bg-gray-100 text-left text-sm font-semibold text-gray-700"> </th>
                                    <th className="py-3 px-6 bg-gray-100 text-right text-sm font-semibold text-gray-700">Actions</th>
                                </tr>
                            </thead>
                            <tbody>

                                {(dData as any).map((data) => (
                                    <tr className="border-b" key={data.uuid}>
                                        <td className="py-3 px-6 text-gray-700">{data.domain}</td>
                                        <td className="py-3 px-6 text-red-500 flex items-center">

                                            {data.verified ? (
                                                <>
                                                    <span className="inline-block w-2 h-2 bg-green-500 rounded-full mr-2"></span>
                                                    <span className="text-green-500">Authenticated</span>
                                                </>
                                            ) : (
                                                <>
                                                    <span className="inline-block w-2 h-2 bg-red-500 rounded-full mr-2"></span>
                                                    <span className="text-red-500">Not authenticated</span>
                                                </>
                                            )}

                                        </td>
                                        <td className="py-3 px-6 text-right">

                                            {data.verified ? (<>
                                                <button className="text-blue-500 hover:text-blue-700 mr-2" onClick={() => handleRedirect(data.id)}> View Configuration </button>
                                            </>) : (<>
                                                <button className="text-blue-500 hover:text-blue-700 mr-2" onClick={() => handleRedirect(data.id)}> Authenticate </button>
                                            </>)}

                                        </td>

                                        <td className="py-3 px-6 text-right">

                                            <button className="text-red-500 hover:text-red-700" onClick={() => handleDelete(data.id)}>
                                                <i className="bi bi-trash3-fill"></i>
                                            </button>
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
                                total={domainData?.payload?.total || 0} // Ensure your API returns a total count
                                onChange={onPageChange}
                                showSizeChanger
                                pageSizeOptions={["10", "20", "50", "100"]}
                                showTotal={(total) => `Total ${total} domains`}
                            />
                        </div>
                    </>) : (
                    <EmptyState
                        title="You have not created any Domain"
                        description="Add a domain to customise your email for your audience"
                        icon={<i className="bi bi-emoji-frown text-xl"></i>}
                    />
                )}


            </div>

        </>
    );


};

export default DomainDashboardComponent;
