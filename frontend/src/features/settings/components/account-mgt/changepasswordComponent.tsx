import React, { useState } from "react";
import { Lock, Key, Shield, CheckCircle } from 'lucide-react';
import useAuthStore from "../../../auth/store/auth.store";

const ChangePasswordComponent: React.FC = () => {
    const {
        changePassword,
        setChangePasswordValues,
    } = useAuthStore();

    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [passwordStrength, setPasswordStrength] = useState<number>(0);
    const [formData, setFormData] = useState({
        old_password: '',
        new_password: '',
        confirm_password: ''
    });
    const [errors, setErrors] = useState<Record<string, string>>({});
    const [showSuccess, setShowSuccess] = useState(false);

    const calculatePasswordStrength = (password: string): number => {
        let strength = 0;
        if (password.length >= 8) strength += 25;
        if (/[A-Z]/.test(password)) strength += 25;
        if (/[a-z]/.test(password)) strength += 25;
        if (/[0-9]/.test(password)) strength += 12.5;
        if (/[^A-Za-z0-9]/.test(password)) strength += 12.5;
        return Math.min(strength, 100);
    };

    const getStrengthColor = (strength: number): string => {
        if (strength < 30) return 'bg-red-500';
        if (strength < 60) return 'bg-yellow-500';
        if (strength < 80) return 'bg-blue-500';
        return 'bg-green-500';
    };

    const getStrengthText = (strength: number): string => {
        if (strength < 30) return 'Weak';
        if (strength < 60) return 'Fair';
        if (strength < 80) return 'Good';
        return 'Strong';
    };

    const handleInputChange = (field: string, value: string) => {
        setFormData(prev => ({ ...prev, [field]: value }));
        
        if (field === 'new_password') {
            const strength = calculatePasswordStrength(value);
            setPasswordStrength(strength);
        }
        
        // Clear errors when user starts typing
        if (errors[field]) {
            setErrors(prev => ({ ...prev, [field]: '' }));
        }
    };

    const validateForm = (): boolean => {
        const newErrors: Record<string, string> = {};

        if (!formData.old_password) {
            newErrors.old_password = 'Current password is required';
        }

        if (!formData.new_password) {
            newErrors.new_password = 'New password is required';
        } else if (formData.new_password.length < 8) {
            newErrors.new_password = 'Password must be at least 8 characters long';
        }

        if (!formData.confirm_password) {
            newErrors.confirm_password = 'Please confirm your password';
        } else if (formData.new_password !== formData.confirm_password) {
            newErrors.confirm_password = 'Passwords must match';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        
        if (!validateForm()) return;

        try {
            setIsLoading(true);
            setChangePasswordValues(formData);
            await changePassword();
            
            // Reset form and show success
            setFormData({ old_password: '', new_password: '', confirm_password: '' });
            setPasswordStrength(0);
            setShowSuccess(true);
            setTimeout(() => setShowSuccess(false), 5000);
        } catch (error) {
            console.error("Failed to change password:", error);
            setErrors({ submit: 'Failed to change password. Please try again.' });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="max-w-4xl mx-auto p-8">
            {/* Header Section */}
            <div className="mb-8">
                <div className="flex items-center gap-3 mb-4">
                    <div className="p-2 bg-blue-100 rounded-lg">
                        <Lock className="w-6 h-6 text-blue-600" />
                    </div>
                    <div>
                        <h1 className="text-3xl font-bold text-gray-900">Change Password</h1>
                        <p className="text-gray-500 mt-1">Update your account security</p>
                    </div>
                </div>
            </div>

            {/* Success Message */}
            {showSuccess && (
                <div className="mb-8 bg-green-50 border border-green-200 rounded-2xl p-4">
                    <div className="flex items-center gap-3">
                        <CheckCircle className="w-5 h-5 text-green-600" />
                        <p className="text-green-800 font-medium">Password changed successfully!</p>
                    </div>
                </div>
            )}

            {/* Security Tips Card */}
            <div className="bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded-2xl p-6 mb-8 shadow-sm">
                <div className="flex items-start gap-4">
                    <div className="flex-shrink-0">
                        <div className="w-10 h-10 bg-blue-200 rounded-full flex items-center justify-center">
                            <Shield className="w-5 h-5 text-blue-700" />
                        </div>
                    </div>
                    <div className="flex-1">
                        <h3 className="text-lg font-semibold text-blue-900 mb-2">
                            Password Security Tips
                        </h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-blue-800 text-sm">
                            <div className="flex items-center gap-2">
                                <div className="w-1.5 h-1.5 bg-blue-600 rounded-full"></div>
                                Use at least 8 characters
                            </div>
                            <div className="flex items-center gap-2">
                                <div className="w-1.5 h-1.5 bg-blue-600 rounded-full"></div>
                                Include uppercase and lowercase letters
                            </div>
                            <div className="flex items-center gap-2">
                                <div className="w-1.5 h-1.5 bg-blue-600 rounded-full"></div>
                                Add numbers and special characters
                            </div>
                            <div className="flex items-center gap-2">
                                <div className="w-1.5 h-1.5 bg-blue-600 rounded-full"></div>
                                Avoid personal information
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Form Card */}
            <div className="bg-white border border-gray-200 rounded-2xl p-8 shadow-sm">
                <div className="max-w-md">
                    <div className="space-y-6">
                        {/* Current Password */}
                        <div>
                            <label className="block text-gray-700 font-medium mb-2 flex items-center gap-2">
                                <Key className="w-4 h-4" />
                                Current Password
                            </label>
                            <input
                                type="password"
                                value={formData.old_password}
                                onChange={(e) => handleInputChange('old_password', e.target.value)}
                                className="w-full h-12 px-4 rounded-xl border border-gray-300 hover:border-blue-400 focus:border-blue-500 focus:outline-none transition-colors"
                                placeholder="Enter your current password"
                            />
                            {errors.old_password && (
                                <p className="mt-1 text-red-600 text-sm">{errors.old_password}</p>
                            )}
                        </div>

                        {/* New Password */}
                        <div>
                            <label className="block text-gray-700 font-medium mb-2 flex items-center gap-2">
                                <Lock className="w-4 h-4" />
                                New Password
                            </label>
                            <input
                                type="password"
                                value={formData.new_password}
                                onChange={(e) => handleInputChange('new_password', e.target.value)}
                                className="w-full h-12 px-4 rounded-xl border border-gray-300 hover:border-blue-400 focus:border-blue-500 focus:outline-none transition-colors"
                                placeholder="Enter your new password"
                            />
                            
                            {/* Password Strength Indicator */}
                            {passwordStrength > 0 && (
                                <div className="mt-3">
                                    <div className="flex items-center justify-between mb-1">
                                        <span className="text-xs text-gray-600">Password Strength</span>
                                        <span className={`text-xs font-medium ${
                                            passwordStrength < 30 ? 'text-red-600' :
                                            passwordStrength < 60 ? 'text-yellow-600' :
                                            passwordStrength < 80 ? 'text-blue-600' : 'text-green-600'
                                        }`}>
                                            {getStrengthText(passwordStrength)}
                                        </span>
                                    </div>
                                    <div className="w-full bg-gray-200 rounded-full h-2">
                                        <div
                                            className={`h-2 rounded-full transition-all duration-300 ${getStrengthColor(passwordStrength)}`}
                                            style={{ width: `${passwordStrength}%` }}
                                        ></div>
                                    </div>
                                </div>
                            )}
                            
                            {errors.new_password && (
                                <p className="mt-1 text-red-600 text-sm">{errors.new_password}</p>
                            )}
                        </div>

                        {/* Confirm Password */}
                        <div>
                            <label className="block text-gray-700 font-medium mb-2 flex items-center gap-2">
                                <CheckCircle className="w-4 h-4" />
                                Confirm New Password
                            </label>
                            <input
                                type="password"
                                value={formData.confirm_password}
                                onChange={(e) => handleInputChange('confirm_password', e.target.value)}
                                className="w-full h-12 px-4 rounded-xl border border-gray-300 hover:border-blue-400 focus:border-blue-500 focus:outline-none transition-colors"
                                placeholder="Confirm your new password"
                            />
                            {errors.confirm_password && (
                                <p className="mt-1 text-red-600 text-sm">{errors.confirm_password}</p>
                            )}
                        </div>

                        {/* Submit Error */}
                        {errors.submit && (
                            <div className="bg-red-50 border border-red-200 rounded-xl p-3">
                                <p className="text-red-800 text-sm">{errors.submit}</p>
                            </div>
                        )}

                        {/* Submit Button */}
                        <button
                            type="button"
                            onClick={handleSubmit}
                            disabled={isLoading}
                            className="w-full h-12 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl transition-all duration-200 flex items-center justify-center gap-2 disabled:cursor-not-allowed"
                        >
                            {isLoading ? (
                                <>
                                    <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                                    Updating Password...
                                </>
                            ) : (
                                <>
                                    <Lock className="w-5 h-5" />
                                    Change Password
                                </>
                            )}
                        </button>
                    </div>
                </div>
            </div>

            {/* Additional Security Section */}
            <div className="mt-8 text-center">
                <p className="text-gray-500 text-sm mb-2">
                    Want to enhance your account security even more?
                </p>
                <a href="#" className="text-blue-600 hover:text-blue-700 font-medium text-sm hover:underline">
                    Set up two-factor authentication
                </a>
            </div>
        </div>
    );
};

export default ChangePasswordComponent;