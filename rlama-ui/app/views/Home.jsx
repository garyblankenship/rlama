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
  PlusOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { healthService, ragService } from '../services/api';
import api from '../services/api';

const { Title, Text, Paragraph } = Typography;

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
  const [loading, setLoading] = useState(true);
  const [rlamaStatus, setRlamaStatus] = useState('loading');
  const [ollamaStatus, setOllamaStatus] = useState('loading');
  const [modelsStatus, setModelsStatus] = useState('loading');
  const [embeddingsStatus, setEmbeddingsStatus] = useState('loading');
  const [refreshing, setRefreshing] = useState(false);
  const refreshingRef = useRef(false);
  const [error, setError] = useState(null);
  const [serviceDetails, setServiceDetails] = useState({
    rlama: "Checking...",
    ollama: "Checking...",
    models: "Checking...",
    embeddings: "Checking..."
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
  const [rlamaCli, setRlamaCli] = useState({ status: 'loading', details: 'Checking...' });
  const [ollamaCli, setOllamaCli] = useState({ status: 'loading', details: 'Checking...' });
  const navigate = useNavigate();

  const checkSystemHealth = async () => {
    if (refreshingRef.current) {
      console.log("Health check already in progress. Skipping.");
      return;
    }
    
    refreshingRef.current = true;
    setRefreshing(true);
    setLoading(true);
    setError(null);

    // Reset all statuses to loading
    setRlamaStatus('loading');
    setOllamaStatus('loading');
    setModelsStatus('loading');
    setEmbeddingsStatus('loading');
    
    try {
      // Check RLAMA version
      try {
        const rlamaVersionOutput = await api.get(`/exec?command=${encodeURIComponent('rlama -v')}`);
        console.log("RLAMA version output:", rlamaVersionOutput);
        
        if (rlamaVersionOutput && rlamaVersionOutput.data && rlamaVersionOutput.data.stdout && !rlamaVersionOutput.data.stderr) {
          setRlamaStatus('ok');
          setServiceDetails(prev => ({
            ...prev,
            rlama: `RLAMA installed: ${rlamaVersionOutput.data.stdout.trim()}`
          }));
          setLogs(prev => ({ 
            ...prev, 
            rlama: `RLAMA version check:\n${rlamaVersionOutput.data.stdout}`
          }));
          setRlamaCli({
            status: 'ok',
            version: rlamaVersionOutput.data.stdout.trim(),
            details: 'RLAMA CLI available',
            output: rlamaVersionOutput.data.stdout
          });
        } else {
          throw new Error('RLAMA not installed or version not detected');
        }
      } catch (err) {
        console.error("RLAMA check error:", err);
        setRlamaStatus('error');
        setServiceDetails(prev => ({ 
          ...prev, 
          rlama: `RLAMA not installed or inaccessible: ${err.message}` 
        }));
        setLogs(prev => ({ 
          ...prev, 
          rlama: `Error: ${err.message}`
        }));
      }

      // Check Ollama version
      try {
        const ollamaVersionOutput = await api.get(`/exec?command=${encodeURIComponent('ollama -v')}`);
        console.log("Ollama version output:", ollamaVersionOutput);
        
        if (ollamaVersionOutput && ollamaVersionOutput.data && ollamaVersionOutput.data.stdout && !ollamaVersionOutput.data.stderr) {
          setOllamaStatus('ok');
          setServiceDetails(prev => ({
            ...prev,
            ollama: `Ollama installed: ${ollamaVersionOutput.data.stdout.trim()}`
          }));
          setLogs(prev => ({ 
            ...prev, 
            ollama: `Ollama version check:\n${ollamaVersionOutput.data.stdout}`
          }));
          setOllamaCli({
            status: 'ok',
            version: ollamaVersionOutput.data.stdout.trim(),
            details: 'Ollama CLI available',
            output: ollamaVersionOutput.data.stdout
          });
          
          // Si Ollama est disponible, on considère que les modèles sont disponibles
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
      } catch (err) {
        console.error("Ollama check error:", err);
        setOllamaStatus('error');
        setServiceDetails(prev => ({ 
          ...prev, 
          ollama: `Ollama not installed or inaccessible: ${err.message}` 
        }));
        setLogs(prev => ({ 
          ...prev, 
          ollama: `Error: ${err.message}`
        }));
        
        // Si Ollama n'est pas disponible, les modèles ne le sont pas non plus
        setModelsStatus('error');
        setEmbeddingsStatus('error');
        setServiceDetails(prev => ({
          ...prev,
          models: 'Ollama not available - LLM models inaccessible',
          embeddings: 'Ollama not available - embedding models inaccessible'
        }));
      }
    } catch (globalError) {
      setError(`System check error: ${globalError.message}`);
      console.error("System check global error:", globalError);
      setRlamaStatus('error');
      setOllamaStatus('error');
      setModelsStatus('error');
      setEmbeddingsStatus('error');
    } finally {
      setLoading(false);
      setRefreshing(false);
      refreshingRef.current = false;
    }
  };

  useEffect(() => {
    checkSystemHealth();
    
    const intervalId = setInterval(() => {
      checkSystemHealth();
    }, 30000);
    
    return () => {
      clearInterval(intervalId);
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

  return (
    <div>
      <div className="text-center mb-8">
        <Title level={1} style={{ color: 'var(--primary-800)', margin: 0 }}>RLAMA Dashboard</Title>
        <Text type="secondary">Retrieval-Augmented Generation Platform</Text>
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
      
      <Row gutter={[16, 16]} className="mb-8">
        <Col xs={24} sm={12} md={8} style={{ display: 'flex' }}>
          <Card 
            title="RAG Systems" 
            className="shadow-md rounded-lg text-center w-full"
            bodyStyle={{ 
              display: 'flex', 
              flexDirection: 'column', 
              height: '100%', 
              padding: '24px',
              gap: '16px',
              overflow: 'hidden'
            }}
            hoverable
            onClick={() => navigate('/systems')}
          >
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '16px' }}>
              <DatabaseOutlined style={{ fontSize: 48, color: 'var(--primary-600)' }} />
              <Paragraph className="mt-2 mb-0">Manage your existing RAG systems</Paragraph>
            </div>
            <div style={{ marginTop: 'auto', paddingTop: '16px', width: '100%' }}>
              <Button type="primary" size="large" block>
                View my systems
              </Button>
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} style={{ display: 'flex' }}>
          <Card 
            title="Create RAG" 
            className="shadow-md rounded-lg text-center w-full"
            bodyStyle={{ 
              display: 'flex', 
              flexDirection: 'column', 
              height: '100%', 
              padding: '24px',
              gap: '16px',
              overflow: 'hidden'
            }}
            hoverable
            onClick={() => navigate('/create')}
          >
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '16px' }}>
              <PlusOutlined style={{ fontSize: 48, color: 'var(--accent-green)' }} />
              <Paragraph className="mt-2 mb-0">Configure a new RAG system</Paragraph>
            </div>
            <div style={{ marginTop: 'auto', paddingTop: '16px', width: '100%' }}>
              <Button type="primary" icon={<PlusOutlined />} size="large" block>
                Create RAG
              </Button>
            </div>
          </Card>
        </Col>
        <Col xs={24} md={8} style={{ display: 'flex' }}>
          <Card 
            title="Documentation" 
            className="shadow-md rounded-lg text-center w-full"
            bodyStyle={{ 
              display: 'flex', 
              flexDirection: 'column', 
              height: '100%', 
              padding: '24px',
              gap: '16px',
              overflow: 'hidden'
            }}
            hoverable
            onClick={() => window.open('https://github.com/DonTizi/rlama', '_blank')}
          >
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '16px' }}>
              <RobotOutlined style={{ fontSize: 48, color: 'var(--primary-800)' }} />
              <Paragraph className="mt-2 mb-0">View documentation</Paragraph>
            </div>
            <div style={{ marginTop: 'auto', paddingTop: '16px', width: '100%' }}>
              <Button type="default" size="large" block>
                Access documentation
              </Button>
            </div>
          </Card>
        </Col>
      </Row>

      <Modal
        title={`Details of ${currentServiceDetails?.name}`}
        open={detailsModalVisible}
        onCancel={() => setDetailsModalVisible(false)}
        width={800}
        footer={[
          <Button 
            key="debug" 
            type="primary" 
            danger
            onClick={() => {
              // Open console and log all available debug info
              console.log("==== DEBUG INFO ====");
              console.log("Service Details:", serviceDetails);
              console.log("Service Logs:", logs);
              console.log("RLAMA Status:", rlamaStatus);
              console.log("Ollama Status:", ollamaStatus);
              console.log("Current Service:", currentServiceDetails);
            }}
          >
            Log Debug Info
          </Button>,
          <Button key="close" onClick={() => setDetailsModalVisible(false)}>
            Close
          </Button>
        ]}
      >
        {currentServiceDetails && (
          <>
            <Title level={5}>Status: <Text style={{ color: currentServiceDetails.status === 'ok' ? 'var(--accent-green)' : currentServiceDetails.status === 'warning' ? 'var(--accent-yellow)' : 'var(--accent-red)'}}>{currentServiceDetails.statusText}</Text></Title>
            
            {currentServiceDetails.name === 'RLAMA API' && rlamaCli.version && (
              <Paragraph>
                <Text strong>RLAMA Version:</Text> {rlamaCli.version}
                {rlamaCli.output && (
                  <pre style={{ 
                    whiteSpace: 'pre-wrap', 
                    background: 'var(--neutral-100)', 
                    padding: '8px', 
                    borderRadius: '4px',
                    fontSize: '0.85em',
                    marginTop: '8px'
                  }}>
                    {rlamaCli.output}
                  </pre>
                )}
              </Paragraph>
            )}
            
            {currentServiceDetails.name === 'Ollama' && ollamaCli.version && (
              <Paragraph>
                <Text strong>Ollama Version:</Text> {ollamaCli.version}
                {ollamaCli.output && (
                  <pre style={{ 
                    whiteSpace: 'pre-wrap', 
                    background: 'var(--neutral-100)', 
                    padding: '8px', 
                    borderRadius: '4px',
                    fontSize: '0.85em',
                    marginTop: '8px'
                  }}>
                    {ollamaCli.output}
                  </pre>
                )}
              </Paragraph>
            )}
            
            <Paragraph>{currentServiceDetails.details}</Paragraph>

            {(currentServiceDetails.isLlmType && currentServiceDetails.modelList?.length > 0) && (
              <>
                <Paragraph strong>Local LLM models detected :</Paragraph>
                <List
                  size="small"
                  bordered
                  dataSource={currentServiceDetails.modelList}
                  renderItem={item => <List.Item>{typeof item === 'object' ? item.name : item}</List.Item>}
                  style={{ marginBottom: '16px', maxHeight: '150px', overflowY: 'auto' }}
                />
              </>
            )}
             {(currentServiceDetails.isEmbType && currentServiceDetails.modelList?.length > 0) && (
              <>
                <Paragraph strong>Local embedding models detected :</Paragraph>
                <List
                  size="small"
                  bordered
                  dataSource={currentServiceDetails.modelList}
                  renderItem={item => <List.Item>{typeof item === 'object' ? item.name : item}</List.Item>}
                  style={{ marginBottom: '16px', maxHeight: '150px', overflowY: 'auto' }}
                />
              </>
            )}

            {currentServiceDetails.isEmbType && currentServiceDetails.modelList?.length === 0 && (currentServiceDetails.status === 'warning' || currentServiceDetails.status === 'error') && (
              <Alert 
                message="No local embedding model found or endpoint inaccessible"
                description={
                  <>
                    Check that your RLAMA backend exposes a functional `/embedding-models` endpoint.
                    You can explore and download models (including embeddings) from the official Ollama site.{' '}
                    <a href="https://ollama.com/models" target="_blank" rel="noopener noreferrer">
                      Visit Ollama Models
                    </a>.
                  </>
                }
                type="info"
                showIcon
                style={{ marginBottom: '16px' }}
              />
            )}
            
            {currentServiceDetails.name === 'Ollama' && (currentServiceDetails.status === 'warning' || currentServiceDetails.status === 'error') && (
                <Alert
                    message="Ollama detection issue"
                    description={
                        <>
                            Automatic Ollama detection (via `ollama list`) failed or returned a warning.
                            If Ollama is installed and `ollama serve` is running on your machine, try refreshing.
                            You can also check the Ollama documentation at{' '}
                            <a href="https://ollama.com/" target="_blank" rel="noopener noreferrer">
                                Ollama.com
                            </a> or manually confirm if you're sure it's functional.
                        </>
                    }
                    type="warning"
                    showIcon
                    style={{ marginBottom: '16px' }}
                    action={
                         <Space direction="vertical" style={{width: "100%"}}>
                             <Button type="default" onClick={manuallyConfirmOllama} block>
                                I've verified, Ollama is functional
                            </Button>
                             <Button type="primary" onClick={() => window.open('https://ollama.com/', '_blank')} block>
                                Visit Ollama.com for help/installation
                            </Button>
                         </Space>
                    }
                />
            )}

            {currentServiceDetails.logs && (
              <>
                <Divider>Logs / Technical Information</Divider>
                <pre style={{ 
                  whiteSpace: 'pre-wrap', 
                  wordBreak: 'break-all', 
                  background: 'var(--neutral-100)', 
                  padding: '10px', 
                  borderRadius: '4px',
                  maxHeight: '200px',
                  overflowY: 'auto',
                  color: 'var(--neutral-700)',
                  fontSize: '0.85em'
                }}>
                  {currentServiceDetails.logs}
                </pre>
              </>
            )}

            <Divider>Debug Information</Divider>
            <Paragraph>
              <Text strong>API URL:</Text> {window.location.origin}/api
            </Paragraph>
            <Paragraph>
              <Text strong>Last check:</Text> {new Date().toLocaleString()}
            </Paragraph>
            <Button 
              onClick={() => {
                const debug = {
                  apiUrl: window.location.origin + '/api',
                  timestamp: new Date().toISOString(),
                  statuses: {
                    rlamaApi: rlamaStatus,
                    ollama: ollamaStatus,
                    models: modelsStatus,
                    embeddings: embeddingsStatus
                  },
                  details: serviceDetails,
                  logs: logs
                };
                
                // Create a downloadable file with debug info
                const blob = new Blob([JSON.stringify(debug, null, 2)], {type: 'application/json'});
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'rlama-debug.json';
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
              }}
            >
              Download debug logs
            </Button>
          </>
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