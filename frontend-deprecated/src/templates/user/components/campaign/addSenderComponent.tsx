import { FormEvent, useEffect, useState } from "react";
import { Modal } from "../../../../components";
import { BaseEntity } from "../../../../interface/baseentity.interface";
import useCampaignStore, { Campaign } from "../../../../store/userstore/campaignStore";
import useSenderStore, { Sender } from "../../../../store/userstore/senderStore";

interface Props {
    isOpen: boolean;
    onClose: () => void;
    campaign: (Campaign & BaseEntity) | null;
}

const AddSenderComponent: React.FC<Props> = ({ isOpen, onClose, campaign }) => {
    const { createCampaignValues, setCreateCampaignValues, updateCampaign } = useCampaignStore();
    const { getSenders, senderData } = useSenderStore();
    const [selectedSender, setSelectedSender] = useState<Sender | null>(null);

    useEffect(() => {
        const fetchSender = async () => { await getSenders(); };
        fetchSender();
    }, [getSenders]);

    useEffect(() => {
        if (createCampaignValues.sender) {
            const sender = (senderData as Sender[]).find(s => s.email === createCampaignValues.sender);
            if (sender) {
                setSelectedSender(sender);
                setCreateCampaignValues({ ...createCampaignValues, sender_from_name: sender.name });
            }
        }
    }, [createCampaignValues.sender, senderData]);

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        await updateCampaign(campaign?.uuid as string);
        onClose();
    };

    const handleEmailChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const selectedEmail = e.target.value;
        const sender = (senderData as Sender[]).find(s => s.email === selectedEmail);
        if (sender) {
            setSelectedSender(sender);
            setCreateCampaignValues({
                sender: selectedEmail,
                sender_from_name: sender.name
            });
        }
    };

    const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setCreateCampaignValues({ ...createCampaignValues, sender_from_name: e.target.value });
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Sender Details">
            <>
                <h1 className="mt-4 text-lg font-semibold mb-4">Who is sending this email campaign?</h1>

                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label htmlFor="sender_email" className="block text-sm font-medium mb-3 text-gray-700">
                            Email Address <span className="text-red-500">*</span>
                        </label>
                        <select
                            id="sender_email"
                            value={createCampaignValues.sender || ""}
                            onChange={handleEmailChange}
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            required
                        >
                            <option value="" disabled>Select an email...</option>
                            {Array.isArray(senderData) && senderData.map(sender => (
                                <option key={sender.uuid} value={sender.email}>
                                    {sender.email}
                                </option>
                            ))}
                        </select>
                    </div>
                    <div className="mb-4">
                        <label htmlFor="sender_name" className="block text-sm font-medium mb-3 text-gray-700">
                            Name
                        </label>
                        <input
                            type="text"
                            id="sender_name"
                            value={createCampaignValues.sender_from_name || ""}
                            onChange={handleNameChange}
                            placeholder="Enter sender name..."
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
                            disabled={!createCampaignValues.sender || !createCampaignValues.sender_from_name}
                        >
                            Save
                        </button>
                    </div>
                </form>
            </>
        </Modal>
    );
};

export default AddSenderComponent;