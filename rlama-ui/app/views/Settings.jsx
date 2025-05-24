import React, { useState, useEffect } from 'react';
import { 
  Card, Form, Input, Button, Select, 
  Typography, Tabs, message, Space, 
  Divider, Switch, Alert, List, 
  Popconfirm, Modal, Tooltip, Table,
  Tag, Badge
} from 'antd';
import { 
  SettingOutlined, 
  SaveOutlined, 
  PlusOutlined,
  DeleteOutlined,
  EditOutlined,
  KeyOutlined,
  UserOutlined,
  EyeOutlined, 
  EyeInvisibleOutlined,
  GlobalOutlined,
  RobotOutlined,
  DollarOutlined,
  InfoCircleOutlined,
  ApiOutlined
} from '@ant-design/icons';
import { settingsService } from '../services/api';

const { Title, Text, Paragraph } = Typography;
const { Option } = Select;
const { TextArea } = Input;

// Updated OpenAI models with current pricing (January 2025)
const OPENAI_MODELS = [
  {
    category: "Reasoning Models (o-series)",
    models: [
      {
        name: "o3-mini",
        id: "o3-mini",
        description: "Latest reasoning model, 93% cheaper than o1",
        inputPrice: 1.10,
        outputPrice: 4.40,
        contextWindow: "200K tokens",
        recommended: true,
        new: true
      },
      {
        name: "o1-pro", 
        id: "o1-pro",
        description: "Most powerful reasoning model from OpenAI",
        inputPrice: 150.00,
        outputPrice: 600.00,
        contextWindow: "200K tokens",
        enterprise: true
      },
      {
        name: "o1",
        id: "o1", 
        description: "Advanced reasoning model",
        inputPrice: 15.00,
        outputPrice: 60.00,
        contextWindow: "200K tokens"
      }
    ]
  },
  {
    category: "GPT-4 Series",
    models: [
      {
        name: "GPT-4.5",
        id: "gpt-4.5",
        description: "Natural conversation, emotional intelligence",
        inputPrice: 75.00,
        outputPrice: 150.00,
        contextWindow: "128K tokens",
        new: true
      },
      {
        name: "GPT-4.1",
        id: "gpt-4.1",
        description: "Latest GPT-4 version with 1M context",
        inputPrice: 30.00,
        outputPrice: 60.00,
        contextWindow: "1M tokens",
        new: true
      },
      {
        name: "GPT-4.1-nano",
        id: "gpt-4.1-nano",
        description: "Lightweight version of GPT-4.1",
        inputPrice: 5.00,
        outputPrice: 15.00,
        contextWindow: "128K tokens",
        new: true
      },
      {
        name: "GPT-4",
        id: "gpt-4",
        description: "Latest GPT-4 model with vision support",
        inputPrice: 5.00,
        outputPrice: 15.00,
        contextWindow: "128K tokens",
        popular: true
      },
      {
        name: "GPT-4-turbo",
        id: "gpt-4-turbo", 
        description: "Efficient version of GPT-4",
        inputPrice: 0.15,
        outputPrice: 0.60,
        contextWindow: "128K tokens",
        recommended: true
      }
    ]
  },
  {
    category: "GPT-3.5 Series", 
    models: [
      {
        name: "GPT-3.5 Turbo",
        id: "gpt-3.5-turbo",
        description: "Fast and economical model",
        inputPrice: 0.50,
        outputPrice: 1.50,
        contextWindow: "16K tokens",
        budget: true
      }
    ]
  }
];

const Settings = () => {
  const [activeKey, setActiveKey] = useState('default-keys');
  const [loading, setLoading] = useState(false);
  const [profiles, setProfiles] = useState([]);
  const [apiKeys, setApiKeys] = useState({});
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [editingProfile, setEditingProfile] = useState(null);
  const [availableModels, setAvailableModels] = useState([]);
  
  const [profileForm] = Form.useForm();
  const [apiKeyForm] = Form.useForm();
  const [generalForm] = Form.useForm();

  // Load data on component mount
  useEffect(() => {
    loadSettings();
  }, []);

  // Debug: Log profiles state changes
  useEffect(() => {
    console.log('📊 Settings: Profiles state changed:', {
      length: profiles.length,
      profiles: profiles,
      isArray: Array.isArray(profiles)
    });
  }, [profiles]);

  const loadSettings = async () => {
    try {
      setLoading(true);
      
      console.log('🔄 Settings: Loading settings...');
      
      // Load profiles
      console.log('🔍 Settings: Loading profiles...');
      const profilesData = await settingsService.getProfiles();
      console.log('📋 Settings: Profiles data received:', profilesData);
      console.log('📋 Settings: Profiles type:', typeof profilesData);
      console.log('📋 Settings: Profiles is array?', Array.isArray(profilesData));
      
      // Now the API always returns an array, so we can simplify this
      if (Array.isArray(profilesData)) {
        console.log('✅ Settings: Setting', profilesData.length, 'profiles to state');
        setProfiles(profilesData);
      } else {
        console.log('⚠️ Settings: Received non-array data, setting empty array');
        setProfiles([]);
      }
      
      // Load API keys
      const apiKeysData = await settingsService.getApiKeys();
      setApiKeys(apiKeysData || {});
      
      // Load general settings
      const generalData = await settingsService.getGeneralSettings();
      
      // Load available models
      const modelsData = await settingsService.getAvailableModels();
      setAvailableModels(modelsData || []);
      
      // Set form values
      apiKeyForm.setFieldsValue(apiKeysData || {});
      generalForm.setFieldsValue(generalData || {
        auto_save: true,
        show_notifications: true,
        default_model: 'gpt-4o',
        default_embedding_model: 'text-embedding-3-small'
      });
      
      console.log('✅ Settings: All settings loaded successfully');
      
    } catch (error) {
      console.error('❌ Settings: Error loading settings:', error);
      message.error('Error loading settings');
    } finally {
      setLoading(false);
    }
  };

  // Profile management
  const showProfileModal = (profile = null) => {
    setEditingProfile(profile);
    setIsModalVisible(true);
    
    if (profile) {
      profileForm.setFieldsValue(profile);
    } else {
      profileForm.resetFields();
    }
  };

  const handleProfileSubmit = async (values) => {
    try {
      setLoading(true);
      
      console.log('📝 Profile form submitted with values:', values);
      
      // Client-side validation
      if (!values.name || !values.provider || !values.api_key) {
        console.error('❌ Missing required fields:', { name: !!values.name, provider: !!values.provider, api_key: !!values.api_key });
        throw new Error('Please fill in all required fields');
      }
      
      // Validate API key format - more lenient check
      if (values.provider === 'openai' && !values.api_key.startsWith('sk-')) {
        console.error('❌ Invalid API key format:', values.api_key.substring(0, 10) + '...');
        throw new Error('OpenAI API keys must start with "sk-"');
      }
      
      // More lenient API key length check
      if (values.provider === 'openai' && values.api_key.length < 10) {
        console.error('❌ API key too short:', values.api_key.length);
        throw new Error('OpenAI API key appears to be too short');
      }
      
      // Validate profile name (no spaces, special characters except hyphens and underscores)
      const nameRegex = /^[a-zA-Z0-9_-]+$/;
      if (!nameRegex.test(values.name)) {
        console.error('❌ Invalid profile name:', values.name);
        throw new Error('Profile name can only contain letters, numbers, hyphens, and underscores');
      }
      
      console.log('✅ Client-side validation passed');
      
      if (editingProfile) {
        console.log('🔄 Updating existing profile:', editingProfile.name);
        await settingsService.updateProfile(editingProfile.name, values);
        message.success('Profile updated successfully');
      } else {
        console.log('🆕 Creating new profile:', values.name);
        console.log('📤 Sending profile data to backend:', {
          name: values.name,
          provider: values.provider,
          api_key: `${values.api_key.substring(0, 15)}...`,
          description: values.description
        });
        
        const result = await settingsService.createProfile(values);
        console.log('✅ Profile creation result:', result);
        console.log('✅ Profile creation result type:', typeof result);
        console.log('✅ Profile creation result keys:', Object.keys(result || {}));
        
        message.success('Profile created successfully');
        
        // Wait a moment before reloading to ensure backend persistence
        console.log('⏳ Waiting 1 second before reloading profiles...');
        setTimeout(async () => {
          console.log('🔄 Now reloading profiles after creation...');
          await loadSettings();
        }, 1000);
      }
      
      setIsModalVisible(false);
      profileForm.resetFields();
      setEditingProfile(null);
      
      // Don't reload immediately for new profiles - we set a timeout above
      if (editingProfile) {
        loadSettings();
      }
      
    } catch (error) {
      console.error('❌ Error saving profile:', error);
      console.error('❌ Error details:', {
        message: error.message,
        response: error.response?.data,
        status: error.response?.status,
        headers: error.response?.headers
      });
      
      // Show user-friendly error message
      const errorMessage = error.message || 'Unknown error occurred';
      message.error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const deleteProfile = async (profileName) => {
    try {
      setLoading(true);
      await settingsService.deleteProfile(profileName);
      message.success('Profile deleted successfully');
      loadSettings();
    } catch (error) {
      console.error('Error deleting profile:', error);
      message.error('Error deleting profile');
    } finally {
      setLoading(false);
    }
  };

  // API Keys management
  const handleApiKeysSubmit = async (values) => {
    try {
      setLoading(true);
      
      // Save API keys
      await settingsService.saveApiKeys(values);
      setApiKeys(values);
      
      // If OpenAI API key is provided, set it as environment variable for CLI commands
      if (values.openai_api_key) {
        try {
          const command = `export OPENAI_API_KEY="${values.openai_api_key}"`;
          await settingsService.setEnvironmentVariable('OPENAI_API_KEY', values.openai_api_key);
        } catch (envError) {
          console.warn('Could not set environment variable:', envError);
        }
      }
      
      message.success('API keys saved successfully');
    } catch (error) {
      console.error('Error saving API keys:', error);
      message.error('Error saving API keys');
    } finally {
      setLoading(false);
    }
  };

  // General settings management
  const handleGeneralSettingsSubmit = async (values) => {
    try {
      setLoading(true);
      await settingsService.saveGeneralSettings(values);
      message.success('General settings saved successfully');
    } catch (error) {
      console.error('Error saving general settings:', error);
      message.error('Error saving general settings');
    } finally {
      setLoading(false);
    }
  };

  const PasswordInput = ({ value, onChange, placeholder }) => {
    const [visible, setVisible] = useState(false);
    
    return (
      <Input
        type={visible ? 'text' : 'password'}
        value={value}
        onChange={onChange}
        placeholder={placeholder}
        suffix={
          <Button
            type="text"
            size="small"
            icon={visible ? <EyeInvisibleOutlined /> : <EyeOutlined />}
            onClick={() => setVisible(!visible)}
          />
        }
      />
    );
  };

  // Function to render model pricing columns
  const modelColumns = [
    {
      title: 'Model',
      dataIndex: 'name',
      key: 'name',
      render: (text, record) => (
        <Space direction="vertical" size="small">
          <Space>
            <strong>{text}</strong>
            {record.new && <Tag color="green">New</Tag>}
            {record.recommended && <Tag color="blue">Recommended</Tag>}
            {record.popular && <Tag color="orange">Popular</Tag>}
            {record.enterprise && <Tag color="purple">Enterprise</Tag>}
            {record.budget && <Tag color="cyan">Budget</Tag>}
          </Space>
          <Text type="secondary" style={{ fontSize: '12px' }}>{record.description}</Text>
        </Space>
      )
    },
    {
      title: 'Input Price',
      dataIndex: 'inputPrice',
      key: 'inputPrice',
      render: (price) => <Text strong>${price}/1M tokens</Text>
    },
    {
      title: 'Output Price',
      dataIndex: 'outputPrice', 
      key: 'outputPrice',
      render: (price) => <Text strong>${price}/1M tokens</Text>
    },
    {
      title: 'Context',
      dataIndex: 'contextWindow',
      key: 'contextWindow'
    }
  ];

  const tabItems = [
    {
      key: 'default-keys',
      label: (
        <span>
          <KeyOutlined />
          Default API Keys
        </span>
      ),
      children: (
        <div>
          <Title level={4}>Default API Keys Configuration</Title>
          
          <Alert
            message="⚠️ Development Notice"
            description="OpenAI and other external APIs are currently not functional and are under development. This feature will be available in a future release."
            type="warning"
            showIcon
            style={{ marginBottom: 24 }}
          />
          
          <Alert
            message="Default API Keys"
            description={
              <div>
                <Paragraph>
                  Configure your default API keys here. These keys will be automatically used when you run RLAMA commands without specifying a profile.
                  <strong> This is the recommended approach for most users.</strong>
                </Paragraph>
                <Paragraph>
                  <strong>Usage examples with default keys:</strong>
                  <br />
                  • <code>rlama rag o3-mini my-rag ./docs</code> - Uses default OpenAI key
                  <br />
                  • <code>rlama update-model my-rag gpt-4o</code> - Uses default OpenAI key
                  <br />
                  • <code>rlama run my-rag</code> - Uses default OpenAI key for inference
                </Paragraph>
                <Paragraph>
                  <Text type="secondary">
                    💡 <strong>When to use Default API Keys vs Named Profiles:</strong>
                    <br />
                    • <strong>Default API Keys:</strong> You have one OpenAI account → Use this
                    <br />
                    • <strong>Named Profiles:</strong> You have multiple OpenAI accounts → Use profiles with --profile flag
                  </Text>
                </Paragraph>
              </div>
            }
            type="info"
            showIcon
            style={{ marginBottom: 24 }}
          />

          <Form
            form={apiKeyForm}
            layout="vertical"
            onFinish={handleApiKeysSubmit}
          >
            <Card title={<><RobotOutlined /> OpenAI</>} style={{ marginBottom: 16 }}>
              <Form.Item
                name="openai_api_key"
                label="OpenAI API Key"
                help="Used for GPT models (ChatGPT, GPT-4, o3-mini, etc.)"
                rules={[
                  {
                    pattern: /^sk-/,
                    message: 'Please enter a valid OpenAI API key (starts with sk-)'
                  }
                ]}
              >
                <PasswordInput placeholder="sk-proj-..." />
              </Form.Item>

              <Form.Item
                name="openai_organization"
                label="OpenAI Organization (optional)"
                help="Only needed if you're part of an OpenAI organization"
              >
                <Input placeholder="org-..." />
              </Form.Item>
              
              <Alert
                message="How to get your OpenAI API Key"
                description={
                  <div>
                    1. Go to <a href="https://platform.openai.com/api-keys" target="_blank" rel="noopener noreferrer">platform.openai.com/api-keys</a>
                    <br />
                    2. Click "Create new secret key"
                    <br />
                    3. Copy the key (starts with <code>sk-proj-</code> or <code>sk-</code>)
                    <br />
                    4. Paste it above and save
                  </div>
                }
                type="info"
                style={{ marginTop: 8 }}
              />
            </Card>

            <Card title={<><GlobalOutlined /> Google Search</>} style={{ marginBottom: 16 }}>
              <Form.Item
                name="google_api_key"
                label="Google API Key"
                help="Google Cloud API key for search services"
              >
                <PasswordInput placeholder="AIza..." />
              </Form.Item>

              <Form.Item
                name="google_search_engine_id"
                label="Custom Search Engine ID"
                help="Your Custom Search Engine (CSE) ID"
              >
                <Input placeholder="..." />
              </Form.Item>
            </Card>

            <Card title={<><RobotOutlined /> Anthropic</>} style={{ marginBottom: 24 }}>
              <Form.Item
                name="anthropic_api_key"
                label="Anthropic API Key"
                help="Used for Claude models"
              >
                <PasswordInput placeholder="sk-ant-..." />
              </Form.Item>
            </Card>

            <Form.Item>
              <Button 
                type="primary" 
                htmlType="submit" 
                icon={<SaveOutlined />}
                loading={loading}
                size="large"
              >
                Save Default API Keys
              </Button>
            </Form.Item>
          </Form>
        </div>
      ),
    },
    {
      key: 'profiles',
      label: (
        <span>
          <ApiOutlined />
          Named Profiles
        </span>
      ),
      children: (
        <div>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
            <Title level={4}>Named OpenAI Profiles</Title>
            <Space>
              <Button 
                type="default" 
                icon={<InfoCircleOutlined />}
                onClick={async () => {
                  console.log('🔍 Debug: Testing raw CLI output...');
                  try {
                    const result = await settingsService.debugProfilesRaw();
                    console.log('🔍 Debug result:', result);
                    message.info('Debug output logged to console - open DevTools to see details');
                  } catch (error) {
                    console.error('❌ Debug failed:', error);
                    message.error('Debug failed: ' + error.message);
                  }
                }}
              >
                Debug CLI
              </Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showProfileModal()}
              >
                New Profile
              </Button>
            </Space>
          </div>
          
          <Alert
            message="⚠️ Development Notice"
            description="OpenAI and other external APIs are currently not functional and are under development. This feature will be available in a future release."
            type="warning"
            showIcon
            style={{ marginBottom: 24 }}
          />
          
          <Alert
            message="Named Profiles Management"
            description={
              <div>
                <Paragraph>
                  Named profiles allow you to manage different OpenAI API keys for specific projects or environments.
                  Use these when you need to switch between different accounts or organizations.
                </Paragraph>
                <Paragraph>
                  <strong>Available CLI commands:</strong>
                  <br />
                  • <code>rlama profile add [name] openai [api-key]</code> - Add a profile
                  <br />
                  • <code>rlama profile list</code> - List all profiles
                  <br />
                  • <code>rlama profile delete [name]</code> - Delete a profile
                </Paragraph>
                <Paragraph>
                  <strong>Usage examples with named profiles:</strong>
                  <br />
                  • <code>rlama rag o3-mini my-rag ./docs --profile work-account</code>
                  <br />
                  • <code>rlama update-model my-rag gpt-4o --profile personal-account</code>
                </Paragraph>
              </div>
            }
            type="info"
            showIcon
            style={{ marginBottom: 24 }}
          />

          {console.log('🎯 About to render List with profiles:', profiles, 'loading:', loading)}
          
          <List
            loading={loading}
            dataSource={profiles}
            renderItem={(profile) => {
              console.log('🎨 Rendering profile:', profile);
              return (
                <List.Item
                  actions={[
                    <Tooltip title="Edit">
                      <Button 
                        type="text" 
                        icon={<EditOutlined />}
                        onClick={() => showProfileModal(profile)}
                      />
                    </Tooltip>,
                    <Popconfirm
                      title={<span style={{ color: '#ffffff' }}>Are you sure you want to delete this profile?</span>}
                      onConfirm={() => deleteProfile(profile.name)}
                      okText="Yes"
                      cancelText="No"
                    >
                      <Tooltip title="Delete">
                        <Button 
                          type="text" 
                          danger 
                          icon={<DeleteOutlined />}
                        />
                      </Tooltip>
                    </Popconfirm>
                  ]}
                >
                  <List.Item.Meta
                    avatar={<RobotOutlined style={{ fontSize: '20px', color: '#1890ff' }} />}
                    title={
                      <Space>
                        <strong style={{ fontSize: '16px', color: '#ffffff' }}>{profile.name}</strong>
                        <Badge status="success" text={profile.provider} />
                      </Space>
                    }
                    description={
                      <Space direction="vertical" size="small">
                        <Text type="secondary" style={{ color: '#bfbfbf' }}>Provider: {profile.provider}</Text>
                        <Text type="secondary" style={{ color: '#bfbfbf' }}>
                          Created: {profile.created_on ? new Date(profile.created_on).toLocaleDateString() : 'Unknown'}
                        </Text>
                        {profile.last_used && profile.last_used !== 'never' && (
                          <Text type="secondary" style={{ color: '#bfbfbf' }}>
                            Last used: {new Date(profile.last_used).toLocaleDateString()}
                          </Text>
                        )}
                        {profile.last_used === 'never' && (
                          <Text type="secondary" style={{ color: '#bfbfbf' }}>Last used: Never</Text>
                        )}
                        {profile.description && (
                          <Text type="secondary" italic style={{ color: '#bfbfbf' }}>{profile.description}</Text>
                        )}
                      </Space>
                    }
                  />
                </List.Item>
              );
            }}
          />

          <Modal
            title={editingProfile ? "Edit OpenAI Profile" : "New OpenAI Profile"}
            open={isModalVisible}
            onCancel={() => {
              setIsModalVisible(false);
              profileForm.resetFields();
              setEditingProfile(null);
            }}
            footer={null}
            width={600}
          >
            <Form
              form={profileForm}
              layout="vertical"
              onFinish={handleProfileSubmit}
            >
              <Form.Item
                name="name"
                label="Profile Name"
                rules={[
                  { required: true, message: 'Please enter a profile name' },
                  { 
                    pattern: /^[a-zA-Z0-9_-]+$/, 
                    message: 'Profile name can only contain letters, numbers, hyphens, and underscores' 
                  },
                  {
                    min: 2,
                    message: 'Profile name must be at least 2 characters long'
                  }
                ]}
                help="Use letters, numbers, hyphens, and underscores only (e.g., work-account, personal_openai)"
              >
                <Input 
                  placeholder="e.g., work-account, personal-account, project-x" 
                  disabled={!!editingProfile} 
                />
              </Form.Item>

              <Form.Item
                name="provider"
                label="Provider"
                rules={[{ required: true, message: 'Please select a provider' }]}
                initialValue="openai"
              >
                <Select placeholder="Select a provider" disabled>
                  <Option value="openai">
                    <Space>
                      <RobotOutlined />
                      OpenAI
                    </Space>
                  </Option>
                </Select>
              </Form.Item>

              <Form.Item
                name="api_key"
                label="OpenAI API Key"
                rules={[
                  { required: true, message: 'Please enter the API key' },
                  { 
                    pattern: /^sk-/, 
                    message: 'OpenAI API keys start with "sk-"' 
                  },
                  {
                    min: 10,
                    message: 'API key seems too short'
                  }
                ]}
                help="Your OpenAI API key from https://platform.openai.com/api-keys (starts with sk-)"
              >
                <PasswordInput placeholder="sk-proj-..." />
              </Form.Item>

              <Form.Item
                name="description"
                label="Description (optional)"
                help="Optional description to help you remember what this profile is for"
              >
                <TextArea 
                  rows={2} 
                  placeholder="e.g., Personal OpenAI account, Work organization key, Project Alpha..."
                  maxLength={200}
                  showCount
                />
              </Form.Item>

              <Alert
                message="Usage with RLAMA"
                description={
                  <div>
                    Once created, you can use this profile with:
                    <br />
                    • <code>rlama rag [model] [rag-name] [folder] --profile {profileForm.getFieldValue('name') || '[profile-name]'}</code>
                    <br />
                    • <code>rlama update-model [rag-name] [model] --profile {profileForm.getFieldValue('name') || '[profile-name]'}</code>
                    <br />
                    <Text type="secondary" style={{ fontSize: '12px' }}>
                      💡 <strong>Tip:</strong> For most users, using "Default API Keys" is simpler and doesn't require --profile flags.
                    </Text>
                  </div>
                }
                type="info"
                style={{ marginBottom: 16 }}
              />

              <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
                <Space>
                  <Button onClick={() => {
                    setIsModalVisible(false);
                    profileForm.resetFields();
                    setEditingProfile(null);
                  }}>
                    Cancel
                  </Button>
                  <Button type="primary" htmlType="submit" loading={loading}>
                    {editingProfile ? 'Update' : 'Create Profile'}
                  </Button>
                </Space>
              </Form.Item>
            </Form>
          </Modal>
        </div>
      ),
    },
    {
      key: 'openai-models',
      label: (
        <span>
          <DollarOutlined />
          Models & Pricing
        </span>
      ),
      children: (
        <div>
          <Title level={4}>OpenAI Models & Pricing</Title>
          
          <Alert
            message="Updated Pricing - January 2025"
            description="Here are the latest OpenAI models available with current pricing. Prices are in USD per million tokens."
            type="info"
            showIcon
            style={{ marginBottom: 24 }}
          />

          {OPENAI_MODELS.map((category, index) => (
            <Card 
              key={index}
              title={
                <Space>
                  <RobotOutlined />
                  {category.category}
                </Space>
              } 
              style={{ marginBottom: 16 }}
            >
              <Table
                dataSource={category.models}
                columns={modelColumns}
                pagination={false}
                size="small"
                rowKey="id"
              />
            </Card>
          ))}

          <Card 
            title={
              <Space>
                <InfoCircleOutlined />
                Usage Tips
              </Space>
            }
            style={{ marginTop: 16 }}
          >
            <div>
              <Paragraph>
                <strong>Choosing the right model:</strong>
              </Paragraph>
              <ul>
                <li><strong>o3-mini</strong> : Excellent choice for reasoning at reduced cost (93% cheaper than o1)</li>
                <li><strong>GPT-4o</strong> : Versatile with multimodal support (images, audio)</li>
                <li><strong>GPT-4o mini</strong> : Budget-friendly for simple tasks</li>
                <li><strong>o1-pro</strong> : For complex reasoning tasks (enterprise)</li>
                <li><strong>GPT-3.5 Turbo</strong> : Most economical for basic tasks</li>
              </ul>
              
              <Paragraph style={{ marginTop: 16 }}>
                <strong>Cost optimization:</strong>
                <br />
                • Use context caching (50% reduction)
                <br />
                • Choose appropriate context window size
                <br />
                • Test multiple models for your use case
              </Paragraph>
            </div>
          </Card>
        </div>
      ),
    },
    {
      key: 'general',
      label: (
        <span>
          <SettingOutlined />
          General
        </span>
      ),
      children: (
        <div>
          <Title level={4}>General Settings</Title>
          
          <Form
            form={generalForm}
            layout="vertical"
            onFinish={handleGeneralSettingsSubmit}
            initialValues={{
              auto_save: true,
              show_notifications: true,
              default_model: 'gpt-4o',
              default_embedding_model: 'text-embedding-3-small'
            }}
          >
            <Card title="Interface" style={{ marginBottom: 16 }}>
              <Form.Item
                name="show_notifications"
                label="Notifications"
                valuePropName="checked"
              >
                <Switch checkedChildren="Enabled" unCheckedChildren="Disabled" />
              </Form.Item>

              <Form.Item
                name="auto_save"
                label="Auto Save"
                valuePropName="checked"
              >
                <Switch checkedChildren="Enabled" unCheckedChildren="Disabled" />
              </Form.Item>
            </Card>

            <Card title="Default Models" style={{ marginBottom: 24 }}>
              <Form.Item
                name="default_model"
                label="Default LLM Model"
              >
                <Select placeholder="Select a model" loading={loading}>
                  {availableModels.map(model => (
                    <Option key={model} value={model}>{model}</Option>
                  ))}
                </Select>
              </Form.Item>

              <Form.Item
                name="default_embedding_model"
                label="Default Embedding Model"
              >
                <Select placeholder="Select an embedding model">
                  <Option value="text-embedding-ada-002">text-embedding-ada-002</Option>
                  <Option value="text-embedding-3-small">text-embedding-3-small</Option>
                  <Option value="text-embedding-3-large">text-embedding-3-large</Option>
                </Select>
              </Form.Item>
            </Card>

            <Form.Item>
              <Button 
                type="primary" 
                htmlType="submit" 
                icon={<SaveOutlined />}
                loading={loading}
                size="large"
              >
                Save General Settings
              </Button>
            </Form.Item>
          </Form>
        </div>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto' }}>
      <Title level={2}>
        <SettingOutlined style={{ marginRight: '8px' }} />
        Settings
      </Title>
      
      <Tabs
        activeKey={activeKey}
        onChange={setActiveKey}
        items={tabItems}
        size="large"
      />
    </div>
  );
};

export default Settings; 