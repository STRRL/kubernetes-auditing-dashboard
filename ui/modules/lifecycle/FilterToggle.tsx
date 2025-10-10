/**
 * FilterToggle Component
 * Checkbox control for toggling read-only event visibility
 */

import React from 'react';

export interface FilterToggleProps {
  /** Whether the toggle is checked (read-only events hidden) */
  checked: boolean;
  /** Callback when toggle state changes */
  onChange: (checked: boolean) => void;
  /** Optional label override (defaults to "Hide read-only events") */
  label?: string;
}

/**
 * Accessible checkbox toggle for filtering read-only events
 *
 * @example
 * <FilterToggle
 *   checked={hideReadOnly}
 *   onChange={setHideReadOnly}
 * />
 */
export function FilterToggle({
  checked,
  onChange,
  label = 'Hide read-only events',
}: FilterToggleProps) {
  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange(event.target.checked);
  };

  return (
    <div className="flex items-center gap-2">
      <input
        type="checkbox"
        id="filter-toggle"
        checked={checked}
        onChange={handleChange}
        aria-checked={checked}
        className="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 cursor-pointer"
      />
      <label
        htmlFor="filter-toggle"
        className="text-sm font-medium text-gray-700 cursor-pointer select-none"
      >
        {label}
      </label>
    </div>
  );
}
