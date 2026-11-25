import { Page } from '../components/Page';

export function Soon({ title }: { title: string }) {
  return (
    <Page ai={{ context: title, message: 'This screen is not built yet.' }}>
      <h1>{title}</h1>
      <p className="k">Not built yet. See docs/10-ui-concept.md for what it will contain.</p>
    </Page>
  );
}
