import React from 'react';
import { MultiSelect } from '@/components/ui/multi-select';

interface FilterSectionProps {
  title: string;
  options: string[];
  selected: string[];
  onChange: (values: string[]) => void;
  colorMap?: { [key: string]: string };
  uppercaseLabels?: boolean;
}

export const FilterSection: React.FC<FilterSectionProps> = ({
  title,
  options,
  selected,
  onChange,
  colorMap = {},
  uppercaseLabels = false,
}) => {
  const multiSelectOptions = options.map((option) => ({
    label: uppercaseLabels ? option.toUpperCase() : option,
    value: option,
    style: colorMap[option.toLowerCase()]
      ? {
        badgeColor: colorMap[option.toLowerCase()],
      }
      : undefined,
  }));

  return (
    <div className="space-y-2 flex-1">
      <h3 className="text-sm font-medium text-gray-700">{title}</h3>
      <div className="">
        <MultiSelect
          className='min-h-[42px]'
          options={multiSelectOptions}
          onValueChange={onChange}
          defaultValue={selected}
          placeholder={`Select ${title.toLowerCase()}...`}
          maxCount={3}
          searchable={true}
        />
      </div>
    </div>
  );
};
