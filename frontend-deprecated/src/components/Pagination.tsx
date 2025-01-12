import React from 'react';
import { PaginatedResponse } from '../interface/pagination.interface';


type PaginationProps = {
    paginationInfo: Omit<PaginatedResponse<any>, "data">
    handlePageChange: (newPage: number) => void
    item:string
}

const Pagination: React.FC<PaginationProps> = ({ paginationInfo, handlePageChange,item }) => {
    return (
        <div className="mt-4 flex justify-between items-center">
            <div>Total {item}: {paginationInfo.total_count}</div>
            <div className="flex space-x-2">
                <button
                    onClick={() => handlePageChange(paginationInfo.current_page - 1)}
                    disabled={paginationInfo.current_page === 1}
                    className="px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300"
                >
                    Previous
                </button>
                <span className="py-2">
                    Page {paginationInfo.current_page} of {paginationInfo.total_pages}
                </span>
                <button
                    onClick={() => handlePageChange(paginationInfo.current_page + 1)}
                    disabled={paginationInfo.current_page === paginationInfo.total_pages}
                    className="px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300"
                >
                    Next
                </button>
            </div>
        </div>
    );
};

export default Pagination;
