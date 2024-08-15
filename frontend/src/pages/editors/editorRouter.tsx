import { useParams } from "react-router-dom";

import { Navigate } from "react-router-dom";
import DragAndDropEditorPage from "./dragAndDropEditorPage";
import CodeEditorPage from "./codeEditor";

const EditorRouter = () => {
    const { editorType } = useParams();

    switch (editorType) {
        case '1':
            return <DragAndDropEditorPage />;
        case '2':
            return <CodeEditorPage />;
        // case 'blockEditor':
        //     return <BlockEditor />;
        default:
            return <Navigate to="/404" replace />;
    }
};
export default EditorRouter