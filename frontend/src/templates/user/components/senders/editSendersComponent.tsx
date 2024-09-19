import React, { useState, useEffect } from 'react';
import { Modal } from "../../../../components";
import { Sender } from "../../../../store/userstore/senderStore";

type Props = {
    isOpen: boolean;
    onClose: () => void;
    Sender: Sender;
}

const EditSenderComponent: React.FC<Props> = ({ isOpen, onClose, Sender }) => {
    const [formValues, setFormValues] = useState({
        name: Sender.name,
        email: Sender.email
    });

    useEffect(() => {
        // Update local state when Sender prop changes
        setFormValues({
            name: Sender.name,
            email: Sender.email
        });
    }, [Sender]);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormValues(prev => ({ ...prev, [name]: value }));
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        // Handle form submission here
        console.log('Submitting:', formValues);
        // You would typically call an update function from your store here
        onClose();
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Edit Sender">
            <form onSubmit={handleSubmit}>
                <p className="mb-4 text-gray-600">Edit sender details</p>
                <div className="mb-4">
                    <label
                        htmlFor="name"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Sender Name
                    </label>
                    <input
                        type="text"
                        id="name"
                        name="name"
                        value={formValues.name}
                        onChange={handleInputChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                </div>
                <div className="mb-4">
                    <label
                        htmlFor="email"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Sender Email
                    </label>
                    <input
                        type="email"
                        id="email"
                        name="email"
                        value={formValues.email}
                        onChange={handleInputChange}
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
                        Save Changes
                    </button>
                </div>
            </form>
        </Modal>
    );
}

export default EditSenderComponent;