package discordnotifier

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/model"
	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/notifier"
	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Token            string
	TargetChannelIDs []string // TODO: 微妙な設計
}

func New(cfg Config) (*discordNotifier, error) {
	sess, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}
	return &discordNotifier{
		sess:             sess,
		targetChannelIDs: cfg.TargetChannelIDs,
	}, nil
}

type discordNotifier struct {
	sess             *discordgo.Session
	targetChannelIDs []string
	logger           *slog.Logger
}

func (d *discordNotifier) Open() error {
	return d.sess.Open()
}

func (d *discordNotifier) Close() {
	d.sess.Close()
}

// NewBook implements notifier.BookNotifier.
func (d *discordNotifier) NewBook(ctx context.Context, books ...*model.Book) error {
	var errs []error
	for _, channelID := range d.targetChannelIDs {
		if err := d.sendBookNotify(ctx, channelID, books...); err != nil {
			d.logger.LogAttrs(context.TODO(), slog.LevelError, "d.sendBookNotify error", slog.String("error", err.Error()), slog.String("channelID", channelID))
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (d *discordNotifier) sendBookNotify(ctx context.Context, channelID string, books ...*model.Book) error {
	if _, err := d.sess.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds: booksToEmbeds(books...),
	}); err != nil {
		return err
	}
	return nil
}

var _ notifier.BookNotifier = (*discordNotifier)(nil)

// convert domain.Book to discordgo.MessageEmbed

func booksToEmbeds(books ...*model.Book) []*discordgo.MessageEmbed {
	embeds := make([]*discordgo.MessageEmbed, 0, len(books))
	for _, book := range books {
		embeds = append(embeds, bookToDiscordEmbed(book))
	}
	return embeds
}

func bookToDiscordEmbed(book *model.Book) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		URL:   book.URL(),
		Type:  discordgo.EmbedTypeArticle,
		Title: book.Title(),
		// Description: book.Description(),
		// Timestamp: time.Now().Format(time.RFC3339),
		// Color:       0,
		// Footer:      &discordgo.MessageEmbedFooter{},
		// Image: &discordgo.MessageEmbedImage{
		// 	URL: book.Cover(),
		// },
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: book.Cover(),
		},
		// Video:     &discordgo.MessageEmbedVideo{},
		Provider: &discordgo.MessageEmbedProvider{
			// URL:  "",
			Name: strings.Join(book.Publishers(), ", "),
		},
		// Author: &discordgo.MessageEmbedAuthor{},
		// Fields:    []*discordgo.MessageEmbedField{},
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "Description",
		Value: book.Description(),
	})

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Publication Date",
		Value:  book.PublishedAt().Format("2006-01-02"),
		Inline: true,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Pages",
		Value:  fmt.Sprintf("%d", book.Page()),
		Inline: true,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Authors",
		Value:  strings.Join(book.Authors(), ", "),
		Inline: true,
	})

	return embed
}
