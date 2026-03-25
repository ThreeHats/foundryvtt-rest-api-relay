import React, { useState, useCallback, useRef } from 'react';

interface InteractiveJsonProps {
  data: any;
  className?: string;
}

function CopiedToast({ text, position }: { text: string; position: { x: number; y: number } }) {
  return (
    <span
      className="interactive-json__toast"
      style={{
        left: position.x,
        top: position.y,
      }}
    >
      Copied: {text.length > 40 ? text.slice(0, 40) + '...' : text}
    </span>
  );
}

function JsonValue({
  value,
  path,
  onCopy,
}: {
  value: any;
  path: string;
  onCopy: (text: string, e: React.MouseEvent) => void;
}) {
  if (value === null) {
    return (
      <span
        className="interactive-json__value interactive-json__null"
        title={`Click to copy value`}
        onClick={(e) => { e.stopPropagation(); onCopy('null', e); }}
      >
        null
      </span>
    );
  }

  if (typeof value === 'boolean') {
    return (
      <span
        className="interactive-json__value interactive-json__boolean"
        title={`Click to copy value`}
        onClick={(e) => { e.stopPropagation(); onCopy(String(value), e); }}
      >
        {String(value)}
      </span>
    );
  }

  if (typeof value === 'number') {
    return (
      <span
        className="interactive-json__value interactive-json__number"
        title={`Click to copy value`}
        onClick={(e) => { e.stopPropagation(); onCopy(String(value), e); }}
      >
        {String(value)}
      </span>
    );
  }

  if (typeof value === 'string') {
    return (
      <span
        className="interactive-json__value interactive-json__string"
        title={`Click to copy value`}
        onClick={(e) => { e.stopPropagation(); onCopy(value, e); }}
      >
        "{value}"
      </span>
    );
  }

  return null;
}

function JsonNode({
  keyName,
  value,
  path,
  isLast,
  onCopy,
  depth,
}: {
  keyName?: string | number;
  value: any;
  path: string;
  isLast: boolean;
  onCopy: (text: string, e: React.MouseEvent) => void;
  depth: number;
}) {
  const [collapsed, setCollapsed] = useState(depth >= 3);
  const indent = '  '.repeat(depth);
  const childIndent = '  '.repeat(depth + 1);
  const comma = isLast ? '' : ',';

  const isArray = Array.isArray(value);
  const isObject = value !== null && typeof value === 'object' && !isArray;
  const isComplex = isArray || isObject;

  const renderKey = () => {
    if (keyName === undefined) return null;
    const displayKey = typeof keyName === 'string' ? `"${keyName}"` : String(keyName);
    return (
      <>
        <span
          className="interactive-json__key"
          title={path}
          onClick={(e) => { e.stopPropagation(); onCopy(path, e); }}
        >
          {displayKey}
        </span>
        <span className="interactive-json__punct">: </span>
      </>
    );
  };

  if (!isComplex) {
    return (
      <div className="interactive-json__line">
        {indent}{renderKey()}
        <JsonValue value={value} path={path} onCopy={onCopy} />
        <span className="interactive-json__punct">{comma}</span>
      </div>
    );
  }

  const entries = isArray ? value : Object.keys(value);
  const openBrace = isArray ? '[' : '{';
  const closeBrace = isArray ? ']' : '}';

  if (entries.length === 0) {
    return (
      <div className="interactive-json__line">
        {indent}{renderKey()}
        <span className="interactive-json__punct">{openBrace}{closeBrace}{comma}</span>
      </div>
    );
  }

  if (collapsed) {
    const count = entries.length;
    const label = isArray ? `${count} item${count !== 1 ? 's' : ''}` : `${count} key${count !== 1 ? 's' : ''}`;
    return (
      <div className="interactive-json__line">
        {indent}{renderKey()}
        <span
          className="interactive-json__collapse-toggle"
          onClick={(e) => { e.stopPropagation(); setCollapsed(false); }}
          title="Click to expand"
        >
          {openBrace} <span className="interactive-json__collapsed-label">{label}</span> {closeBrace}
        </span>
        <span className="interactive-json__punct">{comma}</span>
      </div>
    );
  }

  return (
    <>
      <div className="interactive-json__line">
        {indent}{renderKey()}
        <span
          className="interactive-json__collapse-toggle"
          onClick={(e) => { e.stopPropagation(); setCollapsed(true); }}
          title="Click to collapse"
        >
          {openBrace}
        </span>
      </div>
      {isArray
        ? value.map((item: any, i: number) => (
            <JsonNode
              key={i}
              keyName={i}
              value={item}
              path={`${path}[${i}]`}
              isLast={i === value.length - 1}
              onCopy={onCopy}
              depth={depth + 1}
            />
          ))
        : Object.keys(value).map((key, i, arr) => (
            <JsonNode
              key={key}
              keyName={key}
              value={value[key]}
              path={path ? `${path}.${key}` : key}
              isLast={i === arr.length - 1}
              onCopy={onCopy}
              depth={depth + 1}
            />
          ))}
      <div className="interactive-json__line">
        {indent}
        <span className="interactive-json__punct">{closeBrace}{comma}</span>
      </div>
    </>
  );
}

export default function InteractiveJson({ data, className = '' }: InteractiveJsonProps) {
  const [toast, setToast] = useState<{ text: string; position: { x: number; y: number } } | null>(null);
  const containerRef = useRef<HTMLPreElement>(null);
  const toastTimeout = useRef<ReturnType<typeof setTimeout>>();

  const handleCopy = useCallback((text: string, e: React.MouseEvent) => {
    navigator.clipboard.writeText(text).catch(() => {});

    // Position toast near cursor, relative to the container
    const rect = containerRef.current?.getBoundingClientRect();
    if (rect) {
      setToast({
        text,
        position: {
          x: e.clientX - rect.left,
          y: e.clientY - rect.top - 28,
        },
      });
    }

    if (toastTimeout.current) clearTimeout(toastTimeout.current);
    toastTimeout.current = setTimeout(() => setToast(null), 1500);
  }, []);

  // If data is a string (non-JSON text response), just render it plain
  if (typeof data === 'string') {
    return <pre className={`interactive-json ${className}`}>{data}</pre>;
  }

  return (
    <pre className={`interactive-json ${className}`} ref={containerRef}>
      {toast && <CopiedToast text={toast.text} position={toast.position} />}
      <code className="interactive-json__code">
        <JsonNode
          value={data}
          path=""
          isLast={true}
          onCopy={handleCopy}
          depth={0}
        />
      </code>
    </pre>
  );
}
