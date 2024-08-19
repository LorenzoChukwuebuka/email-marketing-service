import { useState, useRef } from "react";
import { Modal } from "../../../../components";
import useContactStore from "../../../../store/userstore/contactStore";
import * as Yup from 'yup';



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

    const { setSelectedCSVFile, getAllContacts, batchContactUpload, isLoading } = useContactStore();

    const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file) {
            try {
                await fileValidationSchema.validate({ file });
                setSelectedFile(file);

                setError(null);
            } catch (validationError: any) {
                setError(validationError.message);
                setSelectedFile(null);
            }
        }
    };

    const handleButtonClick = () => {
        fileInputRef.current?.click();
    };

    const submitFile = async (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
        if (!selectedFile) {
            setError('Please select a valid CSV file.');
            return;
        }
        setSelectedCSVFile(selectedFile);
        await batchContactUpload();
        await getAllContacts();
        onClose()
        location.reload()
        setSelectedCSVFile(null)
    };

    const handleDownload = () => {
        // Create a blob with the CSV content
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
        <Modal isOpen={isOpen} onClose={onClose} title="Upload Contact CSV">
            <>
                <p className="text-lg font-semibold mb-2">Select .csv or .xls file to import</p>
                <h5 className="text-blue-500 cursor-pointer mb-2" onClick={handleDownload}>How to format your .csv or excel file. Download the sample CSV below.</h5>

                {error && <p className="text-red-600">{error}</p>}

                {selectedFile ? (
                    <div className="mb-4">
                        <p className="text-green-600">Selected file: {selectedFile.name}</p>
                        <button
                            className="mt-2 bg-gray-300 px-4 py-2 rounded-md transition duration-300"
                            onClick={handleButtonClick}
                        >
                            Choose a different file
                        </button>
                    </div>
                ) : (
                    <button
                        className="bg-gray-300 px-4 py-2 rounded-md transition duration-300"
                        onClick={handleButtonClick}
                    >
                        Select File
                    </button>
                )}

                <input
                    type="file"
                    ref={fileInputRef}
                    className="hidden"
                    accept=".csv, .xls, .xlsx"
                    onChange={handleFileChange}
                />

                <div className="flex justify-end space-x-2 mt-4">
                    {isLoading ? (
                        <button
                            type="button"
                            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
                        >
                            Please wait ...
                        </button>
                    ) : (
                        <button
                            type="button"
                            onClick={submitFile}
                            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
                        >
                            Upload file
                        </button>
                    )}
                    <button
                        type="button"
                        onClick={onClose}
                        className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
                    >
                        Cancel
                    </button>
                </div>
            </>
        </Modal>
    );
};

export default ContactUpload;
