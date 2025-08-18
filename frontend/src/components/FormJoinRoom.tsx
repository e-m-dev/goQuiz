import { useState } from "react";
import { joinRoom } from "../lib/api";
import { useLocalStorage } from "@uidotdev/usehooks";
import { Button } from "./ui/button";

const FormJoinRoom: React.FC = () => {
    const [roomCode, setRoomCode] = useState('');
    const [session, setSession] = useLocalStorage("session", { roomCode: "", playerId: "" });
    const playerName = localStorage.getItem("playerName") || '';

    const handleRoom = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRoomCode(event.target.value);
    }

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            const data = await joinRoom(roomCode, playerName);
            setSession({ roomCode, playerId: data.id });
            console.log(data);
        } catch (error) {
            console.log(error);
        }
    }
    return(
        <div className="w-full max-w-sm p-4 bg-white border border-gray-200 rounded-lg shadow-sm sm:p-6 md:p-8 dark:bg-gray-800 dark:border-gray-700">
            <form className="space-y-6" onSubmit={handleSubmit}>
                <h5 className="text-xl font-medium text-gray-900 dark:text-white">Join Room</h5>
                <div>
                    <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Room Code</label>
                    <input type="text" value={roomCode} onChange={handleRoom} name="rName" id="rName" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white" placeholder="XXXXXX" required />
                </div>
                <Button type="submit" variant="outline">Join</Button>
            </form>
        </div>
    );
};

export default FormJoinRoom;