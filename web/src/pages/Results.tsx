import { useState } from 'react';
import { api } from '../api/client';
import { useData } from '../api/useData';
import { Page } from '../components/Page';
import { SamplePill } from '../components/Pill';
import { FixCard } from '../components/FixCard';

const VIEWS = [
  { id: 'fix', label: 'By fix' },
  { id: 'asset', label: 'By asset' },
  { id: 'cve', label: 'By CVE' },
] as const;

export function Results() {
  const { data, error } = useData(api.fixes);
  const [view, setView] = useState<(typeof VIEWS)[number]['id']>('fix');
  const ai = {
    context: 'Results',
    message: 'Fixes are ordered by how much risk they remove. The first one also cuts attack paths through your jump hosts.',
  };

  return (
    <Page ai={ai}>
      <div className="head">
        <div>
          <h1>Results: action plan</h1>
          <p className="k">10.0.0.0/24, Full profile</p>
        </div>
        <SamplePill />
      </div>
      <div className="tabs" role="group" aria-label="Group results">
        {VIEWS.map((v) => (
          <button key={v.id} type="button" aria-pressed={view === v.id} className={view === v.id ? 'btn' : 'btn ghost'} onClick={() => setView(v.id)}>
            {v.label}
          </button>
        ))}
      </div>
      {error && <p role="alert">Could not load data: {error}</p>}
      {!data && !error && <p>Loading…</p>}
      {data && view === 'fix' && data.map((f) => <FixCard key={f.id} fix={f} />)}
      {data && view !== 'fix' && (
        <section className="card">
          <p>Grouping by {view === 'asset' ? 'asset' : 'CVE'} needs the findings API. Not built yet.</p>
        </section>
      )}
    </Page>
  );
}
