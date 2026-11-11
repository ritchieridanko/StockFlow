import React from "react";

interface TextFieldProps {
  type: "text" | "number";
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
  min?: number;
  max?: number;
  step?: number;
}

const TextField: React.FC<TextFieldProps> = ({
  type,
  placeholder,
  value,
  onChange,
  min,
  max,
  step,
}: TextFieldProps) => {
  return (
    <input
      type={type}
      placeholder={placeholder}
      className={`
        px-2 py-1 w-full flex flex-row justify-between items-center gap-2
        rounded-sm outline-2
        font-sans font-normal text-blue-500 text-base
        transition-all duration-200 ease-in-out
        focus:bg-white focus:outline-blue-500
        placeholder:font-sans placeholder:font-normal placeholder:text-gray-700 placeholder:text-base
        
        ${value.length > 0 ? "bg-white outline-blue-500" : "bg-gray-200 outline-gray-200"}
      `}
      value={value}
      onChange={(e) => onChange(e.target.value)}
      min={min}
      max={max}
      step={step}
    />
  );
};

export default TextField;
