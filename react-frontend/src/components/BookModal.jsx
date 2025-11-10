import React from "react";

export default function BookModal({ book, onClose }) {
  const authors = book.authors || [];

  return (
    <div className="modal-overlay" role="dialog" aria-modal="true" aria-label={`Details for ${book.title}`} onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onClose} aria-label="Close details">×</button>

        <div className="modal-body">
          <div className="modal-cover">
            {book.cover_image ? (
              <img src={book.cover_image} alt={`Cover of ${book.title}`} />
            ) : (
              <div className="no-cover large">No Cover</div>
            )}
          </div>

          <div className="modal-info">
            <h2 className="modal-title">{book.title}</h2>
            {authors.length > 0 && (
              <div className="modal-authors">
                <strong>Authors:</strong>
                <ul>
                  {authors.map((a, i) => (
                    <li key={i}>{a.name} {a.birth_year ? `(${a.birth_year}${a.death_year ? ' - ' + a.death_year : ''})` : ''}</li>
                  ))}
                </ul>
              </div>
            )}

            {book.subjects && book.subjects.length > 0 && (
              <div className="modal-section">
                <strong>Subjects:</strong>
                <p>{book.subjects.join(', ')}</p>
              </div>
            )}

            {book.languages && book.languages.length > 0 && (
              <div className="modal-section">
                <strong>Languages:</strong>
                <p>{book.languages.map((l) => l.code || l).join(', ')}</p>
              </div>
            )}

            <div className="modal-section">
              <strong>Downloads:</strong>
              <p>{book.download_count || 0}</p>
            </div>

            {book.plain_text && (
              <div className="modal-actions">
                <a href={book.plain_text} target="_blank" rel="noreferrer" className="btn primary">Read the book</a>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
