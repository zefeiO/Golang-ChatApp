import { Navigate } from "react-router-dom";
import { useUsername } from "../hooks/useUsername"

export const ProtectedRoute = ({ children }) => {
    const { username } = useUsername();
    if (username === "") {
        return <Navigate to="/" />
    }
    return children;
}