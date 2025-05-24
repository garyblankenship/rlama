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
      // ATTENTION: Relies on the backend /exec endpoint.
      // This is generally NOT recommended for production without proper safeguards.
      console.log("Attempting to check Ollama CLI via /exec?command=ollama+-v");
      const response = await api.get(`/exec?command=${encodeURIComponent('ollama -v')}`);
      
      // A successful execution of 'ollama -v' usually prints version to stdout.
      // Check if stdout is present and doesn't contain common error indicators.
      if (response.data && response.data.stdout && 
          !response.data.stdout.toLowerCase().includes("command not found") &&
          !response.data.stdout.toLowerCase().includes("n'est pas reconnu") && // French for "is not recognized"
          (!response.data.stderr || response.data.stderr.trim() === "") // Prefer empty stderr
      ) {
        console.log("Ollama CLI found:", response.data.stdout.trim());
        return { available: true, version: response.data.stdout.trim(), details: response.data.stdout.trim() };
      } else {
        const errorDetails = response.data.stderr || response.data.stdout || "Ollama CLI not found or error during execution via /exec.";
        console.warn("Ollama CLI not found or error:", errorDetails);
        return { available: false, details: errorDetails };
      }
    } catch (error) {
      const errorMsg = error.response?.data?.detail || error.response?.data?.message || error.message;
      console.error('Ollama CLI check (/exec) catastrophically failed:', errorMsg);
      return { available: false, details: `Error executing 'ollama -v' via /exec: ${errorMsg}` };
    }
  },

  // NOUVELLE FONCTION: Vérification de l'existence et version du CLI RLAMA via /exec
  checkRlamaCli: async () => {
    try {
      // ATTENTION: Relies on the backend /exec endpoint with the same security caveats.
      console.log("Attempting to check RLAMA CLI via /exec?command=rlama+-v");
      const response = await api.get(`/exec?command=${encodeURIComponent('rlama -v')}`);
      
      if (response.data && response.data.stdout &&
          !response.data.stdout.toLowerCase().includes("command not found") &&
          !response.data.stdout.toLowerCase().includes("n'est pas reconnu") &&
          (!response.data.stderr || response.data.stderr.trim() === "")
      ) {
        console.log("RLAMA CLI found:", response.data.stdout.trim());
        return { available: true, version: response.data.stdout.trim(), details: response.data.stdout.trim() };
      } else {
        const errorDetails = response.data.stderr || response.data.stdout || "RLAMA CLI not found or error during execution via /exec.";
        console.warn("RLAMA CLI not found or error:", errorDetails);
        return { available: false, details: errorDetails };
      }
    } catch (error) {
      const errorMsg = error.response?.data?.detail || error.response?.data?.message || error.message;
      console.error('RLAMA CLI check (/exec) catastrophically failed:', errorMsg);
      return { available: false, details: `Error executing 'rlama -v' via /exec: ${errorMsg}` };
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

export default api; 