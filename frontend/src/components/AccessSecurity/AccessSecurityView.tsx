import { useState, useEffect, useCallback } from 'react';
import { Shield, Database, Activity, Users, Section } from "lucide-react";

function getActivityIcon(type) {
  const icons = {
    'upload': '📤',
    'download': '📥',
    'share': '🔗',
    'delete': '🗑️',
    'login': '🔐',
    'settings': '⚙️'
  };
  return icons[type] || '📋';
}

function getScoreRecommendation(stats) {
  if (!stats.mfaEnabled) {
    return 'Enable MFA to improve your security score';
  }
  if (!stats.recoveryBackedUp) {
    return 'Backup your recovery phrase to improve your security score';
  }
  if (stats.securityScore >= 90) {
    return 'Your vault is highly secure. Keep up the good practices!';
  }
  return 'Review recommendations above to improve your security';
}

async function verifyActivity(activity) {
  // Open verification modal for this activity
  const commit = await fetchCommitFromStellar(activity.commitHash);
  // Show modal with commit details (reuse EncryptionVerificationModal)
}

function fetchCommitFromStellar(commitHash) {
  return commitHash;
}
const AccessSecurityView = () => {
  // Mock data with Wails safety checks
  const [stats, setStats] = useState({
    mfaEnabled: false,
    lastLogin: {
      time: '2 hours',
      location: 'San Francisco, CA'
    },
    recoveryBackedUp: false
  });
  const [isWailsReady, setIsWailsReady] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // Safe Wails backend check
  const fetchSecurityStats = useCallback(async () => {
    if (typeof window !== 'undefined' && window.go && window.go?.Security) {
      try {
        setIsWailsReady(true);
        const securityStats = await window.go?.Security.GetSecurityStats();
        setStats(securityStats || stats);
      } catch (error) {
        console.warn('Security stats fetch failed:', error);
        // Keep mock data
      }
    } else {
      console.log('Wails not ready - using mock security data');
      // Enhanced mock data
      setStats({
        mfaEnabled: Math.random() > 0.7, // 30% chance disabled for demo
        lastLogin: {
          time: ['1 min', '2 hours', '1 day', '3 days'][Math.floor(Math.random() * 4)],
          location: ['San Francisco, CA', 'New York, NY', 'London, UK', 'Berlin, DE'][Math.floor(Math.random() * 4)]
        },
        recoveryBackedUp: Math.random() > 0.3 // 70% chance backed up
      });
    }
    setIsLoading(false);
  }, []);

//   useEffect(() => {
//     fetchSecurityStats();
    
//     // Poll every 30s for real-time security updates
//     const interval = setInterval(fetchSecurityStats, 30000);
//     return () => clearInterval(interval);
//   }, [fetchSecurityStats]);

  if (isLoading) {
    return (
      <Section className="security-section">
        <div className="section-header">
          <span className="icon">🔐</span>
          <h2>Access Security</h2>
        </div>
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-yellow-400"></div>
          <span className="ml-2 text-sm text-muted-foreground">Loading security status...</span>
        </div>
      </Section>
    );
  }

  return (
    <Section className="security-section">
      <div className="section-header">
        <span className="icon">🔐</span>
        <h2>Access Security</h2>
        <div className={`ml-auto px-3 py-1 rounded-full text-xs font-medium ${
          isWailsReady 
            ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30' 
            : 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/30'
        }`}>
          {isWailsReady ? '✅ Live Data' : '🔧 Mock Mode'}
        </div>
      </div>
      
      <ul className="status-list">
        {/* MFA Status */}
        <li className={`status-item ${stats.mfaEnabled ? 'success' : 'warning'}`}>
          <span className="status-icon">{stats.mfaEnabled ? '✅' : '⚠️'}</span>
          <span className="status-text">
            MFA (Multi-Factor Authentication) {stats.mfaEnabled ? 'enabled' : 'disabled'}
          </span>
          {!stats.mfaEnabled && (
            <button 
              className="action-btn"
              onClick={() => window.location.href = '/dashboard/settings'}
            >
              Enable MFA
            </button>
          )}
        </li>

        {/* Last Login */}
        <li className="status-item success">
          <span className="status-icon">✅</span>
          <span className="status-text">
            Last login: <span className="font-semibold">{stats.lastLogin.time}</span> ago from{' '}
            <span className="font-semibold">{stats.lastLogin.location}</span>
          </span>
        </li>

        {/* Recovery Phrase */}
        <li className={`status-item ${stats.recoveryBackedUp ? 'success' : 'warning'}`}>
          <span className="status-icon">{stats.recoveryBackedUp ? '✅' : '⚠️'}</span>
          <span className="status-text">
            Recovery phrase {stats.recoveryBackedUp ? 'backed up' : 'not backed up'}
          </span>
          {!stats.recoveryBackedUp && (
            <button 
              className="action-btn warning"
              onClick={() => window.location.href = '/dashboard/settings'}
            >
              Backup Now
            </button>
          )}
        </li>
      </ul>
    </Section>
  );
};

export default AccessSecurityView;
