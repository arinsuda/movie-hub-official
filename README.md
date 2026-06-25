# Movie Hub

A full-stack web application for browsing and reviewing movies, built with Vue 3 on the frontend and Go on the backend.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Vue 3 + TypeScript |
| Backend | Go |
| Database | PostgreSQL 16 |
| Object Storage | MinIO |
| Containerization | Docker / Docker Compose |

---

## Project Structure

```
movie-hub-official/
├── client/          # Vue 3 + TypeScript frontend
├── server/          # Go backend API
└── docker-compose.yml
```

---

## Prerequisites

- [Docker](https://www.docker.com/) & Docker Compose
- [Node.js](https://nodejs.org/) (for local frontend development)
- [Go](https://go.dev/) 1.21+ (for local backend development)

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/arinsuda/movie-hub-official.git
cd movie-hub-official
```

### 2. Set up environment variables

สร้างไฟล์ `.env` ที่ root ของ project:

```env
# PostgreSQL
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=moviereview
POSTGRES_PORT=5432

# MinIO
MINIO_ROOT_USER=your_minio_user
MINIO_ROOT_PASSWORD=your_minio_password
MINIO_PORT=9000
MINIO_CONSOLE_PORT=9001
```

### 3. Start infrastructure services

```bash
docker compose up -d
```

Services ที่จะรันขึ้นมา:

- **PostgreSQL** → `localhost:5432`
- **MinIO** → `localhost:9000` (API), `localhost:9001` (Console)

### 4. Run the backend

```bash
cd server
go mod tidy
go run .
```

### 5. Run the frontend

```bash
cd client
npm install
npm run dev
```

แอปพลิเคชันจะพร้อมใช้งานที่ `http://localhost:5173`

---

## Services

### MinIO Console

เข้าถึง MinIO web console สำหรับจัดการ object storage ได้ที่ `http://localhost:9001`  
ใช้ `MINIO_ROOT_USER` และ `MINIO_ROOT_PASSWORD` ที่กำหนดไว้ใน `.env`

### PostgreSQL

ต่อเชื่อมกับ database ได้ที่ `localhost:5432` ด้วย credentials ใน `.env`

---

## Development

### Frontend

```bash
cd client
npm run dev      # Development server
npm run build    # Production build
npm run preview  # Preview production build
```

### Backend

```bash
cd server
go run .         # Run server
go build .       # Build binary
go test ./...    # Run tests
```

---

## License

This project is open source. See [LICENSE](LICENSE) for details.
