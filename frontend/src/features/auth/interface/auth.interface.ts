import { BaseEntity } from '../../../interface/baseentity.interface';

export type FormValues = {
    fullname: string;
    company: string;
    email: string;
    password: string;
    confirmPassword: string;
    token: string;
    phonenumber: string;
}

export interface Company {
    company_id: string;
    companyname: string;
    company_created_at: string; // ISO 8601 datetime string
    company_updated_at: string;
    company_deleted_at: string | null;
}


export type LoginValues = Pick<FormValues, 'email' | 'password'>

export type ForgetPasswordValues = Pick<FormValues, 'email'>

export type ResetPasswordValues = Pick<FormValues, 'email' | 'confirmPassword' | 'password' | 'token'>

export type EditFormValues = Omit<FormValues, 'password' | 'confirmPassword' | 'token'>;

export type UserDetails = {
    fullname: string
    email: string
    company: string
    phonenumber: string
    verified: boolean,
    blocked: boolean,
} & BaseEntity & Company

export type ChangePasswordValues = {
    old_password: string;
    new_password: string;
    confirm_password: string;
}

export type OtpValue = Pick<FormValues, "token">

export type SignUpAPIData = {
    userId: string,
    message: string
}

export interface VerifyLoginFormData {
    token: string;
    user_id:string
}