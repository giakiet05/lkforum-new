import type { MessageResponse } from "../dtos/message-dto";

const ALGORITHM = "AES-GCM";
const KEY_VERSION = "demo-v1";
const DECRYPTION_FAILED_CONTENT = "[Encrypted message could not be decrypted]";

function bytesToBase64(bytes: Uint8Array): string {
    let binary = "";
    bytes.forEach((byte) => {
        binary += String.fromCharCode(byte);
    });
    return btoa(binary);
}

function base64ToBytes(value: string): Uint8Array {
    const binary = atob(value);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
        bytes[i] = binary.charCodeAt(i);
    }
    return bytes;
}

async function getChannelKey(channelId: string): Promise<CryptoKey> {
    const material = new TextEncoder().encode(`lkforum-e2ee-demo:${channelId}`);
    const digest = await crypto.subtle.digest("SHA-256", material);
    return crypto.subtle.importKey("raw", digest, ALGORITHM, false, ["encrypt", "decrypt"]);
}

export async function encryptMessage(channelId: string, plaintext: string) {
    const key = await getChannelKey(channelId);
    const nonce = crypto.getRandomValues(new Uint8Array(12));
    const encrypted = await crypto.subtle.encrypt(
        { name: ALGORITHM, iv: nonce },
        key,
        new TextEncoder().encode(plaintext),
    );

    return {
        ciphertext: bytesToBase64(new Uint8Array(encrypted)),
        nonce: bytesToBase64(nonce),
        algorithm: ALGORITHM,
        key_version: KEY_VERSION,
    };
}

export function isEncryptedMessage(message: Pick<MessageResponse, "ciphertext" | "nonce">): boolean {
    return Boolean(message.ciphertext && message.nonce);
}

export async function decryptMessage(message: MessageResponse): Promise<MessageResponse> {
    if (!isEncryptedMessage(message)) {
        return message;
    }

    try {
        const key = await getChannelKey(message.channel_id);
        const decrypted = await crypto.subtle.decrypt(
            { name: message.algorithm || ALGORITHM, iv: base64ToBytes(message.nonce) },
            key,
            base64ToBytes(message.ciphertext),
        );

        return {
            ...message,
            content: new TextDecoder().decode(decrypted),
        };
    } catch {
        return {
            ...message,
            content: DECRYPTION_FAILED_CONTENT,
        };
    }
}

export async function decryptMessages(messages: MessageResponse[]): Promise<MessageResponse[]> {
    return Promise.all(messages.map(decryptMessage));
}
