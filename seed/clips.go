package seed

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ritikmayank/catch-it/repository"
)

type clipDef struct {
	Text     string
	Language string
}

// pool of sentences across languages — 5 clips per day, cycling through the pool
var pool = []clipDef{
	// English
	{"The weather forecast says it will rain tomorrow afternoon.", "en-US"},
	{"She opened the window to let in the fresh morning air.", "en-GB"},
	{"The train arrives at platform three in five minutes.", "en-US"},
	{"He forgot to bring his umbrella on a rainy day.", "en-IN"},
	{"The children played in the park until sunset.", "en-GB"},
	{"Please press the button on the left side of the screen.", "en-US"},
	{"The shop closes at nine o'clock every evening.", "en-IN"},
	{"We need to finish this report before the meeting tomorrow.", "en-GB"},
	{"The new phone model comes in three different colors.", "en-US"},
	{"She has been studying English for the past two years.", "en-IN"},
	// Hindi
	{"आज मौसम बहुत अच्छा है।", "hi-IN"},
	{"कृपया यहाँ अपना नाम लिखें।", "hi-IN"},
	{"बाजार सुबह आठ बजे खुलता है।", "hi-IN"},
	{"उसने कल एक नई किताब खरीदी।", "hi-IN"},
	{"क्या आप मुझे रास्ता बता सकते हैं?", "hi-IN"},
	// Spanish
	{"El tren llega a las tres de la tarde.", "es-ES"},
	{"Por favor, cierre la puerta al salir.", "es-ES"},
	{"Mañana tenemos una reunión importante.", "es-ES"},
	{"El mercado abre a las ocho de la mañana.", "es-ES"},
	{"Necesito encontrar una farmacia cercana.", "es-ES"},
	// French
	{"Le soleil se couche à l'ouest chaque soir.", "fr-FR"},
	{"Pouvez-vous répéter la question, s'il vous plaît?", "fr-FR"},
	{"La réunion commence dans dix minutes.", "fr-FR"},
	{"Il faut tourner à gauche au prochain carrefour.", "fr-FR"},
	{"Le restaurant ferme à vingt-deux heures.", "fr-FR"},
	// More English (varied complexity)
	{"The instructions were written in very small print.", "en-US"},
	{"Can you speak a little more slowly, please?", "en-GB"},
	{"The battery on my phone is running low.", "en-US"},
	{"She smiled when she heard the good news.", "en-IN"},
	{"The bus stop is just around the corner.", "en-GB"},
}

// RunClips seeds 5 clips per day for 30 days starting from startDate.
func RunClips(ctx context.Context, repo *repository.ClipRepository, startDate time.Time) {
	total := 0
	for day := 0; day < 30; day++ {
		date := startDate.AddDate(0, 0, day).Format("2006-01-02")
		for i := 0; i < 5; i++ {
			clip := pool[(day*5+i)%len(pool)]
			if err := repo.Upsert(ctx, date, i+1, clip.Text, clip.Language); err != nil {
				log.Error().Err(err).Str("date", date).Int("index", i+1).Msg("seed clip failed")
			} else {
				total++
			}
		}
	}
	log.Info().Int("count", total).Msg("clips seeded")
}
