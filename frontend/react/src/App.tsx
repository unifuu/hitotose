import './App.css';
import { Route, Routes } from 'react-router-dom';
import useToken from './useToken';
import SignIn from './Components/SignIn/SignIn';
import Blog from './Components/Blog/Blog';

export default function App() {
    const { token, setToken } = useToken()

    if (!token) {
        return (
            <div className="App">
                <Routes>
                    <Route path="/" element={<Blog authed={false} />} />
                    <Route path="/fuu" element={<SignIn setToken={setToken} />} />
                </Routes>
            </div>
        );
    } else {
        return (
            <div className="App">
                <Routes>
                    <Route path="/" element={<Blog authed={true} />} />
                </Routes>
            </div>
        );
    }
}