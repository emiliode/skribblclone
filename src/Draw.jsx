import React, { useState, useEffect, useRef } from 'react'
import { useLocation, useParams } from 'react-router-dom'
import './Draw.css'
import Chat from './components/chat'
import CanvasDraw from 'react-canvas-draw'
import { CompactPicker } from 'react-color'
import { connect, sendMsg } from './api'

const Draw = (props) => {
  const addPlayer = (newplayer) => setPlayer([...players, newplayer])
  let params = useParams()
  let location = useLocation()
  const [name, setName] = useState('')
  const [chatHistory, setChatHistory] = useState([])
  const [color, setColor] = useState('#444')
  const saveableCanvas = useRef(null)
  const [isDrawer, setDrawing] = useState(false)
  const [players, setPlayer] = useState([])

  useEffect(() => {
    console.log('connecting to game: ' + params.gameId)
    connect(params.gameId, handleMessage)
    if (location.state !== null) {
      console.log(location.state.name)
      setName(location.state.name)
      sendMsg(location.state.name)
    }
  }, [])

  const handleMessage = (msg) => {
    let json = JSON.parse(msg.data)
    switch (json.event) {
      case 'join':
        addPlayer(json.client)
        break
      case 'JOINCOMPLETE':
        let data = JSON.parse(json.body)
        console.log(data)
        setPlayer(data)
        break
      case 'DRAWUPDATE':
        saveableCanvas.current.loadSaveData(JSON.stringify(json.body), true)
        break
      case 'GUESS':
        setChatHistory((prevchatHistory) => [
          ...prevchatHistory,
          json.client.id + ': ' + json.body,
        ])
        console.log(json.client.id + ' guessed: ' + json.body)
        break
      default:
        console.log(json)
        break
    }
  }
  const handleChange = (event) => {
    setName(event.target.value)
  }
  const handleSubmit = (event) => {
    sendMsg(name)
    if (name === 'emil') {
      setDrawing(true)
    }
    event.preventDefault()
  }
  const handleChangeComplete = (color, event) => {
    setColor(color.hex)
  }
  const sendDrawing = () => {
    if (isDrawer) {
      //   console.log(saveableCanvas.getSaveData())
      sendMsg(saveableCanvas.current.getSaveData())
    }
  }
  return (
    <div>
      {players.length == 0 ? (
        <form onSubmit={handleSubmit}>
          <label>
            Name:
            <input type='text' value={name} onChange={handleChange} />
          </label>
          <input type='submit' value='Submit' />
        </form>
      ) : (
        <div className='container'>
          <div className={'main'}>
            <p>A simple drawing Canvas</p>
            <CompactPicker onChangeComplete={handleChangeComplete} />
            <div
              style={{
                display: 'flex',
                justify_content: 'space-between',
                width: '400px',
              }}
            >
              <button
                onClick={() => {
                  saveableCanvas.current.eraseAll()
                }}
              >
                {' '}
                Erase
              </button>
            </div>
            <CanvasDraw
              disabled={!isDrawer}
              hideGrid={false}
              brushColor={color}
              ref={saveableCanvas}
              onChange={sendDrawing}
              canvasWidth={800}
              canvasHeight={600}
            />
          </div>
          <Chat
            chatHistory={chatHistory}
            name={name}
            setChatHistory={setChatHistory}
          />
        </div>
      )}
    </div>
  )
}

export default Draw
