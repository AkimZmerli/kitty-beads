import { useCallback, useRef, useEffect, useState } from 'react';

interface ResizeHandleProps {
  onResize: (deltaY: number) => void;
  minHeight: number;
  maxHeightPercent: number;
}

export function ResizeHandle({ onResize }: ResizeHandleProps) {
  const [isDragging, setIsDragging] = useState(false);
  const lastYRef = useRef<number>(0);

  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    e.preventDefault();
    setIsDragging(true);
    lastYRef.current = e.clientY;
  }, []);

  useEffect(() => {
    if (!isDragging) return;

    const handleMouseMove = (e: MouseEvent) => {
      const deltaY = lastYRef.current - e.clientY;
      lastYRef.current = e.clientY;
      onResize(deltaY);
    };

    const handleMouseUp = () => {
      setIsDragging(false);
    };

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);

    // Prevent text selection while dragging
    document.body.style.userSelect = 'none';
    document.body.style.cursor = 'ns-resize';

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
      document.body.style.userSelect = '';
      document.body.style.cursor = '';
    };
  }, [isDragging, onResize]);

  return (
    <div
      className={`
        absolute top-0 left-0 right-0 h-1 cursor-ns-resize z-10
        transition-colors duration-150
        ${isDragging ? 'bg-grassy-green' : 'hover:bg-grassy-green/50'}
      `}
      onMouseDown={handleMouseDown}
    />
  );
}
