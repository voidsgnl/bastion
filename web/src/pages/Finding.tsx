import { useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { api } from '../api/client';
import { useData } from '../api/useData';
import type { FixTab } from '../api/types';
import { Page } from '../components/Page';
import { ConfidencePill, SamplePill } from '../components/Pill';

const TABS: FixTab[] = ['Linux', 'Windows', 'Cloud', 'Code'];

export function Finding() {
  const { id = '' } = useParams();
  const { data, error } = useData(() => api.finding(id), [id]);
  const [tab, setTab] = useState<FixTab>('Linux');
  const [note, setNote] = useState('');
  const ai = {
    context: 'Finding',
    message: 'This host is the entry point for 5 attack paths. Upgrading it is the single change that breaks all of them.',
  };

  if (error) return <Page ai={ai}><p role="alert">Could not load data: {error}</p></Page>;
  if (!data) {
    return (
      <Page ai={ai}>
        <p className="k"><Link to="/findings">Results</Link></p>
        <p>Loading, or no detail exists for this fix in the sample data.</p>
      </Page>
    );
  }

  const text = data.fix[tab];
  const copy = () =>
    navigator.clipboard.writeText(text).then(
      () => setNote('Copied.'),
      () => setNote('Copy failed. Select the text and copy it by hand.'),
    );
  const download = () => {
    const url = URL.createObjectURL(new Blob([text], { type: 'text/plain' }));
    const a = document.createElement('a');
    a.href = url;
    a.download = `fix-${data.fixId}-${tab.toLowerCase()}.txt`;
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <Page ai={ai}>
      <p className="k"><Link to="/findings">Results</Link> / {data.fixTitle}</p>
      <div className="head">
        <h1>{data.title}</h1>
        <div className="acts"><SamplePill /><ConfidencePill value={data.confidence} /></div>
      </div>
      <p className="lede">{data.summary}</p>
      <div className="grid2">
        <div className="col">
          <section className="card">
            <h2>Risk breakdown</h2>
            {data.risk.map(([k, v]) => (
              <div className="row" key={k}><span className="grow">{k}</span><b>{v}</b></div>
            ))}
          </section>
          <section className="card">
            <h2 className="mb">Evidence</h2>
            <pre className="dark">{data.evidence}</pre>
          </section>
        </div>
        <div className="col">
          <section className="card">
            <h2 className="mb">Fix</h2>
            <div className="tabs mb" role="group" aria-label="Platform">
              {TABS.map((t) => (
                <button key={t} type="button" aria-pressed={tab === t} className={tab === t ? 'btn' : 'btn ghost'} onClick={() => setTab(t)}>{t}</button>
              ))}
            </div>
            <pre className="mb">{text}</pre>
            <div className="acts">
              <button className="btn ghost" type="button" onClick={copy}>Copy</button>
              <button className="btn ghost" type="button" onClick={download}>Download</button>
              <button className="btn" type="button" onClick={() => setNote('Re-test needs the scan engine, which is not connected yet.')}>Re-test</button>
            </div>
            <p className="k mt2" role="status">{note || 'Re-test runs live and shows the result here.'}</p>
          </section>
          <section className="card">
            <h2>History</h2>
            {data.history.map(([k, v]) => (
              <div className="row" key={k}><span className="grow">{k}</span><span>{v}</span></div>
            ))}
          </section>
        </div>
      </div>
    </Page>
  );
}
