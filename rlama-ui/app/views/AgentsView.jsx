import React, { useState, useRef, useEffect } from 'react';
import { 
  Card, 
  Input, 
  Button, 
  Select, 
  Switch, 
  Typography, 
  Space, 
  Spin, 
  Alert,
  Divider,
  Badge,
  List,
  Avatar,
  Tag
} from 'antd';
import { 
  SendOutlined, 
  RobotOutlined, 
  ThunderboltOutlined, 
  EyeOutlined, 
  LoadingOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  ClockCircleOutlined,
  DatabaseOutlined,
  SearchOutlined
} from '@ant-design/icons';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { agentService, ragService } from '../services/api';
import api from '../services/api';
import './AgentsView.css';

const { TextArea } = Input;
const { Option } = Select;
const { Title, Text, Paragraph } = Typography;

const AgentsView = () => {
  // State for form
  const [query, setQuery] = useState('');
  const [model, setModel] = useState('');
  const [ragName, setRagName] = useState('');
  const [webSearch, setWebSearch] = useState(false);
  const [availableModels, setAvailableModels] = useState([]);
  const [availableRags, setAvailableRags] = useState([]);
  
  // State for execution
  const [isExecuting, setIsExecuting] = useState(false);
  const [tasks, setTasks] = useState([]);
  const [currentStep, setCurrentStep] = useState('');
  const [currentStepType, setCurrentStepType] = useState(''); // 'thinking', 'searching', 'analyzing', etc.
  const [response, setResponse] = useState('');
  const [progress, setProgress] = useState([]);
  const [error, setError] = useState(null);
  const [executionSteps, setExecutionSteps] = useState([]); // New state for structured steps
  const [debugMessages, setDebugMessages] = useState([]); // Debug messages from backend
  const [showDebug, setShowDebug] = useState(false); // Toggle debug panel
  
  // Refs
  const responseRef = useRef(null);
  const streamCleanupRef = useRef(null);
  const progressLogRef = useRef(null);

  // Track agent execution with proper timing
  const [agentSteps, setAgentSteps] = useState([]);
  const [detailedActions, setDetailedActions] = useState([]); // For detailed actions like "Read progress.go"
  const [executionStartTime, setExecutionStartTime] = useState(null);
  const [completedTasks, setCompletedTasks] = useState(0);
  const [totalTasks, setTotalTasks] = useState(0);
  const [isProcessingResponse, setIsProcessingResponse] = useState(false);
  const [forceRefresh, setForceRefresh] = useState(0); // Force UI updates

  // Load initial data
  useEffect(() => {
    loadInitialData();
  }, []);

  // Auto scroll progress log
  useEffect(() => {
    if (progressLogRef.current) {
      progressLogRef.current.scrollTop = progressLogRef.current.scrollHeight;
    }
  }, [progress]);

  // Helper function to clean text from debug info but preserve useful content
  const cleanText = (text) => {
    if (typeof text !== 'string') return '';
    
    let cleaned = text;
    
    // Remove only obvious debug patterns, keep most content
    cleaned = cleaned.replace(/\[[0-9]+[A-Z]\]/g, ''); // [2A], [1B], etc.
    cleaned = cleaned.replace(/\[K\]/g, ''); // [K]
    cleaned = cleaned.replace(/\([0-9.]+s\)/g, ''); // (0.8s)
    cleaned = cleaned.replace(/\[[0-9]+B\]/g, ''); // [2B]
    
    // Remove SSE data prefixes but keep content
    cleaned = cleaned.replace(/^data:\s*/, '');
    cleaned = cleaned.replace(/^event:\s*\w+\s*/, '');
    
    // Normalize excessive whitespace but preserve structure
    cleaned = cleaned.replace(/\s+/g, ' ');
    cleaned = cleaned.replace(/\n\s*\n/g, '\n');
    
    return cleaned.trim();
  };

  // Enhanced function to parse real agent outputs - more aggressive parsing
  const parseAgentOutput = (message) => {
    const cleaned = cleanText(message);
    
    // Log EVERYTHING for debugging
    console.log('🔍 Raw message:', message);
    console.log('🧹 Cleaned message:', cleaned);
    
    // More flexible patterns - catch more cases
    
    // Starting patterns - multiple variations
    if (cleaned.toLowerCase().includes('starting') || 
        cleaned.toLowerCase().includes('agent') || 
        cleaned.includes('stream started')) {
      console.log('✅ Detected: Starting agent');
      return { 
        type: 'starting', 
        phase: 'init',
        message: '◦ Starting agent system',
        emoji: '◦'
      };
    }
    
    // Progress patterns - very broad
    if (cleaned.toLowerCase().includes('progress') || 
        cleaned.toLowerCase().includes('processing') ||
        cleaned.toLowerCase().includes('working')) {
      console.log('✅ Detected: Progress');
      return { 
        type: 'progress', 
        phase: 'processing',
        message: '∘ Processing request',
        emoji: '∘'
      };
    }
    
    // Analysis patterns
    if (cleaned.toLowerCase().includes('analyz') || 
        cleaned.toLowerCase().includes('decompos') ||
        cleaned.toLowerCase().includes('break') ||
        cleaned.toLowerCase().includes('complex')) {
      console.log('✅ Detected: Analysis/Decomposition');
      return { 
        type: 'decomposing', 
        phase: 'decomposition',
        message: '◈ Analyzing and decomposing query',
        emoji: '◈'
      };
    }
    
    // Task patterns - very flexible
    if (/task/i.test(cleaned) || /step/i.test(cleaned)) {
      console.log('✅ Detected: Task-related');
      
      // Try to extract task numbers
      const numberMatch = cleaned.match(/(\d+)/);
      if (numberMatch) {
        const taskNum = parseInt(numberMatch[1]);
        
        // Check if it's task completion
        if (cleaned.toLowerCase().includes('complet') || 
            cleaned.toLowerCase().includes('finish') ||
            cleaned.toLowerCase().includes('done') ||
            cleaned.toLowerCase().includes('terminé')) {
          return { 
            type: 'tasks_completed', 
            phase: 'completion',
            count: taskNum,
            message: `◉ ${taskNum} tasks completed`,
            emoji: '◉'
          };
        } else {
          // Task execution
          return { 
            type: 'task_execution', 
            phase: 'execution',
            taskNumber: taskNum,
            taskName: 'Executing task',
            message: `◦ Task ${taskNum}: In progress`,
            emoji: '◦'
          };
        }
      } else {
        // Generic task processing
        return { 
          type: 'task_execution', 
          phase: 'execution',
          message: '⚡ Processing tasks',
          emoji: '⚡'
        };
      }
    }
    
    // Search patterns
    if (cleaned.toLowerCase().includes('search') || 
        cleaned.toLowerCase().includes('query') ||
        cleaned.toLowerCase().includes('find') ||
        cleaned.toLowerCase().includes('look')) {
      console.log('✅ Detected: Search');
      return { 
        type: 'search', 
        phase: 'search',
        message: '◎ Searching for information',
        emoji: '◎'
      };
    }
    
    // Orchestration patterns
    if (cleaned.toLowerCase().includes('orchestr') || 
        cleaned.toLowerCase().includes('coordinat') ||
        cleaned.toLowerCase().includes('manag') ||
        cleaned.toLowerCase().includes('organiz')) {
      console.log('✅ Detected: Orchestration');
      return { 
        type: 'orchestration', 
        phase: 'coordination',
        message: '◈ Orchestrating execution',
        emoji: '◈'
      };
    }
    
    // Compilation patterns
    if (cleaned.toLowerCase().includes('compil') || 
        cleaned.toLowerCase().includes('collect') ||
        cleaned.toLowerCase().includes('gather') ||
        cleaned.toLowerCase().includes('assembl')) {
      console.log('✅ Detected: Compilation');
      return { 
        type: 'compilation', 
        phase: 'compilation',
        message: '◈ Compiling results',
        emoji: '◈'
      };
    }
    
    // Synthesis patterns
    if (cleaned.toLowerCase().includes('synthes') || 
        cleaned.toLowerCase().includes('generat') ||
        cleaned.toLowerCase().includes('final') ||
        cleaned.toLowerCase().includes('response') ||
        cleaned.toLowerCase().includes('answer')) {
      console.log('✅ Detected: Synthesis');
      return { 
        type: 'synthesis', 
        phase: 'synthesis',
        message: '◌ Synthesizing response',
        emoji: '◌'
      };
    }
    
    // Completion patterns
    if (cleaned.toLowerCase().includes('complet') || 
        cleaned.toLowerCase().includes('finish') ||
        cleaned.toLowerCase().includes('done') ||
        cleaned.toLowerCase().includes('success') ||
        cleaned.toLowerCase().includes('terminé')) {
      console.log('✅ Detected: Completion');
      return { 
        type: 'tasks_completed', 
        phase: 'completion',
        message: '◉ Tasks completed',
        emoji: '◉'
      };
    }
    
    // Detect detailed actions (Read files, Thinking, etc.)
    if (cleaned.toLowerCase().includes('read') && 
        (cleaned.includes('.go') || cleaned.includes('.py') || cleaned.includes('.js') || cleaned.includes('.md'))) {
      console.log('✅ Detected: File reading');
      return { 
        type: 'detailed_action', 
        actionType: 'read',
        phase: 'file_operation',
        message: `◎ ${cleaned}`,
        emoji: '◎'
      };
    }
    
    if (cleaned.toLowerCase().includes('thought') || cleaned.toLowerCase().includes('thinking')) {
      const timeMatch = cleaned.match(/(\d+)\s*seconds?/);
      const duration = timeMatch ? timeMatch[1] : '';
      console.log('✅ Detected: Thinking');
      return { 
        type: 'detailed_action', 
        actionType: 'thinking',
        phase: 'reasoning',
        message: `◌ Thought for ${duration} seconds`,
        emoji: '◌',
        duration: duration
      };
    }
    
    if (cleaned.toLowerCase().includes('executing') || cleaned.toLowerCase().includes('running')) {
      console.log('✅ Detected: Execution');
      return { 
        type: 'detailed_action', 
        actionType: 'execution',
        phase: 'execution',
        message: `⚡ ${cleaned}`,
        emoji: '⚡'
      };
    }
    
    // Fallback for any meaningful content - be very permissive
    if (cleaned.length > 10 && 
        !cleaned.toLowerCase().includes('sse') &&
        !cleaned.toLowerCase().includes('chunk') &&
        !cleaned.toLowerCase().includes('data:')) {
      console.log('✅ Detected: Generic activity');
      return { 
        type: 'processing', 
        phase: 'active',
        message: `∘ ${cleaned.substring(0, 50)}${cleaned.length > 50 ? '...' : ''}`,
        emoji: '∘'
      };
    }
    
    console.log('❌ No pattern matched, returning null');
    return null;
  };

  // Helper function to get typing text based on current step
  const getTypingText = () => {
    if (currentStep && currentStep !== 'Starting agent...') {
      return currentStep;
    }
    
    switch (currentStepType) {
      case 'init': return 'Starting agent...';
      case 'starting': return '🚀 Starting agent system...';
      case 'progress': return '⚙️ Processing request...';
      case 'decomposing': return '📋 Breaking down tasks...';
      case 'tasks_identified': return '📋 Tasks identified...';
      case 'task_execution': return '⚡ Executing task...';
      case 'orchestration': return '🎯 Orchestrating...';
      case 'search': return '🔍 Searching...';
      case 'compilation': return '📊 Compiling results...';
      case 'synthesis': return '🔄 Synthesizing response...';
      case 'processing': return '⚙️ Processing information...';
      case 'tasks_completed': return '✅ Tasks completed!';
      case 'completed': return '✅ Done!';
      default: return 'Processing...';
    }
  };

  const loadInitialData = async () => {
    try {
      // Load available models - using the ragService instead since agentService.getAvailableModels might fail
      try {
        const models = await ragService.getAvailableModels();
        setAvailableModels(models);
        if (models.length > 0) {
          setModel(models[0]);
        }
      } catch (modelError) {
        console.error('Error loading models:', modelError);
        // Try to get models through health check as fallback
        try {
          const modelsResponse = await api.get('/models');
          if (modelsResponse.data && modelsResponse.data.models) {
            setAvailableModels(modelsResponse.data.models);
            if (modelsResponse.data.models.length > 0) {
              setModel(modelsResponse.data.models[0]);
            }
          }
        } catch (fallbackError) {
          console.error('Fallback model loading also failed:', fallbackError);
        }
      }

      // Load available RAGs
      try {
        const rags = await ragService.getAllRags();
        const ragNames = rags.map(r => r.name);
        setAvailableRags(ragNames);
      } catch (ragError) {
        console.error('Error loading RAGs:', ragError);
        // Set empty array if loading fails
        setAvailableRags([]);
      }
    } catch (error) {
      console.error('Error loading initial data:', error);
    }
  };

  const handleExecute = async () => {
    if (!query.trim()) return;
    
    setIsExecuting(true);
    setError(null);
    setTasks([]);
    setProgress([]);
    setResponse('');
    setExecutionSteps([]);
    setAgentSteps([]);
    setDetailedActions([]);
    setCompletedTasks(0);
    setTotalTasks(0);
    setIsProcessingResponse(false);
    setCurrentStep('Starting agent...');
    setCurrentStepType('init');
    setDebugMessages([]);
    setExecutionStartTime(new Date());
    setForceRefresh(0); // Reset force refresh

    try {
      const requestData = {
        query: query.trim(),
        model: model || undefined,
        rag_name: ragName || undefined,
        web_search: webSearch,
        verbose: true
      };

      if (streamCleanupRef.current) {
        streamCleanupRef.current();
      }

      streamCleanupRef.current = agentService.runAgentStream(requestData, {
        onProgress: (progressData) => {
          const progressStr = typeof progressData === 'string' ? progressData : 
                            (progressData?.message || progressData?.content || JSON.stringify(progressData));
          
          console.log('🎯 Raw progress data:', progressStr);
          
          setDebugMessages(prev => [...prev, {
            id: Date.now() + Math.random(),
            message: progressStr,
            timestamp: new Date()
          }]);
          
          const agentInfo = parseAgentOutput(progressStr);
          
          if (agentInfo) {
            console.log('✅ Parsed agent step:', agentInfo);
            
            // Force UI refresh
            setForceRefresh(prev => prev + 1);
            
            // Always update current step display - force refresh
            setCurrentStep(agentInfo.message);
            setCurrentStepType(agentInfo.type);
            
            // Handle different types of agent steps with forced updates
            if (agentInfo.type === 'starting') {
              setCurrentStep('◦ Starting agent system');
              setCurrentStepType('starting');
            } else if (agentInfo.type === 'progress') {
              setCurrentStep('∘ Processing request');
              setCurrentStepType('progress');
            } else if (agentInfo.type === 'decomposing') {
              setCurrentStep('◈ Decomposing tasks');
              setCurrentStepType('decomposing');
            } else if (agentInfo.type === 'tasks_identified') {
              setTotalTasks(agentInfo.count);
              setCurrentStep(`◈ ${agentInfo.count} tasks identified`);
              setCurrentStepType('tasks_identified');
            } else if (agentInfo.type === 'task_execution') {
              setCurrentStep(`◦ Executing task ${agentInfo.taskNumber || ''}`);
              setCurrentStepType('task_execution');
              // Track task completion
              if (agentInfo.taskNumber) {
                setCompletedTasks(prev => Math.max(prev, agentInfo.taskNumber));
              }
            } else if (agentInfo.type === 'tasks_completed') {
              setCompletedTasks(agentInfo.count || completedTasks + 1);
              setCurrentStep(`◉ ${agentInfo.count || 'Tasks'} completed`);
            } else if (agentInfo.type === 'synthesis') {
              setIsProcessingResponse(true);
              setCurrentStep('◌ Synthesizing final response');
              setCurrentStepType('synthesis');
            } else if (agentInfo.type === 'processing') {
              setCurrentStep(agentInfo.message || '∘ Processing information');
              setCurrentStepType('processing');
            } else if (agentInfo.type === 'search') {
              setCurrentStep(agentInfo.message || '◎ Searching for information');
              setCurrentStepType('search');
            } else if (agentInfo.type === 'orchestration') {
              setCurrentStep(agentInfo.message || '◈ Orchestrating execution');
              setCurrentStepType('orchestration');
            } else if (agentInfo.type === 'compilation') {
              setCurrentStep(agentInfo.message || '◈ Compiling results');
              setCurrentStepType('compilation');
            }
            
            // Handle detailed actions vs main agent steps
            if (agentInfo.type === 'detailed_action') {
              // Add to detailed actions timeline
              setDetailedActions(prev => {
                const actionStartTime = new Date();
                const newAction = {
                  id: Date.now() + Math.random(),
                  ...agentInfo,
                  timestamp: actionStartTime
                };
                
                // Keep only last 10 detailed actions to avoid clutter
                const updatedActions = [...prev, newAction];
                return updatedActions.slice(-10);
              });
            } else {
              // Add to main agent steps timeline
              setAgentSteps(prev => {
                // Less strict duplicate checking - only check message similarity
                const isDuplicate = prev.some(step => 
                  step.message === agentInfo.message &&
                  Math.abs(step.timestamp.getTime() - new Date().getTime()) < 2000 // within 2 seconds
                );
                
                if (!isDuplicate) {
                  const stepStartTime = new Date();
                  const newStep = {
                    id: Date.now() + Math.random(),
                    ...agentInfo,
                    timestamp: stepStartTime,
                    status: agentInfo.type === 'tasks_completed' ? 'completed' : 'running',
                    duration: null,
                    startTime: stepStartTime
                  };
                  
                  // Mark previous steps of same type as completed (but allow multiple of same type)
                  const updatedSteps = prev.map(step => 
                    step.type === agentInfo.type && step.status === 'running' && step.id !== newStep.id
                      ? { ...step, status: 'completed', duration: Math.round((stepStartTime.getTime() - step.startTime.getTime()) / 1000) }
                      : step
                  );
                  
                  return [...updatedSteps, newStep];
                }
                return prev;
              });
            }
          } else {
            console.log('❌ Failed to parse, showing fallback');
            // FALLBACK: Always show SOMETHING even if parsing fails
            const cleaned = cleanText(progressStr);
            if (cleaned.length > 5) {
              // Force UI refresh even for unparsed content
              setForceRefresh(prev => prev + 1);
              
              // Show the raw content as fallback
              const fallbackMessage = `∘ ${cleaned.substring(0, 80)}${cleaned.length > 80 ? '...' : ''}`;
              setCurrentStep(fallbackMessage);
              setCurrentStepType('processing');
              
              // Add generic step to timeline too
              setAgentSteps(prev => {
                const stepStartTime = new Date();
                const fallbackStep = {
                  id: Date.now() + Math.random(),
                  type: 'processing',
                  phase: 'active',
                  message: fallbackMessage,
                  emoji: '∘',
                  timestamp: stepStartTime,
                  status: 'running',
                  duration: null,
                  startTime: stepStartTime
                };
                
                                 // Avoid too many generic steps
                 const recentGeneric = prev.filter(step => 
                   step.type === 'processing' && 
                   Math.abs(step.timestamp.getTime() - stepStartTime.getTime()) < 5000
                 );
                 
                 if (recentGeneric.length < 2) {
                   return [...prev, fallbackStep];
                 }
                 return prev;
               });
            }
          }
        },
        onTaskUpdate: (taskData) => {
          console.log('Task update:', taskData);
          
          // Update execution steps with task completion
          if (taskData.status === 'completed' || taskData.status === 'failed') {
            setExecutionSteps(prev => prev.map(step => {
              // Match by description or task ID
              if (step.status === 'running' && 
                  (step.description?.includes(taskData.description) || 
                   taskData.description?.includes(step.description) ||
                   step.id === taskData.task_id)) {
                return {
                  ...step,
                  status: taskData.status,
                  result: taskData.result,
                  tool: taskData.tool,
                  error: taskData.error
                };
              }
              return step;
            }));
          }
          
          // Update tasks list if needed
          setTasks(prev => {
            const existingIndex = prev.findIndex(t => t.task_id === taskData.task_id);
            if (existingIndex >= 0) {
              const updated = [...prev];
              updated[existingIndex] = { ...updated[existingIndex], ...taskData };
              return updated;
            } else {
              return [...prev, taskData];
            }
          });
        },
        onAnswerChunk: (chunk) => {
          const chunkStr = typeof chunk === 'string' ? chunk : JSON.stringify(chunk);
          const cleanedChunk = cleanText(chunkStr);
          
          console.log('🎯 Answer chunk received:', chunkStr);
          console.log('🧹 Cleaned chunk:', cleanedChunk);
          
          // Check if this chunk contains step information (emojis/progress)
          const stepEmojis = ['🚀', '📋', '⚡', '🔍', '🎯', '📊', '🔄', '✅'];
          const hasStepEmoji = stepEmojis.some(emoji => chunkStr.includes(emoji) || cleanedChunk.includes(emoji));
          
          if (hasStepEmoji) {
            console.log('✅ Chunk contains step emoji, treating as progress');
            // Treat this as progress, not final response
            const agentInfo = parseAgentOutput(cleanedChunk);
            
            if (agentInfo) {
              console.log('✅ Parsed step from chunk:', agentInfo);
              
              // Force UI refresh
              setForceRefresh(prev => prev + 1);
              
              // Update current step
              setCurrentStep(agentInfo.message);
              setCurrentStepType(agentInfo.type);
              
              // Add to timeline
              setAgentSteps(prev => {
                const isDuplicate = prev.some(step => 
                  step.message === agentInfo.message &&
                  Math.abs(step.timestamp.getTime() - new Date().getTime()) < 3000
                );
                
                if (!isDuplicate) {
                  const stepStartTime = new Date();
                  const newStep = {
                    id: Date.now() + Math.random(),
                    ...agentInfo,
                    timestamp: stepStartTime,
                    status: agentInfo.type === 'tasks_completed' ? 'completed' : 'running',
                    duration: null,
                    startTime: stepStartTime
                  };
                  
                  // Mark previous steps as completed if needed
                  const updatedSteps = prev.map(step => 
                    step.type === agentInfo.type && step.status === 'running' && step.id !== newStep.id
                      ? { ...step, status: 'completed', duration: Math.round((stepStartTime.getTime() - step.startTime.getTime()) / 1000) }
                      : step
                  );
                  
                  return [...updatedSteps, newStep];
                }
                return prev;
              });
            }
            return; // Don't process as final response
          }
          
          // Only process as final response if it doesn't contain step indicators
          if (cleanedChunk && 
              cleanedChunk.length > 30 && 
              !hasStepEmoji &&
              !cleanedChunk.includes('Auto-detected') &&
              !cleanedChunk.includes('Starting agent') &&
              !cleanedChunk.includes('Analyzing complex') &&
              !cleanedChunk.includes('Decomposed into') &&
              !cleanedChunk.includes('Task ') &&
              !cleanedChunk.includes('Orchestration') &&
              !cleanedChunk.includes('Search for') &&
              !cleanedChunk.includes('Compile and present') &&
              !cleanedChunk.includes('Synthesizing') &&
              !cleanedChunk.includes('terminées') &&
              !cleanedChunk.includes('RAG system') &&
              !cleanedChunk.includes('model ') &&
              !cleanedChunk.includes('Loading')) {
            
            console.log('✅ Adding to final response:', cleanedChunk);
            
            setResponse(prev => {
              // Avoid duplicate content
              if (prev.includes(cleanedChunk.substring(0, 50))) {
                return prev;
              }
              
              // Clean join with proper spacing
              const needsSpace = prev && !prev.endsWith(' ') && !prev.endsWith('\n') && !cleanedChunk.startsWith(' ');
              return prev + (needsSpace ? ' ' : '') + cleanedChunk;
            });
            
            if (responseRef.current) {
              responseRef.current.scrollTop = responseRef.current.scrollHeight;
            }
          }
        },
        onError: (errorMsg) => {
          // Ensure error is a string
          const errorStr = typeof errorMsg === 'string' ? errorMsg : 
                          (errorMsg?.message || errorMsg?.detail || JSON.stringify(errorMsg));
          setError(errorStr);
          setIsExecuting(false);
          setCurrentStep('');
          setCurrentStepType('');
        },
        onDone: () => {
          setIsExecuting(false);
          setIsProcessingResponse(false);
          
          // Calculate total execution time in seconds (not timestamp)
          const endTime = new Date();
          const totalDuration = executionStartTime 
            ? Math.round((endTime.getTime() - executionStartTime.getTime()) / 1000)
            : 0;
          
          // Mark all running steps as completed with individual durations
          setAgentSteps(prev => prev.map(step => {
            if (step.status === 'running') {
              const stepDuration = step.startTime 
                ? Math.round((endTime.getTime() - step.startTime.getTime()) / 1000)
                : totalDuration;
              return {
                ...step,
                status: 'completed',
                duration: stepDuration
              };
            }
            return step;
          }));
          
          // Add final completion step if not already present
          setAgentSteps(prev => {
            const hasCompletionStep = prev.some(step => step.type === 'completion');
            if (!hasCompletionStep) {
              return [...prev, {
                id: Date.now(),
                type: 'completion',
                phase: 'done',
                message: `◉ All tasks completed in ${totalDuration}s`,
                emoji: '◉',
                timestamp: endTime,
                status: 'completed',
                duration: totalDuration,
                startTime: executionStartTime || endTime
              }];
            }
            return prev;
          });
          
          setCurrentStep(`◉ Completed in ${totalDuration} seconds`);
          setCurrentStepType('completed');
          
          setTimeout(() => {
            setCurrentStep('');
            setCurrentStepType('');
          }, 3000);
        }
      });

    } catch (error) {
      const errorStr = error?.message || error?.detail || 'An error occurred';
      setError(errorStr);
      setIsExecuting(false);
    }
  };

  const getStepIcon = (type, status) => {
    const iconStyle = { fontSize: '18px' };
    
    if (status === 'completed') {
      return <CheckCircleOutlined style={{ ...iconStyle, color: 'var(--status-success)' }} />;
    } else if (status === 'failed') {
      return <ExclamationCircleOutlined style={{ ...iconStyle, color: 'var(--status-error)' }} />;
    }
    
    // Icons based on step type
    switch (type) {
      case 'init':
        return <RobotOutlined style={{ ...iconStyle, color: 'var(--accent-primary)' }} />;
      case 'analyzing':
        return <EyeOutlined style={{ ...iconStyle, color: 'var(--status-info)' }} />;
      case 'searching':
        return <SearchOutlined style={{ ...iconStyle, color: 'var(--status-warning)' }} />;
      case 'orchestrating':
        return <ThunderboltOutlined style={{ ...iconStyle, color: 'var(--accent-primary)' }} />;
      case 'compiling':
        return <DatabaseOutlined style={{ ...iconStyle, color: 'var(--status-success)' }} />;
      case 'synthesizing':
        return <LoadingOutlined spin style={{ ...iconStyle, color: 'var(--accent-primary)' }} />;
      default:
        return <LoadingOutlined spin style={{ ...iconStyle, color: 'var(--text-secondary)' }} />;
    }
  };

  const getCurrentStepTitle = () => {
    switch (currentStepType) {
      case 'init': return 'Initializing';
      case 'starting': return 'Starting';
      case 'progress': return 'Processing';
      case 'decomposing': return 'Decomposing Tasks';
      case 'tasks_identified': return 'Tasks Identified';
      case 'task_execution': return 'Executing Tasks';
      case 'orchestration': return 'Orchestrating';
      case 'search': return 'Searching';
      case 'compilation': return 'Compiling Results';
      case 'synthesis': return 'Synthesizing Response';
      case 'processing': return 'Processing';
      case 'tasks_completed': return 'Tasks Completed';
      case 'completed': return 'Completed';
      default: return 'Processing';
    }
  };

  const getTaskIcon = (status) => {
    switch (status) {
      case 'completed': return <CheckCircleOutlined style={{ color: 'var(--status-success)' }} />;
      case 'failed': return <ExclamationCircleOutlined style={{ color: 'var(--status-error)' }} />;
      case 'running': return <LoadingOutlined spin style={{ color: 'var(--status-info)' }} />;
      default: return <ClockCircleOutlined style={{ color: 'var(--status-warning)' }} />;
    }
  };

  const getTaskColor = (status) => {
    switch (status) {
      case 'completed': return 'success';
      case 'failed': return 'error';
      case 'running': return 'processing';
      default: return 'default';
    }
  };

  return (
    <div className="agents-container">
      {/* Header */}
      <div className="agents-header">
        <div className="header-content">
          <RobotOutlined className="header-icon" />
          <Title level={2} style={{ margin: 0 }}>AI Agents</Title>
          <Text type="secondary">Intelligent task orchestration with real-time feedback</Text>
        </div>
      </div>

      <div className="agents-layout">
        {/* Left Panel - Input Form */}
        <Card className="input-panel" title="Agent Configuration">
          <Space direction="vertical" style={{ width: '100%' }} size="large">
            <div>
              <Text strong>Query</Text>
              <TextArea
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Describe what you want the agent to do..."
                rows={4}
                disabled={isExecuting}
              />
            </div>

            <div>
              <Text strong>Model</Text>
              <Select
                value={model}
                onChange={setModel}
                style={{ width: '100%' }}
                placeholder="Select a model"
                disabled={isExecuting}
              >
                {availableModels.map(m => (
                  <Option key={m} value={m}>{m}</Option>
                ))}
              </Select>
            </div>

            <div>
              <Text strong>RAG System (Optional)</Text>
              <Select
                value={ragName}
                onChange={setRagName}
                style={{ width: '100%' }}
                placeholder="Select a RAG system"
                allowClear
                disabled={isExecuting}
              >
                {availableRags.map(rag => (
                  <Option key={rag} value={rag}>
                    <DatabaseOutlined /> {rag}
                  </Option>
                ))}
              </Select>
            </div>

            <div className="switch-item">
              <Text strong>Web Search</Text>
              <Switch
                checked={webSearch}
                onChange={setWebSearch}
                disabled={isExecuting}
                checkedChildren={<SearchOutlined />}
                unCheckedChildren={<SearchOutlined />}
              />
            </div>

            <Button
              type="primary"
              size="large"
              icon={isExecuting ? <LoadingOutlined /> : <SendOutlined />}
              onClick={handleExecute}
              disabled={!query.trim() || isExecuting}
              block
              className="execute-button"
            >
              {isExecuting ? 'Executing...' : 'Execute Agent'}
            </Button>
          </Space>
        </Card>

        {/* Right Panel - Dynamic Process Display */}
        <div className="process-panel">
          {/* Current Step Display */}
          {(isExecuting || currentStep) && (
            <Card className="current-step-card" key={`current-step-${forceRefresh}`}>
              <div className="current-step">
                <div className="step-header">
                  <div className="step-indicator">
                    {getStepIcon(currentStepType, isExecuting ? 'running' : 'completed')}
                  </div>
                  <div className="step-content">
                    <Text strong className="step-title">
                      {getCurrentStepTitle()}
                    </Text>
                    <Text type="secondary" className="step-description">
                      {currentStep || 'Ready to execute'}
                    </Text>
                    {isExecuting && (
                      <div className="progress-bar">
                        <div className="progress-fill"></div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </Card>
          )}

          {/* Agent Steps Timeline - Exact Match to Design */}
          {agentSteps.length > 0 && (
            <Card title="Agent Progress" className="compact-agent-steps-card" key={`agent-steps-${forceRefresh}`}>
              <div className="compact-agent-timeline">
                {agentSteps.map((step, index) => (
                  <div key={`${step.id}-${forceRefresh}`} className={`compact-agent-step ${step.status}`}>
                    <div className="compact-step-icon">
                      <span className="step-emoji">{step.emoji}</span>
                      {step.status === 'completed' && (
                        <CheckCircleOutlined className="step-check" />
                      )}
                    </div>
                    <div className="compact-step-message">
                                            <Text className="step-text">
                        {step.message.replace(/^[🚀📋⚡🔍🎯📊🔄✅⚙️◦◈◎◌◉∘]\s*/, '')}
                      </Text>
                    </div>
                    <div className="compact-step-time">
                      <Text type="secondary">
                        {step.timestamp.toLocaleTimeString([], {
                          hour: 'numeric',
                          minute: '2-digit',
                          second: '2-digit',
                          hour12: true
                        })}
                      </Text>
                    </div>
                    {step.status === 'running' && (
                      <div className="compact-step-indicator">
                        <LoadingOutlined spin style={{ fontSize: '12px', color: 'var(--accent-primary)' }} />
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </Card>
          )}

          {/* Detailed Actions Timeline - Like in the image */}
          {detailedActions.length > 0 && (
            <Card className="detailed-actions-card" size="small">
              <div className="detailed-actions-list">
                {detailedActions.map((action, index) => (
                  <div key={`${action.id}-${forceRefresh}`} className="detailed-action-item">
                    <span className="action-emoji">{action.emoji}</span>
                    <span className="action-text">{action.message.replace(/^[👁🧠⚡◎◌◦]\s*/, '')}</span>
                  </div>
                ))}
              </div>
            </Card>
          )}

          {/* Task Progress Summary */}
          {totalTasks > 0 && (
            <Card className="task-summary-card" key={`task-summary-${forceRefresh}`}>
              <div className="task-summary">
                <Text strong>Task Progress: </Text>
                <Text>{completedTasks} / {totalTasks} tasks completed</Text>
                <div className="task-progress-bar">
                  <div 
                    className="task-progress-fill" 
                    style={{ width: `${(completedTasks / totalTasks) * 100}%` }}
                  ></div>
                </div>
              </div>
            </Card>
          )}

          {/* Response Display with Markdown */}
          {response && (
            <Card title="Final Response" className="response-card">
              <div 
                ref={responseRef}
                className="response-content"
              >
                <ReactMarkdown 
                  className="response-markdown"
                  remarkPlugins={[remarkGfm]}
                >
                  {response}
                </ReactMarkdown>
              </div>
            </Card>
          )}

          {/* Processing Indicator */}
          {isProcessingResponse && !response && (
            <Card className="processing-card">
              <div className="processing-indicator">
                <LoadingOutlined spin style={{ fontSize: '24px', color: 'var(--accent-primary)' }} />
                <Text style={{ marginLeft: '1rem' }}>◌ Synthesizing final response...</Text>
              </div>
            </Card>
          )}

          {/* Error Display */}
          {error && (
            <Card className="error-card">
              <Alert
                message="Execution Error"
                description={error}
                type="error"
                showIcon
                closable
                onClose={() => setError(null)}
              />
            </Card>
          )}

          {/* Debug Panel */}
          {(debugMessages.length > 0 || isExecuting) && (
            <Card 
              title={
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <span>Debug Messages</span>
                  <Button 
                    size="small" 
                    onClick={() => setShowDebug(!showDebug)}
                    type={showDebug ? "primary" : "default"}
                  >
                    {showDebug ? 'Hide' : 'Show'} Debug
                  </Button>
                </div>
              }
              className="debug-card"
            >
              {showDebug && (
                <div className="debug-content">
                  <div className="debug-log">
                    {debugMessages.map((msg) => (
                      <div key={msg.id} className="debug-message">
                        <Text type="secondary" className="debug-timestamp">
                          {msg.timestamp.toLocaleTimeString()}
                        </Text>
                        <Text className="debug-text" code>
                          {msg.message}
                        </Text>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </Card>
          )}

          {/* Welcome Message */}
          {!isExecuting && !response && !error && agentSteps.length === 0 && (
            <Card className="welcome-card">
              <div className="welcome-content">
                <RobotOutlined className="welcome-icon" />
                <Title level={3}>Ready to Assist</Title>
                <Paragraph type="secondary">
                  Configure your agent settings and describe what you'd like to accomplish. 
                  The agent will break down complex tasks and provide real-time feedback 
                  as it works through each step with detailed progress tracking.
                </Paragraph>
                <div className="feature-list">
                  <div className="feature-item">
                    <ThunderboltOutlined /> Real-time task orchestration
                  </div>
                  <div className="feature-item">
                    <EyeOutlined /> Live progress monitoring with emojis
                  </div>
                  <div className="feature-item">
                    <DatabaseOutlined /> RAG integration support
                  </div>
                  <div className="feature-item">
                    <SearchOutlined /> Web search capabilities
                  </div>
                </div>
              </div>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
};

export default AgentsView; 