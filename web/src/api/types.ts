export type Confidence = 'Confirmed' | 'Likely' | 'Unverified';
export type FixTab = 'Linux' | 'Windows' | 'Cloud' | 'Code';

export interface Fix {
  id: string;
  title: string;
  removes: number;
  paths: number;
  confidence: Confidence;
}

export interface HomeData {
  risk: number;
  exposure: { internet: number; internal: number };
  fixes: Fix[];
  scans: { name: string; status: 'Complete' | 'Running' | 'Scheduled' }[];
  confidence: Record<Confidence, number>;
}

export interface FindingDetail {
  fixId: string;
  fixTitle: string;
  title: string;
  confidence: Confidence;
  summary: string;
  risk: [string, string][];
  evidence: string;
  fix: Record<FixTab, string>;
  history: [string, string][];
}
