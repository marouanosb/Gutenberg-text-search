import React, { useState } from "react";

export default function SearchForm({ urlBase, onSearch, setLoading }) {
  const [keyword, setKeyword] = useState("");
  const [title, setTitle] = useState("");
  const [author, setAuthor] = useState("");
  const [language, setLanguage] = useState("en");
  const [sort, setSort] = useState("download_count");
  const [order, setOrder] = useState("ascending");

  const buildUrl = () => {
    let url = `${urlBase}sort=${sort}&order=${order}`;
    if (language !== "all") url += `&languages=${language}`;
    if (author) url += `&author_name=${encodeURIComponent(author)}`;
    if (title) url += `&title=${encodeURIComponent(title)}`;
    if (keyword) url += `&keyword=${encodeURIComponent(keyword)}`;
    return url;
  };

  const submit = (e) => {
    e?.preventDefault();
    const url = buildUrl();
    onSearch(url, { keyword, title, author });
  };

  return (
    <form className="card form" onSubmit={submit}>
      <h3>Simple Search</h3>
      <div className="row">
        <input placeholder="Keyword" value={keyword} onChange={(e) => setKeyword(e.target.value)} />
        <input placeholder="Title" value={title} onChange={(e) => setTitle(e.target.value)} />
        <input placeholder="Author" value={author} onChange={(e) => setAuthor(e.target.value)} />
      </div>

      <div className="row small">
        <select value={language} onChange={(e) => setLanguage(e.target.value)}>
          <option value="en">English</option>
          <option value="fr">French</option>
          <option value="all">All</option>
        </select>

        <select value={sort} onChange={(e) => setSort(e.target.value)}>
          <option value="download_count">Download Count</option>
          <option value="closeness">Closeness</option>
          <option value="betweenness">Betweenness</option>
        </select>

        <select value={order} onChange={(e) => setOrder(e.target.value)}>
          <option value="ascending">Ascending</option>
          <option value="descending">Descending</option>
        </select>
      </div>

      <div className="row actions">
        <button type="submit" className="btn primary">Search</button>
        <button type="button" className="btn" onClick={() => { setKeyword(""); setTitle(""); setAuthor(""); }}>Clear</button>
      </div>
    </form>
  );
}
