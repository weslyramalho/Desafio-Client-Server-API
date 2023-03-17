# Desafio-Client-Server-API

Primeiro desafio feito durante o curso de Go Expert

Neste desafio aplicamos os conhecimentos sobre weberver http, conextos, banco de dados e manipulação de arquivos co Go.

Foi criado dois sistemas em Go:

-Client.go
-Server.go

O cliente.go realiza uma requisição HTTP no server.go, solicitando a cotação do dólar.

O server.go consome a API https://economia.awesomeapi.com.br/json/last/USD-BRL contendo o câmbio do dolar e em seguida retorna no formato JSON o resultado para o cliente.

O server.go registra no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar é de 200ms e o timeout máximo para conseguir persistir os dados  no banco é de 10ms.

O client.go recebe do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go tem um timeout máximo de 300ms para receber o resultado do server.go.

O client.go salva a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}

O endpoint gerado pelo server.go para este desafio é /cotacao e a porta a ser utilizada pelo servidor HTTP é a 8080.


