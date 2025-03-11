# API-File 📂

![Go](https://img.shields.io/badge/Go-1.17-blue)
![Fiber](https://img.shields.io/badge/Fiber-2.0-green)
![Gorm](https://img.shields.io/badge/Gorm-1.21.12-orange)

## 📜 Description

API-File is a powerful API designed for managing files and documents. It provides endpoints for uploading, retrieving, updating, and deleting images and documents. Additionally, it supports real-time upload progress tracking via WebSocket.

## 🌐 WebSocket

Uploads can be recorded and tracked in real-time using the WebSocket routes provided in `websocket_routes.go`.

## 📋 Endpoints

### Private Routes

- **Storage Paths**
    - `GET /v1/storage-paths/` - Get all storage paths
    - `POST /v1/storage-paths/` - Create a new storage path
    - `GET /v1/storage-paths/:id` - Get a specific storage path
    - `PUT /v1/storage-paths/:id` - Update a specific storage path

- **Folders**
    - `POST /v1/folders/` - Create a new folder
    - `GET /v1/folders/:id` - Get a specific folder
    - `PUT /v1/folders/:id` - Update a specific folder
    - `DELETE /v1/folders/:id` - Delete a specific folder
    - `PUT /v1/folders/:id/restore` - Restore a deleted folder

- **Images**
    - `POST /v1/images/` - Upload a new image
    - `GET /v1/images/:id` - Get a specific image
    - `PUT /v1/images/:id` - Update a specific image
    - `DELETE /v1/images/:id` - Delete a specific image
    - `PUT /v1/images/:id/restore` - Restore a deleted image

- **Documents**
    - `POST /v1/documents/` - Upload a new document
    - `GET /v1/documents/:id` - Get a specific document
    - `PUT /v1/documents/:id` - Update a specific document
    - `DELETE /v1/documents/:id` - Delete a specific document
    - `PUT /v1/documents/:id/restore` - Restore a deleted document

- **WebSocket**
    - `GET /v1/handshake` - Handshake route for WebSocket

### Public Routes

- **Image**
    - `GET /v1/image/:id` - Get a specific image file
    - `GET /v1/image/:id/:size` - Get a specific image file with size

- **Document**
    - `GET /v1/document/:id` - Get a specific document file

- **WebSocket**
  -  `WS /v1/ws/progress` - WebSocket route for real-time upload progress tracking

## 🚀 Getting Started

The project can be easily started with Docker by using the `dev` or `prod` environment.

### Development

```sh
docker compose up -d dev
```

### Production

```sh
docker compose up -d prod
```

## 🤝 Contributing
We welcome contributions! Please fork the repository and submit a pull request.

## 📝 License

This project is licensed under the MIT License.

## 📞 Contact

For any questions or support, please contact [arnold.molenaar@webmi.nl](mailto:arnold.molenaar@webmi.nl).
<hr></hr> Made with ❤️ by Arnold Molenaar