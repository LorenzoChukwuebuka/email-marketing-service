import { useParams } from "react-router-dom";

import { Navigate } from "react-router-dom";
import CodeEditor from "../templates/codeEditor";
import DragAndDropEditor from "../templates/dragAndDrop";
import RichTextEditor from "../templates/richTextEditor";


const EditorRouter = () => {
    const { editorType } = useParams();

    switch (editorType) {
        case '1':
            return <DragAndDropEditor />;
        case '2':
            return <CodeEditor />;
        case '3':
            return <RichTextEditor />;
        default:
            return <Navigate to="/404" replace />;
    }
};
export default EditorRouter


