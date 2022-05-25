import PropTypes from 'prop-types';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

async function loginUser(credentials) {
    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
    return fetch('http://localhost:3030/api/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'TimeZone': timezone,
        },
        body: JSON.stringify(credentials)
    }).then(data => data.json())
}

export default function Login({ setAuthDetails }) {
    const [username, setUserName] = useState();
    const [password, setPassword] = useState();
    const [error, setError] = useState();

    const navigate = useNavigate();

    const handleSubmit = async e => {
        e.preventDefault();
        const authDetails = await loginUser({
            username,
            password
        });

        if (authDetails.error) {
            if (authDetails.error === 'sql: no rows in result set') {
                setError('Akun tidak ditemukan');
                setTimeout(() => setError(null), 3000);
            }
        }

        setAuthDetails(authDetails);
        navigate('/dashboard');
    }
    
    return (
        <>
            <div className="flex flex-col items-center h-screen">
                <div className="bg-white shadow-lg rounded-lg px-8 pt-6 pb-8 mb-4">
                    {error &&
                        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mt-2 mb-2" role="alert" id="error-section">
                            <strong className="font-bold">Error! </strong>
                            <span className="block sm:inline">{error}</span>
                        </div>
                    }
                    <div className="flex flex-col items-center justify-center">


                        <h1 className="text-2xl font-bold text-center mb-10">Login</h1>
                        <form className="w-full max-w-lg" onSubmit={handleSubmit}>
                            <div className="flex flex-wrap -mx-3 mb-6">
                                <div className="w-full md:w-1/2 px-3 mb-6 md:mb-0">
                                    <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2" htmlFor="grid-first-name">
                                        Username
                                    </label>
                                    <input className="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500" id="grid-first-name" type="text" placeholder="Username" onChange={e => setUserName(e.target.value)} />
                                </div>
                                <div className="w-full md:w-1/2 px-3">
                                    <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2" htmlFor="grid-last-name">
                                        Password
                                    </label>
                                    <input className="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 leading-tight focus:outline-none focus:bg-white focus:border-gray-500" id="grid-last-name" type="password" placeholder="Password" onChange={e => setPassword(e.target.value)} />
                                </div>
                            </div>
                            <div className="flex items-center justify-between">
                                <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">
                                    Login
                                </button>
                                <a className="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800" href="index.js">
                                    Forgot Password?
                                </a>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </>
    );
}

Login.propTypes = {
    setAuthDetails: PropTypes.func.isRequired
}