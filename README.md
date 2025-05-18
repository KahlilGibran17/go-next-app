# go-next-app
This is the take home test created using by go as backend and next.js as frontend

# Go + Frontend Project with Serverless API on Vercel

This repository contains a frontend project integrated with a Go backend API deployed as serverless functions on Vercel.  
The backend handles user login with JWT authentication and CRUD operations for items. The frontend consumes these APIs.

---


---

## Features

- Go serverless API for:
  - User login with JWT (`/api/login`)
  - Item CRUD (`/api/items`)
- Frontend consumes these APIs
- Easy local development with `vercel dev`
- Deployment-ready for Vercel hosting

---

## Prerequisites

- Go 1.20+ ([Install Go](https://go.dev/dl/))
- Node.js & npm ([Install Node.js](https://nodejs.org/))
- Vercel CLI ([Install Vercel CLI](https://vercel.com/docs/cli))

---

## Setup & Installation

1. Clone the repo:
2 .Download Go module dependencies:
   go mod tidy
3. (Optional) Install Vercel CLI globally:
   npm install -g vercel

Running Locally
Backend (Go API)
From the frontend/api directory, run:
vercel dev
This starts the backend serverless API at http://localhost:3000/api/login, http://localhost:3000/api/items, etc.

Frontend
From the frontend folder, run your frontend framework command (e.g., Next.js):
npm install
npm run dev
The frontend will be available at http://localhost:3000 (or specified port).

Deployment

Backend
From the frontend/api folder:
vercel
Follow prompts to link or create your Vercel project.

Deploy production:
vercel --prod

Backend API will be live at:
https://your-vercel-project.vercel.app/api/

Frontend
From the root frontend folder, deploy your frontend separately:
vercel
vercel --prod

API Endpoints
POST /api/login
Authenticate user and get JWT token.

Request body:
  "email": "user@example.com",
  "password": "password"
  
Response:
{
  "token": "jwt.token.here"
}


