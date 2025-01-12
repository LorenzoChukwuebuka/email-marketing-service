import { Modal } from "../../../../components"
import { PlanData } from "../../../../store/admin/planStore";

interface Props {
    isOpen: boolean;
    onClose: () => void;
    plan: PlanData | null
}

const PaymentComponent: React.FC<Props> = ({ isOpen, onClose }) => {
    return <>

        <Modal isOpen={isOpen} onClose={onClose} title="Make Payment">

            <> Hello world  </>
        </Modal>

    </>
}

export default PaymentComponent