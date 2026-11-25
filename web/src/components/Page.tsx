import { useContext, type ReactNode } from 'react';
import { AiPanel } from './AiPanel';
import { AppShellContext } from './AppShellContext';

export function Page({ ai, children }: { ai: { context: string; message: string }; children: ReactNode }) {
  const { showAi } = useContext(AppShellContext);
  return (
    <>
      <main className="main">{children}</main>
      {showAi && <AiPanel {...ai} />}
    </>
  );
}
