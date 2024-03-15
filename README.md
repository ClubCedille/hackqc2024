<p align="center">
  <img src="template/static/logo_munis.png" alt="MUNIS Logo" width="200"/>
</p>

<h1 align="center">MUNIS - HackQC 2024</h1>
<p align="center">Votre bouée communautaire au Québec.</p>

<p align="center">
  <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License Badge"/>
  <img src="https://github.com/ClubCedille/hackqc2024/actions/workflows/main.yml/badge.svg" alt="Workflow Badge"/>
</p>

## Description

MUNIS permet d'extraire les évènements d’urgence de multiples sources de données, permettant une réaction rapide et informée. En parallèle, la plateforme invite les utilisateurs à contribuer, permettant l'ajout de ressources d’aide qui viennent enrichir notre communauté de soutien. L'application comprend différentes sections dédiées à l'aide offerte et aux moyens de contribuer aux initiatives, facilitant ainsi l'engagement des utilisateurs.

## Dépendances

- [Go 1.22](https://go.dev/doc/install)
- [Air](https://github.com/cosmtrek/air)

## Configuration

### Nix

Vous pouvez utiliser nix pour configurer les dépendances : `nix develop`

### Utilisation de Make (alternative)

Alternativement, vous pouvez initialiser manuellement le projet en utilisant le Makefile : `make init`

## Comment exécuter

Dans la racine du projet, exécutez l'une des commandes suivantes :

Exécution Go :

`go run .`

Rechargement en direct :

`air`

### Utilisation du Makefile

```bash
make build
make start
```

#### Avec Docker

Construisez l'image Docker et exécutez le conteneur sur le port 8080

```bash
make docker-build
make docker-run
```
