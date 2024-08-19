import { FormEvent, useState } from "react";
import { Modal } from "../../../../components"
import useTemplateStore from "../../../../store/userstore/templateStore";
import * as Yup from "yup";

interface Props {
    isOpen: boolean;
    onClose: () => void;
    editorType: "drag-and-drop" | "html-editor" | "rich-text";
}

const CreateTransactionalTemplate: React.FC<Props> = ({ isOpen, onClose, editorType }) => {

    const { setFormValues, formValues, createTemplate } = useTemplateStore()
    const [errors, setErrors] = useState<{ [key: string]: string }>({});

    const validationSchema = Yup.object().shape({
        template_name: Yup.string().required("template name is required"),
        tags: Yup.string()
            .test("valid-tags", "Tags must be non-empty, separated by commas, and contain no spaces", (value) => {
                if (!value) return false;
                return value.split(',').every(tag => tag.trim() !== "" && !tag.includes(' '));
            })
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { id, value } = e.target;
        setFormValues({ ...formValues, [id]: value });
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()
        try {
            await validationSchema.validate(formValues, { abortEarly: false });
            setFormValues({
                ...formValues,
                type: "transactional",
                editor_type: editorType,
                is_editable: true,
                is_gallery_template: false,
                is_published: false,
            });
            await createTemplate()
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
    return <>

        <Modal isOpen={isOpen} onClose={onClose} title="Create  Template">
            <>
                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label
                            htmlFor="first_name"
                            className="block text-sm font-medium text-gray-700"
                        >
                            Enter Template Name
                        </label>
                        <input
                            type="text"
                            id="template_name"
                            placeholder="template ..."
                            value={formValues.template_name}
                            onChange={handleChange}
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            required
                        />

                        {errors.templateName && (
                            <div style={{ color: "red" }}>{errors.templateName}</div>
                        )}

                    </div>

                    <div className="mb-4">
                        <label
                            htmlFor="last_name"
                            className="block text-sm font-medium text-gray-700"
                        >
                            Enter Tag
                        </label>
                        <input
                            type="text"
                            id="tags"
                            value={formValues.tags}
                            onChange={handleChange}
                            placeholder="E.g. signup, welcome"
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            required
                        />

                        {errors.tags && (
                            <div style={{ color: "red" }}>{errors.tags}</div>
                        )}

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
            </>
        </Modal>

    </>
}


export default CreateTransactionalTemplate