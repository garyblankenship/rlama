const { contextBridge, ipcRenderer } = require('electron');

// Exposer des APIs sécurisées au renderer process
contextBridge.exposeInMainWorld('electron', {
  // Utilitaires de fichiers
  selectDirectory: () => ipcRenderer.invoke('select-directory'),
  
  // Gestion du backend
  isBackendReady: () => ipcRenderer.invoke('is-backend-ready'),
  
  // Contrôles de fenêtre
  windowMinimize: () => ipcRenderer.invoke('window-minimize'),
  windowMaximize: () => ipcRenderer.invoke('window-maximize'),
  windowClose: () => ipcRenderer.invoke('window-close'),
  windowIsMaximized: () => ipcRenderer.invoke('window-is-maximized'),
  
  // Événements
  onBackendReady: (callback) => {
    ipcRenderer.on('backend-ready', (_, value) => callback(value));
    return () => ipcRenderer.removeAllListeners('backend-ready');
  },
  executeCommand: (command, args) => ipcRenderer.invoke('execute-command', command, args)

});

contextBridge.exposeInMainWorld('electronAPI', {
  // ... existing APIs ...
  
  // APIs de mise à jour
  checkForUpdates: () => ipcRenderer.invoke('check-for-updates'),
  quitAndInstall: () => ipcRenderer.invoke('quit-and-install'),
  getAppVersion: () => ipcRenderer.invoke('get-app-version'),
  
  // Écouter les événements de mise à jour
  onDownloadProgress: (callback) => {
    ipcRenderer.on('download-progress', (event, progress) => callback(progress));
  },
  
  removeDownloadProgressListener: () => {
    ipcRenderer.removeAllListeners('download-progress');
  }
}); 
