// Web Crypto API encryption using AES-256-GCM
// Matches the Go implementation's encryption scheme

const AUTH_TOKEN = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NjMwOTg5MjEsImlkIjoiNjc4YjAxZDQtYzExZi00OTMyLTk2MzktZDUxODNlZTVmMTI3IiwicmlkIjoiNzc0ZDk3NmItNTMzYS00NzE1LTlkZmItY2RmYjU1N2M2ZmRjIn0.4hZBJRmkoSMutMqpGVYhmIWmo2-5lTKle9p5QkSyF6B-OHGS1xpdoFi_wEZVn53qcHEaUEHAzg8rPiy0ZsjFCQ";

// Derive encryption key from auth token using SHA-256
export async function deriveKey(token = AUTH_TOKEN) {
    const encoder = new TextEncoder();
    const data = encoder.encode(token);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);

    return await crypto.subtle.importKey(
        'raw',
        hashBuffer,
        { name: 'AES-GCM', length: 256 },
        false,
        ['encrypt', 'decrypt']
    );
}

// Encrypt plaintext using AES-256-GCM
export async function encrypt(plaintext, key) {
    const encoder = new TextEncoder();
    const data = encoder.encode(plaintext);

    // Generate random 12-byte nonce (same as Go's GCM default)
    const nonce = crypto.getRandomValues(new Uint8Array(12));

    const ciphertext = await crypto.subtle.encrypt(
        {
            name: 'AES-GCM',
            iv: nonce,
        },
        key,
        data
    );

    // Combine nonce + ciphertext (matches Go implementation)
    const combined = new Uint8Array(nonce.length + ciphertext.byteLength);
    combined.set(nonce, 0);
    combined.set(new Uint8Array(ciphertext), nonce.length);

    // Base64 encode (matches Go's base64.StdEncoding)
    return btoa(String.fromCharCode(...combined));
}

// Decrypt ciphertext using AES-256-GCM
export async function decrypt(ciphertextB64, key) {
    try {
        // Decode base64
        const combined = Uint8Array.from(atob(ciphertextB64), c => c.charCodeAt(0));

        // Extract nonce and ciphertext
        const nonceSize = 12;
        if (combined.length < nonceSize) {
            throw new Error('Ciphertext too short');
        }

        const nonce = combined.slice(0, nonceSize);
        const ciphertext = combined.slice(nonceSize);

        const plaintext = await crypto.subtle.decrypt(
            {
                name: 'AES-GCM',
                iv: nonce,
            },
            key,
            ciphertext
        );

        const decoder = new TextDecoder();
        return decoder.decode(plaintext);
    } catch (error) {
        console.error('Decryption failed:', error);
        return null;
    }
}

export { AUTH_TOKEN };
