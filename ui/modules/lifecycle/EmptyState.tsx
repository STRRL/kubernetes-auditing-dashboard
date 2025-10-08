import React from 'react';

export const EmptyState: React.FC = () => {
  return (
    <div className="flex flex-col items-center justify-center min-h-[400px] text-center">
      <div className="text-6xl mb-4">📋</div>
      <h3 className="text-xl font-semibold text-gray-700 mb-2">No audit event record</h3>
      <p className="text-gray-500">
        This resource has no recorded lifecycle events in the audit history.
      </p>
    </div>
  );
};
