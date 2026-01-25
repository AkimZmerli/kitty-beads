import { NavLink } from 'react-router-dom';
import { useTerminal } from '../context/TerminalContext';
import type { Feature } from '../types/api';
import {
  LayoutDashboard,
  Map,
  FileText,
  ClipboardList,
  ListTodo,
  Target,
  FlaskConical,
  Rocket,
  Database,
  Terminal,
  Activity,
  CheckCircle2,
  Circle,
  type LucideIcon,
} from 'lucide-react';

interface SidebarProps {
  currentFeature: Feature | null;
}

interface NavItemProps {
  to: string;
  icon: LucideIcon;
  label: string;
  showCheck?: boolean;
  isComplete?: boolean;
}

function NavItem({ to, icon: Icon, label, showCheck, isComplete }: NavItemProps) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        `flex items-center gap-3 px-5 py-3 border-l-3 transition-all ${
          isActive
            ? 'bg-night-bg-highlight text-neon-cyan border-neon-cyan font-semibold'
            : 'border-transparent text-text-secondary hover:bg-night-bg-highlight hover:text-text-primary'
        }`
      }
    >
      {({ isActive }) => (
        <>
          <span className="w-6 flex justify-center">
            {showCheck ? (
              isComplete ? (
                <CheckCircle2 className="w-5 h-5 text-neon-green" />
              ) : (
                <Circle className="w-5 h-5 text-text-muted" />
              )
            ) : (
              <Icon
                className={`w-5 h-5 ${
                  isActive
                    ? 'drop-shadow-[0_0_6px_rgba(125,207,255,0.5)]'
                    : ''
                }`}
              />
            )}
          </span>
          <span>{label}</span>
        </>
      )}
    </NavLink>
  );
}

interface TerminalNavItemProps {
  icon: LucideIcon;
  label: string;
}

function TerminalNavItem({ icon: Icon, label }: TerminalNavItemProps) {
  const { openPanel, isOpen } = useTerminal();

  return (
    <button
      onClick={openPanel}
      className={`w-full flex items-center gap-3 px-5 py-3 border-l-3 transition-all ${
        isOpen
          ? 'bg-night-bg-highlight text-neon-cyan border-neon-cyan font-semibold'
          : 'border-transparent text-text-secondary hover:bg-night-bg-highlight hover:text-text-primary'
      }`}
    >
      <span className="w-6 flex justify-center">
        <Icon
          className={`w-5 h-5 ${
            isOpen ? 'drop-shadow-[0_0_6px_rgba(125,207,255,0.5)]' : ''
          }`}
        />
      </span>
      <span>{label}</span>
    </button>
  );
}

function SidebarSection({ title }: { title: string }) {
  return (
    <div className="px-5 pt-4 pb-2 text-xs uppercase tracking-wider text-text-muted font-semibold">
      {title}
    </div>
  );
}

// Header height constant (must match Header component)
const HEADER_HEIGHT = 73; // px - approximately py-4 + content

export function Sidebar({ currentFeature }: SidebarProps) {
  const artifacts = currentFeature?.artifacts || {};

  return (
    <aside
      className="fixed left-0 bottom-0 w-56 bg-sidebar-bg border-r border-night-border py-5 overflow-y-auto z-20"
      style={{ top: `${HEADER_HEIGHT}px` }}
    >
      <SidebarSection title="Workflow" />
      <NavItem to="/" icon={LayoutDashboard} label="Overview" />
      <NavItem to="/roadmap" icon={Map} label="Roadmap" />
      <NavItem to="/specify" icon={FileText} label="Specify" showCheck isComplete={artifacts.spec} />
      <NavItem to="/plan" icon={ClipboardList} label="Plan" showCheck isComplete={artifacts.plan} />
      <NavItem to="/tasks" icon={ListTodo} label="Tasks" showCheck isComplete={artifacts.tasks} />
      <NavItem to="/kanban" icon={Target} label="Implement" showCheck isComplete={artifacts.tasks} />

      <SidebarSection title="Artifacts" />
      <NavItem to="/research" icon={FlaskConical} label="Research" />
      <NavItem to="/quickstart" icon={Rocket} label="Quickstart" />
      <NavItem to="/data-model" icon={Database} label="Data Model" />

      <SidebarSection title="System" />
      <TerminalNavItem icon={Terminal} label="Terminal" />
      <NavItem to="/diagnostics" icon={Activity} label="Diagnostics" />
    </aside>
  );
}
