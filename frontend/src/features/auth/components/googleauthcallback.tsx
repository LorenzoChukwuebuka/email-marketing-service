// src/components/auth/GoogleCallback.tsx
import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
//import useAuthStore from '../store/auth.store';

const GoogleCallback = () => {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
   // const { setLoginValues, loginUser } = useAuthStore();

    useEffect(() => {
        console.log('you are here')
        const token = searchParams.get('token');

        console.log(token)
        // if (token) {
        //     // Use your existing auth store to handle the login
        //     // setLoginValues({ token });
        //     // loginUser();
        //     navigate('/app'); // or wherever you redirect after login
        // } else {
            navigate('/'); // redirect back to login if no token
        // }
    }, []);

    return <div>Completing login...</div>;
};

export default GoogleCallback;