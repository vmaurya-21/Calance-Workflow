# Database Setup Quick Start

## Problem
If you see the error `password authentication failed for user "postgres"`, your `.env` file password doesn't match docker-compose.yml.

## Solution

Update your `.env` file to match the docker-compose.yml defaults:

```env
DB_PASSWORD=postgres
```

OR update docker-compose.yml to match your desired password:

```yaml
environment:
  POSTGRES_PASSWORD: your_desired_password
```

## Recommended Setup

For development, use the default credentials from docker-compose.yml:

1. Copy `.env.example` to `.env`:
   ```bash
   copy .env.example .env   # Windows
   cp .env.example .env     # Linux/Mac
   ```

2. The `.env.example` is already configured with the correct defaults matching docker-compose.yml

3. Start your containers:
   ```bash
   docker-compose up -d
   ```

4. Run your application:
   ```bash
   go run cmd/server/main.go
   ```

## Current docker-compose.yml Settings

```yaml
POSTGRES_USER: postgres
POSTGRES_PASSWORD: postgres
POSTGRES_DB: calance_workflow
```

Make sure your `.env` file has:
```env
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=calance_workflow
```
