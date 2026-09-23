import React, { createContext, useContext, useState } from "react";

export type UserRole =
  | "PM"
  | "ARCHITECT"
  | "ENGINEER"
  | "SUPPLIER"
  | "LOGISTICS"
  | "SITE_SUPERINTENDENT"
  | "QA";

export interface RoleConfig {
  id: UserRole;
  label: string;
  description: string;
  icon: string;
  channels: string[];
}

export const ROLES: Record<UserRole, RoleConfig> = {
  PM: {
    id: "PM",
    label: "Project Manager",
    description: "Project coordination, stakeholders, decisions & overall execution",
    icon: "manage_accounts",
    channels: ["General", "Design", "Procurement", "Logistics", "Site Operations", "Quality", "Decisions"]
  },
  ARCHITECT: {
    id: "ARCHITECT",
    label: "Architect",
    description: "Technical definition, design specifications & requirements",
    icon: "architecture",
    channels: ["General", "Design", "Procurement"]
  },
  ENGINEER: {
    id: "ENGINEER",
    label: "Engineer",
    description: "Technical assessment, structural integrity & quality reviews",
    icon: "engineering",
    channels: ["General", "Design", "Quality"]
  },
  SUPPLIER: {
    id: "SUPPLIER",
    label: "Supplier",
    description: "Material sourcing, procurement offers & delivery commitments",
    icon: "inventory_2",
    channels: ["Procurement", "Logistics"]
  },
  LOGISTICS: {
    id: "LOGISTICS",
    label: "Logistics Operator",
    description: "Freight transport, route planning, delay reporting & carrier ops",
    icon: "local_shipping",
    channels: ["Logistics", "Site Operations"]
  },
  SITE_SUPERINTENDENT: {
    id: "SITE_SUPERINTENDENT",
    label: "Site Superintendent",
    description: "On-site receiving, field operations & daily superintendent tasks",
    icon: "construction",
    channels: ["General", "Logistics", "Site Operations", "Quality"]
  },
  QA: {
    id: "QA",
    label: "QA / Inspector",
    description: "Quality inspection, weld tests & structural acceptance sign-off",
    icon: "fact_check",
    channels: ["General", "Site Operations", "Quality"]
  }
};

interface RoleContextType {
  activeRole: UserRole;
  setRole: (role: UserRole) => void;
  roleConfig: RoleConfig;
}

const RoleContext = createContext<RoleContextType | undefined>(undefined);

export const RoleProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [activeRole, setActiveRole] = useState<UserRole>("PM");

  return (
    <RoleContext.Provider
      value={{
        activeRole,
        setRole: setActiveRole,
        roleConfig: ROLES[activeRole]
      }}
    >
      {children}
    </RoleContext.Provider>
  );
};

export const useRoleContext = (): RoleContextType => {
  const context = useContext(RoleContext);
  if (!context) {
    throw new Error("useRoleContext must be used within a RoleProvider");
  }
  return context;
};
