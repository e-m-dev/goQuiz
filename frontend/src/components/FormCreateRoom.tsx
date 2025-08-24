import React, { useState } from 'react';
import { useLocalStorage } from '@uidotdev/usehooks';
import { createRoom } from '../lib/api';
import { joinRoom } from '../lib/api';
import { Button } from './ui/button';
import { Input } from "./ui/input";
import { Label } from './ui/label';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "./ui/card";
import BackButton from './ui/BackButton';
import { Copy, LogIn, Check } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const FormCreateRoom: React.FC = () => {
    const navigate = useNavigate();
    const [roomCode, setRoomCode] = useState('');
    const [roomName, setRoomName] = useState('');
    const [creating, setCreating] = useState(false);
    const [copied, setCopied] = useState(false);
    const [session, setSession] = useLocalStorage("session", { roomCode: "", playerId: "" });
    const rawName = localStorage.getItem("playerName");
    const playerName = rawName ? JSON.parse(rawName) : "";      // Strip JSON encoding from useLocalStorage hook -> Backend expects raw string


    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setCreating(true);

        try {
            const data = await createRoom(roomName);
            console.log('Created Room:', data);
            setRoomCode(data.code);
        } catch (error) {
            console.error(error);
            setCreating(false);
        } 

    }

    const handleJoin = async () => {
        try {
            const data = await joinRoom(roomCode, playerName);
            setSession({ roomCode, playerId: data.id });
            navigate("/room");
            console.log(data);
        } catch (error) {
            console.log(error);
        }
    }

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRoomName(event.target.value);
    }

    const handleCopy = async () => {
        if(!roomCode) return;
        await navigator.clipboard.writeText(roomCode);
        setCopied(true)
        setTimeout(() => {
            setCopied(false);
        }, 2000);
    }

    return (
        <div className='min-h-dvh grid place-items-center p-4'>
            <BackButton location='/' />
            <Card className='w-full max-w-sm'>
                <CardHeader>
                    <CardTitle>Create Room</CardTitle>
                    <CardDescription>
                        Enter the room details to create!
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form className="space-y-4" onSubmit={handleSubmit}>
                        <Label>Room Name</Label>
                        <Input type="text" value={roomName} onChange={handleChange} required />
                        <div className='flex justify-center gap-2'>
                            <Button type="submit" disabled={creating}>Create Room</Button>
                        </div>
                        <Card>
                            <CardContent>
                                <div className='w-full flex items-center justify-between text-left'>
                                    <span className='font-mono tracking-wider text-xl select-all'>{roomCode}</span>
                                    <div className='flex items-center gap-2'>
                                        <Button type="button" size="icon" variant="ghost" onClick={handleCopy} disabled={!roomCode} aria-label={copied ? "Copied" : "Copy room code"}>
                                            {copied ? <Check className="h-4 w-4 text-green-500" /> : <Copy className="h-4 w-4" />}
                                        </Button>

                                        <Button type='button' size="icon" variant="ghost" onClick={handleJoin} disabled={!roomCode}>
                                            <LogIn className='h-4 w-4' />
                                        </Button>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
};

export default FormCreateRoom;