import { useLocalStorage } from '@uidotdev/usehooks';
import { useNavigate } from 'react-router-dom';
import React from 'react';
import { Button } from './ui/button';

const Landing: React.FC = () => {
    const [name, saveName] = useLocalStorage("playerName", '');
    const navigate = useNavigate();

    const handleName = (event: React.ChangeEvent<HTMLInputElement>) => {
        saveName(event.target.value);
    }

    return (
        <form className="max-w-sm mx-auto">
            <div className="mb-5">
                <label htmlFor="base-input" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Player Name</label>
                <input type="text" id="base-input" value={name} onChange={handleName} className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
            </div>
            <Button onClick={() => {
                if (!name.trim()) {
                alert("Please enter your name first!");
                return;
                }
                navigate("/create");
            }}>Create Room</Button>
            <Button onClick={() => {
                if (!name.trim()) {
                    alert("Please enter your name first!");
                    return;
                }
                navigate("/join");
            }}>Join Room</Button>
        </form>
        
    );
};

export default Landing;