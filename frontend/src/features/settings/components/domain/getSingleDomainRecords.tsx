import { useParams } from "react-router-dom";
import { useEffect, useMemo, useState } from "react";
import { Form, Input, Button, Card, Typography, message } from "antd";
import { CopyOutlined, CheckOutlined } from "@ant-design/icons";
import { DomainRecord } from '../../interface/domain.interface';
import useDomainStore from "../../store/domain.store";
import { useSingleDomainQuery } from "../../hooks/useDomainQuery";

const { Title } = Typography;

type RecordInputProps = {
    label: string;
    type: string;
    recordNamePlaceholder: string;
    value: string;
};

const RecordInput = ({ label, type, recordNamePlaceholder, value }: RecordInputProps) => {
    const [copied, setCopied] = useState({ name: false, value: false });

    const handleCopy = async (text: string, field: "name" | "value") => {
        try {
            await navigator.clipboard.writeText(text);
            setCopied(prev => ({ ...prev, [field]: true }));
            message.success("Copied to clipboard");
            setTimeout(() => setCopied(prev => ({ ...prev, [field]: false })), 2000);
        } catch (err) {
            console.error("Failed to copy text: ", err);
            message.error("Failed to copy");
        }
    };

    return (
        <Card className="mb-4">
            <Title level={5}>{label}</Title>

            <Form layout="vertical">
                <Form.Item label="Type">
                    <Input
                        value={type}
                        readOnly
                        className="bg-gray-50"
                    />
                </Form.Item>

                <Form.Item label="Record name">
                    <Input
                        value={recordNamePlaceholder}
                        readOnly
                        addonAfter={
                            <Button
                                type="text"
                                onClick={() => handleCopy(recordNamePlaceholder, "name")}
                                icon={copied.name ? <CheckOutlined /> : <CopyOutlined />}
                            />
                        }
                    />
                </Form.Item>

                <Form.Item label="Value">
                    <Input
                        value={value}
                        readOnly
                        className="bg-gray-50"
                        addonAfter={
                            <Button
                                type="text"
                                onClick={() => handleCopy(value, "value")}
                                icon={copied.value ? <CheckOutlined /> : <CopyOutlined />}
                            />
                        }
                    />
                </Form.Item>
            </Form>
        </Card>
    );
};

const DNSAuthenticationRecords: React.FC = () => {
    const { authenticateDomain } = useDomainStore();
    const [domainD, setDomainD] = useState<DomainRecord | null>(null);
    const { id } = useParams<{ id: string }>() as { id: string };
    const { data: domainData } = useSingleDomainQuery(id);
    const dData = useMemo(() => domainData?.payload || null, [domainData]);

    useEffect(() => {
        if (dData) {
            setDomainD(dData as DomainRecord);
        }
    }, [dData]);

    const authenticate = async () => {
        try {
            await authenticateDomain(id);
            message.success("Domain authentication initiated");
        } catch (error) {
            console.log(error)
            message.error("Failed to authenticate domain");
        }
    };

    return (
        <div className="max-w-xl p-6">
            <Title level={4} className="mb-6">
                DNS records for domain authentication
            </Title>

            <RecordInput
                label="CrabMailer Code"
                type="TXT"
                recordNamePlaceholder="Leave this field blank"
                value={domainD?.txt_record as string}
            />
            <RecordInput
                label="DKIM record"
                type="TXT"
                recordNamePlaceholder={domainD?.dkim_selector as string + "._domainkey"}
                value={"v=DKIM1; k=rsa; p=" + domainD?.dkim_public_key as string}
            />
            <RecordInput
                label="DMARC record"
                type="TXT"
                recordNamePlaceholder="_dmarc"
                value={domainD?.dmarc_record as string}
            />
            <RecordInput
                label="SPF record"
                type="TXT"
                recordNamePlaceholder="@"
                value={domainD?.spf_record as string}
            />
            <RecordInput
                label="MX record"
                type="TXT"
                recordNamePlaceholder="@"
                value={domainD?.mx_record as string}
            />

            <div className="flex justify-between items-center mt-6">
                <Button
                    type="primary"
                    onClick={authenticate}
                    className="bg-blue-600"
                >
                    Authenticate Domain
                </Button>

                <Button
                    onClick={() => window.history.back()}
                >
                    Go back
                </Button>
            </div>
        </div>
    );
};

export default DNSAuthenticationRecords;