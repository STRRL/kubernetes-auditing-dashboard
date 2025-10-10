/**
 * Custom React hook for localStorage persistence
 * Based on usehooks-ts pattern with graceful degradation
 */

import { useState, useEffect, Dispatch, SetStateAction } from 'react';

/**
 * Custom hook for persisting state in localStorage
 *
 * @param key - localStorage key to store the value under
 * @param initialValue - Default value if no stored value exists
 * @returns Tuple of [storedValue, setValue] matching React.useState API
 *
 * @example
 * const [hideReadOnly, setHideReadOnly] = useLocalStorage('filter-state', true);
 *
 * Features:
 * - Reads from localStorage on mount
 * - Writes to localStorage on value change
 * - Supports JSON serialization/deserialization
 * - Gracefully degrades if localStorage unavailable (private browsing)
 * - Handles corrupted localStorage values
 * - Supports functional updates like setState
 */
export function useLocalStorage<T>(
  key: string,
  initialValue: T
): [T, Dispatch<SetStateAction<T>>] {
  // State to store our value
  // Pass initial state function to useState so logic is only executed once
  const [storedValue, setStoredValue] = useState<T>(() => {
    if (typeof window === 'undefined') {
      return initialValue;
    }

    try {
      // Get from localStorage by key
      const item = window.localStorage.getItem(key);
      // Parse stored json or return initialValue if null
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      // If error (corrupted data, localStorage unavailable), return initialValue
      console.warn(`Error reading localStorage key "${key}":`, error);
      return initialValue;
    }
  });

  // Return a wrapped version of useState's setter function that
  // persists the new value to localStorage
  const setValue: Dispatch<SetStateAction<T>> = (value) => {
    try {
      // Allow value to be a function so we have same API as useState
      const valueToStore =
        value instanceof Function ? value(storedValue) : value;

      // Save state
      setStoredValue(valueToStore);

      // Save to localStorage
      if (typeof window !== 'undefined') {
        window.localStorage.setItem(key, JSON.stringify(valueToStore));
      }
    } catch (error) {
      // If error (quota exceeded, localStorage unavailable), still update state
      console.warn(`Error setting localStorage key "${key}":`, error);

      // Update state even if localStorage fails (graceful degradation)
      const valueToStore =
        value instanceof Function ? value(storedValue) : value;
      setStoredValue(valueToStore);
    }
  };

  return [storedValue, setValue];
}
