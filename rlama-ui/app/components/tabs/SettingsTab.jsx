import React, { useState, useEffect } from 'react';
import { 
  Card, Form, Input, Button, Select, 
  InputNumber, Switch, Space, Typography, 
  Tabs, Divider, message, Alert, Descriptions
} from 'antd';
import { 
  FolderOpenOutlined, 
  SaveOutlined, 
  SyncOutlined,
  EyeOutlined, 
  EyeInvisibleOutlined,
  GlobalOutlined
} from '@ant-design/icons';
import { ragService } from '../../services/api';

const { Title, Text, Paragraph } = Typography;
const { Option } = Select;

const SettingsTab = ({ ragName, ragInfo }) => {
  const [watchForm] = Form.useForm();
  const [webWatchForm] = Form.useForm();
  const [modelForm] = Form.useForm();
  
  const [loading, setLoading] = useState(false);
  const [models, setModels] = useState([]);
  const [loadingModels, setLoadingModels] = useState(true);
  const [activeKey, setActiveKey] = useState('watch');
  const [isFolderWatchActive, setIsFolderWatchActive] = useState(false);
  const [isWebWatchActive, setIsWebWatchActive] = useState(false);
  const [watchSettings, setWatchSettings] = useState(null);
  const [webWatchSettings, setWebWatchSettings] = useState(null);

  // Load available models and check monitoring status
  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoadingModels(true);
        
        // Get available models
        const availableModels = await ragService.getAvailableModels();
        setModels(availableModels);
        
        // Check folder monitoring status
        try {
          const watchStatus = await ragService.getWatchStatus(ragName);
          if (watchStatus && watchStatus.active) {
            setIsFolderWatchActive(true);
            setWatchSettings(watchStatus);
            
            // Pre-fill the form with current settings
            watchForm.setFieldsValue({
              folder_path: watchStatus.folder_path,
              interval: watchStatus.interval || 0
            });
          }
        } catch (watchError) {
          console.warn("Error checking folder monitoring status:", watchError);
        }
        
        // Check web monitoring status
        try {
          const webWatchStatus = await ragService.getWebWatchStatus(ragName);
          if (webWatchStatus && webWatchStatus.active) {
            setIsWebWatchActive(true);
            setWebWatchSettings(webWatchStatus);
            
            // Pre-fill the form with current settings
            webWatchForm.setFieldsValue({
              url: webWatchStatus.url,
              interval: webWatchStatus.interval || 0,
              depth: webWatchStatus.depth || 1
            });
          }
        } catch (webWatchError) {
          console.warn("Error checking web monitoring status:", webWatchError);
        }
      } catch (error) {
        console.error('Error loading data:', error);
      } finally {
        setLoadingModels(false);
      }
    };
    
    fetchData();
  }, [ragName]);

  // Handle folder selection
  const handleSelectFolder = async () => {
    if (window.electron) {
      try {
        const folderPath = await window.electron.selectDirectory();
        if (folderPath) {
          watchForm.setFieldsValue({ folder_path: folderPath });
        }
      } catch (error) {
        console.error('Error selecting folder:', error);
        message.error('Unable to select folder');
      }
    } else {
      message.warning('Folder selection via interface is not available');
    }
  };

  // Configure folder monitoring
  const setupWatch = async (values) => {
    try {
      setLoading(true);
      
      const watchData = {
        rag_name: ragName,
        folder_path: values.folder_path,
        interval: values.interval || 0
      };
      
      const result = await ragService.setupWatch(watchData);
      
      if (result) {
        message.success('Folder monitoring configured successfully');
        setIsFolderWatchActive(true);
        setWatchSettings(watchData);
        
        // Check for new documents immediately after activation
        try {
          message.info('Checking for new documents...');
          await ragService.checkWatched(ragName);
          
          // Refresh interface after a short delay
          setTimeout(() => {
            window.location.reload();
          }, 2000);
        } catch (checkError) {
          console.error('Error checking for documents:', checkError);
        }
      } else {
        message.error('Error configuring monitoring');
      }
    } catch (error) {
      console.error('Error setting up folder monitoring:', error);
      message.error(`Unable to configure monitoring: ${error.response?.data?.detail || error.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Configure web monitoring
  const setupWebWatch = async (values) => {
    try {
      setLoading(true);
      
      const webWatchData = {
        rag_name: ragName,
        url: values.url,
        interval: values.interval || 0,
        depth: values.depth || 1
      };
      
      const result = await ragService.setupWebWatch(webWatchData);
      
      if (result) {
        message.success('Web monitoring configured successfully');
        setIsWebWatchActive(true);
        setWebWatchSettings(webWatchData);
      } else {
        message.error('Error configuring web monitoring');
      }
    } catch (error) {
      console.error('Error setting up web monitoring:', error);
      message.error(`Unable to configure web monitoring: ${error.response?.data?.detail || error.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Disable folder monitoring
  const disableWatch = async () => {
    try {
      setLoading(true);
      const result = await ragService.disableWatch(ragName);
      
      if (result) {
        message.success('Folder monitoring disabled');
        setIsFolderWatchActive(false);
        setWatchSettings(null);
        watchForm.resetFields();
        
        // Refresh interface after disabling
        setTimeout(() => {
          window.location.reload();
        }, 1000);
      } else {
        message.error('Error disabling monitoring');
      }
    } catch (error) {
      console.error('Error disabling folder monitoring:', error);
      message.error(`Unable to disable monitoring: ${error.response?.data?.detail || error.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Disable web monitoring
  const disableWebWatch = async () => {
    try {
      setLoading(true);
      const result = await ragService.disableWebWatch(ragName);
      
      if (result) {
        message.success('Web monitoring disabled');
        setIsWebWatchActive(false);
        setWebWatchSettings(null);
        webWatchForm.resetFields();
      } else {
        message.error('Error disabling web monitoring');
      }
    } catch (error) {
      console.error('Error disabling web monitoring:', error);
      message.error(`Unable to disable web monitoring: ${error.response?.data?.detail || error.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Update LLM model
  const updateModel = async (values) => {
    try {
      setLoading(true);
      await ragService.updateModel(ragName, values.model);
      message.success(`Model updated to "${values.model}"`);
      // In a real app, we would need to update ragInfo
    } catch (error) {
      console.error('Error updating model:', error);
      message.error(`Unable to update model: ${error.response?.data?.detail || error.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Check monitored resources
  const checkWatched = async () => {
    try {
      setLoading(true);
      const result = await ragService.checkWatched(ragName);
      
      if (result) {
        message.success('Verification completed successfully');
        
        // Force application refresh
        setTimeout(() => {
          window.location.reload(); // Refresh page to see new documents
        }, 1500);
      } else {
        message.error('Error checking monitored resources');
      }
    } catch (error) {
      console.error('Error checking monitored resources:', error);
      message.error(`Unable to check resources: ${error.response?.data?.detail || error.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Style pour forcer le texte blanc
  const whiteTextStyle = {
    color: 'var(--text-primary)',
  };

  const strongWhiteTextStyle = {
    color: '#fafafa !important',
  };

  // Settings tabs
  const tabItems = [
    {
      key: 'watch',
      label: (
        <span style={whiteTextStyle}>
          <FolderOpenOutlined />
          Folder Monitoring
        </span>
      ),
      children: (
        <div style={whiteTextStyle}>
          <Alert
            message="Feature in development"
            description="Automatic folder monitoring is currently under development and will be available in a future version. Thank you for your patience."
            type="warning"
            showIcon
            style={{ marginBottom: 24 }}
          />
          
          <div style={{ opacity: 0.6, pointerEvents: 'none', ...whiteTextStyle }}>
            <Paragraph style={whiteTextStyle}>
              Configure automatic monitoring of a folder to update this RAG 
              when new files are added.
            </Paragraph>
            
            <Form
              form={watchForm}
              layout="vertical"
            >
              <Form.Item
                name="folder_path"
                label={<span style={whiteTextStyle}>Folder to monitor</span>}
                rules={[{ required: true, message: 'Please select a folder to monitor' }]}
              >
                <Input 
                  placeholder="/path/to/folder" 
                  style={whiteTextStyle}
                  addonAfter={
                    <Button 
                      type="text" 
                      icon={<FolderOpenOutlined />}
                      disabled
                    />
                  } 
                  disabled
                />
              </Form.Item>
              
              <Form.Item
                name="interval"
                label={<span style={whiteTextStyle}>Check interval (minutes)</span>}
                help={<span style={whiteTextStyle}>0 = only when used</span>}
                initialValue={0}
              >
                <InputNumber min={0} max={1440} style={{ width: 200, ...whiteTextStyle }} disabled />
              </Form.Item>
              
              <Form.Item>
                <Button 
                  type="primary" 
                  htmlType="submit" 
                  icon={<EyeOutlined />} 
                  disabled
                >
                  Enable monitoring
                </Button>
              </Form.Item>
            </Form>
          </div>
        </div>
      ),
    },
    {
      key: 'webwatch',
      label: (
        <span style={whiteTextStyle}>
          <GlobalOutlined />
          Web Monitoring
        </span>
      ),
      children: (
        <div style={whiteTextStyle}>
          <Alert
            message="Feature in development"
            description="Automatic web monitoring is currently under development and will be available in a future version. Thank you for your patience."
            type="warning"
            showIcon
            style={{ marginBottom: 24 }}
          />
          
          <div style={{ opacity: 0.6, pointerEvents: 'none', ...whiteTextStyle }}>
            <Paragraph style={whiteTextStyle}>
              Configure automatic monitoring of a website to update this RAG 
              when new pages are added or modified.
            </Paragraph>
            
            <Form
              form={webWatchForm}
              layout="vertical"
            >
              <Form.Item
                name="url"
                label={<span style={whiteTextStyle}>URL to monitor</span>}
                rules={[
                  { required: true, message: 'Please enter a URL to monitor' },
                  { type: 'url', message: 'Please enter a valid URL' }
                ]}
              >
                <Input placeholder="https://example.com" style={whiteTextStyle} disabled />
              </Form.Item>
              
              <Form.Item
                name="interval"
                label={<span style={whiteTextStyle}>Check interval (minutes)</span>}
                help={<span style={whiteTextStyle}>0 = only when used</span>}
                initialValue={0}
              >
                <InputNumber min={0} max={1440} style={{ width: 200, ...whiteTextStyle }} disabled />
              </Form.Item>
              
              <Form.Item
                name="depth"
                label={<span style={whiteTextStyle}>Crawl depth</span>}
                initialValue={1}
              >
                <InputNumber min={1} max={10} style={{ width: 200, ...whiteTextStyle }} disabled />
              </Form.Item>
              
              <Form.Item>
                <Button 
                  type="primary" 
                  htmlType="submit" 
                  icon={<EyeOutlined />} 
                  disabled
                >
                  Enable web monitoring
                </Button>
              </Form.Item>
            </Form>
          </div>
        </div>
      ),
    },
    {
      key: 'model',
      label: (
        <span style={whiteTextStyle}>
          <SaveOutlined />
          LLM Model
        </span>
      ),
      children: (
        <div style={whiteTextStyle}>
          <Paragraph style={whiteTextStyle}>
            Update the LLM model used by this RAG. Note that existing embeddings won't be recalculated.
          </Paragraph>
          
          <Alert
            message="Caution"
            description="Changing the LLM model without recalculating embeddings may affect result quality. It's recommended to create a new RAG for a complete model change."
            type="warning"
            showIcon
            style={{ marginBottom: 16 }}
          />
          
          <Form
            form={modelForm}
            layout="vertical"
            onFinish={updateModel}
            initialValues={{
              model: ragInfo.model
            }}
          >
            <Form.Item
              name="model"
              label={<span style={whiteTextStyle}>New LLM model</span>}
              rules={[{ required: true, message: 'Please select a model' }]}
            >
              <Select
                placeholder="Select a model"
                loading={loadingModels}
                disabled={loadingModels || loading}
                style={whiteTextStyle}
              >
                {models.map(model => (
                  <Option key={model} value={model} style={whiteTextStyle}>{model}</Option>
                ))}
              </Select>
            </Form.Item>
            
            <Form.Item>
              <Button 
                type="primary" 
                htmlType="submit" 
                loading={loading}
                disabled={loading || modelForm.getFieldValue('model') === ragInfo.model}
              >
                Update model
              </Button>
            </Form.Item>
          </Form>
        </div>
      ),
    },
    {
      key: 'info',
      label: (
        <span style={whiteTextStyle}>
          <InfoCircleOutlined />
          Information
        </span>
      ),
      children: (
        <div style={whiteTextStyle}>
          <Card title={<span style={whiteTextStyle}>RAG Information</span>}>
            <Descriptions bordered column={1}>
              <Descriptions.Item label={<span style={whiteTextStyle}>Name</span>}>
                <span style={whiteTextStyle}>{ragInfo.name}</span>
              </Descriptions.Item>
              <Descriptions.Item label={<span style={whiteTextStyle}>Model</span>}>
                <span style={whiteTextStyle}>{ragInfo.model}</span>
              </Descriptions.Item>
              <Descriptions.Item label={<span style={whiteTextStyle}>Creation Date</span>}>
                <span style={whiteTextStyle}>{new Date(ragInfo.created_on).toLocaleString()}</span>
              </Descriptions.Item>
              <Descriptions.Item label={<span style={whiteTextStyle}>Document Count</span>}>
                <span style={whiteTextStyle}>{ragInfo.documents_count}</span>
              </Descriptions.Item>
              <Descriptions.Item label={<span style={whiteTextStyle}>Total Size</span>}>
                <span style={whiteTextStyle}>{ragInfo.size}</span>
              </Descriptions.Item>
            </Descriptions>
          </Card>
        </div>
      ),
    },
  ];

  return (
    <div style={whiteTextStyle}>
      <Title level={4} style={whiteTextStyle}>Settings</Title>
      
      <Tabs
        activeKey={activeKey}
        onChange={setActiveKey}
        items={tabItems}
        tabPosition="left"
        style={{ minHeight: '400px', ...whiteTextStyle }}
      />
    </div>
  );
};

// Note: InfoCircleOutlined is missing from imports
// Simulate it here to avoid modifying imports
const InfoCircleOutlined = () => <span>ℹ️</span>;

export default SettingsTab; 