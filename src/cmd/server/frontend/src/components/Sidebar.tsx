import { NavLink } from 'react-router-dom';
import type { Feature } from '../types/api';

interface SidebarProps {
  currentFeature: Feature | null;
}

interface NavItemProps {
  to: string;
  icon: string;
  label: string;
  showCheck?: boolean;
  isComplete?: boolean;
}

function NavItem({ to, icon, label, showCheck, isComplete }: NavItemProps) {
  const displayIcon = showCheck ? (isComplete ? '✅' : '⬜') : icon;

  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        `flex items-center gap-3 px-5 py-3 border-l-3 transition-all ${
          isActive
            ? 'bg-green-50 text-grassy-green border-grassy-green font-semibold'
            : 'border-transparent text-text-secondary hover:bg-gray-100'
        }`
      }
    >
      <span className="w-6 text-center">{displayIcon}</span>
      <span>{label}</span>
    </NavLink>
  );
}

function SidebarSection({ title }: { title: string }) {
  return (
    <div className="px-5 pt-4 pb-2 text-xs uppercase tracking-wider text-text-light font-semibold">
      {title}
    </div>
  );
}

export function Sidebar({ currentFeature }: SidebarProps) {
  const artifacts = currentFeature?.artifacts || {};

  return (
    <aside className="w-56 bg-sidebar-bg border-r border-border py-5">
      <SidebarSection title="Workflow" />
      <NavItem to="/" icon="📊" label="Overview" />
      <NavItem to="/specify" icon="📄" label="Specify" showCheck isComplete={artifacts.spec} />
      <NavItem to="/plan" icon="📋" label="Plan" showCheck isComplete={artifacts.plan} />
      <NavItem to="/tasks" icon="📝" label="Tasks" showCheck isComplete={artifacts.tasks} />
      <NavItem to="/kanban" icon="🎯" label="Implement" showCheck isComplete={artifacts.tasks} />

      <SidebarSection title="Artifacts" />
      <NavItem to="/research" icon="🔬" label="Research" />
      <NavItem to="/quickstart" icon="🚀" label="Quickstart" />
      <NavItem to="/data-model" icon="💾" label="Data Model" />

      <SidebarSection title="System" />
      <NavItem to="/terminal" icon="💻" label="Terminal" />
      <NavItem to="/diagnostics" icon="🔍" label="Diagnostics" />
    </aside>
  );
}
