import CryptoJS from "crypto-js";

const ENCRYPTION_KEY = import.meta.env.VITE_ENC_KEY;

const decryptApiKey = (encryptedKey) => {
    try {
      console.log("Encrypted key:", encryptedKey);
      
      // Decode the base64 URL encoded string
      const ciphertext = CryptoJS.enc.Base64url.parse(encryptedKey);
      console.log("Decoded ciphertext:", ciphertext.toString());
  
      // Convert the encryption key to WordArray
      const key = CryptoJS.enc.Utf8.parse(ENCRYPTION_KEY);
      console.log("Encryption key:", ENCRYPTION_KEY);
      console.log("Key length:", new TextEncoder().encode(ENCRYPTION_KEY).length);
  
      // Extract the IV (first 16 bytes)
      const iv = ciphertext.clone();
      iv.sigBytes = 16;
      iv.clamp();
      console.log("IV:", iv.toString());
  
      // Get the actual ciphertext
      const encryptedBytes = ciphertext.clone();
      encryptedBytes.words.splice(0, 4); // remove IV
      encryptedBytes.sigBytes -= 16;
      console.log("Encrypted bytes:", encryptedBytes.toString());
  
      // Decrypt using AES in CFB mode
      const decrypted = CryptoJS.AES.decrypt(
        { ciphertext: encryptedBytes },
        key,
        { iv: iv, mode: CryptoJS.mode.CFB, padding: CryptoJS.pad.NoPadding }
      );
  
      const result = decrypted.toString(CryptoJS.enc.Utf8);
      console.log("Decrypted result:", result);
      return result;
    } catch (error) {
      console.error("Decryption failed:", error);
      return "Decryption failed";
    }
  };


  export default decryptApiKey