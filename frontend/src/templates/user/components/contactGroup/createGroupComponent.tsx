import { Modal } from "../../../../components"
import * as Yup from "yup";
import { ChangeEvent, FormEvent, useState } from "react";
import useContactGroupStore from "../../../../store/userstore/contactGroupStore";

interface CreateGroupProps {
    isOpen: boolean;
    onClose: () => void;
}

const CreateGroup: React.FC<CreateGroupProps> = ({ isOpen, onClose }) => {
    const [errors, setErrors] = useState<{ [key: string]: string }>({});

    const { formValues, setFormValues, isLoading, createGroup, getAllGroups } = useContactGroupStore()

    const validationSchema = Yup.object().shape({
        group_name: Yup.string()
            .required("group name is required"),
    });

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()

        try {
            await validationSchema.validate(formValues, { abortEarly: false });
            await createGroup()
            await getAllGroups()
            onClose();
            setErrors({})
        } catch (err) {
            const validationErrors: { [key: string]: string } = {};
            if (err instanceof Yup.ValidationError) {
                err.inner.forEach((error) => {
                    validationErrors[error.path || ""] = error.message;
                });
                setErrors(validationErrors);
            }
        }
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { id, value } = e.target;
        setFormValues({ ...formValues, [id]: value });
    };

    return <>
        <Modal isOpen={isOpen} onClose={onClose} title="Create Group">

            <>
                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label
                            htmlFor="first_name"
                            className="block text-sm font-medium text-gray-700"
                        >
                            Group name
                        </label>
                        <input
                            type="text"
                            id="group_name"
                            value={formValues.group_name}
                            onChange={handleChange}
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            required
                        />
                        {errors.group_name && (
                            <div style={{ color: "red" }}>{errors.group_name}</div>
                        )}
                    </div>

                    <div className="mb-4">
                        <label
                            htmlFor="last_name"
                            className="block text-sm font-medium text-gray-700"
                        >
                            Description
                        </label>
                        <textarea
                            id="description"
                            value={formValues.description}
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
                            {isLoading ? "Please wait ..." : "Submit"}
                        </button>
                    </div>
                </form>

            </>

        </Modal>


    </>
}


export default CreateGroup