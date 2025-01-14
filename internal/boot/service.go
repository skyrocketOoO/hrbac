package boot

import (
	"github.com/rs/zerolog/log"
	"github.com/skyrocketOoO/serverx/internal/service/exter/db"
	"github.com/skyrocketOoO/serverx/internal/service/inter/validator"
)

func NewService() error {
	log.Info().Msg("InitService")
	if err := db.New(); err != nil {
		return err
	}

	validator.New()

	return nil
}
