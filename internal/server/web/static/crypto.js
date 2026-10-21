// Password-based encryption for documents. Everything happens in the browser: the server only stores the
// envelope this produces, never the password or the plaintext.
//
// Envelope (JSON, binary fields in standard base64):
//   { "v": 1, "iter": <PBKDF2 iterations>, "salt": <16 bytes>, "iv": <12 bytes>, "ct": <AES-256-GCM ciphertext + tag> }
(function () {
  const ITERATIONS = 600000;

  const encoder = new TextEncoder();
  const decoder = new TextDecoder();

  // Builds the string in chunks, as spreading a large array into fromCharCode overflows the stack
  function toBase64(bytes) {
    let binary = '';

    for (let i = 0; i < bytes.length; i += 0x8000) {
      binary += String.fromCharCode.apply(null, bytes.subarray(i, i + 0x8000));
    }

    return btoa(binary);
  }

  function fromBase64(text) {
    return Uint8Array.from(atob(text), (c) => c.charCodeAt(0));
  }

  async function deriveKey(password, salt, iterations, usage) {
    const material = await crypto.subtle.importKey('raw', encoder.encode(password), 'PBKDF2', false, ['deriveKey']);

    return crypto.subtle.deriveKey(
      { name: 'PBKDF2', hash: 'SHA-256', salt, iterations },
      material,
      { name: 'AES-GCM', length: 256 },
      false,
      [usage],
    );
  }

  // supported reports whether this browser can encrypt. WebCrypto only exists on HTTPS and localhost.
  function supported() {
    return Boolean(globalThis.crypto && crypto.subtle);
  }

  async function encrypt(text, password) {
    const salt = crypto.getRandomValues(new Uint8Array(16));
    const iv = crypto.getRandomValues(new Uint8Array(12));
    const key = await deriveKey(password, salt, ITERATIONS, 'encrypt');
    const ciphertext = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, key, encoder.encode(text));

    return JSON.stringify({
      v: 1,
      iter: ITERATIONS,
      salt: toBase64(salt),
      iv: toBase64(iv),
      ct: toBase64(new Uint8Array(ciphertext)),
    });
  }

  // decrypt rejects if the password is wrong, as GCM authenticates the ciphertext
  async function decrypt(envelope, password) {
    const parsed = JSON.parse(envelope);

    if (parsed.v !== 1) {
      throw new Error('Unsupported envelope version ' + parsed.v);
    }

    const key = await deriveKey(password, fromBase64(parsed.salt), parsed.iter, 'decrypt');
    const plaintext = await crypto.subtle.decrypt({ name: 'AES-GCM', iv: fromBase64(parsed.iv) }, key, fromBase64(parsed.ct));

    return decoder.decode(plaintext);
  }

  globalThis.spacebinCrypto = { supported, encrypt, decrypt };
})();
