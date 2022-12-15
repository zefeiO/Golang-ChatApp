import React, { useEffect, useState } from "react";
import { io } from "socket.io-client";
import { ChatWindow } from "./Components/ChatWindow";
import "./App.css"

const App = () => {
    const [socket, setSocket] = useState(null);

    useEffect(() => {
        const newSocket = io("http://localhost:8080/join");
        setSocket(newSocket);
        return () => newSocket.close();
    }, [])

    return (
        <ChatWindow socket={socket}/>
    )
}

export default App;