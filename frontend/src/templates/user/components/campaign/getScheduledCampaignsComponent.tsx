import { useState } from "react";
import CreateCampaignComponent from "./createCampaignComponent";

const GetScheduledCampaignComponent: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

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
        <CreateCampaignComponent isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />
    </>
}

export default GetScheduledCampaignComponent