# MNLRMap

> [!WARNING]  
> This Application has currently no Authentication

MNLRMap is a self-hosted map extraction and viewing application that allows you to create, manage, and view offline vector maps based on OpenStreetMap data using the PMTiles format.

## Features

- **Map Extraction**: Select areas from OpenStreetMap to create custom vector map extracts
- **Map Management**: View, manage, delete, and recreate your map extracts through a user-friendly interface
- **Offline Vector Maps**: View your maps offline using MapLibre GL with PMTiles support
- **Responsive UI**: Modern interface built with Vue 3, TypeScript, and Quasar Framework
- **Real-time Status Updates**: Track the progress of map extractions in real-time
- **Lightweight Backend**: Powered by PocketBase and Go

## How It Works

MNLRMap uses Protomaps technology to extract and serve vector maps. The application:

1. Extracts map data from OpenStreetMap based on user-defined bounding boxes
2. Processes the data into the PMTiles format for efficient storage and delivery
3. Stores maps in a local database for management
4. Serves maps through a MapLibre GL-powered frontend for smooth vector rendering

## Use Cases

- Create offline maps for hiking, travel, or fieldwork
- Build custom maps for specific regions or points of interest
- Deploy in environments with limited internet connectivity
- Develop custom mapping applications with self-hosted map tiles

## Deployment

### Docker

You can deploy MNLRMap using the official Docker image from GitHub Container Registry:

```bash
docker run -d \
  --name mnlrmap \
  -p 8090:8090 \
  -v ./data:/db \
  ghcr.io/manuelraven/mnlrmap:latest
```

### Docker Compose

For easier deployment, you can use this docker-compose.yml:

```yaml
version: "3"
services:
  mnlrmap:
    image: ghcr.io/manuelraven/mnlrmap:latest
    container_name: mnlrmap
    ports:
      - "8090:8090"
    volumes:
      - ./data:/db
    restart: unless-stopped
```

Start the service with:

```bash
docker-compose up -d
```

## Accessing the Application

After deployment, access MNLRMap through your web browser at:

```
http://localhost:8090
```

## License

This project is licensed under The Unlicense - see the LICENSE file for details.

Packages can have different License.

## Acknowledgements

- [Protomaps](https://protomaps.com/) for the PMTiles technology
- [MapLibre GL](https://maplibre.org/) for the mapping library
- [PocketBase](https://pocketbase.io/) for the backend database
- [OpenStreetMap](https://www.openstreetmap.org/) for the map data
