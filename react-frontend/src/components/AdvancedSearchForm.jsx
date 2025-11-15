import React, { useState } from "react";

export default function AdvancedSearchForm({ urlBase, onSearch, setLoading }) {
  const [keyword, setKeyword] = useState("");
  const [keywordType, setKeywordType] = useState("classique");
  const [title, setTitle] = useState("");
  const [titleType, setTitleType] = useState("classique");
  const [author, setAuthor] = useState("");
  const [authorType, setAuthorType] = useState("classique");
  const [language, setLanguage] = useState("en");
  const [sort, setSort] = useState("download_count");
  const [order, setOrder] = useState("ascending");

  const buildUrl = () => {
    let url = `${urlBase}sort=${sort}&order=${order}`;
    if (language !== "all") url += `&languages=${language}`;
    if (author) url += `&author_name=${encodeURIComponent(author)}&author_name_type=${authorType}`;
    if (title) url += `&title=${encodeURIComponent(title)}&title_name_type=${titleType}`;
    if (keyword) url += `&keyword=${encodeURIComponent(keyword)}&keyword_type=${keywordType}`;
    return url;
  };

  const submit = (e) => {
    e?.preventDefault();
    const url = buildUrl();
    onSearch(url, { keyword, title, author });
  };

  return (
    <form className="card form" onSubmit={submit}>
      <h3>Advanced Search</h3>

      <div className="col">
        <label>Keyword</label>
        <div className="row small">
          <input placeholder="Keyword" value={keyword} onChange={(e) => setKeyword(e.target.value)} />
          <select value={keywordType} onChange={(e) => setKeywordType(e.target.value)}>
            <option value="classique">Classique</option>
            <option value="regex">Regex</option>
          </select>
        </div>
      </div>

      <div className="col">
        <label>Title</label>
        <div className="row small">
          <input placeholder="Title" value={title} onChange={(e) => setTitle(e.target.value)} />
          <select value={titleType} onChange={(e) => setTitleType(e.target.value)}>
            <option value="classique">Classique</option>
            <option value="regex">Regex</option>
          </select>
        </div>
      </div>

      <div className="col">
        <label>Author</label>
        <div className="row small">
          <input placeholder="Author" value={author} onChange={(e) => setAuthor(e.target.value)} />
          <select value={authorType} onChange={(e) => setAuthorType(e.target.value)}>
            <option value="classique">Classique</option>
            <option value="regex">Regex</option>
          </select>
        </div>
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
        <button className="btn primary" type="submit">Search</button>
        <button type="button" className="btn" onClick={() => { setKeyword(""); setTitle(""); setAuthor(""); }}>Clear</button>
      </div>
    </form>
  );
}
