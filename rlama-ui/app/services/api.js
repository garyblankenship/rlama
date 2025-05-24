import axios from 'axios';

const API_URL = 'http://localhost:5001';

// Création d'une instance axios préconfigurée avec timeout augmenté
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 60000, // Increased to 60 seconds (from 10 seconds)
});

// Add request interceptor for debugging
api.interceptors.request.use(config => {
  // Use longer timeout for specific operations
  if (config.url === '/rags' && config.method === 'post') {
    config.timeout = 120000; // 2 minutes for RAG creation
  }
  console.log(`API Request: ${config.method.toUpperCase()} ${config.url}`);
  return config;
}, error => {
  console.error('API Request Error:', error);
  return Promise.reject(error);
});

// Add response interceptor for debugging
api.interceptors.response.use(response => {
  console.log(`API Response: ${response.status} ${response.config.url}`);
  return response;
}, error => {
  if (error.response) {
    console.error(`API Error: ${error.response.status} ${error.config.url}`);
    console.error('Error details:', error.response.data);
  } else if (error.request) {
    console.error(`API Error: No response received for ${error.config?.url || 'unknown URL'}`);
    console.error('Request details:', error.request);
  } else {
    console.error('API Error:', error.message);
  }
  return Promise.reject(error);
});

// Services API pour les RAG
export const ragService = {
  // Récupération de tous les RAG
  getAllRags: async () => {
    const response = await api.get('/rags');
    return response.data;
  },

  // Création d'un nouveau RAG
  createRag: async (ragData) => {
    const response = await api.post('/rags', ragData);
    return response.data;
  },

  // Suppression d'un RAG
  deleteRag: async (ragName) => {
    const response = await api.delete(`/rags/${ragName}`);
    return response.data;
  },

  // Récupération des documents d'un RAG
  getRagDocuments: async (ragName) => {
    const response = await api.get(`/rags/${ragName}/documents`);
    return response.data;
  },

  // Récupération des chunks d'un RAG
  getRagChunks: async (ragName, documentFilter = '', showContent = false) => {
    const params = new URLSearchParams();
    if (documentFilter) params.append('document_filter', documentFilter);
    if (showContent) params.append('show_content', showContent);
    
    const response = await api.get(`/rags/${ragName}/chunks?${params.toString()}`);
    return response.data;
  },

  // Ajout de documents à un RAG
  addDocuments: async (ragName, folderPath) => {
    const formData = new FormData();
    formData.append('folder_path', folderPath);
    
    const response = await api.post(`/rags/${ragName}/documents`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  // Mise à jour du modèle d'un RAG
  updateModel: async (ragName, modelName) => {
    const formData = new FormData();
    formData.append('model_name', modelName);
    
    const response = await api.put(`/rags/${ragName}/model`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  // Configuration de la surveillance de dossier
  setupWatch: async (watchData) => {
    const response = await api.post(`/rags/${watchData.rag_name}/watch`, watchData);
    return response.data;
  },

  // Configuration de la surveillance web
  setupWebWatch: async (webWatchData) => {
    const response = await api.post(`/rags/${webWatchData.rag_name}/web-watch`, webWatchData);
    return response.data;
  },

  // Désactivation de la surveillance de dossier
  disableWatch: async (ragName) => {
    const response = await api.delete(`/rags/${ragName}/watch`);
    return response.data;
  },

  // Désactivation de la surveillance web
  disableWebWatch: async (ragName) => {
    const response = await api.delete(`/rags/${ragName}/web-watch`);
    return response.data;
  },

  // Vérification forcée des dossiers/sites surveillés
  checkWatched: async (ragName) => {
    const response = await api.post(`/rags/${ragName}/check-watched`);
    return response.data;
  },

  // Interrogation d'un RAG
  queryRag: async (queryData) => {
    const response = await api.post('/query', queryData, {
      timeout: 180000,
    });
    return response.data;
  },

  // Récupération des modèles disponibles
  getAvailableModels: async () => {
    const response = await api.get('/models');
    return response.data.models;
  },

  queryRagStream: (queryData, { onProgress, onAnswerChunk, onError, onDone }) => {
    const controller = new AbortController(); // For potential cancellation

    const fetchData = async () => {
      try {
        console.log("Starting SSE stream request", queryData);
        const response = await fetch(`${API_URL}/query-stream`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'text/event-stream' // Important to tell the server we expect SSE
          },
          body: JSON.stringify(queryData),
          signal: controller.signal, // For cancellation
        });

        if (!response.ok) {
          // Try to parse error from backend if response is not OK
          let errorDetail = `Request failed with status ${response.status}`;
          try {
            const errorData = await response.json();
            errorDetail = errorData.detail || JSON.stringify(errorData);
          } catch (e) {
            // If response is not JSON, use status text
            errorDetail = response.statusText || errorDetail;
          }
          console.error("Stream failed with error:", errorDetail);
          onError(errorDetail);
          onDone(); // Ensure onDone is called
          return;
        }

        console.log("SSE stream started successfully");
        
        // Use native EventSource for standard SSE handling if possible
        if (typeof EventSource !== 'undefined' && window.ReadableStream) {
          // Modern approach using native EventSource-like implementation with ReadableStream
          const reader = response.body.getReader();
          const decoder = new TextDecoder();
          let buffer = '';

          while (true) {
            const { done, value } = await reader.read();
            if (done) {
              console.log("Stream complete (done=true)");
              break;
            }

            const chunk = decoder.decode(value, { stream: true });
            console.log("Received SSE chunk:", chunk.substring(0, 50) + (chunk.length > 50 ? '...' : ''));
            buffer += chunk;

            // Process completed SSE messages (separated by double newlines)
            let eventEnd;
            // Look for both \r\n\r\n (standard) and \n\n (sometimes used)
            const standardDelimiter = '\r\n\r\n';
            const altDelimiter = '\n\n';
            
            while ((eventEnd = buffer.indexOf(standardDelimiter)) >= 0 || 
                   (eventEnd = buffer.indexOf(altDelimiter)) >= 0) {
              
              // Determine which delimiter was found
              const foundStandardDelimiter = buffer.indexOf(standardDelimiter) === eventEnd;
              const delimiter = foundStandardDelimiter ? standardDelimiter : altDelimiter;
              const delimiterLength = delimiter.length;
              
              let eventData = buffer.substring(0, eventEnd);
              buffer = buffer.substring(eventEnd + delimiterLength);
              
              // Process SSE data fields
              const lines = eventData.split(/\r?\n/);
              for (const line of lines) {
                if (line.startsWith('data: ')) {
                  try {
                    const jsonData = JSON.parse(line.substring(6));
                    console.log("Parsed SSE event:", jsonData.type);
                    
                    if (jsonData.type === 'progress') {
                      onProgress(jsonData.content);
                    } else if (jsonData.type === 'answer_chunk') {
                      onAnswerChunk(jsonData.content);
                    } else if (jsonData.type === 'error') {
                      console.error('SSE Error Event:', jsonData.content);
                      onError(jsonData.content);
                    } else if (jsonData.type === 'done') {
                      console.log("Received done event from stream");
                      // Don't break here - more data might be present in buffer
                    }
                  } catch (e) {
                    console.error('Error parsing SSE JSON data:', e, "Raw line:", line);
                  }
                }
              }
            }
          }
        } else {
          // Fallback to manual text parsing
          const text = await response.text();
          console.log("Fallback: received complete response text");
          
          // Parse the complete response text as SSE
          const events = text.split(/\r?\n\r?\n/).filter(e => e.trim());
          for (const event of events) {
            if (event.startsWith('data: ')) {
              try {
                const jsonData = JSON.parse(event.substring(6));
                
                if (jsonData.type === 'progress') {
                  onProgress(jsonData.content);
                } else if (jsonData.type === 'answer_chunk') {
                  onAnswerChunk(jsonData.content);
                } else if (jsonData.type === 'error') {
                  onError(jsonData.content);
                }
              } catch (e) {
                console.error('Error parsing SSE event:', e);
              }
            }
          }
        }
        
        console.log("Stream processing complete");
        onDone();
      } catch (err) {
        if (err.name === 'AbortError') {
          console.log('Stream fetch aborted by user');
        } else {
          console.error('SSE connection error:', err);
          onError(err.message || 'Streaming connection failed');
        }
        onDone();
      }
    };
    
    fetchData();

    // Return a cleanup function to abort the fetch request if needed
    return () => {
      console.log("Aborting stream fetch");
      controller.abort();
    };
  },

  // Récupération du statut de la surveillance de dossier
  getWatchStatus: async (ragName) => {
    try {
      const response = await api.get(`/rags/${ragName}/watch-status`);
      return response.data;
    } catch (error) {
      console.error(`Erreur lors de la récupération du statut de surveillance pour ${ragName}:`, error);
      return null;
    }
  },

  // Récupération du statut de la surveillance web
  getWebWatchStatus: async (ragName) => {
    try {
      const response = await api.get(`/rags/${ragName}/web-watch-status`);
      return response.data;
    } catch (error) {
      console.error(`Erreur lors de la récupération du statut de surveillance web pour ${ragName}:`, error);
      return null;
    }
  },
};

// Fonction utilitaire pour extraire la version courte
const extractShortVersion = (stdout) => {
  if (!stdout) return null;
  
  // Pour RLAMA: extraire juste "RLAMA version X.X.X"
  const rlamaMatch = stdout.match(/RLAMA version (\S+)/i);
  if (rlamaMatch) return `v${rlamaMatch[1]}`;
  
  // Pour Ollama: extraire juste "ollama version X.X.X"
  const ollamaMatch = stdout.match(/ollama version (\S+)/i);
  if (ollamaMatch) return `v${ollamaMatch[1]}`;
  
  // Si pas de match, prendre la première ligne non vide et la limiter
  const firstLine = stdout.split('\n')[0]?.trim();
  return firstLine?.length > 50 ? `${firstLine.substring(0, 50)}...` : firstLine;
};

// Service de vérification de la disponibilité du backend et autres dépendances
export const healthService = {
  checkHealth: async () => {
    try {
      // D'abord vérifier si le CLI RLAMA est disponible
      const rlamaCliCheck = await healthService.checkRlamaCli();
      if (!rlamaCliCheck.available) {
        console.error('RLAMA CLI not available:', rlamaCliCheck.details);
        return false;
      }

      // Ensuite vérifier si l'API répond
      const response = await api.get('/health');
      return response.status === 200;
    } catch (error) {
      console.error('Health check failed:', error.message);
      return false;
    }
  },
  
  // Vérification d'Ollama (Server/API status via /check-ollama)
  checkOllama: async () => {
    try {
      const response = await api.get('/check-ollama'); // This endpoint is currently 404ing
      // Assuming a successful response looks like { available: true, status: 'ok', ... }
      return response.data && response.data.available && response.data.status === 'ok';
    } catch (error) {
      console.error('Ollama server check (/check-ollama) failed:', error.message);
      return false; // Indicates the primary check for Ollama server failed
    }
  },
  
  // Vérification des modèles LLM disponibles (via /models)
  checkModels: async () => {
    try {
      const response = await api.get('/models');
      return response.data && Array.isArray(response.data.models) && response.data.models.length > 0;
    } catch (error) {
      console.error('Models check failed:', error.message);
      return false;
    }
  },
  
  // Vérification des modèles d'embeddings
  checkEmbeddings: async () => {
    try {
      const response = await api.get('/embedding-models'); // This endpoint is 404ing
      return response.data && Array.isArray(response.data.models) && response.data.models.length > 0;
    } catch (error) {
      try {
        const configResponse = await api.get('/config'); // Fallback to /config
        return configResponse.data && configResponse.data.embedding_model != null;
      } catch (configError) {
        console.error('Embeddings check failed (both /embedding-models and /config):', error.message);
        return false;
      }
    }
  },

  // NOUVELLE FONCTION: Vérification de l'existence et version du CLI Ollama via /exec
  checkOllamaCli: async () => {
    try {
      console.log("Attempting to check Ollama CLI via /exec?command=ollama+--version");
      const response = await api.get(`/exec?command=${encodeURIComponent('ollama --version')}`);
      
      if (response.data && response.data.stdout && 
          !response.data.stdout.toLowerCase().includes("command not found") &&
          !response.data.stdout.toLowerCase().includes("n'est pas reconnu") &&
          (!response.data.stderr || response.data.stderr.trim() === "")
      ) {
        const shortVersion = extractShortVersion(response.data.stdout);
        console.log("Ollama CLI found:", shortVersion);
        return { available: true, version: shortVersion, details: `Ollama CLI available (${shortVersion})` };
      } else {
        const errorDetails = response.data.stderr || response.data.stdout || "Ollama CLI not found or error during execution via /exec.";
        console.warn("Ollama CLI not found or error:", errorDetails);
        return { available: false, details: errorDetails };
      }
    } catch (error) {
      const errorMsg = error.response?.data?.detail || error.response?.data?.message || error.message;
      console.error('Ollama CLI check (/exec) catastrophically failed:', errorMsg);
      return { available: false, details: `Error executing 'ollama --version' via /exec: ${errorMsg}` };
    }
  },

  // NOUVELLE FONCTION: Vérification de l'existence et version du CLI RLAMA via /exec
  checkRlamaCli: async () => {
    try {
      console.log("Attempting to check RLAMA CLI via /exec?command=rlama+--version");
      const response = await api.get(`/exec?command=${encodeURIComponent('rlama --version')}`);
      
      if (response.data && response.data.stdout &&
          !response.data.stdout.toLowerCase().includes("command not found") &&
          !response.data.stdout.toLowerCase().includes("n'est pas reconnu") &&
          (!response.data.stderr || response.data.stderr.trim() === "")
      ) {
        const shortVersion = extractShortVersion(response.data.stdout);
        console.log("RLAMA CLI found:", shortVersion);
        return { available: true, version: shortVersion, details: `RLAMA CLI available (${shortVersion})` };
      } else {
        const errorDetails = response.data.stderr || response.data.stdout || "RLAMA CLI not found or error during execution via /exec.";
        console.warn("RLAMA CLI not found or error:", errorDetails);
        return { available: false, details: errorDetails };
      }
    } catch (error) {
      const errorMsg = error.response?.data?.detail || error.response?.data?.message || error.message;
      console.error('RLAMA CLI check (/exec) catastrophically failed:', errorMsg);
      return { available: false, details: `Error executing 'rlama --version' via /exec: ${errorMsg}` };
    }
  },
  
  // Exécution de `ollama list` via /exec comme fallback (si /check-ollama n'est pas bon)
  // Cette fonction est déjà référencée dans Home.jsx pour le statut Ollama server.
  execOllamaList: async () => {
    try {
      console.log("Attempting to execute 'ollama list' via /exec");
      const response = await api.get(`/exec?command=${encodeURIComponent('ollama list')}`);
      // La logique pour interpréter `stdout` et `stderr` est dans Home.jsx
      // Ici, on retourne simplement la réponse brute de /exec
      return response.data; 
    } catch (error) {
      console.error("Failed to execute 'ollama list' via /exec:", error);
      // Renvoyer une structure qui peut être gérée par l'appelant dans Home.jsx
      return { 
        stdout: null, 
        stderr: error.response?.data?.detail || error.response?.data?.message || error.message || "Network error or /exec endpoint failure.",
        error: true // Flag to indicate the call itself failed
      };
    }
  }
};

// Service de vérification pour Ollama
export const ollamaService = {
  checkOllama: async () => {
    try {
      // Essayer d'abord via l'API backend
      const response = await api.get('/check-ollama');
      return response.data;
    } catch (error) {
      if (error.response && error.response.status === 404) {
        // Endpoint non implémenté, essayer de lancer une commande
        try {
          const execResponse = await api.get('/exec?command=ollama+list');
          return {
            available: !!execResponse.data.stdout,
            models: [],
            message: 'Detected via shell command'
          };
        } catch (execError) {
          console.error('Error executing command:', execError);
          return { available: false, message: execError.message };
        }
      }
      console.error('Error checking Ollama:', error);
      return { available: false, message: error.message };
    }
  },
  
  // Exécuter une commande arbitraire (pour le développement seulement)
  executeCommand: async (command) => {
    try {
      const response = await api.get(`/exec?command=${encodeURIComponent(command)}`);
      return response.data;
    } catch (error) {
      console.error('Error executing command:', error);
      throw error;
    }
  }
};

// Service de gestion des paramètres et profils
export const settingsService = {
  // Gestion des profils
  getProfiles: async () => {
    try {
      console.log('🔍 API: Getting profiles...');
      const response = await api.get('/profiles');
      console.log('✅ API: Profiles response:', response.data);
      console.log('✅ API: Profiles type:', typeof response.data);
      console.log('✅ API: Profiles array?', Array.isArray(response.data));
      
      // Fix: Backend returns {profiles: [array]}, we need to extract the profiles array
      const profilesArray = response.data.profiles || response.data;
      
      if (Array.isArray(profilesArray)) {
        console.log('✅ API: Found', profilesArray.length, 'profiles');
        profilesArray.forEach((profile, index) => {
          console.log(`📋 Profile ${index}:`, {
            name: profile.name,
            provider: profile.provider,
            created_on: profile.created_on,
            last_used: profile.last_used,
            description: profile.description
          });
        });
      }
      
      // Also log the raw response for debugging
      console.log('🔍 API: Raw response status:', response.status);
      console.log('🔍 API: Raw response headers:', response.headers);
      
      return profilesArray;
    } catch (error) {
      console.error('❌ API: Error getting profiles:', error);
      // Try direct CLI command as fallback
      try {
        console.log('🔄 API: Trying direct CLI command...');
        const execResponse = await api.get('/exec?command=rlama+profile+list');
        console.log('🔄 API: CLI response:', execResponse.data);
        
        if (execResponse.data.stdout) {
          // Parse the tabular output
          const lines = execResponse.data.stdout.trim().split('\n');
          const profiles = [];
          
          // Skip first two lines (header and empty line)
          for (let i = 3; i < lines.length; i++) {
            const line = lines[i].trim();
            if (line && !line.startsWith('-')) {
              const parts = line.split(/\s+/);
              if (parts.length >= 4) {
                profiles.push({
                  name: parts[0],
                  provider: parts[1],
                  created_at: parts[2],
                  last_used_at: parts[3] === 'never' ? null : parts[3]
                });
              }
            }
          }
          
          console.log('✅ API: Parsed CLI profiles:', profiles);
          return profiles;
        }
      } catch (execError) {
        console.warn('❌ API: CLI fallback failed:', execError);
      }
      return [];
    }
  },

  // Debug function to see raw CLI output
  debugProfilesRaw: async () => {
    try {
      console.log('🔍 API: Getting raw CLI output for profiles...');
      const execResponse = await api.get('/exec?command=rlama+profile+list');
      console.log('🔍 API: Raw CLI stdout:', execResponse.data.stdout);
      console.log('🔍 API: Raw CLI stderr:', execResponse.data.stderr);
      console.log('🔍 API: Raw CLI return code:', execResponse.data.returncode);
      
      // Also test the JSON version
      const execJsonResponse = await api.get('/exec?command=rlama+profile+list+--json');
      console.log('🔍 API: Raw CLI JSON stdout:', execJsonResponse.data.stdout);
      console.log('🔍 API: Raw CLI JSON stderr:', execJsonResponse.data.stderr);
      
      return {
        normal: execResponse.data,
        json: execJsonResponse.data
      };
    } catch (error) {
      console.error('❌ API: Error getting raw CLI output:', error);
      return null;
    }
  },

  createProfile: async (profileData) => {
    try {
      console.log('🌐 API: Creating profile with data:', {
        name: profileData.name,
        provider: profileData.provider,
        api_key: profileData.api_key ? `${profileData.api_key.substring(0, 10)}...` : 'MISSING',
        description: profileData.description
      });
      
      const response = await api.post('/profiles', profileData);
      console.log('✅ API: Profile created successfully:', response.data);
      console.log('✅ API: Response status:', response.status);
      console.log('✅ API: Response headers:', response.headers);
      console.log('✅ API: Full response object:', response);
      
      // Immediately test if the profile was actually saved
      console.log('🔍 API: Testing immediate profile retrieval...');
      try {
        const testResponse = await api.get('/profiles');
        console.log('🔍 API: Immediate profiles check:', testResponse.data);
        console.log('🔍 API: Found profiles count:', testResponse.data?.length || 0);
        
        if (testResponse.data && Array.isArray(testResponse.data)) {
          const foundProfile = testResponse.data.find(p => p.name === profileData.name);
          if (foundProfile) {
            console.log('✅ API: Profile found immediately after creation:', foundProfile);
          } else {
            console.log('❌ API: Profile NOT found immediately after creation!');
            console.log('❌ API: Available profiles:', testResponse.data.map(p => p.name));
          }
        }
      } catch (testError) {
        console.warn('⚠️ API: Could not test immediate profile retrieval:', testError);
      }
      
      return response.data;
    } catch (error) {
      console.error('❌ API: Error creating profile via API:', error);
      console.error('❌ API: Error response:', {
        status: error.response?.status,
        statusText: error.response?.statusText,
        data: error.response?.data,
        headers: error.response?.headers
      });
      
      // Don't try CLI fallback for profile creation due to URL length limits with API keys
      // Instead, throw a more descriptive error
      if (error.response?.status === 500) {
        throw new Error('Server error creating profile. Please check your API key format and try again.');
      } else if (error.response?.status === 400) {
        throw new Error('Invalid profile data. Please check all fields and try again.');
      } else if (error.response?.status === 409) {
        throw new Error('A profile with this name already exists. Please choose a different name.');
      } else if (!error.response) {
        throw new Error('Network error. Please check that the RLAMA backend is running on localhost:5001');
      } else {
        throw new Error(`Failed to create profile: ${error.response?.data?.detail || error.message}`);
      }
    }
  },

  updateProfile: async (profileName, profileData) => {
    try {
      const response = await api.put(`/profiles/${profileName}`, profileData);
      return response.data;
    } catch (error) {
      console.error('Error updating profile via API:', error);
      
      // Don't try CLI fallback for profile updates due to URL length limits with API keys
      if (error.response?.status === 404) {
        throw new Error('Profile not found. It may have been deleted.');
      } else if (error.response?.status === 400) {
        throw new Error('Invalid profile data. Please check all fields and try again.');
      } else if (error.response?.status === 500) {
        throw new Error('Server error updating profile. Please try again.');
      } else {
        throw new Error(`Failed to update profile: ${error.response?.data?.detail || error.message}`);
      }
    }
  },

  deleteProfile: async (profileName) => {
    try {
      const response = await api.delete(`/profiles/${profileName}`);
      return response.data;
    } catch (error) {
      console.error('Error deleting profile via API:', error);
      
      // Try CLI fallback for deletion (safer since no long API key involved)
      try {
        const command = `rlama profile delete ${profileName}`;
        const execResponse = await api.get(`/exec?command=${encodeURIComponent(command)}`);
        if (execResponse.data.stderr && execResponse.data.stderr.includes('not found')) {
          throw new Error('Profile not found. It may have already been deleted.');
        }
        return { success: !execResponse.data.stderr };
      } catch (execError) {
        console.error('CLI fallback failed:', execError);
        if (error.response?.status === 404) {
          throw new Error('Profile not found. It may have already been deleted.');
        } else {
          throw new Error(`Failed to delete profile: ${error.response?.data?.detail || error.message}`);
        }
      }
    }
  },

  // Gestion des clés API
  getApiKeys: async () => {
    try {
      const response = await api.get('/settings/api-keys');
      return response.data;
    } catch (error) {
      console.error('Error getting API keys:', error);
      // Try to get from localStorage as fallback
      try {
        const stored = localStorage.getItem('rlama_api_keys');
        return stored ? JSON.parse(stored) : {};
      } catch (storageError) {
        console.warn('Could not read from localStorage:', storageError);
        return {};
      }
    }
  },

  saveApiKeys: async (apiKeysData) => {
    try {
      // Save to backend first
      const response = await api.post('/settings/api-keys', apiKeysData);
      
      // Also save to localStorage as backup
      localStorage.setItem('rlama_api_keys', JSON.stringify(apiKeysData));
      
      return response.data;
    } catch (error) {
      console.error('Error saving API keys to backend:', error);
      
      // Fallback to localStorage only
      try {
        localStorage.setItem('rlama_api_keys', JSON.stringify(apiKeysData));
        console.log('API keys saved to localStorage as fallback');
        return { success: true, source: 'localStorage' };
      } catch (storageError) {
        console.error('Could not save to localStorage either:', storageError);
        throw error;
      }
    }
  },

  // Set environment variable for CLI commands
  setEnvironmentVariable: async (name, value) => {
    try {
      const response = await api.post('/settings/environment', { name, value });
      return response.data;
    } catch (error) {
      console.error('Error setting environment variable:', error);
      // Try via exec command as fallback
      try {
        const command = `export ${name}="${value}"`;
        const execResponse = await api.get(`/exec?command=${encodeURIComponent(command)}`);
        return { success: !execResponse.data.stderr, source: 'exec' };
      } catch (execError) {
        console.warn('Could not set environment variable via exec:', execError);
        throw error;
      }
    }
  },

  // Paramètres généraux
  getGeneralSettings: async () => {
    try {
      const response = await api.get('/settings/general');
      return response.data;
    } catch (error) {
      console.error('Error getting general settings:', error);
      return {
        auto_save: true,
        show_notifications: true,
        default_model: 'gpt-3.5-turbo'
      };
    }
  },

  saveGeneralSettings: async (settingsData) => {
    try {
      const response = await api.post('/settings/general', settingsData);
      return response.data;
    } catch (error) {
      console.error('Error saving general settings:', error);
      throw error;
    }
  },

  // Modèles disponibles
  getAvailableModels: async () => {
    try {
      const response = await api.get('/models');
      return response.data.models;
    } catch (error) {
      console.error('Error getting available models:', error);
      return [];
    }
  }
};

// Services API pour les agents
export const agentService = {
  // Exécution d'un agent
  runAgent: async (queryData) => {
    const response = await api.post('/agent/run', queryData, {
      timeout: 300000, // 5 minutes
    });
    return response.data;
  },

  // Exécution d'un agent avec streaming
  runAgentStream: (queryData, { onProgress, onTaskUpdate, onAnswerChunk, onError, onDone }) => {
    const controller = new AbortController();

    const fetchData = async () => {
      try {
        console.log("Starting agent stream request", queryData);
        const response = await fetch(`${API_URL}/agent/stream`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'text/event-stream'
          },
          body: JSON.stringify(queryData),
          signal: controller.signal,
        });

        if (!response.ok) {
          let errorDetail = `Request failed with status ${response.status}`;
          try {
            const errorData = await response.json();
            errorDetail = errorData.detail || JSON.stringify(errorData);
          } catch (e) {
            errorDetail = response.statusText || errorDetail;
          }
          console.error("Agent stream failed with error:", errorDetail);
          onError(errorDetail);
          onDone();
          return;
        }

        console.log("Agent stream started successfully");
        
        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';

        while (true) {
          const { done, value } = await reader.read();
          if (done) {
            console.log("Agent stream complete (done=true)");
            break;
          }

          const chunk = decoder.decode(value, { stream: true });
          console.log("Received agent SSE chunk:", chunk.substring(0, 50) + (chunk.length > 50 ? '...' : ''));
          buffer += chunk;

          // Process completed SSE messages
          const standardDelimiter = '\r\n\r\n';
          const altDelimiter = '\n\n';
          
          let eventEnd;
          while ((eventEnd = buffer.indexOf(standardDelimiter)) >= 0 || 
                 (eventEnd = buffer.indexOf(altDelimiter)) >= 0) {
            
            const foundStandardDelimiter = buffer.indexOf(standardDelimiter) === eventEnd;
            const delimiter = foundStandardDelimiter ? standardDelimiter : altDelimiter;
            const delimiterLength = delimiter.length;
            
            let eventData = buffer.substring(0, eventEnd);
            buffer = buffer.substring(eventEnd + delimiterLength);
            
            // Process SSE data fields
            const lines = eventData.split(/\r?\n/);
            for (const line of lines) {
              if (line.startsWith('data: ')) {
                try {
                  const jsonData = JSON.parse(line.substring(6));
                  console.log("Parsed agent SSE event:", jsonData.type);
                  
                  if (jsonData.type === 'progress') {
                    onProgress(jsonData.content);
                  } else if (jsonData.type === 'task_update') {
                    onTaskUpdate(jsonData.content);
                  } else if (jsonData.type === 'answer_chunk') {
                    onAnswerChunk(jsonData.content);
                  } else if (jsonData.type === 'error') {
                    console.error('Agent SSE Error Event:', jsonData.content);
                    onError(jsonData.content);
                  } else if (jsonData.type === 'done') {
                    console.log("Received done event from agent stream");
                  }
                } catch (e) {
                  console.error('Error parsing agent SSE JSON data:', e, "Raw line:", line);
                }
              }
            }
          }
        }
        
        console.log("Agent stream processing complete");
        onDone();
      } catch (err) {
        if (err.name === 'AbortError') {
          console.log('Agent stream fetch aborted by user');
        } else {
          console.error('Agent SSE connection error:', err);
          onError(err.message || 'Agent streaming connection failed');
        }
        onDone();
      }
    };
    
    fetchData();

    return () => {
      console.log("Aborting agent stream fetch");
      controller.abort();
    };
  },

  // Récupération des modèles disponibles pour les agents
  getAvailableModels: async () => {
    const response = await api.get('/agent/models');
    return response.data.models;
  },
};

export default api; 