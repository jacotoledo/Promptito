# Deploy Promptito to Koyeb (Free Tier)

## Prerequisites
- A GitHub account
- A forked copy of this repository

---

## Step-by-Step Deployment

### Step 1: Fork the Repository
1. Go to [github.com/jacotoledo/Promptito](https://github.com/jacotoledo/Promptito)
2. Click the **Fork** button (top right)
3. Wait for the fork to complete

> **Warning:** Do NOT use GitHub's "Upload files" feature. Always push changes from your local git repo. Uploading files via the web interface corrupts directory structure and breaks the build.

### Step 2: Create a New Service on Koyeb
1. Go to [app.koyeb.com](https://app.koyeb.com) and sign in
2. Click **Create Service**
3. Select **Deploy from GitHub**

### Step 3: Connect Your GitHub Account (First Time Only)
1. Click **Connect GitHub account**
2. Authorize Koyeb to access your repositories
3. Select your forked repository from the list

### Step 4: Configure the Build (CRITICAL!)
This is where most people get stuck!

1. In the **Builder** dropdown, select **Dockerfile** (NOT Auto-detect or Buildpack)
2. Koyeb will automatically find the `Dockerfile` in the repo root
3. Leave the Dockerfile location as `./Dockerfile`

### Step 5: Configure Instance
1. Under **Instance**, select **CPU Eco**
2. Choose **Free** tier
3. Select a region: **Washington, D.C.** or **Frankfurt** (free tier only)

### Step 6: Configure Port
1. Scroll to **Exposing your service**
2. Set the port to: `8080` (Koyeb free tier doesn't allow privileged ports below 1024)
3. Keep **HTTPS enabled** for best compatibility

### Step 7: Deploy
1. Click **Deploy**
2. Wait for the build to complete (2-5 minutes)
3. Your app will be live at `https://<your-app-name>.koyeb.app`

---

## Troubleshooting

### "Build failed" or "directory not found"
- Make sure you selected **Dockerfile** as the builder (Step 4)
- Check if files were uploaded via GitHub web instead of git push
- Delete and re-fork the repo if necessary

### "Permission denied" on port 80
- Port 80 is privileged on Linux. Always use **port 8080**.

### Buttons or JavaScript not working
- Ensure HTTPS is enabled on your Koyeb service
- Clear browser cache after deployment
- Check browser console for CSP errors

### Health check failing
- The app exposes a `/health` endpoint
- Wait 30-60 seconds for the first health check

---

## Updating Your Deployment

Any push to the `main` branch will automatically trigger a new deployment.

### Manual Redeploy
1. Go to your Service on Koyeb
2. Click **Redeploy**

---

## Custom Domain Setup

1. In Koyeb dashboard, go to your Service **Settings**
2. Click **Domains**
3. Add your domain (e.g., `prompts.yourdomain.com`)
4. Add the CNAME record shown to your DNS provider
5. SSL certificates are automatic

---

## Keep Your App Awake

Koyeb's free tier sleeps apps after 30 days of inactivity.

### Using UptimeRobot (Free)
1. Go to [uptimerobot.com](https://uptimerobot.com)
2. Create a free account
3. Add a new HTTP(s) monitor
4. Enter your Koyeb URL: `https://<your-app-name>.koyeb.app`
5. Set check interval to 5 minutes

---

## Uninstall

1. Go to your Service on Koyeb
2. Click **Settings**
3. Scroll to the bottom and click **Delete Service**

---

**Note:** Koyeb is a third-party service. This project is not affiliated with Koyeb. Use at your own risk.
