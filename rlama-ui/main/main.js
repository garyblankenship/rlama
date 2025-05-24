const { app, BrowserWindow, dialog, ipcMain } = require('electron');
const path = require('path');
const { spawn, exec } = require('child_process');
const fs = require('fs');
const axios = require('axios');
const sudo = require('sudo-prompt');
const AppUpdater = require('./updater');

let mainWindow;
let pythonProcess;
let updater;
const BACKEND_PORT = 5001;
const BACKEND_URL = `http://127.0.0.1:${BACKEND_PORT}`;
let backendReady = false;
let healthCheckAttempts = 0;
const MAX_HEALTH_CHECK_ATTEMPTS = 10;

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    frame: false, // Remove native frame
    titleBarStyle: 'hidden', // Hide title bar on macOS
    vibrancy: 'under-window', // macOS vibrancy effect
    icon: path.join(__dirname, '../public/logo.png'), // Add app icon
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js'),
    },
  });

  // Try different possible paths for the HTML file
  const possiblePaths = [
    path.join(__dirname, '../dist/index.html'),
    path.join(__dirname, '../index.html'),
    path.join(app.getAppPath(), 'dist/index.html'),
    path.join(process.cwd(), 'dist/index.html')
  ];

  // Log all possible paths for debugging
  console.log('Searching for HTML file at these locations:');
  possiblePaths.forEach(p => console.log(`- ${p} (exists: ${fs.existsSync(p)})`));

  // Find the first path that exists
  const htmlPath = possiblePaths.find(p => fs.existsSync(p));

  if (!htmlPath) {
    console.error('ERROR: HTML file not found at any of the expected locations');
    dialog.showErrorBox(
      'Application Error',
      `Cannot find HTML file at any of the expected locations. 
      Current directory: ${__dirname}
      App path: ${app.getAppPath()}
      Working dir: ${process.cwd()}`
    );
    return;
  }

  console.log(`Loading HTML from: ${htmlPath}`);
  mainWindow.loadFile(htmlPath).catch(err => {
    console.error('Failed to load HTML file:', err);
    dialog.showErrorBox(
      'Application Error',
      `Failed to load HTML file: ${err.message}`
    );
  });

  // Open DevTools only in development
  if (process.env.NODE_ENV === 'development') {
    mainWindow.webContents.openDevTools();
  }

  mainWindow.on('closed', () => {
    mainWindow = null;
  });

  // Initialiser le système de mise à jour après la création de la fenêtre
  if (!process.env.NODE_ENV || process.env.NODE_ENV === 'production') {
    updater = new AppUpdater(mainWindow);
  }
}

async function startPythonBackend() {
  const scriptPath = path.join(__dirname, '../backend/app.py');

  console.log(`Looking for Python script at: ${scriptPath}`);
  if (!fs.existsSync(scriptPath)) {
    dialog.showErrorBox(
      'Erreur Backend',
      `Le script Python n'a pas été trouvé: ${scriptPath}`
    );
    app.quit();
    return;
  }

  // Try different Python commands - starting with conda environment
  const pythonCommands = [
    // Try conda environment paths first
    '/Users/dontizi/anaconda3/envs/reMind/bin/python',
    '/Users/dontizi/anaconda3/bin/python',
    // Then try system Python
    'python3',
    'python'
  ];
  let pythonCommand = null;

  for (const cmd of pythonCommands) {
    try {
      // Use spawn sync to test if python command exists
      console.log(`Testing Python command: ${cmd}`);
      const result = require('child_process').spawnSync(cmd, ['--version']);
      if (result.status === 0) {
        pythonCommand = cmd;
        console.log(`Found Python command: ${cmd}`);
        break;
      }
    } catch (error) {
      console.log(`Command ${cmd} not found or failed: ${error.message}`);
    }
  }

  if (!pythonCommand) {
    dialog.showErrorBox(
      'Erreur Python',
      'Python was not found. Please install Python 3.x.'
    );
    app.quit();
    return;
  }

  // Lancer le processus Python avec debug info
  console.log(`Starting Python backend with: ${pythonCommand} ${scriptPath}`);
  pythonProcess = spawn(pythonCommand, [scriptPath]);
  console.log('Python process spawned with PID:', pythonProcess.pid);

  pythonProcess.stdout.on('data', (data) => {
    console.log(`Backend stdout: ${data}`);
    if (data.toString().includes('Backend server started') || data.toString().includes('Uvicorn running on')) {
      checkBackendHealth();
    }
  });

  pythonProcess.stderr.on('data', (data) => {
    console.error(`Backend stderr: ${data}`);
    if (data.toString().includes('Uvicorn running on')) {
      checkBackendHealth();
    }
  });

  pythonProcess.on('close', (code) => {
    console.log(`Backend process exited with code ${code}`);
    if (code !== 0 && mainWindow) {
      dialog.showErrorBox(
        'Erreur Backend',
        `Le processus Python s'est arrêté avec le code ${code}`
      );
    }
  });

  pythonProcess.on('error', (err) => {
    console.error(`Failed to start Python process: ${err}`);
    dialog.showErrorBox(
      'Erreur Backend',
      `Impossible de démarrer le backend Python: ${err.message}`
    );
  });
}

async function checkBackendHealth() {
  console.log('Checking backend health...');
  console.log(`Current status - backendReady: ${backendReady}, attempts: ${healthCheckAttempts}`);
  
  if (backendReady) {
    console.log('Backend is already ready, skipping health check');
    return;
  }

  if (healthCheckAttempts >= MAX_HEALTH_CHECK_ATTEMPTS) {
    console.error('Backend health check failed after maximum attempts');
    dialog.showErrorBox(
      'Erreur Backend',
      'The backend is not responding after several attempts. Please restart the application.'
    );
    return;
  }
  
  try {
    console.log(`Attempting to connect to backend at ${BACKEND_URL}/health`);
    const response = await axios.get(`${BACKEND_URL}/health`);
    console.log('Health check response:', response.data);
    
    if (response.status === 200) {
      console.log('Backend is ready, setting backendReady to true');
      backendReady = true;
      if (mainWindow) {
        console.log('Sending backend-ready event to renderer');
        mainWindow.webContents.send('backend-ready', true);
      } else {
        console.warn('mainWindow is null, cannot send backend-ready event');
      }
    }
  } catch (error) {
    console.error('Backend health check failed:', error.message);
    healthCheckAttempts++;
    console.log(`Retrying in 1s (attempt ${healthCheckAttempts}/${MAX_HEALTH_CHECK_ATTEMPTS})`);
    setTimeout(checkBackendHealth, 1000);
  }
}

app.whenReady().then(() => {
  createWindow();
  startPythonBackend();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('before-quit', () => {
  if (pythonProcess) {
    pythonProcess.kill();
  }
});

ipcMain.handle('select-directory', async () => {
  const result = await dialog.showOpenDialog(mainWindow, {
    properties: ['openDirectory']
  });
  return result.filePaths[0];
});

ipcMain.handle('is-backend-ready', () => {
  return backendReady;
});

ipcMain.handle('execute-command', async (event, command, args) => {
  return new Promise((resolve, reject) => {
    if (command === 'sudo') {
      // Pour les commandes nécessitant des privilèges administrateur
      sudo.exec(args.join(' '), {
        name: 'RLAMA Installer'
      }, (error, stdout, stderr) => {
        if (error) reject(error);
        else resolve(stdout);
      });
    } else {
      // Pour les commandes normales
      exec(`${command} ${args.join(' ')}`, (error, stdout, stderr) => {
        if (error) reject(error);
        else resolve(stdout);
      });
    }
  });
});

// Window control handlers
ipcMain.handle('window-minimize', () => {
  if (mainWindow) {
    mainWindow.minimize();
  }
});

ipcMain.handle('window-maximize', () => {
  if (mainWindow) {
    if (mainWindow.isMaximized()) {
      mainWindow.unmaximize();
    } else {
      mainWindow.maximize();
    }
  }
});

ipcMain.handle('window-close', () => {
  if (mainWindow) {
    mainWindow.close();
  }
});

ipcMain.handle('window-is-maximized', () => {
  return mainWindow ? mainWindow.isMaximized() : false;
});

process.on('uncaughtException', (error) => {
  console.error('Uncaught Exception:', error);
  if (mainWindow) {
    dialog.showErrorBox(
      'Unhandled Error',
      `Une erreur non gérée s'est produite: ${error.message}`
    );
  }
});

// Nouveaux IPC handlers pour les mises à jour
ipcMain.handle('check-for-updates', async () => {
  if (updater) {
    updater.checkForUpdates();
    return true;
  }
  return false;
});

ipcMain.handle('quit-and-install', async () => {
  if (updater) {
    updater.quitAndInstall();
  }
});

ipcMain.handle('get-app-version', () => {
  return app.getVersion();
}); 