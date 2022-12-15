import { useEffect, useState } from "react";


export const ChatWindow = ({ socket }) => {
    const [username, setUsername] = useState("James");
    const [inputText, setInputText] = useState("");
    const [messages, setMessages] = useState([]);

    const handleChange = (event) => {
        setInputText(event.target.value);
    }   

    const sendMessage = () => {
        if (!socket) alert("Not connected");

        // send message through websocket
        socket.emit("message", {sender: username, text: inputText});
        console.log("Sent message to server");

        setMessages((prev) => [...prev, {sender: username, text: inputText}]);
    }

    // setMessages whenever backend sends back a new message
    useEffect(() => {
        const messageListener = (content) => {
            setMessages((prev) => [...prev, JSON.parse(content)]);
        }

        socket.on("message", messageListener);

        return () => {
            socket.off("message", messageListener);
        }
    }, [socket]);


    return (
        <div id="chat-window">
            <div id="messages">{
                messages.map((message) => {
                    return (
                        <div class="message">
                            <span class="sender">{message.sender}</span>
                            <span class="text">{message.text}</span>
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
        <button id="send-button" onClick={sendMessage}>Send</button>
        </div>
    )
}