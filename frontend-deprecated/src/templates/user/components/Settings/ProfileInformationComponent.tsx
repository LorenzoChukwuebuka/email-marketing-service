import React, { useEffect, useState, ChangeEvent, MouseEvent } from "react";
import * as Yup from "yup";
import useAuthStore from "../../../../store/userstore/AuthStore";

const ProfileInformationComponent: React.FC = () => {
    const {
        userData,
        getUserDetails,
        setEditFormValues,
        editFormValues,
        editUserDetails,
    } = useAuthStore();

    const [errors, setErrors] = useState<{ [key: string]: string }>({});

    const initEdit = () => {
        setEditFormValues({
            fullname: userData?.fullname || "",
            email: userData?.email || "",
            company: userData?.company || "",
            phonenumber: userData?.phonenumber || "",
        });
    };

    // Define the validation schema
    const validationSchema = Yup.object().shape({
        phonenumber: Yup.string()
            .required("Phone number is required")
            .min(11, "Phone number must be exactly 11 digits")
            .max(11, "Phone number must be exactly 11 digits"),
    });

    const handleEditInformation = async (e: MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();

        try {
            await validationSchema.validate(
                { phonenumber: editFormValues.phonenumber },
                { abortEarly: false }
            );
            await editUserDetails();
            initEdit();
        } catch (error) {
            if (error instanceof Yup.ValidationError) {
                const validationErrors: { [key: string]: string } = {};
                error.inner.forEach((err) => {
                    validationErrors[err.path || ""] = err.message;
                });
                setErrors(validationErrors);
            }
        }
    };

    useEffect(() => {
        getUserDetails();
    }, [getUserDetails]);

    useEffect(() => {
        if (userData) {
            initEdit();
        }
    }, [userData]);

    return (
        <div className="max-w-3xl ml-5 p-6 ">
            <h1 className="text-2xl font-bold text-gray-800 mb-4">
                Profile Information
            </h1>
            <p className="text-gray-600 mb-6">
                This is the information we have associated with your Crabmailer profile,
                which you can use to access multiple Crabmailer accounts.
                <strong className="block">
                    All contact information is kept strictly confidential.
                </strong>
            </p>
            <div className="space-y-4">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-gray-700 font-bold mb-1">
                            Username
                        </label>
                        <input
                            type="text"
                            name="username"
                            value={editFormValues?.email || ""}
                            readOnly
                            className="w-full bg-gray-100 p-2 rounded-md"
                        />
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-gray-700 font-bold mb-1">
                                FullName
                            </label>
                            <input
                                type="text"
                                name="fullName"
                                value={editFormValues?.fullname || ""}
                                onChange={(event: ChangeEvent<HTMLInputElement>) =>
                                    setEditFormValues({
                                        ...editFormValues,
                                        fullname: event.target.value,
                                    })
                                }
                                className="w-full bg-gray-100 p-2 rounded-md"
                            />
                        </div>
                        <div>
                            <label className="block text-gray-700 font-bold mb-1">
                                Company
                            </label>
                            <input
                                type="text"
                                name="company"
                                value={editFormValues?.company || ""}
                                onChange={(event: ChangeEvent<HTMLInputElement>) =>
                                    setEditFormValues({
                                        ...editFormValues,
                                        company: event.target.value,
                                    })
                                }
                                className="w-full bg-gray-100 p-2 rounded-md"
                            />
                        </div>
                    </div>
                </div>
                <div>
                    <label className="block text-gray-700 font-bold mb-1">Email</label>
                    <input
                        type="email"
                        name="email"
                        value={editFormValues?.email || ""}
                        readOnly
                        className="w-[20em] bg-gray-100 p-2 rounded-md"
                    />
                </div>
                <div>
                    <label className="block text-gray-700 font-bold mb-1">
                        Phone Number
                    </label>
                    <input
                        type="text"
                        name="phoneNumber"
                        value={editFormValues?.phonenumber || ""}
                        maxLength={11}
                        onChange={(event: ChangeEvent<HTMLInputElement>) =>
                            setEditFormValues({
                                ...editFormValues,
                                phonenumber: event.target.value,
                            })
                        }
                        className="w-[20em] bg-gray-100 p-2 rounded-md"
                    />
                    {errors.phonenumber && (
                        <div className="text-red-500">{errors.phonenumber}</div>
                    )}
                </div>
            </div>
            <button
                onClick={handleEditInformation}
                className="mt-6 bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold py-2 px-4 rounded-md"
            >
                Edit Information
            </button>
        </div>
    );
};

export default ProfileInformationComponent;
