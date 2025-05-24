import React, { useState, useEffect } from 'react';
import { Card, Table, Button, Space, Typography, message, Tooltip, Popconfirm, Tag, Badge, Row, Col, Statistic } from 'antd';
import { 
  PlusOutlined, 
  DeleteOutlined, 
  EyeOutlined, 
  SearchOutlined,
  ClockCircleOutlined,
  FileTextOutlined,
  RobotOutlined,
  CalendarOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ExclamationCircleOutlined,
  SyncOutlined,
  ApiOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { ragService, healthService } from '../services/api';

const { Title, Text, Paragraph } = Typography;

// Composant pour afficher l'état de santé d'un service
const HealthStatus = ({ name, status }) => {
  let icon;
  let color;
  
  switch (status?.status) {
    case 'ok':
      icon = <CheckCircleOutlined />;
      color = "var(--accent-green)";
      break;
    case 'warning':
      icon = <ExclamationCircleOutlined />;
      color = "var(--accent-yellow)";
      break;
    case 'error':
      icon = <CloseCircleOutlined />;
      color = "var(--accent-red)";
      break;
    default:
      icon = <SyncOutlined spin />;
      color = "var(--neutral-500)";
  }
  
  return (
    <Card className="mb-4">
      <Statistic 
        title={name}
        value={status?.status === 'ok' ? 'Available' : (status?.status === 'warning' ? 'Attention' : 'Not available')} 
        valueStyle={{ color }}
        prefix={icon}
      />
      {status?.message && (
        <Paragraph type={status.status === 'error' ? 'danger' : 'secondary'} style={{ marginTop: 8 }}>
          {status.message}
        </Paragraph>
      )}
      {status?.details && (
        <Paragraph type="secondary" style={{ fontSize: '0.85rem' }}>
          {status.details}
        </Paragraph>
      )}
    </Card>
  );
};

const Dashboard = () => {
  const [rags, setRags] = useState([]);
  const [loading, setLoading] = useState(true);
  const [systemHealth, setSystemHealth] = useState({
    rlama: { status: 'loading' },
    ollama: { status: 'loading' },
    models: { status: 'loading' },
    embeddings: { status: 'loading' }
  });
  const [refreshingHealth, setRefreshingHealth] = useState(false);
  const navigate = useNavigate();

  // Charger la liste des RAG
  const loadRags = async () => {
    try {
      setLoading(true);
      const data = await ragService.getAllRags();
      setRags(data);
    } catch (error) {
      console.error('Erreur lors du chargement des RAGs:', error);
      message.error('Unable to load RAG systems');
    } finally {
      setLoading(false);
    }
  };

  // Charger l'état de santé du système
  const loadSystemHealth = async () => {
    try {
      setRefreshingHealth(true);
      const health = await healthService.checkSystemHealth();
      setSystemHealth(health);
    } catch (error) {
      console.error('Error checking system status:', error);
    } finally {
      setRefreshingHealth(false);
    }
  };

  // Charger les données au montage du composant
  useEffect(() => {
    loadRags();
    loadSystemHealth();
  }, []);

  // Rafraîchir l'état de santé
  const handleRefreshHealth = () => {
    loadSystemHealth();
  };

  // Supprimer un RAG
  const handleDelete = async (ragName) => {
    try {
      await ragService.deleteRag(ragName);
      message.success(`RAG "${ragName}" supprimé avec succès`);
      loadRags(); // Recharger la liste
    } catch (error) {
      console.error(`Erreur lors de la suppression du RAG "${ragName}":`, error);
      message.error(`Impossible de supprimer le RAG "${ragName}"`);
    }
  };
  
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

  // Définition des colonnes du tableau
  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      render: (text) => (
        <div className="flex items-center gap-2">
          <Badge status="processing" color="var(--primary-600)" />
          <Text strong>{text}</Text>
        </div>
      ),
      sorter: (a, b) => a.name.localeCompare(b.name),
    },
    {
      title: 'Model',
      dataIndex: 'model',
      key: 'model',
      render: (text) => (
        <Tag color="var(--primary-100)" style={{ color: "var(--primary-700)", fontWeight: 500 }}>
          <RobotOutlined style={{ marginRight: 5 }} /> {text}
        </Tag>
      ),
      sorter: (a, b) => a.model.localeCompare(b.model),
    },
    {
      title: 'Creation Date',
      dataIndex: 'created_on',
      key: 'created_on',
      render: (text) => (
        <Tooltip title={`Créé le ${formatDate(text)}`}>
          <span className="flex items-center gap-2">
            <CalendarOutlined style={{ color: "var(--neutral-600)" }} />
            {formatDate(text)}
          </span>
        </Tooltip>
      ),
      sorter: (a, b) => new Date(a.created_on) - new Date(b.created_on),
    },
    {
      title: 'Documents',
      dataIndex: 'documents_count',
      key: 'documents_count',
      render: (count) => (
        <div className="flex items-center gap-2">
          <FileTextOutlined style={{ color: "var(--accent-purple)" }} />
          <span>{count}</span>
        </div>
      ),
      sorter: (a, b) => a.documents_count - b.documents_count,
    },
    {
      title: 'Size',
      dataIndex: 'size',
      key: 'size',
      render: (size) => {
        let color = 'var(--neutral-600)';
        
        if (size.includes('MB')) {
          const value = parseFloat(size);
          if (value > 10) color = 'var(--accent-yellow)';
          if (value > 50) color = 'var(--accent-red)';
        }
        
        return <Text style={{ color }}>{size}</Text>;
      },
      sorter: (a, b) => {
        const getSizeInBytes = (sizeStr) => {
          const num = parseFloat(sizeStr);
          if (sizeStr.includes('KB')) return num * 1024;
          if (sizeStr.includes('MB')) return num * 1024 * 1024;
          if (sizeStr.includes('GB')) return num * 1024 * 1024 * 1024;
          return num;
        };
        return getSizeInBytes(a.size) - getSizeInBytes(b.size);
      },
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space size="small">
          <Button 
            type="primary" 
            onClick={() => navigate(`/rag/${record.name}`)}
            icon={<SearchOutlined />}
          >
            Explore
          </Button>
          <Popconfirm
            title="Supprimer ce RAG"
            description="Are you sure you want to delete this RAG? This action is irreversible."
            onConfirm={() => handleDelete(record.name)}
            okText="Yes"
            cancelText="No"
            okButtonProps={{ danger: true }}
          >
            <Button 
              type="primary" 
              danger 
              icon={<DeleteOutlined />}
            >
              Delete
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <div>
          <Title level={2} style={{ color: 'var(--primary-800)', margin: 0 }}>My RAG Systems</Title>
          <Text type="secondary">Retrieval-Augmented Generation (RAG) system for accessing your documents</Text>
        </div>
        <Button 
          type="primary" 
          icon={<PlusOutlined />} 
          onClick={() => navigate('/create')}
          size="large"
          className="shadow-md"
        >
          Create new RAG
        </Button>
      </div>

      {/* Tableau des RAGs existants */}
      <Card className="shadow-md rounded-lg">
        <Table
          dataSource={rags}
          columns={columns}
          rowKey="name"
          loading={loading}
          pagination={{ 
            pageSize: 10,
            showTotal: (total, range) => `${range[0]}-${range[1]} of ${total} RAGs`,
          }}
          locale={{ 
            emptyText: (
              <div className="flex flex-col items-center justify-center p-8">
                <RobotOutlined style={{ fontSize: 48, color: 'var(--neutral-400)', marginBottom: 16 }} />
                <Text strong>No RAG system found</Text>
                <Text type="secondary" className="mb-4">Create your first RAG system to get started</Text>
                <Button 
                  type="primary" 
                  icon={<PlusOutlined />}
                  onClick={() => navigate('/create')}
                >
                  Create RAG
                </Button>
              </div>
            )
          }}
          rowClassName={() => 'fade-in'}
          onRow={(record) => ({
            onClick: () => {},
            className: 'hover:bg-primary-50 transition-colors cursor-pointer'
          })}
        />
      </Card>
    </div>
  );
};

export default Dashboard; 