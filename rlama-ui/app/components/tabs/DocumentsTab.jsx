import React, { useState, useEffect } from 'react';
import { 
  Table, Card, Button, Input, Space, 
  Typography, Modal, message, Tag, Tooltip, 
  Empty, Spin, Popconfirm
} from 'antd';
import { 
  FolderOpenOutlined, 
  SearchOutlined, 
  FileAddOutlined, 
  FileTextOutlined, 
  SyncOutlined 
} from '@ant-design/icons';
import { ragService } from '../../services/api';

const { Title, Text, Paragraph } = Typography;
const { Search } = Input;

const DocumentsTab = ({ ragName, refreshParent }) => {
  const [documents, setDocuments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchText, setSearchText] = useState('');
  const [selectedDocument, setSelectedDocument] = useState(null);
  const [documentContent, setDocumentContent] = useState('');
  const [modalVisible, setModalVisible] = useState(false);
  const [contentLoading, setContentLoading] = useState(false);
  const [addingDocuments, setAddingDocuments] = useState(false);

  // Charger les documents
  const loadDocuments = async () => {
    try {
      setLoading(true);
      console.log(`Loading documents for RAG "${ragName}"`);
      const docs = await ragService.getRagDocuments(ragName);
      console.log(`Loaded ${docs.length} documents`, docs);
      setDocuments(docs);
      
      // Informer le parent si aucun document n'est présent
      if (docs.length === 0 && refreshParent) {
        message.warning("No document found in this RAG. You can add some using the button above.");
      }
    } catch (error) {
      console.error(`Erreur lors du chargement des documents du RAG "${ragName}":`, error);
      message.error(`Impossible de charger les documents du RAG "${ragName}": ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Charger les documents au montage du composant
  useEffect(() => {
    loadDocuments();
  }, [ragName]);

  // Filtrer les documents par le texte de recherche
  const filteredDocuments = documents.filter(doc => 
    doc.name.toLowerCase().includes(searchText.toLowerCase())
  );

  // Gérer la sélection de dossier et l'ajout de documents
  const handleAddDocuments = async () => {
    if (window.electron) {
      try {
        setAddingDocuments(true);
        const folderPath = await window.electron.selectDirectory();
        
        if (folderPath) {
          console.log(`Adding documents from folder: ${folderPath}`);
          const result = await ragService.addDocuments(ragName, folderPath);
          console.log("Documents added successfully:", result);
          message.success('Documents added successfully');
          
          // Informer le parent qu'un rafraîchissement est nécessaire
          if (refreshParent) {
            refreshParent();
          } else {
            // Attendre un peu avant de recharger pour laisser le temps aux fichiers d'être traités
            setTimeout(() => {
              loadDocuments();
            }, 1000);
          }
        }
      } catch (error) {
        console.error('Erreur lors de l\'ajout de documents:', error);
        message.error(`Impossible d'ajouter les documents: ${error.response?.data?.detail || error.message}`);
      } finally {
        setAddingDocuments(false);
      }
    } else {
      message.warning('Folder selection via interface is not available');
    }
  };

  // Afficher le contenu d'un document
  const showDocumentContent = async (document) => {
    try {
      setSelectedDocument(document);
      setModalVisible(true);
      setContentLoading(true);
      
      console.log(`Loading content for document: ${document.id}`);
      // Récupérer les chunks de ce document pour avoir le contenu
      const chunks = await ragService.getRagChunks(ragName, document.id, true);
      console.log(`Loaded ${chunks.length} chunks for document ${document.id}`);
      
      // Si aucun chunk n'est trouvé
      if (!chunks || chunks.length === 0) {
        setDocumentContent("Aucun contenu disponible pour ce document.");
        return;
      }
      
      // Trier les chunks par position
      chunks.sort((a, b) => {
        // Try to extract position numbers (format "X of Y")
        try {
          const posA = parseInt(a.position.split(' ')[0]);
          const posB = parseInt(b.position.split(' ')[0]);
          return posA - posB;
        } catch (e) {
          // En cas d'erreur, ne pas changer l'ordre
          return 0;
        }
      });
      
      // Déterminer le type de contenu pour un formatage approprié
      const contentType = document.content_type.toLowerCase();
      const isPdf = contentType.includes('pdf');
      const isCode = ['json', 'xml', 'html', 'javascript', 'python', 'java', 'c++', 'typescript'].some(
        type => contentType.includes(type)
      );
      
      // Limiter la longueur totale pour les grands documents
      let combinedContent = '';
      const maxChunksToShow = 3; // Nombre de chunks à afficher pour les grands documents
      
      if (chunks.length > maxChunksToShow) {
        // Afficher le premier chunk
        combinedContent += "### Start of document ###\n\n";
        combinedContent += chunks[0].content;
        
        // Indication qu'on a sauté du contenu
        combinedContent += "\n\n[...]\n\n";
        
        // Afficher un chunk du milieu si disponible
        const middleIndex = Math.floor(chunks.length / 2);
        combinedContent += chunks[middleIndex].content;
        
        // Indication qu'on a sauté du contenu
        combinedContent += "\n\n[...]\n\n";
        
        // Afficher le dernier chunk
        combinedContent += chunks[chunks.length - 1].content;
        combinedContent += "\n\n### Fin du document ###\n";
        
        // Ajouter une note sur le contenu tronqué
        combinedContent += `\n\n(Document tronqué: ${chunks.length} chunks au total, affichage partiel)`;
      } else {
        // Pour les petits documents, afficher tous les chunks avec un séparateur
        combinedContent = chunks.map(chunk => chunk.content).join('\n\n...\n\n');
      }
      
      // Formater le contenu en fonction du type de fichier
      setDocumentContent(combinedContent || "Contenu vide.");
    } catch (error) {
      console.error('Erreur lors du chargement du contenu:', error);
      setDocumentContent(`Erreur lors du chargement du contenu: ${error.message}`);
    } finally {
      setContentLoading(false);
    }
  };

  // Colonnes du tableau
  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      sorter: (a, b) => a.name.localeCompare(b.name),
      render: (text, record) => (
        <span style={{ 
          display: 'flex', 
          alignItems: 'center',
          maxWidth: '250px', 
          overflow: 'hidden', 
          textOverflow: 'ellipsis', 
          whiteSpace: 'nowrap' 
        }}>
          <FileTextOutlined style={{ marginRight: '8px', color: 'var(--primary-600)' }} />
          <Tooltip title={text}>{text}</Tooltip>
        </span>
      ),
      ellipsis: true,
    },
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      ellipsis: true,
      width: '20%',
      render: (text) => (
        <Tooltip title={text}>
          <span style={{ display: 'inline-block', maxWidth: '100%', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
            {text.length > 15 ? `${text.substring(0, 15)}...` : text}
          </span>
        </Tooltip>
      ),
    },
    {
      title: 'Type',
      dataIndex: 'content_type',
      key: 'content_type',
      render: (text) => {
        // Simplifier le type de contenu pour un affichage plus concis
        let type = text.split('/')[1] || text;
        
        // Raccourcir les types très longs
        if (type.length > 25) {
          // Pour les types complexes comme vnd.openxmlformats-officedocument.wordprocessingml.document
          if (type.includes('vnd.')) {
            // Extraire l'extension de fichier principale
            if (type.includes('wordprocessing')) return <Tag color="blue">docx</Tag>;
            if (type.includes('spreadsheet')) return <Tag color="green">xlsx</Tag>;
            if (type.includes('presentation')) return <Tag color="purple">pptx</Tag>;
            // Format générique pour les autres vnd
            return <Tag color="orange">{type.split('.').pop() || 'vnd'}</Tag>;
          }
          // Tronquer les autres types longs
          return (
            <Tooltip title={type}>
              <Tag color="blue">{type.substring(0, 15)}...</Tag>
            </Tooltip>
          );
        }
        
        // Coloration par type de fichier
        let color = 'blue';
        if (['pdf'].includes(type)) color = 'red';
        if (['docx', 'doc', 'txt', 'md'].includes(type)) color = 'green';
        if (['xlsx', 'xls', 'csv'].includes(type)) color = 'orange';
        if (['pptx', 'ppt'].includes(type)) color = 'purple';
        if (['json', 'xml', 'html', 'htm'].includes(type)) color = 'cyan';
        
        return <Tag color={color}>{type}</Tag>;
      },
    },
    {
      title: 'Size',
      dataIndex: 'size',
      key: 'size',
      sorter: (a, b) => {
        const getSizeValue = (size) => {
          const value = parseFloat(size);
          if (size.includes('KB')) return value * 1024;
          if (size.includes('MB')) return value * 1024 * 1024;
          return value;
        };
        return getSizeValue(a.size) - getSizeValue(b.size);
      },
    },
  ];

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={4}>Documents ({documents.length})</Title>
        <Space>
          <Search
            placeholder="Search for a document"
            allowClear
            onChange={(e) => setSearchText(e.target.value)}
            style={{ width: 250 }}
          />
          <Button 
            type="primary" 
            icon={<FileAddOutlined />} 
            onClick={handleAddDocuments}
            loading={addingDocuments}
            disabled={addingDocuments}
          >
            Add documents
          </Button>
          <Button 
            icon={<SyncOutlined />} 
            onClick={() => {
              loadDocuments();
              if (refreshParent) refreshParent();
            }}
            disabled={loading}
          >
            Refresh
          </Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={filteredDocuments}
        rowKey="id"
        loading={loading}
        locale={{ 
          emptyText: searchText ? 'No matching document found' : 'No document in this RAG' 
        }}
        pagination={{ 
          pageSize: 10,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100'],
        }}
      />
    </div>
  );
};

export default DocumentsTab; 