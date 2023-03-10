package interactions

import (
	"github.com/andersfylling/disgord"
	"github.com/yoonaowo/discord_verifier/internal/translations"

	"github.com/yoonaowo/discord_verifier/internal/models"
	"github.com/yoonaowo/discord_verifier/internal/models/discord"
	"github.com/yoonaowo/discord_verifier/internal/utils"

	"context"
	"sync"
)

var interactions = &[]*models.Interaction{}
var mapInteractions = make(map[string]*models.Interaction)
var mutex sync.Mutex

func AddInteraction(interaction *models.Interaction) {

	mutex.Lock()
	defer mutex.Unlock()

	*interactions = append(*interactions, interaction)
	mapInteractions[interaction.CommandDefinition.Name] = interaction
}

func GetInteraction(name string) (interactionData *models.Interaction, err error) {

	mutex.Lock()
	defer mutex.Unlock()

	err = utils.ErrInteractionNotFound

	interaction, ok := mapInteractions[name]

	if !ok {
		return
	}

	return interaction, nil
}

var failedResponse = &disgord.CreateInteractionResponse{
	Type: disgord.InteractionCallbackChannelMessageWithSource,
	Data: &disgord.CreateInteractionResponseData{
		Embeds: []*disgord.Embed{
			{
				Title:       translations.Get("ERROR"),
				Description: translations.Get("SOMETHING_WENT_WRONG"),
			},
		},
	},
}

func handle(session disgord.Session, interactionCreate *disgord.InteractionCreate) {
	interaction, err := GetInteraction(interactionCreate.Data.Name)

	if err != nil {
		_ = discordModels.GetClient().SendInteractionResponse(context.Background(), interactionCreate, failedResponse)
		return
	}

	interaction.Callback(session, interactionCreate)
}

func Handle(session disgord.Session, interactionCreate *disgord.InteractionCreate) {
	go handle(session, interactionCreate) // memes
}

func Setup(client *disgord.Client) {
	AddInteraction(verifyStruct)

	GuildID, err := disgord.GetSnowflake(utils.FlagDiscordGuild)

	if err != nil {
		utils.Logger().Fatalln("get snowflake interactions ->", err)
		return
	}

	for _, interaction := range *interactions {
		if err := client.ApplicationCommand(0).Guild(GuildID).Create(interaction.CommandDefinition); err != nil {
			utils.Logger().Fatalf(err.Error())
		}
	}

}
