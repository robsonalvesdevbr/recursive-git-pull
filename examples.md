# Exemplos Práticos - RGP

## Scripts de Exemplo

### 1. Script para Pull em Workspace Completo

```bash
#!/bin/bash
# pull-all.sh - Pull em todos os repositórios do workspace

echo "Fazendo pull em todos os repositórios..."
rgp -path ~/workspace -command pull -verbose -ignore-dirty

echo "Pull concluído!"
```

### 2. Script para Verificar Status de Todos os Projetos

```bash
#!/bin/bash
# check-status.sh - Verifica status de todos os repositórios

echo "Verificando status de todos os projetos..."
rgp -path ~/projects -command "status --short" -parallel=false

echo "Verificando repositórios sujos..."
rgp -path ~/projects -command "status --porcelain" | grep -v "^$"
```

### 3. Script para Atualizar Apenas Microserviços

```bash
#!/bin/bash
# update-services.sh - Atualiza apenas repositórios de microserviços

echo "Atualizando microserviços..."
rgp -path ~/workspace -include "*-service,*-api" -exclude "test-*" -command pull -verbose

echo "Microserviços atualizados!"
```

### 4. Script para Backup de Branches Locais

```bash
#!/bin/bash
# backup-branches.sh - Lista todas as branches locais de todos os repos

echo "Fazendo backup da lista de branches..."
rgp -path ~/workspace -command "branch -a" > branches-backup-$(date +%Y%m%d).txt

echo "Backup salvo em branches-backup-$(date +%Y%m%d).txt"
```

## Cenários Comuns

### Desenvolvimento Diário

```bash
# Manhã: verificar o que mudou
rgp -command "fetch --all"
rgp -command "status --short"

# Pull de tudo que está limpo
rgp -command pull -ignore-dirty

# Ver commits recentes
rgp -command "log --oneline -5"
```

### Manutenção de Repositórios

```bash
# Limpar branches mergeadas
rgp -command "branch --merged | grep -v main | xargs -n 1 git branch -d"

# Verificar repositórios que precisam de push  
rgp -command "status --porcelain --branch" | grep "ahead"

# Fazer fetch de todos os remotes
rgp -command "fetch --all --prune"
```

### Trabalho com Features

```bash
# Criar branch em todos os repos
rgp -command "checkout -b feature/nova-funcionalidade"

# Verificar em qual branch cada repo está
rgp -command "branch --show-current"

# Fazer push de todas as branches atuais
rgp -command "push -u origin HEAD"
```

### Análise de Código

```bash
# Contar linhas de código em todos os projetos
rgp -command "ls-files | xargs wc -l"

# Encontrar TODOs em todos os repos
rgp -command "grep -r 'TODO\|FIXME' --include='*.go' --include='*.js' ."

# Ver últimos commits de cada repo
rgp -command "log --oneline -1"
```

## Automatização com Cron

### Backup Automático Diário

```bash
# Adicionar ao crontab: crontab -e
# 0 2 * * * /path/to/daily-backup.sh

#!/bin/bash
# daily-backup.sh
cd ~/workspace
rgp -command "fetch --all" > /dev/null 2>&1
rgp -command "log --oneline -10" > ~/backups/git-log-$(date +%Y%m%d).txt
```

### Notificação de Repositórios Sujos

```bash
#!/bin/bash
# check-dirty.sh - Notifica sobre repositórios com mudanças não commitadas

DIRTY=$(rgp -command "status --porcelain" | grep -v "^$")

if [ ! -z "$DIRTY" ]; then
    echo "AVISO: Repositórios com mudanças não commitadas:"
    echo "$DIRTY"
    # Enviar notificação (ex: via email, Slack, etc.)
fi
```

## Relatórios e Monitoramento

### Relatório de Atividade

```bash
#!/bin/bash
# activity-report.sh - Gera relatório de atividade dos últimos 7 dias

echo "Relatório de Atividade (últimos 7 dias)"
echo "========================================="

rgp -command "log --oneline --since='7 days ago'" | \
  awk '{print $1}' | sort | uniq -c | sort -nr | \
  head -10 > activity-report.txt

echo "Relatório salvo em activity-report.txt"
```

### Status Dashboard

```bash
#!/bin/bash
# dashboard.sh - Dashboard simples de status

clear
echo "Git Dashboard"
echo "============="

echo -e "\nTotal de repositórios:"
rgp -command "status" | grep "repositories processed" | tail -1

echo -e "\nBranches atuais:"
rgp -command "branch --show-current" | sort | uniq -c

echo -e "\nRepositórios sujos:"
DIRTY_COUNT=$(rgp -command "status --porcelain" | grep -c ".")
echo "Total: $DIRTY_COUNT repositórios"

echo -e "\nÚltima atualização: $(date)"
```

## Customização Avançada

### Configuração por Tipo de Projeto

```bash
# Para projetos Go
rgp -include "*-go,*-golang" -command "mod tidy"

# Para projetos Node.js
rgp -include "*-node,*-js,*-ts" -command "status"

# Para projetos Python
rgp -include "*-py,*-python" -command "status"
```

### Integração com Ferramentas de CI/CD

```bash
#!/bin/bash
# ci-check.sh - Verifica se todos os repos estão prontos para CI

echo "Verificando repositórios para CI..."

# Verificar se há mudanças não commitadas
DIRTY=$(rgp -command "status --porcelain" | grep -v "^$")
if [ ! -z "$DIRTY" ]; then
    echo "ERRO: Repositórios com mudanças não commitadas encontrados!"
    exit 1
fi

# Verificar se estão atualizados com remote
rgp -command "fetch --all" > /dev/null 2>&1
BEHIND=$(rgp -command "status --porcelain --branch" | grep "behind")
if [ ! -z "$BEHIND" ]; then
    echo "AVISO: Alguns repositórios estão desatualizados!"
    echo "$BEHIND"
    exit 1
fi

echo "OK: Todos os repositórios estão prontos para CI!"
```

## Dicas de Performance

### Para Muitos Repositórios
```bash
# Aumentar número de workers
rgp -workers 8 -command pull

# Reduzir timeout para comandos rápidos
rgp -timeout 10s -command "branch --show-current"
```

### Para Comandos Lentos
```bash
# Executar sequencialmente
rgp -parallel=false -command "log --all --graph"

# Aumentar timeout
rgp -timeout 60s -command "fsck"
```