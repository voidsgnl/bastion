import { Link } from 'react-router-dom';
import { api } from '../api/client';
import { useData } from '../api/useData';
import { Page } from '../components/Page';
import { ConfidencePill, SamplePill } from '../components/Pill';
import { impact } from '../components/FixCard';
import type { Confidence } from '../api/types';

const CONF_NOTE: Record<Confidence, string> = {
  Confirmed: 'Verified with evidence',
  Likely: 'Strong signals, not proven',
  Unverified: 'Needs a check',
};

const STATUS_CLASS = { Complete: 'c', Running: 'l', Scheduled: 'u' } as const;

export function Home() {
  const { data, error } = useData(api.home);
  const ai = {
    context: 'Home',
    message: 'Your top fix is the first one in the list. It removes the most risk for the least effort. Want the steps?',
  };
  if (error) return <Page ai={ai}><p role="alert">Could not load data: {error}</p></Page>;
  if (!data) return <Page ai={ai}><p>Loading…</p></Page>;
  const total = data.exposure.internet + data.exposure.internal;
  const pct = (n: number) => `${Math.round((n / total) * 100)}%`;

  return (
    <Page ai={ai}>
      <div className="head">
        <h1>Home</h1>
        <div className="acts">
          <SamplePill />
          <button className="btn" type="button" disabled title="Needs the scan engine">New scan</button>
        </div>
      </div>
      <div className="grid2">
        <section className="card">
          <p className="k">Overall risk</p>
          <p className="big">{data.risk} <span>/ 100</span></p>
          <p className="k">Risk score breakdown is shown on each finding.</p>
        </section>
        <section className="card">
          <p className="k">Exposure</p>
          <p className="lab">Internet-facing: {data.exposure.internet} assets</p>
          <div className="bar"><i style={{ width: pct(data.exposure.internet) }} /></div>
          <p className="lab">Internal: {data.exposure.internal} assets</p>
          <div className="bar"><i className="alt" style={{ width: pct(data.exposure.internal) }} /></div>
        </section>
      </div>
      <section className="card">
        <h2>Top 5 actions</h2>
        <p className="k mb">The fixes that remove the most risk. Start at the top.</p>
        {data.fixes.slice(0, 5).map((f, i) => (
          <div className="row" key={f.id}>
            <b className="num">{i + 1}</b>
            <div className="grow">
              <p className="t">{f.title}</p>
              <p className="k">{impact(f)}</p>
            </div>
            <Link className="btn ghost" to={`/findings/${f.id}`}>View fix</Link>
          </div>
        ))}
      </section>
      <div className="grid2">
        <section className="card">
          <h2>Recent scans</h2>
          {data.scans.map((s) => (
            <div className="row" key={s.name}>
              <span className="grow">{s.name}</span>
              <span className={`pill ${STATUS_CLASS[s.status]}`}>{s.status}</span>
            </div>
          ))}
        </section>
        <section className="card">
          <h2>Findings by confidence</h2>
          {(Object.keys(data.confidence) as Confidence[]).map((c) => (
            <div className="row" key={c}>
              <ConfidencePill value={c} />
              <span className="grow">{CONF_NOTE[c]}</span>
              <b>{data.confidence[c]}</b>
            </div>
          ))}
        </section>
      </div>
    </Page>
  );
}
