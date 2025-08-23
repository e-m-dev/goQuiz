import { Button } from "./button";
import { ArrowLeft } from "lucide-react";
import { useNavigate } from "react-router-dom";

interface BackButtonProps {
    location: string
}

const BackButton: React.FC<BackButtonProps> = ({ location }) => {
    const nav = useNavigate();
    return (
        <Button variant="ghost" size="icon" onClick={() => nav(location)} className="absolute top-4 left-4">
            <ArrowLeft className="h-5 w-5" />
        </Button>
    )
}

export default BackButton;