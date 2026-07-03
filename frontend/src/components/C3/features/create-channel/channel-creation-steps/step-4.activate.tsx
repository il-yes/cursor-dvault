import { C3styles } from '@/components/C3/styles/styles';
import React, { useEffect, useState } from 'react';
import { Department, VaultAssignment } from '@/components/C3/domain/channel/channel.types';
import { fetchDepartments } from '@/components/C3/domain/channel/channel.repository';
import { CreateChannelDraft } from '../types';


interface Step4Props {

  data: CreateChannelDraft;
  onBack: () => void;
  onCreate: (
    values: Partial<CreateChannelDraft>
  ) => void;

}

export const C3Step4 = ({ data, onBack, onCreate }: Step4Props) => {
  const getAssignment = (vault: string) => {
    return data.assignments?.find(
      a => a.vault === vault
    );
  };


  return (
    <div className="modal">
      <C3styles />
      <div className="modal-header">
        <div>
          <div className="modal-title">Activate</div>
          <div className="modal-subtitle">{data.template?.id}</div>
        </div>
        <div className="step-indicator">
          <div className="step-label">Step 4 of 4</div>
          <div className="step-dots">
            <div className="sdot-i done" />
            <div className="sdot-i done" />
            <div className="sdot-i done" />
            <div className="sdot-i active" />
          </div>
        </div>
      </div>
      <div className="modal-body">
        {/* Summary block */}
        <div className="summary-block">
          <div className="summary-row">
            <div className="sum-key">Channel</div>
            <div className="sum-val">{data.channelName}</div>
          </div>
          <div className="summary-row">
            <div className="sum-key">Template</div>
            <div className="sum-val">{data.template?.id}</div>
          </div>
          <div className="summary-row">
            <div className="sum-key">Custom property</div>
            <div className="sum-val">
              {data.properties?.map((p, i) => (
                <div key={i} className="custom-kv">
                  <span className="kv-key">{p.key}</span>
                  <span className="kv-sep">:</span>
                  <span className="kv-val">{p.value}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
        {/* Flow */}
        <div className="flow-section">
          <div className="section-label">Flow</div>
          <div className="flow-visual">
            {data?.slots?.map((slot, i) => {
              const assignment = getAssignment(slot.vault);

              return (
                <React.Fragment key={slot.id}>
                  <div className="fv-vault">
                    <div className="fv-dot" style={{ background: assignment?.vaultColor }} />
                    {slot.vault}
                  </div>
                  {i !== data?.slots?.length - 1 && <div className="fv-arrow">→</div>}
                </React.Fragment>
              )
            })}
          </div>
        </div>
        {/* Slots */}
        <div className="section-label">Slots</div>
        <table className="slots-summary">
          <thead>
            <tr>
              <th style={{ width: "38%" }}>Slot</th>
              <th style={{ width: "32%" }}>Vault</th>
              <th style={{ width: "15%" }}>Access</th>
              <th style={{ width: "15%" }} />
            </tr>
          </thead>
          <tbody>
            {data.slots?.map(slot => {
              const assignment = getAssignment(slot.vault);

              return (
                <tr key={slot.id}>
                  <td>
                    <div className="vault-chip">
                      <div className="vc-dot" style={{ background: assignment?.vaultColor }} />
                      {slot.name}
                    </div>
                  </td>
                  <td style={{ color: "#aaa", fontSize: 11 }}>
                    {slot.vault} <span className="text-gray-400">-</span>  {assignment?.owner.label}
                  </td>
                  <td style={{ color: "#aaa", fontSize: 11 }}>Write</td>
                  <td>

                    {slot.gated && (
                      <span className="gate-on">
                        gate ●
                      </span>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
        {/* Status rows */}
        <div className="status-section">
          <div className="status-row">
            <div className="sr-left">
              <span className="sr-icon">⛓</span>
              <span className="sr-label">C3 Extension</span>
            </div>
            <div className="sr-val">
              Not set &nbsp;<span className="c3-ext-badge">⛓+</span>
            </div>
          </div>
          <div className="status-row">
            <div className="sr-left">
              <span className="sr-icon">✦</span>
              <span className="sr-label">Stellar anchoring</span>
            </div>
            <div className="sr-val positive">
              <div className="stellar-dot" />
              ON — every commit anchored
            </div>
          </div>
        </div>
        {/* Post-activation ledger preview */}
        <div className="post-preview">
          <div className="pp-label">After activation — ledger row</div>
          <div className="pp-row">
            <div className="pp-sdot" />
            <div className="pp-name">contract-execution — Cipla_India</div>
            <div className="pp-meta">
              <span>just now</span>
              <div className="pp-pipe">
                <div className="pp-seg" />
                <div className="pp-seg" />
                <div className="pp-seg" />
              </div>
              <span className="pp-c3">⛓+</span>
            </div>
          </div>
        </div>
      </div>
      <div className="modal-footer">
        <button className="btn " onClick={() => onBack()}>← Back</button>
        <button className="btn-activate" onClick={() => onCreate(data)}>⚡ Start Channel</button>
      </div>
    </div>
  );
}



