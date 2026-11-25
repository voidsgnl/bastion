import type { Confidence } from '../api/types';
import { SAMPLE_DATA } from '../api/mock';

const KIND: Record<Confidence, string> = { Confirmed: 'c', Likely: 'l', Unverified: 'u' };

export function ConfidencePill({ value }: { value: Confidence }) {
  return <span className={`pill ${KIND[value]}`}>{value}</span>;
}

export function SamplePill() {
  return SAMPLE_DATA ? <span className="pill u">Sample data</span> : null;
}
