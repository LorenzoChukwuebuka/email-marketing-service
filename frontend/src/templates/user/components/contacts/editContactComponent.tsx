import { FormEvent, useEffect } from "react";
import { Modal } from "../../../../components"
import useContactStore, { Contact } from "../../../../store/userstore/contactStore";
import { ChangeEvent } from 'react';


interface EditContactProps {
    isOpen: boolean;
    onClose: () => void;
    contact: Contact | null
}
const EditContact: React.FC<EditContactProps> = ({ isOpen, onClose, contact }) => {
    const { editContactValues, setEditContactValues, editContact, getAllContacts } = useContactStore()

    const initEdit = () => {
        setEditContactValues({
            uuid: contact?.uuid as string || "",
            first_name: contact?.first_name as string || "",
            last_name: contact?.last_name as string || "",
            email: contact?.email as string || "",
            from: contact?.from as string || "",
            is_subscribed: contact?.is_subscribed
        })
    };

    const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        editContact()
        getAllContacts()
        onClose()
    }

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
        const { id, value } = e.target;
        setEditContactValues({ ...editContactValues, [id]: value });
    };

    useEffect(() => {
        if (contact) {
            initEdit()
        }

    }, [contact])

    return (
        <>
            <Modal isOpen={isOpen}
                onClose={onClose}
                title="Edit Contact">
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
                            value={editContactValues.first_name}
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
                            value={editContactValues.last_name}
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
                            value={editContactValues.email}
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
                            value={editContactValues.from}
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
                </form> </Modal>


        </>
    )
}

export default EditContact