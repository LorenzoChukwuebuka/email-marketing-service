import CryptoJS from "crypto-js";

const ENCRYPTION_KEY: string = import.meta.env.VITE_ENC_KEY;

const decryptApiKey = (encryptedKey: string): string => {
    try {
        console.log("Encrypted key:", encryptedKey);

        // Decode the base64 URL encoded string
        const ciphertext: CryptoJS.lib.WordArray = CryptoJS.enc.Base64url.parse(encryptedKey);
        console.log("Decoded ciphertext:", ciphertext.toString());

        // Convert the encryption key to WordArray
        const key: CryptoJS.lib.WordArray = CryptoJS.enc.Utf8.parse(ENCRYPTION_KEY);
        console.log("Encryption key:", ENCRYPTION_KEY);
        console.log("Key length:", new TextEncoder().encode(ENCRYPTION_KEY).length);

        // Extract the IV (first 16 bytes)
        const iv: CryptoJS.lib.WordArray = ciphertext.clone();
        iv.sigBytes = 16;
        iv.clamp();
        console.log("IV:", iv.toString());

        // Get the actual ciphertext
        const encryptedBytes: CryptoJS.lib.WordArray = ciphertext.clone();
        encryptedBytes.words.splice(0, 4); // remove IV
        encryptedBytes.sigBytes -= 16;
        console.log("Encrypted bytes:", encryptedBytes.toString());

        // Decrypt using AES in CFB mode
        const decrypted: CryptoJS.lib.WordArray = CryptoJS.AES.decrypt(
            CryptoJS.lib.CipherParams.create({
                ciphertext: encryptedBytes
            }),
            key,
            { iv: iv, mode: CryptoJS.mode.CFB, padding: CryptoJS.pad.NoPadding }
        );

        const result: string = decrypted.toString(CryptoJS.enc.Utf8);
        console.log("Decrypted result:", result);
        return result;
    } catch (error) {
        console.error("Decryption failed:", error);
        return "Decryption failed";
    }
};

export default decryptApiKey;