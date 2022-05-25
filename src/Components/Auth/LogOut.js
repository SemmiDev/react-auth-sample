import { useNavigate } from "react-router-dom";

export default function LogOut() {
    localStorage.removeItem('authDetails');
    const navigate = useNavigate();
    navigate('/login');
}