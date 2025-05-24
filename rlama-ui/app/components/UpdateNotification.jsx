import React, { useState, useEffect } from 'react';
import { Button, Progress, Modal, Typography, Space } from 'antd';
import { DownloadOutlined, ReloadOutlined } from '@ant-design/icons';

const { Text, Title } = Typography;

const UpdateNotification = () => {
  const [updateAvailable, setUpdateAvailable] = useState(false);
  const [downloading, setDownloading] = useState(false);
  const [downloadProgress, setDownloadProgress] = useState(0);
  const [updateReady, setUpdateReady] = useState(false);
  const [currentVersion, setCurrentVersion] = useState('');

  useEffect(() => {
    // Obtenir la version actuelle
    window.electronAPI?.getAppVersion().then(version => {
      setCurrentVersion(version);
    });

    // Écouter les événements de progression du téléchargement
    if (window.electronAPI?.onDownloadProgress) {
      window.electronAPI.onDownloadProgress((progress) => {
        setDownloading(true);
        setDownloadProgress(Math.round(progress.percent));
        
        if (progress.percent === 100) {
          setDownloading(false);
          setUpdateReady(true);
        }
      });
    }

    // Cleanup
    return () => {
      if (window.electronAPI?.removeDownloadProgressListener) {
        window.electronAPI.removeDownloadProgressListener();
      }
    };
  }, []);

  const handleCheckUpdates = async () => {
    if (window.electronAPI?.checkForUpdates) {
      await window.electronAPI.checkForUpdates();
    }
  };

  const handleInstallUpdate = async () => {
    if (window.electronAPI?.quitAndInstall) {
      await window.electronAPI.quitAndInstall();
    }
  };

  return (
    <div style={{ padding: '16px' }}>
      <Space direction="vertical" style={{ width: '100%' }}>
        <div>
          <Text strong>Version actuelle: {currentVersion}</Text>
        </div>
        
        <Button 
          type="primary" 
          icon={<ReloadOutlined />}
          onClick={handleCheckUpdates}
          disabled={downloading}
        >
          Vérifier les mises à jour
        </Button>

        {downloading && (
          <div>
            <Text>Téléchargement en cours...</Text>
            <Progress percent={downloadProgress} />
          </div>
        )}

        {updateReady && (
          <Modal
            title="Mise à jour prête"
            open={true}
            onOk={handleInstallUpdate}
            onCancel={() => setUpdateReady(false)}
            okText="Installer et redémarrer"
            cancelText="Plus tard"
          >
            <p>La mise à jour a été téléchargée et est prête à être installée.</p>
            <p>L'application va redémarrer pour appliquer la mise à jour.</p>
          </Modal>
        )}
      </Space>
    </div>
  );
};

export default UpdateNotification; 