// TicketReplyForm.tsx
import { useState, FormEvent } from "react";
import * as yup from 'yup';

interface TicketReplyFormProps {
    user: string;
    email: string;
    onSubmit: (message: string, files: File[]) => Promise<void>;
}

const MAX_FILES = 3;

const TicketReplyForm: React.FC<TicketReplyFormProps> = ({ user, email, onSubmit }) => {
    const [files, setFiles] = useState<File[]>([]);
    const [fileError, setFileError] = useState<string | null>(null);
    const [errors, setErrors] = useState<{ [key: string]: string }>({});
    const [formData, setFormData] = useState({
        name: user,
        email: email,
        message: '',
    });

    const validationSchema = yup.object().shape({
        message: yup.string().required('Message is required').min(10, 'Message should be at least 10 characters'),
    });

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files) {
            const newFiles = Array.from(e.target.files);
            if (files.length + newFiles.length > MAX_FILES) {
                setFileError(`You can only upload a maximum of ${MAX_FILES} files.`);
                return;
            }
            setFiles(prevFiles => [...prevFiles, ...newFiles].slice(0, MAX_FILES));
            setFileError(null);
        }
    };

    const removeFile = (index: number) => {
        setFiles(prevFiles => prevFiles.filter((_, i) => i !== index));
        setFileError(null);
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        try {
            await validationSchema.validate(formData, { abortEarly: false });
            setErrors({});
            await onSubmit(formData.message, files);
        } catch (err) {
            if (err instanceof yup.ValidationError) {
                const newErrors: { [key: string]: string } = {};
                err.inner.forEach((validationError) => {
                    if (validationError.path) {
                        newErrors[validationError.path] = validationError.message;
                    }
                });
                setErrors(newErrors);
            }
        }
    };

    return (
        <div className="bg-white p-6 rounded-lg shadow mt-5 mb-10" id="replyTicket">
            <h2 className="text-2xl font-bold mb-6">Reply</h2>
            <form onSubmit={handleSubmit}>
                <div className="grid grid-cols-2 gap-6 mb-6">
                    <div>
                        <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">Name</label>
                        <input
                            type="text"
                            id="name"
                            name="name"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            className="w-full p-2 border border-gray-300 rounded-md bg-gray-100"
                        />
                    </div>
                    <div>
                        <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">Email Address</label>
                        <input
                            type="email"
                            id="email"
                            name="email"
                            value={formData.email}
                            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                            className="w-full p-2 border border-gray-300 rounded-md bg-gray-100"
                        />
                    </div>
                </div>

                <div className="mb-6">
                    <label htmlFor="message" className="block text-sm font-medium text-gray-700 mb-1">Message</label>
                    <div className="border border-gray-300 rounded-md">
                        <textarea
                            id="message"
                            name="message"
                            value={formData.message}
                            onChange={(e) => setFormData({ ...formData, message: e.target.value })}
                            rows={6}
                            className="w-full p-2 border-none focus:ring-0"
                        ></textarea>
                    </div>
                    {errors.message && <p className="text-red-500 text-sm mt-1">{errors.message}</p>}
                </div>

                <div className="mb-6">
                    <label className="block text-sm font-medium text-gray-700 mb-1">Attachments (Max {MAX_FILES} files)</label>
                    <div className="flex items-center">
                        <input
                            type="file"
                            onChange={handleFileChange}
                            multiple
                            className="flex-grow p-2 border border-gray-300 rounded-md"
                            disabled={files.length >= MAX_FILES}
                        />
                    </div>
                    {fileError && <small className="mt-2 text-red-500">{fileError}</small>}
                    <div className="mt-2">
                        {files.map((file, index) => (
                            <div key={index} className="flex items-center justify-between text-sm text-gray-600 mb-1">
                                <span>{file.name}</span>
                                <button type="button" onClick={() => removeFile(index)} className="text-red-500 hover:text-red-700">Remove</button>
                            </div>
                        ))}
                    </div>
                    <p className="text-xs text-gray-500 mt-1">Allowed File Extensions: .jpg, .gif, .jpeg, .png, .txt, .pdf (Max file size: 1024MB)</p>
                </div>

                <div className="flex items-center">
                    <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">Submit</button>
                    <button type="button" className="ml-4 px-4 py-2 text-gray-700 rounded-md hover:bg-gray-100">Cancel</button>
                </div>
            </form>
        </div>
    );
};

export default TicketReplyForm