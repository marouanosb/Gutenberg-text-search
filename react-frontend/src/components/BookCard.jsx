import React, { useState } from "react";
import BookModal from "./BookModal";

export default function BookCard({ book }) {
  const cover = book.cover_image || null;
  const author = book.authors && book.authors.length > 0 ? book.authors[0].name : "Unknown";
  const downloads = book.download_count || 0;
  const [open, setOpen] = useState(false);

  const openDetails = (e) => {
    e.preventDefault();
    setOpen(true);
  };

  return (
    <>
      <article
        className="card book-card interactive"
        onClick={openDetails}
        role="button"
        tabIndex={0}
        onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') setOpen(true); }}
        aria-label={`Open details for ${book.title}`}
        style={{position: 'relative'}}
      >
        {/* Downloads badge top-right */}
        <div className="downloads-badge" aria-hidden>
          ⬇ {downloads}
        </div>

        <div className="cover">
          {cover ? (
            // eslint-disable-next-line jsx-a11y/img-redundant-alt
            <img src={cover} alt={`Cover of ${book.title}`} loading="lazy" />
          ) : (
            <div className="no-cover">No Cover</div>
          )}
          <div className="cover-gradient" />
        </div>

        <div className="meta">
          <h3 className="title" title={book.title}>{book.title}</h3>
          <p className="author" title={author}>{author}</p>

          <div className="meta-footer">
            {book.plain_text && (
              <a
                href={book.plain_text}
                target="_blank"
                rel="noreferrer"
                className="btn small outline"
                onClick={(e) => e.stopPropagation()}
                title="Open book in new tab"
              >
                Read
              </a>
            )}
          </div>
        </div>
      </article>

      {open && <BookModal book={book} onClose={() => setOpen(false)} />}
    </>
  );
}
