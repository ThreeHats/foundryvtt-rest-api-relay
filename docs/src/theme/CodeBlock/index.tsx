import React from 'react';
import OriginalCodeBlock from '@theme-original/CodeBlock';
import InteractiveJson from '@site/src/components/InteractiveJson';

export default function CodeBlock(props: any) {
  // Only enhance JSON code blocks
  if (props.className === 'language-json' || props.metastring?.includes('interactive')) {
    const rawContent = typeof props.children === 'string' ? props.children.trim() : '';
    if (rawContent) {
      try {
        const parsed = JSON.parse(rawContent);
        return <InteractiveJson data={parsed} className="interactive-json--static" />;
      } catch {
        // Fall through to original if invalid JSON
      }
    }
  }

  return <OriginalCodeBlock {...props} />;
}
