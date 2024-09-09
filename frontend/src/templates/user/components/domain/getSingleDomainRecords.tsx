import { useParams } from "react-router-dom"
import useDomainStore, { DomainRecord } from "../../../../store/userstore/domainStore"
import { FormEvent, useEffect, useState } from "react"

type Props = {
    label: string
    type: string
    recordNamePlaceholder: string
    value: string
}

const RecordInput = ({ label, type, recordNamePlaceholder, value }: Props) => {
    const [copied, setCopied] = useState({ name: false, value: false });

    const handleCopy = (text: string, field: "name" | "value") => {
        navigator.clipboard.writeText(text)
            .then(() => setCopied(prev => ({ ...prev, [field]: true })))
            .catch(err => console.error("Failed to copy text: ", err));

        setTimeout(() => setCopied(prev => ({ ...prev, [field]: false })), 2000);
    };

    return (
        <div className="border p-4 rounded-lg bg-white mb-6">
            <label className="block text-gray-700 text-sm font-bold mb-2">{label}</label>
            <div className="mb-4 relative">
                <span className="block text-gray-600 text-sm mb-1">Type</span>
                <input
                    type="text"
                    value={type}
                    readOnly
                    className="bg-gray-100 border border-gray-300 rounded-md py-2 px-4 w-full"
                />
            </div>
            <div className="mb-4 relative">
                <span className="block text-gray-600 text-sm mb-1">Record name</span>
                <div className="flex items-center">
                    <input
                        type="text"
                        value={recordNamePlaceholder}
                        readOnly
                        className="bg-white border border-gray-300 rounded-md py-2 px-4 w-full pr-16" // Added padding to the right
                    />
                    <button
                        onClick={() => handleCopy(recordNamePlaceholder, "name")}
                        className="ml-2 text-blue-500 text-sm flex-shrink-0" // Adjusted margin to the left
                    >
                        {copied.name ? "Copied!" : "Copy"}
                    </button>
                </div>
            </div>
            <div className="relative">
                <span className="block text-gray-600 text-sm mb-1">Value</span>
                <div className="flex items-center">
                    <input
                        type="text"
                        value={value}
                        readOnly
                        className="bg-gray-100 border border-gray-300 rounded-md py-2 px-4 w-full pr-16" // Added padding to the right
                    />
                    <button
                        onClick={() => handleCopy(value, "value")}
                        className="ml-2 text-blue-500 text-sm flex-shrink-0" // Adjusted margin to the left
                    >
                        {copied.value ? "Copied!" : "Copy"}
                    </button>
                </div>
            </div>
        </div>
    );
};


const DNSAuthenticationRecords: React.FC = () => {
    const { getDomain, domainData, authenticateDomain } = useDomainStore()
    const [domainD, setDomainD] = useState<DomainRecord | null>(null)

    const { id } = useParams<{ id: string }>() as { id: string }

    useEffect(() => {
        const fetchD = async () => {
            await getDomain(id)
        }
        fetchD()
    }, [getDomain])


    useEffect(() => {
        if (domainData) {
            setDomainD(domainData as DomainRecord)
        }
    }, [domainData])


    const authenticate = async (e: FormEvent<HTMLButtonElement>) => {
        e.preventDefault()
        await authenticateDomain(id)
    }


    return (
        <div className="max-w-xl  p-6">

            <h1 className="font-semibold text-lg mb-4"> DNS records for domain authentication </h1>

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
                value={"v=DKIM1; k=rsa; p="+domainD?.dkim_public_key as string}
            />

            <RecordInput
                label="DMARC record"
                type="TXT"
                recordNamePlaceholder="_dmarc"
                value={domainD?.dmarc_record as string}
            />

            <div className="flex justify-between items-center mt-4">
                <button className="py-2 px-4 rounded-md bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 transition duration-200" onClick={(e) => authenticate(e)}>
                    Authenticate Domain
                </button>

                <button className="py-2 px-4 rounded-md bg-gray-200 text-gray-700 text-sm font-medium hover:bg-gray-300 transition duration-200" onClick={() => window.history.back()}>
                    Go back
                </button>
            </div>


        </div>


    );
};

export default DNSAuthenticationRecords;
