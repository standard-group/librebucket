# Quick Start Guide

Get up and running with Librebucket in just a few minutes! This guide will walk you through creating your first repository and making your first commit.

## Step 1: Start Librebucket

After [installing Librebucket](installation.md), start the server:

```bash
./librebucket
```

You should see output similar to:
```
Working dir: /path/to/librebucket
DB initialized at: /path/to/librebucket/config/data/users.db
Server starting on :3000
```

## Step 2: Create a User Account

First, register a new user account using the API:

```bash
curl -X POST http://localhost:3000/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "secure_password123"
  }'
```

**Response:**
```json
{
  "status": "success",
  "token": "your-auth-token-here",
  "user": {
    "id": 1,
    "username": "alice",
    "is_admin": false
  }
}
```

!!! important "Save Your Token"
    Save the returned token - you'll need it for API requests!

## Step 3: Create Your First Repository

Create a new repository using your authentication token:

```bash
curl -X POST http://localhost:3000/api/v1/git/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-auth-token-here" \
  -d '{
    "username": "alice",
    "reponame": "my-first-repo",
    "public": false
  }'
```

## Step 4: Clone and Use Your Repository

Now you can clone your new repository:

```bash
git clone http://localhost:3000/alice/my-first-repo.git
cd my-first-repo
```

## Step 5: Make Your First Commit

Add some content and push to your repository:

```bash
# Create a README
echo "# My First Librebucket Repository" > README.md

# Add and commit
git add README.md
git commit -m "Initial commit: Add README"

# Push to your Librebucket server
git push -u origin main
```

## Step 6: Explore the Web Interface

Open your browser and navigate to `http://localhost:3000` to:

- Browse your repositories
- View commit history
- Explore file contents
- Manage your account

## What's Next?

Now that you have Librebucket running, explore these features:

### Repository Management
- [Create and manage repositories](../user-guide/repositories.md)
- [Browse code and history](../user-guide/web-interface.md)
- [Set up Git hooks](../user-guide/git-hooks.md)

### User Management
- [Create additional users](../user-guide/users.md)
- [Set up admin accounts](../user-guide/administration.md)
- [Configure permissions](../user-guide/permissions.md)

### API Integration
- [Explore the full API](../api/overview.md)
- [Automate repository creation](../api/repositories.md)
- [Integrate with CI/CD](../deployment/ci-cd.md)

### Production Deployment
- [Configure for production](../deployment/self-hosting.md)
- [Set up Docker deployment](../deployment/docker.md)
- [Configure reverse proxy](../deployment/reverse-proxy.md)

## Example Workflows

### Development Team Setup

```bash
# Admin creates team repositories
curl -X POST http://localhost:3000/api/v1/git/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer admin-token" \
  -d '{
    "username": "team",
    "reponame": "backend-api",
    "public": false
  }'

curl -X POST http://localhost:3000/api/v1/git/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer admin-token" \
  -d '{
    "username": "team",
    "reponame": "frontend-app",
    "public": false
  }'
```

### Personal Projects

```bash
# Create multiple personal repositories
for repo in "dotfiles" "scripts" "notes"; do
  curl -X POST http://localhost:3000/api/v1/git/create \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer your-token" \
    -d "{
      \"username\": \"alice\",
      \"reponame\": \"$repo\",
      \"public\": false
    }"
done
```

## Troubleshooting

### Can't Connect to Server
- Ensure Librebucket is running: `ps aux | grep librebucket`
- Check the port: `netstat -ln | grep 3000`
- Try a different port: `LIBREBUCKET_PORT=8080 ./librebucket`

### Authentication Errors
- Verify your token is correct
- Check token format: `Authorization: Bearer <token>`
- Try using the X-Auth-Token header instead

### Repository Creation Fails
- Ensure the username exists and matches your account
- Check repository name doesn't already exist
- Verify you have proper permissions

### Git Clone/Push Issues
- Ensure the repository was created successfully
- Check the clone URL format: `http://localhost:3000/username/reponame.git`
- Verify network connectivity to the server

!!! tip "Need More Help?"
    - Check the [full documentation](../user-guide/web-interface.md)
    - Visit our [troubleshooting guide](../about/troubleshooting.md)
    - Create an [issue on GitHub](https://github.com/standard-group/librebucket/issues)
