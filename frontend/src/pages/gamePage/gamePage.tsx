import { useState, useRef, useEffect} from 'react';
import {useParams} from 'react-router-dom';
import './gamePage.css';
import { MoveButton } from '../../components/MoveButton/MoveButton';

import { ErrorCode } from '../../constants/ErrorCodes';
const { GAME_FOUND, PLAYER_JOINED, PLAYER_LEFT, GAME_COUNTDOWN } = ErrorCode


function GamePage() {
  const [playerMove, setPlayerMove] = useState<string>("")
  const { id } = useParams();
  const [connectedPlayers, setConnectedPlayers] = useState<number>(0) // TODO: get from backend
  const [maxPlayers, setMaxPlayers] = useState<number>(2) // TODO: get from backend
  const [status, setStatus] = useState<string>("Disconnected")
  const ws = useRef<WebSocket | null>(null)
  const [gameMessage, setGameMessage] = useState<string>("")
  
  const lookingForGame = () => {
    setStatus("Looking for game...")
  }


  useEffect(() => {
    if(!id) {
      console.error("No groupID provided")
      return;
    }
    ws.current = new WebSocket(`ws://${process.env.REACT_APP_WEBSITE_NAME}:8000/game?groupID=${id}`);
    
    ws.current.onopen = () => {
      console.log('ws opened');
      lookingForGame();
    };

    ws.current.onmessage = (message) => {
      const data = JSON.parse(message.data);
      console.log(data);
      
      switch(Number(data.code)) {
        case GAME_FOUND:
          setStatus(`Connected to: ${data.groupID}`)
          break;
        case PLAYER_JOINED:
        case PLAYER_LEFT:
          setConnectedPlayers(data.data.current)
          setMaxPlayers(data.data.max)
          break;
        case GAME_COUNTDOWN:
          setGameMessage(data.info)
          console.log(data.info)
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

  const copyToClipboard = () => {
    navigator.clipboard.writeText(window.location.href)
    .then(() => {
      const copyFeedback = document.getElementById('copy-feedback');
      if (copyFeedback) {
        copyFeedback.classList.add('fade-in');
        setTimeout(() => {
          copyFeedback.classList.remove('fade-in');
        }, 1000);
      }
    })
    .catch(err => {
      console.error('Failed to copy page URL: ', err);
    });
  }



  const options = ["rock", "paper", "scissors"] // gotta match backend fyi
  return (
    <div>
      <h1>{status}</h1>
      <h3>Connected players: {connectedPlayers}/{maxPlayers}</h3>
      <button id="clipboard-copy" onClick={copyToClipboard}>Copy game link</button> <span id="copy-feedback">Copied</span>
      <h2>{gameMessage}</h2>
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
