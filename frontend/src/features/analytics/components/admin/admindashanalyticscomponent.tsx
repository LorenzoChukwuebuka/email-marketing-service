import { useEffect, useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { 
    Card, 
    Row, 
    Col, 
    Statistic, 
    Progress, 
    Typography, 
    Space, 
    Badge,
    Alert,
    Spin
} from 'antd';
import {
    DatabaseOutlined,
    HddOutlined,
    DashboardOutlined,
    PlayCircleOutlined,
    WarningOutlined
} from '@ant-design/icons';
import { 
    Server,
   
} from 'lucide-react';

const { Title, Text } = Typography;

const SystemMonitorDashboard = () => {
    const [historicalData, setHistoricalData] = useState<Record<string, any>[]>([]);
    const [isConnected, setIsConnected] = useState(false);
    const [currentData, setCurrentData] = useState({
        system: {
            operating_system: '',
            platform: '',
            hostname: '',
            num_processes: 0,
            total_memory: '',
            free_memory: '',
            used_memory_percent: ''
        },
        disk: {
            total_disk_space: '',
            used_disk_space: '',
            free_disk_space: '',
            used_disk_percent: ''
        },
        cpu: {
            model_name: '',
            family: '',
            speed: '',
            cpu_usage: []
        }
    });

    useEffect(() => {
        const ws = new WebSocket(import.meta.env.VITE_SOCKET_SERVER as string);

        ws.onopen = () => {
            setIsConnected(true);
        };

        ws.onclose = () => {
            setIsConnected(false);
        };

        ws.onerror = () => {
            setIsConnected(false);
        };

        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            setCurrentData(data);

            setHistoricalData(prev => {
                const newData = [...prev, {
                    time: new Date().toLocaleTimeString(),
                    memoryUsage: parseFloat(data.system.used_memory_percent),
                    diskUsage: parseFloat(data.disk.used_disk_percent),
                    cpuUsage: parseFloat(data.cpu.cpu_usage[0]?.split(': ')[1] || '0')
                }].slice(-20);
                return newData;
            });
        };

        return () => ws.close();
    }, []);

    const getUsageColor = (usage: number) => {
        if (usage < 50) return '#52c41a';
        if (usage < 80) return '#faad14';
        return '#ff4d4f';
    };

    const getUsageStatus = (usage: number) => {
        if (usage < 50) return 'success';
        if (usage < 80) return 'warning';
        return 'exception';
    };

    const memoryUsage = parseFloat(currentData.system.used_memory_percent) || 0;
    const diskUsage = parseFloat(currentData.disk.used_disk_percent) || 0;
    //@ts-ignore
    const cpuUsage = parseFloat(currentData.cpu.cpu_usage[0]?.split(': ')[1] || '0') || 0;

    return (
        <div className="p-6 max-w-7xl mx-auto space-y-6">
            {/* Header */}
            <div className="flex justify-between items-center mb-6">
                <div>
                    <Title level={2} className="!mb-2">
                        <DashboardOutlined className="mr-3" />
                        System Resource Monitor
                    </Title>
                    <Text type="secondary">Real-time system performance monitoring</Text>
                </div>
                <div className="flex items-center space-x-4">
                    <Badge 
                        status={isConnected ? "processing" : "error"} 
                        text={isConnected ? "Connected" : "Disconnected"} 
                    />
                    {!isConnected && (
                        <Alert
                            message="Connection Lost"
                            description="Attempting to reconnect..."
                            type="warning"
                            showIcon
                            closable
                        />
                    )}
                </div>
            </div>

            {/* Stats Grid */}
            <Row gutter={[16, 16]}>
                <Col xs={24} sm={12} lg={6}>
                    <Card className="h-full hover:shadow-lg transition-shadow duration-300">
                        <Statistic
                            title={
                                <Space>
                                    <DatabaseOutlined className="text-blue-500" />
                                    Memory Usage
                                </Space>
                            }
                            value={memoryUsage}
                            suffix="%"
                            precision={1}
                            valueStyle={{ color: getUsageColor(memoryUsage) }}
                        />
                        <Progress
                            percent={memoryUsage}
                            strokeColor={getUsageColor(memoryUsage)}
                            size="small"
                            className="mt-2"
                        />
                        <Text type="secondary" className="text-xs block mt-2">
                            Free: {currentData.system.free_memory} / Total: {currentData.system.total_memory}
                        </Text>
                    </Card>
                </Col>

                <Col xs={24} sm={12} lg={6}>
                    <Card className="h-full hover:shadow-lg transition-shadow duration-300">
                        <Statistic
                            title={
                                <Space>
                                    <HddOutlined className="text-green-500" />
                                    Disk Usage
                                </Space>
                            }
                            value={diskUsage}
                            suffix="%"
                            precision={1}
                            valueStyle={{ color: getUsageColor(diskUsage) }}
                        />
                        <Progress
                            percent={diskUsage}
                            strokeColor={getUsageColor(diskUsage)}
                            size="small"
                            className="mt-2"
                        />
                        <Text type="secondary" className="text-xs block mt-2">
                            Free: {currentData.disk.free_disk_space} / Total: {currentData.disk.total_disk_space}
                        </Text>
                    </Card>
                </Col>

                <Col xs={24} sm={12} lg={6}>
                    <Card className="h-full hover:shadow-lg transition-shadow duration-300">
                        <Statistic
                            title={
                                <Space>
                                    <PlayCircleOutlined className="text-orange-500" />
                                    CPU Usage
                                </Space>
                            }
                            value={cpuUsage}
                            suffix="%"
                            precision={1}
                            valueStyle={{ color: getUsageColor(cpuUsage) }}
                        />
                        <Progress
                            percent={cpuUsage}
                            strokeColor={getUsageColor(cpuUsage)}
                            size="small"
                            className="mt-2"
                        />
                        <Text type="secondary" className="text-xs block mt-2">
                            {currentData.cpu.model_name || 'Unknown CPU'}
                        </Text>
                    </Card>
                </Col>

                <Col xs={24} sm={12} lg={6}>
                    <Card className="h-full hover:shadow-lg transition-shadow duration-300">
                        <Statistic
                            title={
                                <Space>
                                    <Server className="text-purple-500" />
                                    System Info
                                </Space>
                            }
                            value={currentData.system.platform || 'Unknown'}
                            valueStyle={{ fontSize: '18px' }}
                        />
                        <div className="mt-4">
                            <Text type="secondary" className="text-xs block">
                                Processes: <span className="font-semibold">{currentData.system.num_processes}</span>
                            </Text>
                            <Text type="secondary" className="text-xs block">
                                Hostname: <span className="font-semibold">{currentData.system.hostname}</span>
                            </Text>
                        </div>
                    </Card>
                </Col>
            </Row>

            {/* Chart Section */}
            <Card 
                title={
                    <Space>
                        <WarningOutlined className="text-blue-500" />
                        Resource Usage History
                    </Space>
                }
                className="shadow-lg"
            >
                {historicalData.length > 0 ? (
                    <ResponsiveContainer width="100%" height={400}>
                        <LineChart
                            data={historicalData}
                            margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                        >
                            <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                            <XAxis 
                                dataKey="time" 
                                tick={{ fontSize: 12 }}
                                stroke="#666"
                            />
                            <YAxis 
                                tick={{ fontSize: 12 }}
                                stroke="#666"
                            />
                            <Tooltip 
                                contentStyle={{
                                    backgroundColor: '#fff',
                                    border: '1px solid #d9d9d9',
                                    borderRadius: '6px',
                                    boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                }}
                            />
                            <Legend />
                            <Line
                                type="monotone"
                                dataKey="memoryUsage"
                                stroke="#1890ff"
                                strokeWidth={2}
                                name="Memory Usage %"
                                dot={{ fill: '#1890ff', strokeWidth: 2, r: 3 }}
                            />
                            <Line
                                type="monotone"
                                dataKey="diskUsage"
                                stroke="#52c41a"
                                strokeWidth={2}
                                name="Disk Usage %"
                                dot={{ fill: '#52c41a', strokeWidth: 2, r: 3 }}
                            />
                            <Line
                                type="monotone"
                                dataKey="cpuUsage"
                                stroke="#faad14"
                                strokeWidth={2}
                                name="CPU Usage %"
                                dot={{ fill: '#faad14', strokeWidth: 2, r: 3 }}
                            />
                        </LineChart>
                    </ResponsiveContainer>
                ) : (
                    <div className="flex justify-center items-center h-64">
                        <Spin size="large" />
                        <Text className="ml-4">Loading chart data...</Text>
                    </div>
                )}
            </Card>

            {/* CPU Details */}
            <Card 
                title={
                    <Space>
                        <PlayCircleOutlined className="text-orange-500" />
                        CPU Details
                    </Space>
                }
                className="shadow-lg"
            >
                <Row gutter={[24, 24]}>
                    <Col xs={24} md={12}>
                        <Title level={4} className="mb-4">CPU Information</Title>
                        <div className="space-y-4">
                            <Card size="small" className="bg-gray-50">
                                <Statistic
                                    title="Model"
                                    value={currentData.cpu.model_name || 'Unknown'}
                                    valueStyle={{ fontSize: '16px', fontWeight: 'bold' }}
                                />
                            </Card>
                            <Card size="small" className="bg-gray-50">
                                <Statistic
                                    title="Family"
                                    value={currentData.cpu.family || 'Unknown'}
                                    valueStyle={{ fontSize: '16px', fontWeight: 'bold' }}
                                />
                            </Card>
                            <Card size="small" className="bg-gray-50">
                                <Statistic
                                    title="Speed"
                                    value={currentData.cpu.speed || 'Unknown'}
                                    valueStyle={{ fontSize: '16px', fontWeight: 'bold' }}
                                />
                            </Card>
                        </div>
                    </Col>
                    <Col xs={24} md={12}>
                        <Title level={4} className="mb-4">Core Usage</Title>
                        <div className="space-y-3">
                            {currentData.cpu.cpu_usage.map((usage, index) => {
                                //@ts-ignore
                                const usageValue = parseFloat(usage.split(': ')[1] || '0');
                                return (
                                    <div key={index} className="flex items-center gap-4">
                                        <Text className="w-16 text-sm font-medium">
                                            CPU {index}
                                        </Text>
                                        <div className="flex-1">
                                            <Progress
                                                percent={usageValue}
                                                strokeColor={getUsageColor(usageValue)}
                                                size="small"
                                                status={getUsageStatus(usageValue) as any}
                                            />
                                        </div>
                                        <Text className="w-16 text-sm font-mono">
                                            {usageValue.toFixed(1)}%
                                        </Text>
                                    </div>
                                );
                            })}
                        </div>
                    </Col>
                </Row>
            </Card>
        </div>
    );
};

export default SystemMonitorDashboard;