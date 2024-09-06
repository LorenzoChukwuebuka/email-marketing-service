type Props = {
    email: string;
    name: string;
    dkim: string;
    dmarc: string;
    onEdit: () => void;
    onDelete: () => void
}

const EmailCard = ({ email, name, dkim, dmarc, onEdit, onDelete }: Props) => {
    return (
        <div className="p-4 bg-white shadow rounded-lg mb-4">
            <div className="flex items-start">
                <div className="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center text-green-500 mr-4">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className="w-6 h-6">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16.707 9.293a1 1 0 00-1.414 0L12 12.586 8.707 9.293a1 1 0 10-1.414 1.414l4 4a1 1 0 001.414 0l4-4a1 1 0 000-1.414z" />
                    </svg>
                </div>
                <div className="flex-1">
                    <h4 className="font-semibold text-gray-800">{name} <span className="text-gray-600">({email})</span></h4>
                    <p className="text-sm text-gray-600">Verified • <span className="text-blue-600">{email}</span> has verified it via email.</p>
                    <div className="flex items-center text-sm mt-2">
                        <span className="mr-6">
                            <span className="font-medium">IP address:</span> <span className="text-gray-700">Shared IP</span>
                        </span>
                        <span className="mr-6">
                            <span className="font-medium">DKIM signature:</span> <span className="text-yellow-500">{dkim}</span>
                        </span>
                        <span>
                            <span className="font-medium">DMARC:</span> <span className="text-yellow-500">{dmarc}</span>
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
    const handleEdit = (email: string) => {
        console.log(`Edit ${email}`);
        // Logic to handle editing the email
        // You might want to open a modal or redirect to an edit page
    };

    const handleDelete = (email: string) => {
        console.log(`Delete ${email}`);
        // Logic to handle deleting the email
        // Confirm with the user and then proceed to delete
    };

    return <>
        <div className="p-6 bg-gray-100 min-h-screen">
            <EmailCard
                email="enzobyte.tech@gmail.com"
                name="hello"
                dkim="Default ⚠️"
                dmarc="Freemail domain is not recommended ⚠️"
                onEdit={() => handleEdit("enzobyte.tech@gmail.com")}
                onDelete={() => handleDelete("enzobyte.tech@gmail.com")}
            />
            <EmailCard
                email="kampus360ng@gmail.com"
                name="My Company"
                dkim="Default ⚠️"
                dmarc="Freemail domain is not recommended ⚠️"
                onEdit={() => handleEdit("kampus360ng@gmail.com")}
                onDelete={() => handleDelete("kampus360ng@gmail.com")}
            />
        </div>
    </>
}

export default SendersDashComponent