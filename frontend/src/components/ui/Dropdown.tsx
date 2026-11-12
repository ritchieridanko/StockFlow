import React from "react";

interface DropdownProps {
  selectedOption: string;
  defaultOption: string;
  options: string[];
  onSelect: (value: string) => void;
}

const Dropdown: React.FC<DropdownProps> = ({
  selectedOption,
  defaultOption,
  options,
  onSelect,
}: DropdownProps) => {
  return (
    <select
      value={selectedOption}
      className={`
        px-2 py-1 w-full flex flex-row justify-between items-center gap-2
        rounded-sm outline-2
        font-sans font-normal text-base
        cursor-pointer
        transition-all duration-200 ease-in-out
        focus:bg-white focus:outline-blue-500
        
        ${selectedOption.length > 0 ? "bg-white outline-blue-500 text-blue-500" : "bg-gray-200 outline-gray-200 text-gray-500"}
      `}
      onChange={(e) => onSelect(e.target.value)}
    >
      {["", ...options].map((o, index) => (
        <option key={index} value={o.length > 1 ? o.toLowerCase() : ""}>
          {o.length > 1 ? o : defaultOption}
        </option>
      ))}
    </select>
  );
};

export default Dropdown;
