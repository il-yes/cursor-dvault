import React from "react";

export const DocumentsPage: React.FC = () => {
  const documents = [
    {
      id: "DOC-2024-08",
      name: "Structural Blueprints v4.2",
      type: "Architectural Drawing",
      updated: "Sep 15, 2026",
      size: "24.5 MB",
    },
    {
      id: "DOC-2024-05",
      name: "OSHA Safety Compliance Certificate",
      type: "Permit & License",
      updated: "Aug 28, 2026",
      size: "1.8 MB",
    },
    {
      id: "DOC-2024-02",
      name: "Concrete Mixture Specification C3",
      type: "Material Standard",
      updated: "Aug 12, 2026",
      size: "4.2 MB",
    },
  ];

  return (
    <div className="flex flex-col w-full gap-4 px-4 py-4 max-w-2xl mx-auto pb-8">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-[#181c1e]">Project Documents</h1>
        <button className="bg-[#041627] text-white px-3 py-1.5 rounded-lg text-xs font-semibold hover:bg-[#041627]/90 transition-colors flex items-center gap-1">
          <span className="material-symbols-outlined text-[16px]">upload_file</span>
          Upload
        </button>
      </div>

      <div className="flex flex-col gap-3">
        {documents.map((doc) => (
          <div
            key={doc.id}
            className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex items-center justify-between group hover:border-[#041627]/30 transition-colors"
          >
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-[#ebeef0] text-[#041627] flex items-center justify-center">
                <span className="material-symbols-outlined text-[20px]">
                  description
                </span>
              </div>
              <div className="flex flex-col">
                <h3 className="text-sm font-semibold text-[#181c1e]">
                  {doc.name}
                </h3>
                <span className="text-xs text-[#44474c]">
                  {doc.type} • {doc.size}
                </span>
              </div>
            </div>
            <button className="text-[#44474c] hover:text-[#041627] p-2">
              <span className="material-symbols-outlined text-[20px]">
                download
              </span>
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default DocumentsPage;
