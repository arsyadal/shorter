# 🚀 Deploy Full-Stack ke Vercel

## Opsi 1: Full-Stack Vercel (Recommended)

### 1. Restructure untuk Vercel
```bash
# Pindah frontend ke root dan backend ke API routes
mkdir api
cp -r backend/* api/
rm -rf backend

# Update struktur:
# /pages atau /app (frontend)
# /api (backend serverless functions)
```

### 2. Convert Go Backend ke Vercel Functions
```javascript
// api/shorten.js - Convert dari Go ke JavaScript/TypeScript
export default async function handler(req, res) {
  if (req.method === 'POST') {
    // Logic URL shortening
    // Connect ke PostgreSQL
    // Return JSON response
  }
}
```

### 3. Environment Variables di Vercel
```bash
# Via Vercel Dashboard atau CLI
vercel env add DATABASE_URL
vercel env add REDIS_URL
vercel env add CUSTOM_DOMAIN
```

---

## Opsi 2: Hybrid Deployment (Vercel + Backend Terpisah)

### Frontend di Vercel
```bash
cd frontend
npx vercel --prod
# Auto deploy dengan domain: your-app.vercel.app
```

### Backend di Platform Lain
- **Railway**: Deploy Go backend
- **Heroku**: Deploy Go backend  
- **VPS**: Deploy dengan Docker

### Update Frontend Config
```javascript
// frontend/src/utils/api.ts
const API_URL = process.env.NODE_ENV === 'production' 
  ? 'https://your-backend.railway.app/api'
  : 'http://localhost:8080/api'
```

---

## Opsi 3: Full Vercel dengan Next.js API Routes

### Struktur Project Baru
```
shorter-vercel/
├── pages/
│   ├── api/
│   │   ├── shorten.js
│   │   ├── [code].js (redirect)
│   │   └── stats/[code].js
│   ├── index.js
│   └── qr/[code].js
├── lib/
│   ├── database.js
│   └── redis.js
└── vercel.json
```

### Database Setup
```bash
# Vercel Postgres
vercel storage create postgres --name url-shortener-db

# Atau External
# Supabase, PlanetScale, Railway Postgres
```

---

## ⚡ PERBANDINGAN PLATFORM:

| Platform | Frontend | Backend | Database | Setup | Cost |
|----------|----------|---------|----------|-------|------|
| **Vercel** | ✅ Excellent | ✅ Serverless | ✅ Built-in | 🟢 Easy | 💰 Generous Free |
| **Railway** | ✅ Good | ✅ Full Go | ✅ PostgreSQL | 🟢 Easy | 💰 $5/mo |
| **Heroku** | ✅ Good | ✅ Full Go | ✅ PostgreSQL | 🟡 Medium | 💰 $7/mo |
| **VPS** | ✅ Good | ✅ Full Control | ✅ Any DB | 🔴 Complex | 💰 $5-20/mo |

---

## 🎯 RECOMMENDATION:

### Untuk Beginners: **Vercel Full-Stack**
- Paling mudah setup
- Auto SSL & CDN
- Generous free tier
- Integrated monitoring

### Untuk Production: **Railway**
- Keep existing Go code
- Better for complex backend logic
- Predictable pricing
- Easy database management

### Untuk Enterprise: **VPS/AWS**
- Full control
- Custom configurations
- Better for high traffic
- More complex setup 