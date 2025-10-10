/**
 * Tests for useLocalStorage hook
 * Following TDD: These tests MUST FAIL initially (hook not yet implemented)
 */

import { renderHook, act } from '@testing-library/react';
import { useLocalStorage } from './useLocalStorage';

describe('useLocalStorage', () => {
  beforeEach(() => {
    // Clear localStorage before each test
    localStorage.clear();
    jest.clearAllMocks();
  });

  it('should initialize with default value when localStorage is empty', () => {
    const { result } = renderHook(() => useLocalStorage('test-key', 'default-value'));

    expect(result.current[0]).toBe('default-value');
  });

  it('should read existing value from localStorage on mount', () => {
    // Pre-populate localStorage
    localStorage.setItem('test-key', JSON.stringify('stored-value'));

    const { result } = renderHook(() => useLocalStorage('test-key', 'default-value'));

    expect(result.current[0]).toBe('stored-value');
  });

  it('should update localStorage when value changes', () => {
    const { result } = renderHook(() => useLocalStorage<string>('test-key', 'initial'));

    act(() => {
      result.current[1]('updated');
    });

    expect(result.current[0]).toBe('updated');
    expect(localStorage.getItem('test-key')).toBe(JSON.stringify('updated'));
  });

  it('should handle JSON serialization/deserialization', () => {
    const complexValue = { foo: 'bar', nested: { count: 42 } };
    const { result } = renderHook(() => useLocalStorage('test-key', complexValue));

    act(() => {
      result.current[1]({ foo: 'baz', nested: { count: 99 } });
    });

    expect(result.current[0]).toEqual({ foo: 'baz', nested: { count: 99 } });
    const stored = JSON.parse(localStorage.getItem('test-key') || '');
    expect(stored).toEqual({ foo: 'baz', nested: { count: 99 } });
  });

  it('should gracefully degrade if localStorage is unavailable (private browsing)', () => {
    // Mock localStorage.getItem to throw
    const getItemSpy = jest.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new Error('localStorage unavailable');
    });

    const { result } = renderHook(() => useLocalStorage('test-key', 'fallback'));

    // Should fall back to initial value
    expect(result.current[0]).toBe('fallback');

    getItemSpy.mockRestore();
  });

  it('should handle corrupted localStorage values', () => {
    // Store invalid JSON
    localStorage.setItem('test-key', 'not-valid-json{');

    const { result } = renderHook(() => useLocalStorage('test-key', 'fallback'));

    // Should fall back to default value on parse error
    expect(result.current[0]).toBe('fallback');
  });

  it('should support functional updates like setState', () => {
    const { result } = renderHook(() => useLocalStorage<number>('counter', 0));

    act(() => {
      result.current[1]((prev) => prev + 1);
    });

    expect(result.current[0]).toBe(1);

    act(() => {
      result.current[1]((prev) => prev + 10);
    });

    expect(result.current[0]).toBe(11);
  });

  it('should handle localStorage quota exceeded errors', () => {
    const setItemSpy = jest.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
      throw new Error('QuotaExceededError');
    });

    const { result } = renderHook(() => useLocalStorage('test-key', 'initial'));

    // Attempt to update (should not crash)
    act(() => {
      result.current[1]('new-value');
    });

    // State should still update even if localStorage fails
    expect(result.current[0]).toBe('new-value');

    setItemSpy.mockRestore();
  });
});
