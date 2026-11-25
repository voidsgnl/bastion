import type { ReactNode } from 'react';
import { AiPanel } from './AiPanel';

export function Page({ ai, children }: { ai: { context: string; message: string }; children: ReactNode }) {
  return (
    <>
      <main>{children}</main>
      <AiPanel {...ai} />
    </>
  );
}
