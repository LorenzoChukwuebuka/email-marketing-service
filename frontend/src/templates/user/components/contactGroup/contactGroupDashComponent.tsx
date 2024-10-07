import { useEffect, useState } from "react";
import useContactGroupStore from "../../../../store/userstore/contactGroupStore"
import GetAllContactGroups from "./getAllContactGroupComponent"
import CreateGroup from "./createGroupComponent";
import useDebounce from "../../../../hooks/useDebounce";

const ContactGroupDash: React.FC = () => {

    const { selectedGroupIds, deleteGroup, searchGroup, getAllGroups } = useContactGroupStore()
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(false)
    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query

    // Debounce the search query
    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay


    const handleSearch = (query: string) => {
        setSearchQuery(query)
    };



    const deleteGrp = async () => {
        const confirmResult = confirm("Do you want to delete group?");

        if (confirmResult) {
            deleteGroup()
        }
    }

    useEffect(() => {
        const fetchG = async () => {
            setIsLoading(true)
            await getAllGroups()
            await new Promise(resolve => setTimeout(resolve, 1000))
            setIsLoading(false)
        }

        fetchG()
    }, [getAllGroups])

    useEffect(() => {
        if (debouncedSearchQuery !== "") {
            searchGroup(debouncedSearchQuery);
        } else {
            getAllGroups(); // Reset to all groups when search query is empty
        }
    }, [debouncedSearchQuery, searchGroup, getAllGroups]);

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
                            Create Group
                        </button>


                        {selectedGroupIds.length > 0 && (
                            <>
                                <button
                                    className="bg-red-200 px-4 py-2 rounded-md transition duration-300"
                                    onClick={() => deleteGrp()}
                                >

                                    <span className="text-red-500"> Delete Group </span>
                                    <i className="bi bi-trash text-red-500"></i>
                                </button>

                            </>

                        )}
                    </div>

                    <div className="ml-3">
                        <input
                            type="text"
                            placeholder="Search..."
                            className="bg-gray-100 px-4 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-300"
                            onChange={(e) => handleSearch(e.target.value)}
                        />
                    </div>

                </div>

                <GetAllContactGroups />

                <CreateGroup isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />

            </>
        )}


    </>
}

export default ContactGroupDash