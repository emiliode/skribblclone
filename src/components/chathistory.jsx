import React from 'react'
import './chathistory.css'
const ChatHistory = (props) => {
  const messages = props.chatHistory.map((msg, index) => (
    <div className='chatMessage' key={index}>
      <span className='guessname'>{msg.name + ':'}</span>
      <span className='guess'>{msg.guess}</span>
    </div>
  ))

  return (
    <div className='ChatHistory'>
      <h2>Chat History</h2>
      {messages}
    </div>
  )
}
export default ChatHistory
