import ChatHistory from './chathistory'
import { sendMsg } from './../api/index'
import './chat.css'
const Chat = (props) => {
  const handleInput = (e) => {
    if (e.key === 'Enter') {
      makeGuess(e.target.value)
      e.preventDefault()
    }
  }

  const makeGuess = (guess) => {
    sendMsg(JSON.stringify({ event: 'GUESS', body: guess }))
    props.setChatHistory((prevchatHistory) => [
      ...prevchatHistory,
      { name: props.name, guess: guess },
    ])
  }
  return (
    <div className='chat'>
      <ChatHistory chatHistory={props.chatHistory} />
      <input className='chatinput' type='text' onKeyDown={handleInput} />
    </div>
  )
}
export default Chat
