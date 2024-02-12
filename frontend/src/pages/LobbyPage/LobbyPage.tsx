import { useState, useRef, useEffect} from 'react';
import {useNavigate} from 'react-router-dom';
import './LobbyPage.css';
import { MoveButton } from '../../components/MoveButton/MoveButton';
import { ErrorCode } from '../../constants/ErrorCodes';
const { LOBBIES, PLAYER_JOINED, PLAYER_LEFT, GAME_COUNTDOWN } = ErrorCode


function LobbyPage() {
  const navigate = useNavigate();
  const ws = useRef<WebSocket | null>(null)
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


    const getLobbies = () => {

    }


    useEffect(() => {
      ws.current = new WebSocket(`ws://${process.env.REACT_APP_WEBSITE_NAME}:8000/lobby`);
      
      ws.current.onopen = () => {
        console.log('ws opened');
        getLobbies();
      };
  
      ws.current.onmessage = (message) => {
        const data = JSON.parse(message.data);
        console.log(data);
        switch(Number(data.code)) {
          case LOBBIES:
            break;
          default:
            console.error("UNKNOWN ERROR", data.code)
        }
        
      }
  
      ws.current.onclose = () => {
        console.log('ws closed');
      }
      ws.current.onerror = (err) => {
        console.error(
          'Socket encountered error: ',
          err,
          'Closing socket'
        );
        ws.current?.close();
      }
  
  
      return () => {
        if(ws.current) {
          ws.current.close();
        }
      };
    }, []);

  return (
    <div>
    <h1>Lobby</h1>
    <button onClick={startGame}>Start Game</button>
    </div>
  );
}
export default LobbyPage;
