# Waktu Solat Teluk Intan

This project provides a simple, elegant web application to display daily prayer times for Teluk Intan, Perak, Malaysia. It highlights the current prayer time and updates the current time in real-time.

---

## ✨ Features

* **Real-time Clock:** Displays the current time, updated every second.
* **Dynamic Prayer Times:** Fetches and displays daily prayer times for Teluk Intan.
* **Current Prayer Highlight:** Automatically highlights the current prayer time based on the actual time.
* **Responsive Design:** Optimized for various screen sizes, including mobile devices.
* **Light/Dark Mode:** Adapts to the user's system preference for light or dark mode.
* **Refresh Button:** Manually refresh prayer time data from the backend.

---

## 💻 Technologies Used

* **Frontend:**
    * **HTML5:** Structure of the web page (single file).
    * **CSS3:** Styling (inline within HTML).
    * **JavaScript (Vanilla JS):** Client-side logic for fetching data, updating time, and handling UI interactions (inline within HTML).
* **Backend:**
    * **PHP:** Server-side scripting for fetching prayer times from an external API (`e-solat.gov.my`).
    * **cURL (PHP Extension):** Used for making HTTP requests to external APIs.
* **Web Server:**
    * **Nginx:** Serves the HTML frontend and acts as a reverse proxy for the PHP backend using PHP-FPM.
* **PHP FastCGI Process Manager (PHP-FPM):** Executes PHP scripts.

---

## 🚀 Setup Guide

Follow these steps to set up and run the application on your Debian-based server (e.g., Orange Pi Zero 3).

### Prerequisites

Before you begin, ensure you have the following installed on your server:

* **Nginx:** A high-performance web server.
* **PHP 8.2 and PHP-FPM:** PHP runtime and its FastCGI Process Manager.
* **PHP cURL extension:** Required by the PHP backend to make external API calls.

You can install them using:

```bash
sudo apt update
sudo apt install nginx php8.2-fpm php8.2-curl -y
```

### 1. Backend Setup (`api.php`)

The PHP script fetches prayer times.

1.  **Create the web directory:**
    ```bash
    sudo mkdir -p /var/www/html/solat
    ```
2.  **Create the `api.php` file:**
    ```bash
    sudo nano /var/www/html/solat/api.php
    ```
3.  **Paste the PHP code (from our previous conversation) into `api.php`.**
4.  **Set permissions for `api.php`:**
    ```bash
    sudo chmod 644 /var/www/html/solat/api.php
    sudo chown www-data:www-data /var/www/html/solat/api.php
    ```
5.  **Restart PHP-FPM:**
    ```bash
    sudo systemctl restart php8.2-fpm
    ```

### 2. Frontend Setup (`index.html`)

The frontend is a single HTML file containing all HTML, CSS, and JavaScript.

1.  **Create the `index.html` file:**
    ```bash
    sudo nano /var/www/html/solat/index.html
    ```
2.  **Paste the HTML code (from our previous conversation) into `index.html`.**
3.  **Set permissions for `index.html`:**
    ```bash
    sudo chmod 644 /var/www/html/solat/index.html
    sudo chown www-data:www-data /var/www/html/solat/index.html
    ```

### 3. Nginx Configuration

Configure Nginx to serve your HTML file and pass PHP requests to PHP-FPM.

1.  **Edit your Nginx site configuration file:**
    ```bash
    sudo nano /etc/nginx/sites-available/default # Or your custom site file, e.g., /etc/nginx/sites-available/solat.home
    ```
2.  **Paste the following Nginx configuration:**

    ```nginx
    server {
        listen 80;
        server_name solat.home; # Replace with your domain or server IP

        location / {
            root /var/www/html/solat;
            index index.html index.htm;
            try_files $uri $uri/ =404;
        }

        location ~ \.php$ {
            root /var/www/html/solat;
            try_files $uri =404;

            fastcgi_split_path_info ^(.+\.php)(/.+)$;
            
            fastcgi_pass unix:/var/run/php/php8.2-fpm.sock; # Ensure this matches your PHP-FPM socket
            fastcgi_index index.php;
            include fastcgi_params;
            fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
            fastcgi_read_timeout 120;
        }
    }
    ```
    * **Important:** Replace `solat.home` with your actual domain name or the IP address of your Orange Pi if you're accessing it directly by IP.
    * Ensure `fastcgi_pass unix:/var/run/php/php8.2-fpm.sock;` correctly points to your PHP 8.2 FPM socket.

3.  **Test and reload Nginx:**
    ```bash
    sudo nginx -t
    sudo systemctl reload nginx
    ```

### 4. Client-side Hosts File (if using `solat.home`)

If you're accessing the application using `http://solat.home`, you need to map this hostname to your Orange Pi's IP address on your client machine (the computer you're browsing from).

* **On Linux/macOS:** Edit `/etc/hosts`
* **On Windows:** Edit `C:\Windows\System32\drivers\etc\hosts` (as administrator)

Add the following line, replacing `YOUR_ORANGE_PI_IP` with your Orange Pi's actual IP address:

```
YOUR_ORANGE_PI_IP solat.home
```

---

## 🚀 Usage

After completing the setup:

1.  Open your web browser.
2.  Navigate to `http://solat.home` (or your Orange Pi's IP address if you used it in the Nginx config).

The page should load, display the current date and time, and fetch the prayer times from the backend. The current prayer time will be highlighted.

---

## 🔍 Troubleshooting

* **404 Not Found for `api.php`:**
    * Verify `api.php` exists at `/var/www/html/solat/api.php`.
    * Check file and directory permissions (`sudo chmod 644 /var/www/html/solat/api.php` and `sudo chown www-data:www-data /var/www/html/solat/api.php`).
    * Ensure PHP-FPM is running (`sudo systemctl status php8.2-fpm`).
    * Check Nginx error logs (`sudo tail -f /var/log/nginx/error.log`) for specific errors.
* **Empty Prayer Times / API returns empty `prayers` array:**
    * The external `e-solat.gov.my` API might be down or returning unexpected data. Try accessing `https://www.e-solat.gov.my/index.php?r=esolatApi/takwimsolat&zone=PRK05&period=today` directly in your browser or with `curl` from your server to verify its status.
    * Check PHP error logs (if configured) for issues during API fetching.
* **"undefined" next to the time / incorrect time format:**
    * Ensure your PHP backend is returning the time in `HH:MM:SS` format and your JavaScript is parsing it correctly (which it should be with the current code).
    * Clear your browser cache.

---

## 📄 License

This project is open-source and available under the MIT License.
