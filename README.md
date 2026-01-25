Link to my application:
https://cvwo-chatit.onrender.com

AI declaration:

I used AI only in ways allowed as a learning and research aid:

To clarify concepts and best practices related to full-stack development.

To understand and compare technologies (e.g., SQLite vs PostgreSQL, Redux state management, and React component design).

At no point did I ask AI to generate or modify the actual code for the forum application. All coding decisions and implementations were completed independently, and I can fully explain and justify every part of the code.


Instructions on deploying the website locally:

Local Deployment Instructions

These instructions guide you through setting up the website locally using version dc8b1f6. (Note: The latest version contains experimental fixes for Next.js issues and may not be stable.)

1. Clone the Repository
git clone <repo-url>
git checkout dc8b1f6

2. Backend Setup

Navigate to the backend folder:

cd backend


Install dependencies:

go mod tidy


Set the required environment variables:

PORT — port for the backend server (e.g., 8080)

DATABASE_URL — URL for your database

JWT_SECRET — secret key for JWT authentication

Modify the following backend files for local development:

Cookie settings:

File: handler/auth.go → loginHandler

Disable secure cookie settings for local testing.

CORS settings:

File: middleware/auth.go

Update Access-Control-Allow-Origin to allow requests from your frontend (usually http://localhost:3000).

Start the backend server:

go run main.go

3. Frontend Setup

Navigate to the frontend folder:

cd frontend


Install dependencies:

npm install


Set the environment variable:

NEXT_PUBLIC_API_URL=http://localhost:8080


Start the frontend server:

npm run dev


Your frontend should now be accessible at http://localhost:3000 (or the port shown in your terminal).

4. Known Issues

Search functionality:
The useSearchParams hook may not work correctly in some places. Specifically:

Searching posts while in the Topics page, or searching topics while in the Posts page, may fail.

Workaround: Navigate to a different page, then try the search again.

Search under a single post still works as expected.
