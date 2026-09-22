import React from "react";

export const ProfilePage: React.FC = () => {
  const team = [
    {
      name: "Jean Dupont",
      role: "Site Manager / Lead Engineer",
      company: "Eiffage Construction",
      contact: "j.dupont@eiffage.fr",
      avatarInitials: "JD",
    },
    {
      name: "Claire Martin",
      role: "Structural Inspector",
      company: "Bureau Veritas",
      contact: "c.martin@bureauveritas.com",
      avatarInitials: "CM",
    },
    {
      name: "Antoine Moreau",
      role: "Procurement Lead",
      company: "BuildFlow Corp",
      contact: "a.moreau@buildflow.io",
      avatarInitials: "AM",
    },
  ];

  return (
    <div className="flex flex-col w-full gap-4 px-4 py-4 max-w-2xl mx-auto pb-8">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-[#181c1e]">Project Stakeholders</h1>
        <span className="text-xs font-semibold text-[#041627] bg-[#041627]/10 px-2.5 py-1 rounded-full">
          3 Active Leads
        </span>
      </div>

      <div className="flex flex-col gap-3">
        {team.map((member, idx) => (
          <div
            key={idx}
            className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex items-center gap-4"
          >
            <div className="w-12 h-12 rounded-full bg-[#041627] text-white flex items-center justify-center font-bold text-sm shrink-0">
              {member.avatarInitials}
            </div>
            <div className="flex flex-col flex-1">
              <h3 className="text-base font-semibold text-[#181c1e]">
                {member.name}
              </h3>
              <span className="text-xs text-[#006c49] font-medium">
                {member.role}
              </span>
              <span className="text-xs text-[#44474c] mt-0.5">
                {member.company} • {member.contact}
              </span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ProfilePage;
