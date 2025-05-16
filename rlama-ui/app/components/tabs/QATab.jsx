import React, { useState, useEffect, useRef } from 'react';
import { 
  Input, Button, List, Avatar, Spin, Select, 
  InputNumber, Tooltip, Typography, message, 
  Divider, Tag, Card, Space
} from 'antd';
import { 
  SendOutlined, 
  ClearOutlined, 
  CopyOutlined, 
  CheckOutlined, 
  QuestionCircleOutlined,
  RobotOutlined,
  UserOutlined,
  SettingOutlined,
  InfoCircleOutlined,
  LinkOutlined
} from '@ant-design/icons';
import { ragService } from '../../services/api';
import ReactMarkdown from 'react-markdown';
import rehypeHighlight from 'rehype-highlight';
import hljs from 'highlight.js';
import 'highlight.js/styles/github-dark.css';

const { TextArea } = Input;
const { Option } = Select;
const { Paragraph, Text, Title } = Typography;

const QATab = ({ ragName, defaultModel }) => {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [progressMessages, setProgressMessages] = useState([]);
  const [selectedModel, setSelectedModel] = useState(defaultModel);
  const [availableModels, setAvailableModels] = useState([]);
  const [contextSize, setContextSize] = useState(5);
  const messagesEndRef = useRef(null);
  const [copiedStates, setCopiedStates] = useState({});
  const currentStreamController = useRef(null);
  const chatContainerRef = useRef(null);

  // Fonction pour nettoyer les balises <think> des réponses
  const cleanThinkTags = (text) => {
    if (!text) return "";
    const thinkRegex = /<think>[\s\S]*?<\/think>/g;
    return text.replace(thinkRegex, "").trim();
  };

  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages, progressMessages]);

  useEffect(() => {
    setSelectedModel(defaultModel);
  }, [defaultModel]);

  useEffect(() => {
    const fetchModels = async () => {
      try {
        const models = await ragService.getAvailableModels();
        setAvailableModels(models || []);
      } catch (error) {
        console.error('Failed to fetch models:', error);
        message.error('Failed to load available LLM models.');
      }
    };
    fetchModels();
    return () => {
      if (currentStreamController.current) {
        currentStreamController.current();
        currentStreamController.current = null;
      }
    };
  }, []);

  const handleCopy = (text, id) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopiedStates(prev => ({ ...prev, [id]: true }));
      setTimeout(() => setCopiedStates(prev => ({ ...prev, [id]: false })), 2000);
    }).catch(err => {
      console.error('Failed to copy text: ', err);
      message.error('Failed to copy text.');
    });
  };

  // Composant personnalisé pour remplacer les citations de source dans Markdown
  const SourceReference = ({ children }) => {
    // On vérifie si le texte correspond à une citation de source
    const text = String(children);
    const sourceMatch = text.match(/\(Source: ([^)]+)\)/);
    
    if (sourceMatch) {
      const sourceName = sourceMatch[1];
      return (
        <span className="source-reference">
          <LinkOutlined style={{ marginRight: 4 }} />
          {sourceName}
        </span>
      );
    }
    
    return <span>{text}</span>;
  };

  // Composant personnalisé pour le rendu des éléments de code
  const CodeBlock = ({ node, inline, className, children, ...props }) => {
    const match = /language-(\w+)/.exec(className || '');
    const language = match ? match[1] : 'plaintext';
    
    if (!inline) {
      return (
        <div className="code-content">
          <div className="flex justify-between items-center mb-1 text-neutral-400 text-xs">
            <span>{language}</span>
            <Button 
              type="text" 
              size="small" 
              icon={<CopyOutlined />} 
              onClick={() => handleCopy(children, `code-${language}-${Date.now()}`)}
              className="text-neutral-400 hover:text-white"
            />
          </div>
          <pre className={className} {...props}>
            <code>{children}</code>
          </pre>
        </div>
      );
    }
    
    return (
      <code className={className} {...props}>
        {children}
      </code>
    );
  };

  const handleSubmit = async () => {
    if (!input.trim()) return;
    if (loading && currentStreamController.current) {
      console.log("Cancelling previous stream...");
      currentStreamController.current();
      currentStreamController.current = null;
      setLoading(false);
    }

    const question = input.trim();
    const userMessage = { role: 'user', content: question, key: `user-${Date.now()}` };    
    const initialAssistantMessage = { role: 'assistant', content: '', isError: false, key: `assistant-${Date.now()}` };
    
    setMessages(prevMessages => [...prevMessages, userMessage, initialAssistantMessage]);
    setInput('');
    setProgressMessages([]);
    setLoading(true);

    let accumulatedAnswer = "";

    currentStreamController.current = ragService.queryRagStream(
      {
        rag_name: ragName,
        prompt: question,
        context_size: contextSize,
        model: selectedModel !== defaultModel ? selectedModel : undefined,
      },
      {
        onProgress: (progressContent) => {
          setProgressMessages(prev => [...prev, progressContent]);
        },
        onAnswerChunk: (answerChunk) => {
          const cleanedChunk = cleanThinkTags(answerChunk);
          if (cleanedChunk) {
            accumulatedAnswer += cleanedChunk;
            setMessages(prevMessages =>
              prevMessages.map(msg =>
                msg.key === initialAssistantMessage.key
                  ? { ...msg, content: accumulatedAnswer }
                  : msg
              )
            );
          }
        },
        onError: (errorContent) => {
          console.error('Erreur stream RAG:', errorContent);
          setMessages(prevMessages =>
            prevMessages.map(msg =>
              msg.key === initialAssistantMessage.key
                ? { ...msg, content: `Erreur: ${errorContent}`, isError: true }
                : msg
            )
          );
          setLoading(false);
          setProgressMessages([]);
          currentStreamController.current = null;
        },
        onDone: () => {
          setLoading(false);
          setProgressMessages([]);
          
          const finalCleanedAnswer = cleanThinkTags(accumulatedAnswer);
          
          setMessages(prevMessages =>
            prevMessages.map(msg => {
              if (msg.key === initialAssistantMessage.key) {
                if (!msg.isError && finalCleanedAnswer.trim() === "") {
                  return { 
                    ...msg, 
                    content: "No response received or empty response. Check that the RAG contains documents relevant to your question.", 
                    isError: true 
                  };
                }
                return { ...msg, content: finalCleanedAnswer.trim() };
              }
              return msg;
            })
          );
          currentStreamController.current = null;
          accumulatedAnswer = "";
          
          setTimeout(() => {
            setMessages(prevMessages => [...prevMessages]);
          }, 100);
        },
      }
    );
  };

  const handleClearChat = () => {
    setMessages([]);
    setProgressMessages([]);
    if (loading && currentStreamController.current) {
      currentStreamController.current();
      currentStreamController.current = null;
      setLoading(false);
    }
  };

  // Composants personnalisés pour le rendu Markdown
  const markdownComponents = {
    root: ({ children }) => <div className="markdown-content">{children}</div>,
    code: CodeBlock,
    // Traiter les références sources dans les paragraphes
    p: ({ children }) => {
      if (typeof children === 'string' && children.includes('(Source:')) {
        // Pour un paragraphe contenant des références de source, on le divise et on traite chaque partie
        const parts = [];
        let lastIndex = 0;
        const sourceRegex = /\(Source: ([^)]+)\)/g;
        let match;
        let key = 0;
        
        const text = children;
        while ((match = sourceRegex.exec(text)) !== null) {
          if (match.index > lastIndex) {
            parts.push(<span key={`text-${key++}`}>{text.substring(lastIndex, match.index)}</span>);
          }
          
          parts.push(
            <span className="source-reference" key={`source-${key++}`}>
              <LinkOutlined style={{ marginRight: 4 }} />
              {match[1]}
            </span>
          );
          
          lastIndex = match.index + match[0].length;
        }
        
        if (lastIndex < text.length) {
          parts.push(<span key={`text-end-${key++}`}>{text.substring(lastIndex)}</span>);
        }
        
        return <p>{parts}</p>;
      }
      return <p>{children}</p>;
    },
    // Personnaliser d'autres éléments si nécessaire
    h1: ({ children }) => <h1 style={{ fontWeight: 600, fontSize: '1.8em', marginTop: '0.5em', marginBottom: '0.5em' }}>{children}</h1>,
    h2: ({ children }) => <h2 style={{ fontWeight: 600, fontSize: '1.5em', marginTop: '0.5em', marginBottom: '0.5em' }}>{children}</h2>,
    h3: ({ children }) => <h3 style={{ fontWeight: 600, fontSize: '1.3em', marginTop: '0.5em', marginBottom: '0.5em' }}>{children}</h3>,
    h4: ({ children }) => <h4 style={{ fontWeight: 600, fontSize: '1.2em', marginTop: '0.5em', marginBottom: '0.5em' }}>{children}</h4>,
    ul: ({ children }) => <ul style={{ paddingLeft: '1.5em', marginBottom: '1em' }}>{children}</ul>,
    ol: ({ children }) => <ol style={{ paddingLeft: '1.5em', marginBottom: '1em' }}>{children}</ol>,
    li: ({ children }) => <li style={{ marginBottom: '0.5em' }}>{children}</li>,
    blockquote: ({ children }) => (
      <blockquote style={{ 
        borderLeft: '3px solid var(--primary-300)', 
        paddingLeft: '1em', 
        margin: '1em 0', 
        color: 'var(--neutral-700)',
        backgroundColor: 'var(--neutral-100)',
        padding: '0.5em 1em',
        borderRadius: '0.25em'
      }}>
        {children}
      </blockquote>
    ),
    hr: () => <hr style={{ margin: '1em 0', border: 'none', borderTop: '1px solid var(--neutral-300)' }} />,
    strong: ({ children }) => <strong style={{ fontWeight: 600 }}>{children}</strong>,
    em: ({ children }) => <em style={{ fontStyle: 'italic' }}>{children}</em>,
    a: ({ href, children }) => <a href={href} target="_blank" rel="noopener noreferrer" style={{ color: 'var(--primary-600)', textDecoration: 'underline' }}>{children}</a>,
  };

  const renderMessageContent = (message) => {
    if (message.isError) {
      return (
        <div className="p-3 bg-red-50 text-red-800 rounded-lg">
          <InfoCircleOutlined style={{ marginRight: 8 }} />
          {message.content}
        </div>
      );
    }
    
    if (message.role === 'user') {
      return (
        <div className="whitespace-pre-wrap">
          {message.content}
        </div>
      );
    }
    
    return (
      <div className="relative">
        <div className="markdown-body">
          <ReactMarkdown 
            components={markdownComponents}
            remarkPlugins={[]}
            rehypePlugins={[rehypeHighlight]}
          >
            {message.content}
          </ReactMarkdown>
        </div>
        {message.content && (
          <Tooltip title={copiedStates[message.key] ? 'Copied!' : 'Copy response'}>
            <Button 
              type="text"
              icon={copiedStates[message.key] ? <CheckOutlined /> : <CopyOutlined />}
              onClick={() => handleCopy(message.content, message.key)}
              size="small"
              className="absolute top-0 right-0 opacity-70 hover:opacity-100"
            />
          </Tooltip>
        )}
      </div>
    );
  };

  return (
    <div className="flex flex-col h-[calc(100vh-300px)]">
      <div className="flex justify-between items-center mb-4">
        <div className="flex items-center gap-4">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <SettingOutlined style={{ color: 'var(--primary-700)' }} />
              <Text strong>Search Parameters</Text>
            </div>
            <div className="flex items-center gap-4">
              <Tooltip title="Select the LLM model to use for responses">
                <Select
                  value={selectedModel}
                  style={{ width: 200 }}
                  onChange={(value) => setSelectedModel(value)}
                  loading={availableModels.length === 0 && !ragName}
                  suffixIcon={<RobotOutlined />}
                >
                  {defaultModel && (
                    <Option key={defaultModel} value={defaultModel}>
                      <div className="flex items-center gap-2">
                        <RobotOutlined />
                        <span>{defaultModel}</span>
                        <Tag color="var(--primary-100)" style={{ color: 'var(--primary-700)' }}>default</Tag>
                      </div>
                    </Option>
                  )}
                  {availableModels.filter(m => m !== defaultModel).map(model => (
                    <Option key={model} value={model}>
                      <div className="flex items-center gap-2">
                        <RobotOutlined />
                        <span>{model}</span>
                      </div>
                    </Option>
                  ))}
                </Select>
              </Tooltip>
              
              <Tooltip title="Number of context chunks to provide to the LLM">
                <div className="flex items-center gap-2">
                  <span>Context:</span>
                  <InputNumber 
                    min={1} 
                    max={20} 
                    value={contextSize} 
                    onChange={(value) => setContextSize(value || 1)} 
                    className="w-16"
                  />
                </div>
              </Tooltip>
            </div>
          </div>
        </div>
        
        <Button 
          danger 
          type="primary" 
          icon={<ClearOutlined />} 
          onClick={handleClearChat}
          disabled={messages.length === 0}
        >
          Clear conversation
        </Button>
      </div>
      
      <div className="chat-container relative flex-grow mb-4 p-4" ref={chatContainerRef}>
        {messages.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-center text-neutral-500">
            <QuestionCircleOutlined style={{ fontSize: 48, marginBottom: 16 }} />
            <Title level={4}>Ask questions about your documents</Title>
            <Text type="secondary">
              The RAG system will search for the most relevant information in your documents
              <br />and generate a response based on this content.
            </Text>
          </div>
        ) : (
          <div className="space-y-6">
            {messages.map((item, index) => (
              <div 
                key={item.key || index} 
                className={`flex ${item.role === 'user' ? 'justify-end' : 'justify-start'}`}
              >
                <div 
                  className={`max-w-[85%] message ${item.role === 'user' ? 'user-message' : 'assistant-message'}`}
                >
                  {renderMessageContent(item)}
                </div>
              </div>
            ))}
            
            {loading && progressMessages.length > 0 && (
              <div className="p-3 bg-neutral-100 rounded-lg mt-4 mb-2">
                <div className="flex items-center gap-2 mb-2">
                  <Spin size="small" />
                  <Text strong>Searching...</Text>
                </div>
                <div className="pl-6 text-xs font-mono text-neutral-600">
                  {progressMessages.map((msg, index) => (
                    <div key={`prog-${index}`} className="mb-1">
                      {msg}
                    </div>
                  ))}
                </div>
              </div>
            )}
            <div ref={messagesEndRef} />
          </div>
        )}
      </div>
      
      <div className="flex gap-2 items-end">
        <TextArea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Ask questions about your documents..."
          onKeyPress={(e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
              e.preventDefault();
              handleSubmit();
            }
          }}
          autoSize={{ minRows: 1, maxRows: 5 }}
          disabled={loading}
          className="flex-grow rounded-lg shadow-sm"
          autoFocus
        />
        <Button 
          type="primary" 
          icon={<SendOutlined />} 
          onClick={handleSubmit} 
          loading={loading}
          size="large"
          className="mb-[2px] shadow-sm"
          disabled={!input.trim()}
        >
          Send
        </Button>
      </div>
    </div>
  );
};

export default QATab; 