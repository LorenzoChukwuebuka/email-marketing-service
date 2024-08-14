import { useEffect, useState } from "react";
import EmptyState from "../../../components/emptyStateComponent";
import useTemplateStore, { Template } from "../../../store/userstore/templateStore";
import { Link } from "react-router-dom";
import { Modal } from "../../../components";
import { BaseEntity } from "../../../interface/baseentity.interface";


const TransactionalTemplateDash: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { getAllTransactionalTemplates, _templateData } = useTemplateStore()
    const [previewTemplate, setPreviewTemplate] = useState<Template & BaseEntity | null>(null);

    const openPreview = (template: (Template & BaseEntity)) => {
        setPreviewTemplate(template);
        setIsModalOpen(true);
    };

    useEffect(() => {
        getAllTransactionalTemplates()
    }, [])
    return <>

        <div className="flex justify-between items-center rounded-md p-2 bg-white mt-10">
            <div className="space-x-1  h-auto w-full p-2 px-2 ">
                <button
                    className="bg-gray-300 px-2 py-2 rounded-md transition duration-300"
                    onClick={() => setIsModalOpen(true)}
                >
                    Create  Transactional Template
                </button>
            </div>

            <div className="ml-3">
                <input
                    type="text"
                    placeholder="Search..."
                    className="bg-gray-100 px-4 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-300"
                // onChange={(e) => handleSearch(e.target.value)}
                />
            </div>

        </div>

        <div className="mt-4 p-2">

            {Array.isArray(_templateData) && _templateData.length > 0 ? (
                <>
                    <div className="space-y-4">
                        {_templateData.map((template, index) => (
                            <div key={template.uuid || index} className="bg-white p-4 rounded-lg shadow-sm">
                                <div className="flex items-center space-x-4">
                                    <div className="w-12 h-12 bg-gray-300 rounded-lg"></div>
                                    <div className="flex-grow">
                                        <h3 className="text-lg font-semibold text-gray-800">{template.template_name}</h3>
                                        <p className="text-sm text-gray-600">ID - {index + 1}  {new Date(template.created_at).toLocaleString('en-US', {
                                            timeZone: 'UTC',
                                            year: 'numeric',
                                            month: 'long',
                                            day: 'numeric',
                                            hour: 'numeric',
                                            minute: 'numeric',
                                            second: 'numeric'
                                        })}</p>
                                        <div className="flex space-x-2 mt-2">
                                            <button className="text-blue-600 cursor-pointer text-sm" onClick={() => openPreview(template)}>Preview</button>
                                            <Link
                                                to={`/editor/1?type=t&uuid=${template.uuid}`}
                                                className="text-blue-600 cursor-pointer text-sm"
                                            >
                                                Edit
                                            </Link>
                                        </div>
                                    </div>
                                    <div className="flex items-center space-x-2">
                                        <span className="px-2 py-1 bg-gray-200 text-gray-700 text-xs font-medium rounded"> Draft </span>
                                        <button className="text-gray-400 hover:text-gray-600">
                                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                                <path d="M6 10a2 2 0 11-4 0 2 2 0 014 0zM12 10a2 2 0 11-4 0 2 2 0 014 0zM16 12a2 2 0 100-4 2 2 0 000 4z" />
                                            </svg>
                                        </button>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>

                </>) :
                (<>  <EmptyState title="You  have not created any Template"
                    description="Create a easily send marketing email to your audience"
                    icon={<i className="bi bi-emoji-frown text-xl"></i>}
                    buttonText="Create Template"
                // onButtonClick={() => navigate("/user/dash/marketing")}
                /> </>)}

        </div>

        <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title="Preview" >
            <>
                {previewTemplate && (
                    <div className="w-full h-full">
                        <iframe
                            srcDoc={previewTemplate.email_html}
                            title="Template Preview"
                            className="w-full h-[70vh] border-0"
                            sandbox="allow-scripts"
                        />
                    </div>
                )}

            </>
        </Modal>

    </>
}

export default TransactionalTemplateDash