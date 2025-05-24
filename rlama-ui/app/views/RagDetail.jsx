import React, { useState, useEffect } from 'react';
import { 
  Tabs, Card, Typography, Button, message, 
  Breadcrumb, Spin, Result, Divider, Badge, Tag, Tooltip
} from 'antd';
import { 
  ArrowLeftOutlined, 
  FileOutlined, 
  DatabaseOutlined, 
  QuestionCircleOutlined, 
  SettingOutlined,
  RobotOutlined,
  HomeOutlined,
  LoadingOutlined,
  ReloadOutlined,
  CalendarOutlined,
} from '@ant-design/icons';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { ragService } from '../services/api';

// Import des composants d'onglets
import DocumentsTab from '../components/tabs/DocumentsTab';
import ChunksTab from '../components/tabs/ChunksTab';
import QATab from '../components/tabs/QATab';
import SettingsTab from '../components/tabs/SettingsTab';

const { Title, Text } = Typography;

const RagDetail = () => {
  const { ragName } = useParams();
  const navigate = useNavigate();
  
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [ragInfo, setRagInfo] = useState(null);
  const [activeTab, setActiveTab] = useState('documents');
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  
  // Déclencher un rafraîchissement
  const triggerRefresh = () => {
    setRefreshTrigger(prev => prev + 1);
    message.info('Refreshing data...');
  };
  
  // Charger les informations du RAG
  useEffect(() => {
    const fetchRagInfo = async () => {
      try {
        setLoading(true);
        console.log(`Loading RAG info for "${ragName}"`);
        const allRags = await ragService.getAllRags();
        const currentRag = allRags.find(rag => rag.name === ragName);
        
        if (!currentRag) {
          console.error(`RAG "${ragName}" not found in the list:`, allRags);
          setError(`RAG "${ragName}" not found`);
          return;
        }
        
        // Make a second call to get documents to ensure accurate count
        try {
          const docs = await ragService.getRagDocuments(ragName);
          console.log(`Loaded ${docs.length} documents for size calculation`);
          
          // Calculate total size manually if needed
          if (currentRag.size === "0 B" && docs.length > 0) {
            let totalSize = 0;
            docs.forEach(doc => {
              // Parse size strings like "3.70 KB" to bytes
              const sizeStr = doc.size;
              const sizeNum = parseFloat(sizeStr);
              if (sizeStr.includes("KB")) {
                totalSize += sizeNum * 1024;
              } else if (sizeStr.includes("MB")) {
                totalSize += sizeNum * 1024 * 1024;
              } else if (sizeStr.includes("B")) {
                totalSize += sizeNum;
              }
            });
            
            // Format bytes to human-readable
            const formatSize = (bytes) => {
              if (bytes < 1024) return `${bytes} B`;
              if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
              return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
            };
            
            currentRag.size = formatSize(totalSize);
            currentRag.documents_count = docs.length;
          }
        } catch (docError) {
          console.warn(`Could not fetch documents for ${ragName}:`, docError);
        }
        
        console.log(`RAG info loaded:`, currentRag);
        setRagInfo(currentRag);
        setError(null);
      } catch (error) {
        console.error(`Erreur lors du chargement des infos du RAG "${ragName}":`, error);
        setError(`Impossible de charger les informations du RAG "${ragName}"`);
      } finally {
        setLoading(false);
      }
    };
    
    fetchRagInfo();
  }, [ragName, refreshTrigger]);
  
  // Format de la date
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('fr-FR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };
  
  // Si erreur, afficher un message d'erreur
  if (error) {
    return (
      <div className="fade-in">
        <Result
          status="error"
          title="Erreur"
          subTitle={error}
          extra={[
            <Button 
              key="refresh" 
              type="primary" 
              icon={<ReloadOutlined />} 
              onClick={triggerRefresh}
            >
              Retry
            </Button>,
            <Button 
              key="back" 
              onClick={() => navigate('/')}
              icon={<ArrowLeftOutlined />}
            >
             Back to dashboard
            </Button>
          ]}
        />
      </div>
    );
  }
  
  // Si chargement, afficher un indicateur
  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center p-12">
        <Spin indicator={<LoadingOutlined style={{ fontSize: 36, color: 'var(--primary-600)' }} spin />} />
        <Text className="mt-4" type="secondary">{`Chargement du RAG "${ragName}"...`}</Text>
      </div>
    );
  }
  
  // Items des onglets
  const tabItems = [
    {
      key: 'documents',
      label: (
        <span className="flex items-center gap-2">
          <FileOutlined />
          Documents
          <Badge count={ragInfo.documents_count} style={{ backgroundColor: 'var(--primary-600)' }} />
        </span>
      ),
      children: <DocumentsTab ragName={ragName} refreshParent={triggerRefresh} />,
    },
    {
      key: 'chunks',
      label: (
        <span className="flex items-center gap-2">
          <DatabaseOutlined />
          Chunks
        </span>
      ),
      children: <ChunksTab ragName={ragName} refreshParent={triggerRefresh} />,
    },
    {
      key: 'qa',
      label: (
        <span className="flex items-center gap-2">
          <QuestionCircleOutlined />
          Questions / Answers
        </span>
      ),
      children: <QATab ragName={ragName} model={ragInfo.model} />,
    },
    {
      key: 'settings',
      label: (
        <span className="flex items-center gap-2">
          <SettingOutlined />
          Settings
        </span>
      ),
      children: <SettingsTab ragName={ragName} ragInfo={ragInfo} refreshParent={triggerRefresh} />,
    },
  ];
  
  return (
    <div className="fade-in">
      {/* En-tête */}
      <div className="mb-4">
        <Breadcrumb
          items={[
            { 
              title: (
                <Link to="/" className="flex items-center gap-1">
                  <HomeOutlined /> Dashboard
                </Link>
              ) 
            },
            { title: <span className="text-primary-700">{ragName}</span> },
          ]}
        />
      </div>
      
      <Card className="shadow-md rounded-lg mb-6">
        <div className="flex justify-between items-start">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <Badge status="processing" color="var(--primary-600)" />
              <Title level={2} style={{ margin: 0, color: 'var(--primary-800)' }}>{ragName}</Title>
              
              <Tag color="var(--primary-100)" style={{ color: "var(--primary-700)" }}>
                <RobotOutlined style={{ marginRight: 5 }} /> {ragInfo.model}
              </Tag>
            </div>
            
            <div className="flex items-center gap-4 mt-2">
              <Tooltip title="Nombre de documents">
                <div className="flex items-center gap-1">
                  <FileOutlined style={{ color: 'var(--accent-purple)' }} />
                  <Text>{ragInfo.documents_count} documents</Text>
                </div>
              </Tooltip>
              
              <Tooltip title="Taille totale">
                <Text type="secondary">{ragInfo.size}</Text>
              </Tooltip>
              
              <Tooltip title={`Created on ${formatDate(ragInfo.created_on)}`}>
                <div className="flex items-center gap-1">
                  <CalendarOutlined style={{ color: 'var(--neutral-600)' }} />
                  <Text type="secondary">{formatDate(ragInfo.created_on)}</Text>
                </div>
              </Tooltip>
            </div>
            
            {ragInfo.watch_enabled && (
              <div className="status-badge status-active mt-3">
                <span className="pulse" style={{ width: '8px', height: '8px', borderRadius: '50%', display: 'inline-block' }}></span>
                Active monitoring
              </div>
            )}
          </div>
          
          <Button 
            icon={<ArrowLeftOutlined />} 
            onClick={() => navigate('/')}
          >
            Back
          </Button>
        </div>
      </Card>
      
      {/* Onglets */}
      <Card className="shadow-md rounded-lg">
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={tabItems}
          destroyInactiveTabPane
          size="large"
          animated={{ inkBar: true, tabPane: true }}
          tabBarExtraContent={
            <Button
              icon={<ReloadOutlined />}
              onClick={triggerRefresh}
              type="text"
            >
              Refresh
            </Button>
          }
        />
      </Card>
    </div>
  );
};

export default RagDetail; 