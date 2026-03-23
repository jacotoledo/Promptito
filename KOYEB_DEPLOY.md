# Deploy Promptito to Koyeb (Free Tier)

## Prerequisites
- A GitHub account
- A forked copy of this repository

---

## Step-by-Step Deployment

### Step 1: Fork the Repository
1. Go to this repository on GitHub
2. Click the **Fork** button (top right)
3. Wait for the fork to complete

### Step 2: Create a New Service on Koyeb
1. Go to [app.koyeb.com](https://app.koyeb.com) and sign in
2. Click **Create Service**
3. Select **Deploy from GitHub**

### Step 3: Connect Your GitHub Account (First Time Only)
1. Click **Connect GitHub account**
2. Authorize Koyeb to access your repositories
3. Select your forked repository from the list

### Step 4: Configure the Build (CRITICAL!)
**This is where most people get stuck!**

After selecting your repo, look for the **Builder** option:

1. **DO NOT leave it on "Auto-detect" or "Buildpack"**
2. Click the dropdown and select **Dockerfile**
3. Koyeb will automatically find your `Dockerfile` in the repo root
4. Leave the Dockerfile location as `./Dockerfile`

### Step 5: Configure the Port
1. Scroll to **Exposing your service**
2. Set the port to: `80`

### Step 6: Environment Variables (Optional)
If your app requires any secrets:
1. Go to **Environment variables**
2. Add any required variables
3. Mark sensitive ones as secret

### Step 7: Deploy
1. Click **Deploy**
2. Wait for the build to complete (2-5 minutes)
3. Your app will be live at `https://<your-app-name>.koyeb.app`

---

## Troubleshooting Common Issues

### "Build failed" or "Application error"
- Make sure you selected **Dockerfile** as the builder (Step 4)
- Check the build logs by clicking on the deployment in the dashboard
- Verify the port is set to `80`

### "Health check failing"
- The Dockerfile includes a health check on `/health`
- Make sure your app responds to `http://localhost:80/health`
- Wait 30-60 seconds for the first health check

### Koyeb uses Buildpack instead of Dockerfile
This happens when you don't explicitly select Dockerfile:
1. Go to your Service settings
2. Find **Builder** in the configuration
3. Change from `Auto` or `Buildpack` to **Dockerfile**
4. Redeploy

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

## Updating Your Deployment

Any push to the `main` branch will automatically trigger a new deployment.

### Manual Redeploy
1. Go to your Service on Koyeb
2. Click **Redeploy**

---

## Uninstall

1. Go to your Service on Koyeb
2. Click **Settings**
3. Scroll to the bottom and click **Delete Service**

---

**Note:** Koyeb is a third-party service. This project is not affiliated with Koyeb. Use at your own risk.
