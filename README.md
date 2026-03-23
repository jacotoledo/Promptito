# Promptito

<p align="center">
  <a href="https://github.com/jacotoledo/Promptito/stargazers">
    <img src="https://img.shields.io/github/stars/jacotoledo/Promptito?style=social" alt="Stars">
  </a>
  <a href="https://github.com/jacotoledo/Promptito/network/members">
    <img src="https://img.shields.io/github/forks/jacotoledo/Promptito?style=social" alt="Forks">
  </a>
  <img src="https://img.shields.io/github/license/jacotoledo/Promptito" alt="License">
  <img src="https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-blue" alt="Platform">
</p>

<p align="center">
  <strong>PROMPT</strong> + <strong>RAPIDITO</strong> (Spanish for <em>super fast</em>)<br>
  <em>The fastest way to manage AI prompts</em>
</p>

<p align="center">
  <a href="https://github.com/jacotoledo/Promptito">
    <img src="https://raw.githubusercontent.com/jacotoledo/Promptito/main/public/screenshots/Screenshot.png" alt="Promptito UI" width="800">
  </a>
</p>

<p align="center">
  <a href="https://github.com/jacotoledo/Promptito/releases/latest">
    <img src="https://img.shields.io/badge/Download-Windows-blue?style=flat-square&logo=windows" alt="Windows">
  </a>
  <a href="https://app.koyeb.com">
    <img src="https://img.shields.io/badge/Deploy-Free%20on%20Koyeb-black?style=flat-square" alt="Koyeb">
  </a>
  <a href="https://github.com/jacotoledo/Promptito/blob/main/CHANGELOG.md">
    <img src="https://img.shields.io/badge/Changelog-v2.0.0-green?style=flat-square" alt="Changelog">
  </a>
</p>

---

> **Your AI prompts, searchable, copyable, and shareable.**
>
> Built by someone who got tired of scrolling through 47 Notion docs labeled `"v2_FINAL_real.md"` just to find one good prompt.

---

## Why Promptito?

| Feature | Promptito | Notion | GitHub Gists |
|---------|:---------:|:------:|:------------:|
| Zero dependencies | ✅ | ❌ | ❌ |
| Self-hosted | ✅ | ❌ | ❌ |
| No account required | ✅ | ❌ | ❌ |
| Free forever | ✅ | 💰 | ✅ |
| REST API included | ✅ | ❌ | ❌ |
| Search & filter | ✅ | ✅ | ❌ |
| Bundle download | ✅ | ❌ | ❌ |
| Standards metadata | ✅ | ❌ | ❌ |

---

## Get Started

### Option 1: Free Website (60 seconds)

1. Fork this repo (use the **Fork** button, do NOT upload files manually)
2. Go to [app.koyeb.com](https://app.koyeb.com)
3. Sign up with GitHub
4. Click **Create Service** → **Deploy from GitHub** → select your fork
5. **IMPORTANT:** In the Builder dropdown, select **Dockerfile** (not Auto-detect)
6. Under **Instance**, select **CPU Eco** → **Free**
7. Under **Exposing your service**, set port to **8080** and **uncheck HTTPS** (Koyeb free tier has SSL limitations with custom domains; for now use HTTP)
8. Click **Deploy**

Done. Your prompts are live at `https://your-app.koyeb.app`

> **Note:** If deployment fails with "directory not found", you may have uploaded files via GitHub instead of pushing with git. Delete and re-fork the repo.

### Option 2: Run Locally

1. Download from [Releases](https://github.com/jacotoledo/Promptito/releases)
2. Double-click `promptito.exe`
3. Open [http://localhost](http://localhost)

---

## Features

| Feature | What it does |
|---------|--------------|
| **Search** | Find prompts by name, description, or content |
| **Filter** | By category, tags, or SFIA skill level |
| **Copy** | One click to copy any prompt |
| **Bundle** | Download multiple prompts as a ZIP |
| **API** | For AI agents and scripts |

---

## Add Your Own Prompts

1. Create a folder: `public/prompts/my-prompt/`
2. Add a file: `public/prompts/my-prompt/SKILL.md`

```markdown
---
name: My Prompt
description: What it does
---

# Role
You are a helpful AI that...

# Instructions
Do this and that...
```

3. Refresh the page.

That's it. No database. No config files. Just files.

---

## API (For Developers)

```bash
GET /api/skills          # All prompts
GET /api/skills/{slug}  # One prompt
GET /api/search?q=text   # Search
GET /api/tags           # List tags
GET /api/bundle         # Download as ZIP
```

---

## Tech Stack

- **Go** - Zero dependencies, single binary
- **Self-hosted** - Your data stays on your server
- **No tracking** - No analytics, no telemetry, no cookies
- **Standards-compliant** - IPTC, ISO/IEC 5259, NIST AI RMF, SFIA 9

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## Support

If Promptito saves you time, consider leaving a tip:

<a href="https://jtg365.com/tip">
  <img src="https://img.shields.io/badge/Tip-JTG365-orange?style=flat-square&logo=bitcoin" alt="Tip Jar">
</a>

---

## License

MIT — see [LICENSE](LICENSE)

---

<p align="center">
  Built with ❤️ by <a href="https://jtg365.com">Jaco Toledo</a><br>
  <a href="https://github.com/jacotoledo">GitHub</a> · <a href="https://jacotoledo.github.io/Promptito">Website</a> · <a href="https://jtg365.com">jtg365.com</a>
</p>
