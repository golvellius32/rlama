import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createRag } from '../services/api';
import './Create.css';

function Create() {
  const [modelName, setModelName] = useState('llama3');
  const [ragName, setRagName] = useState('');
  const [files, setFiles] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const handleFileChange = (e) => {
    setFiles(Array.from(e.target.files));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!ragName.trim()) {
      setError('Please enter a RAG name');
      return;
    }

    if (files.length === 0) {
      setError('Please select at least one file');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const formData = new FormData();
      formData.append('modelName', modelName);
      formData.append('ragName', ragName);
      files.forEach(file => formData.append('files', file));

      await createRag(formData);
      navigate('/');
    } catch (err) {
      setError(`Error creating RAG: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="create-container">
      <h2>Create a New RAG System</h2>
      {error && <div className="error-message">{error}</div>}
      
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="ragName">RAG Name:</label>
          <input
            id="ragName"
            type="text"
            value={ragName}
            onChange={(e) => setRagName(e.target.value)}
            placeholder="Enter a name for your RAG"
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="modelName">Model:</label>
          <select
            id="modelName"
            value={modelName}
            onChange={(e) => setModelName(e.target.value)}
          >
            <option value="llama3">llama3</option>
            <option value="mistral">mistral</option>
            <option value="gemma">gemma</option>
          </select>
        </div>

        <div className="form-group">
          <label htmlFor="files">Documents:</label>
          <input
            id="files"
            type="file"
            multiple
            onChange={handleFileChange}
            required
          />
          <div className="file-list">
            {files.map((file, index) => (
              <div key={index} className="file-item">
                {file.name}
              </div>
            ))}
          </div>
        </div>

        <button 
          type="submit" 
          disabled={isLoading} 
          className="submit-button"
        >
          {isLoading ? 'Creating...' : 'Create RAG'}
        </button>
      </form>
    </div>
  );
}

export default Create;