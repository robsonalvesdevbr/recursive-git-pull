# Recursive Git Pull (RGP)

Uma ferramenta de linha de comando em Go para executar comandos Git recursivamente em múltiplos repositórios.

## Características

- Execução de comandos Git em todos os repositórios dentro de um diretório
- Execução paralela para melhor performance  
- Suporte a comandos Git personalizados
- Filtros por padrões de nome (include/exclude)
- Opção para ignorar repositórios "sujos" (com mudanças não commitadas)
- Timeout configurável para comandos
- Logs detalhados e relatórios de status
- Funcionalidade especial para pull em todos os branches
- Interface de linha de comando simples e intuitiva

## Instalação

### Compilar do código fonte

```bash
git clone https://github.com/robsonalvesdevbr/recursive-git-pull.git
cd recursive-git-pull
go build -o bin/rgp ./cmd/rgp

# Opcional: Instalar globalmente
go install ./cmd/rgp
```

### Usando Go Install (futuro)

```bash
go install github.com/robsonalvesdevbr/recursive-git-pull/cmd/rgp@latest
```

## Uso

### Sintaxe básica

```bash
rgp [opções]
```

### Opções disponíveis

- `-path string`: Diretório raiz para buscar repositórios Git (padrão: ".")
- `-command string`: Comando Git para executar (padrão: "pull")  
- `-parallel`: Executar comandos em paralelo (padrão: true)
- `-workers int`: Número máximo de workers paralelos (padrão: 4)
- `-timeout string`: Timeout para cada comando (padrão: "30s")
- `-ignore-dirty`: Ignorar repositórios com mudanças não commitadas
- `-include string`: Padrões para incluir repositórios (separados por vírgula)
- `-exclude string`: Padrões para excluir repositórios (separados por vírgula)
- `-all-branches`: Pull todos os branches (funciona apenas com comando pull)
- `-verbose`: Saída detalhada
- `-help, -h`: Mostrar ajuda

### Exemplos

#### 1. Pull básico em todos os repositórios

```bash
rgp -path ./workspace -command pull
```

#### 2. Status de todos os repositórios (sem paralelização)

```bash
rgp -path ./projects -command status -parallel=false
```

#### 3. Pull em todos os branches de todos os repositórios

```bash
rgp -path ./repos -command pull -all-branches
```

#### 4. Usando filtros para incluir apenas serviços

```bash
rgp -include '*-service' -exclude 'test-*'
```

#### 5. Comando personalizado com timeout

```bash
rgp -command "log --oneline -5" -timeout 10s -verbose
```

#### 6. Ignorar repositórios sujos durante pull

```bash
rgp -command pull -ignore-dirty
```

## Estrutura do projeto

```
workspace/
├── repo1/
│   └── .git/
├── repo2/  
│   └── .git/
└── repo3/
    └── .git/
```

A ferramenta encontrará automaticamente todos os repositórios Git (diretórios contendo `.git/`) dentro do diretório especificado.

## Casos de uso comuns

### Desenvolvimento com microserviços
```bash
# Pull de todos os microserviços
rgp -path ./microservices -command pull -verbose

# Status de todos os serviços
rgp -path ./microservices -command status
```

### Manutenção de múltiplos projetos
```bash
# Verificar branches atuais
rgp -command "branch --show-current"

# Fetch de todos os repositórios
rgp -command "fetch --all"
```

### Limpeza e organização
```bash
# Verificar repositórios sujos
rgp -command "status --porcelain"

# Limpar branches merged
rgp -command "branch -d --merged"
```

## Desenvolvimento

### Requisitos

- Go 1.19 ou superior
- Git instalado no sistema

### Compilar e testar

```bash
# Compilar
go build -o bin/rgp ./cmd/rgp

# Testar
./bin/rgp -help

# Executar testes (quando disponíveis)
go test ./...

# Formatar código
go fmt ./...
```

### Estrutura do código

```
├── cmd/rgp/           # Ponto de entrada da aplicação
├── internal/          # Lógica interna da aplicação
│   ├── config/        # Configuração e parsing de flags
│   ├── finder/        # Descoberta de repositórios
│   └── git/           # Execução de comandos Git
├── pkg/types/         # Tipos públicos
└── Makefile           # Scripts de build
```

## Contribuindo

1. Faça fork do projeto
2. Crie sua feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)  
5. Abra um Pull Request

## Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## Problemas conhecidos

- A funcionalidade `-all-branches` pode ser lenta em repositórios com muitos branches
- Alguns comandos Git interativos não são suportados

## Roadmap

- [ ] Suporte a arquivos de configuração (YAML/JSON)
- [ ] Integração com hooks Git
- [ ] Suporte a templates de comandos
- [ ] Interface web opcional
- [ ] Testes automatizados
- [ ] Binários pré-compilados para releases

## Suporte

Para problemas ou sugestões, abra uma issue no [GitHub](https://github.com/robsonalvesdevbr/recursive-git-pull/issues).