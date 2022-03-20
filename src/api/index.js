var socket = null

let connect = (gameid, cb) => {
  socket = new WebSocket('ws://localhost:8080/game/' + gameid)
  console.log('Attempting Connection...')

  socket.onopen = () => {
    console.log('Successfully Connected')
  }

  socket.onmessage = (msg) => {
    console.log(msg)
    cb(msg)
  }

  socket.onclose = (event) => {
    console.log('Socket Closed Connection: ', event)
  }

  socket.onerror = (error) => {
    console.log('Socket Error: ', error)
  }
}

let sendMsg = (msg) => {
  setTimeout(function () {
    if (socket === null) sendMsg()
    else {
      console.log('sending msg: ', msg)
      socket.send(msg)
    }
  }, 500)
}

export { connect, sendMsg }
