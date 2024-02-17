import { useState, useRef, useEffect} from 'react';
import {useNavigate} from 'react-router-dom';
import './LobbyPage.css';
import { ErrorCode } from '../../constants/ErrorCodes';
const { LOBBY_UPDATE } = ErrorCode
const { REACT_APP_WEBSITE_NAME } = process.env;


interface Group {
  clients: string[];
  max: number;
}

function LobbyPage() {
  const navigate = useNavigate();
  const [lobbies, setLobbies] = useState({} as {[key: string]: Group});
  const ws = useRef<WebSocket | null>(null)
    const startGame = () => {
        console.log("start game")
        navigate("/" + genRandomString(10))
    
    }
    console.log(REACT_APP_WEBSITE_NAME)
    const genRandomString = (length: number) =>{
      var chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
      var charLength = chars.length;
      var result = '';
      for ( var i = 0; i < length; i++ ) {
         result += chars.charAt(Math.floor(Math.random() * charLength));
      }
      return result;
    }


    useEffect(() => {
      ws.current = new WebSocket(`ws://${REACT_APP_WEBSITE_NAME}:8000/lobby`);
      
      ws.current.onopen = () => {
        console.log('ws opened');
      };
  
      ws.current.onmessage = (message) => {
        const data = JSON.parse(message.data);
        switch(Number(data.code)) {
          case LOBBY_UPDATE:
            console.log("Connected to lobby");
            setLobbies(data.data);
            break;
          default:
            console.error("UNKNOWN CODE AT WS SWITCH", data.code)
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
    <table className='lobby-container'>
      <thead>
        <tr>
          <th>Lobby</th>
          <th>Players</th>
        </tr>
      </thead>
      <tbody>
      {
      Object.entries(lobbies).map(([key, value]: [string, Group], index) => {
      return (
        <tr key={index} className='row' onClick={() => navigate(`/${key}`)}>
          <td>{key}</td>
          <td>{value.clients}/{value.max}</td>
        </tr>
      );
    })}
      </tbody>
      </table>
    </div>
  );
}
export default LobbyPage;
