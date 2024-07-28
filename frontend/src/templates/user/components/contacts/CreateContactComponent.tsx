
import { ChangeEvent, FormEvent } from "react";
import useContactStore from "../../../../store/userstore/contactStore";
import { Modal } from "../../../../components";

interface CreateContactProps {
    isOpen: boolean;
    onClose: () => void;
}


const CreateContact: React.FC<CreateContactProps> = ({ isOpen, onClose }) => {
    const { setContactFormValues, contactFormValues, createContact, getAllContacts } = useContactStore();

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        await createContact();
        await getAllContacts();
        onClose();
    };

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
        const { id, value } = e.target;
        setContactFormValues({ ...contactFormValues, [id]: value });
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title="Create Contact"
        >
            <form onSubmit={handleSubmit}>
                <div className="mb-4">
                    <label
                        htmlFor="first_name"
                        className="block text-sm font-medium text-gray-700"
                    >
                        First Name
                    </label>
                    <input
                        type="text"
                        id="first_name"
                        value={contactFormValues.first_name}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                </div>

                <div className="mb-4">
                    <label
                        htmlFor="last_name"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Last Name
                    </label>
                    <input
                        type="text"
                        id="last_name"
                        value={contactFormValues.last_name}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                </div>

                <div className="mb-4">
                    <label
                        htmlFor="email"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Email
                    </label>
                    <input
                        type="email"
                        id="email"
                        value={contactFormValues.email}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                </div>

                <div className="mb-4">
                    <label
                        htmlFor="from"
                        className="block text-sm font-medium text-gray-700"
                    >
                        From
                    </label>
                    <input
                        type="text"
                        id="from"
                        value={contactFormValues.from}
                        onChange={handleChange}
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
                    >
                        Submit
                    </button>
                </div>
            </form>
        </Modal>
    );
};

export default CreateContact;