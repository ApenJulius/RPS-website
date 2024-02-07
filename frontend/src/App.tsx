import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import logo from './logo.svg';
import './App.css';
import GamePage from './pages/gamePage/gamePage'; // Make sure to import your GamePage component
import LobbyPage from './pages/LobbyPage/LobbyPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LobbyPage />} />

        <Route path="/:id" element={<GamePage />} />
      </Routes>
    </Router>
  );
}

export default App;