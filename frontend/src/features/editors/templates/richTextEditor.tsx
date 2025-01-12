import { useRef, useState, useEffect } from 'react';
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import { useLocation, useNavigate } from 'react-router-dom';
import 'react-quill/dist/quill.snow.css';
import useCampaignStore from '../../campaign/store/campaign.store';
import useTemplateStore from '../../email-templates/store/template.store';
import { useSingleMarketingTemplateQuery } from '../../email-templates/hooks/useMarketingTemplateQuery';
import { useSingleTransactionalTemplateQuery } from '../../email-templates/hooks/useTransactionTemplateQuery';
import SendTestEmail from '../../email-templates/components/sendTestEmail';

const RichTextEditor = () => {
    const quillRef = useRef<ReactQuill>(null);
    const [autoSaved, setAutoSaved] = useState(false);
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { currentTemplate, updateTemplate, setCurrentTemplate } = useTemplateStore();
    const navigate = useNavigate();
    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const uuid = queryParams.get('uuid');
    const _type = queryParams.get('type');
    const [editorContent, setEditorContent] = useState(currentTemplate?.email_html || '');
    const { updateCampaign, setCreateCampaignValues, currentCampaignId, clearCurrentCampaignId } = useCampaignStore()

    const transactionalQuery = useSingleTransactionalTemplateQuery(uuid as string);
    const marketingQuery = useSingleMarketingTemplateQuery(uuid as string);

    useEffect(() => {
        if (!uuid) return;

        if (_type === "t") {
            if (transactionalQuery.data) {
                setCurrentTemplate(transactionalQuery.data);
            }
        } else {
            if (marketingQuery.data) {
                setCurrentTemplate(marketingQuery.data);
            }
        }
    }, [uuid, _type, transactionalQuery.data, marketingQuery.data, setCurrentTemplate])

    useEffect(() => {
        return () => {
            setCurrentTemplate(null);
        };
    }, [setCurrentTemplate]);

    useEffect(() => {
        const debounce = setTimeout(() => {
            saveDesign();
        }, 1000);
        return () => clearTimeout(debounce);
    }, [editorContent]);


    const saveDesign = async () => {
        if (currentCampaignId) {
            setCreateCampaignValues({ template_id: uuid as string })
            new Promise(resolve => setTimeout(resolve, 3000));
            updateCampaign(currentCampaignId)
        }

        if (uuid && currentTemplate) {
            const updatedTemplate = {
                ...currentTemplate,
                email_html: editorContent
            };
            await updateTemplate(uuid, updatedTemplate)
            setAutoSaved(true);
            console.log("Design saved to database!");
            setTimeout(() => setAutoSaved(false), 3000);
        } else {
            console.log("UUID or template is missing", { uuid, currentTemplate });
        }
    };

    const handleChange = (content: string) => {
        setEditorContent(content);
    };

    const testDesign = () => {
        setIsModalOpen(true)
    }

    if (!currentTemplate) {
        return <div>Loading template...</div>;
    }

    const handleNavigate = () => {
        if (currentCampaignId) {
            clearCurrentCampaignId();
            navigate("/app/campaign/edit/" + currentCampaignId);
        } else {
            if (_type === "t") {
                navigate("/app/templates/transactional");
            } else {
                navigate("/app/templates/marketing");
            }
        }

    };

    return (
        <div className="h-screen flex flex-col p-4">
            <header className="flex items-center justify-between bg-gray-100 px-4 h-[5em] py-2">
                <div className="flex items-center">
                    <button className="mr-2 text-gray-600" onClick={() => handleNavigate()}>
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path fillRule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clipRule="evenodd" />
                        </svg>
                    </button>
                    <div>
                        <h1 className="text-sm font-semibold">{currentTemplate?.template_name}</h1>
                    </div>
                </div>
                <div className="flex items-center space-x-2 text-xs">
                    {autoSaved && (
                        <span className="text-green-600 mr-2">Auto Saved!</span>
                    )}
                    <button className="bg-white text-blue-600 border border-blue-300 px-3 py-1 rounded mr-2" onClick={testDesign}>
                        Send Test
                    </button>
                    <button className="bg-navy-900 text-black border-black text-md cursor-pointer font-semibold px-3 py-1 rounded" onClick={() => { saveDesign(); handleNavigate(); }}>
                        Save and exit
                    </button>
                </div>
            </header>

            <div className="flex-grow">
                <ReactQuill
                    ref={quillRef}
                    value={editorContent}
                    onChange={handleChange}
                    modules={{
                        toolbar: [
                            [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
                            ['bold', 'italic', 'underline', 'strike'],
                            ['blockquote', 'code-block'],
                            [{ 'list': 'ordered' }, { 'list': 'bullet' }],
                            [{ 'script': 'sub' }, { 'script': 'super' }],
                            [{ 'indent': '-1' }, { 'indent': '+1' }],
                            [{ 'direction': 'rtl' }],
                            [{ 'size': ['small', false, 'large', 'huge'] }],
                            [{ 'color': [] }, { 'background': [] }],
                            [{ 'font': [] }],
                            [{ 'align': [] }],
                            ['link', 'image', 'video'],
                            ['clean']
                        ],
                        clipboard: {
                            matchVisual: false,
                        },

                    }}
                    style={{ height: 'calc(100vh - 5em)', overflowY: 'auto' }}
                />
            </div>

            <SendTestEmail isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} template_id={uuid as string} />
        </div>
    );
};

export default RichTextEditor;