import React, { FormEvent, useEffect, useRef, useState } from 'react';
import { ChevronUpIcon, ChevronDownIcon } from 'lucide-react';
import useSupportStore from '../../../../store/userstore/support.store';
import Cookies from "js-cookie";
import { useNavigate } from 'react-router-dom';
import * as Yup from "yup"

const MAX_FILES = 3;

const validationSchema = Yup.object().shape({
    name: Yup.string().required('Name is required'),
    email: Yup.string().email('Invalid email').required('Email is required'),
    priority: Yup.string().oneOf(['Low', 'Medium', 'High'], 'Invalid priority').required('Priority is required'),
    subject: Yup.string().required('Subject is required'),
    message: Yup.string().required('Message is required'),
});

const SupportRequestForm = () => {
    const [isRecentTicketsOpen, setIsRecentTicketsOpen] = useState(false);
    const { getTickets, supportTicketData, createTicket } = useSupportStore()
    const [files, setFiles] = useState<File[]>([]);
    const [fileError, setFileError] = useState<string | null>(null);
    const [errors, setErrors] = useState<{ [key: string]: string }>({});
    const parentRef = useRef<HTMLDivElement>(null);
    const [parentHeight, setParentHeight] = useState<number | null>(null);
    const navigate = useNavigate()

    let cookie = Cookies.get("Cookies");
    let user = cookie ? JSON.parse(cookie)?.details?.fullname : "";
    let email = cookie ? JSON.parse(cookie)?.details?.email : ""

    const [formData, setFormData] = useState({
        name: user,
        email: email,
        priority: 'Medium',
        subject: '',
        message: '',
    });

    useEffect(() => {
        const fetchSupportData = async () => {
            await getTickets()
        }

        fetchSupportData()
    }, [getTickets])

    useEffect(() => {
        if (parentRef.current) {
            setParentHeight(parentRef.current.clientHeight);
        }
    }, [isRecentTicketsOpen]);



    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        const { name, value } = e.target;
        setFormData((prevData) => ({
            ...prevData,
            [name]: value,
        }));

        // Clear the error for this field when the user starts typing
        setErrors(prevErrors => ({ ...prevErrors, [name]: '' }));
    };


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
            await createTicket({ ...formData, priority: formData.priority.toLowerCase() }, files);
            await new Promise(resolve => setTimeout(resolve, 700));
            await getTickets();
            // Reset form or navigate to a success page
        } catch (err) {
            if (err instanceof Yup.ValidationError) {
                const newErrors: { [key: string]: string } = {};
                err.inner.forEach((error) => {
                    if (error.path) {
                        newErrors[error.path] = error.message;
                    }
                });
                setErrors(newErrors);
            }
        }
    };

    const toggleRecentTickets = () => {
        setIsRecentTicketsOpen(!isRecentTicketsOpen);
    };

    const handleNavigation = (uuid: string) => {
        navigate("/user/dash/support/ticket/details/" + uuid)
    }

    return (
        <div ref={parentRef} className="flex max-w-6xl mx-auto space-x-1   mb-10 p-6 mt-10 rounded-lg">

            {/* Sidebar */}
            <div className={`w-64 mr-6 bg-white border-r p-3 transition-all duration-300 ease-in-out`}
                style={{
                    height: isRecentTicketsOpen ? `${parentHeight}px` : '10em',
                    overflow: 'hidden'
                }}>
                <button className="text-blue-600 mr-2" onClick={() => window.history.back()}>
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                    </svg>
                </button>
                <div className="mb-4">
                    <button
                        onClick={toggleRecentTickets}
                        className="w-full flex justify-between items-center p-3 bg-blue-100 text-blue-700 rounded-md"
                    >
                        <span className="font-medium">Your Recent Tickets</span>
                        {isRecentTicketsOpen ? (
                            <ChevronUpIcon className="w-5 h-5" />
                        ) : (
                            <ChevronDownIcon className="w-5 h-5" />
                        )}
                    </button>
                    {isRecentTicketsOpen && (
                        <div className="mt-2 p-3 bg-white border overflow-y-auto border-gray-200 rounded-md shadow-sm">
                            {Array.isArray(supportTicketData) && supportTicketData.length > 0 ? (
                                supportTicketData.map((ticket, index) => (
                                    <div key={ticket.uuid} className="mt-2 p-3 bg-white border border-gray-200 rounded-md cursor-pointer shadow-sm" onClick={() => handleNavigation(ticket.uuid)}>
                                        <div className="mb-2">
                                            <div className="font-medium">#{index + 1 + ticket.ticket_number} - {ticket.subject}</div>
                                            <div className="text-sm text-gray-500">
                                                {ticket.status.charAt(0).toUpperCase() + ticket.status.slice(1)} {ticket.last_reply ? ` - Last reply: ${new Date(ticket.last_reply).toLocaleDateString()}` : ''}
                                            </div>
                                        </div>
                                    </div>
                                ))
                            ) : (
                                <div className="mb-2">
                                    <div className="font-medium">You have not created any tickets</div>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            </div>

            {/* Form */}
            <div className="flex-1 bg-white p-4 rounded-sm" >
                <h2 className="text-2xl font-bold mb-6">Create new Support Request</h2>
                <p className="p-2 text-md font-normal"> If you can't find a solution to your problems in our knowledgebase, you can submit a ticket. An agent will respond to you soonest </p>
                <form onSubmit={handleSubmit}>
                    <div className="grid grid-cols-2 gap-6 mb-6">
                        <div>
                            <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
                                Name
                            </label>
                            <input
                                type="text"
                                id="name"
                                name="name"
                                value={formData.name}
                                onChange={handleInputChange}
                                className={`w-full p-2 border ${errors.name ? 'border-red-500' : 'border-gray-300'} rounded-md bg-gray-100`}
                            />
                            {errors.name && <p className="mt-1 text-xs text-red-500">{errors.name}</p>}
                        </div>
                        <div>
                            <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                                Email Address
                            </label>
                            <input
                                type="email"
                                id="email"
                                name="email"
                                value={formData.email}
                                onChange={handleInputChange}
                                className={`w-full p-2 border ${errors.email ? 'border-red-500' : 'border-gray-300'} rounded-md bg-gray-100`}
                            />
                            {errors.email && <p className="mt-1 text-xs text-red-500">{errors.email}</p>}
                        </div>
                    </div>
                    <div className="mb-6">
                        <label htmlFor="priority" className="block text-sm font-medium text-gray-700 mb-1">
                            Priority
                        </label>
                        <select
                            id="priority"
                            name="priority"
                            value={formData.priority}
                            onChange={handleInputChange}
                            className={`w-full p-2 border ${errors.priority ? 'border-red-500' : 'border-gray-300'} rounded-md`}
                        >
                            <option value="Low">Low</option>
                            <option value="Medium">Medium</option>
                            <option value="High">High</option>
                        </select>
                        {errors.priority && <p className="mt-1 text-xs text-red-500">{errors.priority}</p>}
                    </div>
                    <div className="mb-6">
                        <label htmlFor="subject" className="block text-sm font-medium text-gray-700 mb-1">
                            Subject
                        </label>
                        <input
                            type="text"
                            id="subject"
                            name="subject"
                            value={formData.subject}
                            onChange={handleInputChange}
                            className={`w-full p-2 border ${errors.subject ? 'border-red-500' : 'border-gray-300'} rounded-md`}
                        />
                        {errors.subject && <p className="mt-1 text-xs text-red-500">{errors.subject}</p>}
                    </div>
                    <div className="mb-6">
                        <label htmlFor="message" className="block text-sm font-medium text-gray-700 mb-1">
                            Message
                        </label>
                        <div className={`border ${errors.message ? 'border-red-500' : 'border-gray-300'} rounded-md`}>
                            <textarea
                                id="message"
                                name="message"
                                value={formData.message}
                                onChange={handleInputChange}
                                rows={6}
                                className="w-full p-2 border-none focus:ring-0"
                            ></textarea>
                        </div>
                        {errors.message && <p className="mt-1 text-xs text-red-500">{errors.message}</p>}
                    </div>
                    <div className="mb-6">
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Attachments (Max {MAX_FILES} files)
                        </label>
                        <div className="flex items-center">
                            <input
                                type="file"
                                onChange={handleFileChange}
                                multiple
                                className="flex-grow p-2 border border-gray-300 rounded-md"
                                disabled={files.length >= MAX_FILES}
                            />
                        </div>
                        {fileError && (
                            <small className="mt-2 text-red-500">
                                {fileError}
                            </small>
                        )}
                        <div className="mt-2">
                            {files.map((file, index) => (
                                <div key={index} className="flex items-center justify-between text-sm text-gray-600 mb-1">
                                    <span>{file.name}</span>
                                    <button
                                        type="button"
                                        onClick={() => removeFile(index)}
                                        className="text-red-500 hover:text-red-700"
                                    >
                                        Remove
                                    </button>
                                </div>
                            ))}
                        </div>
                        <p className="text-xs text-gray-500 mt-1">
                            Allowed File Extensions: .jpg, .gif, .jpeg, .png, .txt, .pdf (Max file size: 1024MB)
                        </p>
                    </div>
                    <div className="flex items-center">
                        <button
                            type="submit"
                            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                        >
                            Submit
                        </button>
                        <button
                            type="button"
                            className="ml-4 px-4 py-2 text-gray-700 rounded-md hover:bg-gray-100"
                        >
                            Cancel
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default SupportRequestForm;