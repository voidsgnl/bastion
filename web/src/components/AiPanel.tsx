import { FormEvent, useState } from 'react';

interface Msg {
  from: 'ai' | 'you';
  text: string;
}

export function AiPanel({ context, message }: { context: string; message: string }) {
  const [msgs, setMsgs] = useState<Msg[]>([{ from: 'ai', text: message }]);
  const [q, setQ] = useState('');

  function ask(e: FormEvent) {
    e.preventDefault();
    if (!q.trim()) return;
    setMsgs((m) => [
      ...m,
      { from: 'you', text: q },
      { from: 'ai', text: 'The AI service is not connected yet.' },
    ]);
    setQ('');
  }

  return (
    <aside aria-label="AI analyst">
      <h2>AI analyst</h2>
      <span className="pill u self">Context: {context}</span>
      <div className="thread" aria-live="polite">
        {msgs.map((m, i) => (
          <div key={i} className={m.from === 'ai' ? 'msg' : 'msg you'}>
            {m.text}
          </div>
        ))}
      </div>
      <form onSubmit={ask}>
        <label className="k" htmlFor="ask">
          Ask about this page
        </label>
        <input id="ask" value={q} onChange={(e) => setQ(e.target.value)} placeholder="Why is this fix first?" />
      </form>
    </aside>
  );
}
