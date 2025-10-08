import React from 'react';

interface DiffEntry {
  path: string;
  oldValue: string;
  newValue: string;
}

interface ResourceDiff {
  added?: string | null;
  removed?: string | null;
  modified: DiffEntry[];
}

interface DiffViewerProps {
  diff: ResourceDiff | null | undefined;
}

export const DiffViewer: React.FC<DiffViewerProps> = ({ diff }) => {
  if (!diff) return null;

  const added = diff.added ? JSON.parse(diff.added) : null;
  const removed = diff.removed ? JSON.parse(diff.removed) : null;
  const modified = diff.modified || [];

  const hasChanges = added || removed || modified.length > 0;

  if (!hasChanges) {
    return <div className="text-sm text-gray-500 italic">No changes detected</div>;
  }

  return (
    <div className="mt-2 space-y-2 text-sm">
      {added && (
        <div className="bg-green-50 border-l-4 border-green-500 p-2 rounded">
          <div className="font-semibold text-green-800 mb-1">➕ Added</div>
          <pre className="text-xs text-green-700 overflow-x-auto">
            {JSON.stringify(added, null, 2)}
          </pre>
        </div>
      )}

      {removed && (
        <div className="bg-red-50 border-l-4 border-red-500 p-2 rounded">
          <div className="font-semibold text-red-800 mb-1">➖ Removed</div>
          <pre className="text-xs text-red-700 overflow-x-auto">
            {JSON.stringify(removed, null, 2)}
          </pre>
        </div>
      )}

      {modified.length > 0 && (
        <div className="bg-blue-50 border-l-4 border-blue-500 p-2 rounded">
          <div className="font-semibold text-blue-800 mb-2">✏️ Modified</div>
          <div className="space-y-2">
            {modified.map((entry, idx) => (
              <div key={idx} className="text-xs">
                <div className="font-mono text-blue-700 font-semibold">{entry.path}</div>
                <div className="flex items-center gap-2 ml-2 mt-1">
                  <span className="text-red-600">
                    {JSON.stringify(JSON.parse(entry.oldValue))}
                  </span>
                  <span className="text-gray-500">→</span>
                  <span className="text-green-600">
                    {JSON.stringify(JSON.parse(entry.newValue))}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};
