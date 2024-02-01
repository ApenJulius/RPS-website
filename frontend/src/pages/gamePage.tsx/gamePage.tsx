import { useState, useRef, useEffect} from 'react';

import './gamePage.css';
import { MoveButton } from '../../components/MoveButton/MoveButton';


function GamePage() {
  const [playerMove, setPlayerMove] = useState<string>("")
  const ws = useRef<WebSocket | null>(null)
  
  useEffect(() => {
    ws.current = new WebSocket('ws://localhost:8000/ws');
    return () => {
      if(ws.current) {
        ws.current.close();
      }
    };
  }, []);

  // Sending a message
  const sendMessage = () => {
    if(ws.current) {
      ws.current.send(JSON.stringify({ move: 'test move' }));
    }
  };

  const handleButtonClick = (move: string) => {
    setPlayerMove(move)
    console.log(playerMove)
    if(ws.current) {
      console.log("current")
      ws.current.send(JSON.stringify({ move }));
    }
  }
  const options = ["rock", "paper", "scissors"] // gotta match backend fyi
  return (
      <div className='move-btn-container'>
        {
          options.map((move_option) => {
            return <MoveButton move={move_option} onClick={handleButtonClick} isDisabled={move_option === playerMove}/>
          })
        }
        </div>
  );
}

export default GamePage;
