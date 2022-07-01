package telegram

import (
	"errors"
	"fmt"
	"log"
	"my/bot/storage"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, userName string) error {

	text = strings.TrimSpace(text)

	log.Printf("got new command %s from %s", text, userName)

	if isAddCmd(text) {
		return p.savePage(chatID, text, userName)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, userName)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMassege(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, userName string) error {
	page := &storage.Page{
		URL:      pageURL,
		UserName: userName,
	}

	isExist, err := p.storage.IsExists(page)
	if err != nil {
		return fmt.Errorf("can not do command: %w", err)
	}
	if isExist {
		return p.tg.SendMassege(chatID, msgAlreadyExists)
	}

	err = p.storage.Save(page)
	if err != nil {
		return fmt.Errorf("can not do command: %w", err)
	}

	err = p.tg.SendMassege(chatID, msgSaved)
	if err != nil {
		return fmt.Errorf("can not do command: %w", err)
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, userName string) error {
	page, err := p.storage.PickRandom(userName)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return fmt.Errorf("can not do command: %w", err)
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMassege(chatID, msgNoSavedPages)
	}

	err = p.tg.SendMassege(chatID, page.URL)
	if err != nil {
		return err
	}
	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMassege(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMassege(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
