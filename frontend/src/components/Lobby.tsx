import { useLocalStorage } from '@uidotdev/usehooks';
import React, { useEffect, useState } from 'react';
import { getRoom, leaveRoom } from '../lib/api';
import { allowWS, closeWS, getWS } from '@/lib/WS';
import { Button } from './ui/button';
import { useNavigate } from 'react-router-dom';
/*
import { Input } from "./ui/input";
import { Label } from './ui/label';
*/
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "./ui/card";

const Lobby: React.FC= () => {
    const navigate = useNavigate();
    type Session = { roomCode: string, playerId: string };
    type Player = { id: string, name: string, host: boolean };

    const [session, setSession] = useLocalStorage<Session | null>("session", null);
    const [playerName] = useLocalStorage("playerName", null);
    const [players, setPlayers] = useState<Player[]>([]);

    useEffect(() => {
        if(!session) return;

        (async () => {
            try {
                const data = await getRoom(session.roomCode);
                setPlayers(data.players);
            } catch (err) {
                console.log("failed to fetch room", err);
            }
        })();
    }, [session]);

    useEffect(() => {
        if(!session) return;
        allowWS();
        getWS(`ws://localhost:8080/ws/${session.roomCode}?playerId=${session.playerId}`);
    }, [session]);

    const handleLeave = () => {
        if(session) { leaveRoom(session.roomCode, session.playerId); };
        setSession(null);
        closeWS(1000, "leave");
        navigate("/");
    };

    return (
        <div className='min-h-dvh grid place-items-center p-4'>
            <Card className='w-full max-w-sm'>
                <CardHeader>
                    <CardTitle>ROOM NAME</CardTitle>
                    <CardDescription>
                        Gathering Players! {session?.roomCode}
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <div className='flex justify-center gap-2 space-y-4'>
                        {players.length === 0 ? (
                            <p>No player joined.</p>
                        ) : (
                            <ul>
                                {players.map((player) => (
                                    <li key={player.id}>
                                        {player.name} (ID: {player.id}) {player.host ? 'Host' : ''}
                                    </li>
                                ))}
                                <li><Button variant={'destructive'} onClick={handleLeave}>Leave</Button></li>
                            </ul>
                        )}
                        
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};

export default Lobby;