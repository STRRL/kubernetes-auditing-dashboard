import React from 'react';
import ReactDiffViewer from 'react-diff-viewer-continued';

interface ResourceDiff {
  added?: string | null;
  removed?: string | null;
  modified: Array<{
    path: string;
    oldValue: string;
    newValue: string;
  }>;
}

interface DiffViewerProps {
  currentState?: string;
  previousState?: string;
}

export const DiffViewer: React.FC<DiffViewerProps> = ({ currentState, previousState }) => {
  if (!currentState || !previousState) return null;

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
