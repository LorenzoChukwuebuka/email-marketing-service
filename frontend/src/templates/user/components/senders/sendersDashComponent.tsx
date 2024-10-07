import { useEffect, useState } from "react";
import useSenderStore from "../../../../store/userstore/senderStore";
import EmptyState from "../../../../components/emptyStateComponent";
import EditSenderComponent from "./editSendersComponent";
import { Sender } from "../../../../store/userstore/senderStore";



type Props = {
    email: string;
    name: string;
    dkim: string;
    dkimSigned: boolean;
    dmarc: string;
    verified: boolean;
    verificationText: string;
    onEdit: () => void;
    onDelete: () => void;
}

const EmailCard = ({ email, name, dkim, dkimSigned, dmarc, verified, verificationText, onEdit, onDelete }: Props) => {

    const getBgColorClass = (dkimSigned: boolean, verified: boolean) => {
        if (dkimSigned && verified) {
            return 'bg-green-100 text-green-500';
        } else if (!dkimSigned && !verified) {
            return 'bg-red-100 text-red-500';
        } else if (!dkimSigned && verified) {
            return 'bg-gray-100 text-gray-500';
        } else {
            return 'bg-red-100 text-red-500'; // Default case
        }
    };


    const bgColorClass = getBgColorClass(dkimSigned, verified);


    return (
        <div className="p-6 bg-white shadow rounded-lg mb-4">
            <div className="flex items-start">
                <div className={`w-12 h-12 rounded-full ${bgColorClass} flex items-center justify-center text-xl mr-4`}>
                    <i className={`bi ${verified ? 'bi-person-fill-check' : 'bi-person-fill-x'}`}></i>
                </div>
                <div className="flex-1">
                    <h4 className="font-semibold text-gray-800">{name} <span className="text-gray-600">({email})</span></h4>
                    <p className="text-sm text-gray-600">{verified ? "Verified" : "Unverified"} • <span className="text-blue-600">{verificationText}</span></p>
                    <div className="flex items-center text-sm mt-2">
                        <span className="mr-6">
                            <span className="font-medium">DKIM signature:</span>
                            <span className={dkimSigned ? 'text-green-500' : 'text-yellow-500'}>
                                {dkim}
                            </span>
                        </span>
                        <span>
                            <span className="font-medium">DMARC:</span>
                            <span className={dkimSigned ? 'text-green-500' : 'text-yellow-500'}>
                                {dmarc}
                            </span>
                        </span>
                    </div>
                    <div className="mt-2 text-sm">
                        <button
                            onClick={onEdit}
                            className="text-blue-600 mr-4"
                        >
                            Edit
                        </button>
                        <button
                            onClick={onDelete}
                            className="text-blue-600"
                        >
                            Delete
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

const SendersDashComponent: React.FC = () => {
    const { getSenders, senderData, deleteSender } = useSenderStore();
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [selectedSender, setSelectedSender] = useState<Sender | null>(null);

    useEffect(() => {
        const fetchSender = async () => {
            await getSenders();
        };
        fetchSender();
    }, [getSenders]);



    const handleEdit = (sender: Sender) => {
        setSelectedSender(sender);
        setIsEditModalOpen(true);
    };

    const handleDelete = async (id: string) => {
        let confirmResult = confirm("Do you want to delete contact(s)?");

        if (confirmResult) {
            await deleteSender(id);
        }

        await new Promise(resolve => setTimeout(resolve, 2000));

        await getSenders();
    };

    const handleCloseEditModal = () => {
        setIsEditModalOpen(false);
        setSelectedSender(null);
    };

    return (
        <div className="p-6 bg-gray-100 min-h-screen">
            <h1 className="text-2xl font-bold mb-4">Sender</h1>

            {Array.isArray(senderData) && senderData.length > 0 ? (
                <>
                    {senderData.map((sender, index) => (
                        <EmailCard
                            key={index}
                            email={sender.email}
                            name={sender.name}
                            dkim={sender.is_signed ? 'DKIM is signed' : 'Default ⚠️'}
                            dkimSigned={sender.is_signed}
                            dmarc={sender.is_signed ? "Dmarc is verified" : "Freemail domain is not recommended ⚠️"}
                            verified={sender.verified}
                            verificationText={`${sender.email} has been verified.`}
                            onEdit={() => handleEdit(sender)}
                            onDelete={() => handleDelete(sender.uuid)}
                        />
                    ))}
                </>
            ) : (
                <EmptyState
                    title="You have not created any Senders"
                    description="Create Senders To easily send mails from your domain"
                    icon={<i className="bi bi-emoji-frown text-xl"></i>}
                />
            )}

            {selectedSender && (
                <EditSenderComponent
                    isOpen={isEditModalOpen}
                    onClose={handleCloseEditModal}
                    Sender={selectedSender}
                />
            )}
        </div>
    );
}


export default SendersDashComponent;
