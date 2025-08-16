export async function createRoom(name: string) {
    const res = await fetch('http://localhost:8080/rooms', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name })
    });
    if(!res.ok) throw new Error('Failed to create room');
    return res.json();
}

export async function joinRoom(code: string, Name: string) {
    const safeCode = encodeURIComponent(code.trim().toUpperCase());
    const res = await fetch(`http://localhost:8080/rooms/${safeCode}/join`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ Name })
    });
    if(!res.ok) {
        throw new Error((await res.text()) || 'Failed to join room');
    }
    return res.json();
}

export async function getRoom(code: string) {
    const safeCode = encodeURIComponent(code.trim().toUpperCase());
    const res = await fetch(`http://localhost:8080/rooms/${safeCode}`, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
    });
    if(!res.ok) {
        throw new Error((await res.text()) || 'Failed to get room');
    }
    return res.json();
}