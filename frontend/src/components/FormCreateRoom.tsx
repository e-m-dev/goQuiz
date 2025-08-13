import React, { useState } from 'react';
import { createRoom } from '../lib/api';

const FormCreateRoom: React.FC = () => {
    const [roomCode, setRoomCode] = useState('');
    const [roomName, setRoomName] = useState('');

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            const data = await createRoom(roomName);
            console.log('Created Room:', data);
            setRoomCode(data.code);
        } catch (error) {
            console.error(error);
        }

    }

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRoomName(event.target.value);
    }

    return (
        <div className="w-full max-w-sm p-4 bg-white border border-gray-200 rounded-lg shadow-sm sm:p-6 md:p-8 dark:bg-gray-800 dark:border-gray-700">
            <form className="space-y-6" onSubmit={handleSubmit}>
                <h5 className="text-xl font-medium text-gray-900 dark:text-white">Create Room</h5>
                <div>
                    <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Room Name</label>
                    <input type="text" value={roomName} onChange={handleChange} name="rName" id="rName" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white" placeholder="Friday Pub Quiz" required />
                </div>
                <div>
                    <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Round Options, Coming Soon</label>
                    <input type="text" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white" />
                </div>
                <div className="flex items-start">
                    <div className="flex items-start">
                        <div className="flex items-center h-5">
                            <input id="remember" type="checkbox" value="" className="w-4 h-4 border border-gray-300 rounded-sm bg-gray-50 focus:ring-3 focus:ring-blue-300 dark:bg-gray-700 dark:border-gray-600 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800" />
                        </div>
                        <label className="ms-2 text-sm font-medium text-gray-900 dark:text-gray-300">Do Something</label>
                    </div>
                </div>
                <button type="submit" className="w-full text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Create Room</button>
            </form>
            <div className="block max-w-sm p-6 bg-white border border-gray-200 rounded-lg shadow-sm hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700">
                <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">Room Details</h5>
                <p className="font-normal text-gray-700 dark:text-gray-400">{roomCode}</p>
            </div>
        </div>
    );
};

export default FormCreateRoom;