package main

import "github.com/bwmarrin/discordgo"

func handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Message.ChannelID == cfg.Discord.MembershipChannelID {
		// Button clicked
		if i.Type == discordgo.InteractionMessageComponent {
			switch i.MessageComponentData().CustomID {
			case "monthly":
				data := membershipPaymentPrompt("monthly")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseModal, Data: data})
			case "yearly":
				data := membershipPaymentPrompt("yearly")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseModal, Data: data})
			case "activate":
				if getEmail(i.Member.User.ID) == "" {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral, Content: "Please submit your email first!"}})
					return
				}
				if isSubscribed(getEmail(i.Member.User.ID)) {
					if !makeMember(i.Member.User.ID) {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral, Content: "Something went wrong, please create support ticket"}})
						return
					}
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral, Content: "Subscription Activated! \n \n Welcome to the community!"}})
					return
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral, Content: "You do not have an active subscription. \n \n Please purchase a membership first! \n \n If you have already purchased a membership, please open a support ticket"}})
				return

			}

		}
		// Email submitted
		if i.Type == discordgo.InteractionModalSubmit {
			if !idEmailValid(i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral, Content: "Invalid Email, Please try again or open a support ticket"}})
				return
			}
			if i.Member.User.Bot {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral, Content: "Bots can't buy memberships. \n \n If you are a human, please open a support ticket"}})
				return
			}
			if emailExists(i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value, i.Member.User.ID) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral, Content: "Email already in use by another account, Please try again or open a support ticket"}})
				return
			}
			insertToDb(i.Member.User.ID, i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value, i.Member.User.Username, i.Member.User.Discriminator)
			switch i.ModalSubmitData().CustomID {
			case "email-monthly":
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: getRedirectButton("monthly")})

			case "email-yearly":
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: getRedirectButton("yearly")})
			}
		}

	}
}
