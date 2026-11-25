import { NavLink, Route, Routes } from 'react-router-dom';
import { Home } from '../pages/Home';
import { Results } from '../pages/Results';
import { Finding } from '../pages/Finding';
import { Soon } from '../pages/Soon';

const NAV = [
  { to: '/', label: 'Home', end: true },
  { to: '/targets', label: 'Targets' },
  { to: '/scans', label: 'Scans' },
  { to: '/findings', label: 'Findings' },
  { to: '/attack-paths', label: 'Attack paths' },
  { to: '/reports', label: 'Reports' },
  { to: '/integrations', label: 'Integrations' },
  { to: '/settings', label: 'Settings' },
];

export function App() {
  return (
    <div className="app">
      <nav aria-label="Main">
        <p className="logo">Bastion</p>
        {NAV.map((n) => (
          <NavLink key={n.to} to={n.to} end={n.end}>
            {n.label}
          </NavLink>
        ))}
      </nav>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/findings" element={<Results />} />
        <Route path="/findings/:id" element={<Finding />} />
        <Route path="/targets" element={<Soon title="Targets" />} />
        <Route path="/scans" element={<Soon title="Scans" />} />
        <Route path="/attack-paths" element={<Soon title="Attack paths" />} />
        <Route path="/reports" element={<Soon title="Reports" />} />
        <Route path="/integrations" element={<Soon title="Integrations" />} />
        <Route path="/settings" element={<Soon title="Settings" />} />
        <Route path="*" element={<Soon title="Not found" />} />
      </Routes>
    </div>
  );
}
