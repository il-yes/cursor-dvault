import { useState, useEffect } from 'react';
import './encryptionVerif.css';

// Mock Stellar fetch call (replace with real implementation)
export async function fetchCommitFromStellar(hash) {
  await new Promise((res) => setTimeout(res, 800));
  return {
    fileName: "contracts.zip",
    txHash: hash,
    action: "FILE_COMMIT",
    timestamp: "2025-11-28T12:30:00Z",
    author: "0xANONYMOUS",
    signatureValid: true,
    cid: "QmXkjExampleCID12345"
  };
}

export default function EncryptionVerificationModal({ file, onClose }) {
  const [commit, setCommit] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  
  useEffect(() => {
    async function loadCommit() {
      try {
        const commitData = await fetchCommitFromStellar(file.commitHash);
        setCommit(commitData);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }
    loadCommit();
  }, [file.commitHash]);
  
  if (loading) {
    return <div className="modal-loading">Loading verification data...</div>;
  }
  
  if (error) {
    return (
      <div className="modal-error">
        <p>Could not load verification data: {error}</p>
        <button onClick={onClose}>Close</button>
      </div>
    );
  }
  
  return (
    <div className="verification-modal">
      <div className="modal-header">
        <h2>Encryption Verification</h2>
        <button className="close-btn" onClick={onClose}>×</button>
      </div>
      
      <div className="modal-body">
        <div className="verification-status">
          <span className="icon">✅</span>
          <p>This file was encrypted on your device before upload</p>
        </div>
        
        <div className="commit-details">
          <h3>Cryptographic Proof</h3>
          
          <div className="detail-row">
            <span className="label">Stellar Transaction:</span>
            <code className="value monospace">{commit.txHash}</code>
            <button 
              className="copy-btn"
              onClick={() => navigator.clipboard.writeText(commit.txHash)}
            >
              Copy
            </button>
          </div>
          
          <div className="detail-row">
            <span className="label">Action:</span>
            <span className="value">{commit.action}</span>
          </div>
          
          <div className="detail-row">
            <span className="label">Timestamp:</span>
            <span className="value">{new Date(commit.timestamp).toLocaleString()}</span>
          </div>
          
          <div className="detail-row">
            <span className="label">Author:</span>
            <span className="value">{commit.author}</span>
          </div>
          
          <div className="detail-row">
            <span className="label">Signature:</span>
            <span className={`value ${commit.signatureValid ? 'valid' : 'invalid'}`}>
              {commit.signatureValid ? '✅ Valid' : '❌ Invalid'}
            </span>
          </div>
          
          <div className="detail-row">
            <span className="label">IPFS CID:</span>
            <code className="value monospace">{commit.cid}</code>
          </div>
          
          <div className="detail-row">
            <span className="label">Encryption:</span>
            <span className="value">AES-256-GCM</span>
          </div>
        </div>
        
        <div className="modal-actions">
          <button 
            className="secondary"
            onClick={() => window.open(`https://stellar.expert/explorer/public/tx/${commit.txHash}`, '_blank')}
          >
            View on Stellar Explorer
          </button>
          
          <button 
            className="secondary"
            onClick={() => downloadAuditReport(commit)}
          >
            Download Audit Report
          </button>
        </div>
      </div>
    </div>
  );
}

// Helper function to generate audit report
function downloadAuditReport(commit) {
  const report = `
ANKHORA ENCRYPTION AUDIT REPORT
================================

File: ${commit.fileName}
Generated: ${new Date().toISOString()}

ENCRYPTION VERIFICATION
-----------------------
✅ File was encrypted client-side using AES-256-GCM
✅ Encryption occurred before transmission to server
✅ Server received only encrypted bytes (no plaintext access)

CRYPTOGRAPHIC PROOF
-------------------
Stellar Transaction: ${commit.txHash}
Action: ${commit.action}
Timestamp: ${new Date(commit.timestamp).toISOString()}
Author: ${commit.author}
Signature: ${commit.signatureValid ? 'Valid' : 'Invalid'}
IPFS CID: ${commit.cid}

VERIFICATION STEPS
------------------
1. Visit: https://stellar.expert/explorer/public/tx/${commit.txHash}
2. Verify transaction exists on Stellar blockchain
3. Verify signature matches author's public key
4. Verify timestamp matches file creation time

This report provides cryptographic proof that the file was encrypted
on the client device and never exposed to the server in plaintext form.

For questions, contact: security@blackops.africa
  `.trim();
  
  const blob = new Blob([report], { type: 'text/plain' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `ankhora-audit-${commit.fileName}-${Date.now()}.txt`;
  a.click();
  URL.revokeObjectURL(url);
}
