import { useState, useRef } from "react";
import { Modal, Button } from "antd";
import * as Yup from 'yup';
import useContactStore from "../../store/contact.store";

// Yup schema for file validation
const fileValidationSchema = Yup.object().shape({
    file: Yup.mixed()
        .required('A file is required')
        .test(
            'fileType',
            'Only CSV files are supported',
            (value) => {
                const file = value as File;
                return file && file.type === 'text/csv';
            }
        )
});

interface ContactUploadProps {
    isOpen: boolean;
    onClose: () => void;
}

const ContactUpload: React.FC<ContactUploadProps> = ({ isOpen, onClose }) => {
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [error, setError] = useState<string | null>(null);
    const fileInputRef = useRef<HTMLInputElement | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(false)

    const { setSelectedCSVFile, batchContactUpload } = useContactStore();

    const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file) {
            try {
                setIsLoading(true)
                await fileValidationSchema.validate({ file });
                setSelectedFile(file);
                setError(null);
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
            } catch (validationError: any) {
                setError(validationError.message);
                setSelectedFile(null);
            } finally {
                setIsLoading(false)
            }
        }
    };

    const handleButtonClick = () => {
        fileInputRef.current?.click();
    };

    const submitFile = async () => {
        if (!selectedFile) {
            setError('Please select a valid CSV file.');
            return;
        }
        setSelectedCSVFile(selectedFile);
        await batchContactUpload();
        new Promise((resolve) => setTimeout(resolve, 3000));

        onClose();
        location.reload();
        setSelectedCSVFile(null);
    };

    const handleDownload = () => {
        const csvContent = 'data:text/csv;charset=utf-8,' + encodeURIComponent(
            'First Name,Last Name,Email,From\n' +
            'John,Doe,john.doe@example.com,Website\n' +
            'Jane,Smith,jane.smith@example.com,Referral\n' +
            'Michael,Johnson,michael.johnson@example.com,Trade Show\n' +
            'Emily,Brown,emily.brown@example.com,Social Media\n' +
            'David,Wilson,david.wilson@example.com,Newsletter\n' +
            'Sarah,Taylor,sarah.taylor@example.com,Website\n' +
            'Robert,Anderson,robert.anderson@example.com,Referral\n' +
            'Jennifer,Thomas,jennifer.thomas@example.com,Trade Show\n' +
            'William,Jackson,william.jackson@example.com,Social Media\n' +
            'Elizabeth,White,elizabeth.white@example.com,Newsletter'
        );
        const link = document.createElement('a');
        link.href = csvContent;
        link.download = 'sample.csv';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return (
        <Modal
            title="Upload Contact CSV"
            open={isOpen}
            onCancel={onClose}
            footer={[
                <Button key="cancel" onClick={onClose}>
                    Cancel
                </Button>,
                <Button
                    key="submit"
                    onClick={submitFile}
                    loading={isLoading}
                >
                    {isLoading ? 'Please wait...' : 'Upload file'}
                </Button>
            ]}
        >
            <div className="space-y-4">
                <p className="text-lg font-semibold">Select .csv or .xls file to import</p>

                <h5
                    className="text-blue-500 cursor-pointer"
                    onClick={handleDownload}
                >
                    How to format your .csv or excel file. Download the sample CSV below.
                </h5>

                {error && <p className="text-red-600">{error}</p>}

                {selectedFile ? (
                    <div>
                        <p className="text-green-600">Selected file: {selectedFile.name}</p>
                        <Button
                            className="mt-2"
                            onClick={handleButtonClick}
                        >
                            Choose a different file
                        </Button>
                    </div>
                ) : (
                    <Button onClick={handleButtonClick}>
                        Select File
                    </Button>
                )}

                <input
                    type="file"
                    ref={fileInputRef}
                    className="hidden"
                    accept=".csv, .xls, .xlsx"
                    onChange={handleFileChange}
                />
            </div>
        </Modal>
    );
};

export default ContactUpload;