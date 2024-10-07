// src/components/NavBar.tsx

import React, { useState, useEffect } from 'react';
import { Link } from "react-router-dom";
import { Menu, X } from 'lucide-react';
import Cookies from "js-cookie";

const NavBar: React.FC = () => {
    const [isMenuOpen, setIsMenuOpen] = useState(false);
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    useEffect(() => {
        const userCookie = Cookies.get("Cookies"); // Replace "Cookies" with the actual name of your cookie
        if (userCookie) {
            setIsAuthenticated(true); // Set authenticated if cookie exists
        } else {
            setIsAuthenticated(false);
        }
    }, []);

    const apiName = import.meta.env.VITE_API_NAME;
    const firstFourLetters = apiName.slice(0, 4);
    const remainingLetters = apiName.slice(4);

    return (
        <header className="bg-white py-4 shadow-sm sticky top-0 z-50">
            <div className="container mx-auto px-4 flex justify-between items-center">
                <h1 className="text-center text-2xl font-bold">
                    <Link to="/">
                        <span className="text-indigo-700">{firstFourLetters}</span>
                        <span className="text-gray-700">{remainingLetters}</span>
                        <i className="bi bi-mailbox2-flag text-indigo-700 ml-2"></i>
                    </Link>
                </h1>
                <nav className="hidden md:flex space-x-6">
                    <a href="#features" className="text-gray-600 hover:text-gray-800">Features</a>
                    <a href="#pricing" className="text-gray-600 hover:text-gray-800">Pricing</a>
                    <a href="#faq" className="text-gray-600 hover:text-gray-800">FAQ</a>
                </nav>
                <div className="hidden md:flex items-center">
                    {!isAuthenticated ? (
                        <>
                            <Link to="/auth/login" className="text-gray-600 hover:text-gray-800 px-3 py-2">Login</Link>
                            <Link to="/auth/sign-up" className="bg-blue-900 hover:bg-blue-700 text-white px-4 py-2 rounded-md ml-4 transition duration-300">Sign up</Link>
                        </>
                    ) : (
                        <Link to="/user/dash" className="bg-blue-900 hover:bg-blue-700 text-white px-4 py-2 rounded-md ml-4 transition duration-300">My Dashboard</Link>
                    )}
                </div>
                <button className="md:hidden" onClick={() => setIsMenuOpen(!isMenuOpen)}>
                    {isMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
                </button>
            </div>
            {isMenuOpen && (
                <div className="md:hidden bg-white py-4 px-4">
                    <a href="#features" className="block py-2 text-gray-600 hover:text-gray-800">Features</a>
                    <a href="#pricing" className="block py-2 text-gray-600 hover:text-gray-800">Pricing</a>
                    <a href="#faq" className="block py-2 text-gray-600 hover:text-gray-800">FAQ</a>
                    <Link to="/auth/login" className="block py-2 text-gray-600 hover:text-gray-800">Login</Link>
                    <Link to="/auth/sign-up" className="block py-2 text-blue-600 hover:text-blue-700">Sign up</Link>
                </div>
            )}
        </header>
    );
};

export default NavBar;
