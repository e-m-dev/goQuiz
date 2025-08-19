import React, { useState } from 'react';
import { createRoom } from '../lib/api';
import { joinRoom } from '../lib/api';
import { Button } from './ui/button';
import { Input } from "./ui/input";
import { Label } from './ui/label';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "./ui/card";

const FormCreateRoom: React.FC = () => {
    const [roomCode, setRoomCode] = useState('');
    const [roomName, setRoomName] = useState('');

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            const data = await createRoom(roomName);
            console.log('Created Room:', data);
            setRoomCode(data.code);
            //TODO: join and navigate to room
        } catch (error) {
            console.error(error);
        }

    }

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRoomName(event.target.value);
    }

    return (
        <div className='min-h-dvh grid place-items-center p-4'>
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
                            <Button type="submit">Create Room</Button>
                        </div>
                        <Card>
                            <CardHeader>
                                <CardTitle>Room Details</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <p>{roomCode}</p>
                            </CardContent>
                        </Card>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
};

export default FormCreateRoom;