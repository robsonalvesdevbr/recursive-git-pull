# CLAUDE

## Overview

O objetivo deste projeto é criar uma ferramenta que facilite a execução de comandos Git em múltiplos repositórios de forma recursiva.

## Funcionalidades

- Execução de comandos Git em todos os repositórios dentro de um diretório especificado.
  Neste exemplo, todos os repositórios dentro da pasta "workspace" serão afetados.
  Exemplo:
  - workspace/
    - repo1/
      - .git/
    - repo2/
      - .git/
    - repo3/
      - .git/
- Suporte a comandos Git personalizados.
- Relatórios de status da execução dos comandos.
- Opção para ignorar repositórios que não estejam em um estado "limpo" (sem mudanças não commitadas).
- Opção para executar comandos em paralelo para melhorar a performance.
- Opção para incluir ou excluir repositórios com base em padrões de nome.
- Opção para definir um tempo limite para a execução dos comandos.
- Logs detalhados da execução dos comandos.
- Executar o pull em todos os branches de todos os repositórios.
