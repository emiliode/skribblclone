import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
const Home = (props) => {
  const [name, setName] = useState('')
  const [gameID, setGameID] = useState('')
  let navigate = useNavigate()
  const handleChange = (event) => {
    setName(event.target.value)
  }

  const handleSubmit = (event) => {
    fetch('/creategame', {
      method: 'POST',
      headers: { 'Content-Type': 'application/text' },
      body: name,
    })
      .then((res) => res.json())
      .then((data) => {
        console.log(data)
        setGameID(data.game)
        navigate('game/' + data.game, { replace: true, state: { name: name } })
      })
    event.preventDefault()
  }

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <label>
          Name:
          <input type='text' value={name} onChange={handleChange} />
        </label>
        <input type='submit' value='Create Game' />
      </form>

      <p>{gameID}</p>
    </div>
  )
}
export default Home
