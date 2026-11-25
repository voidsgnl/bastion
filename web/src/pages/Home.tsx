import { Link } from 'react-router-dom';
import { api } from '../api/client';
import { useData } from '../api/useData';
import { Page } from '../components/Page';
import { SamplePill } from '../components/Pill';

const STATUS_CLASS = { Complete: 'complete', Running: 'running', Scheduled: 'scheduled' } as const;

export function Home() {
  const { data, error } = useData(api.home);
  const ai = { context: 'Scans', message: 'I can see the page you’re on. Pick a question below.' };

  if (error) return <Page ai={ai}><p role="alert">Could not load data: {error}</p></Page>;
  if (!data) return <Page ai={ai}><p>Loading…</p></Page>;

  return (
    <Page ai={ai}>
      <div className="head">
        <h1>My Scans</h1>
        <div className="acts">
          <SamplePill />
        </div>
      </div>
      <section className="card scan-table">
        <div className="table-scroll">
          <table>
            <thead><tr><th>Name</th><th>Schedule</th><th>Last Modified</th><th>Status</th></tr></thead>
            <tbody>
              {data.scans.map((scan, index) => {
                const [target, profile] = scan.name.split(',').map((part) => part.trim());
                return (
                  <tr key={scan.name}>
                    <td><Link className="scan-name" to="/findings">{scan.name}</Link><span className="scan-target">{target}</span></td>
                    <td>On demand</td>
                    <td>{index === 0 ? 'Today, 9:12 AM' : index === 1 ? 'Nov 17, 3:40 PM' : 'Nov 16, 2:00 AM'}</td>
                    <td><span className={`scan-status ${STATUS_CLASS[scan.status]}`}><i />{scan.status === 'Complete' ? 'Completed' : scan.status}</span></td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </section>
      <p className="k">Select a completed scan to review its findings and remediation plan.</p>
    </Page>
  );
}
