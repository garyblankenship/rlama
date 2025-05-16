import React, { useState, useEffect } from 'react';
import { 
  Table, Card, Button, Input, Space, 
  Typography, Modal, message, Tag, Tooltip, 
  Empty, Spin, Checkbox, Select
} from 'antd';
import { 
  SearchOutlined, 
  SyncOutlined, 
  EyeOutlined 
} from '@ant-design/icons';
import { ragService } from '../../services/api';

const { Title, Text } = Typography;
const { Search } = Input;
const { Option } = Select;

const ChunksTab = ({ ragName }) => {
  const [chunks, setChunks] = useState([]);
  const [documents, setDocuments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchText, setSearchText] = useState('');
  const [documentFilter, setDocumentFilter] = useState('');
  const showContent = true;
  const [modalVisible, setModalVisible] = useState(false);
  const [selectedChunk, setSelectedChunk] = useState(null);

  // Charger les chunks
  const loadChunks = async () => {
    try {
      setLoading(true);
      const chunksData = await ragService.getRagChunks(ragName, documentFilter, showContent);
      setChunks(chunksData);
    } catch (error) {
      console.error(`Erreur lors du chargement des chunks du RAG "${ragName}":`, error);
      message.error(`Impossible de charger les chunks du RAG "${ragName}"`);
    } finally {
      setLoading(false);
    }
  };

  // Charger les documents pour le filtre
  const loadDocuments = async () => {
    try {
      const docs = await ragService.getRagDocuments(ragName);
      setDocuments(docs);
    } catch (error) {
      console.error(`Erreur lors du chargement des documents du RAG "${ragName}":`, error);
    }
  };

  // Charger les données au montage du composant
  useEffect(() => {
    loadDocuments();
    loadChunks();
  }, [ragName]);

  // Recharger les chunks quand les filtres changent
  useEffect(() => {
    loadChunks();
  }, [documentFilter, showContent]);

  // Filtrer les chunks par le texte de recherche
  const filteredChunks = chunks.filter(chunk => {
    const searchTarget = showContent 
      ? chunk.content?.toLowerCase() || ''
      : chunk.id.toLowerCase();
    
    return searchTarget.includes(searchText.toLowerCase());
  });

  // Afficher le contenu d'un chunk dans un modal
  const showChunkDetails = (chunk) => {
    setSelectedChunk(chunk);
    setModalVisible(true);
  };

  // Si aucun contenu n'est disponible, charger le contenu d'un chunk spécifique
  const loadChunkContent = async (chunkId) => {
    try {
      const chunkWithContent = await ragService.getRagChunks(ragName, chunkId, true);
      if (chunkWithContent && chunkWithContent.length > 0) {
        // Mettre à jour le chunk sélectionné avec son contenu
        setSelectedChunk({
          ...selectedChunk,
          content: chunkWithContent[0].content
        });
      }
    } catch (error) {
      console.error('Erreur lors du chargement du contenu du chunk:', error);
      message.error('Impossible de charger le contenu du chunk');
    }
  };

  // Colonnes du tableau
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      ellipsis: true,
      width: '15%',
      render: (text) => {
        // Extraire le numéro de chunk à partir de l'ID
        let chunkNumber = '';
        const match = text.match(/_chunk_(\d+)$/);
        if (match && match[1]) {
          chunkNumber = match[1];
        } else {
          // Si le format est différent, afficher l'ID complet
          chunkNumber = text;
        }
        
        return (
          <Tooltip title={text}>
            <span style={{ display: 'inline-block', maxWidth: '100%', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
              {chunkNumber}
            </span>
          </Tooltip>
        );
      },
    },
    {
      title: 'Document',
      dataIndex: 'document_id',
      key: 'document_id',
      ellipsis: true,
      width: '35%',
      render: (text, record) => {
        // Extraire le nom du document à partir de l'ID du chunk
        let documentName = '';
        
        // Essayer d'abord de trouver le document dans la liste des documents
        const doc = documents.find(d => d.id === text);
        if (doc) {
          documentName = doc.name;
        } else {
          // Sinon, extraire le nom du fichier de l'ID du chunk
          const chunkId = record.id;
          const match = chunkId.match(/^(.+?)_chunk_\d+$/);
          if (match && match[1]) {
            documentName = match[1];
          } else {
            // Si aucun pattern ne correspond, utiliser l'ID du document tel quel
            documentName = text;
          }
        }
        
        return (
          <Tooltip title={documentName}>
            <span style={{ display: 'inline-block', maxWidth: '100%', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
              {documentName}
            </span>
          </Tooltip>
        );
      },
    },
    {
      title: 'Content',
      dataIndex: 'content',
      key: 'content',
      ellipsis: true,
      width: '40%',
      render: (text) => (
        <div style={{ 
          maxHeight: '60px', 
          overflow: 'hidden', 
          textOverflow: 'ellipsis',
          display: '-webkit-box',
          WebkitLineClamp: 2,
          WebkitBoxOrient: 'vertical'
        }}>
          {text || 'Non available'}
        </div>
      ),
    },
    {
      title: 'Actions',
      key: 'actions',
      width: '10%',
      align: 'center',
      render: (_, record) => (
        <Button 
          type="link" 
          onClick={() => showChunkDetails(record)}
          icon={<EyeOutlined />}
        >
          View
        </Button>
      ),
    },
  ];

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={4}>Chunks ({chunks.length})</Title>
        
        <Space>
          <Select
            placeholder="Filter by document"
            allowClear
            style={{ width: 250 }}
            onChange={(value) => setDocumentFilter(value || '')}
          >
            {documents.map(doc => (
              <Option key={doc.id} value={doc.id}>{doc.name}</Option>
            ))}
          </Select>
          
          <Search
            placeholder="Search for a chunk"
            allowClear
            onChange={(e) => setSearchText(e.target.value)}
            style={{ width: 250 }}
          />
          
          <Button 
            icon={<SyncOutlined />} 
            onClick={loadChunks}
            disabled={loading}
          >
            Refresh
          </Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={filteredChunks}
        rowKey="id"
        loading={loading}
        locale={{ 
          emptyText: searchText || documentFilter 
            ? 'No matching chunk found' 
            : 'No chunk in this RAG' 
        }}
        pagination={{ 
          pageSize: 10,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100'],
        }}
        style={{ overflowX: 'auto' }}
      />

      {/* Modal pour afficher le contenu détaillé d'un chunk */}
      <Modal
        title={`Chunk: ${selectedChunk?.id}`}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setModalVisible(false)}>
            Close
          </Button>
        ]}
        width={800}
      >
        {selectedChunk ? (
          <div>
            <p><strong>Document:</strong> {selectedChunk.document_id}</p>
            <Divider />
            <div className="code-content" style={{ maxHeight: '400px', overflow: 'auto' }}>
              {selectedChunk.content ? (
                selectedChunk.content
              ) : (
                <div>
                  <p>Contenu non disponible. Cliquez pour charger.</p>
                  <Button onClick={() => loadChunkContent(selectedChunk.id)}>
                    Charger le contenu
                  </Button>
                </div>
              )}
            </div>
          </div>
        ) : (
          <Empty description="No chunk selected" />
        )}
      </Modal>
    </div>
  );
};

// Composant Divider manquant dans les imports
const Divider = ({ children, ...props }) => (
  <div 
    style={{ 
      width: '100%', 
      height: '1px', 
      backgroundColor: '#f0f0f0', 
      margin: '16px 0',
      position: 'relative',
    }}
    {...props}
  >
    {children && (
      <div style={{ 
        position: 'absolute', 
        backgroundColor: '#fff', 
        padding: '0 8px', 
        top: '-10px', 
        left: '16px' 
      }}>
        {children}
      </div>
    )}
  </div>
);

export default ChunksTab; 