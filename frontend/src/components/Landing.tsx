import { useLocalStorage } from '@uidotdev/usehooks';
import { useNavigate } from 'react-router-dom';
import React from 'react';
import { Button } from './ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Input } from './ui/input';

const Landing: React.FC = () => {
    const [name, saveName] = useLocalStorage("playerName", '');
    const navigate = useNavigate();

    const handleName = (event: React.ChangeEvent<HTMLInputElement>) => {
        saveName(event.target.value);
    }

    return (
        <div className='min-h-dvh grid place-items-center p-4'>
            <Card className='w-full max-w-sm'>
                <CardHeader>
                    <CardTitle>Player Name</CardTitle>
                    <CardDescription>
                        Your name will be displayed in game!
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form className="space-y-4">
                        <Input type="text" value={name} onChange={handleName} required />
                        <div className='flex justify-center gap-2'>
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
                        </div>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
};

export default Landing;