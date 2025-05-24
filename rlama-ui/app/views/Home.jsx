import React, { useState, useEffect, useRef } from 'react';
import { Card, Button, Typography, Row, Col, Statistic, Spin, Alert, Space, Divider, Modal, List } from 'antd';
import { 
  CheckCircleOutlined, 
  CloseCircleOutlined, 
  ExclamationCircleOutlined, 
  SyncOutlined,
  DatabaseOutlined,
  RobotOutlined,
  ApiOutlined,
  PlusOutlined,
  FileTextOutlined,
  BlockOutlined,
  EyeOutlined,
  GlobalOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { healthService, ragService } from '../services/api';
import api from '../services/api';

const { Title, Text, Paragraph } = Typography;

// Composant pour afficher les statistiques du dashboard
const DashboardStats = () => {
  const [stats, setStats] = useState({
    totalRags: 0,
    totalDocuments: 0,
    totalChunks: 0,
    watchedFolders: 0,
    watchedWebsites: 0,
    loading: true,
    error: null
  });

  const loadDashboardStats = async () => {
    try {
      setStats(prev => ({ ...prev, loading: true, error: null }));
      
      // Récupérer la liste des RAGs
      const ragsResponse = await api.get('/rags');
      const rags = ragsResponse.data || [];
      
      let totalDocuments = 0;
      let totalChunks = 0;
      let watchedFolders = 0;
      let watchedWebsites = 0;
      
      // Pour chaque RAG, récupérer les détails
      for (const rag of rags) {
        try {
          // Compter les documents
          const docsResponse = await api.get(`/rags/${rag.name}/documents`);
          const documents = docsResponse.data || [];
          totalDocuments += documents.length;
          
          // Compter les chunks
          const chunksResponse = await api.get(`/rags/${rag.name}/chunks`);
          const chunks = chunksResponse.data || [];
          totalChunks += chunks.length;
          
          // Vérifier le statut de surveillance des dossiers
          try {
            const watchStatusResponse = await api.get(`/rags/${rag.name}/watch-status`);
            if (watchStatusResponse.data && watchStatusResponse.data.active) {
              watchedFolders++;
            }
          } catch (e) {
            // Ignore les erreurs de statut de surveillance
          }
          
          // Vérifier le statut de surveillance web
          try {
            const webWatchStatusResponse = await api.get(`/rags/${rag.name}/web-watch-status`);
            if (webWatchStatusResponse.data && webWatchStatusResponse.data.active) {
              watchedWebsites++;
            }
          } catch (e) {
            // Ignore les erreurs de statut de surveillance web
          }
        } catch (e) {
          console.warn(`Erreur lors de la récupération des détails pour le RAG ${rag.name}:`, e);
        }
      }
      
      setStats({
        totalRags: rags.length,
        totalDocuments,
        totalChunks,
        watchedFolders,
        watchedWebsites,
        loading: false,
        error: null
      });
    } catch (error) {
      console.error('Erreur lors du chargement des statistiques:', error);
      setStats(prev => ({
        ...prev,
        loading: false,
        error: 'Failed to load dashboard statistics'
      }));
    }
  };

  useEffect(() => {
    loadDashboardStats();
  }, []);

  if (stats.loading) {
    return (
      <div style={{ textAlign: 'center', padding: '40px 0' }}>
        <Spin size="large" />
        <Text type="secondary" style={{ display: 'block', marginTop: '16px' }}>
          Loading statistics...
        </Text>
      </div>
    );
  }

  if (stats.error) {
    return (
      <Alert
        message="Error loading statistics"
        description={stats.error}
        type="error"
        showIcon
        action={
          <Button size="small" onClick={loadDashboardStats}>
            Retry
          </Button>
        }
      />
    );
  }

  return (
    <Row gutter={[16, 16]}>
      <Col xs={12} sm={8} md={6} lg={4}>
        <Statistic
          title="RAG Systems"
          value={stats.totalRags}
          prefix={<DatabaseOutlined style={{ color: 'var(--primary-600)' }} />}
          valueStyle={{ color: 'var(--primary-700)' }}
        />
      </Col>
      <Col xs={12} sm={8} md={6} lg={4}>
        <Statistic
          title="Documents"
          value={stats.totalDocuments}
          prefix={<FileTextOutlined style={{ color: 'var(--accent-blue)' }} />}
          valueStyle={{ color: 'var(--accent-blue)' }}
        />
      </Col>
      <Col xs={12} sm={8} md={6} lg={4}>
        <Statistic
          title="Chunks"
          value={stats.totalChunks}
          prefix={<BlockOutlined style={{ color: 'var(--accent-purple)' }} />}
          valueStyle={{ color: 'var(--accent-purple)' }}
        />
      </Col>
      <Col xs={12} sm={8} md={6} lg={4}>
        <Statistic
          title="Watched Folders"
          value={stats.watchedFolders}
          prefix={<EyeOutlined style={{ color: 'var(--accent-green)' }} />}
          valueStyle={{ color: 'var(--accent-green)' }}
        />
      </Col>
      <Col xs={12} sm={8} md={6} lg={4}>
        <Statistic
          title="Watched Websites"
          value={stats.watchedWebsites}
          prefix={<GlobalOutlined style={{ color: 'var(--accent-orange)' }} />}
          valueStyle={{ color: 'var(--accent-orange)' }}
        />
      </Col>
    </Row>
  );
};

// Composant pour afficher l'état d'un service avec bouton aligné
const ServiceStatus = ({ name, status, icon, details, logs, onShowDetails, version }) => {
  let statusIcon;
  let color;
  let statusText;
  
  switch (status) {
    case 'ok':
      statusIcon = <CheckCircleOutlined />;
      color = "var(--accent-green)";
      statusText = "Available";
      break;
    case 'warning':
      statusIcon = <ExclamationCircleOutlined />;
      color = "var(--accent-yellow)";
      statusText = "Warning";
      break;
    case 'error':
      statusIcon = <CloseCircleOutlined />;
      color = "var(--accent-red)";
      statusText = "Not available";
      break;
    default:
      statusIcon = <SyncOutlined spin />;
      color = "var(--neutral-500)";
      statusText = "Checking...";
  }
  
  return (
    <Card 
      className="service-card" 
      bodyStyle={{ 
        height: '100%', 
        display: 'flex', 
        flexDirection: 'column',
        padding: '24px'
      }}
      hoverable
    >
      <div style={{ flex: 1 }}>
        <Statistic 
          title={name}
          value={statusText}
          valueStyle={{ color }}
          prefix={statusIcon}
        />
        <div style={{ marginTop: 16, fontSize: 24, color: 'var(--primary-700)' }}>
          {icon}
        </div>
        {details && (
          <Paragraph type="secondary" style={{ marginTop: 8, fontSize: '0.85em' }}>
            {details}
          </Paragraph>
        )}
        {version && (
          <Paragraph type="secondary" style={{ marginTop: 4, fontSize: '0.75em', color: 'var(--primary-600)' }}>
            Version: {version}
          </Paragraph>
        )}
      </div>
      
      <div style={{ marginTop: 'auto', paddingTop: '16px' }}>
        <Button 
          type="primary" 
          block
          onClick={(e) => {
            e.stopPropagation();
            onShowDetails();
          }}
        >
          View details
        </Button>
      </div>
    </Card>
  );
};

const Home = () => {
  const [initialLoading, setInitialLoading] = useState(true);
  const [loading, setLoading] = useState(false);
  const [rlamaStatus, setRlamaStatus] = useState('loading');
  const [ollamaStatus, setOllamaStatus] = useState('loading');
  const [modelsStatus, setModelsStatus] = useState('loading');
  const [embeddingsStatus, setEmbeddingsStatus] = useState('loading');
  const [refreshing, setRefreshing] = useState(false);
  const refreshingRef = useRef(false);
  const [error, setError] = useState(null);
  const [serviceDetails, setServiceDetails] = useState({
    rlama: "Initializing system check...",
    ollama: "Initializing system check...",
    models: "Initializing system check...",
    embeddings: "Initializing system check..."
  });
  const [llmModelList, setLlmModelList] = useState([]);
  const [embeddingModelList, setEmbeddingModelList] = useState([]);
  const [logs, setLogs] = useState({
    rlama: null,
    ollama: null,
    models: null,
    embeddings: null
  });
  const [detailsModalVisible, setDetailsModalVisible] = useState(false);
  const [currentServiceDetails, setCurrentServiceDetails] = useState(null);
  const [rlamaCli, setRlamaCli] = useState({ status: 'loading', details: 'Initializing...' });
  const [ollamaCli, setOllamaCli] = useState({ status: 'loading', details: 'Initializing...' });
  const navigate = useNavigate();

  // Fonction utilitaire pour extraire la version courte
  const extractVersion = (stdout) => {
    if (!stdout) return null;
    // Pour RLAMA: extraire juste "RLAMA version X.X.X"
    const rlamaMatch = stdout.match(/RLAMA version (\S+)/i);
    if (rlamaMatch) return `v${rlamaMatch[1]}`;
    
    // Pour Ollama: extraire juste "ollama version X.X.X"
    const ollamaMatch = stdout.match(/ollama version (\S+)/i);
    if (ollamaMatch) return `v${ollamaMatch[1]}`;
    
    // Si pas de match, prendre la première ligne non vide et la nettoyer
    const lines = stdout.split('\n').filter(line => line.trim());
    if (lines.length > 0) {
      const firstLine = lines[0].trim();
      // Éviter les messages comme "Initialized services:"
      if (firstLine.toLowerCase().includes('initialized') || 
          firstLine.toLowerCase().includes('checking') ||
          firstLine.length > 50) {
        return 'Available';
      }
      return firstLine;
    }
    return 'Available';
  };

  const checkSystemHealth = async (isInitial = false, retryCount = 0) => {
    const maxRetries = 2;
    
    if (refreshingRef.current && retryCount === 0) {
      console.log("Health check already in progress. Skipping.");
      return;
    }
    
    console.log(`Starting system health check (isInitial: ${isInitial}, retry: ${retryCount})`);
    
    refreshingRef.current = true;
    if (isInitial) {
      setInitialLoading(true);
    } else {
      setRefreshing(true);
      setLoading(true);
    }
    setError(null);

    // Reset all statuses to loading
    setRlamaStatus('loading');
    setOllamaStatus('loading');
    setModelsStatus('loading');
    setEmbeddingsStatus('loading');
    
    try {
      // Petit délai pour s'assurer que le composant est bien monté
      if (isInitial) {
        await new Promise(resolve => setTimeout(resolve, 100));
      }
      
      // Check RLAMA version
      try {
        console.log("🔍 Checking RLAMA CLI...");
        const rlamaVersionOutput = await api.get(`/exec?command=${encodeURIComponent('rlama --version')}`);
        console.log("✅ RLAMA version API response:", {
          status: rlamaVersionOutput.status,
          stdout: rlamaVersionOutput.data?.stdout,
          stderr: rlamaVersionOutput.data?.stderr,
          returncode: rlamaVersionOutput.data?.returncode
        });
        
        if (rlamaVersionOutput && rlamaVersionOutput.data && rlamaVersionOutput.data.stdout && !rlamaVersionOutput.data.stderr) {
          const shortVersion = extractVersion(rlamaVersionOutput.data.stdout);
          console.log("✅ RLAMA CLI found, version:", shortVersion);
          setRlamaStatus('ok');
          setServiceDetails(prev => ({
            ...prev,
            rlama: `RLAMA CLI installed and functional`
          }));
          setLogs(prev => ({ 
            ...prev, 
            rlama: `RLAMA version check successful: ${shortVersion}`
          }));
          setRlamaCli({
            status: 'ok',
            version: shortVersion,
            details: 'RLAMA CLI available',
            output: rlamaVersionOutput.data.stdout
          });
        } else {
          // Fallback: essayer avec l'ancienne commande -v
          console.log("🔄 Trying RLAMA fallback command...");
          try {
            const fallbackOutput = await api.get(`/exec?command=${encodeURIComponent('rlama -v')}`);
            if (fallbackOutput?.data?.stdout) {
              const shortVersion = extractVersion(fallbackOutput.data.stdout);
              console.log("✅ RLAMA CLI found via fallback, version:", shortVersion);
              setRlamaStatus('ok');
              setServiceDetails(prev => ({
                ...prev,
                rlama: `RLAMA CLI installed and functional`
              }));
              setLogs(prev => ({ 
                ...prev, 
                rlama: `RLAMA version check successful: ${shortVersion}`
              }));
              setRlamaCli({
                status: 'ok',
                version: shortVersion,
                details: 'RLAMA CLI available',
                output: fallbackOutput.data.stdout
              });
            } else {
              throw new Error('RLAMA not installed or version not detected');
            }
          } catch (fallbackError) {
            throw new Error('RLAMA not installed or version not detected');
          }
        }
      } catch (err) {
        console.error("❌ RLAMA check error:", err);
        setRlamaStatus('error');
        setServiceDetails(prev => ({ 
          ...prev, 
          rlama: `RLAMA not installed or inaccessible` 
        }));
        setLogs(prev => ({ 
          ...prev, 
          rlama: `Error: ${err.message}`
        }));
        setRlamaCli({
          status: 'error',
          version: null,
          details: `Error: ${err.message}`
        });
      }

      // Check Ollama version
      try {
        console.log("🔍 Checking Ollama CLI...");
        const ollamaVersionOutput = await api.get(`/exec?command=${encodeURIComponent('ollama --version')}`);
        console.log("✅ Ollama version API response:", {
          status: ollamaVersionOutput.status,
          stdout: ollamaVersionOutput.data?.stdout,
          stderr: ollamaVersionOutput.data?.stderr,
          returncode: ollamaVersionOutput.data?.returncode
        });
        
        if (ollamaVersionOutput && ollamaVersionOutput.data && ollamaVersionOutput.data.stdout && !ollamaVersionOutput.data.stderr) {
          const shortVersion = extractVersion(ollamaVersionOutput.data.stdout);
          console.log("✅ Ollama CLI found, version:", shortVersion);
          setOllamaStatus('ok');
          setServiceDetails(prev => ({
            ...prev,
            ollama: `Ollama CLI installed and functional`
          }));
          setLogs(prev => ({ 
            ...prev, 
            ollama: `Ollama version check successful: ${shortVersion}`
          }));
          setOllamaCli({
            status: 'ok',
            version: shortVersion,
            details: 'Ollama CLI available',
            output: ollamaVersionOutput.data.stdout
          });
          
          // Si Ollama est disponible, on considère que les modèles sont disponibles
          console.log("✅ Setting models and embeddings as available");
          setModelsStatus('ok');
          setEmbeddingsStatus('ok');
          setServiceDetails(prev => ({
            ...prev,
            models: 'Ollama detected - LLM models available',
            embeddings: 'Ollama detected - embedding models available'
          }));
        } else {
          // Fallback: essayer avec l'ancienne commande -v
          console.log("🔄 Trying Ollama fallback command...");
          try {
            const fallbackOutput = await api.get(`/exec?command=${encodeURIComponent('ollama -v')}`);
            if (fallbackOutput?.data?.stdout) {
              const shortVersion = extractVersion(fallbackOutput.data.stdout);
              console.log("✅ Ollama CLI found via fallback, version:", shortVersion);
              setOllamaStatus('ok');
              setServiceDetails(prev => ({
                ...prev,
                ollama: `Ollama CLI installed and functional`
              }));
              setLogs(prev => ({ 
                ...prev, 
                ollama: `Ollama version check successful: ${shortVersion}`
              }));
              setOllamaCli({
                status: 'ok',
                version: shortVersion,
                details: 'Ollama CLI available',
                output: fallbackOutput.data.stdout
              });
              
              // Si Ollama est disponible, on considère que les modèles sont disponibles
              console.log("✅ Setting models and embeddings as available (fallback)");
              setModelsStatus('ok');
              setEmbeddingsStatus('ok');
              setServiceDetails(prev => ({
                ...prev,
                models: 'Ollama detected - LLM models available',
                embeddings: 'Ollama detected - embedding models available'
              }));
            } else {
              throw new Error('Ollama not installed or version not detected');
            }
          } catch (fallbackError) {
            throw new Error('Ollama not installed or version not detected');
          }
        }
      } catch (err) {
        console.error("❌ Ollama check error:", err);
        setOllamaStatus('error');
        setServiceDetails(prev => ({ 
          ...prev, 
          ollama: `Ollama not installed or inaccessible` 
        }));
        setLogs(prev => ({ 
          ...prev, 
          ollama: `Error: ${err.message}`
        }));
        setOllamaCli({
          status: 'error',
          version: null,
          details: `Error: ${err.message}`
        });
        
        // Si Ollama n'est pas disponible, les modèles ne le sont pas non plus
        console.log("❌ Setting models and embeddings as error due to Ollama failure");
        setModelsStatus('error');
        setEmbeddingsStatus('error');
        setServiceDetails(prev => ({
          ...prev,
          models: 'Ollama not available - LLM models inaccessible',
          embeddings: 'Ollama not available - embedding models inaccessible'
        }));
      }
    } catch (globalError) {
      console.error("System check global error:", globalError);
      
      // Si c'est un chargement initial et qu'on peut encore retry
      if (isInitial && retryCount < maxRetries) {
        console.log(`Retrying system check in 2 seconds (attempt ${retryCount + 1}/${maxRetries})`);
        setTimeout(() => {
          refreshingRef.current = false; // Reset pour permettre le retry
          checkSystemHealth(true, retryCount + 1);
        }, 2000);
        return; // Ne pas exécuter le finally pour ce retry
      }
      
      setError(`System check error: ${globalError.message}`);
      setRlamaStatus('error');
      setOllamaStatus('error');
      setModelsStatus('error');
      setEmbeddingsStatus('error');
    } finally {
      // Only reset states if we're not going to retry
      if (!isInitial || retryCount >= maxRetries) {
        setLoading(false);
        setRefreshing(false);
        setInitialLoading(false);
        refreshingRef.current = false;
        console.log("System health check completed");
      }
    }
  };

  useEffect(() => {
    console.log("🚀 RLAMA Dashboard mounting - starting initial system check");
    
    // Lancer le check du système au démarrage avec le flag initial
    checkSystemHealth(true);
    
    // Timeout de sécurité au cas où le premier check reste bloqué
    const safetyTimeout = setTimeout(() => {
      if (refreshingRef.current) {
        console.warn("⚠️ System check seems stuck, forcing reset");
        refreshingRef.current = false;
        setInitialLoading(false);
        setLoading(false);
        setRefreshing(false);
      }
    }, 30000); // 30 secondes de timeout
    
    return () => {
      clearTimeout(safetyTimeout);
      console.log("🏁 RLAMA Dashboard unmounting");
    };
  }, []);

  const handleRefresh = () => {
    checkSystemHealth();
  };

  const manuallyConfirmOllama = () => {
    setOllamaStatus('ok');
    setServiceDetails(prev => ({
      ...prev,
      ollama: 'Ollama manually confirmed as functional'
    }));
    setLogs(prev => ({
      ...prev,
      ollama: `Ollama manually confirmed as functional by user.`
    }));
    if (detailsModalVisible) {
        setDetailsModalVisible(false);
    }
  };

  const showDetailsModal = (serviceName) => {
    const statusMap = { ok: "Available", warning: "Warning", error: "Not available", loading: "Checking..." };
    let title = serviceName;
    let status = 'loading';
    let detailsText = '';
    let serviceLogs = '';
    let specificModelList = [];
    let isLlm = false;
    let isEmb = false;

    switch (serviceName) {
      case 'RLAMA API':
        status = rlamaStatus;
        detailsText = serviceDetails.rlama;
        serviceLogs = logs.rlama;
        break;
      case 'Ollama':
        status = ollamaStatus;
        detailsText = serviceDetails.ollama;
        serviceLogs = logs.ollama;
        break;
      case 'LLM Models':
        status = modelsStatus;
        detailsText = serviceDetails.models;
        serviceLogs = logs.models;
        specificModelList = llmModelList;
        isLlm = true;
        break;
      case "Embedding Models":
        status = embeddingsStatus;
        detailsText = serviceDetails.embeddings;
        serviceLogs = logs.embeddings;
        specificModelList = embeddingModelList;
        isEmb = true;
        break;
      default:
        return;
    }

    setCurrentServiceDetails({
      name: title,
      status: status,
      statusText: statusMap[status] || status,
      details: detailsText,
      logs: serviceLogs,
      modelList: specificModelList,
      isLlmType: isLlm,
      isEmbType: isEmb,
    });
    setDetailsModalVisible(true);
  };

  // Écran de chargement initial
  if (initialLoading) {
    return (
      <div style={{ 
        display: 'flex', 
        flexDirection: 'column', 
        alignItems: 'center', 
        justifyContent: 'center', 
        minHeight: '70vh',
        textAlign: 'center'
      }}>
        <div style={{ marginBottom: '32px' }}>
          <Title level={1} style={{ color: 'var(--primary-800)', margin: 0 }}>
            Initializing RLAMA
          </Title>
          <Text type="secondary">
            Checking system components...
            {rlamaStatus === 'error' && ollamaStatus === 'error' ? ' (retrying...)' : ''}
          </Text>
        </div>
        
        <Spin size="large" style={{ marginBottom: '24px' }} />
        
        <div style={{ width: '400px', textAlign: 'left' }}>
          <div style={{ marginBottom: '12px' }}>
            <Text type="secondary">
              {rlamaStatus === 'loading' ? '🔄' : rlamaStatus === 'ok' ? '✅' : '❌'} 
              {' '}Checking RLAMA CLI...
            </Text>
          </div>
          <div style={{ marginBottom: '12px' }}>
            <Text type="secondary">
              {ollamaStatus === 'loading' ? '🔄' : ollamaStatus === 'ok' ? '✅' : '❌'} 
              {' '}Checking Ollama...
            </Text>
          </div>
          <div style={{ marginBottom: '12px' }}>
            <Text type="secondary">
              {modelsStatus === 'loading' ? '🔄' : modelsStatus === 'ok' ? '✅' : '❌'} 
              {' '}Verifying models...
            </Text>
          </div>
          <div style={{ marginBottom: '12px' }}>
            <Text type="secondary">
              {embeddingsStatus === 'loading' ? '🔄' : embeddingsStatus === 'ok' ? '✅' : '❌'} 
              {' '}Checking embeddings...
            </Text>
          </div>
        </div>
        
        <div style={{ marginTop: '24px', textAlign: 'center' }}>
          <Text type="secondary" style={{ fontSize: '12px', display: 'block', marginBottom: '16px' }}>
            This usually takes a few seconds...
          </Text>
          <Button 
            type="text" 
            size="small" 
            onClick={() => {
              console.log("User skipped initial loading");
              setInitialLoading(false);
              refreshingRef.current = false;
            }}
          >
            Skip and continue →
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="text-center mb-8">
        <Title level={1} style={{ color: 'var(--primary-800)', margin: 0 }}>RLAMA Dashboard</Title>
        <Text type="secondary">Retrieval-Augmented Generation & Agents Platform</Text>
      </div>
      
      {error && (
        <Alert
          message="System check error"
          description={error}
          type="error"
          showIcon
          className="mb-6"
          closable
          onClose={() => setError(null)}
        />
      )}
      
      <Card 
        title="System Status" 
        className="shadow-md rounded-lg mb-8"
        extra={
          <Button 
            icon={<SyncOutlined />} 
            loading={refreshing || (loading && !refreshing)}
            onClick={handleRefresh}
            disabled={refreshing}
          >
            Actualize
          </Button>
        }
      >
        <Spin spinning={loading && !refreshing} tip="Checking system status...">
          <Row gutter={[16, 16]}>
            <Col xs={24} sm={12} md={6} style={{ display: 'flex' }}>
              <ServiceStatus 
                name="RLAMA API" 
                status={rlamaStatus} 
                icon={<ApiOutlined />} 
                details={serviceDetails.rlama}
                version={rlamaCli.version}
                onShowDetails={() => showDetailsModal('RLAMA API')}
              />
            </Col>
            <Col xs={24} sm={12} md={6} style={{ display: 'flex' }}>
              <ServiceStatus 
                name="Ollama" 
                status={ollamaStatus} 
                icon={<RobotOutlined />} 
                details={serviceDetails.ollama}
                version={ollamaCli.version}
                onShowDetails={() => showDetailsModal('Ollama')}
              />
            </Col>
            <Col xs={24} sm={12} md={6} style={{ display: 'flex' }}>
              <ServiceStatus 
                name="LLM Models" 
                status={modelsStatus} 
                icon={<RobotOutlined />} 
                details={serviceDetails.models}
                onShowDetails={() => showDetailsModal('LLM Models')}
              />
            </Col>
            <Col xs={24} sm={12} md={6} style={{ display: 'flex' }}>
              <ServiceStatus 
                name="Embedding Models" 
                status={embeddingsStatus} 
                icon={<DatabaseOutlined />} 
                details={serviceDetails.embeddings}
                onShowDetails={() => showDetailsModal("Embedding Models")}
              />
            </Col>
          </Row>
        </Spin>
      </Card>
      
      <Divider />
      
      {/* Dashboard Statistics */}
      <Card 
        title="Dashboard Statistics" 
        className="shadow-md rounded-lg mb-8"
        extra={
          <Button 
            type="link" 
            onClick={() => navigate('/systems')}
            style={{ padding: 0 }}
          >
            View all systems →
          </Button>
        }
      >
        <DashboardStats />
      </Card>
      
      {/* Quick Actions */}
      <Row gutter={[16, 16]} className="mb-8">
        <Col xs={24} sm={8}>
          <Card className="text-center h-full" hoverable onClick={() => navigate('/systems')}>
            <DatabaseOutlined style={{ fontSize: 32, color: 'var(--primary-600)', marginBottom: 8 }} />
            <Title level={4} style={{ margin: '8px 0 4px 0' }}>Manage RAGs</Title>
            <Text type="secondary">View and manage your RAG systems</Text>
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card className="text-center h-full" hoverable onClick={() => navigate('/create')}>
            <PlusOutlined style={{ fontSize: 32, color: 'var(--accent-green)', marginBottom: 8 }} />
            <Title level={4} style={{ margin: '8px 0 4px 0' }}>Create RAG</Title>
            <Text type="secondary">Set up a new RAG system</Text>
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card className="text-center h-full" hoverable onClick={() => window.open('https://github.com/DonTizi/rlama', '_blank')}>
            <RobotOutlined style={{ fontSize: 32, color: 'var(--primary-800)', marginBottom: 8 }} />
            <Title level={4} style={{ margin: '8px 0 4px 0' }}>Documentation</Title>
            <Text type="secondary">Learn more about RLAMA</Text>
          </Card>
        </Col>
      </Row>

      <Modal
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <div style={{ 
              color: currentServiceDetails?.status === 'ok' ? 'var(--accent-green)' : 
                     currentServiceDetails?.status === 'warning' ? 'var(--accent-yellow)' : 'var(--accent-red)',
              fontSize: '20px'
            }}>
              {currentServiceDetails?.status === 'ok' ? <CheckCircleOutlined /> : 
               currentServiceDetails?.status === 'warning' ? <ExclamationCircleOutlined /> : <CloseCircleOutlined />}
            </div>
            <span>{currentServiceDetails?.name} Details</span>
          </div>
        }
        open={detailsModalVisible}
        onCancel={() => setDetailsModalVisible(false)}
        width={700}
        footer={[
          <Button key="close" type="primary" onClick={() => setDetailsModalVisible(false)}>
            Close
          </Button>
        ]}
      >
        {currentServiceDetails && (
          <div style={{ padding: '8px 0' }}>
            {/* Status Section */}
            <div style={{ 
              background: currentServiceDetails.status === 'ok' ? '#f6ffed' : 
                         currentServiceDetails.status === 'warning' ? '#fffbe6' : '#fff2f0',
              border: `1px solid ${currentServiceDetails.status === 'ok' ? '#b7eb8f' : 
                                    currentServiceDetails.status === 'warning' ? '#ffe58f' : '#ffccc7'}`,
              borderRadius: '8px',
              padding: '16px',
              marginBottom: '20px'
            }}>
              <Title level={5} style={{ margin: '0 0 8px 0' }}>
                Status: <Text style={{ 
                  color: currentServiceDetails.status === 'ok' ? 'var(--accent-green)' : 
                         currentServiceDetails.status === 'warning' ? 'var(--accent-yellow)' : 'var(--accent-red)'
                }}>
                  {currentServiceDetails.statusText}
                </Text>
              </Title>
              
              {/* Version info */}
              {currentServiceDetails.name === 'RLAMA API' && rlamaCli.version && (
                <Text style={{ display: 'block', marginBottom: '8px' }}>
                  <Text strong>Version:</Text> {rlamaCli.version}
                </Text>
              )}
              
              {currentServiceDetails.name === 'Ollama' && ollamaCli.version && (
                <Text style={{ display: 'block', marginBottom: '8px' }}>
                  <Text strong>Version:</Text> {ollamaCli.version}
                </Text>
              )}
              
              <Text type="secondary">{currentServiceDetails.details}</Text>
            </div>

            {/* Models Section */}
            {(currentServiceDetails.isLlmType && currentServiceDetails.modelList?.length > 0) && (
              <div style={{ marginBottom: '20px' }}>
                <Title level={5}>Available LLM Models</Title>
                <List
                  size="small"
                  bordered
                  dataSource={currentServiceDetails.modelList}
                  renderItem={item => (
                    <List.Item style={{ padding: '8px 16px' }}>
                      {typeof item === 'object' ? item.name : item}
                    </List.Item>
                  )}
                  style={{ maxHeight: '150px', overflowY: 'auto' }}
                />
              </div>
            )}
            
            {(currentServiceDetails.isEmbType && currentServiceDetails.modelList?.length > 0) && (
              <div style={{ marginBottom: '20px' }}>
                <Title level={5}>Available Embedding Models</Title>
                <List
                  size="small"
                  bordered
                  dataSource={currentServiceDetails.modelList}
                  renderItem={item => (
                    <List.Item style={{ padding: '8px 16px' }}>
                      {typeof item === 'object' ? item.name : item}
                    </List.Item>
                  )}
                  style={{ maxHeight: '150px', overflowY: 'auto' }}
                />
              </div>
            )}

            {/* Alerts */}
            {currentServiceDetails.isEmbType && currentServiceDetails.modelList?.length === 0 && (currentServiceDetails.status === 'warning' || currentServiceDetails.status === 'error') && (
              <Alert 
                message="No embedding models found"
                description={
                  <>
                    No local embedding models detected. You can download models from{' '}
                    <a href="https://ollama.com/models" target="_blank" rel="noopener noreferrer">
                      Ollama Models
                    </a>.
                  </>
                }
                type="info"
                showIcon
                style={{ marginBottom: '20px' }}
              />
            )}
            
            {currentServiceDetails.name === 'Ollama' && (currentServiceDetails.status === 'warning' || currentServiceDetails.status === 'error') && (
              <Alert
                message="Ollama Connection Issue"
                description="Ollama CLI detection failed. Make sure Ollama is installed and running."
                type="warning"
                showIcon
                style={{ marginBottom: '20px' }}
                action={
                  <Space direction="vertical" style={{width: "100%"}}>
                    <Button type="default" onClick={manuallyConfirmOllama} block>
                      Mark as functional
                    </Button>
                    <Button type="primary" onClick={() => window.open('https://ollama.com/', '_blank')} block>
                      Install Ollama
                    </Button>
                  </Space>
                }
              />
            )}

            {/* Technical Details (collapsible) */}
            {(currentServiceDetails.logs || (currentServiceDetails.name === 'RLAMA API' && rlamaCli.output) || (currentServiceDetails.name === 'Ollama' && ollamaCli.output)) && (
              <div style={{ marginBottom: '20px' }}>
                <Title level={5}>Technical Information</Title>
                <div style={{ 
                  background: '#fafafa', 
                  border: '1px solid #d9d9d9',
                  borderRadius: '6px',
                  padding: '12px',
                  fontSize: '12px',
                  fontFamily: 'monospace',
                  maxHeight: '120px',
                  overflowY: 'auto'
                }}>
                  {currentServiceDetails.name === 'RLAMA API' && rlamaCli.output && (
                    <div><strong>RLAMA Output:</strong><br/>{rlamaCli.output}</div>
                  )}
                  {currentServiceDetails.name === 'Ollama' && ollamaCli.output && (
                    <div><strong>Ollama Output:</strong><br/>{ollamaCli.output}</div>
                  )}
                  {currentServiceDetails.logs && (
                    <div><strong>Logs:</strong><br/>{currentServiceDetails.logs}</div>
                  )}
                </div>
              </div>
            )}

            {/* Connection Info */}
            <div style={{ 
              background: '#f9f9f9', 
              padding: '12px', 
              borderRadius: '6px',
              fontSize: '13px' 
            }}>
              <Text strong>Connection Info:</Text><br/>
              <Text type="secondary">Backend API: http://localhost:5001</Text><br/>
              <Text type="secondary">Last check: {new Date().toLocaleString()}</Text>
            </div>
          </div>
        )}
      </Modal>

      <style jsx>{`
        .service-card {
          flex: 1;
          display: flex;
          flex-direction: column;
        }
        .service-card:hover {
          transform: translateY(-5px);
          transition: transform 0.3s ease;
        }
      `}</style>
    </div>
  );
};

export default Home; 