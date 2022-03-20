import React from 'react'
import './App.css'
import Draw from './Draw'
import { BrowserRouter, Route, Routes, Link, Outlet } from 'react-router-dom'
import Home from './components/home'
function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path='game'
          element={
            <div>
              {' '}
              <h1>Wrong</h1> <Outlet />{' '}
            </div>
          }
        >
          <Route path=':gameId' element={<Draw />} />
        </Route>
        <Route path='/' element={<Home />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
