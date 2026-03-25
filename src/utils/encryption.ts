import crypto from 'crypto';
import { log } from './logger';

const ALGORITHM = 'aes-256-gcm';
const IV_LENGTH = 12;
const AUTH_TAG_LENGTH = 16;

function getEncryptionKey(): Buffer | null {
  const keyEnv = process.env.CREDENTIALS_ENCRYPTION_KEY;
  if (!keyEnv) return null;

  // Support hex (64 chars) or base64 (44 chars) encoded 32-byte keys
  if (keyEnv.length === 64 && /^[0-9a-fA-F]+$/.test(keyEnv)) {
    return Buffer.from(keyEnv, 'hex');
  }
  const decoded = Buffer.from(keyEnv, 'base64');
  if (decoded.length === 32) {
    return decoded;
  }

  log.error('CREDENTIALS_ENCRYPTION_KEY must be a 32-byte key encoded as hex (64 chars) or base64 (44 chars)');
  return null;
}

export interface EncryptedData {
  ciphertext: string;
  iv: string;
  authTag: string;
}

/**
 * Encrypt plaintext using AES-256-GCM
 */
export function encrypt(plaintext: string): EncryptedData {
  const key = getEncryptionKey();
  if (!key) {
    throw new Error('CREDENTIALS_ENCRYPTION_KEY environment variable is not set or invalid');
  }

  const iv = crypto.randomBytes(IV_LENGTH);
  const cipher = crypto.createCipheriv(ALGORITHM, key, iv, { authTagLength: AUTH_TAG_LENGTH });

  let encrypted = cipher.update(plaintext, 'utf8', 'hex');
  encrypted += cipher.final('hex');
  const authTag = cipher.getAuthTag();

  return {
    ciphertext: encrypted,
    iv: iv.toString('hex'),
    authTag: authTag.toString('hex')
  };
}

/**
 * Decrypt ciphertext using AES-256-GCM
 */
export function decrypt(ciphertext: string, iv: string, authTag: string): string {
  const key = getEncryptionKey();
  if (!key) {
    throw new Error('CREDENTIALS_ENCRYPTION_KEY environment variable is not set or invalid');
  }

  const decipher = crypto.createDecipheriv(ALGORITHM, key, Buffer.from(iv, 'hex'), { authTagLength: AUTH_TAG_LENGTH });
  decipher.setAuthTag(Buffer.from(authTag, 'hex'));

  let decrypted = decipher.update(ciphertext, 'hex', 'utf8');
  decrypted += decipher.final('utf8');

  return decrypted;
}

/**
 * Check if encryption is available (key is configured)
 */
export function isEncryptionAvailable(): boolean {
  return getEncryptionKey() !== null;
}
