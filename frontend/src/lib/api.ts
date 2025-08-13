export async function createRoom(name: string) {
    const res = await fetch('http://localhost:8080/rooms', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name })
    });
    if(!res.ok) throw new Error('Failed to create room');
    return res.json();
}