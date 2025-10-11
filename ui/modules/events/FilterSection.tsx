import React from 'react';
import { Button } from '@/components/ui/button';

interface FilterSectionProps {
  title: string;
  options: string[];
  selected: string[];
  onToggle: (value: string) => void;
  colorMap?: { [key: string]: string };
}

export const FilterSection: React.FC<FilterSectionProps> = ({
  title,
  options,
  selected,
  onToggle,
  colorMap = {},
}) => {
  return (
    <div className="space-y-2">
      <h3 className="text-sm font-medium text-gray-700">{title}</h3>
      <div className="flex flex-wrap gap-2">
        {options.map((option) => {
          const isSelected = selected.includes(option);
          const color = colorMap[option.toLowerCase()] || 'bg-gray-500';

          return (
            <Button
              key={option}
              variant={isSelected ? "default" : "outline"}
              size="sm"
              onClick={() => onToggle(option)}
              className={isSelected ? `${color} text-white hover:opacity-80` : ''}
            >
              {option}
            </Button>
          );
        })}
      </div>
    </div>
  );
};
