Link to my application:
https://cvwo-chatit.onrender.com

AI declaration:

I used AI only in ways allowed as a learning and research aid:

To clarify concepts and best practices related to full-stack development.

To understand and compare technologies (e.g., SQLite vs PostgreSQL, Redux state management, and React component design).

At no point did I ask AI to generate or modify the actual code for the forum application. All coding decisions and implementations were completed independently, and I can fully explain and justify every part of the code.

Instructions on deploying the website locally:

clone the repo(please use version: dc8b1f6, the latest one was s desperate attempt to fix next.js issue) and:

cd to backend, after that run go mod tidy. There are a few environmental variables that need to be setup: 
PORT,
DATABASE_URL,
JWT_SECRET,
and need to:
modify the cookie setup in handler/auth, loginhandler to disable secure settings for cookie,
change the access control allow origin field in backend/middleware/auth.go enable cors to the corresponding localhost for frontend
(sorry for hardcoding those),
type go run main.go in terminal and enter.

cd to frontend, after that run npm install. Env var:
NEXT_PUBLIC_API_URL=http://localhost:8080(usually),
Then run npm run dev in terminal and you should be set!

Note that there are some bugs as I was unable to use useSearchParams, dk why. So you might(will) encounter issues when trying to search for posts or topics, basically if it involves changing url path query parameters. This can be fixed by going to another page (different url path) and search function should work(search posts when you are in topics and search topics when your in posts). This is only limited to search functions for topics and posts in app bar, search function under a post still works as per normal.
