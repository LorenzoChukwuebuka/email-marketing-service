import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import LoadingSpinnerComponent from '../../../components/loadingSpinnerComponent';
import authApi from '../api/auth.api';
import Cookies from 'js-cookie';

const GoogleCallback = () => {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    // const { setLoginValues, loginUser } = useAuthStore();

    useEffect(() => {
        const token = searchParams.get('token');
        const refresh_token = searchParams.get('refresh_token');
        const details_key = searchParams.get('details_key');

        localStorage.setItem('token', token as string);
        localStorage.setItem('refresh_token', refresh_token as string);
        localStorage.setItem('key', details_key as string)

        // Fetch user details
        const fetchUserDetails = async () => {
            if (!details_key) return;

            try {
                const response = await authApi.getGoogleAuthLoginDetails(details_key);

                const cookieData = {
                    refresh_token,
                    token,
                    details: response.payload, // Ensure the details are nested properly
                };
                // Store in cookies
                Cookies.set('Cookies', JSON.stringify(cookieData), {
                    expires: 3,
                    sameSite: 'Strict',
                    secure: true
                });
                // Delay navigation slightly to ensure storage is cleared before reloading
                setTimeout(() => {
                    localStorage.removeItem('token');
                    localStorage.removeItem('refresh_token');
                    localStorage.removeItem('key');

                    window.location.href = '/app';
                }, 200); // 200ms to ensure storage removal

            } catch (error) {
                console.error('Error fetching user details:', error);
            }
        };

        fetchUserDetails();
    }, [searchParams, navigate]); // Dependencies



    return <LoadingSpinnerComponent />;
};

export default GoogleCallback;