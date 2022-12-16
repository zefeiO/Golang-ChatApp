import { useEffect, useState } from "react";
import { useUsername } from "../hooks/useUsername";


export const ChatPage = ({ socket }) => {
    const { username } = useUsername();
    const [inputText, setInputText] = useState("");
    const [messages, setMessages] = useState([]);

    const handleChange = (event) => {
        setInputText(event.target.value);
    }   

    const sendMessage = () => {
        console.log(socket);
        // send message through websocket
        socket.send(JSON.stringify({
            username: username,
            text: inputText
        }))
        console.log("Sent message to server");
    }

    // setMessages whenever backend sends back a new message
    useEffect(() => {
        if (socket) {
            socket.onmessage = (msg) => {
                setMessages((prev) => [...prev, JSON.parse(msg.data)]);
            }
    
            return () => {
                socket.close();
            }
        }
    }, [socket]);


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
            value={inputText}
        />
        { socket ? 
        <button id="send-button" onClick={sendMessage}>Send</button> :
        <button id="send-button" >Not Connected</button>
        }
        
        </div>
    )
}