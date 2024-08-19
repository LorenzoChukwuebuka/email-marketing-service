import { FormEvent, useEffect } from "react";
import { Modal } from "../../../../components";
import useContactGroupStore, { ContactGroupData } from "../../../../store/userstore/contactGroupStore";

interface EditContactProps {
    isOpen: boolean;
    onClose: () => void;
    group: ContactGroupData | null
}
const EditGroupComponent: React.FC<EditContactProps> = ({ isOpen, onClose, group }) => {
    const { setEditValues, editValues, isLoading, updateGroup } = useContactGroupStore()
    const iniEdit = () => {
        setEditValues({
            uuid: group?.uuid as string,
            group_name: group?.group_name as string,
            description: group?.description as string
        })
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { id, value } = e.target;
        setEditValues({ ...editValues, [id]: value });
    };


    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()
        await updateGroup()
        onClose()
    }

    useEffect(() => {
        if (group) {
            iniEdit()
        }
    }, [group])
    return <>
        <Modal isOpen={isOpen} onClose={onClose} title="Edit Group"> <>

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
                        value={editValues.group_name}
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
                        Description
                    </label>
                    <textarea
                        id="description"
                        value={editValues.description}
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


        </> </Modal>
    </>
}

export default EditGroupComponent