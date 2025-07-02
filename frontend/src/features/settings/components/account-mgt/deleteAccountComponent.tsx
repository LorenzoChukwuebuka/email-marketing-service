import React, { useState } from 'react';
import { AlertTriangle, Trash2, X, ShieldCheck } from 'lucide-react';
import useAuthStore from "../../../auth/store/auth.store";

const DeleteAccountComponent: React.FC = () => {
    const { deleteUser, cancelDelete } = useAuthStore();
    const [isDeleting, setIsDeleting] = useState(false);
    const [isCanceling, setIsCanceling] = useState(false);
    const [showConfirmation, setShowConfirmation] = useState(false);

    const handleDeleteAccount = async () => {
        if (!showConfirmation) {
            setShowConfirmation(true);
            return;
        }

        setIsDeleting(true);
        try {
            await deleteUser();
        } catch (error) {
            console.log(error);
        } finally {
            setIsDeleting(false);
        }
    };

    const handleCancelDelete = async () => {
        setIsCanceling(true);
        try {
            await cancelDelete();
        } catch (error) {
            console.log(error);
        } finally {
            setIsCanceling(false);
        }
    };

    const resetConfirmation = () => {
        setShowConfirmation(false);
    };

    return (
        <div className="max-w-4xl mx-auto p-8">
            {/* Header Section */}
            <div className="mb-8">
                <div className="flex items-center gap-3 mb-4">
                    <div className="p-2 bg-red-100 rounded-lg">
                        <AlertTriangle className="w-6 h-6 text-red-600" />
                    </div>
                    <div>
                        <h1 className="text-3xl font-bold text-gray-900">Delete Account</h1>
                        <p className="text-gray-500 mt-1">Permanently remove your account and data</p>
                    </div>
                </div>
            </div>

            {/* Warning Card */}
            <div className="bg-gradient-to-r from-red-50 to-red-100 border border-red-200 rounded-2xl p-8 mb-8 shadow-sm">
                <div className="flex items-start gap-4">
                    <div className="flex-shrink-0">
                        <div className="w-12 h-12 bg-red-200 rounded-full flex items-center justify-center">
                            <AlertTriangle className="w-6 h-6 text-red-700" />
                        </div>
                    </div>
                    <div className="flex-1">
                        <h3 className="text-xl font-semibold text-red-900 mb-3">
                            This action cannot be undone
                        </h3>
                        <div className="space-y-3 text-red-800">
                            <p className="leading-relaxed">
                                Deleting your account will permanently remove all data associated with your account, including:
                            </p>
                            <ul className="space-y-2 ml-4">
                                <li className="flex items-center gap-2">
                                    <div className="w-1.5 h-1.5 bg-red-600 rounded-full"></div>
                                    All projects and their associated data
                                </li>
                                <li className="flex items-center gap-2">
                                    <div className="w-1.5 h-1.5 bg-red-600 rounded-full"></div>
                                    Inboxes and communication history
                                </li>
                                <li className="flex items-center gap-2">
                                    <div className="w-1.5 h-1.5 bg-red-600 rounded-full"></div>
                                    Domain configurations and settings
                                </li>
                                <li className="flex items-center gap-2">
                                    <div className="w-1.5 h-1.5 bg-red-600 rounded-full"></div>
                                    User preferences and customizations
                                </li>
                            </ul>
                            <div className="mt-4 p-4 bg-red-200/50 rounded-lg border border-red-300">
                                <p className="text-sm font-medium">
                                    <ShieldCheck className="w-4 h-4 inline mr-2" />
                                    After clicking "Delete Account", we'll send you a confirmation email to complete the process.
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Action Buttons */}
            <div className="bg-white border border-gray-200 rounded-2xl p-8 shadow-sm">
                {!showConfirmation ? (
                    <div className="space-y-6">
                        <h3 className="text-lg font-semibold text-gray-900">Ready to proceed?</h3>
                        <div className="flex flex-col sm:flex-row gap-4">
                            <button
                                onClick={handleDeleteAccount}
                                disabled={isDeleting}
                                className="flex items-center justify-center gap-2 px-6 py-3 bg-red-600 hover:bg-red-700 text-white font-semibold rounded-xl transition-all duration-200 transform hover:scale-105 hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
                            >
                                <Trash2 className="w-5 h-5" />
                                {isDeleting ? 'Processing...' : 'Delete Account'}
                            </button>
                            
                            <button
                                onClick={handleCancelDelete}
                                disabled={isCanceling}
                                className="flex items-center justify-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl transition-all duration-200 transform hover:scale-105 hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
                            >
                                <X className="w-5 h-5" />
                                {isCanceling ? 'Processing...' : 'Cancel Delete Request'}
                            </button>
                        </div>
                    </div>
                ) : (
                    <div className="space-y-6">
                        <div className="text-center">
                            <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
                                <AlertTriangle className="w-8 h-8 text-red-600" />
                            </div>
                            <h3 className="text-xl font-bold text-gray-900 mb-2">Are you absolutely sure?</h3>
                            <p className="text-gray-600 max-w-md mx-auto">
                                This will permanently delete your account and all associated data. This action cannot be reversed.
                            </p>
                        </div>
                        
                        <div className="flex flex-col sm:flex-row gap-4 justify-center">
                            <button
                                onClick={handleDeleteAccount}
                                disabled={isDeleting}
                                className="flex items-center justify-center gap-2 px-8 py-3 bg-red-600 hover:bg-red-700 text-white font-bold rounded-xl transition-all duration-200 shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                <Trash2 className="w-5 h-5" />
                                {isDeleting ? 'Deleting Account...' : 'Yes, Delete Forever'}
                            </button>
                            
                            <button
                                onClick={resetConfirmation}
                                className="flex items-center justify-center gap-2 px-8 py-3 bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold rounded-xl transition-all duration-200"
                            >
                                <X className="w-5 h-5" />
                                Cancel
                            </button>
                        </div>
                    </div>
                )}
            </div>

            {/* Support Section */}
            <div className="mt-8 text-center">
                <p className="text-gray-500 text-sm">
                    Need help or have questions? 
                    <a href="#" className="text-blue-600 hover:text-blue-700 font-medium ml-1 hover:underline">
                        Contact our support team
                    </a>
                </p>
            </div>
        </div>
    );
};

export default DeleteAccountComponent;