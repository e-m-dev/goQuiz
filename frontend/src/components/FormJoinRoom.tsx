import { useState } from "react";
import { joinRoom } from "../lib/api";
import { useLocalStorage } from "@uidotdev/usehooks";
import { useNavigate } from "react-router-dom";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "./ui/card";
import BackButton from "./ui/BackButton";

const FormJoinRoom: React.FC = () => {
    const navigate = useNavigate();
    const [roomCode, setRoomCode] = useState('');
    const [session, setSession] = useLocalStorage("session", { roomCode: "", playerId: "" });
    const rawName = localStorage.getItem("playerName");
    const playerName = rawName ? JSON.parse(rawName) : "";      // Strip JSON encoding from useLocalStorage hook -> Backend expects raw string

    const handleRoom = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRoomCode(event.target.value);
    }

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            const data = await joinRoom(roomCode, playerName);
            setSession({ roomCode, playerId: data.id });
            navigate("/room");
            console.log(data);
        } catch (error) {
            console.log(error);
        }
    }
    return(
        <div className='min-h-dvh grid place-items-center p-4'>
            <BackButton location='/' />
            <Card className='w-full max-w-sm'>
                <CardHeader>
                    <CardTitle>Join Game</CardTitle>
                    <CardDescription>
                        Enter the room code to join!
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form className="space-y-4" onSubmit={handleSubmit}>
                        <Input type="text" value={roomCode} onChange={handleRoom} required />
                        <div className='flex justify-center gap-2'>
                            <Button type="submit">Join Room</Button>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
};

export default FormJoinRoom;