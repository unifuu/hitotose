import './App.css';
import { Route, Routes } from 'react-router-dom';
import Game from './Components/Game/Game';

function App() {
  return (
    <div className="App">
      <Routes>
        <Route path="/" element={<Game />} />
        <Route path="/game" element={<Game />} />
      </Routes>
    </div>
  );
}

export default App;
