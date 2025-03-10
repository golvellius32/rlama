import React from 'react';
import { Link } from 'react-router-dom';
import './Header.css';

function Header() {
  return (
    <header className="app-header">
      <div className="logo">
        <Link to="/">RLAMA</Link>
      </div>
      <nav>
        <Link to="/">Dashboard</Link>
        <Link to="/create">Create</Link>
      </nav>
    </header>
  );
}

export default Header;