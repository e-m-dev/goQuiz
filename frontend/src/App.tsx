import './App.css'
import FormCreateRoom from './components/FormCreateRoom'
import FormJoinRoom from './components/FormJoinRoom'
import Landing from './components/Landing'

function App() {

  return (
    <div>
      <Landing></Landing>
      <FormCreateRoom></FormCreateRoom>
      <FormJoinRoom></FormJoinRoom>
    </div>
  )
}

export default App
