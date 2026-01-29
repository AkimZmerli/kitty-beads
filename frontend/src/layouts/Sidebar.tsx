import { NavLink } from "react-router-dom";
import { useTerminal } from "../features/terminal";
import {
  Map,
  GitBranch,
  Target,
  PenTool,
  Activity,
  Terminal,
  Settings,
  Code2,
  type LucideIcon,
} from "lucide-react";

interface NavItemProps {
  to: string;
  icon: LucideIcon;
  label: string;
}

function NavItem({ to, icon: Icon, label }: NavItemProps) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        `flex items-center gap-3 px-5 py-3 border-l-3 transition-all ${
          isActive
            ? "bg-night-bg-highlight text-neon-cyan border-neon-cyan font-semibold"
            : "border-transparent text-text-secondary hover:bg-night-bg-highlight hover:text-text-primary"
        }`
      }
    >
      {({ isActive }) => (
        <>
          <span className="w-6 flex justify-center">
            <Icon
              className={`w-5 h-5 ${
                isActive ? "drop-shadow-[0_0_6px_rgba(125,207,255,0.5)]" : ""
              }`}
            />
          </span>
          <span>{label}</span>
        </>
      )}
    </NavLink>
  );
}

function TerminalNavItem({
  icon: Icon,
  label,
}: {
  icon: LucideIcon;
  label: string;
}) {
  const { togglePanel, isOpen } = useTerminal();

  return (
    <button
      onClick={togglePanel}
      className={`w-full flex items-center gap-3 px-5 py-3 border-l-3 transition-all ${
        isOpen
          ? "bg-night-bg-highlight text-neon-cyan border-neon-cyan font-semibold"
          : "border-transparent text-text-secondary hover:bg-night-bg-highlight hover:text-text-primary"
      }`}
    >
      <span className="w-6 flex justify-center">
        <Icon
          className={`w-5 h-5 ${
            isOpen ? "drop-shadow-[0_0_6px_rgba(125,207,255,0.5)]" : ""
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
const HEADER_HEIGHT = 73;

export function Sidebar() {
  return (
    <aside
      className="fixed left-0 bottom-0 w-56 bg-sidebar-bg border-r border-night-border py-5 overflow-y-auto z-20"
      style={{ top: `${HEADER_HEIGHT}px` }}
    >
      <SidebarSection title="Views" />
      <NavItem to="/roadmap" icon={Map} label="Roadmap" />
      <NavItem to="/tree" icon={GitBranch} label="Tree" />
      <NavItem to="/kanban" icon={Target} label="Kanban" />
      <NavItem to="/whiteboard" icon={PenTool} label="Whiteboard" />

      <SidebarSection title="Tools" />
      <NavItem to="/editor" icon={Code2} label="Editor" />
      <TerminalNavItem icon={Terminal} label="Terminal" />

      <SidebarSection title="Collaborate" />
      <NavItem to="/activity" icon={Activity} label="Activity" />

      <SidebarSection title="System" />
      <NavItem to="/diagnostics" icon={Settings} label="Diagnostics" />
    </aside>
  );
}
