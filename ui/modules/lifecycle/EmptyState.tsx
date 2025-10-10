import React from 'react';

export interface EmptyStateProps {
  /** Display variant for different empty state scenarios */
  variant?: 'no-events' | 'filtered-empty';
  /** Optional custom message override */
  message?: string;
}

export const EmptyState: React.FC<EmptyStateProps> = ({
  variant = 'no-events',
  message
}) => {
  // Determine message based on variant or custom override
  const displayMessage = message || (
    variant === 'filtered-empty'
      ? "No events found. Toggle 'Hide read-only events' to see read-only operations."
      : 'This resource has no recorded lifecycle events in the audit history.'
  );

  return (
    <div className="flex flex-col items-center justify-center min-h-[400px] text-center">
      <div className="text-6xl mb-4">📋</div>
      <h3 className="text-xl font-semibold text-gray-700 mb-2">
        {variant === 'filtered-empty' ? 'No events found' : 'No audit event record'}
      </h3>
      <p className="text-gray-500">
        {displayMessage}
      </p>
    </div>
  );
};
