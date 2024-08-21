import { FormEvent, useEffect } from "react"
import { Modal } from "../../../../components"
import useCampaignStore, { Campaign } from "../../../../store/userstore/campaignStore"
import { BaseEntity } from "../../../../interface/baseentity.interface"

interface Props {
    isOpen: boolean
    onClose: () => void
    campaign: (Campaign & BaseEntity) | null
}

const AddCampaignSubjectComponent: React.FC<Props> = ({ isOpen, onClose, campaign }) => {

    const { createCampaignValues, setCreateCampaignValues, updateCampaign } = useCampaignStore()

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()
        await updateCampaign(campaign?.uuid as string)
        onClose()
    }

    const handleSubjectChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setCreateCampaignValues({ subject: e.target.value, preview_text: e.target.value });
    }

    const initAddSubject = () => {
        setCreateCampaignValues({ subject: campaign?.subject ?? "", preview_text: campaign?.preview_text ?? "" })
    }

    useEffect(() => {
        if (campaign) {
            initAddSubject()
        }
    }, [campaign])
    return <>

        <Modal isOpen={isOpen} onClose={onClose} title="Add Campaign Subject">
            <form onSubmit={handleSubmit}>
                <div className="mb-4">
                    <label
                        htmlFor="first_name"
                        className="block text-sm font-medium mb-3 text-gray-700"
                    >
                        Subject <span className="text-red-500"> * </span>
                    </label>
                    <input
                        type="text"
                        id="subject"
                        value={createCampaignValues.subject}
                        onChange={handleSubjectChange}
                        placeholder="add a subject..."
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                    <small className="text-center "> Subject is what your audience see in the title of your email</small>
                </div>
                <div className="mb-4">
                    <label
                        htmlFor="first_name"
                        className="block text-sm font-medium mb-3 text-gray-700"
                    >
                        Preview text
                    </label>
                    <input
                        type="text"
                        id="subject"
                        value={createCampaignValues.preview_text}
                        onChange={handleSubjectChange}
                        placeholder="add a preview text..."
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                    <small className="text-center "> Preview Text tells your audience more about the mail </small>
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

        </Modal>

    </>
}

export default AddCampaignSubjectComponent