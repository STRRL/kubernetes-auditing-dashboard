import React from 'react';

interface TimestampDisplayProps {
  timestamp: string;
}

export const TimestampDisplay: React.FC<TimestampDisplayProps> = ({ timestamp }) => {
  const date = new Date(timestamp);

  const formatted = new Intl.DateTimeFormat('en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date);

  return <span className="text-sm text-gray-700">{formatted}</span>;
};
