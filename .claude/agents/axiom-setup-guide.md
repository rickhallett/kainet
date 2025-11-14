# Axiom Setup Guide
## Configuring Axiom to Receive Logs from Next.js Application

**Date Created:** 2025-11-13
**Purpose:** Step-by-step guide for setting up Axiom account and configuring log ingestion

---

## Step 1: Create Axiom Account

1. **Visit Axiom Website**
   - Go to https://axiom.co
   - Click "Sign Up" or "Get Started"

2. **Choose Plan**
   - **Free Tier**: 500MB/day (perfect for initial testing)
   - **Hobby Plan**: $25/month for 50GB/month (recommended for MVP)
   - You can start with free tier and upgrade later

3. **Sign Up Options**
   - Sign up with email
   - Or use GitHub/Google OAuth for quick signup

4. **Verify Email**
   - Check your email for verification link
   - Click to verify your account

---

## Step 2: Create Organization and Dataset

1. **Create Organization** (if prompted)
   - Organization name: `becoming-diamond` (or your preferred name)
   - This is your workspace for all datasets

2. **Create Your First Dataset**
   - After logging in, click **"New Dataset"** or **"Create Dataset"**
   - Dataset name: `becoming-diamond-prod` (for production logs)
   - Description: "Production logs for Becoming Diamond Next.js app"
   - Click **"Create"**

   **Note**: Datasets are like database tables - they store your log data. You can create multiple datasets for different environments:
   - `becoming-diamond-dev` - Development logs
   - `becoming-diamond-staging` - Staging logs
   - `becoming-diamond-prod` - Production logs

3. **Note Your Dataset Name**
   - You'll need this exact name for the `AXIOM_DATASET` environment variable

---

## Step 3: Generate API Token

1. **Navigate to Settings**
   - Click your profile icon in top right
   - Select **"Settings"** or **"Organization Settings"**

2. **Go to API Tokens**
   - In left sidebar, click **"API Tokens"** or **"Tokens"**
   - Click **"Create Token"** or **"New Token"**

3. **Configure Token**
   - **Name**: `becoming-diamond-nextjs-app` (or descriptive name)
   - **Description**: "API token for Next.js application logging"
   - **Permissions**: Select **"Ingest"** permission (write-only for security)
   - **Datasets**: Select `becoming-diamond-prod` (or all datasets if you want flexibility)
   - **Expiration**: Choose "Never" for production (or set expiration for added security)

4. **Copy Token**
   - Click **"Create Token"**
   - **IMPORTANT**: Copy the token immediately and store it securely
   - Format: `xaat-xxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`
   - **You won't be able to see it again!**

5. **Store Token Securely**
   - Save in password manager
   - Add to `.env.local` file (never commit to Git)
   - Add to deployment platform (Vercel) environment variables

---

## Step 4: Configure Environment Variables

### Local Development (.env.local)

Create or update `.env.local` in your project root:

```bash
# Axiom Configuration
AXIOM_TOKEN=xaat-your-actual-token-here
AXIOM_DATASET=becoming-diamond-prod

# Optional: Organization ID (usually not required)
# AXIOM_ORG_ID=your-org-id
```

**Security Note**:
- `.env.local` should be in `.gitignore` (it already is by default in Next.js)
- Never commit API tokens to version control

### Production Deployment (Vercel)

1. **Go to Vercel Dashboard**
   - Navigate to your project
   - Click **"Settings"** tab
   - Click **"Environment Variables"** in sidebar

2. **Add Axiom Variables**

   **Variable 1:**
   - Name: `AXIOM_TOKEN`
   - Value: `xaat-your-actual-token-here` (paste from Axiom)
   - Environments: Check **"Production"**, **"Preview"**, **"Development"**
   - Click **"Save"**

   **Variable 2:**
   - Name: `AXIOM_DATASET`
   - Value: `becoming-diamond-prod`
   - Environments: Check **"Production"**, **"Preview"**, **"Development"**
   - Click **"Save"**

3. **Redeploy Application**
   - After adding environment variables, redeploy your app
   - Go to **"Deployments"** tab
   - Click **"..."** on latest deployment → **"Redeploy"**
   - Or push a new commit to trigger deployment

---

## Step 5: Verify Log Ingestion

### Test Locally First

1. **Start Development Server**
   ```bash
   npm run dev
   ```

2. **Trigger Log Events**
   - Visit `http://localhost:3003`
   - Navigate to `/app` (triggers authentication)
   - Try signing in (triggers auth logs)
   - Visit profile page (triggers profile logs)

3. **Check Axiom Dashboard**
   - Go to Axiom dashboard
   - Click on `becoming-diamond-prod` dataset
   - Click **"Stream"** tab to see live logs
   - You should see logs appearing in real-time

4. **Expected Log Entries**
   - Middleware request logs (automatic from `withAxiom`)
   - Authentication logs (from `auth.ts`)
   - Profile API logs (from `/api/profile`)
   - Any errors caught by ErrorBoundary

### Verify in Production

1. **After Deployment**
   - Visit your production URL
   - Perform key actions (sign in, navigate, etc.)

2. **Check Axiom Stream**
   - Logs should appear within seconds
   - Real-time streaming with millisecond latency

3. **Test Error Tracking**
   - Trigger an error (if safe in development)
   - Check if error appears in Axiom with proper context

---

## Step 6: Explore Axiom Features

### Stream View (Real-time Logs)

1. **Navigate to Dataset**
   - Click `becoming-diamond-prod`
   - Default view is "Stream" tab

2. **Live Tail**
   - Logs appear in real-time as they're ingested
   - Use filters to narrow down specific logs

3. **Filtering**
   - Click on any field to filter by that value
   - Example: Click on a userId to see all logs for that user
   - Use search bar for complex queries

### Query Builder (Advanced Search)

1. **Click "Analytics" or "Query" Tab**

2. **Use APL (Axiom Processing Language)**
   - Similar to KQL (Kusto Query Language)
   - Example queries:

   **Find all errors in last hour:**
   ```apl
   ['becoming-diamond-prod']
   | where level == 'error'
   | where _time > ago(1h)
   | project _time, message, error, userId
   | sort by _time desc
   ```

   **Track specific user journey:**
   ```apl
   ['becoming-diamond-prod']
   | where userId == 'user_123'
   | where _time > ago(1d)
   | project _time, message, level
   | sort by _time asc
   ```

   **Payment success rate:**
   ```apl
   ['becoming-diamond-prod']
   | where eventType startswith 'checkout' or eventType startswith 'payment'
   | where _time > ago(24h)
   | summarize
       success = countif(eventType == 'checkout.session.completed'),
       failed = countif(eventType == 'payment_intent.payment_failed'),
       total = count()
   | extend success_rate = (success * 100.0) / total
   ```

3. **Save Queries**
   - Click **"Save Query"** to reuse later
   - Name it descriptively: "Last Hour Errors", "User Journey Tracking"

### Dashboards

1. **Create Dashboard**
   - Click **"Dashboards"** in sidebar
   - Click **"New Dashboard"**
   - Name: "Production Health" or "Business Metrics"

2. **Add Widgets**
   - Click **"Add Widget"**
   - Choose visualization type:
     - **Time Series**: For trends over time (request rate, error rate)
     - **Single Stat**: For current values (active users, total signups)
     - **Table**: For detailed records (recent errors, top users)
     - **Gauge**: For percentages (success rate, uptime)
     - **Pie Chart**: For distributions (sign-in methods, error types)

3. **Example Dashboard Widgets**

   **Widget 1: Request Rate (Time Series)**
   ```apl
   ['becoming-diamond-prod']
   | where _time > ago(1h)
   | summarize count() by bin(_time, 1m)
   ```

   **Widget 2: Error Rate (Gauge)**
   ```apl
   ['becoming-diamond-prod']
   | where _time > ago(1h)
   | summarize
       errors = countif(level == 'error'),
       total = count()
   | extend error_rate = (errors * 100.0) / total
   ```

   **Widget 3: Active Users (Single Stat)**
   ```apl
   ['becoming-diamond-prod']
   | where _time > ago(1h)
   | summarize dcount(userId)
   ```

### Alerts (Monitors)

1. **Create Monitor**
   - Click **"Monitors"** in sidebar
   - Click **"New Monitor"**

2. **Configure Alert**
   - **Name**: "Payment Failures Spike"
   - **Description**: "Alert when 2+ payment failures in 5 minutes"
   - **Query**:
     ```apl
     ['becoming-diamond-prod']
     | where eventType == 'payment_intent.payment_failed'
     | where _time > ago(5m)
     | count by customerEmail
     | where _count > 1
     ```
   - **Threshold**: Results > 0
   - **Frequency**: Check every 5 minutes

3. **Notification Channels**
   - **Email**: Add your email
   - **Slack**: Connect Slack workspace and select channel
   - **Webhook**: For custom integrations
   - **PagerDuty**: For on-call rotations

4. **Recommended Alerts** (from recommendations doc)
   - Payment failures (2+ in 5 minutes)
   - Email delivery failures (5+ in 15 minutes)
   - Authentication errors (3+ in 10 minutes)
   - Database connection failures (any)
   - Stripe webhook signature failures (3+ in 10 minutes)

---

## Step 7: Connect Slack (Optional but Recommended)

1. **Navigate to Integrations**
   - Go to Settings → Integrations
   - Find **"Slack"** integration

2. **Connect Slack Workspace**
   - Click **"Connect Slack"**
   - Authorize Axiom in your Slack workspace
   - Select which channels to use for alerts

3. **Create Alert Channels**
   - Recommended channels:
     - `#alerts-payments` - Payment and Stripe webhooks
     - `#alerts-auth` - Authentication failures
     - `#alerts-email` - Email delivery issues
     - `#alerts-database` - Database errors
     - `#alerts-general` - Catch-all for other errors

4. **Configure Monitors to Use Slack**
   - Edit each monitor
   - Add Slack channel as notification destination
   - Test the alert to verify it works

---

## Step 8: Monitor Usage and Costs

### Check Data Ingestion

1. **Go to Settings → Usage**
   - See daily ingestion volume
   - Monitor your usage against plan limits

2. **Free Tier Limits**
   - 500MB/day (approximately 1M small events)
   - 30-day retention
   - Good for initial testing

3. **Upgrade When Needed**
   - Hobby: $25/month for 50GB/month (~1.6GB/day)
   - Pro: $99/month for 250GB/month (~8GB/day)

### Optimize Costs

If approaching limits:

1. **Reduce Log Volume**
   - Filter out debug logs in production
   - Sample high-volume routes (log 10% of requests)
   - Remove verbose attributes from logs

2. **Adjust Retention**
   - Default: 30 days
   - Reduce to 7 days for lower tiers if needed
   - Export old logs before deletion

3. **Use Different Datasets**
   - High-priority logs: `becoming-diamond-critical` (longer retention)
   - Low-priority logs: `becoming-diamond-verbose` (shorter retention)

---

## Step 9: Test All Log Flows

Run through this checklist to ensure all logging is working:

### Authentication Logs
- [ ] User sign-in with magic link → Check for "User sign-in attempt" log
- [ ] User sign-in with Google OAuth → Check for provider="google" log
- [ ] Failed sign-in → Check for error log with reason
- [ ] New user creation → Check for "User created" log

### Payment Logs
- [ ] Successful checkout → Check for "checkout.session.completed" log with amount
- [ ] Failed payment → Check for "payment_intent.payment_failed" log
- [ ] Webhook signature failure → Check for error log (test with invalid signature)

### Email Logs
- [ ] Welcome email sent → Check for "Welcome email sent successfully" log
- [ ] Email retry → Check for "Scheduling email retry" log
- [ ] Email failure → Check for error log with SMTP error

### Profile Logs
- [ ] Profile fetch → Check for "User profile loaded" log
- [ ] Profile update → Check for PUT request log with fields updated

### Video Logs
- [ ] Video token request → Check for "Video token request" log
- [ ] Unauthorized access → Check for warning log with IP address

### Client-Side Logs
- [ ] Component error → Check for "Client component error" log with stack trace
- [ ] User context error → Check for profile load failure log

---

## Troubleshooting

### Logs Not Appearing

**Problem**: No logs showing up in Axiom

**Solutions**:
1. **Check Environment Variables**
   ```bash
   # In your app, log the config (remove in production!)
   console.log('AXIOM_TOKEN exists:', !!process.env.AXIOM_TOKEN)
   console.log('AXIOM_DATASET:', process.env.AXIOM_DATASET)
   ```

2. **Verify Token Permissions**
   - Go to Axiom Settings → API Tokens
   - Check token has "Ingest" permission
   - Check token is active (not expired)

3. **Check Network**
   - Open browser DevTools → Network tab
   - Look for requests to `*.axiom.co`
   - Check for 401 (auth error) or 403 (permission error) responses

4. **Restart Application**
   - Environment variables only load on app start
   - Restart dev server or redeploy production app

### Token Authentication Errors

**Error**: 401 Unauthorized or 403 Forbidden

**Solutions**:
1. **Regenerate Token**
   - Old token may be expired or revoked
   - Create new token in Axiom dashboard
   - Update environment variables

2. **Check Token Format**
   - Should start with `xaat-`
   - No extra spaces or quotes
   - Full token copied correctly

3. **Verify Dataset Name**
   - Must match exactly (case-sensitive)
   - Check for typos in `AXIOM_DATASET`

### Logs Delayed or Missing

**Problem**: Some logs appear, others don't

**Solutions**:
1. **Check Async Calls**
   - Ensure `await log.info()` not `log.info()` without await
   - Logs may be lost if process exits before flushing

2. **Buffering**
   - Axiom SDK buffers logs for efficiency
   - Logs may appear in batches (usually < 1 second delay)

3. **Check Log Levels**
   - Ensure not filtering out info/debug logs
   - Check Axiom UI filters not hiding logs

### High Data Usage

**Problem**: Approaching plan limits too quickly

**Solutions**:
1. **Identify High-Volume Routes**
   ```apl
   ['becoming-diamond-prod']
   | where _time > ago(1h)
   | summarize count() by route
   | sort by _count desc
   ```

2. **Sample Requests**
   - Log only 10% of successful requests
   - Always log errors
   - Reduce verbosity of attributes

3. **Use Conditional Logging**
   ```typescript
   // Only log in production for specific routes
   if (process.env.NODE_ENV === 'production' && isHighPriorityRoute) {
     await log.info('Request', { ... });
   }
   ```

---

## Quick Reference

### Essential URLs
- **Axiom Dashboard**: https://app.axiom.co
- **Documentation**: https://axiom.co/docs
- **API Tokens**: https://app.axiom.co/settings/api-tokens
- **Datasets**: https://app.axiom.co/datasets
- **Usage & Billing**: https://app.axiom.co/settings/billing

### Environment Variables
```bash
AXIOM_TOKEN=xaat-xxxxx-xxxx-xxxx-xxxx
AXIOM_DATASET=becoming-diamond-prod
```

### First Query to Try
```apl
['becoming-diamond-prod']
| where _time > ago(1h)
| limit 100
```

### Support
- **Documentation**: https://axiom.co/docs
- **Community Slack**: https://axiom.co/slack
- **Email Support**: support@axiom.co
- **Status Page**: https://status.axiom.co

---

## Next Steps After Setup

Once logs are flowing:

1. **Create Production Health Dashboard**
   - Request rate, error rate, response time
   - Active users, top routes

2. **Set Up Critical Alerts**
   - Payment failures
   - Authentication errors
   - Email delivery failures

3. **Review Logs Daily**
   - Check for anomalies
   - Monitor error trends
   - Identify optimization opportunities

4. **Document Common Queries**
   - Save frequently used queries
   - Share with team
   - Create runbooks for common issues

5. **Iterate on Logging**
   - Add more context to logs as needed
   - Remove noisy logs
   - Refine log levels

---

**Setup Complete!**

You now have production-grade observability for your Next.js application.
