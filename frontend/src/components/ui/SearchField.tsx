import React, { useState } from "react";

import { Magnifier } from "../../assets/icons";

interface SearchFieldProps {
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
  onPressEnter: () => void;
}

const SearchField: React.FC<SearchFieldProps> = ({
  placeholder,
  value,
  onChange,
  onPressEnter,
}: SearchFieldProps) => {
  const [focused, setFocused] = useState<boolean>(false);

  const hasValue: boolean = value.length > 0;

  return (
    <form
      className="relative w-100"
      onSubmit={(e) => {
        e.preventDefault();
        onPressEnter();
      }}
    >
      <input
        type="text"
        placeholder={placeholder}
        className={`
          p-2 w-full flex flex-row justify-between items-center gap-2
          rounded-sm outline-2
          font-sans font-normal text-blue-500 text-base
          transition-all duration-200 ease-in-out
          focus:bg-white focus:outline-blue-500
          placeholder:font-sans placeholder:font-normal placeholder:text-gray-500 placeholder:text-base
          
          ${hasValue ? "bg-white outline-blue-500" : "bg-gray-200 outline-gray-200"}
        `}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onFocus={() => setFocused(true)}
        onBlur={() => setFocused(false)}
      />
      <span
        className={`
          absolute top-2.5 right-2
          ${focused || hasValue ? "text-blue-500" : "text-gray-500"}
        `}
      >
        <Magnifier size="size-5" strokeWidth={2} />
      </span>
    </form>
  );
};

export default SearchField;
