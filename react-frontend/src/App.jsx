import React, { useState } from "react";
import SearchForm from "./components/SearchForm";
import AdvancedSearchForm from "./components/AdvancedSearchForm";
import Results from "./components/Results";
import Loader from "./components/Loader";

const URL_BASE = "http://localhost:8000/server/books/?";

export default function App() {
  const [tab, setTab] = useState("simple");
  const [booksData, setBooksData] = useState({ result: [], suggestions: [], searchParams: null, hasSearched: false });
  const [loading, setLoading] = useState(false);

  const handleSearch = async (queryUrl, searchParams) => {
    setLoading(true);
    try {
      const res = await fetch(queryUrl, { mode: "cors" });
      const data = await res.json();
      setBooksData({ ...data, searchParams, hasSearched: true });
    } catch (err) {
      console.error("Fetch error:", err);
      setBooksData({ result: [], suggestions: [], searchParams, hasSearched: true });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app-root">
      <header className="hero">
        <h1>Gutenberg Search</h1>
        <p className="subtitle">A modern React frontend — search books by keyword, title or author</p>
      </header>

      <main className="container">
        <div className="tabs">
          <button className={tab === "simple" ? "active" : ""} onClick={() => setTab("simple")}>Simple</button>
          <button className={tab === "advanced" ? "active" : ""} onClick={() => setTab("advanced")}>Advanced</button>
        </div>

        <div className="form-area">
          {tab === "simple" ? (
            <SearchForm urlBase={URL_BASE} onSearch={handleSearch} setLoading={setLoading} />
          ) : (
            <AdvancedSearchForm urlBase={URL_BASE} onSearch={handleSearch} setLoading={setLoading} />
          )}
        </div>

        {loading ? <Loader /> : <Results booksData={booksData} />}
      </main>

      <footer className="footer">
        <small>React frontend — connects to backend at <code>http://localhost:8000</code></small>
      </footer>
    </div>
  );
}
