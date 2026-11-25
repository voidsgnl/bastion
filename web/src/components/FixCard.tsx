import { Link } from 'react-router-dom';
import type { Fix } from '../api/types';
import { ConfidencePill } from './Pill';

export function impact(f: Fix) {
  return `Removes ${f.removes} findings` + (f.paths ? `, breaks ${f.paths} attack paths` : '');
}

export function FixCard({ fix }: { fix: Fix }) {
  return (
    <section className="card">
      <div className="head">
        <h2>{fix.title}</h2>
        <ConfidencePill value={fix.confidence} />
      </div>
      <p className="mt">{impact(fix)}</p>
      <div className="acts">
        <Link className="btn" to={`/findings/${fix.id}`}>
          View fix
        </Link>
        <Link className="btn ghost" to={`/findings/${fix.id}`}>
          Show evidence
        </Link>
      </div>
    </section>
  );
}
