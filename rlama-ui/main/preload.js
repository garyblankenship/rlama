const { contextBridge, ipcRenderer } = require('electron');

// Exposer des APIs sécurisées au renderer process
contextBridge.exposeInMainWorld('electron', {
  // Utilitaires de fichiers
  selectDirectory: () => ipcRenderer.invoke('select-directory'),
  
  // Gestion du backend
  isBackendReady: () => ipcRenderer.invoke('is-backend-ready'),
  
  // Événements
  onBackendReady: (callback) => {
    ipcRenderer.on('backend-ready', (_, value) => callback(value));
    return () => ipcRenderer.removeAllListeners('backend-ready');
  },
  executeCommand: (command, args) => ipcRenderer.invoke('execute-command', command, args)

}); 
