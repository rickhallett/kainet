// Crypto functions for AES-256-GCM encryption

export async function deriveKey(authToken) {
  const encoder = new TextEncoder();
  const data = encoder.encode(authToken);
  const hashBuffer = await crypto.subtle.digest('SHA-256', data);

  return crypto.subtle.importKey(
    'raw',
    hashBuffer,
    { name: 'AES-GCM' },
    false,
    ['encrypt', 'decrypt']
  );
}

export async function encrypt(key, plaintext) {
  const encoder = new TextEncoder();
  const data = encoder.encode(plaintext);

  // Generate random nonce (12 bytes for GCM)
  const nonce = crypto.getRandomValues(new Uint8Array(12));

  // Encrypt
  const ciphertext = await crypto.subtle.encrypt(
    {
      name: 'AES-GCM',
      iv: nonce
    },
    key,
    data
  );

  // Combine nonce + ciphertext
  const combined = new Uint8Array(nonce.length + ciphertext.byteLength);
  combined.set(nonce, 0);
  combined.set(new Uint8Array(ciphertext), nonce.length);

  // Convert to base64
  return btoa(String.fromCharCode(...combined));
}

export async function decrypt(key, ciphertext) {
  // Decode from base64
  const combined = Uint8Array.from(atob(ciphertext), c => c.charCodeAt(0));

  // Extract nonce and ciphertext
  const nonce = combined.slice(0, 12);
  const data = combined.slice(12);

  // Decrypt
  const plaintext = await crypto.subtle.decrypt(
    {
      name: 'AES-GCM',
      iv: nonce
    },
    key,
    data
  );

  const decoder = new TextDecoder();
  return decoder.decode(plaintext);
}
