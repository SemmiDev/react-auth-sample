import Dashboard from "./Components/Dashboard/Dashboard";
import Preferences from "./Components/Preferences/Preferences";
import { Link, Route, Routes } from "react-router-dom";
import Login from "./Components/Auth/Login";
import Register from "./Components/Auth/Register";
import LogOut from "./Components/Auth/LogOut";

const authDetailsKey = "authDetails";

function setAuthDetails(authDetails) {
  localStorage.setItem(authDetailsKey, JSON.stringify(authDetails));
}

function getAuthDetails() {
  return JSON.parse(localStorage.getItem(authDetailsKey));
}

function App() {
  const authDetails = getAuthDetails();
  return (
    <div>

      {authDetails ? (
        <nav>
          <ul>
            <li>
              <Link to="/login">Login</Link>
            </li>
            <li>
              <Link to="/register">Register</Link>
            </li>
          </ul>
        </nav>
      ) : (
        <nav>
          <ul>
            <li>
              <Link to="/dashboard">Dashboard</Link>
            </li>
            <li>
              <Link to="/preferences">Preferences</Link>
            </li>
            <li>
              <Link to="/logout">Log Out</Link>
            </li>
          </ul>
        </nav>
      )}

      <Routes>
        <Route path="/login" element={<Login setAuthDetails={setAuthDetails} />} />
        <Route path="/register" element={<Register setAuthDetails={setAuthDetails} />} />
        <Route path='/logout' element={<LogOut />}></Route>
        <Route path='/dashboard' element={<Dashboard />}></Route>
        <Route path='/preferences' element={<Preferences />}></Route>
        <Route path='*' element={<div>Not found</div>}></Route>
      </Routes>
    </div>
  );
}


export default App;
