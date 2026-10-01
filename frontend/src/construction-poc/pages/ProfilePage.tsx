import React, { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { CONSTRUCTION_ROUTES } from "../constants/routes";

export const ProfilePage: React.FC = () => {
  const navigate = useNavigate();

  // Field switches local state
  const [offlineSync, setOfflineSync] = useState<boolean>(true);
  const [criticalLogistics, setCriticalLogistics] = useState<boolean>(true);
  const [highContrast, setHighContrast] = useState<boolean>(false);
  const [satelliteMode, setSatelliteMode] = useState<boolean>(false);

  // Synchronization simulation state
  const [isSyncing, setIsSyncing] = useState<boolean>(false);
  const [syncComplete, setSyncComplete] = useState<boolean>(false);
  const [lastCloudSyncText, setLastCloudSyncText] = useState<string>("2m ago");

  const syncTimerRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    return () => {
      if (syncTimerRef.current) {
        clearTimeout(syncTimerRef.current);
      }
    };
  }, []);

  const handleSyncClick = () => {
    if (isSyncing) return;
    setIsSyncing(true);
    setSyncComplete(false);

    syncTimerRef.current = setTimeout(() => {
      setIsSyncing(false);
      setSyncComplete(true);
      setLastCloudSyncText("Just now");

      // Reset button success text after 3 seconds while keeping Last Cloud Sync as "Just now"
      syncTimerRef.current = setTimeout(() => {
        setSyncComplete(false);
      }, 3000);
    }, 1500);
  };

  const credentials = [
    {
      id: "cred-1",
      title: "Eurocodes EN 1992-1 / 1993",
      authority: "Structural Concrete & Steel Authority",
      status: "Verified",
      statusStyle: "bg-emerald-50 text-emerald-700 border-emerald-200",
      detail: "Full Sign-off Privilege",
      validity: "Valid to Dec 2026",
      icon: "verified",
    },
    {
      id: "cred-2",
      title: "Bureau Inspection Assessor",
      authority: "Level 3 Certified Technical Auditor",
      status: "Accredited",
      statusStyle: "bg-blue-50 text-blue-700 border-blue-200",
      detail: "Audit Registry: ID-8840-X",
      validity: "Annual Renewal Sync OK",
      icon: "fact_check",
    },
    {
      id: "cred-3",
      title: "Site Safety Pass: Category A",
      authority: "High Risk Infrastructure / Viaduct Zones",
      status: "Site-Clear",
      statusStyle: "bg-emerald-50 text-emerald-700 border-emerald-200",
      detail: "Site-Clear",
      validity: "All Sectors Authorized",
      icon: "health_and_safety",
    },
  ];

  const operationalAccess = [
    {
      id: "acc-1",
      category: "Project",
      name: "PRJ-001 Metro Line 4",
      description: "Full Engineering Sign-off & Inspection Authority",
      icon: "architecture",
      tag: null,
    },
    {
      id: "acc-2",
      category: "Site",
      name: "SITE-001 Riverside Tower",
      description: "Zone Access: Viaduct & Deep Foundation Sector",
      icon: "location_on",
      tag: null,
    },
    {
      id: "acc-3",
      category: "TraceCore",
      name: "TraceCore Milestone Notarization",
      description: "Active Cryptographic Delegated Signer",
      icon: "lock",
      tag: "HSM Tier 1",
    },
    {
      id: "acc-4",
      category: "C3",
      name: "C3 Collaboration Feeds",
      description: "#Logistics, #Foundation, #Technical-Review",
      icon: "hub",
      tag: "3 Linked",
    },
  ];

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6 space-y-6 pb-20">
      {/* User Profile Card */}
      <div className="p-6 rounded-2xl bg-white border border-[#e0e3e5] shadow-sm space-y-6">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            {/* Profile Avatar with Online Indicator */}
            <div className="relative shrink-0">
              <div className="w-16 h-16 rounded-full bg-[#041627] text-white flex items-center justify-center font-bold text-xl ring-4 ring-[#041627]/10 shadow-md">
                MV
              </div>
              <span
                className="w-4 h-4 rounded-full bg-emerald-500 border-2 border-white absolute bottom-0 right-0 shadow-sm"
                title="Online"
                aria-label="Online indicator"
              />
            </div>

            {/* Identity Info */}
            <div>
              <div className="flex items-center gap-2">
                <h1 className="text-xl font-bold text-[#041627]">Manuel Vincent</h1>
                <span
                  className="material-symbols-outlined text-emerald-600 text-[20px]"
                  title="Verified User"
                  aria-label="Verified user icon"
                >
                  verified
                </span>
              </div>
              <p className="text-sm font-semibold text-[#006c49]">Lead Structural Engineer</p>
              <p className="text-xs text-[#44474c] mt-0.5 font-medium">
                Partner • Engineering Partners
              </p>
            </div>
          </div>

          {/* Edit Profile Affordance */}
          <button
            onClick={() => alert("Profile editing functionality is presentation-only in frontend mode.")}
            className="self-start sm:self-auto p-2.5 rounded-full border border-[#e0e3e5] bg-white text-[#44474c] hover:text-[#041627] hover:bg-[#f0f4f8] transition-colors shadow-sm"
            aria-label="Edit Profile"
            title="Edit Profile"
          >
            <span className="material-symbols-outlined text-[18px]">edit</span>
          </button>
        </div>

        {/* Metadata Cards: Workspace ID & Active Site */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-4 border-t border-[#e0e3e5]">
          <div className="p-3.5 rounded-xl bg-[#f7fafc] border border-[#e0e3e5] flex items-center justify-between">
            <div>
              <span className="text-[10px] uppercase font-bold tracking-wider text-[#707478]">
                Workspace ID
              </span>
              <p className="text-sm font-bold text-[#041627] mt-0.5">SW-FR-7501</p>
            </div>
            <span className="px-2.5 py-1 bg-white border border-[#e0e3e5] text-[#041627] text-xs font-bold rounded-lg shadow-2xs">
              Mesh
            </span>
          </div>

          <div className="p-3.5 rounded-xl bg-[#f7fafc] border border-[#e0e3e5] flex items-center justify-between">
            <div>
              <span className="text-[10px] uppercase font-bold tracking-wider text-[#707478]">
                Active Site
              </span>
              <p className="text-sm font-bold text-[#041627] mt-0.5">PRJ-001</p>
            </div>
            <span className="px-2.5 py-1 bg-[#041627] text-white text-xs font-bold rounded-lg shadow-2xs">
              Metro Line 4
            </span>
          </div>
        </div>
      </div>

      {/* Credentials & Sign-off Section */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h2 className="text-base font-bold text-[#041627]">Credentials &amp; Sign-off</h2>
          <span className="px-2.5 py-0.5 bg-[#006c49]/10 text-[#006c49] text-xs font-bold rounded-full">
            3 Active
          </span>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
          {credentials.map((cred) => (
            <div
              key={cred.id}
              className="p-4 rounded-xl bg-white border border-[#e0e3e5] shadow-sm flex flex-col justify-between space-y-3"
            >
              <div className="space-y-2">
                <div className="flex items-start justify-between gap-2">
                  <span className="material-symbols-outlined text-[#041627] text-[22px]">
                    {cred.icon}
                  </span>
                  <span
                    className={`px-2 py-0.5 text-[10px] font-extrabold rounded border ${cred.statusStyle}`}
                  >
                    {cred.status}
                  </span>
                </div>
                <div>
                  <h3 className="text-sm font-bold text-[#041627] leading-snug">{cred.title}</h3>
                  <p className="text-xs text-[#44474c] mt-0.5 leading-tight">{cred.authority}</p>
                </div>
              </div>

              <div className="pt-2 border-t border-[#f0f4f8] text-[11px] font-semibold text-[#041627] space-y-0.5">
                <p>{cred.detail}</p>
                <p className="text-[#707478] font-normal">{cred.validity}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Operational Access Section */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h2 className="text-base font-bold text-[#041627]">Operational Access</h2>
          <span className="px-2.5 py-0.5 bg-[#041627]/10 text-[#041627] text-xs font-bold rounded-full">
            Sovereign Rights
          </span>
        </div>

        <div className="bg-white border border-[#e0e3e5] rounded-xl shadow-sm divide-y divide-[#e0e3e5]">
          {operationalAccess.map((acc) => (
            <div key={acc.id} className="p-4 flex items-center justify-between gap-4">
              <div className="flex items-center gap-3">
                <div className="w-9 h-9 rounded-lg bg-[#f0f4f8] text-[#041627] flex items-center justify-center shrink-0">
                  <span className="material-symbols-outlined text-[20px]">{acc.icon}</span>
                </div>
                <div>
                  <span className="text-[10px] font-extrabold uppercase tracking-wider text-[#707478]">
                    {acc.category}
                  </span>
                  <h3 className="text-sm font-bold text-[#041627]">{acc.name}</h3>
                  <p className="text-xs text-[#44474c] font-medium">{acc.description}</p>
                </div>
              </div>

              {acc.tag && (
                <span className="px-2.5 py-1 bg-[#041627] text-white text-[11px] font-bold rounded-md shrink-0">
                  {acc.tag}
                </span>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Field Configuration Section */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h2 className="text-base font-bold text-[#041627]">Field Configuration</h2>
          <span className="px-2.5 py-0.5 bg-[#f0f4f8] text-[#041627] text-xs font-bold rounded-full border border-[#e0e3e5]">
            Profile v2.4
          </span>
        </div>

        <div className="bg-white border border-[#e0e3e5] rounded-xl shadow-sm divide-y divide-[#e0e3e5]">
          {/* Switch 1: Offline Auto-Sync */}
          <div className="p-4 flex items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="w-9 h-9 rounded-lg bg-[#f0f4f8] text-[#041627] flex items-center justify-center shrink-0">
                <span className="material-symbols-outlined text-[20px]">sync_desktop</span>
              </div>
              <div>
                <h3 className="text-sm font-bold text-[#041627]">Offline Auto-Sync</h3>
                <p className="text-xs text-[#44474c] font-medium">
                  Automatic telemetry caching via dual WiFi + 5G failover
                </p>
              </div>
            </div>

            <button
              type="button"
              role="switch"
              aria-checked={offlineSync}
              aria-label="Offline Auto-Sync"
              onClick={() => setOfflineSync(!offlineSync)}
              className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-[#041627] focus:ring-offset-2 ${
                offlineSync ? "bg-[#041627]" : "bg-[#d0d4d8]"
              }`}
            >
              <span
                aria-hidden="true"
                className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                  offlineSync ? "translate-x-5" : "translate-x-0"
                }`}
              />
            </button>
          </div>

          {/* Switch 2: Critical Logistics Dispatch */}
          <div className="p-4 flex items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="w-9 h-9 rounded-lg bg-[#f0f4f8] text-[#041627] flex items-center justify-center shrink-0">
                <span className="material-symbols-outlined text-[20px]">local_shipping</span>
              </div>
              <div>
                <h3 className="text-sm font-bold text-[#041627]">Critical Logistics Dispatch</h3>
                <p className="text-xs text-[#44474c] font-medium">
                  High-priority alerts (DEL-1042 batching &amp; corridor delays)
                </p>
              </div>
            </div>

            <button
              type="button"
              role="switch"
              aria-checked={criticalLogistics}
              aria-label="Critical Logistics Dispatch"
              onClick={() => setCriticalLogistics(!criticalLogistics)}
              className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-[#041627] focus:ring-offset-2 ${
                criticalLogistics ? "bg-[#041627]" : "bg-[#d0d4d8]"
              }`}
            >
              <span
                aria-hidden="true"
                className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                  criticalLogistics ? "translate-x-5" : "translate-x-0"
                }`}
              />
            </button>
          </div>

          {/* Switch 3: High-Contrast Sunlight Mode */}
          <div className="p-4 flex items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="w-9 h-9 rounded-lg bg-[#f0f4f8] text-[#041627] flex items-center justify-center shrink-0">
                <span className="material-symbols-outlined text-[20px]">contrast</span>
              </div>
              <div>
                <h3 className="text-sm font-bold text-[#041627]">High-Contrast Sunlight Mode</h3>
                <p className="text-xs text-[#44474c] font-medium">
                  Maximizes blueprint line weight under direct midday sun
                </p>
              </div>
            </div>

            <button
              type="button"
              role="switch"
              aria-checked={highContrast}
              aria-label="High-Contrast Sunlight Mode"
              onClick={() => setHighContrast(!highContrast)}
              className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-[#041627] focus:ring-offset-2 ${
                highContrast ? "bg-[#041627]" : "bg-[#d0d4d8]"
              }`}
            >
              <span
                aria-hidden="true"
                className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                  highContrast ? "translate-x-5" : "translate-x-0"
                }`}
              />
            </button>
          </div>

          {/* Switch 4: Satellite / Low-Bandwidth Mode */}
          <div className="p-4 flex items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="w-9 h-9 rounded-lg bg-[#f0f4f8] text-[#041627] flex items-center justify-center shrink-0">
                <span className="material-symbols-outlined text-[20px]">satellite_alt</span>
              </div>
              <div>
                <h3 className="text-sm font-bold text-[#041627]">Satellite / Low-Bandwidth Mode</h3>
                <p className="text-xs text-[#44474c] font-medium">
                  Optimizes packets for remote sub-terrain tunnel inspection
                </p>
              </div>
            </div>

            <button
              type="button"
              role="switch"
              aria-checked={satelliteMode}
              aria-label="Satellite / Low-Bandwidth Mode"
              onClick={() => setSatelliteMode(!satelliteMode)}
              className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-[#041627] focus:ring-offset-2 ${
                satelliteMode ? "bg-[#041627]" : "bg-[#d0d4d8]"
              }`}
            >
              <span
                aria-hidden="true"
                className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                  satelliteMode ? "translate-x-5" : "translate-x-0"
                }`}
              />
            </button>
          </div>
        </div>
      </div>

      {/* System Health & Storage Section */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h2 className="text-base font-bold text-[#041627]">System Health &amp; Storage</h2>
          <span className="px-2.5 py-0.5 bg-emerald-50 text-emerald-700 text-xs font-bold rounded-full border border-emerald-200">
            Optimal
          </span>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div className="p-4 rounded-xl bg-white border border-[#e0e3e5] shadow-sm flex items-center justify-between">
            <div className="space-y-1">
              <span className="text-[10px] uppercase font-bold tracking-wider text-[#707478]">
                Local Cache
              </span>
              <p className="text-lg font-extrabold text-[#041627]">42.8 MB</p>
              <p className="text-xs text-[#44474c] font-medium">14 Vector Blueprints</p>
            </div>
            <div className="w-10 h-10 rounded-lg bg-[#f0f4f8] text-[#041627] flex items-center justify-center">
              <span className="material-symbols-outlined text-[22px]">sd_card</span>
            </div>
          </div>

          <div className="p-4 rounded-xl bg-white border border-[#e0e3e5] shadow-sm flex items-center justify-between">
            <div className="space-y-1">
              <span className="text-[10px] uppercase font-bold tracking-wider text-[#707478]">
                Last Cloud Sync
              </span>
              <p className="text-lg font-extrabold text-[#041627]">{lastCloudSyncText}</p>
              <p className="text-xs text-emerald-700 font-bold flex items-center gap-1">
                <span className="w-2 h-2 rounded-full bg-emerald-500 inline-block" />
                Encrypted Mesh OK
              </p>
            </div>
            <div className="w-10 h-10 rounded-lg bg-[#f0f4f8] text-[#041627] flex items-center justify-center">
              <span className="material-symbols-outlined text-[22px]">cloud_sync</span>
            </div>
          </div>
        </div>

        <p className="text-xs text-[#707478] font-medium text-right pt-1">
          v2.4.1 (Sovereign Mobile Edition)
        </p>
      </div>

      {/* Synchronize All Offline Data Interaction */}
      <div className="pt-2">
        <button
          onClick={handleSyncClick}
          disabled={isSyncing}
          className={`w-full py-3.5 px-4 rounded-xl font-bold text-sm transition-all shadow-sm flex items-center justify-center gap-2 ${
            isSyncing
              ? "bg-[#e0e3e5] text-[#707478] cursor-not-allowed"
              : syncComplete
              ? "bg-emerald-600 text-white"
              : "bg-[#041627] hover:bg-[#006c49] text-white"
          }`}
        >
          <span
            className={`material-symbols-outlined text-[20px] ${
              isSyncing ? "animate-spin" : ""
            }`}
          >
            {syncComplete ? "check_circle" : "sync"}
          </span>
          <span>
            {isSyncing
              ? "Synchronizing Telemetry..."
              : syncComplete
              ? "Synchronized Just Now"
              : "Synchronize All Offline Data"}
          </span>
        </button>
      </div>

      {/* Workspace Security Audit Log Action */}
      <div>
        <button
          onClick={() => navigate(CONSTRUCTION_ROUTES.HISTORY)}
          className="w-full p-4 rounded-xl border border-[#e0e3e5] bg-white hover:border-[#041627] transition-colors shadow-sm flex items-center justify-between group"
        >
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-lg bg-[#f0f4f8] text-[#041627] flex items-center justify-center group-hover:bg-[#041627] group-hover:text-white transition-colors">
              <span className="material-symbols-outlined text-[20px]">shield</span>
            </div>
            <div className="text-left">
              <h3 className="text-sm font-bold text-[#041627]">Workspace Security Audit Log</h3>
              <p className="text-xs text-[#44474c] font-medium">
                Cryptographic notarization &amp; immutable event history
              </p>
            </div>
          </div>
          <span className="material-symbols-outlined text-[#707478] group-hover:text-[#041627] transition-colors text-[20px]">
            chevron_right
          </span>
        </button>
      </div>

      {/* Footer */}
      <footer className="pt-6 border-t border-[#e0e3e5] text-center">
        <p className="text-xs font-semibold text-[#707478]">
          BuildFlow Field OS • Verified Encrypted Node #881-A
        </p>
      </footer>
    </div>
  );
};

export default ProfilePage;
