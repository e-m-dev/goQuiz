import { useLocalStorage } from '@uidotdev/usehooks';
import React, { useEffect, useState } from 'react';
import { getRoom } from '../lib/api';

const Lobby: React.FC= () => {
    type Session = { roomCode: string, playerId: string };
    type Player = { id: string, name: string, host: boolean };
    const [session] = useLocalStorage<Session | null>("session", null);
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


    return (
        <div className="lobby-container">
            <h2>Lobby</h2>
            <h3>Room Code: {session?.roomCode}</h3>
            <div className="players-list">
                <h4>Players:</h4>
                {players.length === 0 ? (
                    <p>No players have joined yet.</p>
                ) : (
                    <ul>
                        {players.map((player) => (
                            <li key={player.id}>
                                {player.name} (ID: {player.id}) {player.host ? '(Host)' : ''}
                            </li>
                        ))}
                    </ul>
                )}
            </div>
        </div>
    );
};

export default Lobby;