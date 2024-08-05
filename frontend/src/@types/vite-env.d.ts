/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_API_URL: string
    readonly VITE_API_NAME: string
    readonly VITE_ENC_KEY: string
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}