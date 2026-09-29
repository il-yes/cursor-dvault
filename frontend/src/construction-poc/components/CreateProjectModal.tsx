import React, { useState } from "react";
import { createProject } from "../data";

interface CreateProjectModalProps {
  isOpen: boolean;
  onClose: () => void;
  onProjectCreated?: () => void;
}

export const CreateProjectModal: React.FC<CreateProjectModalProps> = ({
  isOpen,
  onClose,
  onProjectCreated,
}) => {
  const [projectCode, setProjectCode] = useState("");
  const [projectName, setProjectName] = useState("");
  const [projectType, setProjectType] = useState("INFRASTRUCTURE");
  const [location, setLocation] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submittedNotice, setSubmittedNotice] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!projectCode.trim() || !projectName.trim()) return;

    setIsSubmitting(true);
    setSubmittedNotice(null);
    try {
      const created = await createProject({
        code: projectCode.trim(),
        name: projectName.trim(),
        type: projectType.trim(),
        location: location.trim(),
      });

      setSubmittedNotice(
        `Project "${created.name}" (${created.code}) created successfully in Cloud under workspace context.`
      );
      if (onProjectCreated) {
        onProjectCreated();
      }
    } catch (err: any) {
      console.error("Create project failed:", err);
      setSubmittedNotice(`Failed to create project in Cloud: ${err?.message || err}`);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div className="bg-white rounded-xl max-w-lg w-full p-6 shadow-2xl border border-[#e0e3e5]">
        <div className="flex items-center justify-between pb-3 border-b border-[#f0f3f5] mb-4">
          <div className="flex items-center gap-2">
            <span className="material-symbols-outlined text-[#041627]">add_business</span>
            <h2 className="text-lg font-bold text-[#181c1e]">Create Construction Project</h2>
          </div>
          <button
            onClick={onClose}
            className="text-[#44474c] hover:text-[#181c1e] text-lg font-bold"
          >
            ✕
          </button>
        </div>

        {submittedNotice ? (
          <div className="space-y-4">
            <div className="p-4 bg-[#fff3e0] border border-[#ffe0b2] rounded-lg text-xs text-[#181c1e] space-y-2 whitespace-pre-wrap">
              <span className="font-bold text-[#b76e00] block text-sm flex items-center gap-1">
                <span className="material-symbols-outlined text-[18px]">info</span>
                Application Boundary Notice
              </span>
              <p>{submittedNotice}</p>
            </div>
            <div className="flex justify-end">
              <button
                onClick={() => {
                  setSubmittedNotice(null);
                  onClose();
                }}
                className="px-4 py-2 bg-[#041627] text-white text-xs font-bold rounded-lg hover:bg-[#041627]/90"
              >
                Close Modal
              </button>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4 text-xs">
            <div>
              <label className="block text-[#44474c] font-bold mb-1">Project Identifier Code</label>
              <input
                type="text"
                placeholder="e.g. PRJ-002"
                value={projectCode}
                onChange={(e) => setProjectCode(e.target.value)}
                required
                className="w-full px-3 py-2 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg text-sm text-[#181c1e] outline-none focus:border-[#041627]"
              />
            </div>

            <div>
              <label className="block text-[#44474c] font-bold mb-1">Project Name</label>
              <input
                type="text"
                placeholder="e.g. Harbor Viaduct Phase 2"
                value={projectName}
                onChange={(e) => setProjectName(e.target.value)}
                required
                className="w-full px-3 py-2 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg text-sm text-[#181c1e] outline-none focus:border-[#041627]"
              />
            </div>

            <div>
              <label className="block text-[#44474c] font-bold mb-1">Project Category</label>
              <select
                value={projectType}
                onChange={(e) => setProjectType(e.target.value)}
                className="w-full px-3 py-2 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg text-sm text-[#181c1e] outline-none focus:border-[#041627]"
              >
                <option value="INFRASTRUCTURE">INFRASTRUCTURE</option>
                <option value="COMMERCIAL">COMMERCIAL</option>
                <option value="RESIDENTIAL">RESIDENTIAL</option>
                <option value="INDUSTRIAL">INDUSTRIAL</option>
              </select>
            </div>

            <div>
              <label className="block text-[#44474c] font-bold mb-1">Site Location</label>
              <input
                type="text"
                placeholder="e.g. North Bay Terminal, Chicago, IL"
                value={location}
                onChange={(e) => setLocation(e.target.value)}
                required
                className="w-full px-3 py-2 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg text-sm text-[#181c1e] outline-none focus:border-[#041627]"
              />
            </div>

            <div className="pt-3 border-t border-[#f0f3f5] flex items-center justify-between">
              <span className="text-[10px] text-[#44474c]">
                UX presentation boundary for project creation
              </span>
              <div className="flex gap-2">
                <button
                  type="button"
                  onClick={onClose}
                  className="px-3 py-2 bg-[#f0f3f5] text-[#041627] font-bold rounded-lg hover:bg-[#e0e3e5]"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="px-4 py-2 bg-[#041627] text-white font-bold rounded-lg hover:bg-[#041627]/90 disabled:opacity-50"
                >
                  {isSubmitting ? "Creating..." : "Create Project"}
                </button>
              </div>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};

export default CreateProjectModal;
