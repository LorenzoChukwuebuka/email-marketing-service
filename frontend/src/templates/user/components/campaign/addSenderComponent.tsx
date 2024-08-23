import { FormEvent, useEffect } from "react";
import { Modal } from "../../../../components"
import { BaseEntity } from "../../../../interface/baseentity.interface";
import useCampaignStore, { Campaign, CampaignGroup } from "../../../../store/userstore/campaignStore";
import Cookies from 'js-cookie'


interface Props {
    isOpen: boolean;
    onClose: () => void;
    campaign: (Campaign & BaseEntity) | null;
}

const AddSenderComponent: React.FC<Props> = ({ isOpen, onClose, campaign }) => {

    const { createCampaignValues, setCreateCampaignValues, updateCampaign } = useCampaignStore()

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()
        await updateCampaign(campaign?.uuid as string)
        onClose()
    }

    const handleSubjectChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setCreateCampaignValues({ sender_from_name: e.target.value });
    }

    const initChangeRecipients = () => {
        let cookie: any = Cookies.get("Cookies");
        let user = JSON.parse(cookie)?.details?.email;
        setCreateCampaignValues({ sender_from_name: campaign?.sender_from_name ?? "", sender: user })
    }

    useEffect(() => {
        if (campaign) {
            initChangeRecipients()
        }
    }, [campaign])

    return <>
        <Modal isOpen={isOpen} onClose={onClose} title="Sender Details">

            <>

                <h1 className="mt-4 text-lg font-semibold mb-4"> Who is sending this email campaign? </h1>

                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label
                            htmlFor="first_name"
                            className="block text-sm font-medium mb-3 text-gray-700"
                        >
                            Email Address <span className="text-red-500"> * </span>
                        </label>
                        <input
                            type="text"
                            id="sender_from_name"
                            value={createCampaignValues.sender}
                            placeholder="add a subject..."
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            required
                            readOnly
                        />

                    </div>
                    <div className="mb-4">
                        <label
                            htmlFor="first_name"
                            className="block text-sm font-medium mb-3 text-gray-700"
                        >
                            Name
                        </label>
                        <input
                            type="text"
                            id="subject"
                            value={createCampaignValues.sender_from_name}
                            onChange={handleSubjectChange}
                            placeholder="add a preview text..."
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            required
                        />

                    </div>


                    <div className="flex justify-end space-x-2">
                        <button
                            type="button"
                            onClick={onClose}
                            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                            disabled={createCampaignValues.subject === ""}
                        >
                            Save
                        </button>
                    </div>
                </form>

            </>
        </Modal>
    </>
}

export default AddSenderComponent