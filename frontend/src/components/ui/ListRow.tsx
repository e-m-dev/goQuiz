import React from "react";

interface ListRowProps {
    left?: React.ReactNode;
    text: string;
    right?: React.ReactNode;
    className?: string;
}

const ListRow: React.FC<ListRowProps> = ({ left, text, right, className}) => {
    return (
        <div className={`flex items-center justify-between gap-3 py-2 ${className ?? ""}`}>
            <div className="shrink-0">{left}</div>
            <div className="flex-1 truncate">{text}</div>
            <div className="shrink-0">{right}</div>
        </div>
    )
}

export default ListRow;