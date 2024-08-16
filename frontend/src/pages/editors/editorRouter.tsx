import { useParams } from "react-router-dom";

import { Navigate } from "react-router-dom";
import DragAndDropEditorPage from "./dragAndDropEditorPage";
import CodeEditorPage from "./codeEditor";
import RichTextEditorPage from "./richTextEditor";

const EditorRouter = () => {
    const { editorType } = useParams();

    switch (editorType) {
        case '1':
            return <DragAndDropEditorPage />;
        case '2':
            return <CodeEditorPage />;
        case '3':
            return <RichTextEditorPage />;
        default:
            return <Navigate to="/404" replace />;
    }
};
export default EditorRouter