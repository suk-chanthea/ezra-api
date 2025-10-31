# Environment Configuration Guide

This document describes all environment variables required to run the Ezra API.

## Quick Start

Create a `.env` file in the project root with the following variables:

```env
# Server
PORT=8080

# Database
POSTGRES_URL=postgres://postgres:secret@localhost:5432/ezradb?sslmode=disable

# Security
SECRET=your_secret_key_change_this_in_production

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id_here

# Firebase
FIREBASE_CREDENTIALS_PATH=/path/to/firebase-credentials.json

# Payway Payment Gateway
PAYWAY_MERCHANT_ID=your_merchant_id
PAYWAY_API_KEY=your_api_key
PAYWAY_API_USERNAME=your_api_username
PAYWAY_BASE_URL=https://api-sandbox.payway.com.kh
PAYWAY_RETURN_URL=http://localhost:3000/donation/complete
PAYWAY_CONTINUE_URL=http://localhost:3000/donation/success
PAYWAY_CALLBACK_URL=http://localhost:8080/webhooks/payway

# SMTP Email (Gmail)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-gmail-app-password
SMTP_FROM=your-email@gmail.com
SMTP_SECURE=starttls  # options: starttls (default), ssl, plain

# OTP Settings
OTP_EXPIRY_MINUTES=10
```

## Detailed Configuration

### Server Configuration

#### PORT
- **Description:** HTTP server port
- **Default:** `8080`
- **Example:** `PORT=8080`
- **Required:** No (defaults to 8080 if not set)

---

### Database Configuration

#### POSTGRES_URL
- **Description:** PostgreSQL connection string
- **Format:** `postgres://username:password@host:port/database?sslmode=disable`
- **Example (Local):** `postgres://postgres:secret@localhost:5432/ezradb?sslmode=disable`
- **Example (Docker):** `postgres://postgres:secret@postgres:5432/ezradb?sslmode=disable`
- **Example (Production):** `postgres://user:pass@prod-host:5432/ezradb?sslmode=require`
- **Required:** Yes
- **Security Note:** Use `sslmode=require` in production

---

### Security Configuration

#### SECRET
- **Description:** Secret key for JWT token signing
- **Format:** String (minimum 32 characters recommended)
- **Example:** `SECRET=your_very_long_random_secret_key_here_change_this`
- **Required:** Yes
- **Security Notes:**
  - Must be kept secret
  - Should be different for each environment
  - Generate using: `openssl rand -base64 32`
  - Minimum 32 characters recommended
  - Change regularly in production

---

### Google OAuth Configuration

#### GOOGLE_CLIENT_ID
- **Description:** Google OAuth 2.0 Client ID for Google Sign-In
- **How to Get:**
  1. Go to [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
  2. Create or select a project
  3. Go to "Credentials"
  4. Create "OAuth 2.0 Client ID"
  5. Configure consent screen
  6. Add authorized redirect URIs
  7. Copy the Client ID
- **Example:** `GOOGLE_CLIENT_ID=123456789-abc123def456.apps.googleusercontent.com`
- **Required:** Yes (if using Google Login)
- **Documentation:** [Google OAuth Documentation](https://developers.google.com/identity/protocols/oauth2)

---

### Firebase Configuration

#### FIREBASE_CREDENTIALS_PATH
- **Description:** Path to Firebase service account JSON credentials file
- **How to Get:**
  1. Go to [Firebase Console](https://console.firebase.google.com/)
  2. Select your project
  3. Go to Project Settings > Service Accounts
  4. Click "Generate New Private Key"
  5. Download the JSON file
  6. Save to secure location
  7. Set path in environment variable
- **Example:** `FIREBASE_CREDENTIALS_PATH=/etc/secrets/firebase-credentials.json`
- **Required:** Yes (for push notifications)
- **Security Note:** Keep credentials file secure, never commit to git

---

### Payway Payment Gateway Configuration

#### PAYWAY_MERCHANT_ID
- **Description:** Your Payway merchant account ID
- **How to Get:** Contact Payway or check merchant dashboard
- **Example:** `PAYWAY_MERCHANT_ID=merchant123`
- **Required:** Yes (for donations/payments)

#### PAYWAY_API_KEY
- **Description:** Payway API authentication key
- **How to Get:** Payway merchant dashboard
- **Example:** `PAYWAY_API_KEY=sk_live_abc123def456`
- **Required:** Yes (for donations/payments)
- **Security Note:** Keep secret, never expose in client code

#### PAYWAY_API_USERNAME
- **Description:** Payway API username
- **How to Get:** Payway merchant dashboard
- **Example:** `PAYWAY_API_USERNAME=api_user_123`
- **Required:** Yes (for donations/payments)

#### PAYWAY_BASE_URL
- **Description:** Payway API endpoint URL
- **Sandbox (Testing):** `https://api-sandbox.payway.com.kh`
- **Production:** `https://api.payway.com.kh`
- **Example:** `PAYWAY_BASE_URL=https://api-sandbox.payway.com.kh`
- **Required:** Yes (defaults to sandbox if not set)
- **Note:** Change to production URL when going live

#### PAYWAY_RETURN_URL
- **Description:** URL where users return after payment
- **Example:** `PAYWAY_RETURN_URL=https://yourdomain.com/donation/complete`
- **Required:** Yes
- **Note:** Must be accessible by user's browser

#### PAYWAY_CONTINUE_URL
- **Description:** Success page URL after payment completion
- **Example:** `PAYWAY_CONTINUE_URL=https://yourdomain.com/donation/success`
- **Required:** Yes
- **Note:** Must be accessible by user's browser

#### PAYWAY_CALLBACK_URL
- **Description:** Webhook URL for Payway payment notifications
- **Example:** `PAYWAY_CALLBACK_URL=https://api.yourdomain.com/webhooks/payway`
- **Required:** Yes
- **Security Notes:**
  - Must be publicly accessible
  - Should use HTTPS in production
  - Implement signature verification
  - Set up in Payway merchant dashboard

---

### SMTP Email Configuration

Required for sending OTP verification emails.

#### SMTP_HOST
- **Description:** SMTP server hostname
- **Gmail:** `smtp.gmail.com`
- **SendGrid:** `smtp.sendgrid.net`
- **Amazon SES:** `email-smtp.us-east-1.amazonaws.com`
- **Mailgun:** `smtp.mailgun.org`
- **Example:** `SMTP_HOST=smtp.gmail.com`
- **Required:** Yes (if using OTP)

#### SMTP_PORT
- **Description:** SMTP server port
- **TLS:** `587` (recommended)
- **SSL:** `465`
- **Example:** `SMTP_PORT=587`
- **Required:** Yes (defaults to 587)

#### SMTP_USERNAME
- **Description:** SMTP authentication username
- **Gmail:** Your Gmail address
- **Others:** Provided by email service
- **Example:** `SMTP_USERNAME=your-email@gmail.com`
- **Required:** Yes (if using OTP)

#### SMTP_PASSWORD
- **Description:** SMTP authentication password
- **Gmail:** App Password (NOT regular Gmail password)
- **Others:** API key or password from provider
- **Example:** `SMTP_PASSWORD=abcd efgh ijkl mnop` (remove spaces)
- **Required:** Yes (if using OTP)

**Gmail App Password Setup:**
1. Go to [Google Account Security](https://myaccount.google.com/security)
2. Enable "2-Step Verification"
3. Go to [App Passwords](https://myaccount.google.com/apppasswords)
4. Select "Mail" and "Other (Custom name)"
5. Enter "Ezra API"
6. Click "Generate"
7. Copy the 16-character password
8. Remove spaces: `abcdefghijklmnop`
9. Use in `SMTP_PASSWORD`

#### SMTP_SECURE
- **Description:** Connection security method
- **Options:**
  - `starttls` (default, port 587)
  - `ssl` (implicit TLS, port 465)
  - `plain` (no TLS; not recommended)
- **Auto-defaults:** If `SMTP_PORT=465`, defaults to `ssl`; otherwise `starttls`.
- **Example:** `SMTP_SECURE=starttls`

#### SMTP_FROM
- **Description:** Email address shown as sender
- **Example:** `SMTP_FROM=noreply@yourdomain.com`
- **Required:** No (defaults to SMTP_USERNAME)
- **Best Practice:** Use no-reply address for automated emails

---

### OTP Configuration

#### OTP_EXPIRY_MINUTES
- **Description:** OTP code validity duration in minutes
- **Default:** `10`
- **Recommended:** `5-15` minutes
- **Example:** `OTP_EXPIRY_MINUTES=10`
- **Required:** No (defaults to 10)
- **Security Note:** Shorter is more secure, but less user-friendly

---

## Alternative SMTP Providers

### SendGrid

```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your_sendgrid_api_key
SMTP_FROM=noreply@yourdomain.com
```

**Setup:**
1. Sign up at [SendGrid](https://sendgrid.com/)
2. Create API key with "Mail Send" permission
3. Use "apikey" as username
4. Use API key as password

### Amazon SES

```env
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
SMTP_USERNAME=your_ses_smtp_username
SMTP_PASSWORD=your_ses_smtp_password
SMTP_FROM=noreply@yourdomain.com
```

**Setup:**
1. Sign up for [AWS](https://aws.amazon.com/)
2. Go to SES console
3. Verify domain or email address
4. Create SMTP credentials
5. Use region-specific SMTP endpoint

### Mailgun

```env
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USERNAME=postmaster@your-domain.com
SMTP_PASSWORD=your_mailgun_smtp_password
SMTP_FROM=noreply@yourdomain.com
```

**Setup:**
1. Sign up at [Mailgun](https://www.mailgun.com/)
2. Add and verify domain
3. Get SMTP credentials from dashboard
4. Configure DNS records (SPF, DKIM)

### Mailtrap (Development Only)

```env
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USERNAME=your_mailtrap_username
SMTP_PASSWORD=your_mailtrap_password
SMTP_FROM=test@example.com
```

**Setup:**
1. Sign up at [Mailtrap](https://mailtrap.io/)
2. Create inbox
3. Copy SMTP credentials
4. Use for testing (emails won't be delivered to real addresses)

---

## Environment-Specific Configurations

### Development (.env.development)

```env
PORT=8080
POSTGRES_URL=postgres://postgres:dev@localhost:5432/ezradb_dev?sslmode=disable
SECRET=development_secret_key_not_for_production
PAYWAY_BASE_URL=https://api-sandbox.payway.com.kh
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
OTP_EXPIRY_MINUTES=15
```

**Characteristics:**
- Local database
- Sandbox payment gateway
- Test email service (Mailtrap)
- Longer OTP expiry for testing
- Debug logging enabled

### Staging (.env.staging)

```env
PORT=8080
POSTGRES_URL=postgres://staginguser:password@staging-db:5432/ezradb_staging?sslmode=require
SECRET=staging_secret_key_different_from_production
PAYWAY_BASE_URL=https://api-sandbox.payway.com.kh
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
OTP_EXPIRY_MINUTES=10
```

**Characteristics:**
- Staging database
- Sandbox payment gateway
- Real email service
- Production-like configuration
- Testing with real services

### Production (.env.production)

```env
PORT=8080
POSTGRES_URL=postgres://produser:strongpassword@prod-db:5432/ezradb?sslmode=require
SECRET=super_secure_random_production_secret_key_min_32_chars
PAYWAY_BASE_URL=https://api.payway.com.kh
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
OTP_EXPIRY_MINUTES=10
```

**Characteristics:**
- Production database with SSL
- Production payment gateway
- Professional email service
- Maximum security settings
- Monitoring and logging

---

## Production Security Checklist

Before deploying to production:

### ✅ Secret Management
- [ ] Generate strong SECRET key (min 32 characters)
  ```bash
  openssl rand -base64 32
  ```
- [ ] Use different secrets for each environment
- [ ] Store secrets in secure vault (AWS Secrets Manager, HashiCorp Vault)
- [ ] Rotate secrets regularly
- [ ] Never commit .env to version control

### ✅ Database Security
- [ ] Use strong database password
- [ ] Enable SSL/TLS (sslmode=require)
- [ ] Set up regular backups
- [ ] Configure database firewall rules
- [ ] Enable audit logging
- [ ] Use read-only replicas for analytics

### ✅ HTTPS Configuration
- [ ] All URLs use HTTPS
- [ ] SSL certificates are valid
- [ ] HSTS headers enabled
- [ ] Redirect HTTP to HTTPS
- [ ] Use Let's Encrypt or commercial SSL

### ✅ API Security
- [ ] Rate limiting enabled
- [ ] CORS properly configured
- [ ] Input validation on all endpoints
- [ ] SQL injection protection
- [ ] XSS protection
- [ ] CSRF tokens where applicable

### ✅ Email Security
- [ ] Use professional domain (@yourdomain.com)
- [ ] Configure SPF records
- [ ] Configure DKIM records
- [ ] Configure DMARC records
- [ ] Monitor bounce rates
- [ ] Set up abuse handling

### ✅ Monitoring & Logging
- [ ] Application logging enabled
- [ ] Error tracking (Sentry, Rollbar)
- [ ] Performance monitoring (New Relic, Datadog)
- [ ] Uptime monitoring
- [ ] Alert system configured
- [ ] Log rotation configured

### ✅ OTP Security
- [ ] Rate limiting (max 3 OTPs per hour per email)
- [ ] Failed verification attempt limits
- [ ] IP-based rate limiting
- [ ] Suspicious activity detection
- [ ] Email delivery monitoring

### ✅ Payment Security
- [ ] PCI DSS compliance
- [ ] Webhook signature verification
- [ ] Transaction logging
- [ ] Fraud detection
- [ ] Refund handling
- [ ] Chargeback monitoring

---

## Troubleshooting

### Database Connection Issues

**Error:** `failed to connect database`

**Solutions:**
1. Check POSTGRES_URL format is correct
2. Verify database is running
3. Check firewall allows connection
4. Verify credentials are correct
5. Check SSL mode setting

### Email Sending Issues

**Error:** `failed to send email`

**Solutions:**
1. Verify SMTP credentials
2. Check Gmail App Password (not regular password)
3. Enable 2-Step Verification for Gmail
4. Check firewall allows port 587
5. Verify email service is not blocked

### OTP Not Received

**Solutions:**
1. Check spam/junk folder
2. Verify email service logs
3. Check SMTP credentials
4. Verify email address format
5. Check rate limiting

### Payment Gateway Issues

**Error:** Payment webhook not received

**Solutions:**
1. Verify PAYWAY_CALLBACK_URL is publicly accessible
2. Check Payway merchant dashboard webhook settings
3. Verify webhook endpoint is working
4. Check firewall rules
5. Review webhook logs

---

## Additional Resources

- **[ROUTEMAP.md](./ROUTEMAP.md)** - Complete API documentation
- **[README.md](./README.md)** - Project overview and quick start
- **[ARCHITECTURE.md](./ARCHITECTURE.md)** - System architecture
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Deployment guide

---

## Support

For issues or questions:
- GitHub Issues: [your-repo-url]/issues
- Email: support@yourdomain.com
- Documentation: https://docs.yourdomain.com

