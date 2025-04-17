# Sistema de Gestão Residencial Portarius

Sistema de gestão para portaria de condomínio residencial desenvolvido com Electron, Go e Postgres.

## Funcionalidades

- Cadastro de moradores
- Reserva de salão de festas
- Inventário de veículos e pets
- Gestão de encomendas
- Integração com WhatsApp 
- Autenticação de usuários

## Requisitos

- Node.js (v18 ou superior)
- Go (v1.21 ou superior)
- Postgres (v15 ou superior)

## Estrutura do Projeto

```
.
├── frontend/           # Aplicação Electron
├── backend/           # Servidor Go
└── docs/             # Documentação
```

## Configuração do Ambiente

1. Instale as dependências do frontend:
```bash
cd frontend
npm install
```

2. Instale as dependências do backend:
```bash
cd backend
go mod init portarius
go mod tidy
```

3. Configure as variáveis de ambiente:
```bash
cp .env.example .env
```

## Executando o Projeto

1. Inicie o servidor backend:
```bash
cd backend
go run main.go
```

2. Inicie a aplicação Electron:
```bash
cd frontend
npm start
```

## Licença

MIT 