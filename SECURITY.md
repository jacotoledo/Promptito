# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 2.x.x   | :white_check_mark: |
| 1.x.x   | :x:                |

---

## Report a Vulnerability

We take security seriously. If you find a vulnerability, please report it responsibly.

### How to Report

1. **Do NOT** create a public GitHub Issue for security vulnerabilities
2. Email directly: **security@jtg365.com**
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Any suggested fixes (optional)

### What to Expect

- **Acknowledgment**: Within 48 hours
- **Initial Assessment**: Within 7 days
- **Fix Timeline**: Depends on severity, but we'll work with you on disclosure

---

## Security Design Principles

Promptito is designed with these principles:

### Read-Only by Design

Promptito is a **read-only** server. It:
- Only serves files - never writes or modifies
- Has no database - data lives in flat files
- Exposes no admin endpoints
- Accepts no user-uploaded content

### Defense in Depth

| Protection | Implementation |
|------------|----------------|
| **Input Validation** | All slugs validated against strict regex |
| **Path Traversal** | `..` sequences blocked, paths cleaned |
| **DoS Protection** | Request body limited to 1MB |
| **Rate Limiting** | 100 requests/second per IP |
| **Error Sanitization** | Internal paths never exposed to clients |
| **Security Headers** | CSP, HSTS, X-Frame-Options, and more |
| **Thread Safety** | Concurrent-safe with RWMutex |

### What Promptito Does NOT Do

- **No Authentication**: Designed for trusted networks or local use
- **No Encryption**: Traffic should be protected at the network level (HTTPS via reverse proxy)
- **No User Content**: Can't upload or modify prompts through the server
- **No Analytics**: Zero tracking or telemetry

---

## For Self-Hosted Deployments

### Recommended Security Practices

1. **Run behind a reverse proxy** (nginx, Caddy) with HTTPS
2. **Firewall the server** - only expose port 80 (or 443 for HTTPS) to necessary networks
3. **Regular updates** - pull latest releases for security patches
4. **File permissions** - ensure prompt files are readable but not writable by the server process
5. **Monitor logs** - watch for unusual access patterns

### HTTPS Setup (Koyeb)

Koyeb provides HTTPS automatically. For other hosts:

```nginx
# nginx example
server {
    listen 443 ssl;
    server_name your-domain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:80;
    }
}
```

---

## Responsible Disclosure

We follow responsible disclosure practices:
- We will credit researchers who report valid vulnerabilities (with permission)
- We ask that you give us reasonable time to fix issues before public disclosure
- We will work with you on disclosure timeline for critical issues

---

**Thank you for helping keep Promptito secure!**
