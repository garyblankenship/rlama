import React from 'react';
import { Card, Button, Typography, Space, Alert, Steps, message } from 'antd';
import { ExclamationCircleOutlined, CheckCircleOutlined, PlayCircleOutlined, GlobalOutlined, DatabaseOutlined } from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;
const { Step } = Steps;

const InstallationHelper = ({ missingTools = [], onInstallComplete }) => {

  const getInstallationUrl = (tool) => {
    if (tool === 'ollama') {
      return 'https://ollama.com/';
    }
    
    if (tool === 'rlama') {
      return 'https://rlama.dev/download';
    }
  };

  const getInstallInstructions = (tool) => {
    if (tool === 'ollama') {
      return {
        title: 'Install Ollama',
        steps: [
          'Visit the official Ollama website',
          'Download the installer for your operating system',
          'Follow the installation instructions provided',
          'After installation, run: ollama pull llama3.2:3b (or any model you prefer)',
          'Verify installation with "ollama --version"'
        ]
      };
    }
    
    if (tool === 'rlama') {
      return {
        title: 'Install RLAMA',
        steps: [
          'Visit the official RLAMA download page',
          'Follow the installation instructions for your OS',
          'Make sure Ollama is installed first (requirement)',
          'Verify installation with "rlama --version"'
        ]
      };
    }
  };

  const handleVisitWebsite = (tool) => {
    const url = getInstallationUrl(tool);
    window.open(url, '_blank');
    message.success(`Redirected to ${tool.toUpperCase()} official website`);
  };



  const renderToolCard = (tool) => {
    const instructions = getInstallInstructions(tool);
    
    return (
      <Card
        key={tool}
        title={
          <Space>
            <ExclamationCircleOutlined style={{ color: '#faad14' }} />
            <span>{tool.toUpperCase()} Not Found</span>
          </Space>
        }
        style={{ marginBottom: 16 }}
        extra={
          <Button 
            type="primary" 
            icon={<PlayCircleOutlined />}
            onClick={() => handleVisitWebsite(tool)}
          >
            Install {tool.toUpperCase()}
          </Button>
        }
      >
        <Paragraph>
          <Text strong>{tool.toUpperCase()}</Text> is required for RLAMA to function properly.
          {tool === 'ollama' && ' Ollama provides local LLM capabilities.'}
          {tool === 'rlama' && ' RLAMA CLI provides RAG system management.'}
        </Paragraph>

        {/* Special note for Ollama about models */}
        {tool === 'ollama' && (
          <Alert
            message="Don't forget to install a model!"
            description={
              <>
                After installing Ollama, you'll need to download at least one model to use with RLAMA.
                <br />
                <Text code>ollama pull llama3.2:3b</Text> is a good lightweight option to start with.
              </>
            }
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
          />
        )}

        <Space direction="vertical" style={{ width: '100%' }}>
          <Title level={5}>{instructions.title}</Title>
          <Steps size="small" direction="vertical">
            {instructions.steps.map((step, index) => (
              <Step 
                key={index} 
                title={step}
                status="wait"
              />
            ))}
          </Steps>
          
          <Space>
            <Button 
              type="primary" 
              icon={<GlobalOutlined />}
              onClick={() => handleVisitWebsite(tool)}
            >
              Go to {tool.toUpperCase()} Website
            </Button>
            
            {tool === 'ollama' && (
              <Button 
                icon={<DatabaseOutlined />}
                onClick={() => window.open('https://ollama.com/models', '_blank')}
              >
                Browse Models
              </Button>
            )}
          </Space>
        </Space>
      </Card>
    );
  };

  if (!missingTools || missingTools.length === 0) {
    return (
      <Alert
        message="All Dependencies Available"
        description="RLAMA and Ollama are properly installed and accessible."
        type="success"
        icon={<CheckCircleOutlined />}
        showIcon
      />
    );
  }

  return (
    <>
      <Alert
        message="Missing Dependencies"
        description={`The following tools need to be installed: ${missingTools.join(', ').toUpperCase()}`}
        type="warning"
        showIcon
        style={{ marginBottom: 16 }}
      />
      
      {missingTools.map(tool => renderToolCard(tool))}
    </>
  );
};

export default InstallationHelper; 