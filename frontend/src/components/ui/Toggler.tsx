import React from "react";

interface TogglerProps<T> {
  state1: {
    label: T;
    onClick: () => void;
  };
  state2: {
    label: T;
    onClick: () => void;
  };
  currentState: T;
}

const Toggler = <T extends React.ReactNode>({
  state1,
  state2,
  currentState,
}: TogglerProps<T>) => {
  return (
    <div className="flex flex-row justify-center items-center">
      <button
        className={`
          px-4 py-2 w-25 rounded-l-sm cursor-pointer transition-all duration-200 ease-in-out
          ${currentState === state1.label ? "bg-blue-500" : "bg-gray-200 hover:bg-gray-300"}
        `}
        onClick={state1.onClick}
      >
        <p
          className={`
            font-sans font-medium text-base
            ${currentState === state1.label ? "text-white" : "text-gray-700"}
          `}
        >
          {state1.label}
        </p>
      </button>
      <button
        className={`
          px-4 py-2 w-25 rounded-r-sm cursor-pointer transition-all duration-200 ease-in-out
          ${currentState === state2.label ? "bg-blue-500" : "bg-gray-200 hover:bg-gray-300"}
        `}
        onClick={state2.onClick}
      >
        <p
          className={`
            font-sans font-medium text-base
            ${currentState === state2.label ? "text-white" : "text-gray-700"}
          `}
        >
          {state2.label}
        </p>
      </button>
    </div>
  );
};

export default Toggler;
