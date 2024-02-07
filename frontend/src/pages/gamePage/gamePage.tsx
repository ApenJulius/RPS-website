import { useState, useRef, useEffect} from 'react';
import {useParams} from 'react-router-dom';
import './gamePage.css';
import { MoveButton } from '../../components/MoveButton/MoveButton';

enum ErrorCode {
  GAME_FOUND = 10,
  PLAYER_JOINED = 20,
  PLAYER_LEFT = 30,
}



function GamePage() {
  const [playerMove, setPlayerMove] = useState<string>("")
  const { id } = useParams();
  const [connectedPlayers, setConnectedPlayers] = useState<number>(0) // TODO: get from backend
  const [maxPlayers, setMaxPlayers] = useState<number>(2) // TODO: get from backend
  const [status, setStatus] = useState<string>("Disconnected")
  const ws = useRef<WebSocket | null>(null)
  
  const lookingForGame = () => {
    setStatus("Looking for game...")
  }


  useEffect(() => {
    if(!id) {
      console.error("No groupID provided")
      return;
    }
    ws.current = new WebSocket(`ws://localhost:8000/ws?groupID=${id}`);
    
    ws.current.onopen = () => {
      console.log('ws opened');
      lookingForGame();
    };

    ws.current.onmessage = (message) => {
      const data = JSON.parse(message.data);
      console.log(data);
      console.log(data.data);
      
      switch(Number(data.code)) {
        case ErrorCode.GAME_FOUND:
          setStatus(`Connected to: ${data.groupID}`)
          break;
        case ErrorCode.PLAYER_JOINED:
        case ErrorCode.PLAYER_LEFT:
          setConnectedPlayers(data.data.current)
          setMaxPlayers(data.data.max)
          break;
        default:
          console.error("Unknown message received", data)
        }
      };

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

  // Sending a message
  const sendMessage = (move: string) => {
    if(ws.current) {
      ws.current.send(JSON.stringify({ move }));
    }
  };

  const handleButtonClick = (move: string) => {
    setPlayerMove(move)
    console.log(playerMove)
    sendMessage(move)
  }
  const options = ["rock", "paper", "scissors"] // gotta match backend fyi
  return (
    <div>
      <h1>{status}</h1>
      <h3>Connected players: {connectedPlayers}/{maxPlayers}</h3>
      <div className='move-btn-container'>
        {
          options.map((move_option) => {
            return <MoveButton move={move_option} onClick={handleButtonClick} isDisabled={move_option === playerMove}/>
          })
        }
        </div>
      </div>
  );
}

export default GamePage;
