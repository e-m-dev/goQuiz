let sock: WebSocket | null = null;

export function getWS(url: string) {
    if(url.includes("undefined") || url.includes("null")) return null;

    if (sock && (sock.readyState == WebSocket.OPEN || sock.readyState == WebSocket.CONNECTING)) {
        return sock;
    }

    sock = new WebSocket(url);
    sock.onopen = () => console.log("[WS] Open");
    sock.onclose = (e) => console.log("[WS] Close", e.code);
    sock.onerror = (e) => console.error("[WS] Error", e);
    sock.onmessage = (m) => console.log("[WS] Message", m.data);
    return sock;
}