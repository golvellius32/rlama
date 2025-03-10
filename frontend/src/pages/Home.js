import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { listRags, deleteRag } from '../services/api';
import './Home.css';

function Home() {
  const [rags, setRags] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  const loadRags = async () => {
    setIsLoading(true);
    try {
      const ragList = await listRags();
      setRags(ragList);
      setError(null);
    } catch (err) {
      setError(`Error loading RAGs: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadRags();
  }, []);

  const handleDelete = async (name) => {
    if (window.confirm(`Are you sure you want to delete "${name}"?`)) {
      try {
        await deleteRag(name);
        loadRags();
      } catch (err) {
        setError(`Error deleting RAG: ${err.message}`);
      }
    }
  };

  return (
    <div className="home-container">
      <div className="header-actions">
        <h2>Available RAG Systems</h2>
        <Link to="/create" className="create-button">
          Create New RAG
        </Link>
      </div>

      {error && <div className="error-message">{error}</div>}

      {isLoading ? (
        <div className="loading">Loading RAG systems...</div>
      ) : rags.length === 0 ? (
        <div className="empty-state">
          <p>No RAG systems found.</p>
          <p>Get started by creating your first RAG system.</p>
          <Link to="/create" className="create-button">
            Create RAG
          </Link>
        </div>
      ) : (
        <div className="rag-list">
          {rags.map((rag) => (
            <div key={rag.name} className="rag-card">
              <div className="rag-info">
                <h3>{rag.name}</h3>
                <p>Model: {rag.model_name}</p>
                <p>Documents: {rag.documents.length}</p>
                <p className="rag-date">Created: {new Date(rag.created_at).toLocaleString()}</p>
              </div>
              <div className="rag-actions">
                <Link to={`/chat/${rag.name}`} className="chat-button">
                  Chat
                </Link>
                <button
                  className="delete-button"
                  onClick={() => handleDelete(rag.name)}
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default Home;