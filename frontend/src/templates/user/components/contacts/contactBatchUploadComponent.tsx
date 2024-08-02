import { useState, useRef } from "react";
import { Modal } from "../../../../components";

interface ContactUploadProps {
    isOpen: boolean;
    onClose: () => void;

}

const ContactUpload: React.FC<ContactUploadProps> = ({ isOpen, onClose }) => {
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const fileInputRef = useRef<HTMLInputElement | null>(null);

    const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file) {
            setSelectedFile(file);
        }
    };

    const handleButtonClick = () => {
        fileInputRef.current?.click();
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Upload Contact CSV">
            <>
                <p className="text-lg font-semibold mb-2">Select .csv or .xls file to import</p>
                <h5 className="text-blue-500 mb-2">How to format your .csv or excel file. Download the sample CSV below.</h5>

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
