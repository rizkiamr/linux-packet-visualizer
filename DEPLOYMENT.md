# Cloudflare Pages Deployment Guide

This project is configured to run from the `frontend` directory.

## Build Configuration

When connecting your GitHub repository to Cloudflare Pages, use the following settings:

| Setting | Value | Description |
|---------|-------|-------------|
| **Framework Preset** | `None` | Select **None** if Vite is not listed |
| **Build Command** | `npm run build` | Command to compile the app |
| **Build Output Directory** | `dist` | Folder containing built assets |
| **Root Directory** | `frontend` | **Important:** The app lives in this subdirectory |

## Environment Variables

No special environment variables are required for the default build.

## Manual Deployment (Wrangler)

If deploying via CLI:

```bash
# Preview Deployment
npx wrangler pages deploy frontend/dist --project-name linux-packet-visualizer

# Production Deployment (main branch)
npx wrangler pages deploy frontend/dist --project-name linux-packet-visualizer --branch main
```
