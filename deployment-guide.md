# 🚀 Production Deployment Guide

## Option 1: Deploy ke VPS (Digital Ocean/AWS/GCP)

### 1. Setup Server
```bash
# SSH ke server
ssh root@your-server-ip

# Install Docker & Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo curl -L "https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2. Upload Code
```bash
# Upload kode ke server
scp -r . root@your-server-ip:/var/www/shorter/
```

### 3. Production Environment
```bash
# Buat file .env production
cat > .env << 'EOF'
# Database
DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=your-strong-password
DB_NAME=shorter_db
DB_PORT=5432
DB_SSLMODE=disable

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# Server
PORT=8080
GIN_MODE=release
CUSTOM_DOMAIN=yourdomain.com

# Frontend
NEXT_PUBLIC_API_URL=https://yourdomain.com/api
EOF
```

### 4. SSL Certificate dengan Nginx
```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Setup Nginx reverse proxy
sudo nano /etc/nginx/sites-available/shorter
```

## Option 2: Deploy ke Heroku

### 1. Setup Heroku
```bash
# Install Heroku CLI
npm install -g heroku

# Login dan create app
heroku login
heroku create your-app-name
```

### 2. Add-ons
```bash
# PostgreSQL
heroku addons:create heroku-postgresql:mini

# Redis
heroku addons:create heroku-redis:mini
```

## Option 3: Deploy ke Railway/Vercel

### Railway (Recommended untuk Full-Stack)
1. Connect GitHub repository
2. Auto-deploy backend + database
3. Environment variables otomatis

### Vercel (Frontend Only)
1. Deploy frontend ke Vercel
2. Backend tetap di VPS/Railway

---

## 📊 Monitoring & Analytics

### Setup Monitoring
```bash
# Add monitoring tools
docker run -d \
  --name prometheus \
  -p 9090:9090 \
  prom/prometheus

# Grafana for dashboards  
docker run -d \
  --name grafana \
  -p 3001:3000 \
  grafana/grafana
```

---

## 🔒 Security Checklist

- [ ] SSL Certificate (HTTPS)
- [ ] Rate limiting
- [ ] Input validation
- [ ] CORS configuration
- [ ] Environment variables
- [ ] Database security
- [ ] Regular backups 