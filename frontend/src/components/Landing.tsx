import { useLocalStorage } from '@uidotdev/usehooks';
import React from 'react';

const Landing: React.FC = () => {
    const [name, saveName] = useLocalStorage("playerName", '');

    const handleName = (event: React.ChangeEvent<HTMLInputElement>) => {
        saveName(event.target.value);
    }

    return (
        <form className="max-w-sm mx-auto">
            <div className="mb-5">
                <label htmlFor="base-input" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Player Name</label>
                <input type="text" id="base-input" value={name} onChange={handleName} className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
            </div>
        </form>
    );
};

export default Landing;