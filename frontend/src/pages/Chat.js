import React, { useState, useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';
import { queryRag, getRag } from '../services/api';
import './Chat.css';

function Chat() {
  const { ragName } = useParams();
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [ragInfo, setRagInfo] = useState(null);
  const messagesEndRef = useRef(null);

  // Fetch RAG information
  useEffect(() => {
    const loadRagInfo = async () => {
      try {
        const info = await getRag(ragName);
        setRagInfo(info);
        // Add welcome message
        setMessages([
          {
            text: `Welcome! You can ask questions about documents in the "${ragName}" collection.`,
            sender: 'system'
          }
        ]);
      } catch (error) {
        setMessages([
          {
            text: `Error loading RAG "${ragName}": ${error.message}`,
            sender: 'error'
          }
        ]);
      }
    };

    loadRagInfo();
  }, [ragName]);

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSendMessage = async (e) => {
    e.preventDefault();
    if (!input.trim() || isLoading) return;

    const userMessage = input;
    setInput('');
    setMessages(prev => [...prev, { text: userMessage, sender: 'user' }]);
    setIsLoading(true);

    try {
      const response = await queryRag(ragName, userMessage);
      setMessages(prev => [...prev, { text: response, sender: 'assistant' }]);
    } catch (error) {
      setMessages(prev => [
        ...prev, 
        { text: `Error: ${error.message}`, sender: 'error' }
      ]);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="chat-container">
      <div className="chat-header">
        <h2>{ragName}</h2>
        {ragInfo && (
          <div className="rag-info">
            <span>Model: {ragInfo.model_name}</span>
            <span>Documents: {ragInfo.document_count}</span>
          </div>
        )}
      </div>

      <div className="chat-messages">
        {messages.map((msg, index) => (
          <div key={index} className={`message ${msg.sender}`}>
            <div className="message-content">{msg.text}</div>
          </div>
        ))}
        {isLoading && (
          <div className="message assistant loading">
            <div className="message-content">
              <div className="typing-indicator">
                <span></span>
                <span></span>
                <span></span>
              </div>
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      <form className="chat-input-form" onSubmit={handleSendMessage}>
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Ask a question..."
          disabled={isLoading}
        />
        <button type="submit" disabled={isLoading || !input.trim()}>
          Send
        </button>
      </form>
    </div>
  );
}

export default Chat;