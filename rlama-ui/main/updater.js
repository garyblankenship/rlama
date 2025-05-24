const { autoUpdater } = require('electron-updater');
const { dialog, BrowserWindow } = require('electron');
const isDev = process.env.NODE_ENV === 'development';

class AppUpdater {
  constructor(mainWindow) {
    this.mainWindow = mainWindow;
    this.setupUpdater();
  }

  setupUpdater() {
    // Configuration pour le développement (optionnel)
    if (isDev) {
      // Vous pouvez configurer un serveur de développement ici
      autoUpdater.updateConfigPath = path.join(__dirname, 'dev-app-update.yml');
    }

    // Configuration des événements
    autoUpdater.checkForUpdatesAndNotify();

    // Nouvelle mise à jour disponible
    autoUpdater.on('update-available', (info) => {
      console.log('Mise à jour disponible:', info);
      this.showUpdateAvailableDialog(info);
    });

    // Pas de mise à jour disponible
    autoUpdater.on('update-not-available', (info) => {
      console.log('Aucune mise à jour disponible:', info);
    });

    // Téléchargement en cours
    autoUpdater.on('download-progress', (progressObj) => {
      const message = `Vitesse: ${progressObj.bytesPerSecond} - Téléchargé ${progressObj.percent}% (${progressObj.transferred}/${progressObj.total})`;
      console.log(message);
      
      // Envoyer le progrès au renderer process
      if (this.mainWindow && !this.mainWindow.isDestroyed()) {
        this.mainWindow.webContents.send('download-progress', progressObj);
      }
    });

    // Mise à jour téléchargée
    autoUpdater.on('update-downloaded', (info) => {
      console.log('Mise à jour téléchargée:', info);
      this.showUpdateDownloadedDialog(info);
    });

    // Erreur de mise à jour
    autoUpdater.on('error', (err) => {
      console.error('Erreur de mise à jour:', err);
      this.showUpdateErrorDialog(err);
    });
  }

  showUpdateAvailableDialog(info) {
    dialog.showMessageBox(this.mainWindow, {
      type: 'info',
      title: 'Mise à jour disponible',
      message: `Une nouvelle version (${info.version}) est disponible!`,
      detail: `Version actuelle: ${autoUpdater.currentVersion}\nNouvelle version: ${info.version}\n\nLa mise à jour va commencer automatiquement.`,
      buttons: ['OK']
    });
  }

  showUpdateDownloadedDialog(info) {
    const response = dialog.showMessageBoxSync(this.mainWindow, {
      type: 'info',
      title: 'Mise à jour prête',
      message: 'La mise à jour a été téléchargée. Redémarrer maintenant?',
      detail: `Version ${info.version} est prête à être installée.`,
      buttons: ['Redémarrer maintenant', 'Plus tard'],
      defaultId: 0,
      cancelId: 1
    });

    if (response === 0) {
      // Fermer tous les processus avant de redémarrer
      autoUpdater.quitAndInstall(false, true);
    }
  }

  showUpdateErrorDialog(error) {
    dialog.showErrorBox(
      'Erreur de mise à jour',
      `Une erreur s'est produite lors de la vérification des mises à jour:\n\n${error.message}`
    );
  }

  // Méthode pour vérifier manuellement les mises à jour
  checkForUpdates() {
    if (isDev) {
      console.log('Mode développement - vérification des mises à jour désactivée');
      return;
    }
    autoUpdater.checkForUpdatesAndNotify();
  }

  // Méthode pour forcer l'installation
  quitAndInstall() {
    autoUpdater.quitAndInstall();
  }
}

module.exports = AppUpdater; 