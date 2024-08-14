import { useParams } from "react-router-dom";

import { Navigate } from "react-router-dom";
import DragAndDropEditorPage from "./dragAndDropEditorPage";

const EditorRouter = () => {
    const { editorType } = useParams();

    switch (editorType) {
        case '1':
            return <DragAndDropEditorPage />;
        // case 'codeEditor':
        //     return <CodeEditor />;
        // case 'blockEditor':
        //     return <BlockEditor />;
        default:
            return <Navigate to="/404" replace />;
    }
};


export default EditorRouter