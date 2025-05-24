import React, { useState, useEffect } from 'react';
import { 
  Card, Form, Input, Button, Select, 
  InputNumber, Switch, Space, Typography, 
  Divider, message, Spin, Alert, Steps, Tooltip, Badge
} from 'antd';
import { 
  FolderOpenOutlined, 
  SaveOutlined, 
  ArrowLeftOutlined, 
  InfoCircleOutlined,
  AppstoreOutlined,
  RobotOutlined,
  SettingOutlined,
  FileTextOutlined,
  HomeOutlined
} from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import { ragService } from '../services/api';

const { Title, Text, Paragraph } = Typography;
const { Option } = Select;

const CreateRag = () => {
  const [form] = Form.useForm();
  const [models, setModels] = useState([]);
  const [loading, setLoading] = useState(false);
  const [loadingModels, setLoadingModels] = useState(true);
  const [enableReranker, setEnableReranker] = useState(true);
  const navigate = useNavigate();

  // Charger les modèles disponibles au chargement
  useEffect(() => {
    const fetchModels = async () => {
      try {
        const availableModels = await ragService.getAvailableModels();
        setModels(availableModels);
      } catch (error) {
        console.error('Error loading models:', error);
        message.error('Unable to load available LLM models');
      } finally {
        setLoadingModels(false);
      }
    };
    
    fetchModels();
  }, []);

  // Gérer la sélection de dossier
  const handleSelectFolder = async () => {
    if (window.electron) {
      try {
        const folderPath = await window.electron.selectDirectory();
        if (folderPath) {
          form.setFieldsValue({ folder_path: folderPath });
        }
      } catch (error) {
        console.error('Error selecting folder:', error);
        message.error('Unable to select folder');
      }
    } else {
      message.warning('Folder selection via interface is not available');
    }
  };

  // Soumettre le formulaire
  const handleSubmit = async (values) => {
    try {
      setLoading(true);
      
      // Préparer les données pour la création du RAG
      const ragData = {
        name: values.name,
        model: values.model,
        folder_path: values.folder_path,
        chunk_size: values.chunk_size,
        chunk_overlap: values.chunk_overlap,
        enable_reranker: values.enable_reranker,
        reranker_weight: values.reranker_weight
      };
      
      // Appeler l'API pour créer le RAG
      await ragService.createRag(ragData);
      
      message.success(`RAG "${values.name}" created successfully`);
      navigate('/'); // Rediriger vers le tableau de bord
    } catch (error) {
      console.error('Error creating RAG:', error);
      message.error(`Unable to create RAG: ${error.response?.data?.detail || error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const steps = [
    {
      title: 'Configuration',
      description: 'Basic Parameters',
      icon: <AppstoreOutlined />
    },
    {
      title: 'Indexation',
      description: 'Traitement des documents',
      icon: <FileTextOutlined />
    },
    {
      title: 'Utilisation',
      description: 'Ready-to-use RAG',
      icon: <RobotOutlined />
    }
  ];

  return (
    <div className="fade-in create-rag-container">
      <div className="mb-4">
        <div className="flex items-center gap-1">
          <Link to="/" className="flex items-center gap-1">
            <HomeOutlined /> Dashboard
          </Link>
          {' / '}
          <span className="text-primary-700">New RAG</span>
        </div>
      </div>

      <div className="flex justify-between items-center mb-6">
        <div>
          <Title level={2} style={{ color: 'var(--primary-800)', margin: 0 }}>Create a new RAG system</Title>
          <Text type="secondary">
            Configure your Retrieval-Augmented Generation system to query your documents with AI
          </Text>
        </div>
        <Button 
          icon={<ArrowLeftOutlined />} 
          onClick={() => navigate('/')}
        >
          Back
        </Button>
      </div>
      
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <Card className="shadow-md rounded-lg">
            <Form
              form={form}
              layout="vertical"
              onFinish={handleSubmit}
              initialValues={{
                chunk_size: 1000,
                chunk_overlap: 200,
                enable_reranker: true,
                reranker_weight: 0.7,
                model: loadingModels ? undefined : (models.length > 0 ? models[0] : undefined)
              }}
              className="max-w-3xl"
            >
              <Title level={4} className="flex items-center gap-2 mb-4">
                <Badge status="processing" color="var(--primary-600)" />
                RAG Configuration
              </Title>
              
              <Form.Item
                name="name"
                label="RAG Name"
                rules={[
                  { required: true, message: 'Please enter a unique name for your RAG system' },
                  { pattern: /^[a-zA-Z0-9_-]+$/, message: 'The name must contain only letters, numbers, dashes, and underscores' }
                ]}
                tooltip="A unique name to identify your RAG system"
              >
                <Input placeholder="mon_rag" className="rounded-md" />
              </Form.Item>
              
              <Form.Item
                name="model"
                label={
                  <span className="flex items-center gap-2">
                    LLM Model
                    {loadingModels && <Spin size="small" />}
                  </span>
                }
                rules={[{ required: true, message: 'Please select an LLM model' }]}
                tooltip="Language model used for embeddings and response generation"
              >
                <Select
                  placeholder="Select a model"
                  loading={loadingModels}
                  disabled={loadingModels}
                  className="rounded-md"
                  suffixIcon={<RobotOutlined />}
                >
                  {models.map(model => (
                    <Option key={model} value={model}>{model}</Option>
                  ))}
                </Select>
              </Form.Item>
              
              <Form.Item
                name="folder_path"
                label="Source Folder"
                rules={[{ required: true, message: 'Please select a source folder' }]}
                tooltip="Folder containing documents to index"
              >
                <Input 
                  placeholder="/path/to/documents" 
                  className="rounded-md"
                  addonAfter={
                    <Button 
                      type="text" 
                      icon={<FolderOpenOutlined />} 
                      onClick={handleSelectFolder}
                    />
                  } 
                />
              </Form.Item>
              
              <Divider className="my-6" />
              
              <div className="flex items-center gap-2 mb-4">
                <SettingOutlined style={{ color: 'var(--primary-700)' }} />
                <Title level={4} style={{ margin: 0 }}>Advanced Parameters</Title>
                <Tooltip title="These parameters refine how documents are indexed and searched">
                  <InfoCircleOutlined style={{ color: 'var(--neutral-600)' }} />
                </Tooltip>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Form.Item
                  name="chunk_size"
                  label="Chunk size (characters)"
                  rules={[{ required: true, message: 'Please define a chunk size' }]}
                  tooltip="Number of characters per text segment"
                >
                  <InputNumber min={100} max={10000} step={100} className="w-full" />
                </Form.Item>
                
                <Form.Item
                  name="chunk_overlap"
                  label="Chunk overlap (characters)"
                  rules={[{ required: true, message: 'Please define an overlap' }]}
                  tooltip="Number of characters that overlap between segments"
                >
                  <InputNumber min={0} max={1000} step={50} className="w-full" />
                </Form.Item>
              </div>
              
              <div className="bg-primary-50 p-4 rounded-lg mb-6">
                <Form.Item
                  name="enable_reranker"
                  label={
                    <span className="flex items-center gap-2">
                      Enable result reranking
                      <Tooltip title="Reranking improves search result relevance">
                        <InfoCircleOutlined style={{ color: 'var(--primary-700)' }} />
                      </Tooltip>
                    </span>
                  }
                  valuePropName="checked"
                  className="mb-2"
                >
                  <Switch 
                    onChange={(checked) => setEnableReranker(checked)} 
                    defaultChecked={true} 
                  />
                </Form.Item>
                
                {enableReranker && (
                  <Form.Item
                    name="reranker_weight"
                    label="Reranker weight (0-1)"
                    rules={[{ required: enableReranker, message: 'Please define a weight' }]}
                    tooltip="Influence of reranker on final results (0 = none, 1 = maximum)"
                  >
                    <InputNumber 
                      min={0} 
                      max={1} 
                      step={0.1} 
                      className="w-full" 
                    />
                  </Form.Item>
                )}
              </div>
              
              <Alert
                message="Information on indexing"
                description="Creating a RAG may take several minutes depending on the size of the corpus to index. Please wait while the documents are being processed."
                type="info"
                showIcon
                className="mb-6"
              />
              
              <Form.Item>
                <div className="flex gap-3">
                  <Button 
                    type="primary" 
                    htmlType="submit" 
                    icon={<SaveOutlined />} 
                    loading={loading}
                    disabled={loading}
                    size="large"
                    className="shadow-md"
                  >
                    Create RAG
                  </Button>
                  <Button 
                    onClick={() => navigate('/')} 
                    disabled={loading}
                    size="large"
                  >
                    Cancel
                  </Button>
                </div>
              </Form.Item>
            </Form>
          </Card>
        </div>
        
        <div>
          <Card className="shadow-md rounded-lg">
            <div className="mb-4">
              <Title level={5} style={{ color: 'var(--text-primary)' }}>Creation Process</Title>
              <Text type="secondary" style={{ color: 'var(--text-secondary)' }}>Your RAG will be created following these steps</Text>
            </div>
            
            <Steps 
              direction="vertical" 
              current={0} 
              items={steps}
              className="mb-4"
              style={{ color: 'var(--text-primary)' }}
            />
            
            <div className="p-4 rounded-lg mt-6" style={{ background: 'var(--bg-tertiary)' }}>
              <Title level={5} className="mb-2" style={{ color: 'var(--text-primary)' }}>What is a RAG system?</Title>
              <Paragraph style={{ color: 'var(--text-secondary)' }}>
                A RAG (Retrieval-Augmented Generation) system combines document search and text generation to produce accurate answers based on your own documents.
              </Paragraph>
              <Paragraph style={{ color: 'var(--text-secondary)' }}>
                Once created, you can ask questions in natural language and get contextualized answers with source references.
              </Paragraph>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default CreateRag; 