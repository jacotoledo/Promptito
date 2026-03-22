# Promptito

<p align="center">
  <strong>PROMPT</strong> + <strong>RAPIDITO</strong> (Spanish for <em>super fast</em>)<br>
  <em>The fastest way to manage AI prompts</em>
</p>

<p align="center">
  Made with ❤️ by <a href="https://jtg365.com">Jaco Toledo</a><br>
  <a href="https://github.com/jacotoledo">GitHub</a> · <a href="https://jtg365.com">jtg365.com</a>
</p>

---

> **Your AI prompts, searchable, copyable, and shareable.**
> 
> Built by someone who got tired of scrolling through 47 Notion docs labeled `"v2_FINAL_real.md"` just to find one good prompt.

---

## The Problem

Your prompt library probably looks like this:

```
prompts/
├── chatgpt_prompts_v3_USE_THIS_ONE.txt
├── chatgpt_prompts_v2.txt
├── chatgpt_prompts_v2_actual.txt
├── chatgpt_prompts_OLD_DONT_USE.txt
└── why_am_i_like_this.jpg
```

**Promptito is the antidote.** Simple folders, simple format, one search bar.

---

## Get Started

### Option 1: Run Locally (30 seconds)

1. Download from [Releases](https://github.com/jtg365/promptito/releases)
2. Double-click `promptito.exe`
3. Open [http://localhost:8080](http://localhost:8080)

### Option 2: Free Website (Koyeb)

1. Fork this repo
2. Go to [app.koyeb.com](https://app.koyeb.com)
3. Sign up with GitHub
4. Click **Create Service** → **Deploy from GitHub** → select your fork → **Deploy**

Done. Your prompts are live at `https://your-app.koyeb.app`

*No credit card. No server. No kidding.*

---

## Features

| Feature | What it does |
|---------|--------------|
| **Search** | Find prompts by name, description, or content |
| **Filter** | By category, tags, or skill level |
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
GET /api/skills         # All prompts
GET /api/skills/{name}  # One prompt
GET /api/search?q=text  # Search
```

---

## Tech Stuff

- **Zero dependencies** - It's just a Go binary
- **Self-hosted** - Your data stays on your server
- **No tracking** - No analytics, no telemetry, no cookies
- **Standards-compliant** - Uses real AI metadata standards

---

## License

MIT — see [LICENSE](LICENSE)

---

<p align="center">
Built with slightly too much coffee and a genuine belief that AI prompts deserve to be organized better than my desktop.
</p>
