import { useEffect, useState } from "react";
import { Tabs } from "antd";
import type { TabsProps } from "antd";
import useMetadata from "../../../hooks/useMetaData";
import { Helmet, HelmetProvider } from "react-helmet-async";
import MarketingTemplateDash from './marketing/marketingDashTemplate';
import TransactionalTemplateDash from './transactional/transactionalDashTemplate';

//type Tabtype = 'Transactional' | 'Marketing'

const TemplateBuilderDashComponent: React.FC = () => {
    const [activeTab, setActiveTab] = useState<string>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "Transactional" || storedTab === "Marketing") ? storedTab : "Transactional";
    });

    const metaData = useMetadata("TransactionalTemplates")

    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab);
        }
    }, []);

    useEffect(() => {
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);

    const handleTabChange = (key: string) => {
        setActiveTab(key);
    };

    const items: TabsProps['items'] = [
        {
            key: 'Transactional',
            label: 'Transactional Templates',
            children: <TransactionalTemplateDash />,
        },
        {
            key: 'Marketing',
            label: 'Marketing Templates',
            children: <MarketingTemplateDash />,
        },
    ];

    return (
        <HelmetProvider>
            <Helmet 
                {...metaData} 
                title={activeTab === "Marketing" ? "Marketing Templates - CrabMailer" : "Transactional Templates - CrabMailer"} 
            />

            <div className="p-6 max-w-full">
                <Tabs
                    activeKey={activeTab}
                    items={items}
                    onChange={handleTabChange}
                    size="large"
                    type="line"
                    className="custom-tabs"
                    tabBarStyle={{
                        marginBottom: 24,
                        borderBottom: '1px solid #f0f0f0'
                    }}
                />
            </div>
        </HelmetProvider>
    );
};

export default TemplateBuilderDashComponent;