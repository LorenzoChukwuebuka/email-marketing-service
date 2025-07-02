import React, { useEffect, useState } from "react";
import { User, Mail, Building, Phone, Shield, Edit3, Save, X } from 'lucide-react';
import useAuthStore from "../../../auth/store/auth.store";
import { useUserDetailsQuery } from "../../../auth/hooks/useUserDetailsQuery";

const ProfileInformationComponent: React.FC = () => {
    const {
        setEditFormValues,
        editFormValues,
        editUserDetails,
    } = useAuthStore();

    const { data: userData } = useUserDetailsQuery();
    
    const [isEditing, setIsEditing] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [showSuccess, setShowSuccess] = useState(false);
    const [formData, setFormData] = useState({
        fullname: '',
        email: '',
        company: '',
        phonenumber: '',
    });
    const [errors, setErrors] = useState<Record<string, string>>({});

    const initEdit = () => {
        const newFormData = {
            fullname: userData?.payload.fullname || "",
            email: userData?.payload.email || "",
            company: userData?.payload.company?.companyname || userData?.payload.company || "",
            phonenumber: userData?.payload.phonenumber || "",
        };
        setFormData(newFormData);
        setEditFormValues(newFormData);
    };

    const handleInputChange = (field: string, value: string) => {
        setFormData(prev => ({ ...prev, [field]: value }));
        
        // Clear errors when user starts typing
        if (errors[field]) {
            setErrors(prev => ({ ...prev, [field]: '' }));
        }
    };

    const validateForm = (): boolean => {
        const newErrors: Record<string, string> = {};

        if (!formData.fullname.trim()) {
            newErrors.fullname = 'Full name is required';
        }

        if (!formData.phonenumber.trim()) {
            newErrors.phonenumber = 'Phone number is required';
        } else if (formData.phonenumber.length !== 11) {
            newErrors.phonenumber = 'Phone number must be exactly 11 digits';
        } else if (!/^\d+$/.test(formData.phonenumber)) {
            newErrors.phonenumber = 'Phone number must contain only digits';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSave = async () => {
        if (!validateForm()) return;

        setIsLoading(true);
        try {
            setEditFormValues(formData);
            await editUserDetails();
            setIsEditing(false);
            setShowSuccess(true);
            setTimeout(() => setShowSuccess(false), 3000);
            initEdit();
        } catch (error) {
            console.error("Failed to update user details:", error);
            setErrors({ submit: 'Failed to update profile. Please try again.' });
        } finally {
            setIsLoading(false);
        }
    };

    const handleCancel = () => {
        setIsEditing(false);
        setErrors({});
        initEdit(); // Reset to original values
    };

    useEffect(() => {
        if (userData) {
            initEdit();
        }
    }, [userData]);

    return (
        <div className="max-w-5xl mx-auto p-8">
            {/* Header Section */}
            <div className="mb-8">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <div className="p-2 bg-purple-100 rounded-lg">
                            <User className="w-6 h-6 text-purple-600" />
                        </div>
                        <div>
                            <h1 className="text-3xl font-bold text-gray-900">Profile Information</h1>
                            <p className="text-gray-500 mt-1">Manage your account details and preferences</p>
                        </div>
                    </div>
                    
                    {!isEditing && (
                        <button
                            onClick={() => setIsEditing(true)}
                            className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white font-medium rounded-xl transition-all duration-200 shadow-lg hover:shadow-xl"
                        >
                            <Edit3 className="w-4 h-4" />
                            Edit Profile
                        </button>
                    )}
                </div>
            </div>

            {/* Success Message */}
            {showSuccess && (
                <div className="mb-6 bg-green-50 border border-green-200 rounded-2xl p-4">
                    <div className="flex items-center gap-3">
                        <Save className="w-5 h-5 text-green-600" />
                        <p className="text-green-800 font-medium">Profile updated successfully!</p>
                    </div>
                </div>
            )}

            {/* Privacy Notice */}
            <div className="bg-gradient-to-r from-purple-50 to-indigo-50 border border-purple-200 rounded-2xl p-6 mb-8 shadow-sm">
                <div className="flex items-start gap-4">
                    <div className="flex-shrink-0">
                        <div className="w-10 h-10 bg-purple-200 rounded-full flex items-center justify-center">
                            <Shield className="w-5 h-5 text-purple-700" />
                        </div>
                    </div>
                    <div className="flex-1">
                        <h3 className="text-lg font-semibold text-purple-900 mb-2">
                            Privacy & Security
                        </h3>
                        <p className="text-purple-800 leading-relaxed">
                            This information is associated with your Crabmailer profile and can be used to access multiple Crabmailer accounts. 
                            <span className="font-semibold block mt-1">All contact information is kept strictly confidential.</span>
                        </p>
                    </div>
                </div>
            </div>

            {/* Profile Form */}
            <div className="bg-white border border-gray-200 rounded-2xl shadow-sm overflow-hidden">
                {/* Form Header */}
                <div className="bg-gray-50 px-8 py-4 border-b border-gray-200">
                    <div className="flex items-center justify-between">
                        <h2 className="text-xl font-semibold text-gray-900">Personal Information</h2>
                        {isEditing && (
                            <div className="flex items-center gap-3">
                                <button
                                    onClick={handleCancel}
                                    className="flex items-center gap-2 px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-700 font-medium rounded-lg transition-colors"
                                >
                                    <X className="w-4 h-4" />
                                    Cancel
                                </button>
                                <button
                                    onClick={handleSave}
                                    disabled={isLoading}
                                    className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-700 disabled:bg-purple-400 text-white font-medium rounded-lg transition-colors disabled:cursor-not-allowed"
                                >
                                    {isLoading ? (
                                        <>
                                            <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                                            Saving...
                                        </>
                                    ) : (
                                        <>
                                            <Save className="w-4 h-4" />
                                            Save Changes
                                        </>
                                    )}
                                </button>
                            </div>
                        )}
                    </div>
                </div>

                {/* Form Body */}
                <div className="p-8">
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                        {/* Left Column */}
                        <div className="space-y-6">
                            {/* Full Name */}
                            <div>
                                <label className="block text-gray-700 font-medium mb-2 flex items-center gap-2">
                                    <User className="w-4 h-4" />
                                    Full Name
                                </label>
                                <input
                                    type="text"
                                    value={formData.fullname}
                                    onChange={(e) => handleInputChange('fullname', e.target.value)}
                                    disabled={!isEditing}
                                    className={`w-full h-12 px-4 rounded-xl border transition-colors ${
                                        isEditing 
                                            ? 'border-gray-300 hover:border-purple-400 focus:border-purple-500 focus:outline-none' 
                                            : 'border-gray-200 bg-gray-50 text-gray-600'
                                    }`}
                                    placeholder="Enter your full name"
                                />
                                {errors.fullname && (
                                    <p className="mt-1 text-red-600 text-sm">{errors.fullname}</p>
                                )}
                            </div>

                            {/* Email */}
                            <div>
                                <label className="block text-gray-700 font-medium mb-2 flex items-center gap-2">
                                    <Mail className="w-4 h-4" />
                                    Email Address
                                </label>
                                <input
                                    type="email"
                                    value={formData.email}
                                    disabled
                                    className="w-full h-12 px-4 rounded-xl border border-gray-200 bg-gray-50 text-gray-600"
                                />
                                <p className="mt-1 text-gray-500 text-sm">Email cannot be changed</p>
                            </div>
                        </div>

                        {/* Right Column */}
                        <div className="space-y-6">
                            {/* Company */}
                            <div>
                                <label className="block text-gray-700 font-medium mb-2 flex items-center gap-2">
                                    <Building className="w-4 h-4" />
                                    Company
                                </label>
                                <input
                                    type="text"
                                    value={formData.company}
                                    onChange={(e) => handleInputChange('company', e.target.value)}
                                    disabled={!isEditing}
                                    className={`w-full h-12 px-4 rounded-xl border transition-colors ${
                                        isEditing 
                                            ? 'border-gray-300 hover:border-purple-400 focus:border-purple-500 focus:outline-none' 
                                            : 'border-gray-200 bg-gray-50 text-gray-600'
                                    }`}
                                    placeholder="Enter your company name"
                                />
                            </div>

                            {/* Phone Number */}
                            <div>
                                <label className="block text-gray-700 font-medium mb-2 flex items-center gap-2">
                                    <Phone className="w-4 h-4" />
                                    Phone Number
                                </label>
                                <input
                                    type="tel"
                                    value={formData.phonenumber}
                                    onChange={(e) => handleInputChange('phonenumber', e.target.value)}
                                    disabled={!isEditing}
                                    maxLength={11}
                                    className={`w-full h-12 px-4 rounded-xl border transition-colors ${
                                        isEditing 
                                            ? 'border-gray-300 hover:border-purple-400 focus:border-purple-500 focus:outline-none' 
                                            : 'border-gray-200 bg-gray-50 text-gray-600'
                                    }`}
                                    placeholder="Enter your phone number"
                                />
                                {errors.phonenumber && (
                                    <p className="mt-1 text-red-600 text-sm">{errors.phonenumber}</p>
                                )}
                                {isEditing && (
                                    <p className="mt-1 text-gray-500 text-sm">Must be exactly 11 digits</p>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Submit Error */}
                    {errors.submit && (
                        <div className="mt-6 bg-red-50 border border-red-200 rounded-xl p-4">
                            <p className="text-red-800 text-sm">{errors.submit}</p>
                        </div>
                    )}
                </div>
            </div>

            {/* Account Stats */}
            <div className="mt-8 grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-white border border-gray-200 rounded-xl p-6 text-center">
                    <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-3">
                        <User className="w-6 h-6 text-blue-600" />
                    </div>
                    <h3 className="font-semibold text-gray-900 mb-1">Account Status</h3>
                    <p className="text-green-600 font-medium">Active</p>
                </div>
                
                <div className="bg-white border border-gray-200 rounded-xl p-6 text-center">
                    <div className="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center mx-auto mb-3">
                        <Shield className="w-6 h-6 text-purple-600" />
                    </div>
                    <h3 className="font-semibold text-gray-900 mb-1">Security Level</h3>
                    <p className="text-purple-600 font-medium">Standard</p>
                </div>
                
                <div className="bg-white border border-gray-200 rounded-xl p-6 text-center">
                    <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-3">
                        <Mail className="w-6 h-6 text-green-600" />
                    </div>
                    <h3 className="font-semibold text-gray-900 mb-1">Email Verified</h3>
                    <p className="text-green-600 font-medium">Verified</p>
                </div>
            </div>
        </div>
    );
};

export default ProfileInformationComponent;