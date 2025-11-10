import React from "react";
import BookCard from "./BookCard";

export default function Results({ booksData }) {
  const { result = [], suggestions = [], searchParams, hasSearched } = booksData || {};

  const terms = [];
  if (searchParams?.keyword) terms.push(`keyword "${searchParams.keyword}"`);
  if (searchParams?.title) terms.push(`title "${searchParams.title}"`);
  if (searchParams?.author) terms.push(`author "${searchParams.author}"`);

  return (
    <section className="results">
      {result && result.length > 0 && (
        <div>
          <h2>Results ({result.length})</h2>
          {terms.length > 0 && <p className="muted">Related to {terms.join(", ")}</p>}
          <div className="grid">
            {result.map((b, i) => <BookCard key={i} book={b} />)}
          </div>
        </div>
      )}

      {suggestions && suggestions.length > 0 && (
        <div className="mt">
          <h2>You Might Also Like ({suggestions.length})</h2>
          <div className="grid">
            {suggestions.map((b, i) => <BookCard key={i} book={b} />)}
          </div>
        </div>
      )}

      {hasSearched && result.length === 0 && suggestions.length === 0 && (
        <div className="empty">No books found for your query.</div>
      )}

      {!hasSearched && (
        <div className="empty">Use the form above to search books.</div>
      )}
    </section>
  );
}
