import "./MoveButton.css"


interface MoveButtonProps {
    move: string
    onClick: (move: string) => void
    isDisabled: boolean
}

export const MoveButton: React.FC<MoveButtonProps> = ({move, onClick, isDisabled}) => {
    return (
            <button className="move-btn" onClick={() => onClick(move)} disabled={isDisabled}>{move}</button>
    )
}