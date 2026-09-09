# r6-dissect
[![](https://discordapp.com/api/guilds/936737628756271114/widget.png?style=shield)](https://discord.gg/XdEXWQZZAa)
[![Go Reference](https://pkg.go.dev/badge/github.com/redraskal/r6-dissect.svg)](https://pkg.go.dev/github.com/redraskal/r6-dissect)

Match Replay API/CLI for Rainbow Six: Siege's Dissect (.rec) format.

**This is a work in progress. The data format is subject to change until a stable version is released.**

## Sobre este fork (ItaloRossySantiago/r6-dissect)

Fork de manutenção para o [r6lobby-agent](https://github.com/ItaloRossySantiago/r6lobby-agent).
A `main` upstream está parada desde 16/09/2025 (commit `e6c2ca8`); este fork
parte exatamente desse commit e aplica três mudanças:

**1. `Operator.Role()` não usa mais `panic`.** Todo operador novo que a
Ubisoft lança causa `panic: role unknown for operator ID N` até alguém
reverse-engenheirar o ID e mandar PR — issues [#66](https://github.com/redraskal/r6-dissect/issues/66),
[#103](https://github.com/redraskal/r6-dissect/issues/103) e
[#123](https://github.com/redraskal/r6-dissect/issues/123) são o mesmo
crash, uma vez por temporada. Um processo que lê replay em produção não pode
morrer no meio de uma sessão por causa disso. `Role()` agora devolve
`Unknown` e loga um aviso — ver `dissect/operator_roles.go`.

**2. `x/tools` atualizado de v0.27.0 para a mais recente.** A versão antiga
não é compatível com toolchains Go mais novas (`internal error: package
"strings" without types was imported`) — quebrava `go test` mesmo sem
nenhuma mudança de código.

**3. `Test_operatorsMissing` / `Test_operatorsRedundant` pulados.** Os dois
dependem de `ubi.GetOperatorMap()`, que faz scraping da página de operadores
da Ubisoft. Em set/2026 confirmamos que essa página foi reconstruída em
React — o conteúdo só existe depois do JS rodar num navegador real, e
qualquer cliente HTTP puro (inclusive o `net/http` que a lib usa) recebe só
a casca vazia. Ver comentário em `dissect/test/operators_missing_test.go`.

Também mantemos `dissect/unidentified_operators.go`: um registro dos IDs de
operador vistos em replay real que ainda não têm nome — hoje, um ID visto em
partidas de 2026-09-07 (Y11S3), que por eliminação contra a lista oficial da
Ubisoft só pode ser **Noor** (Defender) ou **Solid Snake** (Attacker) — são
os dois únicos operadores da lista oficial ausentes do `header.go` atual.
Falta confirmar qual dos dois é, e então declarar o ID de verdade.

Download the latest version here: https://github.com/redraskal/r6-dissect/releases

## Current Features
- Match Info (Game version, map, gamemode, match type, teams, players)
- Match Feedback (Kills, headshots, objective locates, defuser plants/disables, BattlEye bans, DCs)
- JSON or Excel output

## Planned Features
- UI alternative
- Track bullet hits/misses
- Track movement packets
- Track other player statistics

### See roadmap at https://github.com/users/redraskal/projects/1.

## CLI Usage
Print a match overview by specifying a match folder or .rec file:
```bash
r6-dissect --info Match-2023-03-13_23-23-58-199
# or
r6-dissect --info Match-2023-03-13_23-23-58-199-R01.rec
```
```
5:20PM INF Version:          Y8S1/7422506
5:20PM INF Recording Player: redraskal [1f63af29-7ebe-48e7-b570-e820632d9565]
5:20PM INF Match ID:         d74d2685-193f-4fee-831f-41f8c7792250
5:20PM INF Timestamp:        2023-03-13 13:00:08 -0500 CDT
5:20PM INF Match Type:       QuickMatch
5:20PM INF Game Mode:        Bomb
5:20PM INF Map:              House
```
You can export round stats to a JSON file:
```bash
r6-dissect Match-2023-03-13_23-23-58-199-R01.rec -o round.json
```
Example:
```json
{
  "gameVersion": "Y8S1",
  "codeVersion": 7422506,
  "timestamp": "2023-03-13T23:25:46Z",
  "matchType": {
    "name": "Ranked",
    "id": 2
  },
  "map": {
    "name": "Villa",
    "id": 88107330328
  },
  "site": "2F Aviator Room, 2F Games Room",
  "recordingPlayerID": 15451868541914624436,
  "recordingProfileID": "1f63af29-7ebe-48e7-b570-e820632d9565",
  "additionalTags": "423855620",
  "gamemode": {
    "name": "Bomb",
    "id": 327933806
  },
...
  "teams": [
    {
      "name": "YOUR TEAM",
      "score": 1,
      "won": true,
      "winCondition": "KilledOpponents",
      "role": "Attack"
    },
    {
      "name": "OPPONENTS",
      "score": 0,
      "won": false,
      "role": "Defense"
    }
  ],
  "players": [
    {
      "id": 1830934665040226621,
      "profileID": "f33396d4-714b-442d-b110-9237e291cc71",
      "username": "IanFiftyForty",
      "teamIndex": 1,
      "operator": {
        "name": "Oryx",
        "id": 104189664155
      },
      "heroName": 243632506966,
      "alliance": 0,
      "roleImage": 104189664090,
      "roleName": "ORYX",
      "rolePortrait": 258649622576
    },
...
  "matchFeedback": [
    {
      "type": "Other",
      "time": "2:59",
      "timeInSeconds": 179,
      "message": "Friendly Fire is now active"
    },
    {
      "type": "Kill",
      "username": "ReithYT",
      "target": "Ambatakum.",
      "headshot": false,
      "time": "1:51",
      "timeInSeconds": 111
    },
...
```
Or the entire match:
```bash
r6-dissect Match-2023-03-13_23-23-58-199 -o match.json
```
Export an Excel spreadsheet by swapping .json with .xlsx.
```bash
r6-dissect Match-2023-03-13_23-23-58-199-R01 -o match.xlsx
```
Output JSON to the console (stdout) with the following syntax:
```bash
# entire match
r6-dissect Match-2023-03-13_23-23-58-199-R01
# or single round
r6-dissect Match-2023-03-13_23-23-58-199-R01/Match-2023-03-13_23-23-58-199-R01.rec
```

See example outputs in [/examples](https://github.com/redraskal/r6-dissect/tree/main/examples).

## Importing a .rec file
```go
package main

import (
	"log"
	"os"

	"github.com/redraskal/r6-dissect/dissect"
)

func main() {
	f, err := os.Open("Match-2023-03-13_23-23-58-199-R01.rec")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r, err := dissect.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	// Use r.ReadPartial() for faster reads with less data (designed to fill in data gaps in the header)
	// dissect.Ok(err) returns true if the error only pertains to EOF (read was successful)
	if err := r.Read(); !dissect.Ok(err) {
		log.Fatal(err)
	}
	print(r.Header.GameVersion) // Y8S1
}
```

#
I would like to thank [stnokott](https://github.com/stnokott) for their work on r6-dissect, along with [draguve](https://github.com/draguve) & other contributors at [draguve/R6-Replays](https://github.com/draguve/R6-Replays) for their additional reverse engineering work.
