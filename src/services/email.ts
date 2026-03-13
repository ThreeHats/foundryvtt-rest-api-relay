import nodemailer from 'nodemailer';
import Handlebars from 'handlebars';
import { log } from '../utils/logger';

const SMTP_HOST = process.env.SMTP_HOST;
const SMTP_PORT = parseInt(process.env.SMTP_PORT || '587');
const SMTP_USER = process.env.SMTP_USER;
const SMTP_PASS = process.env.SMTP_PASS;
const SMTP_FROM = process.env.SMTP_FROM || 'noreply@foundryvtt-relay.com';
const SMTP_SECURE = process.env.SMTP_SECURE === 'true';
const FRONTEND_URL = process.env.FRONTEND_URL || 'http://localhost:3010';

let transporter: nodemailer.Transporter | null = null;

if (SMTP_HOST) {
  transporter = nodemailer.createTransport({
    host: SMTP_HOST,
    port: SMTP_PORT,
    secure: SMTP_SECURE,
    auth: (SMTP_USER && SMTP_PASS) ? {
      user: SMTP_USER,
      pass: SMTP_PASS
    } : undefined
  });
}

const resetEmailTemplate = Handlebars.compile(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body { font-family: Arial, sans-serif; background: #f4f4f4; padding: 20px; }
    .container { max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; padding: 30px; }
    .btn { display: inline-block; background: #4a6fa5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; margin: 20px 0; }
    .footer { margin-top: 30px; font-size: 12px; color: #999; }
  </style>
</head>
<body>
  <div class="container">
    <h2>Password Reset Request</h2>
    <p>You requested a password reset for your Foundry REST API account.</p>
    <p>Click the button below to reset your password. This link expires in 1 hour.</p>
    <a href="{{resetUrl}}" class="btn">Reset Password</a>
    <p>If the button doesn't work, copy and paste this URL into your browser:</p>
    <p style="word-break: break-all; font-size: 14px; color: #666;">{{resetUrl}}</p>
    <div class="footer">
      <p>If you didn't request this, you can safely ignore this email.</p>
      <p>Foundry REST API Relay</p>
    </div>
  </div>
</body>
</html>
`);

export async function sendPasswordResetEmail(to: string, resetToken: string): Promise<void> {
  const resetUrl = `${FRONTEND_URL}?reset-token=${resetToken}`;

  if (!transporter) {
    log.info('SMTP not configured — password reset URL (for development):');
    log.info(resetUrl);
    return;
  }

  try {
    const html = resetEmailTemplate({ resetUrl });

    await transporter.sendMail({
      from: SMTP_FROM,
      to,
      subject: 'Password Reset — Foundry REST API',
      html
    });

    log.info('Password reset email sent', { to });
  } catch (error) {
    log.error('Failed to send password reset email', { error, to });
  }
}
