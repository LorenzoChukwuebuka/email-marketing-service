import React, { useEffect, useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from 'recharts';
import { Server, HardDrive, Cpu, Database } from 'lucide-react';

const SystemMonitorDashboard = () => {
    const [historicalData, setHistoricalData] = useState([]);
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
        const ws = new WebSocket( import.meta.env.VITE_SOCKET_SERVER as string);

        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            setCurrentData(data);
            //@ts-ignore
            setHistoricalData(prev => {
                const newData = [...prev, {
                    time: new Date().toLocaleTimeString(),
                    memoryUsage: parseFloat(data.system.used_memory_percent),
                    diskUsage: parseFloat(data.disk.used_disk_percent),
                    cpuUsage: parseFloat(data.cpu.cpu_usage[0].split(': ')[1])
                }].slice(-20);
                return newData;
            });
        };

        return () => ws.close();
    }, []);


    type Stat = {
        title: string
        value: any
        icon: any
        secondaryValue: any
    }

    const StatCard = ({ title, value, icon: Icon, secondaryValue = null }: Stat) => (
        <div className="card bg-base-100 shadow">
            <div className="card-body">
                <div className="flex justify-between items-center">
                    <h2 className="card-title text-sm">{title}</h2>
                    <Icon className="w-4 h-4" />
                </div>
                <div className="text-2xl font-bold">{value}</div>
                {secondaryValue && (
                    <p className="text-xs opacity-70">{secondaryValue}</p>
                )}
            </div>
        </div>
    );

    return (
        <div className="p-8 max-w-7xl mx-auto space-y-8">
            <div className="text-3xl font-bold">
                System Resource Monitor
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <StatCard
                    title="Memory Usage"
                    value={currentData.system.used_memory_percent}
                    icon={Database}
                    secondaryValue={`Free: ${currentData.system.free_memory} / Total: ${currentData.system.total_memory}`}
                />
                <StatCard
                    title="Disk Usage"
                    value={currentData.disk.used_disk_percent}
                    icon={HardDrive}
                    secondaryValue={`Free: ${currentData.disk.free_disk_space} / Total: ${currentData.disk.total_disk_space}`}
                />
                <StatCard
                    title="CPU"
                    //@ts-ignore
                    value={currentData.cpu.cpu_usage[0]?.split(': ')[1] || '0%'}
                    icon={Cpu}
                    secondaryValue={currentData.cpu.model_name}
                />
                <StatCard
                    title="System Info"
                    value={currentData.system.platform}
                    icon={Server}
                    secondaryValue={`Processes: ${currentData.system.num_processes}`}
                />
            </div>

            {/* Chart Section */}
            <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                    <h2 className="card-title">Resource Usage History</h2>
                    <LineChart
                        width={800}
                        height={400}
                        data={historicalData}
                        margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                        className="mx-auto"
                    >
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="time" />
                        <YAxis />
                        <Tooltip />
                        <Legend />
                        <Line
                            type="monotone"
                            dataKey="memoryUsage"
                            stroke="#8884d8"
                            name="Memory Usage %"
                        />
                        <Line
                            type="monotone"
                            dataKey="diskUsage"
                            stroke="#82ca9d"
                            name="Disk Usage %"
                        />
                        <Line
                            type="monotone"
                            dataKey="cpuUsage"
                            stroke="#ffc658"
                            name="CPU Usage %"
                        />
                    </LineChart>
                </div>
            </div>

            {/* CPU Details */}
            <div className="card bg-base-100 shadow">
                <div className="card-body">
                    <h2 className="card-title">CPU Details</h2>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <h3 className="font-semibold mb-2">CPU Information</h3>
                            <div className="space-y-2">
                                <div className="stats shadow">
                                    <div className="stat">
                                        <div className="stat-title">Model</div>
                                        <div className="stat-value text-lg">{currentData.cpu.model_name}</div>
                                    </div>
                                </div>
                                <div className="stats shadow">
                                    <div className="stat">
                                        <div className="stat-title">Family</div>
                                        <div className="stat-value text-lg">{currentData.cpu.family}</div>
                                    </div>
                                </div>
                                <div className="stats shadow">
                                    <div className="stat">
                                        <div className="stat-title">Speed</div>
                                        <div className="stat-value text-lg">{currentData.cpu.speed}</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div>
                            <h3 className="font-semibold mb-2">Core Usage</h3>
                            <div className="space-y-2">
                                {currentData.cpu.cpu_usage.map((usage, index) => {
                                    //@ts-ignore
                                    const usageValue = parseFloat(usage.split(': ')[1]);
                                    return (
                                        <div key={index} className="flex items-center gap-2">
                                            <span className="text-sm w-20">{`CPU ${index}`}</span>
                                            <progress
                                                className="progress progress-primary w-full"
                                                value={usageValue}
                                                max="100"
                                            ></progress>
                                            <span className="text-sm w-16">{usageValue.toFixed(1)}%</span>
                                        </div>
                                    );
                                })}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default SystemMonitorDashboard;