import { useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import {
    Modal,
    Pagination,
    Input,
    Button,
    Card,
    Tag,
    Dropdown,
    Space,
    Spin,
    Typography,
    Avatar,
    Flex
} from 'antd';
import {
    PlusOutlined,
    SearchOutlined,
    EyeOutlined,
    EditOutlined,
    DeleteOutlined,
    MoreOutlined,
    FileTextOutlined
} from '@ant-design/icons';
import type { MenuProps } from 'antd';
import useDebounce from "../../../../hooks/useDebounce";
import EmptyState from "../../../../components/emptyStateComponent";
import { Template } from '../../interface/email-templates.interface';
import { BaseEntity } from '../../../../../../frontend/src/interface/baseentity.interface';
import useTemplateStore from "../../store/template.store";
import { useTransactionalTemplateQuery } from "../../hooks/useTransactionTemplateQuery";

const { Search } = Input;
const { Title, Text } = Typography;

const TransactionalTemplateDash: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { deleteTemplate } = useTemplateStore();
    const [previewTemplate, setPreviewTemplate] = useState<Template & BaseEntity | null>(null);
    const [searchQuery, setSearchQuery] = useState<string>("");
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);

    const navigate = useNavigate();
    const debouncedSearchQuery = useDebounce(searchQuery, 300);

    const { data: _templateData, isLoading } = useTransactionalTemplateQuery(currentPage, pageSize, debouncedSearchQuery);

    const tempData = useMemo<any[]>(() => {
        if (!_templateData) return [];

        // If it's a paginated response
        if ('data' in _templateData.payload) {
            return _templateData.payload.data;
        }

        // If it's a single template, wrap it in an array
        return [_templateData.payload];
    }, [_templateData]);

    const openPreview = (template: (Template & BaseEntity)) => {
        setPreviewTemplate(template);
        setIsModalOpen(true);
    };

    const onPageChange = (page: number, size: number) => {
        setCurrentPage(page);
        setPageSize(size);
    };

    const handleSearch = (value: string) => {
        setSearchQuery(value);
    };

    const handleNavigate = (template: (Template & BaseEntity)) => {
        const editorType = template.editor_type;

        let redirectUrl = "";
        switch (editorType) {
            case "html-editor":
                redirectUrl = `/editor/2?type=t&uuid=${template.id}`;
                break;
            case "drag-and-drop":
                redirectUrl = `/editor/1?type=t&uuid=${template.id}`;
                break;
            case "rich-text":
                redirectUrl = `/editor/3?type=t&uuid=${template.id}`;
                break;
            default:
                console.log("Unknown editor type:", editorType);
                return;
        }

        window.location.href = redirectUrl;
    };

    const deleteTempl = async (template: (Template & BaseEntity)) => {
        await deleteTemplate(template.id);
    };

    const getDropdownItems = (template: Template & BaseEntity): MenuProps['items'] => [
        {
            key: 'preview',
            label: 'Preview',
            icon: <EyeOutlined />,
            onClick: () => openPreview(template),
        },
        {
            key: 'edit',
            label: 'Edit',
            icon: <EditOutlined />,
            onClick: () => handleNavigate(template),
        },
        {
            type: 'divider',
        },
        {
            key: 'delete',
            label: 'Delete',
            icon: <DeleteOutlined />,
            danger: true,
            onClick: () => deleteTempl(template),
        },
    ];

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString("en-US", {
            timeZone: "UTC",
            year: "numeric",
            month: "long",
            day: "numeric",
            hour: "numeric",
            minute: "numeric",
            second: "numeric",
        });
    };

    return (
        <>
            {isLoading ? (
                <div className="flex items-center justify-center mt-20">
                    <Spin size="large" />
                </div>
            ) : (
                <>
                    {/* Header Section */}
                    <Card className="mt-6 mb-6">
                        <Flex justify="space-between" align="center" wrap="wrap" gap="middle">
                            <Button
                                type="primary"
                                icon={<PlusOutlined />}
                                size="large"
                            >
                                <Link to="/app/templates/transactional">Create Transactional Template</Link>
                            </Button>

                            <Search
                                placeholder="Search templates..."
                                allowClear
                                size="large"
                                style={{ width: 300 }}
                                onSearch={handleSearch}
                                onChange={(e) => handleSearch(e.target.value)}
                                prefix={<SearchOutlined />}
                            />
                        </Flex>
                    </Card>

                    {/* Templates List */}
                    <div className="space-y-4">
                        {Array.isArray(tempData) && tempData.length > 0 ? (
                            <>
                                {tempData.map((template, index) => (
                                    <Card
                                        key={template.id || index}
                                        hoverable
                                        className="shadow-sm"
                                        bodyStyle={{ padding: '20px' }}
                                    >
                                        <Flex align="center" gap="middle">
                                            <Avatar
                                                size={48}
                                                icon={<FileTextOutlined />}
                                                style={{ backgroundColor: '#f0f0f0', color: '#666' }}
                                            />

                                            <div className="flex-grow">
                                                <Title level={4} style={{ margin: 0, marginBottom: 4 }}>
                                                    {template.template_name}
                                                </Title>
                                                <Text type="secondary">
                                                    ID - {index + 1} • {formatDate(template.created_at)}
                                                </Text>
                                                <div className="mt-2">
                                                    <Space>
                                                        <Button
                                                            type="link"
                                                            size="small"
                                                            icon={<EyeOutlined />}
                                                            onClick={() => openPreview(template)}
                                                        >
                                                            Preview
                                                        </Button>
                                                        <Button
                                                            type="link"
                                                            size="small"
                                                            icon={<EditOutlined />}
                                                            onClick={() => handleNavigate(template)}
                                                        >
                                                            Edit
                                                        </Button>
                                                    </Space>
                                                </div>
                                            </div>

                                            <Flex align="center" gap="middle">
                                                <Tag color="default">Draft</Tag>
                                                <Dropdown
                                                    menu={{ items: getDropdownItems(template) }}
                                                    trigger={['click']}
                                                    placement="bottomRight"
                                                >
                                                    <Button
                                                        type="text"
                                                        icon={<MoreOutlined />}
                                                        size="small"
                                                    />
                                                </Dropdown>
                                            </Flex>
                                        </Flex>
                                    </Card>
                                ))}

                                {/* Pagination */}
                                <div className="flex justify-center items-center mt-6 mb-5">
                                    <Pagination
                                        current={currentPage}
                                        pageSize={pageSize}
                                        total={_templateData?.payload?.total || tempData.length}
                                        onChange={onPageChange}
                                        showSizeChanger
                                        showQuickJumper
                                        showTotal={(total, range) =>
                                            `${range[0]}-${range[1]} of ${total} templates`
                                        }
                                        pageSizeOptions={["10", "20", "50", "100"]}
                                    />
                                </div>
                            </>
                        ) : (
                            <EmptyState
                                title="You have not created any Template"
                                description="Create and easily send transactional emails to your audience"
                                icon={<i className="bi bi-emoji-frown text-xl"></i>}
                                buttonText="Create Template"
                                onButtonClick={() => navigate("/app/templates/transactional")}
                            />
                        )}
                    </div>

                    {/* Preview Modal */}
                    <Modal
                        title="Preview Template"
                        open={isModalOpen}
                        onCancel={() => setIsModalOpen(false)}
                        footer={null}
                        width={800}
                        centered
                    >
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
                    </Modal>
                </>
            )}
        </>
    );
};

export default TransactionalTemplateDash;