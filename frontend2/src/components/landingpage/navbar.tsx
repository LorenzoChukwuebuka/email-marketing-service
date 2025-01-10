import React, { useState, useEffect } from 'react';
import { Link } from "react-router-dom";
import { Menu, X } from 'lucide-react';
import Cookies from "js-cookie";
import renderApiName from '../../utils/render-name';

const NavBar = () => {
    const [isMenuOpen, setIsMenuOpen] = useState(false);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isVisible, setIsVisible] = useState(true);
    const [lastScrollY, setLastScrollY] = useState(0);

    useEffect(() => {
        const userCookie = Cookies.get("Cookies");
        if (userCookie) {
            setIsAuthenticated(true);
        } else {
            setIsAuthenticated(false);
        }
    }, []);

    useEffect(() => {
        const handleScroll = () => {
            const currentScrollY = window.scrollY;
            setIsVisible(lastScrollY > currentScrollY || currentScrollY < 10);
            setLastScrollY(currentScrollY);
        };

        window.addEventListener('scroll', handleScroll, { passive: true });
        return () => window.removeEventListener('scroll', handleScroll);
    }, [lastScrollY]);

    return (
        <header className={`bg-white shadow-sm sticky top-0 z-50 transition-all duration-300 ease-in-out ${isVisible ? 'translate-y-0' : '-translate-y-full'}`}>
            <div className="container mx-auto px-4 flex justify-between items-center h-16">
                <h1 className="text-center text-2xl font-bold">
                    <Link to="/">
                        {renderApiName()}
                    </Link>
                </h1>
                <nav className="hidden md:flex space-x-6">
                    {['features', 'pricing', 'faq'].map((item) => (
                        <a
                            key={item}
                            href={`#${item}`}
                            className="text-gray-600 hover:text-gray-800 relative group"
                        >
                            <span className="capitalize">{item}</span>
                            <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-indigo-700 transition-all duration-300 group-hover:w-full"></span>
                        </a>
                    ))}
                </nav>
                <div className="hidden md:flex items-center space-x-4 transition-all duration-300">
                    {!isAuthenticated ? (
                        <>
                            <Link
                                to="/auth/login"
                                className="text-gray-600 hover:text-gray-800 px-3 py-2 relative group"
                            >
                                <span>Login</span>
                                <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-indigo-700 transition-all duration-300 group-hover:w-full"></span>
                            </Link>
                            <Link
                                to="/auth/sign-up"
                                className="bg-blue-900 hover:bg-blue-700 text-white px-4 py-2 rounded-md transition-colors duration-300"
                            >
                                Sign up
                            </Link>
                        </>
                    ) : (
                        <Link
                            to="/app"
                            className="bg-blue-900 hover:bg-blue-700 text-white px-4 py-2 rounded-md transition-colors duration-300"
                        >
                            My Dashboard
                        </Link>
                    )}
                </div>
                <button
                    className="md:hidden transition-transform duration-200 ease-in-out hover:scale-110"
                    onClick={() => setIsMenuOpen(!isMenuOpen)}
                >
                    {isMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
                </button>
            </div>
            <div
                className={`md:hidden bg-white overflow-hidden transition-all duration-300 ease-in-out ${isMenuOpen ? 'max-h-64 opacity-100' : 'max-h-0 opacity-0'
                    }`}
            >
                <div className="px-4 py-2 space-y-2">
                    {['features', 'pricing', 'faq'].map((item) => (
                        <a
                            key={item}
                            href={`#${item}`}
                            className="block py-2 text-gray-600 hover:text-gray-800 transition-colors duration-200"
                        >
                            <span className="capitalize">{item}</span>
                        </a>
                    ))}
                    <Link
                        to="/auth/login"
                        className="block py-2 text-gray-600 hover:text-gray-800 transition-colors duration-200"
                    >
                        Login
                    </Link>
                    <Link
                        to="/auth/sign-up"
                        className="block py-2 text-blue-600 hover:text-blue-700 transition-colors duration-200"
                    >
                        Sign up
                    </Link>
                </div>
            </div>
        </header>
    );
};

export default NavBar;