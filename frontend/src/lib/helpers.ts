export function displayName(name: string, max = 16): string {
    name = name.trim();
    if(name.length <= max) {
        return name;
    } else {
        return name.slice(0, max - 3) + "...";
    }
}