# Ankhora Bridge Extension - Phase 1 Technical Specification

## 1. Native Messaging Protocol

### 1.1 Overview

Native messaging enables bidirectional communication between the browser extension and the desktop app. Messages are exchanged via **stdin/stdout** using JSON format.

### 1.2 Protocol Specification

**Message Format:**
```json
{
  "id": "uuid-v4",
  "action": "string",
  "payload": {},
  "timestamp": "ISO-8601"
}
```

**Response Format:**
```json
{
  "id": "uuid-v4",
  "success": true/false,
  "data": {},
  "error": "string (optional)",
  "timestamp": "ISO-8601"
}
```

### 1.3 Message Types

#### Extension → Desktop App

**1. detect_login**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "action": "detect_login",
  "payload": {
    "url": "https://github.com/login",
    "domain": "github.com",
    "title": "Sign in to GitHub"
  },
  "timestamp": "2025-12-06T12:00:00Z"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "success": true,
  "data": {
    "credentials": [
      {
        "id": "cred-123",
        "name": "GitHub Work",
        "username": "manuel@blackops.io",
        "vault": "work",
        "verified": true,
        "last_used": "2025-12-05T10:30:00Z"
      }
    ]
  },
  "timestamp": "2025-12-06T12:00:00.150Z"
}
```

**2. request_autofill**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "action": "request_autofill",
  "payload": {
    "credential_id": "cred-123",
    "url": "https://github.com/login"
  },
  "timestamp": "2025-12-06T12:00:01Z"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "success": true,
  "data": {
    "username": "manuel@blackops.io",
    "password": "encrypted_password_here",
    "method": "inject"
  },
  "timestamp": "2025-12-06T12:00:01.080Z"
}
```

**3. save_credential**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "action": "save_credential",
  "payload": {
    "url": "https://github.com/login",
    "domain": "github.com",
    "username": "manuel@blackops.io",
    "password": "user_password",
    "name": "GitHub Work"
  },
  "timestamp": "2025-12-06T12:00:02Z"
}
```

**4. search_vault**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440003",
  "action": "search_vault",
  "payload": {
    "query": "github",
    "context_url": "https://github.com/login"
  },
  "timestamp": "2025-12-06T12:00:03Z"
}
```

#### Desktop App → Extension

**1. autofill_command**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440004",
  "action": "autofill_command",
  "payload": {
    "username": "manuel@blackops.io",
    "password": "decrypted_password",
    "tab_id": 12345,
    "method": "inject"
  },
  "timestamp": "2025-12-06T12:00:04Z"
}
```

**2. show_notification**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440005",
  "action": "show_notification",
  "payload": {
    "message": "Credential saved successfully",
    "type": "success",
    "duration": 3000
  },
  "timestamp": "2025-12-06T12:00:05Z"
}
```

### 1.4 Native Messaging Host Setup

**macOS Configuration:**

**File:** `~/Library/Application Support/Google/Chrome/NativeMessagingHosts/com.ankhora.vault.json`

```json
{
  "name": "com.ankhora.vault",
  "description": "Ankhora Vault Native Messaging Host",
  "path": "/Applications/Ankhora.app/Contents/MacOS/ankhora-native-host",
  "type": "stdio",
  "allowed_origins": [
    "chrome-extension://abcdefghijklmnopqrstuvwxyz123456/"
  ]
}
```

**Windows Configuration:**

**Registry Key:** `HKEY_CURRENT_USER\Software\Google\Chrome\NativeMessagingHosts\com.ankhora.vault`

**Value:** Path to JSON manifest file

**File:** `C:\Program Files\Ankhora\native-host-manifest.json`

```json
{
  "name": "com.ankhora.vault",
  "description": "Ankhora Vault Native Messaging Host",
  "path": "C:\\Program Files\\Ankhora\\ankhora-native-host.exe",
  "type": "stdio",
  "allowed_origins": [
    "chrome-extension://abcdefghijklmnopqrstuvwxyz123456/"
  ]
}
```

**Linux Configuration:**

**File:** `~/.config/google-chrome/NativeMessagingHosts/com.ankhora.vault.json`

```json
{
  "name": "com.ankhora.vault",
  "description": "Ankhora Vault Native Messaging Host",
  "path": "/usr/local/bin/ankhora-native-host",
  "type": "stdio",
  "allowed_origins": [
    "chrome-extension://abcdefghijklmnopqrstuvwxyz123456/"
  ]
}
```

### 1.5 Native Host Implementation (Python)

```python
#!/usr/bin/env python3
import sys
import json
import struct
import logging
from vault_manager import VaultManager

class NativeMessagingHost:
    def __init__(self):
        self.vault = VaultManager()
        logging.basicConfig(filename='/tmp/ankhora-native.log', level=logging.DEBUG)
    
    def read_message(self):
        """Read message from stdin (sent by extension)"""
        raw_length = sys.stdin.buffer.read(4)
        if not raw_length:
            return None
        
        message_length = struct.unpack('=I', raw_length)[0]
        message = sys.stdin.buffer.read(message_length).decode('utf-8')
        return json.loads(message)
    
    def send_message(self, message):
        """Send message to stdout (received by extension)"""
        encoded_message = json.dumps(message).encode('utf-8')
        message_length = struct.pack('=I', len(encoded_message))
        sys.stdout.buffer.write(message_length)
        sys.stdout.buffer.write(encoded_message)
        sys.stdout.buffer.flush()
    
    def handle_message(self, message):
        """Route message to appropriate handler"""
        action = message.get('action')
        payload = message.get('payload', {})
        msg_id = message.get('id')
        
        try:
            if action == 'detect_login':
                credentials = self.vault.search_by_domain(payload['domain'])
                return {
                    'id': msg_id,
                    'success': True,
                    'data': {'credentials': [c.to_dict() for c in credentials]}
                }
            
            elif action == 'request_autofill':
                credential = self.vault.get_credential(payload['credential_id'])
                decrypted = self.vault.decrypt_credential(credential)
                return {
                    'id': msg_id,
                    'success': True,
                    'data': {
                        'username': decrypted.username,
                        'password': decrypted.password,
                        'method': 'inject'
                    }
                }
            
            elif action == 'save_credential':
                credential_id = self.vault.save_credential(
                    url=payload['url'],
                    username=payload['username'],
                    password=payload['password'],
                    name=payload.get('name')
                )
                return {
                    'id': msg_id,
                    'success': True,
                    'data': {'credential_id': credential_id}
                }
            
            elif action == 'search_vault':
                results = self.vault.search(payload['query'])
                return {
                    'id': msg_id,
                    'success': True,
                    'data': {'results': [r.to_dict() for r in results]}
                }
            
            else:
                return {
                    'id': msg_id,
                    'success': False,
                    'error': f'Unknown action: {action}'
                }
        
        except Exception as e:
            logging.error(f"Error handling {action}: {str(e)}")
            return {
                'id': msg_id,
                'success': False,
                'error': str(e)
            }
    
    def run(self):
        """Main event loop"""
        logging.info("Native messaging host started")
        
        while True:
            message = self.read_message()
            if message is None:
                break
            
            logging.debug(f"Received: {message}")
            response = self.handle_message(message)
            logging.debug(f"Sending: {response}")
            self.send_message(response)

if __name__ == '__main__':
    host = NativeMessagingHost()
    host.run()
```

---

## 2. Extension Architecture (Manifest V3)

### 2.1 Manifest Configuration

**File:** `manifest.json`

```json
{
  "manifest_version": 3,
  "name": "Ankhora Vault",
  "version": "1.0.0",
  "description": "Zero-knowledge password manager with cryptographic verification",
  
  "permissions": [
    "nativeMessaging",
    "activeTab",
    "scripting",
    "storage",
    "notifications"
  ],
  
  "host_permissions": [
    "<all_urls>"
  ],
  
  "background": {
    "service_worker": "background.js",
    "type": "module"
  },
  
  "content_scripts": [
    {
      "matches": ["<all_urls>"],
      "js": ["content.js"],
      "run_at": "document_idle",
      "all_frames": false
    }
  ],
  
  "action": {
    "default_popup": "popup.html",
    "default_icon": {
      "16": "icons/icon16.png",
      "32": "icons/icon32.png",
      "48": "icons/icon48.png",
      "128": "icons/icon128.png"
    }
  },
  
  "icons": {
    "16": "icons/icon16.png",
    "32": "icons/icon32.png",
    "48": "icons/icon48.png",
    "128": "icons/icon128.png"
  },
  
  "web_accessible_resources": [
    {
      "resources": ["icons/*", "overlay.html"],
      "matches": ["<all_urls>"]
    }
  ]
}
```

### 2.2 File Structure

```
ankhora-extension/
├── manifest.json
├── background.js          # Service worker (message routing)
├── content.js             # Injected into web pages (form detection, autofill)
├── popup.html             # Extension popup UI
├── popup.js               # Popup logic
├── overlay.html           # Ankhora icon overlay
├── utils/
│   ├── messaging.js       # Native messaging wrapper
│   ├── form-detector.js   # Login form detection
│   └── autofill.js        # Autofill injection logic
├── icons/
│   ├── icon16.png
│   ├── icon32.png
│   ├── icon48.png
│   └── icon128.png
└── styles/
    └── content.css        # Styles for injected elements
```

### 2.3 Background Service Worker

**File:** `background.js`

```javascript
// Native messaging port
let nativePort = null;

// Connect to native host on startup
function connectNative() {
  nativePort = chrome.runtime.connectNative('com.ankhora.vault');
  
  nativePort.onMessage.addListener((message) => {
    console.log('Received from native:', message);
    handleNativeMessage(message);
  });
  
  nativePort.onDisconnect.addListener(() => {
    console.error('Native host disconnected:', chrome.runtime.lastError);
    // Attempt reconnect after 5 seconds
    setTimeout(connectNative, 5000);
  });
}

// Handle messages from native host
function handleNativeMessage(message) {
  const { action, payload } = message;
  
  if (action === 'autofill_command') {
    // Send autofill to active tab
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      if (tabs[0]) {
        chrome.tabs.sendMessage(tabs[0].id, {
          action: 'autofill',
          username: payload.username,
          password: payload.password
        });
      }
    });
  } else if (action === 'show_notification') {
    chrome.notifications.create({
      type: 'basic',
      iconUrl: 'icons/icon48.png',
      title: 'Ankhora',
      message: payload.message
    });
  }
}

// Handle messages from content script
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  console.log('Received from content:', message);
  
  if (message.action === 'detect_login') {
    // Forward to native host
    sendToNative({
      action: 'detect_login',
      payload: {
        url: sender.tab.url,
        domain: new URL(sender.tab.url).hostname,
        title: sender.tab.title
      }
    }).then(response => {
      sendResponse(response);
    });
    return true; // Keep channel open for async response
  }
  
  else if (message.action === 'request_autofill') {
    sendToNative({
      action: 'request_autofill',
      payload: message.payload
    }).then(response => {
      sendResponse(response);
    });
    return true;
  }
  
  else if (message.action === 'save_credential') {
    sendToNative({
      action: 'save_credential',
      payload: message.payload
    }).then(response => {
      sendResponse(response);
    });
    return true;
  }
});

// Send message to native host (with Promise wrapper)
function sendToNative(message) {
  return new Promise((resolve, reject) => {
    const messageId = crypto.randomUUID();
    message.id = messageId;
    message.timestamp = new Date().toISOString();
    
    // Set up one-time listener for response
    const listener = (response) => {
      if (response.id === messageId) {
        nativePort.onMessage.removeListener(listener);
        resolve(response);
      }
    };
    
    nativePort.onMessage.addListener(listener);
    nativePort.postMessage(message);
    
    // Timeout after 5 seconds
    setTimeout(() => {
      nativePort.onMessage.removeListener(listener);
      reject(new Error('Native messaging timeout'));
    }, 5000);
  });
}

// Initialize on install
chrome.runtime.onInstalled.addListener(() => {
  console.log('Ankhora extension installed');
  connectNative();
});

// Connect on startup
connectNative();
```

### 2.4 Content Script

**File:** `content.js`

```javascript
// Detect login forms on page load
window.addEventListener('load', () => {
  const loginForm = detectLoginForm();
  
  if (loginForm) {
    console.log('Login form detected:', loginForm);
    
    // Notify background script
    chrome.runtime.sendMessage({
      action: 'detect_login',
      payload: {
        url: window.location.href,
        domain: window.location.hostname
      }
    }, (response) => {
      if (response && response.success) {
        const credentials = response.data.credentials;
        if (credentials.length > 0) {
          showAnkhoraIcon(loginForm, credentials);
        }
      }
    });
  }
});

// Detect login form
function detectLoginForm() {
  const usernameField = document.querySelector(
    'input[type="email"], ' +
    'input[type="text"][name*="user" i], ' +
    'input[type="text"][name*="email" i], ' +
    'input[type="text"][id*="user" i], ' +
    'input[type="text"][id*="email" i], ' +
    'input[autocomplete="username"], ' +
    'input[autocomplete="email"]'
  );
  
  const passwordField = document.querySelector('input[type="password"]');
  
  if (usernameField && passwordField) {
    return { usernameField, passwordField };
  }
  
  return null;
}

// Show Ankhora icon next to login form
function showAnkhoraIcon(loginForm, credentials) {
  const { usernameField } = loginForm;
  
  // Create icon container
  const icon = document.createElement('div');
  icon.id = 'ankhora-icon';
  icon.innerHTML = `
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
      <path d="M12 2L2 7v10c0 5.5 3.8 10.7 10 12 6.2-1.3 10-6.5 10-12V7l-10-5z" 
            fill="#D4AF37" stroke="#D4AF37" stroke-width="2"/>
    </svg>
  `;
  icon.style.cssText = `
    position: absolute;
    right: 10px;
    top: 50%;
    transform: translateY(-50%);
    cursor: pointer;
    z-index: 999999;
    width: 24px;
    height: 24px;
  `;
  
  // Position relative to username field
  const fieldRect = usernameField.getBoundingClientRect();
  usernameField.parentElement.style.position = 'relative';
  usernameField.parentElement.appendChild(icon);
  
  // Click handler
  icon.addEventListener('click', () => {
    if (credentials.length === 1) {
      autofillCredential(credentials[0].id);
    } else {
      showCredentialPicker(credentials);
    }
  });
}

// Autofill credential
function autofillCredential(credentialId) {
  chrome.runtime.sendMessage({
    action: 'request_autofill',
    payload: { credential_id: credentialId }
  }, (response) => {
    if (response && response.success) {
      const { username, password } = response.data;
      injectCredentials(username, password);
    }
  });
}

// Inject credentials into form
function injectCredentials(username, password) {
  const loginForm = detectLoginForm();
  if (!loginForm) return;
  
  const { usernameField, passwordField } = loginForm;
  
  // Fill username
  usernameField.value = username;
  usernameField.dispatchEvent(new Event('input', { bubbles: true }));
  usernameField.dispatchEvent(new Event('change', { bubbles: true }));
  
  // Fill password
  passwordField.value = password;
  passwordField.dispatchEvent(new Event('input', { bubbles: true }));
  passwordField.dispatchEvent(new Event('change', { bubbles: true }));
  
  // Focus password field (triggers validation on some sites)
  passwordField.focus();
  
  console.log('Credentials autofilled successfully');
}

// Listen for autofill commands from background
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === 'autofill') {
    injectCredentials(message.username, message.password);
    sendResponse({ success: true });
  }
});

// Detect form submission to offer saving credentials
document.addEventListener('submit', (e) => {
  const form = e.target;
  const loginForm = detectLoginForm();
  
  if (loginForm) {
    const { usernameField, passwordField } = loginForm;
    
    if (usernameField.value && passwordField.value) {
      // Offer to save
      chrome.runtime.sendMessage({
        action: 'save_credential',
        payload: {
          url: window.location.href,
          domain: window.location.hostname,
          username: usernameField.value,
          password: passwordField.value
        }
      });
    }
  }
}, true);
```

---

## 3. Desktop App Overlay UI Design

### 3.1 Overlay Architecture

The desktop app displays a **native overlay window** positioned above the browser when credentials are available.

**Technology Options:**
- **Electron:** Cross-platform, web-based UI (HTML/CSS/JS)
- **Qt/PyQt:** Native widgets, better performance
- **Tauri:** Rust + web UI, lightweight alternative to Electron

**Recommended:** **Electron** for rapid development and consistent UX across platforms.

### 3.2 Overlay Window Specifications

**Dimensions:**
- Width: 400px
- Height: Auto (min 150px, max 500px)
- Position: Top-right corner of active browser window (20px margin)

**Appearance:**
- Frameless window (no title bar)
- Always on top (z-index above browser)
- Semi-transparent background (90% opacity)
- Rounded corners (12px border-radius)
- Drop shadow (0 10px 40px rgba(0,0,0,0.3))

**Behavior:**
- Auto-dismiss after 10 seconds of inactivity
- Dismiss on ESC key
- Dismiss on click outside
- Keyboard navigation (arrow keys, Enter)

### 3.3 Overlay UI Components

**Header:**
- Ankhora logo (gold ankh icon)
- Current vault indicator ("Work Vault")
- Close button (X)

**Credential List:**
- Each credential shows:
  - Service name (e.g., "GitHub")
  - Username (e.g., "manuel@blackops.io")
  - Vault badge (e.g., "Work")
  - Verification status (✓ Verified on Stellar)
  - Last used timestamp (e.g., "Used 2 hours ago")
- Hover state (gold highlight)
- Selected state (gold border)

**Footer:**
- Keyboard hints ("↵ Autofill • ⌘C Copy • ESC Close")
- Quick actions (Generate password, Open vault)

### 3.4 Overlay HTML/CSS (Electron)

**File:** `overlay.html`

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Ankhora Overlay</title>
  <style>
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }
    
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
      background: rgba(10, 10, 10, 0.95);
      color: #fff;
      width: 400px;
      border-radius: 12px;
      border: 1px solid rgba(212, 175, 55, 0.3);
      box-shadow: 0 10px 40px rgba(0, 0, 0, 0.5);
      overflow: hidden;
    }
    
    .header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 12px 16px;
      border-bottom: 1px solid rgba(212, 175, 55, 0.2);
    }
    
    .header-left {
      display: flex;
      align-items: center;
      gap: 10px;
    }
    
    .logo {
      width: 24px;
      height: 24px;
    }
    
    .vault-name {
      font-size: 14px;
      font-weight: 600;
      color: #D4AF37;
    }
    
    .close-btn {
      background: none;
      border: none;
      color: #666;
      font-size: 20px;
      cursor: pointer;
      padding: 0;
      width: 24px;
      height: 24px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    
    .close-btn:hover {
      color: #fff;
    }
    
    .credentials {
      max-height: 400px;
      overflow-y: auto;
      padding: 8px;
    }
    
    .credential {
      padding: 12px;
      border-radius: 8px;
      cursor: pointer;
      transition: background 0.15s;
      margin-bottom: 4px;
    }
    
    .credential:hover {
      background: rgba(212, 175, 55, 0.1);
    }
    
    .credential.selected {
      background: rgba(212, 175, 55, 0.15);
      border: 1px solid rgba(212, 175, 55, 0.4);
    }
    
    .credential-name {
      font-size: 14px;
      font-weight: 600;
      margin-bottom: 4px;
      display: flex;
      align-items: center;
      gap: 8px;
    }
    
    .credential-username {
      font-size: 12px;
      color: #999;
      margin-bottom: 4px;
    }
    
    .credential-meta {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 11px;
    }
    
    .vault-badge {
      background: rgba(212, 175, 55, 0.2);
      color: #D4AF37;
      padding: 2px 8px;
      border-radius: 4px;
    }
    
    .verified {
      color: #4CAF50;
    }
    
    .last-used {
      color: #666;
    }
    
    .footer {
      padding: 12px 16px;
      border-top: 1px solid rgba(212, 175, 55, 0.2);
      font-size: 11px;
      color: #666;
      text-align: center;
    }
    
    .empty-state {
      padding: 40px 20px;
      text-align: center;
      color: #666;
    }
    
    .empty-state-icon {
      font-size: 48px;
      margin-bottom: 12px;
      opacity: 0.3;
    }
  </style>
</head>
<body>
  <div class="header">
    <div class="header-left">
      <svg class="logo" viewBox="0 0 24 24" fill="none">
        <path d="M12 2L2 7v10c0 5.5 3.8 10.7 10 12 6.2-1.3 10-6.5 10-12V7l-10-5z" 
              fill="#D4AF37"/>
      </svg>
      <span class="vault-name" id="vaultName">Personal Vault</span>
    </div>
    <button class="close-btn" id="closeBtn">×</button>
  </div>
  
  <div class="credentials" id="credentialsList">
    <!-- Populated by JavaScript -->
  </div>
  
  <div class="footer">
    Press ↵ to autofill • ⌘C to copy • ESC to close
  </div>
  
  <script src="overlay.js"></script>
</body>
</html>
```

**File:** `overlay.js`

```javascript
const { ipcRenderer } = require('electron');

let credentials = [];
let selectedIndex = 0;

// Receive credentials from main process
ipcRenderer.on('show-credentials', (event, data) => {
  credentials = data.credentials;
  renderCredentials();
});

// Render credentials list
function renderCredentials() {
  const container = document.getElementById('credentialsList');
  
  if (credentials.length === 0) {
    container.innerHTML = `
      <div class="empty-state">
        <div class="empty-state-icon">🔒</div>
        <div>No credentials found for this site</div>
      </div>
    `;
    return;
  }
  
  container.innerHTML = credentials.map((cred, index) => `
    <div class="credential ${index === selectedIndex ? 'selected' : ''}" 
         data-index="${index}">
      <div class="credential-name">
        ${cred.name}
      </div>
      <div class="credential-username">${cred.username}</div>
      <div class="credential-meta">
        <span class="vault-badge">${cred.vault}</span>
        ${cred.verified ? '<span class="verified">✓ Verified</span>' : ''}
        <span class="last-used">${formatLastUsed(cred.last_used)}</span>
      </div>
    </div>
  `).join('');
  
  // Add click handlers
  document.querySelectorAll('.credential').forEach(el => {
    el.addEventListener('click', () => {
      const index = parseInt(el.dataset.index);
      autofillCredential(credentials[index]);
    });
  });
}

// Format last used timestamp
function formatLastUsed(timestamp) {
  if (!timestamp) return 'Never used';
  
  const date = new Date(timestamp);
  const now = new Date();
  const diff = now - date;
  
  const hours = Math.floor(diff / (1000 * 60 * 60));
  if (hours < 1) return 'Just now';
  if (hours < 24) return `${hours} hours ago`;
  
  const days = Math.floor(hours / 24);
  if (days < 7) return `${days} days ago`;
  
  return date.toLocaleDateString();
}

// Keyboard navigation
document.addEventListener('keydown', (e) => {
  if (e.key === 'ArrowDown') {
    e.preventDefault();
    selectedIndex = Math.min(selectedIndex + 1, credentials.length - 1);
    renderCredentials();
  } else if (e.key === 'ArrowUp') {
    e.preventDefault();
```javascript
    selectedIndex = Math.max(selectedIndex - 1, 0);
    renderCredentials();
  } else if (e.key === 'Enter') {
    e.preventDefault();
    autofillCredential(credentials[selectedIndex]);
  } else if (e.key === 'Escape') {
    closeOverlay();
  } else if (e.key === 'c' && (e.metaKey || e.ctrlKey)) {
    e.preventDefault();
    copyPassword(credentials[selectedIndex]);
  }
});

// Autofill credential
function autofillCredential(credential) {
  ipcRenderer.send('autofill-credential', credential.id);
  closeOverlay();
}

// Copy password to clipboard
function copyPassword(credential) {
  ipcRenderer.send('copy-password', credential.id);
  showNotification('Password copied to clipboard');
}

// Show notification
function showNotification(message) {
  // Could show a toast notification here
  console.log(message);
}

// Close overlay
function closeOverlay() {
  ipcRenderer.send('close-overlay');
}

// Close button handler
document.getElementById('closeBtn').addEventListener('click', closeOverlay);

// Auto-dismiss after 10 seconds
setTimeout(closeOverlay, 10000);
```

### 3.5 Electron Main Process (Overlay Window Management)

**File:** `main.js` (excerpt)

```javascript
const { app, BrowserWindow, ipcMain, screen } = require('electron');
const path = require('path');

let overlayWindow = null;

// Create overlay window
function createOverlay(credentials, browserWindowBounds) {
  if (overlayWindow) {
    overlayWindow.close();
  }
  
  const primaryDisplay = screen.getPrimaryDisplay();
  const { width, height } = primaryDisplay.workAreaSize;
  
  // Position overlay at top-right of browser window
  const x = browserWindowBounds.x + browserWindowBounds.width - 420; // 400px + 20px margin
  const y = browserWindowBounds.y + 80; // Below browser toolbar
  
  overlayWindow = new BrowserWindow({
    width: 400,
    height: Math.min(credentials.length * 80 + 100, 500),
    x: x,
    y: y,
    frame: false,
    transparent: true,
    alwaysOnTop: true,
    resizable: false,
    skipTaskbar: true,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false,
      preload: path.join(__dirname, 'preload.js')
    }
  });
  
  overlayWindow.loadFile('overlay.html');
  
  // Send credentials to overlay
  overlayWindow.webContents.on('did-finish-load', () => {
    overlayWindow.webContents.send('show-credentials', {
      credentials: credentials
    });
  });
  
  // Close on blur (click outside)
  overlayWindow.on('blur', () => {
    overlayWindow.close();
  });
  
  overlayWindow.on('closed', () => {
    overlayWindow = null;
  });
}

// IPC handlers
ipcMain.on('autofill-credential', (event, credentialId) => {
  // Decrypt credential and send to extension
  const credential = vaultManager.getCredential(credentialId);
  const decrypted = vaultManager.decrypt(credential);
  
  // Send to extension via native messaging
  sendToExtension({
    action: 'autofill_command',
    payload: {
      username: decrypted.username,
      password: decrypted.password
    }
  });
});

ipcMain.on('copy-password', (event, credentialId) => {
  const credential = vaultManager.getCredential(credentialId);
  const decrypted = vaultManager.decrypt(credential);
  
  // Copy to clipboard
  clipboard.writeText(decrypted.password);
  
  // Clear clipboard after 30 seconds
  setTimeout(() => {
    if (clipboard.readText() === decrypted.password) {
      clipboard.clear();
    }
  }, 30000);
});

ipcMain.on('close-overlay', () => {
  if (overlayWindow) {
    overlayWindow.close();
  }
});

// Export for use in native messaging handler
module.exports = { createOverlay };
```

---

## 4. Autofill Implementation: Clipboard vs Accessibility API

### 4.1 Comparison

| Method | **Clipboard** | **Accessibility API** |
|--------|--------------|----------------------|
| **Complexity** | Simple | Complex |
| **Speed** | Fast (50-100ms) | Medium (100-200ms) |
| **Reliability** | High (works everywhere) | Medium (some apps block it) |
| **Security** | Low (clipboard can be monitored) | High (direct input injection) |
| **User Experience** | Requires paste action | Seamless autofill |
| **Cross-platform** | Easy | Platform-specific code |
| **Browser support** | Universal | Limited |

### 4.2 Recommended Approach: **Hybrid**

**Primary:** Extension injection (direct DOM manipulation)  
**Fallback 1:** Clipboard (for compatibility)  
**Fallback 2:** Accessibility API (for native apps)

### 4.3 Method 1: Extension Injection (Primary)

**Already implemented in content.js above**

**Pros:**
- ✅ Fastest (direct DOM access)
- ✅ Most reliable for web forms
- ✅ No clipboard pollution
- ✅ Works with React/Vue/Angular forms

**Cons:**
- ❌ Only works in browser
- ❌ Can be blocked by strict CSP
- ❌ Doesn't work in iframes with different origin

**Implementation:** See `injectCredentials()` in content.js above.

### 4.4 Method 2: Clipboard (Fallback)

**Use when:** Extension injection fails or user manually requests copy.

**Python Implementation (Desktop App):**

```python
import pyperclip
import time

def autofill_via_clipboard(username, password):
    """
    Copy credentials to clipboard and simulate paste.
    """
    # Copy username
    pyperclip.copy(username)
    time.sleep(0.1)
    
    # Simulate Cmd+V (macOS) or Ctrl+V (Windows/Linux)
    if sys.platform == 'darwin':
        keyboard.press(Key.cmd)
        keyboard.press('v')
        keyboard.release('v')
        keyboard.release(Key.cmd)
    else:
        keyboard.press(Key.ctrl)
        keyboard.press('v')
        keyboard.release('v')
        keyboard.release(Key.ctrl)
    
    time.sleep(0.2)
    
    # Tab to password field
    keyboard.press(Key.tab)
    keyboard.release(Key.tab)
    
    time.sleep(0.1)
    
    # Copy password
    pyperclip.copy(password)
    time.sleep(0.1)
    
    # Simulate paste
    if sys.platform == 'darwin':
        keyboard.press(Key.cmd)
        keyboard.press('v')
        keyboard.release('v')
        keyboard.release(Key.cmd)
    else:
        keyboard.press(Key.ctrl)
        keyboard.press('v')
        keyboard.release('v')
        keyboard.release(Key.ctrl)
    
    # Clear clipboard after 30 seconds
    threading.Timer(30.0, lambda: pyperclip.copy('')).start()
```

**Pros:**
- ✅ Works everywhere (browser, native apps, terminals)
- ✅ Simple implementation
- ✅ Cross-platform

**Cons:**
- ❌ Clipboard can be monitored by malware
- ❌ Overwrites user's clipboard
- ❌ Requires Tab key to switch fields (not always reliable)
- ❌ Slower user experience

### 4.5 Method 3: Accessibility API (Advanced Fallback)

**Use when:** Autofilling native desktop apps (not browser).

#### macOS: Accessibility API

```python
from AppKit import NSWorkspace, NSAppleScript
from Quartz import CGEventCreateKeyboardEvent, CGEventPost, kCGHIDEventTap

def autofill_via_accessibility_macos(username, password):
    """
    Use macOS Accessibility API to inject credentials.
    Requires: System Preferences > Security > Privacy > Accessibility
    """
    # Get active window
    active_app = NSWorkspace.sharedWorkspace().activeApplication()
    
    # Type username
    type_string(username)
    
    # Press Tab
    press_key(48)  # Tab key code
    
    # Type password
    type_string(password)

def type_string(text):
    """Type a string using CGEvent"""
    for char in text:
        # Create key down event
        event_down = CGEventCreateKeyboardEvent(None, ord(char), True)
        CGEventPost(kCGHIDEventTap, event_down)
        
        # Create key up event
        event_up = CGEventCreateKeyboardEvent(None, ord(char), False)
        CGEventPost(kCGHIDEventTap, event_up)
        
        time.sleep(0.01)

def press_key(keycode):
    """Press a key by keycode"""
    event_down = CGEventCreateKeyboardEvent(None, keycode, True)
    CGEventPost(kCGHIDEventTap, event_down)
    
    event_up = CGEventCreateKeyboardEvent(None, keycode, False)
    CGEventPost(kCGHIDEventTap, event_up)
```

#### Windows: UI Automation

```python
import pywinauto
from pywinauto.keyboard import send_keys

def autofill_via_uiautomation_windows(username, password):
    """
    Use Windows UI Automation to inject credentials.
    """
    # Get active window
    app = pywinauto.Application().connect(active_only=True)
    
    # Find username field (by control type)
    try:
        username_field = app.top_window().child_window(
            control_type="Edit", 
            found_index=0
        )
        username_field.set_focus()
        send_keys(username)
        
        # Tab to password field
        send_keys('{TAB}')
        
        # Type password
        send_keys(password)
        
    except Exception as e:
        print(f"UI Automation failed: {e}")
        # Fallback to clipboard method
        autofill_via_clipboard(username, password)
```

#### Linux: xdotool

```python
import subprocess

def autofill_via_xdotool_linux(username, password):
    """
    Use xdotool to inject credentials on Linux.
    Requires: apt-get install xdotool
    """
    # Type username
    subprocess.run(['xdotool', 'type', username])
    
    # Press Tab
    subprocess.run(['xdotool', 'key', 'Tab'])
    
    # Type password
    subprocess.run(['xdotool', 'type', password])
```

**Pros:**
- ✅ Works with native desktop apps
- ✅ More secure than clipboard
- ✅ Can target specific UI elements

**Cons:**
- ❌ Requires elevated permissions (Accessibility on macOS)
- ❌ Platform-specific code
- ❌ Can be blocked by security software
- ❌ More complex to implement

### 4.6 Recommended Implementation Strategy

```python
class AutofillManager:
    def __init__(self):
        self.methods = [
            self.autofill_via_extension,  # Primary
            self.autofill_via_clipboard,  # Fallback 1
            self.autofill_via_accessibility  # Fallback 2
        ]
    
    def autofill(self, credential, context):
        """
        Try autofill methods in order until one succeeds.
        """
        for method in self.methods:
            try:
                success = method(credential, context)
                if success:
                    return True
            except Exception as e:
                logging.warning(f"{method.__name__} failed: {e}")
                continue
        
        # All methods failed
        logging.error("All autofill methods failed")
        return False
    
    def autofill_via_extension(self, credential, context):
        """Send autofill command to browser extension"""
        if context['type'] != 'browser':
            return False
        
        send_to_extension({
            'action': 'autofill_command',
            'payload': {
                'username': credential.username,
                'password': credential.password
            }
        })
        return True
    
    def autofill_via_clipboard(self, credential, context):
        """Copy to clipboard and simulate paste"""
        autofill_via_clipboard(credential.username, credential.password)
        return True
    
    def autofill_via_accessibility(self, credential, context):
        """Use OS accessibility API"""
        if sys.platform == 'darwin':
            return autofill_via_accessibility_macos(
                credential.username, 
                credential.password
            )
        elif sys.platform == 'win32':
            return autofill_via_uiautomation_windows(
                credential.username, 
                credential.password
            )
        elif sys.platform == 'linux':
            return autofill_via_xdotool_linux(
                credential.username, 
                credential.password
            )
        return False
```

---

## 5. Security Considerations

### 5.1 Credential Transmission

**Problem:** Credentials must pass from desktop app → extension → web page.

**Solution:**
1. **Desktop app:** Decrypt credentials in memory (never write to disk)
2. **Native messaging:** Credentials transmitted via stdin/stdout (not network)
3. **Extension:** Credentials held in memory for <1 second, then cleared
4. **Web page:** Credentials injected directly into form fields (never stored in extension storage)

### 5.2 Memory Protection

```python
import ctypes
from ctypes import c_char_p, c_size_t

def secure_string(data):
    """
    Store sensitive data in protected memory (prevents swapping to disk).
    """
    if sys.platform == 'darwin' or sys.platform == 'linux':
        # Use mlock to prevent swapping
        libc = ctypes.CDLL('libc.so.6' if sys.platform == 'linux' else 'libc.dylib')
        libc.mlock(c_char_p(data.encode()), c_size_t(len(data)))
    
    return data

def secure_clear(data):
    """
    Overwrite memory before deletion.
    """
    if isinstance(data, str):
        data = '\x00' * len(data)
```

### 5.3 Extension Isolation

**Content Security Policy (CSP):**
- Extension runs in isolated world (separate from web page JavaScript)
- Web page cannot access extension APIs or memory
- Extension cannot be inspected by web page

**Manifest V3 Security:**
- Service worker runs in separate process
- No eval() or inline scripts allowed
- All resources must be declared in manifest

---

## 6. Testing Plan

### 6.1 Unit Tests

**Native Messaging:**
- ✅ Message serialization/deserialization
- ✅ Timeout handling
- ✅ Error responses
- ✅ Reconnection logic

**Form Detection:**
- ✅ Detect standard login forms
- ✅ Detect multi-step logins
- ✅ Detect hidden fields
- ✅ Detect React/Vue/Angular forms

**Autofill:**
- ✅ Inject into standard forms
- ✅ Inject into shadow DOM
- ✅ Trigger validation events
- ✅ Handle auto-submit

### 6.2 Integration Tests

**Desktop App ↔ Extension:**
- ✅ Extension connects to native host
- ✅ Extension sends detect_login → receives credentials
- ✅ Extension requests autofill → credentials injected
- ✅ Extension saves credential → stored in vault

**Overlay UI:**
- ✅ Overlay appears at correct position
- ✅ Keyboard navigation works
- ✅ Autofill on Enter
- ✅ Auto-dismiss after timeout

### 6.3 E2E Tests (Real Websites)

Test autofill on:
- ✅ GitHub (standard form)
- ✅ Google (multi-step login)
- ✅ Facebook (React-based)
- ✅ Twitter (dynamic form)
- ✅ Banking sites (strict CSP)

### 6.4 Performance Tests

- ✅ Hotkey detection latency <10ms
- ✅ Vault search <100ms for 1000+ credentials
- ✅ Overlay render <200ms
- ✅ Autofill injection <50ms

---

## 7. Rollout Timeline

**Week 3:**
- ✅ Native messaging protocol implementation
- ✅ Extension manifest and file structure
- ✅ Basic form detection

**Week 4:**
- ✅ Desktop app overlay UI
- ✅ Autofill injection (extension method)
- ✅ Credential saving flow

**Week 5:**
- ✅ Clipboard fallback method
- ✅ Cross-browser testing (Chrome, Firefox)
- ✅ Bug fixes and polish

**Week 6:**
- ✅ Beta testing with early users
- ✅ Performance optimization
- ✅ Documentation

---

## 8. Success Metrics

- **Adoption:** >80% of desktop users install extension within 1 week
- **Autofill Success Rate:** >95% on top 100 websites
- **Latency:** 95th percentile autofill <500ms
- **User Satisfaction:** >4.5/5 rating
- **Retention:** Extension users have 3x higher 30-day retention

---

**This completes the Phase 1 Bridge Extension technical specification.** Ready to start implementation! 🚀