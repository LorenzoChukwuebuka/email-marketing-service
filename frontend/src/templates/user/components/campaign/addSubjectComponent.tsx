import React, { FormEvent, useEffect, useState } from "react"
import { Modal } from "../../../../components"
import useCampaignStore, { Campaign } from "../../../../store/userstore/campaignStore"
import { BaseEntity } from "../../../../interface/baseentity.interface"

interface Props {
    isOpen: boolean
    onClose: () => void
    campaign: (Campaign & BaseEntity) | null
}

interface FormValues {
    subject: string
    preview_text: string
}

const AddCampaignSubjectComponent: React.FC<Props> = ({ isOpen, onClose, campaign }) => {
    const [formValues, setFormValues] = useState<FormValues>({
        subject: "",
        preview_text: ""
    })

    const { updateCampaign, setCreateCampaignValues } = useCampaignStore()

    useEffect(() => {
        if (campaign) {
            setFormValues({
                subject: campaign.subject ?? "",
                preview_text: campaign.preview_text ?? ""
            })
        }
    }, [campaign])

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()
        if (campaign?.uuid) {
            setCreateCampaignValues({
                subject: formValues.subject,
                preview_text: formValues.preview_text
            })
            await updateCampaign(campaign.uuid)
            onClose()
        }
    }

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target
        setFormValues(prev => ({ ...prev, [name]: value }))
    }

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Add Campaign Subject">
            <form onSubmit={handleSubmit}>
                <div className="mb-4">
                    <label htmlFor="subject" className="block text-sm font-medium mb-3 text-gray-700">
                        Subject <span className="text-red-500"> * </span>
                    </label>
                    <input
                        type="text"
                        id="subject"
                        name="subject"
                        value={formValues.subject}
                        onChange={handleInputChange}
                        placeholder="Add a subject..."
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                    <small className="text-center">Subject is what your audience sees in the title of your email</small>
                </div>
                <div className="mb-4">
                    <label htmlFor="preview_text" className="block text-sm font-medium mb-3 text-gray-700">
                        Preview text
                    </label>
                    <input
                        type="text"
                        id="preview_text"
                        name="preview_text"
                        value={formValues.preview_text}
                        onChange={handleInputChange}
                        placeholder="Add a preview text..."
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                    />
                    <small className="text-center">Preview Text tells your audience more about the mail</small>
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
                        disabled={formValues.subject === ""}
                    >
                        Save
                    </button>
                </div>
            </form>
        </Modal>
    )
}

export default AddCampaignSubjectComponent