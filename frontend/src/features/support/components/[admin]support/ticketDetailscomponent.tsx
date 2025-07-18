import { useState, FormEvent, useMemo } from "react";
import Cookies from "js-cookie";
import * as yup from 'yup';
import { useParams } from "react-router-dom";
import { Helmet, HelmetProvider } from "react-helmet-async";
import { Tooltip } from "antd";
import useAdminSupportStore from "../../store/adminsupport.store";
import { TicketFile, Ticket } from '../../interface/support.interface';
import { useTicketDetailsQuery } from "../../hooks/useSupporTicketQuery";
import { MessageCircleQuestionIcon } from "lucide-react";

const MAX_FILES = 3

const AdminTicketDetails: React.FC = () => {
    const { replyTicket } = useAdminSupportStore()
    const { id } = useParams<{ id: string }>() as { id: string };
    //  const [isLoading, setIsLoading] = useState<boolean>(false)
    const [files, setFiles] = useState<File[]>([]);
    const [fileError, setFileError] = useState<string | null>(null);
    const [errors, setErrors] = useState<{ [key: string]: string }>({});

    const { data: supportData, isLoading } = useTicketDetailsQuery(id)

    const supportTicketData = useMemo(() => supportData?.payload, [supportData])

    const cookie = Cookies.get("Cookies");
    const user = cookie ? JSON.parse(cookie)?.details?.fullname : "";
    const email = cookie ? JSON.parse(cookie)?.details?.email : ""

    const [formData, setFormData] = useState({
        name: user,
        email: email,
        priority: 'Medium',
        subject: '',
        message: '',
    });

    const validationSchema = yup.object().shape({
        message: yup.string().required('Message is required').min(10, 'Message should be at least 10 characters'),
    });


    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: 'numeric',
            minute: 'numeric',
            second: 'numeric',
            timeZone: 'UTC'
        });
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

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        const { name, value } = e.target;
        setFormData((prevData) => ({
            ...prevData,
            [name]: value,
        }));
    };


    // const closeTkt = async () => {
    //     let confirmClose = confirm("Do you want to close this ticket?")
    //     if (confirmClose) {
    //         await closeTicket(id)
    //     }
    // }

    if (!supportTicketData) {
        return (
            <div className="flex items-center justify-center mt-20">
                <p>No ticket details available.</p>
            </div>
        );
    }

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        // Validate form data
        try {
            await validationSchema.validate(formData, { abortEarly: false });
            setErrors({});

            // Send the message and files to the backend
            await replyTicket(id, formData.message, files);

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

    const renderAttachments = (files: TicketFile[]) => {
        if (!files || files.length === 0) return null;

        const baseUrl = import.meta.env.VITE_BASE_API_URL as string;

        return (
            <div className="mt-4">
                <h4 className="text-sm font-semibold mb-2">Attachments:</h4>
                <div className="flex flex-wrap gap-2">
                    {files.map((file, index) => (
                        <a
                            key={index}
                            href={`${baseUrl}/${file.file_path}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-gray-200 text-gray-800 hover:bg-gray-300"
                        >
                            <svg className="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                                <path fillRule="evenodd" d="M8 4a3 3 0 00-3 3v4a5 5 0 0010 0V7a1 1 0 112 0v4a7 7 0 11-14 0V7a5 5 0 0110 0v4a3 3 0 11-6 0V7a1 1 0 012 0v4a1 1 0 102 0V7a3 3 0 00-3-3z" clipRule="evenodd" />
                            </svg>
                            {file.file_name}
                        </a>
                    ))}
                </div>
            </div>
        );
    };

    return (
        <>
            <HelmetProvider>

                <Helmet title={`Ticket Details for #${(supportTicketData as Ticket)?.ticket_number}`} />
                {isLoading ? (<>
                    <div className="flex items-center justify-center mt-20">
                        <span className="loading loading-spinner loading-lg"></span>
                    </div>
                </>) : (<>

                    <div className="flex flex-col lg:flex-row mb-10 mt-5 p-4 bg-gray-100 min-h-screen">
                        {/* Sidebar */}
                        <div className="w-full lg:w-1/4 bg-white rounded-lg p-4 border-r lg:sticky lg:top-5 h-auto lg:h-full">
                            <div className="mb-4">
                                <h2 className="text-lg font-bold">Ticket Information</h2>
                                <p className="text-sm">Requestor</p>
                                <p className="font-semibold">{(supportTicketData as Ticket)?.name} <span className="text-xs text-gray-500">Authorized User</span></p>
                            </div>
                            <div className="mb-4">
                                <p className="text-sm">Submitted</p>
                                <p className="font-semibold">{new Date((supportTicketData as Ticket)?.created_at as string).toLocaleString('en-US', {
                                    timeZone: 'UTC',
                                    year: 'numeric',
                                    month: 'long',
                                    day: 'numeric',
                                    hour: 'numeric',
                                    minute: 'numeric',
                                    second: 'numeric',
                                })}</p>
                            </div>
                            <div className="mb-4">
                                <p className="text-sm">Last Updated</p>
                                <p className="font-semibold">
                                    {(supportTicketData as Ticket)?.last_reply != null
                                        ? new Date((supportTicketData as Ticket)?.last_reply as string).toLocaleString('en-US', {
                                            timeZone: 'UTC',
                                            year: 'numeric',
                                            month: 'long',
                                            day: 'numeric',
                                            hour: 'numeric',
                                            minute: 'numeric',
                                            second: 'numeric',
                                        })
                                        : "Ticket has not been replied to"}
                                </p>
                            </div>
                            <div className="mb-4">
                                <p className="text-sm text-gray-600">Status</p>
                                <p className={`font-semibold text-yellow-500 flex items-center space-x-2`}>
                                    {(supportTicketData as Ticket)?.status}
                                </p>

                                <p className="text-sm text-gray-600 mb-1">Priority</p>
                                <p className="font-semibold   flex items-center space-x-2">
                                    <span className={`text-sm text-white bg-black  ${(supportTicketData as Ticket)?.priority === "high" ? "bg-red-500" : (supportTicketData as Ticket)?.priority === "medium" ? "bg-yellow-500" : (supportTicketData as Ticket)?.priority === "low" ? "bg-blue-500" : ""}  py-1 px-2 rounded-lg`}>
                                        {(supportTicketData as Ticket)?.priority}
                                    </span>
                                </p>
                            </div>
                            <div className="flex space-x-4">
                                <button className="bg-green-500 text-white py-2 px-4 rounded-lg">
                                    <a href="#replyTicket">Reply</a>
                                </button>
                                {/* <button className="bg-red-500 text-white py-2 px-4 rounded-lg" onClick={closeTkt}>Close</button> */}
                            </div>
                        </div>

                        {/* Main Content */}
                        <div className="w-full lg:w-3/4 -mt-5 p-6">

                            {/* Ticket Details */}
                            <div className="bg-white p-6 rounded-lg shadow mb-10">

                                <div className="flex justify-between mb-6">
                                    <button className="text-blue-600 mr-2 tooltip tooltip-right" data-tip="Go Back" onClick={() => window.history.back()}>
                                        <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                                        </svg>
                                    </button>
                                    <h2 className="text-xl font-semibold">View Ticket #{(supportTicketData as Ticket)?.ticket_number}</h2>
                                    <span className={`${(supportTicketData as Ticket)?.status === 'closed' ? 'bg-black text-white' :
                                        (supportTicketData as Ticket)?.status === 'pending' ? 'bg-red-500 text-white' :
                                            (supportTicketData as Ticket)?.status === 'open' ? 'bg-green-500 text-white' :
                                                (supportTicketData as Ticket)?.status === 'resolved' ? 'bg-green-500 text-white' : ''
                                        }  px-2 py-1 rounded`}>{(supportTicketData as Ticket)?.status}</span>

                                    <Tooltip title="You can reopen a ticket by replying to the ticket">
                                        <MessageCircleQuestionIcon />
                                    </Tooltip>


                                </div>

                                {(supportTicketData as Ticket)?.messages &&
                                    (supportTicketData as Ticket)?.messages.map((message) => (
                                        <div key={message.id} className="mb-6 border-t border-b py-7 pt-4">
                                            <div className="flex justify-between bg-gray-200 p-2">
                                                <p className="text-sm space-x-2 text-gray-700">
                                                    Posted by <span className="font-semibold">{message.is_admin ? message.admin.firstname + " " + message.admin.lastname : (supportTicketData as Ticket).name}</span>
                                                    on {formatDate(message.created_at)}
                                                </p>
                                                <span className="text-sm bg-blue-500 text-white rounded-md p-1">
                                                    {message.is_admin ? "admin/operator" : "authorized user"}
                                                </span>
                                            </div>
                                            <p className="mt-2">{message.message}</p>
                                            {renderAttachments(message.files)}
                                        </div>

                                    ))}
                            </div>

                            {/* Reply Section */}
                            <div className="bg-white p-6 rounded-lg shadow mt-5 mb-10" id="replyTicket">
                                <h2 className="text-2xl font-bold mb-6">Reply</h2>
                                <form onSubmit={handleSubmit}>
                                    {/* Name & Email */}
                                    <div className="grid grid-cols-2 gap-6 mb-6">
                                        <div>
                                            <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">Name</label>
                                            <input
                                                type="text"
                                                id="name"
                                                name="name"
                                                value={formData.name}
                                                onChange={handleInputChange}
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
                                                onChange={handleInputChange}
                                                className="w-full p-2 border border-gray-300 rounded-md bg-gray-100"
                                            />
                                        </div>
                                    </div>

                                    {/* Message */}
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

                                    {/* Attachments */}
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

                                    {/* Submit/Cancel Buttons */}
                                    <div className="flex items-center">
                                        <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">Submit</button>
                                        <button type="button" className="ml-4 px-4 py-2 text-gray-700 rounded-md hover:bg-gray-100">Cancel</button>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </div>

                </>)}

            </HelmetProvider>
        </>
    )
}


export default AdminTicketDetails