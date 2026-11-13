import React from "react";

import type { IconProps } from "../../types";

const Minus: React.FC<IconProps> = ({ size, strokeWidth }: IconProps) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      strokeWidth={strokeWidth}
      stroke="currentColor"
      className={size}
    >
      <path strokeLinecap="round" strokeLinejoin="round" d="M5 12h14" />
    </svg>
  );
};

export default Minus;
