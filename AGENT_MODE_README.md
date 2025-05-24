# Agent Mode - Orchestration des Tâches Complexes

## Vue d'ensemble

Le mode **Agent Orchestré** est une nouvelle fonctionnalité de rlama qui permet de décomposer automatiquement des requêtes complexes en sous-tâches simples et de les exécuter dans le bon ordre. C'est exactement comme vous et Cursor travaillez ensemble !

## Fonctionnalités

### 🎯 Décomposition Intelligente des Tâches
L'orchestrateur analyse votre requête et la décompose automatiquement en tâches plus petites :

**Exemple de requête complexe :**
```
"Quand se passe le Snowflake Summit prochainement et où ? Combien ça coûterait pour y participer, aussi cherche-moi des billets d'avions entre Montréal et l'endroit de destination STP"
```

**Décomposition automatique :**
1. **Tâche 1** : Trouver la date et le lieu du prochain Snowflake Summit
2. **Tâche 2** : Trouver le coût de participation au Snowflake Summit
3. **Tâche 3** : Chercher des billets d'avion entre Montréal et la destination
4. **Tâche 4** : Fournir une réponse complète et bien structurée

### 🔧 Outils Disponibles

L'agent orchestré a accès à tous les outils suivants :

- **`web_search`** : Recherche web en temps réel (événements, prix, actualités)
- **`flight_search`** : Recherche de vols entre deux villes
- **`rag_search`** : Recherche dans les bases de connaissances locales
- **`list_dir`** : Liste le contenu des répertoires
- **`read_file`** : Lit le contenu des fichiers
- **`file_search`** : Recherche de fichiers par nom (fuzzy matching)
- **`grep_search`** : Recherche de texte exact dans les fichiers
- **`codebase_search`** : Recherche sémantique dans le code

### 🧠 Intelligence Adaptative

L'agent détermine automatiquement si une requête est :
- **SIMPLE** : Peut être résolue avec un seul outil
- **COMPLEXE** : Nécessite une orchestration de plusieurs tâches

## Utilisation

### Mode par Défaut (Orchestré)
```bash
rlama agent run -w -q "Quand se passe le Snowflake Summit et combien ça coûte depuis Montréal ?"
```

### Mode Explicite
```bash
rlama agent run -w -m orchestrated -q "Trouve les vulnérabilités Python dans mon code et suggère des corrections"
```

### Avec une Base de Connaissances RAG
```bash
rlama agent run my-rag -w -m orchestrated -q "Analyse mes documents et trouve des informations sur la sécurité"
```

## Exemples de Requêtes Complexes

### 1. Planification d'Événement
```bash
rlama agent run -w -q "Trouve la prochaine conférence tech à Paris, le coût d'inscription, et des vols depuis Montréal"
```

**Décomposition automatique :**
- Recherche de conférences tech à Paris
- Recherche des coûts d'inscription
- Recherche de vols Montréal → Paris
- Synthèse complète

### 2. Analyse de Code
```bash
rlama agent run -w -q "Liste les fichiers Python dans mon projet, trouve les problèmes de sécurité, et suggère des améliorations"
```

**Décomposition automatique :**
- Liste des fichiers Python
- Analyse de sécurité du code
- Suggestions d'amélioration
- Rapport consolidé

### 3. Recherche et Analyse
```bash
rlama agent run -w -q "Recherche les dernières tendances en IA, analyse mes documents sur l'IA, et compare les approches"
```

**Décomposition automatique :**
- Recherche web des tendances IA
- Recherche RAG dans les documents locaux
- Analyse comparative
- Synthèse des résultats

## Configuration

### Variables d'Environnement Requises

Pour utiliser la recherche web et les vols :
```bash
export GOOGLE_SEARCH_API_KEY="votre_clé_api"
export GOOGLE_SEARCH_ENGINE_ID="votre_engine_id"
```

### Obtenir les Clés API Google

1. Allez sur [Google Cloud Console](https://console.cloud.google.com/)
2. Créez un projet ou sélectionnez-en un existant
3. Activez l'API "Custom Search API"
4. Créez des identifiants (clé API)
5. Créez un moteur de recherche personnalisé sur [cse.google.com](https://cse.google.com/)

## Modes Disponibles

| Mode | Description | Utilisation |
|------|-------------|-------------|
| `orchestrated` | **Défaut** - Décompose les requêtes complexes | Requêtes multi-étapes |
| `conversation` | Mode conversationnel simple | Requêtes simples |
| `autonomous` | Mode autonome (non implémenté) | Futur |

## Debug et Développement

Pour voir le processus de décomposition en détail :
```bash
rlama agent run -w -v -q "votre requête complexe"
```

Le flag `-v` (verbose) active le mode debug qui montre :
- L'analyse de complexité de la requête
- La décomposition en tâches
- L'exécution de chaque tâche
- Les dépendances entre tâches

## Architecture

```
Requête Utilisateur
        ↓
Analyse de Complexité
        ↓
[SIMPLE] → Exécution Directe
        ↓
[COMPLEXE] → Orchestrateur
        ↓
Décomposition en Tâches
        ↓
Exécution Séquentielle (avec dépendances)
        ↓
Synthèse Finale
        ↓
Réponse Utilisateur
```

## Avantages

✅ **Automatique** : Pas besoin de décomposer manuellement les tâches
✅ **Intelligent** : Comprend les dépendances entre tâches
✅ **Flexible** : S'adapte à différents types de requêtes
✅ **Robuste** : Continue même si certaines tâches échouent
✅ **Transparent** : Mode debug pour comprendre le processus

## Exemples Pratiques

### Planification de Voyage
```bash
rlama agent run -w -q "Je veux aller à la prochaine DockerCon. Trouve quand et où c'est, combien ça coûte, et des vols depuis Québec"
```

### Analyse de Projet
```bash
rlama agent run -w -q "Analyse mon projet Go, trouve les fichiers de configuration, vérifie les vulnérabilités, et suggère des améliorations"
```

### Recherche Comparative
```bash
rlama agent run -w -q "Compare les frameworks JavaScript populaires en 2024, trouve leurs avantages/inconvénients, et recommande le meilleur pour mon projet"
```

---

**Note** : Le mode orchestré est maintenant le mode par défaut. Pour les requêtes simples, l'agent détectera automatiquement qu'aucune orchestration n'est nécessaire et utilisera l'approche directe. 