import { useEffect, useRef, useState } from "react";
import { useUsername } from "../hooks/useUsername";


export const ChatPage = () => {
    const { username } = useUsername();
    const [inputText, setInputText] = useState("");
    const [messages, setMessages] = useState([]);
    const [socket, setSocket] = useState(null);
    const inputRef = useRef(null);

    const handleChange = (event) => {
        setInputText(event.target.value);
        inputRef.current.focus();
    }   

    const handleEnter = (event) => {
        if (event.key === "Enter" && inputText !== "") {
            socket.send(JSON.stringify({
                username: username,
                text: inputText
            }))
            inputRef.current.value = "";
            setInputText("");
            console.log("Sent message to server");
        }
    }

    const sendMessage = () => {
        if (inputText !== "") {
            // send message through websocket
            socket.send(JSON.stringify({
                username: username,
                text: inputText
            }))
            inputRef.current.value = "";
            setInputText("");
            console.log("Sent message to server");
        }
    }

    // setMessages whenever backend sends back a new message
    useEffect(() => {
        const newSocket = new WebSocket("ws://127.0.0.1:8000/join");
        setSocket(newSocket);
        
        newSocket.onmessage = (msg) => {
            setMessages((prev) => [...prev, JSON.parse(msg.data)]);
        }

        return () => {
            newSocket.close();
        }
    }, []);


    return (
        <div id="chat-window">
            <div id="messages">{
                messages.map((message, index) => {
                    return (
                        <div key={"message" + index} className="message">
                            <span className="sender">{message.username}</span>
                            <span className="text">{message.text}</span>
                        </div>
                    )
                })        
            }</div>
        <input 
            type="text" 
            id="message-input" 
            placeholder="Type your message here..." 
            onChange={handleChange}
            onKeyDown={handleEnter}
            value={inputText}
            ref={inputRef}
        />
        { socket ? 
        <button id="send-button" onClick={sendMessage}>Send</button> :
        <button id="send-button" >Not Connected</button>
        }
        
        </div>
    )
}