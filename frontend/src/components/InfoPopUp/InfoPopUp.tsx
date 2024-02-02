import "./InfoPopUp.css"


interface InfoPopUpProps {
    move: string
    onClick: (move: string) => void
    isDisabled: boolean
}

export const InfoPopUp: React.FC<InfoPopUpProps> = ({move, onClick, isDisabled}) => {
    return (
            <div className="info-window">
                
            </div>
    )
}