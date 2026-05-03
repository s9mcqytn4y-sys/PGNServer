FROM php:8.4-cli-bookworm

WORKDIR /var/www/html

ENV DEBIAN_FRONTEND=noninteractive
ENV COMPOSER_ALLOW_SUPERUSER=1
ENV TZ=Asia/Jakarta

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        git \
        unzip \
        curl \
        libicu-dev \
        libpq-dev \
        libzip-dev \
        libonig-dev \
        libxml2-dev \
        libsqlite3-dev \
        postgresql-client \
    && docker-php-ext-install \
        bcmath \
        intl \
        pdo_pgsql \
        pdo_sqlite \
        zip \
    && rm -rf /var/lib/apt/lists/*

COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

EXPOSE 8000

CMD ["sh", "-lc", "php artisan serve --host=0.0.0.0 --port=8000"]
