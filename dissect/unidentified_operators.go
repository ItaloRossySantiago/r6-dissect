package dissect

// UnidentifiedOperators lista IDs de operador vistos em replay real que ainda
// não têm entrada em header.go / operator_roles.go / operator_string.go.
//
// Como identificar um: encontre uma partida em que você sabe de cabeça qual
// operador jogou (o seu, o mais fácil), rode:
//
//	r6lobby-agent scan --json --dir <pasta da partida>
//
// e procure o "profileId" do jogador certo no JSON de saída — o campo
// "operator" ao lado dele é o ID. Depois:
//
//  1. header.go: adicione `NomeDoOperador Operator = <id>` no bloco de
//     constantes Operator.
//  2. operator_string.go: adicione a entrada correspondente em _Operator_map
//     (ou rode `go generate` com o stringer oficial, se disponível).
//  3. operator_roles.go: adicione `<id>: Attack` ou `<id>: Defense` em
//     _operatorRoles, conforme o lado do operador.
//  4. Remova a linha correspondente aqui.
var UnidentifiedOperators = map[Operator]string{
	// Visto em partidas Ranked de 2026-09-07 (r6lobby-agent, Y11S3).
	// Panic original antes do fork: "role unknown for operator ID 444310693746".
	444310693746: "nome ainda não identificado",
}
