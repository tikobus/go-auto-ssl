# Get SSL Certs Auto

`

    # HTTP server
    #
    server {
        listen       80;
        server_name  your.domain;

        location /.well-known/acme-challenge/ {
            proxy_set_header Host "your.domain";
            proxy_pass http://127.0.0.1:20080;
        }

        location / {
            root   html;
            index  index.html index.htm;
        }
    }

    # HTTPS server
    #
    server {
        listen       443 ssl;
        server_name  your.domain;

        ssl_certificate      ssl/your.domain;
        ssl_certificate_key  ssl/your.domain;

        location / {
            root   html;
            index  index.html index.htm;
        }
    }

`