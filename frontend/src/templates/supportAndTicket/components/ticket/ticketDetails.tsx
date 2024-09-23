import { useEffect, useState, FormEvent } from "react";
import useSupportStore, { Ticket } from "../../../../store/userstore/support.store";
import { useParams } from "react-router-dom";
import { parseDate } from '../../../../utils/utils';
import Cookies from "js-cookie";

const MAX_FILES = 3

const TicketDetails: React.FC = () => {

    const { getTicketDetails, supportTicketData } = useSupportStore()
    const { id } = useParams<{ id: string }>() as { id: string };
    const [isLoading, setIsLoading] = useState<boolean>(false)
    const [files, setFiles] = useState<File[]>([]);
    const [fileError, setFileError] = useState<string | null>(null);

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


    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

    };


    if (!supportTicketData) {
        return (
            <div className="flex items-center justify-center mt-20">
                <p>No ticket details available.</p>
            </div>
        );
    }

    useEffect(() => {
        const fetchD = async () => {
            setIsLoading(true);
            await getTicketDetails(id);
            await new Promise(resolve => setTimeout(resolve, 1000))
            setIsLoading(false);
        };

        fetchD();
    }, [getTicketDetails, id]);

    return (
        <>
            {isLoading ? (<>
                <div className="flex items-center justify-center mt-20">
                    <span className="loading loading-spinner loading-lg"></span>
                </div>
            </>) : (<>


                <div className="flex flex-col lg:flex-row h-screen mt-5 p-4 bg-gray-100">
                    {/* Sidebar */}
                    <div className="w-full lg:w-1/4 bg-white h-[25em] rounded-lg p-4 border-r">
                        <div className="mb-4">
                            <h2 className="text-lg font-bold">Ticket Information</h2>
                            <p className="text-sm">Requestor</p>
                            <p className="font-semibold">{(supportTicketData as Ticket)?.name} <span className="text-xs text-gray-500">Authorized User</span></p>
                        </div>
                        {/* <div className="mb-4">
                    <p className="text-sm">Department</p>
                    <p className="font-semibold">Technical Support</p>
                </div> */}
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


                            <p className="font-semibold"> {(supportTicketData as Ticket)?.last_reply != null ? new Date((supportTicketData as Ticket)?.last_reply as string).toLocaleString('en-US', {
                                timeZone: 'UTC',
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric',
                                hour: 'numeric',
                                minute: 'numeric',
                                second: 'numeric',
                            }) : "ticket has not been replied to"}</p>
                        </div>
                        <div className="mb-4">
                            <p className="text-sm text-gray-600 mb-2">Status/Priority</p>
                            <p className="font-semibold text-yellow-500 flex items-center space-x-2">
                                <span>{(supportTicketData as Ticket)?.status}</span>
                                <span className="text-sm text-white bg-black py-1 px-2 rounded-lg">
                                    {(supportTicketData as Ticket)?.priority}
                                </span>
                            </p>
                        </div>

                        <button className="bg-green-500 text-white py-2 px-4 rounded-lg mr-2"> <a href="#replyTicket"> Reply </a>  </button>
                        <button className="bg-red-500 text-white py-2 px-4 rounded-lg">Closed</button>

                        {/* CC Recipients */}
                        {/* <div className="mt-6">
                            <h2 className="text-lg font-bold">CC Recipients</h2>
                            <div className="mt-2 flex">
                                <input
                                    type="email"
                                    className="border p-2 w-full rounded-lg"
                                    placeholder="Enter Email Address"
                                />
                                <button className="bg-blue-500 text-white py-2 px-4 rounded-lg ml-2">Add</button>
                            </div>
                        </div> */}
                    </div>

                    {/* Main Ticket Content */}
                    <div className="w-full lg:w-3/4 p-6 mb-5">

                        <div className="bg-white p-6 rounded-lg shadow mb-5">
                            <div className="flex  justify-between mb-6">
                                <button className="text-blue-600 mr-2" onClick={() => window.history.back()}>
                                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                                    </svg>
                                </button>
                                <h2 className="text-xl font-semibold">View Ticket #{(supportTicketData as Ticket)?.ticket_number}</h2>
                                <span className="bg-red-100 text-red-600 px-2 py-1 rounded">{(supportTicketData as Ticket)?.status}</span>
                            </div>

                            {(supportTicketData as Ticket)?.messages && (supportTicketData as Ticket)?.messages.map((message, index) => (
                                <div key={message.uuid} className="mb-6 border-t border-b py-7 pt-4">
                                    <p className="text-sm text-gray-700 bg-gray-200 p-2 space-x-4" >
                                        Posted by <span className="font-semibold">{message.is_admin ? "Admin" : (supportTicketData as Ticket).name}</span> on {formatDate(message.created_at)}
                                    </p>
                                    <p className="mt-2">{message.message}</p>
                                </div>
                            ))}

                            {/* <div className="flex">
                                <p className="text-sm font-semibold mr-2">Rating:</p>
                                <div className="flex">
                                    {[...Array(5)].map((star, index) => (
                                        <svg key={index} className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path strokeLinecap="round" strokeLinejoin="round" d="M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"></path></svg>
                                    ))}
                                </div>
                            </div> */}
                        </div>



                        <div className="flex-1 p-4 rounded-lg bg-white shadow-sm mb-5" id="replyTicket">
                            <h2 className="text-2xl font-bold mb-6"> Reply </h2>
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
                                            className="w-full p-2 border border-gray-300 rounded-md bg-gray-100"
                                        />
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
                                            className="w-full p-2 border border-gray-300 rounded-md bg-gray-100"
                                        />
                                    </div>
                                </div>

                                <div className="mb-6">
                                    <label htmlFor="message" className="block text-sm font-medium text-gray-700 mb-1">
                                        Message
                                    </label>
                                    <div className="border border-gray-300 rounded-md">
                                        <textarea
                                            id="message"
                                            name="message"
                                            value={formData.message}
                                            onChange={handleInputChange}
                                            rows={6}
                                            className="w-full p-2 border-none focus:ring-0"
                                        ></textarea>
                                    </div>
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


                </div>





            </>)}



        </>
    );
};

export default TicketDetails;
