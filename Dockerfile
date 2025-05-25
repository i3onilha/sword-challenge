FROM golang:1.24.3-bullseye AS development

LABEL maintainer="Jean Bonilha <jeanbonilha.webdev@gmail.com>"

ARG HOME_USER=/home/go

ENV DEBIAN_FRONTEND noninteractive
ENV GOENV="development"

ENV NODE_VERSION v22.12.0
ENV NVM_DIR ${HOME_USER}/.nvm
ENV NPM_FETCH_RETRIES 2
ENV NPM_FETCH_RETRY_FACTOR 10
ENV NPM_FETCH_RETRY_MINTIMEOUT 10000
ENV NPM_FETCH_RETRY_MAXTIMEOUT 60000

ENV SOURCE_CODE ${HOME_USER}/sourcecode

RUN go install golang.org/x/tools/cmd/godoc@v0.5.0
RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.18.0
RUN go install github.com/air-verse/air@v1.52.3
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

# Install Minikube
RUN curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && \
    install minikube-linux-amd64 /usr/local/bin/minikube && \
    rm minikube-linux-amd64

RUN set -xe; \
    apt-get update && \
    apt-get upgrade -yqq && \
    apt-get install -yqq \
    apt-utils \
    gnupg2 \
    git \
    libzip-dev zip unzip \
    default-mysql-client \
    inetutils-ping \
    wget \
    libaio-dev \
    freetds-dev \
    sudo \
    bash-completion \
    curl \
    make \
    ncurses-dev \
    build-essential \
    tree \
    nano \
    tmux \
    tmuxinator \
    xclip \
    apt-transport-https \
    ca-certificates \
    gnupg-agent \
    software-properties-common \
    libssl-dev \
    libgtk-3-dev \
    libwebkit2gtk-4.0-dev \
    nsis \
    ripgrep \
    fontconfig \
    gcc \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /usr/share/fonts/truetype/nerd-fonts \
    && wget -O /tmp/nerd-fonts.zip https://github.com/ryanoasis/nerd-fonts/releases/latest/download/FontPatcher.zip \
    && unzip /tmp/nerd-fonts.zip -d /usr/share/fonts/truetype/nerd-fonts \
    && rm /tmp/nerd-fonts.zip \
    && fc-cache -fv

RUN curl -LO https://github.com/neovim/neovim/releases/download/v0.10.2/nvim-linux64.tar.gz && \
    tar -C /opt -xzf nvim-linux64.tar.gz && \
    rm nvim-linux64.tar.gz

RUN useradd -ms /bin/bash go && echo "go:secret" | chpasswd && adduser go sudo

RUN rm -rf /etc/localtime && \
    ln -s /usr/share/zoneinfo/Europe/Lisbon /etc/localtime

USER go

RUN mkdir -p $NVM_DIR \
    && curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash \
    && . $NVM_DIR/nvm.sh \
    && nvm install ${NODE_VERSION} \
    && nvm use ${NODE_VERSION} \
    && nvm alias ${NODE_VERSION} \
    && npm config set fetch-retries ${NPM_FETCH_RETRIES} \
    && npm config set fetch-retry-factor ${NPM_FETCH_RETRY_FACTOR} \
    && npm config set fetch-retry-mintimeout ${NPM_FETCH_RETRY_MINTIMEOUT} \
    && npm config set fetch-retry-maxtimeout ${NPM_FETCH_RETRY_MAXTIMEOUT} \
    && npm install -g yarn \
    && npm install -g npm \
    && git clone --depth=1 https://github.com/i3onilha/nvim $HOME/.config/nvim \
    && /opt/nvim-linux64/bin/nvim -c 'MasonInstallAll' -c 'qa' \
    && /opt/nvim-linux64/bin/nvim -c 'GoUpdateBinaries' -c 'qa'

RUN git clone --depth 1 https://github.com/junegunn/fzf.git $HOME/.fzf && $HOME/.fzf/install

RUN git clone --bare -b godevenv https://github.com/i3onilha/.dotfiles.git $HOME/.dotfiles && \
    git clone https://github.com/i3onilha/.tmux.git $HOME/.tmux && \
    ln -sf .tmux/.tmux.conf $HOME && \
    cp $HOME/.tmux/.tmux.conf.local $HOME && \
    git --git-dir=$HOME/.dotfiles/ --work-tree=$HOME config --local status.showUntrackedFiles no && \
    git --git-dir=$HOME/.dotfiles/ --work-tree=$HOME reset HEAD . && \
    git --git-dir=$HOME/.dotfiles/ --work-tree=$HOME checkout -- .

RUN export PATH="$HOME/.nvm/versions/node/$NODE_VERSION/bin:$PATH"

WORKDIR $SOURCE_CODE

COPY . .

FROM golang:1.24.3-alpine AS builder

WORKDIR /home/go/sourcecode

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

COPY go.* .

RUN go mod download

COPY . .

# Build with optimizations and minimal binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o service ./cmd/server/server.go

FROM alpine:3.19 AS production

ENV SOURCE_CODE /app

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Europe/Lisbon /etc/localtime && \
    echo "Europe/Lisbon" > /etc/timezone && \
    apk del tzdata

WORKDIR /app

COPY --from=builder /home/go/sourcecode/service /app/service
COPY --from=builder /home/go/sourcecode/.env-prod /app/.env

CMD ["/app/service"]
