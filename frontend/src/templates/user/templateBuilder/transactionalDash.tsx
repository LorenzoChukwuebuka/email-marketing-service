import { useEffect, useRef, useState } from "react";
import EmptyState from "../../../components/emptyStateComponent";
import useTemplateStore, { Template } from "../../../store/userstore/templateStore";
import { Link, useNavigate } from "react-router-dom";
import { Modal } from "../../../components";
import { BaseEntity } from "../../../interface/baseentity.interface";

const TransactionalTemplateDash: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<number | boolean | null>(null);
    const { getAllTransactionalTemplates, _templateData, deleteTemplate } = useTemplateStore()
    const [previewTemplate, setPreviewTemplate] = useState<Template & BaseEntity | null>(null);
    const [searchTerm, setSearchTerm] = useState<string>("");
    const modalRef = useRef<HTMLDivElement>(null)


    const navigate = useNavigate()


    useEffect(() => {
        getAllTransactionalTemplates()
    }, [])


    const openPreview = (template: (Template & BaseEntity)) => {
        setPreviewTemplate(template);
        setIsModalOpen(true);
    };

    const filteredTemplates = (_templateData as (Template & BaseEntity)[]).filter((template) =>
        template.template_name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
        setSearchTerm(e.target.value);
    };

    const handleNavigate = (template: (Template & BaseEntity)) => {
        const editorType = template.editor_type

        let redirectUrl = "";
        switch (editorType) {
            case "html-editor":
                redirectUrl = `/editor/2?type=t&uuid=${template.uuid}`;
                break;
            case "drag-and-drop":
                redirectUrl = `/editor/1?type=t&uuid=${template.uuid}`;
                break;
            case "rich-text":
                redirectUrl = `/editor/3?type=t&uuid=${template.uuid}`;
                break;
            default:
                console.log("Unknown editor type:", editorType);
                return;
        }

        window.location.href = redirectUrl;
    }

    const deleteTempl = async (template: (Template & BaseEntity)) => {
        await deleteTemplate(template.uuid)
        await getAllTransactionalTemplates()
    }

    return <>

        <div className="flex justify-between items-center rounded-md p-2 bg-white mt-10">
            <div className="space-x-1  h-auto w-full p-2 px-2 ">
                <button className="bg-gray-300 px-2 py-2 rounded-md transition duration-300">
                    <Link to="/user/dash/transactional"> Create  Transactional Template </Link>
                </button>
            </div>

            <div className="ml-3">
                <input
                    type="text"
                    placeholder="Search..."
                    className="bg-gray-100 px-4 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-300"
                    onChange={handleSearch}
                />
            </div>

        </div>

        <div className="mt-4 p-2">

            {Array.isArray(filteredTemplates) && filteredTemplates.length > 0 ? (
                <>
                    <div className="space-y-4">
                        {filteredTemplates.map((template, index) => (
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
                                        })}
                                        </p>
                                        <div className="flex space-x-2 mt-2">
                                            <button className="text-blue-600 cursor-pointer text-sm" onClick={() => openPreview(template)}>Preview</button>
                                            <button onClick={() => handleNavigate(template)} className="text-blue-600 cursor-pointer text-sm"
                                            >
                                                Edit
                                            </button>
                                        </div>
                                    </div>
                                    <div className="flex items-center space-x-2">
                                        <span className="px-2 py-1 bg-gray-200 text-gray-700 text-xs font-medium rounded"> Draft </span>
                                        <button className="text-gray-400 hover:text-gray-600" onClick={() => setIsModalOpen(isModalOpen === index ? null : index)}>
                                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                                <path d="M6 10a2 2 0 11-4 0 2 2 0 014 0zM12 10a2 2 0 11-4 0 2 2 0 014 0zM16 12a2 2 0 100-4 2 2 0 000 4z" />
                                            </svg>
                                        </button>
                                        {isModalOpen === index && (
                                            <div
                                                ref={modalRef}
                                                className="absolute right-[2em] mt-[10em] w-28 bg-white  border border-gray-300 rounded-md shadow-lg z-10"
                                            >
                                                <button className="block w-full px-4 py-2 text-  text-sm text-red-700 hover:bg-gray-100" onClick={() => deleteTempl(template)}>
                                                    Delete
                                                </button>

                                            </div>
                                        )}
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
                    onButtonClick={() => navigate("/user/dash/transactional")}
                /> </>)}

        </div>

        <Modal isOpen={isModalOpen as boolean} onClose={() => setIsModalOpen(false)} title="Preview Template" >
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