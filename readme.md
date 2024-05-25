# Rate Limiter Aplicação em Golang

Esta aplicação em Golang implementa um sistema de Rate Limiter, utilizando Redis e memória RAM como mecanismos de armazenamento.

## Pré-requisitos

- Docker
- Go versão 1.22.2

## Variáveis de Ambiente

| Variável | Valor | Descrição |
| --- | --- | --- |
| `TTL` | 30 | Tempo de vida (em segundos) para as requisições |
| `LIMIT` | 5 | Limite de requisições por período de tempo |

## Rodando a Aplicação

### Docker

Para executar a aplicação utilizando Docker, siga os seguintes passos:

1. Inicie o container:
```
docker compose up -d
```

### Local

Para executar a aplicação localmente, siga os seguintes passos:

1. Certifique-se de ter o Go versão 1.22.2 instalado.
2. Execute o comando:
```
go run main.go
```

A aplicação estará escutando na porta `8080`.

## Endpoints

A aplicação possui dois endpoints:

1. `/redis`: Implementa o Rate Limiter utilizando Redis.
2. `/in-memory`: Implementa o Rate Limiter utilizando memória RAM.

Ambos os endpoints são acessíveis via método HTTP GET.

## Testes

Para executar os testes, utilize o seguinte comando:

```
go test ./...
```

## Tecnologias Utilizadas

- Go 1.22.2
- Redis
- Docker