import { useState } from 'react';
import { NavLink, Route, Routes, useLocation } from 'react-router-dom';
import { AppShellContext } from '../components/AppShellContext';
import { Home } from '../pages/Home';
import { Finding } from '../pages/Finding';
import { Results } from '../pages/Results';
import { Soon } from '../pages/Soon';

const TABS = [
  { to: '/', label: 'Scans', end: true },
  { to: '/targets', label: 'Targets' },
  { to: '/reports', label: 'Reports' },
  { to: '/settings', label: 'Settings' },
];

export function App() {
  const [showAi, setShowAi] = useState(false);
  const [theme, setTheme] = useState<'light' | 'dark' | null>(null);
  const location = useLocation();
  const inFindings = location.pathname.startsWith('/findings');

  function toggleTheme() {
    const next = theme ?? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
    const value = next === 'dark' ? 'light' : 'dark';
    document.documentElement.dataset.theme = value;
    setTheme(value);
  }

  return (
    <AppShellContext.Provider value={{ showAi }}>
      <div className="app">
        <header className="topbar">
          <NavLink className="brand" to="/" aria-label="Bastion scans">🛡 Bastion</NavLink>
          <nav className="top-tabs" aria-label="Main navigation">
            {TABS.map((tab) => <NavLink key={tab.to} to={tab.to} end={tab.end}>{tab.label}</NavLink>)}
          </nav>
          <div className="top-actions">
            <button className="top-button" type="button" onClick={toggleTheme} aria-label="Toggle color theme">◐</button>
            <button className={`top-button ${showAi ? 'selected' : ''}`} type="button" onClick={() => setShowAi((open) => !open)} aria-expanded={showAi}>✦ Ask AI</button>
            <button className="top-button primary" type="button" disabled title="Scan engine is not connected">＋ New Scan</button>
          </div>
        </header>
        <div className="workspace">
          <aside className="side-nav" aria-label="Workspace navigation">
            <h2>FOLDERS</h2>
            <NavLink to="/" end><span>My Scans</span><span className="side-count">3</span></NavLink>
            <NavLink to="/scans"><span>All Scans</span><span className="side-count">3</span></NavLink>
            <button type="button" className="side-link" disabled><span>Trash</span><span className="side-count">0</span></button>
            <h2>RESOURCES</h2>
            <NavLink to="/settings">Profiles</NavLink>
            <NavLink to="/settings">Engines</NavLink>
            <NavLink to="/settings">MCP tokens</NavLink>
            {inFindings && <><h2>FINDINGS</h2><NavLink to="/findings">Action plan</NavLink></>}
          </aside>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/findings" element={<Results />} />
            <Route path="/findings/:id" element={<Finding />} />
            <Route path="/targets" element={<Soon title="Targets" />} />
            <Route path="/scans" element={<Soon title="All scans" />} />
            <Route path="/attack-paths" element={<Soon title="Attack paths" />} />
            <Route path="/reports" element={<Soon title="Reports" />} />
            <Route path="/integrations" element={<Soon title="Integrations" />} />
            <Route path="/settings" element={<Soon title="Settings" />} />
            <Route path="*" element={<Soon title="Not found" />} />
          </Routes>
        </div>
      </div>
    </AppShellContext.Provider>
  );
}
