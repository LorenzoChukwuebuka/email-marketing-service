import * as Yup from "yup";
import useAuthStore from "../../store/userstore/AuthStore";
import { useState, ChangeEvent, FormEvent } from "react";
import { Link } from "react-router-dom";

const SignUpTemplate: React.FC = () => {
    const [errors, setErrors] = useState<{ [key: string]: string }>({});
    const { formValues, isLoading, setFormValues, registerUser } = useAuthStore();

    const validationSchema = Yup.object().shape({
        fullname: Yup.string()
            .required("Name is required")
            .min(5, "Name must be at least 5 characters"),
        email: Yup.string()
            .email("Invalid email format")
            .required("Email is required"),
        company: Yup.string().required("Company is required"),
        password: Yup.string()
            .required("Password is required")
            .min(8, "Password must be at least 8 characters")
            .matches(/[a-zA-Z]/, "Password must contain at least one letter")
            .matches(/[0-9]/, "Password must contain at least one number"),
        confirmPassword: Yup.string()
            .oneOf([Yup.ref("password"), undefined], "Passwords must match")
            .required("Confirm Password is required"),
    });

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        try {
            await validationSchema.validate(formValues, { abortEarly: false });
            await registerUser();
            setErrors({});
        } catch (err) {
            const validationErrors: { [key: string]: string } = {};
            if (err instanceof Yup.ValidationError) {
                err.inner.forEach((error) => {
                    validationErrors[error.path || ""] = error.message;
                });
                setErrors(validationErrors);
            }
        }
    };

    const apiName = import.meta.env.VITE_API_NAME;
    const firstFourLetters = apiName.slice(0, 4);
    const remainingLetters = apiName.slice(4);

    return (
        <main className="min-h-screen bg-gradient-to-b from-gray-100 to-gray-200">
            <div className="py-4">
                <h1 className="text-center text-3xl font-bold">
                    <span className="text-indigo-700">{firstFourLetters}</span>
                    <span className="text-gray-700">{remainingLetters}</span>
                    <i className="bi bi-mailbox2-flag text-indigo-700 ml-2"></i>
                </h1>
            </div>

            <div className="container mx-auto px-4 py-8">
                <div className="bg-gray-50 shadow-md rounded-lg p-8 max-w-md mx-auto">
                    <h2 className="text-2xl font-semibold text-center text-gray-700 mb-6">
                        Get Started with {import.meta.env.VITE_API_NAME}
                    </h2>

                    <form className="space-y-4" onSubmit={handleSubmit}>
                        {renderInput("fullname", "Full Name", "text")}
                        {renderInput("email", "Email", "email")}
                        {renderInput("company", "Company", "text")}
                        {renderInput("password", "Password", "password")}
                        {renderInput("confirmPassword", "Confirm Password", "password")}

                        <div className="flex flex-col sm:flex-row justify-between items-center space-y-2 sm:space-y-0">
                            <button
                                type="submit"
                                className={`w-full sm:w-auto bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded-md transition duration-300 ease-in-out ${isLoading ? "opacity-50 cursor-not-allowed" : ""
                                    }`}
                                disabled={isLoading}
                            >
                                {isLoading ? (
                                    <span className="flex items-center justify-center">
                                        <svg className="animate-spin h-5 w-5 mr-3" viewBox="0 0 24 24">
                                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                        </svg>
                                        Creating account...
                                    </span>
                                ) : (
                                    "Create Account"
                                )}
                            </button>
                            <Link
                                to="/auth/login"
                                className="text-indigo-600 hover:text-indigo-800 transition duration-300 ease-in-out"
                            >
                                Already have an account? Log in
                            </Link>
                        </div>
                    </form>

                    <div className="mt-4 text-center text-sm text-gray-600">
                        By signing up, you agree to our{" "}
                        <a href="#" className="text-indigo-600 hover:underline">
                            Terms of Service
                        </a>{" "}
                        and{" "}
                        <a href="#" className="text-indigo-600 hover:underline">
                            Privacy Policy
                        </a>
                    </div>
                </div>
            </div>
        </main>
    );

    function renderInput(name: keyof typeof formValues, label: string, type: string) {
        return (
            <div>
                <label htmlFor={name} className="block text-sm font-medium text-gray-600 mb-1">
                    {label}
                </label>
                <input
                    type={type}
                    id={name}
                    name={name}
                    value={formValues[name]}
                    onChange={(e: ChangeEvent<HTMLInputElement>) =>
                        setFormValues({ ...formValues, [name]: e.target.value })
                    }
                    className={`w-full px-3 py-2 border ${errors[name] ? "border-red-400" : "border-gray-300"
                        } rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-400 transition duration-150 ease-in-out`}
                    placeholder={label}
                />
                {errors[name] && <p className="mt-1 text-sm text-red-500">{errors[name]}</p>}
            </div>
        );
    }
};

export default SignUpTemplate;