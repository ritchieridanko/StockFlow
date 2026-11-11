import type { APIErrResponse } from "../types";

export const isAPIErrResponse = (error: unknown): error is APIErrResponse => {
  return (
    typeof error === "object" &&
    error !== null &&
    "status" in error &&
    typeof (error as any).status === "number" &&
    "error" in error &&
    typeof (error as any).error === "object" &&
    "message" in (error as any).error
  );
};
