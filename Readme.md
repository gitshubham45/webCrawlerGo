# Web Crawler with MongoDB and Redis

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

This project implements a distributed web crawler that uses **Redis** for managing a Bloom filter to track visited URLs and **MongoDB** for storing crawled results. The crawler is designed to efficiently crawl websites, extract product URLs, and store the results in a structured format.

---

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Prerequisites](#prerequisites)
4. [Installation](#installation)

---

## Overview

The web crawler is designed to:
- Use **Redis** as a distributed Bloom filter to track visited URLs and avoid redundant crawling.
- Store crawled results (e.g., product URLs) in **MongoDB** for long-term access.
- Optionally save results to JSON files for offline analysis.

Key components:
- **Redis**: Manages the Bloom filter and ensures efficient tracking of visited URLs.
- **MongoDB**: Stores the crawled results in a database.
- **Go**: Implements the crawler logic and integrates with Redis and MongoDB.

---

## Features

- **Distributed Crawling**: Uses Redis to share the Bloom filter state across multiple workers.
- **Efficient Tracking**: Implements a Bloom filter to minimize memory usage while tracking visited URLs.
- **Persistent Storage**: Crawled data is stored in MongoDB for long-term access.
- **File Backup**: Optionally saves results to JSON files for offline analysis.
- **Scalable**: Designed for distributed systems; can scale horizontally by adding more crawler workers.

---

## Prerequisites

Before running the project, ensure you have the following installed:

- **Go**: Install Go.

---

## Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/yourusername/web-crawler.git
   
   cd web-crawler

## Usage
    ```bash
   docker run -d -p 6379:6379 --name my-redis redis:latest

   docker run -d --name my_mongo -p 27017:27017 mongo

   go run main.go