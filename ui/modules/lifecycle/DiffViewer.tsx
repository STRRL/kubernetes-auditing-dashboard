import React from 'react';
import ReactDiffViewer from 'react-diff-viewer-continued';

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
  currentState?: string;
  previousState?: string;
}

export const DiffViewer: React.FC<DiffViewerProps> = ({ diff, currentState, previousState }) => {
  if (!diff) return null;

  const added = diff.added ? JSON.parse(diff.added) : null;
  const removed = diff.removed ? JSON.parse(diff.removed) : null;
  const modified = diff.modified || [];

  const hasChanges = added || removed || modified.length > 0;

  if (!hasChanges) {
    return <div className="text-sm text-gray-500 italic">No changes detected</div>;
  }

  if (!currentState || !previousState) {
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
  }

  try {
    const oldObj = JSON.parse(previousState);
    const newObj = JSON.parse(currentState);

    const oldYaml = JSON.stringify(oldObj, null, 2);
    const newYaml = JSON.stringify(newObj, null, 2);

    return (
      <div className="mt-2">
        <ReactDiffViewer
          oldValue={oldYaml}
          newValue={newYaml}
          splitView={true}
          leftTitle="Before"
          rightTitle="After"
          showDiffOnly={true}
          extraLinesSurroundingDiff={3}
          useDarkTheme={false}
          styles={{
            variables: {
              light: {
                diffViewerBackground: '#fff',
                addedBackground: '#e6ffed',
                addedColor: '#24292e',
                removedBackground: '#ffeef0',
                removedColor: '#24292e',
                wordAddedBackground: '#acf2bd',
                wordRemovedBackground: '#fdb8c0',
                addedGutterBackground: '#cdffd8',
                removedGutterBackground: '#ffdce0',
                gutterBackground: '#f6f8fa',
                gutterColor: '#57606a',
                highlightBackground: '#fff8c5',
                highlightGutterBackground: '#fff5b1',
              },
            },
            line: {
              fontSize: '12px',
              fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace',
            },
          }}
        />
      </div>
    );
  } catch (error) {
    console.error('Error parsing states for diff:', error);
    return <div className="text-sm text-red-500">Error displaying diff</div>;
  }
};
