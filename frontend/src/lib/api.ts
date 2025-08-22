const API_BASE = import.meta.env.VITE_API_BASE;

export async function createRoom(name: string) {
    const res = await fetch(`${API_BASE}/rooms`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name })
    });
    if(!res.ok) throw new Error('Failed to create room');
    return res.json();
}

export async function joinRoom(code: string, name: string) {
    const safeCode = encodeURIComponent(code.trim().toUpperCase());
    const res = await fetch(`${API_BASE}/rooms/${safeCode}/join`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name })
    });
    if(!res.ok) {
        throw new Error((await res.text()) || 'Failed to join room');
    }
    return res.json();
}

export async function leaveRoom(code: string, id: string) {
    const safeCode = encodeURIComponent(code.trim().toUpperCase());
    const res = await fetch(`${API_BASE}/rooms/${safeCode}/leave`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id })
    });
    if(!res.ok) {
        throw new Error((await res.text()) || 'Failed to leave room');
    }
    return res.json();
}

export async function getRoom(code: string) {
    const safeCode = encodeURIComponent(code.trim().toUpperCase());
    const res = await fetch(`${API_BASE}/rooms/${safeCode}`, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
    });
    if(!res.ok) {
        throw new Error((await res.text()) || 'Failed to get room');
    }
    return res.json();
}