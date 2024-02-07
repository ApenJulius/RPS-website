import { useState, useRef, useEffect} from 'react';
import {useNavigate} from 'react-router-dom';
import './LobbyPage.css';
import { MoveButton } from '../../components/MoveButton/MoveButton';



function LobbyPage() {
  const navigate = useNavigate();
    
    const startGame = () => {
        console.log("start game")
        navigate("/" + genRandomString(10))
    
    }
    const genRandomString = (length: number) =>{
      var chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
      var charLength = chars.length;
      var result = '';
      for ( var i = 0; i < length; i++ ) {
         result += chars.charAt(Math.floor(Math.random() * charLength));
      }
      return result;
    }
    

  return (
    <div>
      <h1>Lobby</h1>
      <button onClick={startGame}>Start Game</button>
      </div>
  );
}

export default LobbyPage;
