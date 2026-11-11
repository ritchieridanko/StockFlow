import React from "react";

interface PaginatorProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (value: number) => void;
}

const Paginator: React.FC<PaginatorProps> = ({
  currentPage,
  totalPages,
  onPageChange,
}: PaginatorProps) => {
  const isPrevClickable: boolean = currentPage > 1;
  const isNextClickable: boolean = currentPage < totalPages;

  return (
    <div className="flex flex-row justify-center items-center gap-2">
      <button
        className={`
          px-4 py-2 rounded-sm transition-all duration-200 ease-in-out
          ${isPrevClickable ? "bg-blue-100 hover:bg-blue-200 cursor-pointer" : "bg-gray-200"}
        `}
        onClick={() => {
          if (isPrevClickable) onPageChange(-1);
        }}
      >
        <p
          className={`
            font-sans font-medium text-base
            ${isPrevClickable ? "text-blue-500" : "text-gray-700"}
          `}
        >
          Prev
        </p>
      </button>
      <p className="px-5 py-2 flex justify-center items-center bg-blue-100 rounded-sm font-sans font-medium text-blue-500 text-base select-none">
        {currentPage} of {totalPages}
      </p>
      <button
        className={`
          px-4 py-2 rounded-sm transition-all duration-200 ease-in-out
          ${isNextClickable ? "bg-blue-100 hover:bg-blue-200 cursor-pointer" : "bg-gray-200"}
        `}
        onClick={() => {
          if (isNextClickable) onPageChange(1);
        }}
      >
        <p
          className={`
            font-sans font-medium text-base
            ${isNextClickable ? "text-blue-500" : "text-gray-700"}
          `}
        >
          Next
        </p>
      </button>
    </div>
  );
};

export default Paginator;
