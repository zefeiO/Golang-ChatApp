import { useRef } from "react";
import { useNavigate } from "react-router-dom"
import { useUsername } from "../hooks/useUsername";

export const LoginPage = () => {
    const navigate = useNavigate();
    const {username, setUsername} = useUsername();
    const inputRef = useRef(null);

    const handleChange = (event) => {
        setUsername(event.target.value);
    }

    const handleEnter = (event) => {
        if (event.key === "Enter" && username !== "") {
            navigate("/chat");
        }
    }

    const startChat = () => {
        if (username === "") return;
        navigate("/chat");
    }

    return (
        <div id="chat-window">
            <h1 id="login-title">Enter your name for the session</h1>
            <input 
                type="text" 
                id="username-input" 
                onChange={handleChange}
                onKeyDown={handleEnter}
                value={username}
                ref={inputRef}
            />
            <button id="username-button" onClick={startChat}>Join</button> 
        </div>
    )
}