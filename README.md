# Waktu Solat Teluk Intan

This is a simple web application that displays current prayer times for a specific location, built with a Go backend and a modern, interactive HTML frontend.

## Features

* **Prayer Time Display**: Shows the five daily prayer times (Subuh, Zohor, Asar, Maghrib, Isyak).
* **Dynamic Data**: Fetches prayer times from a Go API running on the backend.
* **Real-time Clock**: Updates the current time every second.
* **Interactive Highlighting**: Automatically highlights the current or upcoming prayer time.
* **Responsive Design**: Looks great on both desktop and mobile devices.
* **Dark Mode Support**: Adapts to the user's system preference for a comfortable viewing experience.

## Technology Stack

### Frontend:

* HTML5
* JavaScript (Vanilla JS for all logic)
* Tailwind CSS (via CDN for styling)

### Backend:

* Go (Golang) for the API server

## How to Run

### 1. Backend (Go API)

The Go backend must be running to provide the prayer time data. The API is expected to be accessible at `http://127.0.0.1:502/api/prayer-times`.

Build the executable:

```bash
go build -o prayer-times-api main.go
```

Run as a systemd service (recommended):
A systemd service file is provided to manage the application.

Create the `prayer-times-api.service` file to `/etc/systemd/system/`.

Reload the daemon, enable, and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable prayer-times-api.service
sudo systemctl start prayer-times-api.service
```

### 2. Frontend (HTML)

The frontend is a single `index.html` file that communicates with the backend API.

#### Nginx Configuration:

The `index.html` file and other assets (if any) should be served by a web server like Nginx, which also acts as a reverse proxy for the API. An example Nginx configuration is provided to route requests to your Go backend.

```nginx
server {
    listen 80;
    server_name solat.home solat.syazwansaidan.my;
    root /var/www/html/solat;
    index index.html;

    location /api/ {
        proxy_pass http://127.0.0.1:502;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Customization

* **Location**: The prayer times are hardcoded for "Teluk Intan". You can modify the Go backend to accept a location parameter to serve other areas.
* **Styling**: The frontend uses Tailwind CSS via CDN. You can easily adjust the theme and colors by changing the classes or the theme settings.

## Credits

* **Prayer Times Data**: Provided by a Go backend service.
* **Icons**: Unicode character 🕌 for the favicon.
