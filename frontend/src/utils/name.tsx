function renderApiName() {
    let apiName = import.meta.env.VITE_API_NAME;
    const firstFourLetters = apiName.slice(0, 4);
    const remainingLetters = apiName.slice(4);

    return (
        <>
            <span className="text-indigo-700">{firstFourLetters}</span>
            <span className="text-gray-700">{remainingLetters}</span>
            <i className="bi bi-mailbox2-flag text-indigo-700 ml-2"></i>
        </>
    );
}

export default renderApiName