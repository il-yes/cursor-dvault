import { C3styles } from "@/components/C3/styles/styles";
import React, { useEffect, useState } from "react";
import { fetchDepartments } from "@/components/C3/domain/channel/channel.repository";
import { ChannelSlot, Department, VaultAssignment } from "@/components/C3/domain/channel/channel.types";
import { CreateChannelDraft } from "../types";

interface Step3Props {
    data: CreateChannelDraft;
    onBack: () => void;
    onNext: (
        values: Partial<CreateChannelDraft>
    ) => void;
}

interface VaultRowItemProps {
    slot: ChannelSlot;
    departments: Department[];
    updateAssignedUser: (vaultName: string, userId: string) => void;
    updateAssignedPublicKey: (vaultName: string, publicKey: string) => void;
    setAssignmentColor: (vaultName: string, color: string) => void;
}

const VaultRowItem = ({
    slot,
    departments,
    updateAssignedUser,
    updateAssignedPublicKey,
    setAssignmentColor,
}: VaultRowItemProps) => {
    const [user, setUser] = useState<string | undefined>(undefined);
    const [publicKey, setPublicKey] = useState<string | undefined>(undefined);
    const department = departments.find((v: Department) => v.id === slot?.vault);
    const color = department?.color;

    const updateUser = (update: string) => {
        setUser(update);
        updateAssignedUser(slot.vault, update);
    };

    const updatePub = (update: string) => {
        setPublicKey(update);
        updateAssignedPublicKey(slot.vault, update);
    };

    useEffect(() => {
        if (department && color) {
            setAssignmentColor(slot.vault, color);
        }
    }, [department, color, slot.vault, setAssignmentColor]);

    return (
        <div className="vault-select-wrap">
            <div className="vs-dot" style={{ background: color }} />
            <div className="vs-name">{department?.name || slot.vault}</div>
            <div className="vs-arrow">▾</div>
            <div className="vs-inputs flex w-full flex-col gap-y-3 p-2">
                <input type="text" value={user || ""} onChange={(v) => updateUser(v.target.value)} placeholder="User ID / Name" />
                <input type="text" value={publicKey || ""} onChange={(v) => updatePub(v.target.value)} placeholder="Public Key" />
            </div>
        </div>
    );
};

export const Step3 = ({ data, onBack, onNext }: Step3Props) => {
    console.log({ data });
    const [departments, setDepartments] = useState<Department[]>([]);
    const [assignments, setAssignments] = useState<VaultAssignment[]>(() => {
        if (data.assignments && data.assignments.length > 0) {
            return data.assignments;
        }
        const slots = data?.slots ?? data?.template?.slots ?? [];
        return slots.map((slot) => ({
            vault: slot.vault,
            owner: {
                id: "",
                label: "",
                publicKey: "",
            },
        }));
    });

    useEffect(() => {
        getDepartments();
    }, []);

    const getDepartments = async () => {
        try {
            const res = await fetchDepartments();
            setDepartments(res);
        } catch (err) {
            console.error("Failed to fetch departments", err);
        }
    };

    const updateAssignedUser = (vaultName: string, userId: string) => {
        setAssignments((prev) =>
            prev.map((assignment) =>
                assignment.vault === vaultName
                    ? { ...assignment, owner: { ...assignment.owner, id: userId, label: userId } }
                    : assignment
            )
        );
    };

    const updateAssignedPublicKey = (vaultName: string, publicKey: string) => {
        setAssignments((prev) =>
            prev.map((assignment) =>
                assignment.vault === vaultName
                    ? { ...assignment, owner: { ...assignment.owner, publicKey } }
                    : assignment
            )
        );
    };

    const setAssignmentColor = (vaultName: string, color: string) => {
        setAssignments((prev) =>
            prev.map((assignment) =>
                assignment.vault === vaultName
                    ? { ...assignment, vaultColor: color }
                    : assignment
            )
        );
    };

    const slotsToRender = data?.slots?.length ? data.slots : data?.template?.slots ?? [];

    return (
        <div className="modal">
            <C3styles />
            <div className="modal-header">
                <div>
                    <div className="modal-title">Add Vaults</div>
                    <div className="modal-subtitle">{data.channelName}</div>
                </div>
                <div className="step-indicator">
                    <div className="step-label">Step 3 of 4</div>
                    <div className="step-dots">
                        <div className="sdot-i done" />
                        <div className="sdot-i done" />
                        <div className="sdot-i active" />
                        <div className="sdot-i" />
                    </div>
                </div>  
            </div>
            <div className="modal-body">
                <div className="section-label">Internal Vaults</div>
                <div className="section-hint">
                    Roles and access levels are set by the template. Reassign any vault
                    from the dropdown.
                </div>
                <table className="role-table">
                    <thead>
                        <tr>
                            <th style={{ width: "22%" }}>Role</th>
                            <th style={{ width: "52%" }}>Vault assigned</th>
                            <th style={{ width: "26%" }}>Access</th>
                        </tr>
                    </thead>
                    <tbody>
                        {slotsToRender.map((slot, i) => (
                            <tr key={slot.id || i}>
                                <td>
                                    <div className="role-name">{slot.role}</div>
                                </td>
                                <td>
                                    <VaultRowItem
                                        slot={slot}
                                        departments={departments}
                                        updateAssignedUser={updateAssignedUser}
                                        updateAssignedPublicKey={updateAssignedPublicKey}
                                        setAssignmentColor={setAssignmentColor}
                                    />
                                </td>
                                <td>
                                    <span className="role-access">Write</span>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                <hr className="section-divider" />
                {/* C3 Extension section */}
                <div className="c3-optional-box">
                    <div className="c3-box-header">
                        <div className="c3-box-title">
                            <span className="c3-chain-icon">⛓</span>
                            C3 Extension
                        </div>
                        <span className="c3-optional-badge">Optional</span>
                    </div>
                    <div className="c3-box-desc">
                        Invite an external party — a supplier, auditor, or counterparty —
                        to join this channel with a specific role. Can also be added after
                        activation.
                    </div>
                    <div className="c3-add-btn">+ Add external vault</div>
                    <div className="invite-wrap">
                        <input
                            className="invite-input"
                            type="text"
                            placeholder="Paste vault address (ankhora://vault_id)"
                        />
                        <span className="invite-or">or</span>
                        <button className="invite-link-btn">Send invite link →</button>
                    </div>
                    <div className="after-activation-note">
                        No external vault added? The channel activates as internal — ⛓+
                        will appear on the ledger row when ready.
                    </div>
                </div>
            </div>
            <div className="modal-footer">
                <button className="btn " onClick={() =>
                    onBack()
                }>← Back</button>
                <button className="btn btn-primary" onClick={() =>
                    onNext({
                        ...data,
                        assignments
                    })
                }>Next →</button>
            </div>
        </div>
    );
};
