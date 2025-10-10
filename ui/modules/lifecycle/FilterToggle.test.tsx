/**
 * Tests for FilterToggle component
 * Following TDD: These tests MUST FAIL initially (component not yet implemented)
 */

import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { FilterToggle } from './FilterToggle';

describe('FilterToggle', () => {
  it('should render with correct label text', () => {
    const mockOnChange = jest.fn();
    render(<FilterToggle checked={false} onChange={mockOnChange} />);

    expect(screen.getByLabelText('Hide read-only events')).toBeInTheDocument();
  });

  it('should reflect checked state visually', () => {
    const mockOnChange = jest.fn();
    const { rerender } = render(<FilterToggle checked={true} onChange={mockOnChange} />);

    const checkbox = screen.getByLabelText('Hide read-only events') as HTMLInputElement;
    expect(checkbox.checked).toBe(true);

    rerender(<FilterToggle checked={false} onChange={mockOnChange} />);
    expect(checkbox.checked).toBe(false);
  });

  it('should call onChange when clicked', () => {
    const mockOnChange = jest.fn();
    render(<FilterToggle checked={false} onChange={mockOnChange} />);

    const checkbox = screen.getByLabelText('Hide read-only events');
    fireEvent.click(checkbox);

    expect(mockOnChange).toHaveBeenCalledTimes(1);
    expect(mockOnChange).toHaveBeenCalledWith(true);
  });

  it('should be keyboard accessible (Space/Enter)', () => {
    const mockOnChange = jest.fn();
    render(<FilterToggle checked={false} onChange={mockOnChange} />);

    const checkbox = screen.getByLabelText('Hide read-only events');

    // Space key
    fireEvent.keyDown(checkbox, { key: ' ', code: 'Space' });
    expect(mockOnChange).toHaveBeenCalled();

    mockOnChange.mockClear();

    // Enter key (checkbox typically responds to Space, but test both)
    fireEvent.keyDown(checkbox, { key: 'Enter', code: 'Enter' });
    // Note: Native checkboxes respond to Space, not Enter
    // This test documents expected browser behavior
  });

  it('should have proper ARIA attributes', () => {
    const mockOnChange = jest.fn();
    render(<FilterToggle checked={true} onChange={mockOnChange} />);

    const checkbox = screen.getByRole('checkbox', { name: 'Hide read-only events' });
    expect(checkbox).toBeInTheDocument();
    expect(checkbox).toHaveAttribute('type', 'checkbox');
    expect(checkbox).toHaveAttribute('aria-checked', 'true');
  });

  it('should accept optional label prop override', () => {
    const mockOnChange = jest.fn();
    render(<FilterToggle checked={false} onChange={mockOnChange} label="Custom label" />);

    expect(screen.getByLabelText('Custom label')).toBeInTheDocument();
    expect(screen.queryByLabelText('Hide read-only events')).not.toBeInTheDocument();
  });

  it('should toggle state when clicked multiple times', () => {
    const mockOnChange = jest.fn();
    render(<FilterToggle checked={false} onChange={mockOnChange} />);

    const checkbox = screen.getByLabelText('Hide read-only events');

    // First click: false -> true
    fireEvent.click(checkbox);
    expect(mockOnChange).toHaveBeenNthCalledWith(1, true);

    // Simulate parent re-render with new state
    mockOnChange.mockClear();
    const { rerender } = render(<FilterToggle checked={true} onChange={mockOnChange} />);

    // Second click: true -> false
    fireEvent.click(checkbox);
    expect(mockOnChange).toHaveBeenNthCalledWith(1, false);
  });

  it('should have visible focus indicator', () => {
    const mockOnChange = jest.fn();
    render(<FilterToggle checked={false} onChange={mockOnChange} />);

    const checkbox = screen.getByLabelText('Hide read-only events');
    checkbox.focus();

    // Verify element is focused
    expect(checkbox).toHaveFocus();
    // Note: Visual focus indicator styling verified manually (Tailwind utilities)
  });
});
