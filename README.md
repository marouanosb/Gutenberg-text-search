# Moteur de recherche de livre GUTENBERG
- Anis HALFAOUI
- Merouane BOUAFIA

- http://217.182.170.152:8080/
  
# 1. Installation

### 1.1. Windows
Placez-vous dans le dossier `backend`
```bash   
python -m venv env
./env/Scripts/activate
pip install -r req.txt
```
### 1.2. Macos/Linux
Placez-vous dans le dossier `backend`
``` bash   
python3.11 -m venv env
source ./env/bin/activate
pip install -r req.txt
```
Postgres doit être installé sur votre machine (si vous souhaitez l'utiliser).

``` bash   
brew install postgresql

```
La version de Python doit être `< 3.12` et `> 3.6`.

## 2. Backend
Arborescence
 ```
 -> backend
    -> backend
    -> books 
    -> data
    -> keywords
    
```
Remarque : les dossiers `books` / `keywords` sont ignorés par Git. Le plus important est le fichier `db.sqlite3` qui contient toutes les données des livres.
### 2.0. Premiers pas
- Si vous utilisez vos propres données (livres) ou votre propre base de données, exécutez ces commandes (ainsi que celles de la section `2.1`). Si vous gardez les données par défaut, passez directement à la section `2.2`.

Dans le dossier backend
```bash   
python manage.py makemigrations
python manage.py migrate
```
### 2.1. Commandes à exécuter
Ces commandes Django doivent être exécutées dans cet ordre précis pour que le projet fonctionne correctement. Elles peuvent prendre du temps selon votre CPU / débit Internet — prenez un café (et appelez vos amis si vous en avez) ☕.
- Assurez-vous d'être dans le dossier ```./backend```
```sh   
mkdir keywords
python manage.py initBooks  
python manage.py computeKeywords
python manage.py addKeywords
python manage.py createGraphJaccard
python manage.py tfidf
python manage.py final_threshold
python manage.py graphVisualisation
```
### 2.2. Workflow of `./backend/data`
Le répertoire `./backend/data` contient toute la logique de traitement des données, calcul des mots-clés et création du graphe de similarité.
#### 2.2.1. `Commandes`
 - ***Avertissement*** : il est crucial d'exécuter ces commandes dans l'ordre, l'une après l'autre.
##### 2.2.1.1. `initBooks`
* Récupère les données des livres depuis l'API Gutendex `https://gutendex.com/books/` (jusqu'à `5000` threads).
* N'ajoute que les livres français et anglais.
* Extrait les métadonnées (titre, auteur, langue, sujets) et le lien vers le texte (ex: `https://www.gutenberg.org/cache/epub/26184/pg26184.txt`).
* Stocke les métadonnées dans la base de données.
* Stocke le texte du livre dans ```./backend/books``` avec le format `gutenberg_book_id.txt` (ex: `1.txt`). Cela simplifie les calculs (mots-clés, voisins via Jaccard) et garde la base sur les infos pertinentes, pas le texte brut.



##### 2.2.1.2. `computeKeywords`

-   Traite les mots‑clés depuis le répertoire `./backend/books`.
-   Normalise le texte (minuscule, suppression des stopwords, racinisation) avec le package `spacy`.
-   Calcule le nombre d’occurrences de chaque mot‑clé via `Counter`.
-   Produit la liste finale des mots‑clés, chacun associé à son occurrence, et enregistre un fichier JSON dans `./backend/keywords`, par exemple :
   ```
    {
       keyword_1: 12,
       keyword_2: 1 , 
       ...
    }
  ```

En résumé :
-   Garantit qu’il existe un fichier JSON par livre :
  -   Chargement des modèles NLP spécifiques à la langue (anglais et français).
  -   Extraction du texte depuis les fichiers du dossier `books`.
  -   Traitement du texte de chaque livre avec le modèle adapté.
  -   Lemmatisation et filtrage pour extraire des mots‑clés pertinents.
  -   Comptage des occurrences et stockage au format JSON.
  -   Sauvegarde du JSON dans `keywords` en utilisant l’ID du livre comme nom de fichier.

##### 2.2.1.3. `addKeywords`

- Lit les fichiers JSON de mots‑clés pré‑extraits dans le répertoire `keywords`.
- Chaque nom de fichier correspond à la clé primaire d’un livre (ex. : `123.json`).
- Chaque fichier contient un dictionnaire JSON mappant les mots‑clés à leurs occurrences.



1.  **Initialisation :**
  -   Initialise des dictionnaires pour suivre les mots‑clés et leur langue (anglais/français).
  -   Chaque langue possède son propre ensemble de mots‑clés (ex. `data_keywordsenglish`) afin que chaque mot‑clé soit unique.

2.  **Itération sur les fichiers :**
  -   Parcourt chaque fichier de mots‑clés dans `keywords`.
  -   Pour chaque fichier :
    -   **Extraction de l’ID du livre :** extrait l’ID depuis le nom (ex. : `123.json`).
    -   **Chargement JSON :** charge les données de mots‑clés et occurrences.
    -   **Récupération du livre :** récupère l’objet livre correspondant en base.
    -   **Catégorisation par langue :** classe les mots‑clés par langue (anglais ou français) ; chaque langue a sa table de correspondance (ex. : `keywordbook${language}`) pour mapper chaque mot‑clé d’un livre avec son occurrence.
    -   **Mapping des mots‑clés :** associe chaque mot‑clé à des paires (livre, occurrence) → table `(id, occurrence, book_id, keyword_id)`.

##### Opérations base de données


  -   Pour tous les mots‑clés distincts de ${Language} :
    -   Crée des objets `Keywords${Language}`.
    -   Crée des objets de relation `KeywordBook${Language}` reliant :
      -   Livres
      -   Occurrences
      -   Mots‑clés
        
    -  Chaque mot‑clé possède désormais un ID unique.
Au niveau de l’index :
 - Une table d’index est créée mappant l’ID du mot‑clé à l’ID du livre, avec son occurrence. Une table d’index existe pour chaque langue, ex. `keywordbook_${language}`.

##### Problème

- Après plusieurs essais, nous avons constaté que le pré‑traitement standard produisait ~10 millions de tokens en anglais et ~400 k en français. Après 8 heures de calcul continu, seulement ~250 000 tokens anglais sur 10 M avaient été traités. De plus, la table d’index générée contenait ~1 million de lignes, ce qui nous a poussés à mettre en place un mécanisme de seuillage par livre.


##### 2.2.1.4. `createGraphJaccard`

-   Calcule la similarité de Jaccard entre livres à partir de leurs mots‑clés.
-   Crée un graphe où les livres dont la similarité dépasse un seuil (0.5) sont reliés.
-   Stocke les relations de voisinage dans le modèle `Neighbors` :
  -   Charge tous les fichiers JSON de mots‑clés du dossier `keywords`.
  -   Construit un dictionnaire ID‑livre → occurrences de mots‑clés.
  -   Compare chaque livre à tous les autres via la distance de Jaccard.
  -   Si la distance est sous le seuil (donc similaires), connecte les livres comme voisins.
  -   Crée des relations bidirectionnelles en base de données.

<table>
  <tr>
    <td><img src="backend/graph/graph_3d_visualization.png" alt="French Keywords by Threshold" width="100%"></td>
    <td><img src="backend/graph/graph_3d_visualization_simplified.png" alt="French Keywords Reduction" width="100%"></td>
  </tr>

</table>
Visuellement : à gauche tous les nœuds, à droite les 1000 premières arêtes. Résultats :

- `60 composantes connexes`  
- Taille de la plus grande composante connexe : `942`


##### 2.2.1.5. `tfidf`
- Calcule le TF‑IDF pour chaque mot‑clé.

##### 2.2.1.7 `Nombre de tokens, taille des tables d’index`
 - Anglais 
   - nombre de tokens uniques : 36k 
   - taille de la table d’index : 580k lignes
 - Français 
   - nombre de tokens uniques : 6.5k 
   - taille de la table d’index : 40k lignes
 - Graphe :
   - nombre de sommets : 1099 (lignes)
   - nombre d’arêtes bidirectionnelles : 42 956 (lignes)
   - nombre d’arêtes : 21 478
## Calcul du nombre moyen de voisins par sommet

### Données :
- **Number of vertices**: \( n = 1099 \)
- **Total number of edges**: \( E = 21478 \)

### Calcul du degré moyen :
Comme $\sum_{u \in V}deg(v) = 2 \times |E|$ alors :

Formule du degré moyen $d_{\text{avg}}$ : $d_{\text{avg}} = \frac{2E}{n}$.

En remplaçant : $d_{\text{avg}} = \frac{2 \times 21478}{1099} \approx 39.08$.

Donc, en moyenne, **chaque sommet est connecté à ~39 autres sommets**.


## Nombre maximal théorique d'arêtes

Le **nombre maximal d'arêtes** dans un graphe simple (sans boucles ni multi‑arêtes) est :

$
E_{\text{max}} = \frac{n(n-1)}{2}
$

Substituting \( n = 1099 \):

$
E_{\text{max}} = \frac{1099 \times 1098}{2}
$

$
E_{\text{max}} = \frac{1206702}{2} = 603351
$
Donc, **le nombre maximal d'arêtes possible est 603 351**.
 

#### 2.2.2. Similarité de Jaccard

-   Le coefficient de similarité de Jaccard est utilisé pour mesurer la similarité entre deux ensembles de mots‑clés :
    -   `J(A, B) = |A ∩ B| / |A ∪ B|`
    -   Où :
      -   `A` et `B` sont les ensembles de mots‑clés de deux livres.
      -   `|A ∩ B|` est le nombre de mots‑clés communs.
      -   `|A ∪ B|` est le nombre total de mots‑clés uniques dans les deux livres.
-   Un score de Jaccard > 0.4 signifie que deux livres sont assez similaires pour être voisins.
-   Implémentation :
    -   calcule la différence d'occurrence des mots‑clés,
    -   normalise par la différence maximale,
    -   retourne une distance (plus petit = plus similaire),
    -   seuil de distance utilisé : 0.6.

#### 2.2.3. Sérialisation (`serializers.py`)

-   Convertit les modèles Django en JSON pour les réponses API.
-   Sérialise les livres, leurs métadonnées et les voisins :
  -   `LanguageSerializer` : expose le code langue.
  -   `PersonSerializer` : inclut le nom de l’auteur, ses années de naissance et de décès.
  -   `SubjectSerializer` : fournit les sujets.
  -   `BookSerializer` : données complètes du livre, incluant auteurs, langues et sujets.

#### 2.2.4. Configuration (`config.py`)

-   Contient des réglages globaux (seuils Jaccard, etc.).
-   Définit les chemins de stockage des mots‑clés.

#### 2.2.5. Vues (`views.py`)

-   Gère les endpoints :
    -   `server/books/` : renvoie les livres avec filtres (langue, auteur, titre, mot‑clé + langue, tri par téléchargements).
    -   `data/books/neighbors/<int:pk>` : renvoie les voisins d'un livre via centralité (betweenness/closeness), avec détails.
        
#### 2.2.6. `sort.py` : logique de tri et de suggestion

-   **But :** implémente le tri basé sur la centralité et génère des suggestions.
-   **Fonctionnement :**
  1.  **Suggestions :**
    -   Vérifie le cache pour des résultats existants.
    -   Ne traite que les 2 premiers livres de la liste.
    -   Récupère les voisins via l’API `/data/books/neighbors/`.
    -   Utilise un cache par livre pour optimiser les requêtes réseau.
    -   Ajoute un timeout pour éviter les blocages.
    -   Retourne une liste de livres suggérés uniques (jusqu’à 10).
  2.  **Tri par centralité :**
    -   Vérifie le cache pour des calculs existants.
    -   Construit un graphe (pondéré ou non) selon la similarité des sujets.
    -   Calcule les mesures de centralité avec des algorithmes standards ou simplifiés.
    -   Trie les livres selon ces scores.
    -   Met en cache 24 h.
    -   Retourne la liste triée des livres.
   
-   **Représentation du graphe :**
  -   Les livres sont des nœuds du graphe.
  -   Les arêtes représentent la similarité entre livres.
  -   **Poids d’arête :** en graphe pondéré, le poids d’une arête représente le degré de similarité, calculé par le nombre de sujets partagés.
  -   **Ajout de voisins :** ajout selon l’intersection des listes de sujets ; si non vide, création d’une arête entre les nœuds.
#### 2.2.7. `graph.py` : structures et algorithmes de graphe

-   **But :** structures de graphe et algorithmes de centralité.
-   **Représentation :**
  -   Utilise des classes `Node` pour représenter les livres, stocker leurs données et scores de centralité.
  -   Les classes `UnweightedGraph` et `WeightedGraph` fournissent les implémentations de graphe.
  -   Les arêtes sont stockées comme relations de voisinage dans les nœuds.
-   **Algorithmes :**
  -   Algorithme de Brandes pour la betweenness (graphe non pondéré).
-   **Ajout de voisins :**
  -   `UnweightedGraph` : ajoute des nœuds voisins à la liste des voisins.
  -   `WeightedGraph` : ajoute des voisins avec leurs poids d’arête dans le dictionnaire des voisins.
#### 2.2.8. `centrality.py` : calcul des centralités

-   **But :** implémente les algorithmes de closeness et betweenness.
-   **Closeness :** approxime sur les 20 plus proches, mise à l'échelle, inverse de la somme des distances.
-   **Betweenness :** sur grands graphes : top 30 nœuds, par nœud top 20 voisins, timeout 2s, Brandes optimisé.
-   **Fonctions :** calcule les deux centralités en s'appuyant sur `graph.py`.
#### `Stratégie de cache`
- Mise en place d’un cache multi‑niveaux selon les parties de l’application.
- Mise en cache des réponses API, des données de voisins et des calculs de centralité.
- Clés de cache uniques basées sur les paramètres de requête.
- Timeouts raisonnables pour les éléments mis en cache.




## 2.3. Résultats :
- Comparaison betweenness vs closeness. Méthode : à la requête, on prend les voisins déjà calculés en BD, construit le graphe (pondéré pour closeness, non‑pondéré pour betweenness), puis on trie (incluant le nombre de téléchargements). Tests sur `sargon` pour closeness et sur le titre « The c » pour betweenness.
     -  ***Temps de calcul***
       - 2. ***Closeness*** : 0.005 s
       - 3. ***Betweenness*** : 0.0067 s
       - 4. ***Nombre de téléchargements*** : 0.01 s 

## 2.4. Comparaison sur les requêtes API
But : obtenir les livres où le mot `hello` apparaît et mesurer le temps de réponse de chaque méthode (du début de la requête à l’envoi de la réponse).
    -  ***Résultats des requêtes API*** :
      - 2. ***Closeness*** : 1.32 s
      - 3. ***Betweenness*** : 1.44 s
      - 4. ***Nombre de téléchargements*** : 2 s 

## 3. Démarrage du serveur
Dans le dossier ```./backend``` :
```bash
python manage.py runserver
```
L’hôte par défaut est ```localhost:8000```

## 4. Démarrage du frontend

``` bash
cd react-frontend
npm i 
npm run dev
```
Hôte par défaut : ```localhost:3000```

---

## 5. Démarrage avec Docker Compose (Backend + Frontend)

Pré-requis : Docker et Docker Compose v2 installés.

Depuis la racine du projet :

```bash
docker compose build
docker compose up -d
```

Services démarrés :
- Frontend (Nginx) exposé sur le port `8080` de l’hôte → `http://localhost:8080`
- Backend (Django + Gunicorn) exposé sur le port `8000` → `http://localhost:8000`

Le frontend proxy les appels API via le chemin `/server/` vers le backend. Exemple :

```
http://localhost:8080/server/books/?keyword=hello
```

Configuration :
- `ALLOWED_HOSTS` dans `backend/backend/settings.py` doit inclure votre IP/domaine et `backend` (nom du service Docker) pour que Nginx → Django fonctionne.
- Le frontend utilise une base d’URL API relative par défaut (`/server/books/?`). Vous pouvez la surcharger avec `VITE_API_BASE` au besoin.

Arrêt / redémarrage :
```bash
docker compose down
docker compose up -d
```

Déploiement VPS : voir `DEPLOYMENT.md` pour une procédure complète (build, reverse proxy, autodeploy via systemd, etc.).

