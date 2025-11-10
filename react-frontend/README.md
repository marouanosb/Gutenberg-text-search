# Gutenberg React Frontend

This is a minimal React frontend (Vite) that replicates the main search features of the existing Next.js frontend in this repository.

Features
- Simple and advanced search forms (keyword/title/author + sort + language)
- Calls backend at `http://localhost:8000/server/books/?...`
- Displays results and suggestions with cover, author and read link
- Modern responsive design (custom CSS)

Quick start
1. From the project root, open a terminal and go to the new folder:

```cmd
cd react-frontend
```

2. Install dependencies and run the dev server:

```cmd
npm install
npm run dev
```

3. Open the app at the URL printed by Vite (usually http://localhost:5173).

Notes
- The app expects the backend API to be available at `http://localhost:8000/server/books/` as in the original Next.js app.
- If your backend uses a different base URL or port, update `src/App.jsx` (the `URL_BASE` constant) accordingly.

Next steps (optional)
- Add TypeScript or Tailwind for nicer dev DX
- Add pagination and error UI
- Add tests

