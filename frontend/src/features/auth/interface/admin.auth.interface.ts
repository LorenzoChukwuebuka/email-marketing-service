export interface AdminLoginValues {
    email: string;
    password: string;
}

export type AdminData = {
    firstname: string;
    middlename?: string
    lastname:string
    email: string;
    type: "admin";  // Literal type for "admin"
}
