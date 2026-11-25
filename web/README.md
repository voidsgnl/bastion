# web

The Bastion web app: React, TypeScript, Vite. This is the real UI.

```bash
cd web
npm install
npm run dev      # http://localhost:5173
npm run build
```

## What works

- Home, Results (action plan) and Finding detail, with routing between them
- Light and dark theme (follows the system setting)
- Fix tabs, Copy and Download on the finding page
- AI panel UI. Replies are a placeholder until the AI service exists

## What does not work yet

- All data is sample data from `src/api/mock.ts`. Replace `src/api/client.ts` with calls to the REST API in `api/openapi.yaml`
- Re-test needs the scan engine
- Results grouping by asset and by CVE, and the filters
- Not built: Targets, New scan, Live scan, Attack paths, Compare scans, Reports, Integrations, Settings (they show a placeholder page)

The app uses no external fonts or CDNs, so it works on an air-gapped install.
