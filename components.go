package main

import "github.com/bwmarrin/discordgo"

func getMembershipPrompt() *discordgo.MessageSend {
	components := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Monthly",
				CustomID: "monthly",
				Style:    discordgo.SuccessButton,
				Emoji: discordgo.ComponentEmoji{
					Name: "💳",
				},
			},
			discordgo.Button{
				Label:    "Yearly",
				CustomID: "yearly",
				Style:    discordgo.SuccessButton,
				Emoji: discordgo.ComponentEmoji{
					Name: "💳",
				},
			},
			discordgo.Button{
				Label:    "Activate",
				CustomID: "activate",
				Style:    discordgo.SecondaryButton,
				Emoji: discordgo.ComponentEmoji{
					Name: "🔓",
				},
			},
		},
	}
	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Membership",
		Description: "Purchase a montly or yearly membership to support the server and get access to exclusive channels and roles.",
		URL:         "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		Color:       0x00FFFF,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Monthly",
				Value:  "$" + cfg.Variables.MonthlyPrice,
				Inline: false,
			},
			{
				Name:   "Yearly",
				Value:  "$" + cfg.Variables.YearlyPrice,
				Inline: false,
			},
		},
	}
	return &discordgo.MessageSend{Embeds: []*discordgo.MessageEmbed{embed}, Components: []discordgo.MessageComponent{components}}
}

func membershipPaymentPrompt(t string) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Components: []discordgo.MessageComponent{discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "email-" + t,
					Placeholder: "Enter your email",
					MaxLength:   100,
					MinLength:   5,
					Label:       "Email",
					Style:       discordgo.TextInputShort,
					Required:    true,
				},
			},
		},
		},
		CustomID: "email-" + t,
		Title:    "Email for Payment confirmation",
		Flags:    discordgo.MessageFlagsEphemeral,
	}
}

func getRedirectButton(s string) *discordgo.InteractionResponseData {
	switch s {
	case "monthly":
		s = "https://buy.stripe.com/fZedUt8qL2JUg80bIK"
	case "yearly":
		s = "https://buy.stripe.com/bIY3fPayT84ecVO8wx"
	}
	return &discordgo.InteractionResponseData{
		Flags:   discordgo.MessageFlagsEphemeral,
		Content: " ⚠ Please use the SAME EMAIL for contact information when paying !!! \n \n Click here to pay", Components: []discordgo.MessageComponent{discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "Pay",
					Style: discordgo.LinkButton,
					URL:   s,
					Emoji: discordgo.ComponentEmoji{
						Name: "💳",
					},
				},
			},
		},
		},
	}
}
