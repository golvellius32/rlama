const API_BASE_URL = 'http://localhost:3001/api';

export const createRag = async (formData) => {
  const response = await fetch(`${API_BASE_URL}/rag`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to create RAG system');
  }

  return await response.json();
};

export const listRags = async () => {
  const response = await fetch(`${API_BASE_URL}/rag`);
  
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to fetch RAG systems');
  }
  
  return await response.json();
};

export const getRag = async (name) => {
  const response = await fetch(`${API_BASE_URL}/rag/${name}`);
  
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to fetch RAG details');
  }
  
  return await response.json();
};

export const deleteRag = async (name) => {
  const response = await fetch(`${API_BASE_URL}/rag/${name}`, {
    method: 'DELETE',
  });
  
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to delete RAG system');
  }
  
  return true;
};

export const queryRag = async (name, query) => {
  const response = await fetch(`${API_BASE_URL}/query/${name}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ query }),
  });
  
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to query RAG system');
  }
  
  const data = await response.json();
  return data.response;
};