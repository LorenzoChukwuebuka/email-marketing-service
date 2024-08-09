import { useState } from "react";
import EmptyState from "../../../components/emptyStateComponent";
import { Link, useNavigate } from "react-router-dom";

const MarketingTemplateDash: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

    const navigate = useNavigate()
    
    return <>

        <div className="flex justify-between items-center rounded-md p-2 bg-white mt-10">
            <div className="space-x-1  h-auto w-full p-2 px-2 ">
                <button
                    className="bg-gray-300 px-2 py-2 rounded-md transition duration-300"
                    onClick={() => setIsModalOpen(true)}
                >
                    <Link to="/user/dash/marketing">  Create Marketing Template </Link>
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

        <div className="mt-4 p-2">

            <EmptyState title="You  have not created any Template"
                description="Create a easily send marketing email to your audience"
                icon={<i className="bi bi-emoji-frown text-xl"></i>}
                buttonText="Create Template"
                onButtonClick={() =>  navigate("/user/dash/marketing")}
            />

        </div>


    </>
}

export default MarketingTemplateDash