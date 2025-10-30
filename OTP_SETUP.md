# OTP Email Verification Setup Guide

This guide explains how to set up and use the OTP (One-Time Password) email verification feature with Gmail SMTP.

## 📋 Table of Contents

- [Features](#features)
- [Setup Gmail for SMTP](#setup-gmail-for-smtp)
- [Environment Variables](#environment-variables)
- [API Endpoints](#api-endpoints)
- [Usage Examples](#usage-examples)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

## ✨ Features

- **Email Verification**: Verify user email addresses during registration
- **Password Reset**: Send OTP codes for secure password reset
- **Login Verification**: Add extra security with OTP-based login
- **Automatic Expiry**: OTP codes expire after 10 minutes (configurable)
- **Beautiful Email Template**: Professional HTML email with responsive design
- **Rate Limiting**: Prevents spam by deleting old OTPs before sending new ones

## 🔧 Setup Gmail for SMTP

### Step 1: Enable 2-Step Verification

1. Go to [Google Account Security](https://myaccount.google.com/security)
2. Enable **2-Step Verification** if not already enabled
3. Follow the prompts to complete setup

### Step 2: Generate App Password

1. Go to [App Passwords](https://myaccount.google.com/apppasswords)
2. Select **Mail** as the app
3. Select **Other (Custom name)** as the device
4. Enter a name like "Ezra API"
5. Click **Generate**
6. **Copy the 16-character password** (this is your SMTP password)

### Step 3: Configure Environment Variables

Update your `.env` file with the following:

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=xxxx xxxx xxxx xxxx  # 16-character app password from Step 2
SMTP_FROM=your-email@gmail.com
OTP_EXPIRY_MINUTES=10
```

## 🌐 Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SMTP_HOST` | SMTP server host | `smtp.gmail.com` | No |
| `SMTP_PORT` | SMTP server port | `587` | No |
| `SMTP_USERNAME` | Gmail email address | - | **Yes** |
| `SMTP_PASSWORD` | Gmail app password | - | **Yes** |
| `SMTP_FROM` | From email address | Same as `SMTP_USERNAME` | No |
| `OTP_EXPIRY_MINUTES` | OTP expiration time in minutes | `10` | No |

## 📡 API Endpoints

### 1. Send OTP

**Endpoint:** `POST /otp/send`

**Description:** Generate and send an OTP code to the specified email address.

**Request Body:**
```json
{
  "email": "user@example.com",
  "purpose": "email_verification"
}
```

**Purpose Options:**
- `email_verification` - For verifying new email addresses during registration
- `password_reset` - For password reset functionality
- `login` - For two-factor authentication during login

**Success Response (200):**
```json
{
  "message": "OTP sent successfully to your email",
  "email": "user@example.com",
  "expires_at": "2025-10-30T10:45:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid email format or email already registered (for email_verification)
- `400 Bad Request` - Email not found (for password_reset)

### 2. Verify OTP

**Endpoint:** `POST /otp/verify`

**Description:** Verify the OTP code sent to the email address.

**Request Body:**
```json
{
  "email": "user@example.com",
  "code": "123456",
  "purpose": "email_verification"
}
```

**Success Response (200):**
```json
{
  "message": "OTP verified successfully",
  "data": {
    "email": "user@example.com",
    "purpose": "email_verification"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid or expired OTP
- `400 Bad Request` - OTP already used

## 💡 Usage Examples

### Example 1: Email Verification Flow

```bash
# Step 1: Send OTP to user's email
curl -X POST http://localhost:8080/otp/send \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "purpose": "email_verification"
  }'

# Step 2: User receives email with 6-digit code
# User enters: 123456

# Step 3: Verify the OTP code
curl -X POST http://localhost:8080/otp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "code": "123456",
    "purpose": "email_verification"
  }'

# Step 4: If successful, proceed with user registration
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "fullname": "New User",
    "email": "newuser@example.com",
    "password": "securepassword"
  }'
```

### Example 2: Password Reset Flow

```bash
# Step 1: Request password reset OTP
curl -X POST http://localhost:8080/otp/send \
  -H "Content-Type: application/json" \
  -d '{
    "email": "existing@example.com",
    "purpose": "password_reset"
  }'

# Step 2: Verify OTP
curl -X POST http://localhost:8080/otp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "existing@example.com",
    "code": "654321",
    "purpose": "password_reset"
  }'

# Step 3: After verification, allow user to reset password
# (Password reset endpoint would be implemented separately)
```

### Example 3: Two-Factor Authentication

```bash
# Step 1: User provides username/password for login
# Step 2: Backend sends OTP for additional verification
curl -X POST http://localhost:8080/otp/send \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "purpose": "login"
  }'

# Step 3: User enters OTP code
curl -X POST http://localhost:8080/otp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "code": "789012",
    "purpose": "login"
  }'

# Step 4: Complete login process
```

## 🧪 Testing

### Test with cURL

```bash
# Test sending OTP
curl -X POST http://localhost:8080/otp/send \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","purpose":"email_verification"}'

# Check your email for the 6-digit code

# Test verifying OTP
curl -X POST http://localhost:8080/otp/verify \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","code":"123456","purpose":"email_verification"}'
```

### Email Template Preview

The OTP email includes:
- Professional header with branding
- Clear purpose message
- Large, easy-to-read 6-digit code
- Expiration warning (10 minutes)
- Security reminder
- Responsive design for mobile devices

## 🔍 Troubleshooting

### Issue: "Failed to send email"

**Possible Causes:**
1. **Incorrect Gmail credentials**
   - Solution: Double-check SMTP_USERNAME and SMTP_PASSWORD
   - Make sure you're using an App Password, not your regular Gmail password

2. **2-Step Verification not enabled**
   - Solution: Enable 2-Step Verification in Google Account settings

3. **Less secure app access blocked**
   - Solution: Use App Passwords instead (recommended)

4. **Firewall blocking port 587**
   - Solution: Check firewall settings or try port 465 with SSL

### Issue: "OTP has expired"

**Solution:** 
- OTP codes expire after 10 minutes (default)
- Request a new OTP using the `/otp/send` endpoint
- Adjust `OTP_EXPIRY_MINUTES` if needed

### Issue: "Email already registered"

**Solution:**
- This is expected behavior for `email_verification` purpose
- Use a different email or use `password_reset` purpose instead

### Issue: "Invalid or expired OTP"

**Possible Causes:**
1. Incorrect code entered
2. Code has expired (>10 minutes)
3. Code already used
4. Wrong purpose specified

**Solution:** Request a new OTP

## 📊 Database Schema

The OTP feature uses the `otps` table:

```sql
CREATE TABLE IF NOT EXISTS otps (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    code VARCHAR(10) NOT NULL,
    purpose VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
```

## 🔒 Security Features

1. **Automatic Cleanup**: Old OTPs are automatically deleted when sending new ones
2. **One-Time Use**: OTPs are marked as verified after successful use
3. **Time-Limited**: Codes expire after specified duration
4. **Purpose-Specific**: Different purposes prevent code reuse across features
5. **Email Validation**: Ensures proper email format before sending

## 📝 Best Practices

1. **Use HTTPS in production** to protect OTP codes in transit
2. **Implement rate limiting** to prevent spam/abuse
3. **Log OTP operations** for security auditing
4. **Clear sensitive data** from logs (don't log OTP codes)
5. **Set reasonable expiry times** (5-15 minutes recommended)
6. **Add brute-force protection** (limit verification attempts)

## 🚀 Production Recommendations

1. **Use environment variables** for all sensitive configuration
2. **Enable email logging** to track delivery issues
3. **Monitor OTP usage** for suspicious patterns
4. **Set up email alerts** for failed deliveries
5. **Consider SMS alternatives** for critical operations
6. **Implement backup delivery methods** (SMS, voice call)

## 📧 Alternative SMTP Providers

While this guide focuses on Gmail, you can use other SMTP providers:

### SendGrid
```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your_sendgrid_api_key
```

### Amazon SES
```env
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
SMTP_USERNAME=your_ses_username
SMTP_PASSWORD=your_ses_password
```

### Mailgun
```env
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USERNAME=postmaster@your-domain.com
SMTP_PASSWORD=your_mailgun_password
```

## 🆘 Support

If you encounter issues:

1. Check the application logs for detailed error messages
2. Verify all environment variables are set correctly
3. Test SMTP connection independently
4. Review Gmail security settings
5. Check spam folder for OTP emails

---

**Note:** Keep your SMTP credentials secure and never commit them to version control!

