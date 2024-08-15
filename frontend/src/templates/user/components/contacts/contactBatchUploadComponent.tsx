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

    const submitFile = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
        if (!selectedFile) {
            setError('Please select a valid CSV file.');
            return;
        }
        setSelectedCSVFile(selectedFile);
        batchContactUpload();
        getAllContacts();
        onClose()
        location.reload()
        setSelectedCSVFile(null)
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Upload Contact CSV">
            <>
                <p className="text-lg font-semibold mb-2">Select .csv or .xls file to import</p>
                <h5 className="text-blue-500 mb-2">How to format your .csv or excel file. Download the sample CSV below.</h5>

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
