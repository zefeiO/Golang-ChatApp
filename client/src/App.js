import React, { useEffect, useState } from "react";
import { ChatPage } from "./pages/ChatPage";
import { UsernameProvider } from "./hooks/useUsername";
import { LoginPage } from "./pages/LoginPage";
import "./App.css"
import { Route, Routes } from "react-router-dom";
import { ProtectedRoute } from "./Components/ProtectedRoute";


const App = () => {
    return (
        <UsernameProvider>
            <Routes>
                <Route path="/" element={<LoginPage />} />
                <Route 
                    path="/chat" 
                    element={
                        <ProtectedRoute>
                            <ChatPage /> 
                        </ProtectedRoute>
                    } 
                />
            </Routes>
        </UsernameProvider>
    )
}

export default App;