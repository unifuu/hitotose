import './App.css';
import { Route, Routes } from 'react-router-dom';
import Game from './Components/Game/Game';
import useToken from './useToken';
import SignIn from './Components/SignIn/SignIn';

export default function App() {
    const { token, setToken } = useToken()

    if (!token) {
        return (
            <div className="App">
                <Routes>
                    <Route path="/" element={<Game authed={false} />} />
                    <Route path="/fuu" element={<SignIn setToken={setToken} />} />
                    <Route path="/game" element={<Game authed={false} />} />
                </Routes>
            </div>
        );
    } else {
        return (
            <div className="App">
                <Routes>
                    <Route path="/" element={<Game authed={true} />} />
                    <Route path="/fuu" element={<SignIn setToken={setToken} />} />
                    <Route path="/game" element={<Game authed={true} />} />
                </Routes>
            </div>
        );
    }
}