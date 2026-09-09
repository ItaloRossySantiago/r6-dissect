package dissect

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

func readTime(r *Reader) error {
	time, err := r.Uint32()
	if err != nil {
		return err
	}
	r.time = float64(time)
	r.timeRaw = fmt.Sprintf("%d:%02d", time/60, time%60)
	return nil
}

func readY7Time(r *Reader) error {
	time, err := r.String()
	parts := strings.Split(time, ":")
	if len(parts) == 1 {
		seconds, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return err
		}
		r.time = seconds
		r.timeRaw = parts[0]
		return nil
	}
	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	r.time = float64((minutes * 60) + seconds)
	r.timeRaw = time
	return nil
}

func (r *Reader) roundEnd() {
	log.Debug().Msg("round_end")

	planter := -1
	deaths := make(map[int]int)
	sizes := make(map[int]int)
	roles := make(map[int]TeamRole)

	for _, p := range r.Header.Players {
		sizes[p.TeamIndex] += 1
		roles[p.TeamIndex] = r.Header.Teams[p.TeamIndex].Role
	}

	if r.Header.CodeVersion >= Y9S4 {
		team0Won := r.Header.Teams[0].StartingScore < r.Header.Teams[0].Score
		r.Header.Teams[0].Won = team0Won
		r.Header.Teams[1].Won = !team0Won
	}

	for _, u := range r.MatchFeedback {
		switch u.Type {
		case Kill:
			// FORK NOTE (r6lobby): PlayerIndexByUsername devolve -1 quando o
			// alvo não está em r.Header.Players -- acontece de verdade com
			// jogador de backfill que entra no meio da partida e cujo pacote
			// de info nunca chega neste round. Indexar direto crashava com
			// "index out of range [-1]". Ignorar este evento de kill mantém
			// a morte fora da contagem de "time inteiro eliminado" abaixo,
			// que é o único uso de `deaths` nesta função -- consequência
			// aceitável (round pode não detectar vitória por eliminação
			// neste caso raro) frente a derrubar o processo inteiro.
			idx := r.PlayerIndexByUsername(u.Target)
			if idx < 0 {
				break
			}
			i := r.Header.Players[idx].TeamIndex
			deaths[i] = deaths[i] + 1
			// fix killer username
			if len(u.usernameFromScoreboard) > 0 {
				u.Username = u.usernameFromScoreboard
			}
			break
		case Death:
			idx := r.PlayerIndexByUsername(u.Username)
			if idx < 0 {
				break
			}
			i := r.Header.Players[idx].TeamIndex
			deaths[i] = deaths[i] + 1
			break
		case DefuserPlantComplete:
			planter = r.PlayerIndexByUsername(u.Username)
			break
		case DefuserDisableComplete:
			idx := r.PlayerIndexByUsername(u.Username)
			if idx < 0 {
				break
			}
			i := r.Header.Players[idx].TeamIndex
			r.Header.Teams[i].Won = true
			r.Header.Teams[i].WinCondition = DisabledDefuser
			return
		}
	}

	if planter > -1 {
		r.Header.Teams[r.Header.Players[planter].TeamIndex].Won = true
		r.Header.Teams[r.Header.Players[planter].TeamIndex].WinCondition = DefusedBomb
		return
	}

	// skip for now until we have a more reliable way of determining the win condition
	// Y9S4 at least tells us who won now in the header with StartingScore
	if r.Header.CodeVersion >= Y9S4 {
		return
	}

	if deaths[0] == sizes[0] {
		if planter > -1 && roles[0] == Attack { // ignore attackers killed post-plant
			return
		}
		r.Header.Teams[1].Won = true
		r.Header.Teams[1].WinCondition = KilledOpponents
		return
	}
	if deaths[1] == sizes[1] {
		if planter > -1 && roles[1] == Attack { // ignore attackers killed post-plant
			return
		}
		r.Header.Teams[0].Won = true
		r.Header.Teams[0].WinCondition = KilledOpponents
		return
	}

	i := 0
	if roles[1] == Defense {
		i = 1
	}

	r.Header.Teams[i].Won = true
	r.Header.Teams[i].WinCondition = Time
}
