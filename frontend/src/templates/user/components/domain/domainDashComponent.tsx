import { useEffect, useState } from "react";
import useDomainStore from "../../../../store/userstore/domainStore";
import Pagination from "../../../../components/Pagination";
import useDebounce from "../../../../hooks/useDebounce";
import { Link, useNavigate } from "react-router-dom";
import EmptyState from "../../../../components/emptyStateComponent";

const DomainDashboardComponent: React.FC = () => {
    const { getAllDomain, domainData, paginationInfo, searchDomain, deleteDomain } = useDomainStore()
    const [isLoading, setIsLoading] = useState<boolean>(false)

    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query

    // Debounce the search query
    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay

    const navigate = useNavigate()

    useEffect(() => {
        const fetchData = async () => {
            setIsLoading(true)
            await getAllDomain()
            await new Promise(resolve => setTimeout(resolve, 1000))
            setIsLoading(false)
        }

        fetchData()
    }, [getAllDomain])


    useEffect(() => {
        if (debouncedSearchQuery !== "") {
            searchDomain(debouncedSearchQuery);
        } else {
            getAllDomain(); // Reset to all Domain when search query is empty
        }
    }, [debouncedSearchQuery, searchDomain, getAllDomain]);


    const handlePageChange = (newPage: number) => {
        getAllDomain(newPage, paginationInfo.page_size)
    }

    const handleSearchInput = (query: string) => {
        setSearchQuery(query);
    };

    const handleRedirect = (id: string) => {
        navigate("/user/dash/settings/domain/records/" + id)
    }

    const handleDelete = async (id: string) => {

        let confirmResult = confirm("Do you want to delete domain?");

        if (confirmResult) {
            await deleteDomain(id)
        }

        await new Promise(resolve => setTimeout(resolve, 1000))
        await getAllDomain()
    }

    return (

        <>
            <div className="container mx-auto p-6">
                <h1 className="text-2xl font-bold mb-4">Domains</h1>
                {/* <p className="mb-6 text-gray-600">
                An email domain is the part of your email address that comes after the @ symbol. It helps your recipients recognise your brand and trust your emails. For better deliverability results, domain must be authenticated.
            </p>
            <a href="#" className="text-blue-500 underline mb-4 inline-block">Learn what a DKIM or DMARC is.</a> */}

                <div className="w-52 mb-4">
                    <input
                        type="text"
                        placeholder="Search domain by name"
                        className="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        onChange={(e) => handleSearchInput(e.target.value)}
                        value={searchQuery}
                    />
                </div>

                {Array.isArray(domainData) && domainData.length > 0 ? (
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

                                {domainData.map((data) => (
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
                                                <button className="text-blue-500 hover:text-blue-700 mr-2" onClick={() => handleRedirect(data.uuid)}> View Configuration </button>
                                            </>) : (<>
                                                <button className="text-blue-500 hover:text-blue-700 mr-2" onClick={() => handleRedirect(data.uuid)}> Authenticate </button>
                                            </>)}

                                        </td>

                                        <td className="py-3 px-6 text-right">

                                            <button className="text-red-500 hover:text-red-700" onClick={() => handleDelete(data.uuid)}>
                                                <i className="bi bi-trash3-fill"></i>
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                        <Pagination paginationInfo={paginationInfo} handlePageChange={handlePageChange} item="Domains" />
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
