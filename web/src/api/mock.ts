import type { Fix, FindingDetail, HomeData } from './types';

// Sample data. Remove when the REST API (api/openapi.yaml) exists.
export const SAMPLE_DATA = true;

export const fixes: Fix[] = [
  { id: 'openssh', title: 'Upgrade OpenSSH on 6 hosts', removes: 11, paths: 5, confidence: 'Confirmed' },
  { id: 'webserver', title: 'Patch the web server on 2 internet-facing hosts', removes: 7, paths: 0, confidence: 'Likely' },
  { id: 'smbv1', title: 'Disable SMBv1 on 4 hosts', removes: 4, paths: 0, confidence: 'Confirmed' },
  { id: 'admin-cred', title: 'Rotate the exposed admin credential', removes: 3, paths: 0, confidence: 'Confirmed' },
  { id: 'db-port', title: 'Close the open database port on 1 host', removes: 2, paths: 0, confidence: 'Confirmed' },
];

export const home: HomeData = {
  risk: 62,
  exposure: { internet: 18, internal: 74 },
  fixes,
  scans: [
    { name: '10.0.0.0/24, Full', status: 'Complete' },
    { name: 'example.com, Web', status: 'Running' },
    { name: 'Cloud account, Safe', status: 'Scheduled' },
  ],
  confidence: { Confirmed: 41, Likely: 23, Unverified: 9 },
};

export const findings: Record<string, FindingDetail> = {
  openssh: {
    fixId: 'openssh',
    fixTitle: 'Upgrade OpenSSH on 6 hosts',
    title: 'Outdated OpenSSH on the login service',
    confidence: 'Confirmed',
    summary:
      'The login service on this host runs an old version with known flaws. Someone who can reach it may be able to take control of the host. Upgrading closes the gap.',
    risk: [
      ['CVSS', 'High'],
      ['EPSS', 'Elevated'],
      ['CISA KEV', 'Listed'],
      ['Exposure', 'Internet-facing'],
      ['Asset role', 'Jump host'],
    ],
    evidence: 'banner: SSH-2.0-OpenSSH_[VERSION]',
    fix: {
      Linux: 'sudo apt update\nsudo apt install --only-upgrade openssh-server',
      Windows: 'Install the latest Win32-OpenSSH release, then restart the sshd service.',
      Cloud: 'Update the machine image to a patched OpenSSH version and redeploy.',
      Code: 'Not applicable to this finding.',
    },
    history: [
      ['First seen', '[DATE]'],
      ['Last seen', '[DATE]'],
      ['Status', 'Open'],
    ],
  },
};
