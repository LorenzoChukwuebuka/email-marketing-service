import { useState, useMemo } from "react";
import { Table, Card, Typography,  Empty } from "antd";
import type { ColumnsType, TablePaginationConfig } from "antd/es/table";
import { useAllCampaignStatsQuery } from "../hooks/useAnalyticsQuery";
import { parseDate } from "../../../utils/utils";
import useDebounce from "../../../hooks/useDebounce";

const { Title } = Typography;
 


const AnalyticsTableComponent: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState("");
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const debouncedSearchQuery = useDebounce(searchQuery, 300);

    const { data: analyticsData, isLoading } = useAllCampaignStatsQuery(currentPage, pageSize, debouncedSearchQuery);

    const tableData = useMemo(() => analyticsData?.payload?.data, [analyticsData]);
    const totalItems = analyticsData?.payload?.total || 0;

    const handleSearch = (value: string) => {
        setSearchQuery(value);
        setCurrentPage(1);
    };

    const handleTableChange = (pagination: TablePaginationConfig) => {
        setCurrentPage(pagination.current || 1);
        setPageSize(pagination.pageSize || 20);
    };

    const columns: ColumnsType<any> = [
        {
            title: "Campaign Name",
            dataIndex: "name",
            key: "name",
        },
        {
            title: "Bounces",
            dataIndex: "bounces",
            key: "bounces",
        },
        {
            title: "Recipients",
            dataIndex: "recipients",
            key: "recipients",
        },
        {
            title: "Opened",
            dataIndex: "opened",
            key: "opened",
        },
        {
            title: "Clicked",
            dataIndex: "clicked",
            key: "clicked",
        },
        {
            title: "Sent Date",
            dataIndex: "sent_date",
            key: "sent_date",
            render: (date: string | null) =>
                date
                    ? parseDate(date).toLocaleString("en-US", {
                        timeZone: "UTC",
                        year: "numeric",
                        month: "long",
                        day: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                        second: "2-digit",
                    })
                    : "Not Sent",
        },
    ];

    return (
        <div className="p-6 bg-gray-50 mt-10 mb-5">
            <div className="max-w-7xl mx-auto">
                <Title level={2} className="text-gray-800 mb-4">
                    Email Campaign Analytics
                </Title>

                {/* <Card className="mb-6 shadow-sm border-0">
                    <div className="flex justify-between items-center">
                        <Search
                            placeholder="Search campaign name..."
                            value={searchQuery}
                            onChange={(e) => handleSearch(e.target.value)}
                            onSearch={handleSearch}
                            allowClear
                            size="large"
                            style={{ width: 400 }}
                            prefix={<SearchOutlined className="text-gray-400" />}
                        />
                        <div className="text-sm text-gray-500">{totalItems} campaigns total</div>
                    </div>
                </Card> */}

                <Card className="shadow-sm border-0">
                    <Table<any>
                        columns={columns}
                        dataSource={tableData}
                        rowKey={"id"}
                        loading={isLoading}
                        pagination={{
                            current: currentPage,
                            pageSize,
                            total: totalItems,
                            showSizeChanger: true,
                            showQuickJumper: true,
                            pageSizeOptions: ["10", "20", "50", "100"],
                            showTotal: (total, range) => `${range[0]}-${range[1]} of ${total} campaigns`,
                        }}
                        onChange={handleTableChange}
                        locale={{
                            emptyText: (
                                <Empty
                                    description={
                                        <div className="text-gray-500">
                                            {searchQuery
                                                ? `No results for "${searchQuery}"`
                                                : "No analytics data available"}
                                        </div>
                                    }
                                />
                            ),
                        }}
                        scroll={{ x: 800 }}
                        size="middle"
                        className="[&_.ant-table-thead>tr>th]:bg-gray-50"
                    />
                </Card>
            </div>
        </div>
    );
};

export default AnalyticsTableComponent;
