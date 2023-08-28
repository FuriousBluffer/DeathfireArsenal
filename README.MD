# Deathfire Arsenal

Deathfire Arsenal is a web service that allows players to create rooms and join games in various modes. It provides an API for creating players, rooms, joining and leaving rooms, and retrieving mode trends by region and player.

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Documentation](#api-documentation)

## Features

- Create and manage players and rooms.
- Join and leave game rooms with different modes.
- Get mode trends by region and player.
- MongoDB for data storage.
- Redis for caching.

## Prerequisites

Before running the Deathfire Arsenal web service, make sure you have Docker installed on your system.

## Installation

To install and run Deathfire Arsenal, follow these steps:

1. Clone the repository:
2. Run the web service using Docker: You can make use of the Makefile to accomplish this.

   ```bash
   make build
   make run
   ```

   `make build` is equivalent to `docker-compose build`.

   `make run` is equivalent to `docker-compose up`.

## Usage

The Deathfire Arsenal web service will be available at `http://localhost:8080`.

## API Documentation

The API documentation for Deathfire Arsenal is available at [OPEN API Specs](documentation/documentation.yaml). It provides information about the available API endpoints, their input parameters, and expected responses. You can copy and paste the YAML file content into an [online Swagger UI editor](https://editor-next.swagger.io/) to visualize the API documentation in a user-friendly interface.

## A little about the product
- The Game consists of two entities - Players and Gaming Room (An ongoing match where players can join and leave)
- A player shall be identified with an unique ID and will belong to a region.
- A player can create a room for a particular mode of game(like “Team Deathmatch” or “1 V 1”) and share the room id, using which more players can join the room.
- The modes of the game are limited to 5 for the sake of simplicity of testing. Team Deathmatch, Gunsmith, Mayhem, Battle Royale and 1 V 1.
- A single room will have an upper limit based on the mode of the game.
- A player at any given point of time can be playing in a single game or not playing at all, i.e. cannot be playing more than 1 game at a time.
- A room can consist of players from different regions.
- Get Mode Trends By Region shall return the top 3 modes being currently played by the players belonging to the region as provided.