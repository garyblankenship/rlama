import os
import json
import subprocess
import shutil
import time
from pathlib import Path
from typing import List, Dict, Optional, Any
import uvicorn
from fastapi import FastAPI, HTTPException, UploadFile, File, Form, BackgroundTasks, Request
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import httpx
import asyncio
from fastapi.responses import StreamingResponse
from datetime import datetime

# Chemins
HOME_DIR = os.path.expanduser("~")
RLAMA_DATA_DIR = os.path.join(HOME_DIR, ".rlama")
RLAMA_EXECUTABLE = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(__file__))), "rlama")

# Constants for storage
SETTINGS_DIR = os.path.join(RLAMA_DATA_DIR, "settings")
PROFILES_DIR = os.path.join(RLAMA_DATA_DIR, "profiles")
API_KEYS_FILE = os.path.join(SETTINGS_DIR, "api_keys.json")

# Create necessary directories
os.makedirs(SETTINGS_DIR, exist_ok=True)
os.makedirs(PROFILES_DIR, exist_ok=True)

# Modèles de données
class RagInfo(BaseModel):
    name: str
    model: str
    created_on: str
    documents_count: int
    size: str

# Nouveaux modèles pour les agents
class AgentQueryRequest(BaseModel):
    query: str
    model: Optional[str] = None
    rag_name: Optional[str] = None
    web_search: Optional[bool] = False
    profile: Optional[str] = None
    context: Optional[Dict[str, Any]] = None

class TaskProgress(BaseModel):
    task_id: str
    description: str
    status: str  # "pending", "running", "completed", "failed"
    result: Optional[str] = None
    error: Optional[str] = None
    tool: Optional[str] = None

class AgentResponse(BaseModel):
    response: str
    tasks: List[TaskProgress]
    metadata: Optional[Dict[str, Any]] = None

class DocumentInfo(BaseModel):
    id: str
    name: str
    size: str
    content_type: str

class ChunkInfo(BaseModel):
    id: str
    document_id: str
    position: str
    content: Optional[str] = None

class QueryRequest(BaseModel):
    rag_name: str
    prompt: str
    context_size: Optional[int] = None
    model: Optional[str] = None

class CreateRagRequest(BaseModel):
    name: str
    model: str
    folder_path: str
    chunk_size: Optional[int] = 1000
    chunk_overlap: Optional[int] = 200
    enable_reranker: Optional[bool] = True
    reranker_weight: Optional[float] = 0.7
    profile: Optional[str] = None

class WatchRequest(BaseModel):
    rag_name: str
    folder_path: str
    interval: int = 0  # 0 = uniquement à l'utilisation

class WebWatchRequest(BaseModel):
    rag_name: str
    url: str
    interval: int = 0
    depth: Optional[int] = 1

# Nouveau modèle pour la réponse d'exécution de commande
class CommandResponse(BaseModel):
    stdout: Optional[str] = None
    stderr: Optional[str] = None
    returncode: int

# Ajouter cette classe
class WatchStatus(BaseModel):
    active: bool = False
    folder_path: Optional[str] = None
    interval: Optional[int] = None
    last_check: Optional[str] = None

class WebWatchStatus(BaseModel):
    active: bool = False
    url: Optional[str] = None
    interval: Optional[int] = None
    depth: Optional[int] = None
    last_check: Optional[str] = None

# Initialisation de l'API
app = FastAPI(title="RLAMA UI Backend")

# Configuration CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Service RLAMA
class RlamaService:
    def __init__(self):
        # Vérifier que RLAMA est installé
        try:
            result = subprocess.run([RLAMA_EXECUTABLE, "version"], capture_output=True, text=True)
            if result.returncode != 0:
                print(f"WARNING: RLAMA is not properly installed or accessible at {RLAMA_EXECUTABLE}")
        except FileNotFoundError:
            print(f"ERROR: RLAMA not found at {RLAMA_EXECUTABLE}")
            
        # Vérifier que le dossier .rlama existe, sinon le créer
        os.makedirs(RLAMA_DATA_DIR, exist_ok=True)
    
    def get_all_rags(self) -> List[RagInfo]:
        """List all RAGs available"""
        rags = []
        if not os.path.exists(RLAMA_DATA_DIR):
            return rags
            
        for folder in os.listdir(RLAMA_DATA_DIR):
            folder_path = os.path.join(RLAMA_DATA_DIR, folder)
            info_path = os.path.join(folder_path, "info.json")
            
            if os.path.isdir(folder_path) and os.path.exists(info_path):
                try:
                    with open(info_path, "r") as f:
                        info = json.load(f)
                        
                    # Calculer la taille totale des documents
                    total_size = 0
                    docs = info.get("documents", [])
                    for doc in docs:
                        doc_size = doc.get("size", 0)
                        if isinstance(doc_size, str):
                            try:
                                # Tenter de convertir si c'est une chaîne
                                doc_size = float(doc_size.split()[0])
                            except (ValueError, IndexError):
                                doc_size = 0
                        total_size += doc_size
                        
                    size_formatted = self._format_size(total_size)
                    
                    # Assurer que le document_count est correct
                    doc_count = len(docs)
                    
                    # Si aucun document n'est trouvé mais que des chunks existent
                    # (cas où il y a des documents mais l'info est mal formatée)
                    if doc_count == 0 and len(info.get("chunks", [])) > 0:
                        # Essayer d'estimer à partir des chunks
                        unique_docs = set()
                        for chunk in info.get("chunks", []):
                            if "document_id" in chunk and chunk["document_id"]:
                                unique_docs.add(chunk["document_id"])
                                
                        if unique_docs:
                            doc_count = len(unique_docs)
                            print(f"Estimated document count for {folder} from chunks: {doc_count}")
                            
                            # Si la taille est toujours 0, lui donner une valeur par défaut
                            if total_size == 0:
                                total_size = doc_count * 1024  # Estimation approximative
                                size_formatted = self._format_size(total_size)
                    
                    rags.append(RagInfo(
                        name=info.get("name", folder),
                        model=info.get("model_name", "N/A"),
                        created_on=info.get("created_at", "N/A"),
                        documents_count=doc_count,
                        size=size_formatted
                    ))
                except Exception as e:
                    print(f"Error processing RAG {folder}: {str(e)}")
                    
        return rags
    
    def get_rag_documents(self, rag_name: str) -> List[DocumentInfo]:
        """Get documents from a specific RAG"""
        info_path = os.path.join(RLAMA_DATA_DIR, rag_name, "info.json")
        if not os.path.exists(info_path):
            raise HTTPException(status_code=404, detail=f"RAG '{rag_name}' not found")
            
        try:
            with open(info_path, "r") as f:
                info = json.load(f)
                
            # Log RAG info for debugging
            print(f"RAG '{rag_name}' info: Name={info.get('name')}, Model={info.get('model_name')}")
            print(f"Documents count in info.json: {len(info.get('documents', []))}")
                
            documents = []
            for doc in info.get("documents", []):
                doc_id = doc.get("id", "")
                doc_name = doc.get("name", "")
                doc_size = doc.get("size", 0)
                doc_type = doc.get("content_type", "text/plain")
                
                # Log document info for debugging
                print(f"Document found: ID={doc_id}, Name={doc_name}, Size={doc_size}, Type={doc_type}")
                
                documents.append(DocumentInfo(
                    id=doc_id,
                    name=doc_name,
                    size=self._format_size(doc_size),
                    content_type=doc_type
                ))
                
            return documents
        except json.JSONDecodeError as e:
            print(f"JSON decode error for '{info_path}': {str(e)}")
            raise HTTPException(status_code=500, detail=f"Error retrieving documents: {str(e)}")
        except Exception as e:
            print(f"Error retrieving documents for RAG '{rag_name}': {str(e)}")
            raise HTTPException(status_code=500, detail=f"Error retrieving documents: {str(e)}")
    
    def get_rag_chunks(self, rag_name: str, document_filter: Optional[str] = None, show_content: bool = False) -> List[ChunkInfo]:
        """Get chunks from a RAG with optional filtering"""
        info_path = os.path.join(RLAMA_DATA_DIR, rag_name, "info.json")
        if not os.path.exists(info_path):
            raise HTTPException(status_code=404, detail=f"RAG '{rag_name}' not found")
            
        try:
            with open(info_path, "r") as f:
                info = json.load(f)
            
            # Log chunks info for debugging
            print(f"Loading chunks for RAG '{rag_name}', filter: '{document_filter}', show_content: {show_content}")
            print(f"Chunks count in info.json: {len(info.get('chunks', []))}")
                
            chunks = []
            for chunk in info.get("chunks", []):
                # Appliquer le filtre sur le document ID
                doc_id = chunk.get("document_id", "")
                if document_filter and document_filter not in doc_id:
                    continue
                
                chunk_id = chunk.get("id", "")
                position = chunk.get("metadata", {}).get("position", "N/A")
                
                # Log chunk info for debugging
                print(f"Chunk found: ID={chunk_id}, DocID={doc_id}, Position={position}")
                    
                chunk_info = ChunkInfo(
                    id=chunk_id,
                    document_id=doc_id,
                    position=position
                )
                
                if show_content:
                    chunk_info.content = chunk.get("content", "")
                    
                chunks.append(chunk_info)
                
            print(f"Returning {len(chunks)} chunks after filtering")
            return chunks
        except json.JSONDecodeError as e:
            print(f"JSON decode error for '{info_path}': {str(e)}")
            raise HTTPException(status_code=500, detail=f"Error retrieving chunks: {str(e)}")
        except Exception as e:
            print(f"Error retrieving chunks for RAG '{rag_name}': {str(e)}")
            raise HTTPException(status_code=500, detail=f"Error retrieving chunks: {str(e)}")
    
    def create_rag(self, request: CreateRagRequest) -> dict:
        """Create a new RAG system"""
        # Vérifier si le RAG existe déjà
        rag_path = os.path.join(RLAMA_DATA_DIR, request.name)
        if os.path.exists(rag_path):
            raise HTTPException(status_code=400, detail=f"A RAG named '{request.name}' already exists")
            
        # Vérifier si le dossier source existe
        if not os.path.exists(request.folder_path):
            raise HTTPException(status_code=400, detail=f"The source folder '{request.folder_path}' does not exist")
        
        # Vérifier s'il y a des fichiers dans le dossier
        files = []
        for root, _, filenames in os.walk(request.folder_path):
            for filename in filenames:
                # Ignorer les fichiers cachés et les dossiers
                if not filename.startswith('.'):
                    files.append(os.path.join(root, filename))
        
        if not files:
            raise HTTPException(status_code=400, detail=f"The folder '{request.folder_path}' contains no visible files")
            
        print(f"Creating RAG '{request.name}' with model '{request.model}' from folder '{request.folder_path}'")
        print(f"Files found in folder: {len(files)}")
            
        # Préparer la commande
        cmd = [
            RLAMA_EXECUTABLE, "rag", 
            request.model, 
            request.name, 
            request.folder_path,
            "--chunk-size", str(request.chunk_size),
            "--chunk-overlap", str(request.chunk_overlap)
        ]
        
        if not request.enable_reranker:
            cmd.append("--disable-reranker")
        elif request.reranker_weight is not None:
            cmd.extend(["--reranker-weight", str(request.reranker_weight)])
        
        # Add profile if specified
        if request.profile:
            cmd.extend(["--profile", request.profile])
        
        print(f"RLAMA command: {' '.join(cmd)}")
        
        try:
            # Utiliser subprocess.run au lieu de Popen pour une meilleure gestion des erreurs
            result = subprocess.run(
                cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                check=False  # Ne pas lever d'exception, gérer le code de retour manuellement
            )
            
            print(f"RLAMA return code: {result.returncode}")
            print(f"RLAMA stdout: {result.stdout}")
            if result.stderr:
                print(f"RLAMA stderr: {result.stderr}")
                
            # Si le processus échoue
            if result.returncode != 0:
                error_message = result.stderr or result.stdout
                raise HTTPException(
                    status_code=500,
                    detail=f"Error creating RAG: {error_message}"
                )
                
            # Vérifier si le RAG a bien été créé
            if not os.path.exists(rag_path):
                raise HTTPException(
                    status_code=500,
                    detail=f"RAG '{request.name}' not created despite return code 0. Check logs."
                )
                
            info_path = os.path.join(rag_path, "info.json")
            if not os.path.exists(info_path):
                raise HTTPException(
                    status_code=500,
                    detail=f"Missing info.json file for RAG '{request.name}'"
                )
                
            # Lire le fichier info.json pour vérifier qu'il est valide
            try:
                with open(info_path, "r") as f:
                    info = json.load(f)
                doc_count = len(info.get('documents', []))
                chunk_count = len(info.get('chunks', []))
                print(f"RAG info loaded successfully: {info.get('name')}, {doc_count} documents, {chunk_count} chunks")
                
                # Si aucun document n'a été ajouté, forcer la création avec une meilleure commande
                if doc_count == 0 and files:
                    print("No document was added, trying with a specific file...")
                    # Essayer avec un fichier spécifique plutôt que le dossier entier
                    for file_path in files[:3]:  # Essayer avec les 3 premiers fichiers
                        add_cmd = [RLAMA_EXECUTABLE, "add-docs", request.name, file_path]
                        print(f"Tentative with the command: {' '.join(add_cmd)}")
                        add_result = subprocess.run(add_cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
                        print(f"Add result: {add_result.returncode}")
                        if add_result.returncode == 0:
                            print(f"Document {file_path} added successfully")
                            break
            except json.JSONDecodeError as e:
                raise HTTPException(
                    status_code=500,
                    detail=f"Corrupted info.json file for RAG '{request.name}': {str(e)}"
                )
                
            return {"message": "RAG created successfully", "name": request.name}
        except subprocess.CalledProcessError as e:
            print(f"RLAMA process error: returncode={e.returncode}")
            print(f"RLAMA stdout: {e.stdout}")
            print(f"RLAMA stderr: {e.stderr}")
            raise HTTPException(
                status_code=500,
                detail=f"Error creating RAG: {e.stderr or e.stdout}"
            )
        except Exception as e:
            print(f"Exception while creating RAG: {str(e)}")
            raise HTTPException(status_code=500, detail=f"Error creating RAG: {str(e)}")
    
    def add_documents(self, rag_name: str, folder_path: str) -> dict:
        """Add documents to an existing RAG"""
        # Vérifier si le RAG existe
        if not os.path.exists(os.path.join(RLAMA_DATA_DIR, rag_name)):
            raise HTTPException(status_code=404, detail=f"RAG '{rag_name}' not found")
            
        # Vérifier si le dossier source existe
        if not os.path.exists(folder_path):
            raise HTTPException(status_code=400, detail=f"The source folder '{folder_path}' does not exist")
            
        # Exécuter la commande RLAMA
        cmd = [RLAMA_EXECUTABLE, "add-docs", rag_name, folder_path]
        
        try:
            process = subprocess.Popen(
                cmd, 
                stdout=subprocess.PIPE, 
                stderr=subprocess.PIPE,
                text=True
            )
            
            stdout, stderr = process.communicate()
            
            if process.returncode != 0:
                raise HTTPException(
                    status_code=500, 
                    detail=f"Error adding documents: {stderr or stdout}"
                )
                
            return {"message": "Documents added successfully", "rag_name": rag_name}
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error adding documents: {str(e)}")
    
    def delete_rag(self, rag_name: str) -> dict:
        """Delete a RAG system"""
        rag_path = os.path.join(RLAMA_DATA_DIR, rag_name)
        
        # Vérifier si le RAG existe
        if not os.path.exists(rag_path):
            raise HTTPException(status_code=404, detail=f"RAG '{rag_name}' not found")
            
        try:
            # Supprimer le dossier du RAG
            shutil.rmtree(rag_path)
            return {"message": f"RAG '{rag_name}' deleted successfully"}
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error deleting RAG: {str(e)}")
    
    def update_model(self, rag_name: str, model_name: str) -> dict:
        """Update the LLM model of a RAG"""
        # Vérifier si le RAG existe
        info_path = os.path.join(RLAMA_DATA_DIR, rag_name, "info.json")
        if not os.path.exists(info_path):
            raise HTTPException(status_code=404, detail=f"RAG '{rag_name}' not found")
            
        # Exécuter la commande RLAMA
        cmd = [RLAMA_EXECUTABLE, "update-model", rag_name, model_name]
        
        try:
            process = subprocess.Popen(
                cmd, 
                stdout=subprocess.PIPE, 
                stderr=subprocess.PIPE,
                text=True
            )
            
            stdout, stderr = process.communicate()
            
            if process.returncode != 0:
                raise HTTPException(
                    status_code=500, 
                    detail=f"Error updating model: {stderr or stdout}"
                )
                
            return {"message": f"Model updated successfully to '{model_name}'", "rag_name": rag_name}
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error updating model: {str(e)}")
    
    def get_watch_status(self, rag_name: str) -> WatchStatus:
        """Get the folder watch status of a RAG"""
        rag_path = os.path.join(RLAMA_DATA_DIR, rag_name)
        watch_file = os.path.join(rag_path, ".watch")
        
        if not os.path.exists(rag_path):
            raise HTTPException(status_code=404, detail=f"RAG '{rag_name}' not found")
        
        # Vérifier si un fichier .watch existe
        if os.path.exists(watch_file):
            try:
                with open(watch_file, "r") as f:
                    watch_data = json.load(f)
                    
                return WatchStatus(
                    active=True,
                    folder_path=watch_data.get("folder_path", ""),
                    interval=watch_data.get("interval", 0),
                    last_check=watch_data.get("last_check", "")
                )
            except Exception as e:
                print(f"Error reading watch file: {str(e)}")
            
        # Si le fichier n'existe pas ou il y a une erreur
        return WatchStatus(active=False)

    def get_web_watch_status(self, rag_name: str) -> WebWatchStatus:
        """Get the web watch status of a RAG"""
        rag_path = os.path.join(RLAMA_DATA_DIR, rag_name)
        web_watch_file = os.path.join(rag_path, ".web-watch")
        
        if not os.path.exists(rag_path):
            raise HTTPException(status_code=404, detail=f"RAG '{rag_name}' not found")
        
        # Vérifier si un fichier .web-watch existe
        if os.path.exists(web_watch_file):
            try:
                with open(web_watch_file, "r") as f:
                    web_watch_data = json.load(f)
                    
                return WebWatchStatus(
                    active=True,
                    url=web_watch_data.get("url", ""),
                    interval=web_watch_data.get("interval", 0),
                    depth=web_watch_data.get("depth", 1),
                    last_check=web_watch_data.get("last_check", "")
                )
            except Exception as e:
                print(f"Error reading web watch file: {str(e)}")
            
        # Si le fichier n'existe pas ou il y a une erreur
        return WebWatchStatus(active=False)

    def setup_watch(self, request: WatchRequest) -> dict:
        """Configure folder watch for a RAG"""
        # Vérifier si le RAG existe
        rag_path = os.path.join(RLAMA_DATA_DIR, request.rag_name)
        if not os.path.exists(rag_path):
            raise HTTPException(status_code=404, detail=f"RAG '{request.rag_name}' not found")
        
        # Vérifier si le dossier source existe
        if not os.path.exists(request.folder_path):
            raise HTTPException(status_code=400, detail=f"The source folder '{request.folder_path}' does not exist")
        
        # Exécuter la commande RLAMA
        cmd = [RLAMA_EXECUTABLE, "watch", request.rag_name, request.folder_path, str(request.interval)]
        
        try:
            process = subprocess.Popen(
                cmd, 
                stdout=subprocess.PIPE, 
                stderr=subprocess.PIPE,
                text=True
            )
            
            stdout, stderr = process.communicate()
            
            if process.returncode != 0:
                raise HTTPException(
                    status_code=500, 
                    detail=f"Error configuring watch: {stderr or stdout}"
                )
            
            # Sauvegarder les informations de surveillance dans un fichier caché
            watch_info = {
                "folder_path": request.folder_path,
                "interval": request.interval,
                "last_check": time.strftime("%Y-%m-%d %H:%M:%S")
            }
            
            try:
                with open(os.path.join(rag_path, ".watch"), "w") as f:
                    json.dump(watch_info, f)
            except Exception as e:
                print(f"Error saving watch info: {str(e)}")
            
            return {
                "message": "Folder watch configured successfully", 
                "rag_name": request.rag_name,
                "interval": request.interval
            }
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error configuring watch: {str(e)}")

    def disable_watch(self, rag_name: str) -> dict:
        """Disable folder watch"""
        # Vérifier si le RAG existe
        rag_path = os.path.join(RLAMA_DATA_DIR, rag_name)
        if not os.path.exists(rag_path):
            raise HTTPException(status_code=404, detail=f"RAG '{rag_name}' not found")
        
        # Exécuter la commande RLAMA
        cmd = [RLAMA_EXECUTABLE, "watch-off", rag_name]
        
        try:
            process = subprocess.Popen(
                cmd, 
                stdout=subprocess.PIPE, 
                stderr=subprocess.PIPE,
                text=True
            )
            
            stdout, stderr = process.communicate()
            
            if process.returncode != 0:
                raise HTTPException(
                    status_code=500, 
                    detail=f"Error disabling watch: {stderr or stdout}"
                )
            
            # Supprimer le fichier .watch
            watch_file = os.path.join(rag_path, ".watch")
            if os.path.exists(watch_file):
                try:
                    os.remove(watch_file)
                except Exception as e:
                    print(f"Error deleting watch file: {str(e)}")
            
            return {"message": "Folder watch disabled successfully", "rag_name": rag_name}
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error disabling watch: {str(e)}")

    def setup_web_watch(self, request: WebWatchRequest) -> dict:
        """Configure web watch for a RAG"""
        # Vérifier si le RAG existe
        rag_path = os.path.join(RLAMA_DATA_DIR, request.rag_name)
        if not os.path.exists(rag_path):
            raise HTTPException(status_code=404, detail=f"RAG '{request.rag_name}' not found")
        
        # Exécuter la commande RLAMA
        cmd = [RLAMA_EXECUTABLE, "web-watch", request.rag_name, request.url, str(request.interval)]
        
        if request.depth:
            cmd.extend(["--depth", str(request.depth)])
        
        try:
            process = subprocess.Popen(
                cmd, 
                stdout=subprocess.PIPE, 
                stderr=subprocess.PIPE,
                text=True
            )
            
            stdout, stderr = process.communicate()
            
            if process.returncode != 0:
                raise HTTPException(
                    status_code=500, 
                    detail=f"Error configuring web watch: {stderr or stdout}"
                )
            
            # Sauvegarder les informations de surveillance dans un fichier caché
            web_watch_info = {
                "url": request.url,
                "interval": request.interval,
                "depth": request.depth or 1,
                "last_check": time.strftime("%Y-%m-%d %H:%M:%S")
            }
            
            try:
                with open(os.path.join(rag_path, ".web-watch"), "w") as f:
                    json.dump(web_watch_info, f)
            except Exception as e:
                print(f"Error saving web watch info: {str(e)}")
            
            return {
                "message": "Web watch configured successfully", 
                "rag_name": request.rag_name,
                "interval": request.interval
            }
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error configuring web watch: {str(e)}")

    def disable_web_watch(self, rag_name: str) -> dict:
        """Disable web watch"""
        # Vérifier si le RAG existe
        rag_path = os.path.join(RLAMA_DATA_DIR, rag_name)
        if not os.path.exists(rag_path):
            raise HTTPException(status_code=404, detail=f"RAG '{rag_name}' not found")
        
        # Exécuter la commande RLAMA
        cmd = [RLAMA_EXECUTABLE, "web-watch-off", rag_name]
        
        try:
            process = subprocess.Popen(
                cmd, 
                stdout=subprocess.PIPE, 
                stderr=subprocess.PIPE,
                text=True
            )
            
            stdout, stderr = process.communicate()
            
            if process.returncode != 0:
                raise HTTPException(
                    status_code=500, 
                    detail=f"Error disabling web watch: {stderr or stdout}"
                )
            
            # Supprimer le fichier .web-watch
            web_watch_file = os.path.join(rag_path, ".web-watch")
            if os.path.exists(web_watch_file):
                try:
                    os.remove(web_watch_file)
                except Exception as e:
                    print(f"Error deleting web watch file: {str(e)}")
            
            return {"message": "Web watch disabled successfully", "rag_name": rag_name}
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error disabling web watch: {str(e)}")

    def check_watched(self, rag_name: str) -> dict:
        """Force the verification of watched folders/sites"""
        # Vérifier si le RAG existe
        if not os.path.exists(os.path.join(RLAMA_DATA_DIR, rag_name)):
            raise HTTPException(status_code=404, detail=f"RAG '{rag_name}' not found")
            
        # Exécuter la commande RLAMA
        cmd = [RLAMA_EXECUTABLE, "check-watched", rag_name]
        
        try:
            process = subprocess.Popen(
                cmd, 
                stdout=subprocess.PIPE, 
                stderr=subprocess.PIPE,
                text=True
            )
            
            stdout, stderr = process.communicate()
            
            if process.returncode != 0:
                raise HTTPException(
                    status_code=500, 
                    detail=f"Error checking watched resources: {stderr or stdout}"
                )
                
            return {"message": "Verification completed successfully", "rag_name": rag_name}
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error checking: {str(e)}")
    
    def query_rag(self, request: QueryRequest) -> dict:
        """Query a RAG with a question"""
        # Vérifier si le RAG existe
        if not os.path.exists(os.path.join(RLAMA_DATA_DIR, request.rag_name)):
            raise HTTPException(status_code=404, detail=f"RAG '{request.rag_name}' not found")
            
        # Pour RLAMA version 0.1.35, la commande la plus fiable est:
        # rlama run [rag-name] --prompt "question"
        # Streaming must be disabled for capture
        cmd = [
            RLAMA_EXECUTABLE, "run", 
            request.rag_name, 
            "--prompt", request.prompt,
            "--stream=false" # For capture all output at once
        ]
        
        # Add context-size option if specified
        if request.context_size:
            cmd.extend(["--context-size", str(request.context_size)])
            
        # RLAMA doesn't support changing model at query time - Models are set when creating the RAG
        # or using the update-model command
        # if request.model:
        #     cmd.extend(["--model", request.model])
            
        print(f"RLAMA query command: {' '.join(cmd)}")
        
        # Analyze output to detect streaming progress
        def parse_rlama_output(output):
            # Find where answer starts
            answer_start = output.find("--- Answer ---")
            if answer_start >= 0:
                return output[answer_start + 14:].strip()
            return output.strip()
            
        try:
            # Given that LLM can take time,
            # we use a more direct approach with increased timeout
            result = subprocess.run(
                cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                timeout=90  # 90 seconds timeout
            )
            
            print(f"RLAMA query return code: {result.returncode}")
            
            if result.stderr:
                print(f"RLAMA query stderr: {result.stderr}")
                
            # Clean and analyze output
            if result.stdout:
                print(f"RLAMA query stdout: {result.stdout[:200]}... (truncated)")
                clean_answer = parse_rlama_output(result.stdout)
                return {"answer": clean_answer, "rag_name": request.rag_name}
            
            # If no output but no error, try an alternative
            if result.returncode == 0 and not result.stdout.strip():
                # Try directly with CLI without streaming
                cmd_direct = f"{RLAMA_EXECUTABLE} run {request.rag_name} --prompt \"{request.prompt}\" --stream=false"
                if request.context_size:
                    cmd_direct += f" --context-size {request.context_size}"
                
                print(f"Trying direct command: {cmd_direct}")
                direct_result = os.popen(cmd_direct).read()
                
                if direct_result.strip():
                    clean_answer = parse_rlama_output(direct_result)
                    return {"answer": clean_answer, "rag_name": request.rag_name}
            
            # If error, check common messages
            if result.returncode != 0:
                # Check known errors
                known_errors = {
                    "No documents found": "No document found in this RAG. Please add documents before asking questions.",
                    "no chunks": "No chunk found for this search. Try a different question.",
                    "unknown flag: --stream": "RLAMA version incompatible. Version 0.1.35+ required."
                }
                
                for key, message in known_errors.items():
                    if key in result.stderr:
                        return {"answer": message, "rag_name": request.rag_name}
                
                # Otherwise return raw error but formatted
                return {"answer": f"Error: {result.stderr.strip()}", "rag_name": request.rag_name}
            
            # If none of the above, return generic message
            return {
                "answer": "The model did not generate a response. Please try with a different question.",
                "rag_name": request.rag_name
            }
            
        except subprocess.TimeoutExpired:
            print("RLAMA query timeout after 90 seconds")
            return {
                "answer": "Error: The query took too long. Please try with a smaller model or reduce the context size.",
                "rag_name": request.rag_name
            }
        except Exception as e:
            error_msg = str(e)
            print(f"Exception during query: {error_msg}")
            
            if "timeout" in error_msg.lower():
                return {
                    "answer": "Error: timeout of 60000ms exceeded - The query took too long. Please try with a smaller model or reduce the context size.",
                    "rag_name": request.rag_name
                }
                
            return {
                "answer": f"Error: {error_msg}",
                "rag_name": request.rag_name
            }
    
    async def query_rag_stream(self, request_data: QueryRequest, fastapi_request: Request):
        """Query a RAG with a question and stream progress."""
        rag_path = os.path.join(RLAMA_DATA_DIR, request_data.rag_name)
        if not os.path.exists(rag_path):
            # This initial check can raise HTTPException before stream starts
            raise HTTPException(status_code=404, detail=f"RAG '{request_data.rag_name}' not found")

        cmd = [
            RLAMA_EXECUTABLE, "run",
            request_data.rag_name,
            "--prompt", request_data.prompt
        ]
        # Ensure RLAMA streams output. If '--stream=false' was used, remove it.
        # If RLAMA requires an explicit flag like '--stream=true', add it here.
        # Assuming 'rlama run' streams by default when not explicitly set to false.

        if request_data.context_size:
            cmd.extend(["--context-size", str(request_data.context_size)])
        # Remove model param - RLAMA doesn't support changing model at query time
        # if request_data.model: # If a different model is selected for this query
        #     cmd.extend(["--model", request_data.model])

        print(f"Streaming RLAMA query command: {' '.join(cmd)}")

        async def event_generator():
            process = None
            try:
                process = await asyncio.create_subprocess_exec(
                    *cmd,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE
                )

                full_answer_content = ""
                answer_started = False
                answer_delimiter = "--- Answer ---" # RLAMA CLI specific delimiter

                # Stream stdout
                while True:
                    if await fastapi_request.is_disconnected():
                        print("Client disconnected, terminating RLAMA process.")
                        if process and process.returncode is None:
                            process.terminate()
                            await process.wait()
                        break
                    
                    line_bytes = await process.stdout.readline()
                    if not line_bytes:
                        break
                    
                    line = line_bytes.decode('utf-8', errors='replace').strip()

                    if not line: # Skip empty lines
                        continue

                    if answer_delimiter in line:
                        answer_started = True
                        # Send any content before the delimiter as progress
                        progress_part = line.split(answer_delimiter, 1)[0].strip()
                        if progress_part:
                            progress_json = json.dumps({'type': 'progress', 'content': progress_part})
                            yield f"data: {progress_json}\r\n\r\n"
                        
                        # The rest of the line might be the start of the answer
                        answer_part = line.split(answer_delimiter, 1)[1].strip()
                        if answer_part:
                             full_answer_content += answer_part + "\n"
                             # Clean tags <think> before sending to frontend
                             if "<think>" in answer_part:
                                 # If the line contains both an opening and a closing
                                 if "</think>" in answer_part:
                                     think_start = answer_part.find("<think>")
                                     think_end = answer_part.find("</think>") + 8  # +8 to include closing tag
                                     
                                     # Take only parts before and after think tags
                                     clean_line = answer_part[:think_start].strip() + " " + answer_part[think_end:].strip()
                                     if clean_line.strip():  # Verify cleaned line is not empty
                                         answer_json = json.dumps({'type': 'answer_chunk', 'content': clean_line + '\n'})
                                         yield f"data: {answer_json}\r\n\r\n"
                                 else:
                                     # If it's just an opening tag, ignore the whole line
                                     continue
                             elif "</think>" in answer_part:
                                 # If it's just a closing tag, ignore the whole line
                                 continue
                             else:
                                 # If no think tags, send normally
                                 answer_json = json.dumps({'type': 'answer_chunk', 'content': answer_part + '\n'})
                                 yield f"data: {answer_json}\r\n\r\n"
                        continue

                    if answer_started:
                        full_answer_content += line + "\n"
                        # Clean tags <think> before sending to frontend
                        if "<think>" in line:
                            # If the line contains both an opening and a closing
                            if "</think>" in line:
                                think_start = line.find("<think>")
                                think_end = line.find("</think>") + 8  # +8 to include closing tag
                                
                                # Take only parts before and after think tags
                                clean_line = line[:think_start].strip() + " " + line[think_end:].strip()
                                if clean_line.strip():  # Verify cleaned line is not empty
                                    answer_json = json.dumps({'type': 'answer_chunk', 'content': clean_line + '\n'})
                                    yield f"data: {answer_json}\r\n\r\n"
                            else:
                                # If it's just an opening tag, ignore the whole line
                                continue
                        elif "</think>" in line:
                            # If it's just a closing tag, ignore the whole line
                            continue
                        else:
                            # If no think tags, send normally
                            answer_json = json.dumps({'type': 'answer_chunk', 'content': line + '\n'})
                            yield f"data: {answer_json}\r\n\r\n"
                    else:
                        progress_json = json.dumps({'type': 'progress', 'content': line})
                        yield f"data: {progress_json}\r\n\r\n"
                
                # Capture any remaining stderr after stdout is exhausted
                stderr_output = ""
                if process.stderr:
                    stderr_bytes = await process.stderr.read()
                    stderr_output = stderr_bytes.decode('utf-8', errors='replace').strip()
                
                if stderr_output:
                    print(f"RLAMA stderr: {stderr_output}")
                    # Send as an error event if no answer was formed, or as additional info
                    if not answer_started or not full_answer_content.strip():
                         error_json = json.dumps({'type': 'error', 'content': f'RLAMA Error: {stderr_output}'})
                         yield f"data: {error_json}\r\n\r\n"
                    else: # If answer was already streamed, send as progress/warning
                         warn_json = json.dumps({'type': 'progress', 'content': f'RLAMA Info/Error: {stderr_output}'})
                         yield f"data: {warn_json}\r\n\r\n"

                await process.wait() # Ensure process has finished
                
                if process.returncode != 0 and not full_answer_content.strip() and not stderr_output:
                    error_json = json.dumps({'type': 'error', 'content': f'RLAMA process exited with code {process.returncode} but no specific error message.'})
                    yield f"data: {error_json}\r\n\r\n"
                elif not full_answer_content.strip() and not stderr_output and process.returncode == 0 and answer_started is False:
                     # If answer_started is true but full_answer_content is empty, it means an empty answer was given.
                     # This case is for when "--- Answer ---" was never seen.
                     error_json = json.dumps({'type': 'error', 'content': 'RLAMA process completed without error but no answer was produced or delimiter found.'})
                     yield f"data: {error_json}\r\n\r\n"

            except asyncio.CancelledError:
                print("RLAMA stream task cancelled, ensuring process cleanup.")
                if process and process.returncode is None:
                    process.terminate()
                    await process.wait()
                raise # Re-raise CancelledError
            except Exception as e:
                print(f"Error during RLAMA stream: {str(e)}")
                error_json = json.dumps({'type': 'error', 'content': f'Backend streaming error: {str(e)}'})
                yield f"data: {error_json}\r\n\r\n"
            finally:
                if process and process.returncode is None: # Ensure process is killed if generator exits unexpectedly
                    print("Generator exiting, ensuring RLAMA process cleanup.")
                    process.terminate()
                    await process.wait()
                done_json = json.dumps({'type': 'done'})
                yield f"data: {done_json}\r\n\r\n" # Signal end of stream

        return StreamingResponse(event_generator(), media_type="text/event-stream")

    def get_available_models(self) -> List[str]:
        """Get the list of models available via Ollama + OpenAI models"""
        models = []
        
        try:
            # Get Ollama models
            process = subprocess.Popen(
                ["ollama", "list"], 
                stdout=subprocess.PIPE, 
                stderr=subprocess.PIPE,
                text=True
            )
            
            stdout, _ = process.communicate()
            
            if process.returncode == 0:
                lines = stdout.strip().split('\n')
                
                # Ignore the first line (header)
                for line in lines[1:]:
                    if line.strip():
                        # The first column is the model name
                        model_name = line.strip().split()[0]
                        if model_name and not model_name.startswith("NAME"):
                            models.append(model_name)
            else:
                print("Ollama not available or no models installed")
                
        except Exception as e:
            print(f"Error retrieving Ollama models: {str(e)}")
        
        # Add current OpenAI models (2024)
        openai_models = [
            # GPT-4 models
            "gpt-4o",
            "gpt-4o-mini", 
            "gpt-4-turbo",
            "gpt-4",
            "gpt-4-32k",
            
            # GPT-3.5 models
            "gpt-3.5-turbo",
            "gpt-3.5-turbo-16k",
            
            # Legacy models
            "text-davinci-003",
            "text-davinci-002"
        ]
        
        # Add OpenAI models
        models.extend(openai_models)
        
        # Remove duplicates while preserving order
        seen = set()
        unique_models = []
        for model in models:
            if model not in seen:
                seen.add(model)
                unique_models.append(model)
        
        # If no models found, return default
        if not unique_models:
            unique_models = ["gpt-3.5-turbo", "llama2"]
            
        return unique_models
    
    def _format_size(self, size_bytes: int) -> str:
        """Format a size in bytes in a readable format"""
        if size_bytes < 1024:
            return f"{size_bytes} B"
        elif size_bytes < 1024 * 1024:
            return f"{size_bytes / 1024:.2f} KB"
        else:
            return f"{size_bytes / (1024 * 1024):.2f} MB"

# Service Agent
class AgentService:
    def __init__(self):
        self.active_tasks = {}  # Store active task progress
        
    def run_agent(self, request: AgentQueryRequest) -> dict:
        """Execute agent query and return real-time progress"""
        # Build RLAMA agent command
        cmd = [RLAMA_EXECUTABLE, "agent", "run"]
        
        # Add model if specified
        if request.model:
            cmd.extend(["-l", request.model])
            
        # Add RAG if specified
        if request.rag_name:
            cmd.append(request.rag_name)
            
        # Add web search if enabled
        if request.web_search:
            cmd.append("-w")
            
        # Add profile if specified
        if request.profile:
            cmd.extend(["--profile", request.profile])
            
        # Add query
        cmd.extend(["-q", request.query])
        
        print(f"Running agent command: {' '.join(cmd)}")
        
        try:
            # Execute command and capture output
            result = subprocess.run(
                cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                timeout=300  # 5 minutes timeout
            )
            
            if result.returncode != 0:
                raise HTTPException(
                    status_code=500,
                    detail=f"Agent execution failed: {result.stderr or result.stdout}"
                )
                
            return {
                "response": result.stdout,
                "tasks": [],  # Will be enhanced with real-time task tracking
                "metadata": {
                    "command": " ".join(cmd),
                    "execution_time": "N/A"
                }
            }
            
        except subprocess.TimeoutExpired:
            raise HTTPException(status_code=408, detail="Agent execution timeout")
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Agent execution error: {str(e)}")
    
    async def run_agent_stream(self, request: AgentQueryRequest, fastapi_request: Request):
        """Execute agent with real-time streaming progress"""
        import uuid
        import re
        
        # Generate session ID for this agent run
        session_id = str(uuid.uuid4())
        
        # Build command
        cmd = [RLAMA_EXECUTABLE, "agent", "run"]
        if request.model:
            cmd.extend(["-l", request.model])
        if request.rag_name:
            cmd.append(request.rag_name)
        if request.web_search:
            cmd.append("-w")
        if request.profile:
            cmd.extend(["--profile", request.profile])
        cmd.extend(["-q", request.query])
        
        async def event_generator():
            try:
                # Send initial status
                yield f"data: {json.dumps({'type': 'progress', 'content': {'message': 'Starting agent...', 'status': 'running'}})}\n\n"
                
                # Start the process
                process = subprocess.Popen(
                    cmd,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                    text=True,
                    bufsize=1,
                    universal_newlines=True
                )
                
                # Read output line by line
                current_tasks = {}
                final_response = ""
                
                while True:
                    if await fastapi_request.is_disconnected():
                        process.kill()
                        break
                        
                    output = process.stdout.readline()
                    if output == '' and process.poll() is not None:
                        break
                        
                    if output:
                        line = output.strip()
                        print(f"Agent output: {line}")
                        
                        # Parse different types of output
                        if "Orchestrator:" in line:
                            # Task orchestration messages
                            task_info = self._parse_orchestrator_message(line)
                            if task_info:
                                yield f"data: {json.dumps({'type': 'progress', 'content': task_info})}\n\n"
                                
                        elif "Task" in line and ("completed" in line or "failed" in line or "running" in line):
                            # Task status updates
                            task_update = self._parse_task_status(line)
                            if task_update:
                                yield f"data: {json.dumps({'type': 'task_update', 'content': task_update})}\n\n"
                                
                        elif line and not line.startswith("Orchestrator:") and not line.startswith("Task"):
                            # Accumulate final response
                            final_response += line + "\n"
                            # Send partial response
                            yield f"data: {json.dumps({'type': 'answer_chunk', 'content': line})}\n\n"
                
                # Wait for process to complete
                return_code = process.wait()
                
                if return_code != 0:
                    stderr_output = process.stderr.read()
                    yield f"data: {json.dumps({'type': 'error', 'content': f'Agent failed: {stderr_output}'})}\n\n"
                else:
                    # Send final response
                    yield f"data: {json.dumps({'type': 'done', 'content': {'response': final_response.strip()}})}\n\n"
                    
            except Exception as e:
                yield f"data: {json.dumps({'type': 'error', 'content': f'Stream error: {str(e)}'})}\n\n"
                
        return StreamingResponse(
            event_generator(),
            media_type="text/event-stream",
            headers={
                "Cache-Control": "no-cache",
                "Connection": "keep-alive",
                "Access-Control-Allow-Origin": "*",
                "Access-Control-Allow-Headers": "*",
            }
        )
    
    def _parse_orchestrator_message(self, line: str) -> Optional[dict]:
        """Parse orchestrator messages to extract task information"""
        if "Decomposing query" in line:
            return {"message": "Analyzing query and planning tasks...", "status": "planning"}
        elif "Decomposed into" in line and "tasks" in line:
            import re
            match = re.search(r'(\d+) tasks', line)
            if match:
                task_count = match.group(1)
                return {"message": f"Created {task_count} tasks", "status": "planning"}
        elif "Executing task" in line:
            return {"message": "Executing tasks...", "status": "executing"}
        elif "All tasks completed" in line:
            return {"message": "All tasks completed", "status": "completed"}
        elif "Generating final response" in line:
            return {"message": "Generating final response...", "status": "generating"}
        return None
    
    def _parse_task_status(self, line: str) -> Optional[dict]:
        """Parse task status updates"""
        import re
        
        # Match patterns like "Task task1: Description completed successfully"
        task_pattern = r'Task (\w+):\s*(.+?)\s*(completed|failed|running)'
        match = re.search(task_pattern, line, re.IGNORECASE)
        
        if match:
            task_id = match.group(1)
            description = match.group(2)
            status = match.group(3).lower()
            
            return {
                "task_id": task_id,
                "description": description,
                "status": status,
                "timestamp": time.time()
            }
        return None
    
    def get_available_models(self) -> List[str]:
        """Get available models for agents"""
        try:
            result = subprocess.run(
                [RLAMA_EXECUTABLE, "models"],
                capture_output=True,
                text=True
            )
            
            if result.returncode == 0:
                # Parse model output
                models = []
                for line in result.stdout.strip().split('\n'):
                    if line.strip():
                        models.append(line.strip())
                return models
            else:
                return []
        except Exception as e:
            print(f"Error getting models: {str(e)}")
            return []

# Initialize the services
rlama_service = RlamaService()
agent_service = AgentService()

# API routes
@app.get("/health")
def health_check():
    return {"status": "ok", "timestamp": time.time()}

@app.get("/rags", response_model=List[RagInfo])
def list_rags():
    return rlama_service.get_all_rags()

@app.get("/rags/{rag_name}/documents", response_model=List[DocumentInfo])
def get_documents(rag_name: str):
    return rlama_service.get_rag_documents(rag_name)

@app.get("/rags/{rag_name}/chunks", response_model=List[ChunkInfo])
def get_chunks(rag_name: str, document_filter: Optional[str] = None, show_content: bool = False):
    return rlama_service.get_rag_chunks(rag_name, document_filter, show_content)

@app.post("/rags")
def create_rag(request: CreateRagRequest):
    return rlama_service.create_rag(request)

@app.post("/rags/{rag_name}/documents")
def add_documents(rag_name: str, folder_path: str = Form(...)):
    return rlama_service.add_documents(rag_name, folder_path)

@app.delete("/rags/{rag_name}")
def delete_rag(rag_name: str):
    return rlama_service.delete_rag(rag_name)

@app.put("/rags/{rag_name}/model")
def update_model(rag_name: str, model_name: str = Form(...)):
    return rlama_service.update_model(rag_name, model_name)

@app.post("/rags/{rag_name}/watch")
def setup_watch(request: WatchRequest):
    return rlama_service.setup_watch(request)

@app.post("/rags/{rag_name}/web-watch")
def setup_web_watch(request: WebWatchRequest):
    return rlama_service.setup_web_watch(request)

@app.delete("/rags/{rag_name}/watch")
def disable_watch(rag_name: str):
    return rlama_service.disable_watch(rag_name)

@app.delete("/rags/{rag_name}/web-watch")
def disable_web_watch(rag_name: str):
    return rlama_service.disable_web_watch(rag_name)

@app.post("/rags/{rag_name}/check-watched")
def check_watched(rag_name: str):
    return rlama_service.check_watched(rag_name)

@app.post("/query")
def query_rag(request: QueryRequest):
    return rlama_service.query_rag(request)

@app.post("/query-stream")
async def stream_query_rag(request_data: QueryRequest, request: Request):
    return await rlama_service.query_rag_stream(request_data, request)

@app.get("/models")
def get_models():
    return {"models": rlama_service.get_available_models()}

@app.get("/exec", response_model=CommandResponse)
def execute_command(command: str):
    """Execute a shell command and return stdout/stderr"""
    # List of allowed commands for security
    ALLOWED_COMMANDS = ["rlama --version", "ollama --version", "ollama list", "rlama -v", "ollama -v"]
    
    # Verify if the command is allowed
    if command not in ALLOWED_COMMANDS:
        raise HTTPException(status_code=400, detail="Command not allowed")
    
    try:
        # For rlama, we use directly the command without the absolute path
        cmd_parts = command.split()
        if cmd_parts[0] == "rlama":
            # Use directly the rlama command (from PATH)
            cmd = cmd_parts
        else:
            # For other commands, use the command as is
            cmd = cmd_parts
            
        # Execute the command
        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            timeout=10  # 10 seconds timeout
        )
        
        # Log for debug
        print(f"Command executed: {command}")
        print(f"stdout: {result.stdout}")
        print(f"stderr: {result.stderr}")
        print(f"returncode: {result.returncode}")
        
        return CommandResponse(
            stdout=result.stdout,
            stderr=result.stderr,
            returncode=result.returncode
        )
    except subprocess.TimeoutExpired:
        raise HTTPException(status_code=408, detail="Command execution timeout")
    except Exception as e:
        print(f"Error executing command: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/rags/{rag_name}/watch-status", response_model=WatchStatus)
def get_watch_status(rag_name: str):
    return rlama_service.get_watch_status(rag_name)

@app.get("/rags/{rag_name}/web-watch-status", response_model=WebWatchStatus)
def get_web_watch_status(rag_name: str):
    return rlama_service.get_web_watch_status(rag_name)

# Agent endpoints
@app.post("/agent/run")
def run_agent(request: AgentQueryRequest):
    """Run agent query and get response"""
    return agent_service.run_agent(request)

@app.post("/agent/stream")
async def stream_agent(request: AgentQueryRequest, fastapi_request: Request):
    """Run agent with real-time streaming"""
    return await agent_service.run_agent_stream(request, fastapi_request)

@app.get("/agent/models")
def get_agent_models():
    """Get available models for agents"""
    return {"models": agent_service.get_available_models()}

# =======================================
# Endpoints pour la gestion des paramètres
# =======================================

# Modèles pour les paramètres
class ProfileRequest(BaseModel):
    name: str
    provider: str
    api_key: str
    description: Optional[str] = None

class ApiKeys(BaseModel):
    openai_api_key: Optional[str] = None
    openai_organization: Optional[str] = None
    google_api_key: Optional[str] = None
    google_search_engine_id: Optional[str] = None
    anthropic_api_key: Optional[str] = None

class GeneralSettingsRequest(BaseModel):
    auto_save: Optional[bool] = True
    show_notifications: Optional[bool] = True
    default_model: Optional[str] = "gpt-3.5-turbo"
    default_embedding_model: Optional[str] = "text-embedding-ada-002"

# =======================================
# Fonctions utilitaires pour les profils et paramètres
# =======================================

def load_profiles():
    """Load all API profiles from files"""
    profiles = []
    
    try:
        # First try to get profiles from RLAMA CLI
        result = subprocess.run(
            ["rlama", "profile", "list"],
            capture_output=True,
            text=True
        )
        
        if result.returncode == 0 and result.stdout:
            # Parse RLAMA CLI output
            lines = result.stdout.strip().split('\n')
            header_found = False
            
            for line in lines:
                line = line.strip()
                
                # Skip empty lines and header lines
                if not line or line.startswith("Available") or line.startswith("NAME"):
                    if line.startswith("NAME"):
                        header_found = True
                    continue
                
                # Only process data lines after header is found
                if header_found and line:
                    # Split by multiple spaces to handle proper column separation
                    import re
                    parts = re.split(r'\s{2,}', line)  # Split on 2+ spaces
                    
                    if len(parts) >= 4:
                        name = parts[0].strip()
                        provider = parts[1].strip()
                        created_on = parts[2].strip()
                        last_used = parts[3].strip()
                        
                        profiles.append({
                            "name": name,
                            "provider": provider,
                            "created_on": created_on,
                            "last_used": last_used
                        })
                    elif len(parts) >= 3:
                        # Fallback for lines that might have fewer columns
                        name = parts[0].strip()
                        provider = parts[1].strip()
                        # Try to parse the rest as date and last_used
                        remaining = ' '.join(parts[2:]).strip()
                        # Simple heuristic: if 'never' is at the end, split there
                        if remaining.endswith('never'):
                            created_on = remaining[:-5].strip()
                            last_used = 'never'
                        else:
                            # Assume the last word/phrase is last_used
                            parts_remaining = remaining.rsplit(' ', 1)
                            if len(parts_remaining) == 2:
                                created_on = parts_remaining[0].strip()
                                last_used = parts_remaining[1].strip()
                            else:
                                created_on = remaining
                                last_used = 'unknown'
                        
                        profiles.append({
                            "name": name,
                            "provider": provider,
                            "created_on": created_on,
                            "last_used": last_used
                        })
            
        # Also check local profile files
        if os.path.exists(PROFILES_DIR):
            for filename in os.listdir(PROFILES_DIR):
                if filename.endswith('.json'):
                    profile_name = filename[:-5]  # Remove .json extension
                    
                    # Check if this profile is already in the list from CLI
                    existing = any(p['name'] == profile_name for p in profiles)
                    if not existing:
                        profile_path = os.path.join(PROFILES_DIR, filename)
                        try:
                            with open(profile_path, 'r') as f:
                                profile_data = json.load(f)
                                profiles.append({
                                    "name": profile_data.get("name", profile_name),
                                    "provider": profile_data.get("provider", "unknown"),
                                    "created_on": profile_data.get("created_on", "unknown"),
                                    "last_used": profile_data.get("last_used", "never")
                                })
                        except Exception as e:
                            print(f"Error loading profile {filename}: {e}")
    
    except Exception as e:
        print(f"Error loading profiles: {e}")
        # Fallback to empty list
        profiles = []
    
    print(f"Debug: Found {len(profiles)} profiles: {[p['name'] for p in profiles]}")
    return {"profiles": profiles}

def save_profile(profile: ProfileRequest):
    """Save a profile both locally and via RLAMA CLI"""
    try:
        # First, try to save via RLAMA CLI
        result = subprocess.run(
            ["rlama", "profile", "add", profile.name, profile.provider, profile.api_key],
            capture_output=True,
            text=True,
            input="\n"  # Auto-confirm if needed
        )
        
        if result.returncode != 0:
            print(f"RLAMA CLI profile creation failed: {result.stderr}")
            # Continue to save locally anyway
        
        # Also save locally for backup/additional info
        profile_file = os.path.join(PROFILES_DIR, f"{profile.name}.json")
        profile_data = {
            "name": profile.name,
            "provider": profile.provider,
            "api_key": profile.api_key,  # Note: In production, encrypt this
            "description": profile.description,
            "created_on": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
            "last_used": "never"
        }
        
        with open(profile_file, 'w') as f:
            json.dump(profile_data, f, indent=2)
            
    except Exception as e:
        raise Exception(f"Failed to save profile: {str(e)}")

def load_api_keys():
    """Load API keys from settings file"""
    try:
        if os.path.exists(API_KEYS_FILE):
            with open(API_KEYS_FILE, 'r') as f:
                return json.load(f)
        else:
            # Return default empty structure
            return {
                "openai_api_key": None,
                "openai_organization": None,
                "google_api_key": None,
                "google_search_engine_id": None,
                "anthropic_api_key": None
            }
    except Exception as e:
        print(f"Error loading API keys: {e}")
        return {
            "openai_api_key": None,
            "openai_organization": None,
            "google_api_key": None,
            "google_search_engine_id": None,
            "anthropic_api_key": None
        }

def save_api_keys(api_keys: ApiKeys):
    """Save API keys to settings file"""
    try:
        # Convert Pydantic model to dict, filtering out None values
        keys_dict = {}
        for key, value in api_keys.dict().items():
            if value is not None:
                keys_dict[key] = value
        
        # Load existing keys if file exists
        existing_keys = {}
        if os.path.exists(API_KEYS_FILE):
            try:
                with open(API_KEYS_FILE, 'r') as f:
                    existing_keys = json.load(f)
            except:
                existing_keys = {}
        
        # Merge with existing keys (only update provided ones)
        existing_keys.update(keys_dict)
        
        # Save to file
        with open(API_KEYS_FILE, 'w') as f:
            json.dump(existing_keys, f, indent=2)
            
    except Exception as e:
        raise Exception(f"Failed to save API keys: {str(e)}")

@app.get("/profiles")
def get_profiles():
    """Get all API profiles"""
    return load_profiles()

@app.post("/profiles")
def create_profile(request: ProfileRequest):
    """Create a new API profile"""
    try:
        save_profile(request)
        return {"message": f"Profile '{request.name}' created successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.delete("/profiles/{profile_name}")
def delete_profile(profile_name: str):
    """Delete an API profile"""
    try:
        profile_file = os.path.join(PROFILES_DIR, f"{profile_name}.json")
        if os.path.exists(profile_file):
            os.remove(profile_file)
            
        # Also try CLI deletion
        try:
            subprocess.run(
                [RLAMA_EXECUTABLE, "profile", "delete", profile_name],
                input="y\n",
                capture_output=True,
                text=True
            )
        except Exception as e:
            print(f"Warning: CLI profile deletion failed (but file deleted): {e}")
            
        return {"message": f"Profile '{profile_name}' deleted successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/settings/api-keys")
def get_api_keys():
    """Get stored API keys"""
    return load_api_keys()

@app.post("/settings/api-keys")
def save_api_keys_endpoint(api_keys: ApiKeys):
    """Save API keys"""
    try:
        save_api_keys(api_keys)
        return {"message": "API keys saved successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/settings/general")
def get_general_settings():
    """Get general settings"""
    try:
        settings_file = os.path.join(SETTINGS_DIR, "general.json")
        if os.path.exists(settings_file):
            with open(settings_file, "r") as f:
                return json.load(f)
        
        # Retourner paramètres par défaut
        return {
            "auto_save": True,
            "show_notifications": True,
            "default_model": "gpt-3.5-turbo",
            "default_embedding_model": "text-embedding-ada-002"
        }
    except Exception as e:
        print(f"Error getting general settings: {str(e)}")
        return {
            "auto_save": True,
            "show_notifications": True,
            "default_model": "gpt-3.5-turbo",
            "default_embedding_model": "text-embedding-ada-002"
        }

@app.post("/settings/general")
def save_general_settings(request: GeneralSettingsRequest):
    """Save general settings"""
    try:
        settings_file = os.path.join(SETTINGS_DIR, "general.json")
        
        with open(settings_file, "w") as f:
            json.dump(request.dict(), f, indent=2)
        
        return {"message": "General settings saved successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error saving general settings: {str(e)}")

# Start server when this file is executed directly
if __name__ == "__main__":
    print("Backend server started")
    try:
        uvicorn.run(app, host="127.0.0.1", port=5001, log_level="info")
    except Exception as e:
        print(f"ERROR: Failed to start backend server: {e}") 