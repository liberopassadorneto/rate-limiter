# Rate Limiter

## Descrição

Este projeto implementa um Rate Limiter em Go que limita o número de requisições por segundo com base no endereço IP ou em um token de acesso. Utiliza Redis para armazenar as informações de limitação e pode ser configurado via variáveis de ambiente ou arquivo `.env`.

## Funcionalidades

- **Limitação por IP**: Controla o número de requisições por endereço IP.
- **Limitação por Token**: Controla o número de requisições por token de acesso. As configurações de token sobrescrevem as de IP.
- **Middleware**: Pode ser integrado facilmente como middleware em servidores web.
- **Configuração**: Via variáveis de ambiente ou arquivo `.env`.
- **Persistência**: Utiliza Redis para armazenar informações de limitação.
- **Estratégia**: Implementa uma interface que permite trocar facilmente o mecanismo de persistência.

## Configuração

1. **Clone o Repositório**

   ```bash
   git clone https://github.com/seu-usuario/rate-limiter.git
   cd rate-limiter
# rate-limiter
