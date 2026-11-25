import type { FindingDetail, Fix, HomeData } from './types';
import * as mock from './mock';

// The only place the UI gets data. Swap these for fetch calls to the REST API.
export const api = {
  home: (): Promise<HomeData> => Promise.resolve(mock.home),
  fixes: (): Promise<Fix[]> => Promise.resolve(mock.fixes),
  finding: (id: string): Promise<FindingDetail | null> => Promise.resolve(mock.findings[id] ?? null),
};
