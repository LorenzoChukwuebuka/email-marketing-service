import React, { FormEvent, useRef } from 'react';
import ReCAPTCHA from "react-google-recaptcha";

const RecaptchaComponent = () => {
    const recaptchaRef = useRef(null);

    const handleSubmit = async (event: FormEvent) => {
        event.preventDefault();
        //@ts-ignore
        const token = await recaptchaRef.current.executeAsync();
        // Send token to your backend for verification
        console.log(token);
    };

    return (
        <form onSubmit={handleSubmit}>
            <ReCAPTCHA
                ref={recaptchaRef}
                size="invisible"
                sitekey="YOUR_RECAPTCHA_SITE_KEY"
            />
            <button type="submit">Submit</button>
        </form>
    );
};

export default RecaptchaComponent;