import React, { useEffect, useState } from "react";
import { ChatPage } from "./pages/ChatPage";
import { UsernameProvider } from "./hooks/useUsername";
import { LoginPage } from "./pages/LoginPage";
import "./App.css"
import { Route, Routes } from "react-router-dom";
import { ProtectedRoute } from "./Components/ProtectedRoute";


const App = () => {
    const [socket, setSocket] = useState(null);

    useEffect(() => {
        const newSocket = new WebSocket("ws://127.0.0.1:8000/join");
        setSocket(newSocket);
        return () => newSocket.close();
    }, [])

    return (
        <UsernameProvider>
            <Routes>
                <Route path="/" element={<LoginPage />} />
                <Route 
                    path="/chat" 
                    element={
                        <ProtectedRoute>
                            <ChatPage socket={socket}/> 
                        </ProtectedRoute>
                    } 
                />
            </Routes>
        </UsernameProvider>
    )
}

export default App;