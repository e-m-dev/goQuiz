import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import './App.css'
import Landing from './components/Landing'
import FormCreateRoom from './components/FormCreateRoom'
import FormJoinRoom from './components/FormJoinRoom'

function App() {

  return (
    <Router>
      <Routes>
        <Route path='/' element={<Landing />} />
        <Route path='/create' element={<FormCreateRoom />} />
        <Route path='/join' element={<FormJoinRoom />} />
      </Routes>
    </Router>
  )
}

export default App
